const fs = require('fs');
const path = require('path');
const os = require('os');
const { execSync } = require('child_process');
const { getRepoRoot, getWorkspaceRoot, getMemoryDir, getSessionScope } = require('./gep/paths');
const { extractSignals } = require('./gep/signals');
const {
  loadGenes,
  loadCapsules,
  readAllEvents,
  getLastEventId,
  readRecentFailedCapsules,
  ensureAssetFiles,
} = require('./gep/assetStore');
const { selectGeneAndCapsule } = require('./gep/selector');
const { buildGepPrompt, buildReusePrompt, buildHubMatchedBlock } = require('./gep/prompt');
const { hubSearch } = require('./gep/hubSearch');
const { logAssetCall } = require('./gep/assetCallLog');
const { buildCandidatePreviews } = require('./gep/candidateEval');
const memoryAdapter = require('./gep/memoryGraphAdapter');
const {
  getAdvice: getMemoryAdvice,
  recordSignalSnapshot,
  recordHypothesis,
  recordAttempt,
  recordOutcome: recordOutcomeFromState,
  memoryGraphPath,
} = memoryAdapter;
const { readStateForSolidify, writeStateForSolidify } = require('./gep/solidify');
const { fetchTasks, selectBestTask, claimTask, taskToSignals, claimWorkerTask, estimateCommitmentDeadline } = require('./gep/taskReceiver');
const { generateQuestions } = require('./gep/questionGenerator');
const { buildMutation, isHighRiskMutationAllowed } = require('./gep/mutation');
const { selectPersonalityForRun } = require('./gep/personality');
const { clip, writePromptArtifact, renderSessionsSpawnCall } = require('./gep/bridge');
const { getEvolutionDir } = require('./gep/paths');
const { shouldReflect, buildReflectionContext, recordReflection } = require('./gep/reflection');
const { loadNarrativeSummary } = require('./gep/narrativeMemory');
const { maybeReportIssue } = require('./gep/issueReporter');
const { resolveStrategy } = require('./gep/strategy');
const { expandSignals } = require('./gep/learningSignals');

const REPO_ROOT = getRepoRoot();

// 详细日志辅助函数。检查 EVOLVER_VERBOSE 环境常量（由 index.js 中的 --verbose 标志设置）。
function verbose() {
  if (String(process.env.EVOLVER_VERBOSE || '').toLowerCase() !== 'true') return;
  const args = Array.prototype.slice.call(arguments);
  args.unshift('[Verbose]');
  console.log.apply(console, args);
}

// 空闲周期节流：跟踪上次 Hub 获取时间，避免在饱和时进行冗余 API 调用。
// 当 evolver 饱和时（无可操作信号），Hub 调用被限制为最多
// 每 EVOLVER_IDLE_FETCH_INTERVAL_MS（默认 30 分钟）一次，而不是每个周期都调用。
let _lastHubFetchMs = 0;

function shouldSkipHubCalls(signals) {
  if (!Array.isArray(signals)) return false;
  const saturationIndicators = ['force_steady_state', 'evolution_saturation', 'empty_cycle_loop_detected'];
  let hasSaturation = false;
  for (let si = 0; si < saturationIndicators.length; si++) {
    if (signals.indexOf(saturationIndicators[si]) !== -1) { hasSaturation = true; break; }
  }
  if (!hasSaturation) return false;

  const actionablePatterns = [
    'log_error', 'recurring_error', 'capability_gap', 'perf_bottleneck',
    'external_task', 'bounty_task', 'overdue_task', 'urgent',
    'unsupported_input_type',
  ];
  for (let ai = 0; ai < signals.length; ai++) {
    const s = signals[ai];
    if (actionablePatterns.indexOf(s) !== -1) return false;
    if (s.indexOf('errsig:') === 0) return false;
    if (s.indexOf('user_feature_request:') === 0 && s.length > 21) return false;
    if (s.indexOf('user_improvement_suggestion:') === 0 && s.length > 28) return false;
  }
  return true;
}

// 从仓库根目录加载环境变量
try {
  require('dotenv').config({ path: path.join(REPO_ROOT, '.env'), quiet: true });
} catch (e) {
  // dotenv 可能未安装或 .env 缺失，优雅降级
}

// CLI 标志或环境变量配置
const ARGS = process.argv.slice(2);
const IS_REVIEW_MODE = ARGS.includes('--review');
const IS_DRY_RUN = ARGS.includes('--dry-run');
const IS_RANDOM_DRIFT = ARGS.includes('--drift') || String(process.env.RANDOM_DRIFT || '').toLowerCase() === 'true';

// 默认配置
const MEMORY_DIR = getMemoryDir();
const AGENT_NAME = process.env.AGENT_NAME || 'main';
const AGENT_SESSIONS_DIR = path.join(os.homedir(), `.openclaw/agents/${AGENT_NAME}/sessions`);
const CURSOR_TRANSCRIPTS_DIR = process.env.EVOLVER_CURSOR_TRANSCRIPTS_DIR || '';
const TODAY_LOG = path.join(MEMORY_DIR, new Date().toISOString().split('T')[0] + '.md');

// 确保内存目录存在，以便状态/缓存写入正常工作。
try {
  if (!fs.existsSync(MEMORY_DIR)) fs.mkdirSync(MEMORY_DIR, { recursive: true });
} catch (e) {
  console.warn('[Evolver] 创建 MEMORY_DIR 失败（可能导致下游错误）:', e && e.message || e);
}

function formatSessionLog(jsonlContent) {
  const result = [];
  const lines = jsonlContent.split('\n');
  let lastLine = '';
  let repeatCount = 0;

  const flushRepeats = () => {
    if (repeatCount > 0) {
      result.push(`   ... [重复 ${repeatCount} 次] ...`);
      repeatCount = 0;
    }
  };

  for (const line of lines) {
    if (!line.trim()) continue;
    try {
      const data = JSON.parse(line);
      let entry = '';

      if (data.type === 'message' && data.message) {
        const role = (data.message.role || 'unknown').toUpperCase();
        let content = '';
        if (Array.isArray(data.message.content)) {
          content = data.message.content
            .map(c => {
              if (c.type === 'text') return c.text;
              if (c.type === 'toolCall') return `[TOOL: ${c.name}]`;
              return '';
            })
            .join(' ');
        } else if (typeof data.message.content === 'string') {
          content = data.message.content;
        } else {
          content = JSON.stringify(data.message.content);
        }

        // 捕获 errorMessage 字段中的 LLM 错误（如 "Unsupported MIME type: image/gif"）
        if (data.message.errorMessage) {
          const errMsg = typeof data.message.errorMessage === 'string'
            ? data.message.errorMessage
            : JSON.stringify(data.message.errorMessage);
          content = `[LLM 错误] ${errMsg.replace(/\n+/g, ' ').slice(0, 300)}`;
        }

        // 过滤：跳过心跳以减少噪音
        if (content.trim() === 'HEARTBEAT_OK') continue;
        if (content.includes('NO_REPLY') && !data.message.errorMessage) continue;

        // 清理换行符以紧凑阅读
        content = content.replace(/\n+/g, ' ').slice(0, 300);
        entry = `**${role}**: ${content}`;
      } else if (data.type === 'tool_result' || (data.message && data.message.role === 'toolResult')) {
        // 过滤：跳过通用成功结果或简短无信息的结果
        // 只显示错误或有意义的输出
        let resContent = '';

        // 健壮提取：处理缺少 'output' 的结构化工具结果（如 sessions_spawn）
        if (data.tool_result) {
          if (data.tool_result.output) {
            resContent = data.tool_result.output;
          } else {
            resContent = JSON.stringify(data.tool_result);
          }
        }

        if (data.content) resContent = typeof data.content === 'string' ? data.content : JSON.stringify(data.content);

        if (resContent.length < 50 && (resContent.includes('success') || resContent.includes('done'))) continue;
        if (resContent.trim() === '' || resContent === '{}') continue;

        // 改进：显示结果的片段（尤其是错误）而不是隐藏它
        const preview = resContent.replace(/\n+/g, ' ').slice(0, 200);
        entry = `[工具结果] ${preview}${resContent.length > 200 ? '...' : ''}`;
      }

      if (entry) {
        if (entry === lastLine) {
          repeatCount++;
        } else {
          flushRepeats();
          result.push(entry);
          lastLine = entry;
        }
      }
    } catch (e) {
      continue;
    }
  }
  flushRepeats();
  return result.join('\n');
}

function formatCursorTranscript(raw) {
  const lines = raw.split('\n');
  const result = [];
  let skipUntilNextBlock = false;

  for (let i = 0; i < lines.length; i++) {
    const line = lines[i];
    const trimmed = line.trim();

    // 保留用户消息和助手文本响应
    if (trimmed === 'user:' || trimmed.startsWith('A:')) {
      skipUntilNextBlock = false;
      result.push(trimmed);
      continue;
    }

    // 工具调用行：保留为紧凑标记，跳过其参数块
    if (trimmed.startsWith('[Tool call]')) {
      skipUntilNextBlock = true;
      result.push(`[Tool call] ${trimmed.replace('[Tool call]', '').trim()}`);
      continue;
    }

    // 工具结果标记：跳过其内容（通常很大且嘈杂）
    if (trimmed.startsWith('[Tool result]')) {
      skipUntilNextBlock = true;
      continue;
    }

    if (skipUntilNextBlock) continue;

    // 保留用户查询内容和助手文本（跳过 <user_query> 等 XML 标签）
    if (trimmed.startsWith('<') && trimmed.endsWith('>')) continue;
    if (trimmed) {
      result.push(trimmed.slice(0, 300));
    }
  }

  return result.join('\n');
}

function readCursorTranscripts() {
  if (!CURSOR_TRANSCRIPTS_DIR) return '';
  try {
    if (!fs.existsSync(CURSOR_TRANSCRIPTS_DIR)) return '';

    const now = Date.now();
    const ACTIVE_WINDOW_MS = 24 * 60 * 60 * 1000;
    const TARGET_BYTES = 120000;
    const PER_FILE_BYTES = 20000;
    const RECENCY_GUARD_MS = 30 * 1000;

    let files = fs
      .readdirSync(CURSOR_TRANSCRIPTS_DIR)
      .filter(f => f.endsWith('.txt') || f.endsWith('.jsonl'))
      .map(f => {
        try {
          const st = fs.statSync(path.join(CURSOR_TRANSCRIPTS_DIR, f));
          return { name: f, time: st.mtime.getTime(), size: st.size };
        } catch (e) {
          return null;
        }
      })
      .filter(f => f && (now - f.time) < ACTIVE_WINDOW_MS)
      .sort((a, b) => b.time - a.time);

    if (files.length === 0) return '';

    // 跳过最近修改的文件（如果最近 30 秒内被修改）——
    // 它很可能就是触发此 evolver 运行的当前活动会话，
    // 读取它会导致自引用信号噪音。
    if (files.length > 1 && (now - files[0].time) < RECENCY_GUARD_MS) {
      files = files.slice(1);
    }

    const maxFiles = Math.min(files.length, 6);
    const sections = [];
    let totalBytes = 0;

    for (let i = 0; i < maxFiles && totalBytes < TARGET_BYTES; i++) {
      const f = files[i];
      const bytesLeft = TARGET_BYTES - totalBytes;
      const readSize = Math.min(PER_FILE_BYTES, bytesLeft);
      const raw = readRecentLog(path.join(CURSOR_TRANSCRIPTS_DIR, f.name), readSize);
      if (raw.trim() && !raw.startsWith('[缺失]')) {
        const formatted = formatCursorTranscript(raw);
        if (formatted.trim()) {
          sections.push(`--- CURSOR SESSION (${f.name}) ---\n${formatted}`);
          totalBytes += formatted.length;
        }
      }
    }

    return sections.join('\n\n');
  } catch (e) {
    console.warn(`[CursorTranscripts] 读取失败: ${e.message}`);
    return '';
  }
}

function readRealSessionLog() {
  try {
    // 主数据源：OpenClaw 会话日志 (.jsonl)
    if (fs.existsSync(AGENT_SESSIONS_DIR)) {
      const now = Date.now();
      const ACTIVE_WINDOW_MS = 24 * 60 * 60 * 1000; // 24 小时
      const TARGET_BYTES = 120000;
      const PER_SESSION_BYTES = 20000;

      const sessionScope = getSessionScope();

      let files = fs
        .readdirSync(AGENT_SESSIONS_DIR)
        .filter(f => f.endsWith('.jsonl') && !f.includes('.lock'))
        .map(f => {
          try {
            const st = fs.statSync(path.join(AGENT_SESSIONS_DIR, f));
            return { name: f, time: st.mtime.getTime(), size: st.size };
          } catch (e) {
            return null;
          }
        })
        .filter(f => f && (now - f.time) < ACTIVE_WINDOW_MS)
        .sort((a, b) => b.time - a.time);

      if (files.length > 0) {
        let nonEvolverFiles = files.filter(f => !f.name.startsWith('evolver_hand_'));

        if (sessionScope && nonEvolverFiles.length > 0) {
          const scopeLower = sessionScope.toLowerCase();
          const scopedFiles = nonEvolverFiles.filter(f => f.name.toLowerCase().includes(scopeLower));
          if (scopedFiles.length > 0) {
            nonEvolverFiles = scopedFiles;
            console.log(`[SessionScope] 已筛选至 ${scopedFiles.length} 个匹配作用域 "${sessionScope}" 的会话。`);
          } else {
            console.log(`[SessionScope] 没有会话匹配作用域 "${sessionScope}"。使用全部 ${nonEvolverFiles.length} 个会话（回退）。`);
          }
        }

        const activeFiles = nonEvolverFiles.length > 0 ? nonEvolverFiles : files.slice(0, 1);

        const maxSessions = Math.min(activeFiles.length, 6);
        const sections = [];
        let totalBytes = 0;

        for (let i = 0; i < maxSessions && totalBytes < TARGET_BYTES; i++) {
          const f = activeFiles[i];
          const bytesLeft = TARGET_BYTES - totalBytes;
          const readSize = Math.min(PER_SESSION_BYTES, bytesLeft);
          const raw = readRecentLog(path.join(AGENT_SESSIONS_DIR, f.name), readSize);
          const formatted = formatSessionLog(raw);
          if (formatted.trim()) {
            sections.push(`--- SESSION (${f.name}) ---\n${formatted}`);
            totalBytes += formatted.length;
          }
        }

        if (sections.length > 0) {
          return sections.join('\n\n');
        }
      }
    }

    // 回退：Cursor agent-transcripts (.txt)
    const cursorContent = readCursorTranscripts();
    if (cursorContent) {
      console.log('[SessionFallback] 使用 Cursor agent-transcripts 作为会话数据源。');
      return cursorContent;
    }

    return '[未找到会话日志]';
  } catch (e) {
    return `[读取会话日志出错: ${e.message}]`;
  }
}

function readRecentLog(filePath, size = 10000) {
  try {
    if (!fs.existsSync(filePath)) return `[缺失] ${filePath}`;
    const stats = fs.statSync(filePath);
    const start = Math.max(0, stats.size - size);
    const buffer = Buffer.alloc(stats.size - start);
    const fd = fs.openSync(filePath, 'r');
    fs.readSync(fd, buffer, 0, buffer.length, start);
    fs.closeSync(fd);
    return buffer.toString('utf8');
  } catch (e) {
    return `[读取 ${filePath} 出错: ${e.message}]`;
  }
}

function computeAdaptiveStrategyPolicy(opts) {
  const recentEvents = Array.isArray(opts && opts.recentEvents) ? opts.recentEvents : [];
  const selectedGene = opts && opts.selectedGene ? opts.selectedGene : null;
  const signals = Array.isArray(opts && opts.signals) ? opts.signals : [];
  const baseStrategy = resolveStrategy({ signals: signals });

  const tail = recentEvents.slice(-8);
  let repairStreak = 0;
  for (let i = tail.length - 1; i >= 0; i--) {
    if (tail[i] && tail[i].intent === 'repair') repairStreak++;
    else break;
  }
  let failureStreak = 0;
  for (let i = tail.length - 1; i >= 0; i--) {
    if (tail[i] && tail[i].outcome && tail[i].outcome.status === 'failed') failureStreak++;
    else break;
  }

  const antiPatterns = selectedGene && Array.isArray(selectedGene.anti_patterns) ? selectedGene.anti_patterns.slice(-5) : [];
  const learningHistory = selectedGene && Array.isArray(selectedGene.learning_history) ? selectedGene.learning_history.slice(-6) : [];
  const signalTags = new Set(expandSignals(signals, ''));
  const overlappingAntiPatterns = antiPatterns.filter(function (ap) {
    return ap && Array.isArray(ap.learning_signals) && ap.learning_signals.some(function (tag) {
      return signalTags.has(String(tag));
    });
  });
  const hardFailures = overlappingAntiPatterns.filter(function (ap) { return ap && ap.mode === 'hard'; }).length;
  const softFailures = overlappingAntiPatterns.filter(function (ap) { return ap && ap.mode !== 'hard'; }).length;
  const recentSuccesses = learningHistory.filter(function (x) { return x && x.outcome === 'success'; }).length;

  const stagnation = signals.includes('stable_success_plateau') ||
    signals.includes('evolution_saturation') ||
    signals.includes('empty_cycle_loop_detected') ||
    failureStreak >= 3 ||
    repairStreak >= 3;

  const forceInnovate = stagnation && !signals.includes('log_error');
  const highRiskGene = hardFailures >= 1 || (softFailures >= 2 && recentSuccesses === 0);
  const cautiousExecution = highRiskGene || failureStreak >= 2;

  let blastRadiusMaxFiles = selectedGene && selectedGene.constraints && Number.isFinite(Number(selectedGene.constraints.max_files))
    ? Number(selectedGene.constraints.max_files)
    : 12;
  if (cautiousExecution) blastRadiusMaxFiles = Math.max(2, Math.min(blastRadiusMaxFiles, 6));
  else if (forceInnovate) blastRadiusMaxFiles = Math.max(3, Math.min(blastRadiusMaxFiles, 10));

  const directives = [];
  directives.push('Base strategy: ' + baseStrategy.label + ' (' + baseStrategy.description + ')');
  if (forceInnovate) directives.push('Force strategy shift: prefer innovate over repeating repair/optimize.');
  if (highRiskGene) directives.push('Selected gene is high risk for current signals; keep blast radius narrow and prefer smallest viable change.');
  if (failureStreak >= 2) directives.push('检测到近期连续失败；避免重复近期失败的方法。');
  directives.push('Target max files for this cycle: ' + blastRadiusMaxFiles + '.');

  return {
    name: baseStrategy.name,
    label: baseStrategy.label,
    description: baseStrategy.description,
    forceInnovate: forceInnovate,
    cautiousExecution: cautiousExecution,
    highRiskGene: highRiskGene,
    repairStreak: repairStreak,
    failureStreak: failureStreak,
    blastRadiusMaxFiles: blastRadiusMaxFiles,
    directives: directives,
  };
}

function checkSystemHealth() {
  const report = [];
  try {
    // 运行时间与 Node 版本
    const uptime = (os.uptime() / 3600).toFixed(1);
    report.push(`Uptime: ${uptime}h`);
    report.push(`Node: ${process.version}`);

    // 内存使用 (RSS)
    const mem = process.memoryUsage();
    const rssMb = (mem.rss / 1024 / 1024).toFixed(1);
    report.push(`Agent RSS: ${rssMb}MB`);

    // 优化：使用 Node.js 原生 fs.statfsSync 替代 spawn 'df'
    if (fs.statfsSync) {
      const stats = fs.statfsSync('/');
      const total = stats.blocks * stats.bsize;
      const free = stats.bfree * stats.bsize;
      const used = total - free;
      const freeGb = (free / 1024 / 1024 / 1024).toFixed(1);
      const usedPercent = Math.round((used / total) * 100);
      report.push(`Disk: ${usedPercent}% (${freeGb}G free)`);
    }
  } catch (e) {}

  try {
    if (process.platform === 'win32') {
      const wmic = execSync('tasklist /FI "IMAGENAME eq node.exe" /NH', {
        encoding: 'utf8',
        stdio: ['ignore', 'pipe', 'ignore'],
        timeout: 3000,
        windowsHide: true,
      });
      const count = wmic.split('\n').filter(l => l.trim() && !l.includes('INFO:')).length;
      report.push(`Node Processes: ${count}`);
    } else {
      try {
        const pgrep = execSync('pgrep -c node', {
          encoding: 'utf8',
          stdio: ['ignore', 'pipe', 'ignore'],
          timeout: 2000,
        });
        report.push(`Node Processes: ${pgrep.trim()}`);
      } catch (e) {
        const ps = execSync('ps aux | grep node | grep -v grep | wc -l', {
          encoding: 'utf8',
          stdio: ['ignore', 'pipe', 'ignore'],
          timeout: 2000,
        });
        report.push(`Node Processes: ${ps.trim()}`);
      }
    }
  } catch (e) {}

  // 集成健康检查（环境变量）
  try {
    const issues = [];

    // 通用集成状态检查（解耦）
    if (process.env.INTEGRATION_STATUS_CMD) {
      try {
        const status = execSync(process.env.INTEGRATION_STATUS_CMD, {
          encoding: 'utf8',
          stdio: ['ignore', 'pipe', 'ignore'],
          timeout: 2000,
          windowsHide: true,
        });
        if (status.trim()) issues.push(status.trim());
      } catch (e) {}
    }

    if (issues.length > 0) {
      report.push(`Integrations: ${issues.join(', ')}`);
    } else {
      report.push('集成状态: 正常');
    }
  } catch (e) {}

  return report.length ? report.join(' | ') : '健康检查不可用';
}

function getMutationDirective(logContent) {
  // 从近期日志推导出的信号提示。
  const errorMatches = logContent.match(/\[ERROR|Error:|Exception:|FAIL|Failed|"isError":true/gi) || [];
  const errorCount = errorMatches.length;
  const isUnstable = errorCount > 2;
  const recommendedIntent = isUnstable ? 'repair' : 'optimize';

  return `
[信号提示]
- 近期错误数: ${errorCount}
- 稳定性: ${isUnstable ? '不稳定' : '稳定'}
- 建议意图: ${recommendedIntent}
`;
}

const STATE_FILE = path.join(getEvolutionDir(), 'evolution_state.json');
const DORMANT_HYPOTHESIS_FILE = path.join(getEvolutionDir(), 'dormant_hypothesis.json');
const DORMANT_TTL_MS = 3600 * 1000;

function writeDormantHypothesis(data) {
  try {
    const dir = getEvolutionDir();
    if (!fs.existsSync(dir)) fs.mkdirSync(dir, { recursive: true });
    const obj = Object.assign({}, data, { created_at: new Date().toISOString(), ttl_ms: DORMANT_TTL_MS });
    const tmp = DORMANT_HYPOTHESIS_FILE + '.tmp';
    fs.writeFileSync(tmp, JSON.stringify(obj, null, 2) + '\n', 'utf8');
    fs.renameSync(tmp, DORMANT_HYPOTHESIS_FILE);
    console.log('[休眠假设] 已在退避前保存部分状态: ' + (data.backoff_reason || '未知'));
  } catch (e) {
    console.log('[休眠假设] 写入失败（非致命）: ' + (e && e.message ? e.message : e));
  }
}

function readDormantHypothesis() {
  try {
    if (!fs.existsSync(DORMANT_HYPOTHESIS_FILE)) return null;
    const raw = fs.readFileSync(DORMANT_HYPOTHESIS_FILE, 'utf8');
    if (!raw.trim()) return null;
    const obj = JSON.parse(raw);
    const createdAt = obj.created_at ? new Date(obj.created_at).getTime() : 0;
    const ttl = Number.isFinite(Number(obj.ttl_ms)) ? Number(obj.ttl_ms) : DORMANT_TTL_MS;
    if (Date.now() - createdAt > ttl) {
      clearDormantHypothesis();
      console.log('[休眠假设] 已过期（存在时间: ' + Math.round((Date.now() - createdAt) / 1000) + '秒）。已丢弃。');
      return null;
    }
    return obj;
  } catch (e) {
    return null;
  }
}

function clearDormantHypothesis() {
  try {
    if (fs.existsSync(DORMANT_HYPOTHESIS_FILE)) fs.unlinkSync(DORMANT_HYPOTHESIS_FILE);
  } catch (e) {}
}
// 从 WORKSPACE 根目录读取 MEMORY.md 和 USER.md（不是 evolver 插件目录）。
// 这避免了目标文件被临时删除时的符号链接断裂。
const WORKSPACE_ROOT = getWorkspaceRoot();
const ROOT_MEMORY = path.join(WORKSPACE_ROOT, 'MEMORY.md');
const DIR_MEMORY = path.join(MEMORY_DIR, 'MEMORY.md');
const MEMORY_FILE = fs.existsSync(ROOT_MEMORY) ? ROOT_MEMORY : (fs.existsSync(DIR_MEMORY) ? DIR_MEMORY : ROOT_MEMORY);
const USER_FILE = path.join(WORKSPACE_ROOT, 'USER.md');

function readMemorySnippet() {
  try {
    // 会话作用域隔离：当作用域处于活动状态时，优先使用
    // memory/scopes/<scope>/MEMORY.md 中的 MEMORY.md。
    // 如果作用域文件不存在则回退到全局 MEMORY.md
    // （常见：作用域 MEMORY.md 在首次进化时创建）。
    const scope = getSessionScope();
    let memFile = MEMORY_FILE;
    if (scope) {
      const scopedMemory = path.join(MEMORY_DIR, 'scopes', scope, 'MEMORY.md');
      if (fs.existsSync(scopedMemory)) {
        memFile = scopedMemory;
        console.log(`[SessionScope] 正在读取作用域 "${scope}" 的 MEMORY.md。`);
      } else {
        // 首次使用作用域运行：将使用全局 MEMORY.md，但记录一下。
        console.log(`[SessionScope] 作用域 "${scope}" 无 MEMORY.md。使用全局 MEMORY.md。`);
      }
    }
    if (!fs.existsSync(memFile)) return '[MEMORY.md 缺失]';
    const content = fs.readFileSync(memFile, 'utf8');
    // 优化：为现代上下文窗口将限制从 2000 提高到 50000
    return content.length > 50000
      ? content.slice(0, 50000) + `\n... [已截断: 剩余 ${content.length - 50000} 字符]`
      : content;
  } catch (e) {
    return '[读取 MEMORY.md 出错]';
  }
}

function readUserSnippet() {
  try {
    if (!fs.existsSync(USER_FILE)) return '[USER.md 缺失]';
    return fs.readFileSync(USER_FILE, 'utf8');
  } catch (e) {
    return '[读取 USER.md 出错]';
  }
}

function getNextCycleId() {
  let state = { cycleCount: 0, lastRun: 0 };
  try {
    if (fs.existsSync(STATE_FILE)) {
      state = JSON.parse(fs.readFileSync(STATE_FILE, 'utf8'));
    }
  } catch (e) {
    console.warn('[Evolve] 读取状态文件失败:', e && e.message || e);
  }

  state.cycleCount = (state.cycleCount || 0) + 1;
  state.lastRun = Date.now();

  try {
    fs.writeFileSync(STATE_FILE, JSON.stringify(state, null, 2));
  } catch (e) {
    console.warn('[Evolve] 写入状态文件失败:', e && e.message || e);
  }

  return String(state.cycleCount).padStart(4, '0');
}

function performMaintenance() {
  // 自动更新检查（限速，非致命）。
  checkAndAutoUpdate();

  try {
    if (!fs.existsSync(AGENT_SESSIONS_DIR)) return;

    const files = fs.readdirSync(AGENT_SESSIONS_DIR).filter(f => f.endsWith('.jsonl'));

    // 清理 evolver 自身的手动会话文件。
    // 这些是一次性执行器会话，不能累积，
    // 否则会污染 agent 上下文并挤占用户对话。
    const evolverFiles = files.filter(f => f.startsWith('evolver_hand_'));
    for (const f of evolverFiles) {
      try {
        fs.unlinkSync(path.join(AGENT_SESSIONS_DIR, f));
      } catch (_) {}
    }
    if (evolverFiles.length > 0) {
      console.log(`[Maintenance] 已清理 ${evolverFiles.length} 个 evolver 手动会话。`);
    }

    // 当会话数量超过阈值时，归档旧的非 evolver 会话。
    const remaining = files.length - evolverFiles.length;
    if (remaining < 100) return;

    console.log(`[Maintenance] 发现 ${remaining} 个会话日志。正在归档旧日志...`);

    const ARCHIVE_DIR = path.join(AGENT_SESSIONS_DIR, 'archive');
    if (!fs.existsSync(ARCHIVE_DIR)) fs.mkdirSync(ARCHIVE_DIR, { recursive: true });

    const fileStats = files
      .filter(f => !f.startsWith('evolver_hand_'))
      .map(f => {
        try {
          return { name: f, time: fs.statSync(path.join(AGENT_SESSIONS_DIR, f)).mtime.getTime() };
        } catch (e) {
          return null;
        }
      })
      .filter(Boolean)
      .sort((a, b) => a.time - b.time);

    const toArchive = fileStats.slice(0, fileStats.length - 50);

    for (const file of toArchive) {
      const oldPath = path.join(AGENT_SESSIONS_DIR, file.name);
      const newPath = path.join(ARCHIVE_DIR, file.name);
      fs.renameSync(oldPath, newPath);
    }
    if (toArchive.length > 0) {
      console.log(`[Maintenance] 已归档 ${toArchive.length} 个日志到 ${ARCHIVE_DIR}`);
    }
  } catch (e) {
    console.error(`[Maintenance] 错误: ${e.message}`);
  }
}

// --- 自动更新：检查 ClawHub 上 evolver 和 wrapper 的更新版本 ---
function checkAndAutoUpdate() {
  try {
    // 读取配置：默认 autoUpdate = true
    const configPath = path.join(os.homedir(), '.openclaw', 'openclaw.json');
    let autoUpdate = true;
    let intervalHours = 6;
    try {
      if (fs.existsSync(configPath)) {
        const cfg = JSON.parse(fs.readFileSync(configPath, 'utf8'));
        if (cfg.evolver && cfg.evolver.autoUpdate === false) autoUpdate = false;
        if (cfg.evolver && Number.isFinite(Number(cfg.evolver.autoUpdateIntervalHours))) {
          intervalHours = Number(cfg.evolver.autoUpdateIntervalHours);
        }
      }
    } catch (_) {}

    if (!autoUpdate) return;

    // 限速：每个间隔只检查一次
    const stateFile = path.join(MEMORY_DIR, 'evolver_update_check.json');
    const now = Date.now();
    const intervalMs = intervalHours * 60 * 60 * 1000;
    try {
      if (fs.existsSync(stateFile)) {
        const state = JSON.parse(fs.readFileSync(stateFile, 'utf8'));
        if (state.lastCheckedAt && (now - new Date(state.lastCheckedAt).getTime()) < intervalMs) {
          return; // 太早了，跳过
        }
      }
    } catch (_) {}

    let clawhubBin = null;
    const whichCmd = process.platform === 'win32' ? 'where clawhub' : 'which clawhub';
    const candidates = ['clawhub', path.join(os.homedir(), '.npm-global/bin/clawhub'), '/usr/local/bin/clawhub'];
    for (const c of candidates) {
      try {
        if (c === 'clawhub') {
          execSync(whichCmd, { stdio: 'ignore', timeout: 3000, windowsHide: true });
          clawhubBin = 'clawhub';
          break;
        }
        if (fs.existsSync(c)) { clawhubBin = c; break; }
      } catch (_) {}
    }
    if (!clawhubBin) return; // 没有 clawhub CLI 可用

    // 更新 evolver 和 feishu-evolver-wrapper
    const slugs = ['evolver', 'feishu-evolver-wrapper'];
    let updated = false;
    for (const slug of slugs) {
      try {
        const out = execSync(`${clawhubBin} update ${slug} --force`, {
          encoding: 'utf8',
          stdio: ['ignore', 'pipe', 'pipe'],
          timeout: 30000,
          cwd: path.resolve(REPO_ROOT, '..'),
          windowsHide: true,
        });
        if (out && !out.includes('already up to date') && !out.includes('not installed')) {
          console.log(`[AutoUpdate] ${slug}: ${out.trim().split('\n').pop()}`);
          updated = true;
        }
      } catch (e) {
        // 非致命：更新失败不应阻止进化
      }
    }

    // 写入状态
    try {
      const stateData = {
        lastCheckedAt: new Date(now).toISOString(),
        updated,
      };
      fs.writeFileSync(stateFile, JSON.stringify(stateData, null, 2) + '\n');
    } catch (_) {}

    if (updated) {
      console.log('[自动更新] 技能已更新。变更将在下次 wrapper 重启后生效。');
    }
  } catch (e) {
    // 整个自动更新过程都是非致命的
    console.log(`[自动更新] 检查失败（非致命）: ${e.message}`);
  }
}

function sleepMs(ms) {
  const t = Number(ms);
  const n = Number.isFinite(t) ? Math.max(0, t) : 0;
  return new Promise(resolve => setTimeout(resolve, n));
}

// 检查系统负载平均值，通过 os.loadavg() 获取。
// 返回 { load1m, load5m, load15m }。用于负载感知节流。
function getSystemLoad() {
  try {
    const loadavg = os.loadavg();
    return { load1m: loadavg[0], load5m: loadavg[1], load15m: loadavg[2] };
  } catch (e) {
    return { load1m: 0, load5m: 0, load15m: 0 };
  }
}

// 基于 CPU 核心数计算智能默认负载阈值
// 经验法则：
// - 单核：0.8-1.0（使用 0.9）
// - 多核：核心数 x 0.8-1.0（使用 0.9）
// - 生产环境：预留 20% 余量应对突发流量
function getDefaultLoadMax() {
  const cpuCount = os.cpus().length;
  if (cpuCount === 1) {
    return 0.9;
  } else {
    return cpuCount * 0.9;
  }
}

// 检查有多少 agent 会话正在活跃处理中（最近 N 分钟内修改过）。
// 如果 agent 正忙于用户对话，evolver 应该退避。
function getRecentActiveSessionCount(windowMs) {
  try {
    if (!fs.existsSync(AGENT_SESSIONS_DIR)) return 0;
    const now = Date.now();
    const w = Number.isFinite(windowMs) ? windowMs : 10 * 60 * 1000;
    return fs.readdirSync(AGENT_SESSIONS_DIR)
      .filter(f => f.endsWith('.jsonl') && !f.includes('.lock') && !f.startsWith('evolver_hand_'))
      .filter(f => {
        try { return (now - fs.statSync(path.join(AGENT_SESSIONS_DIR, f)).mtimeMs) < w; } catch (_) { return false; }
      }).length;
  } catch (_) { return 0; }
}

function determineBridgeEnabled() {
  const bridgeExplicit = process.env.EVOLVE_BRIDGE;
  if (bridgeExplicit !== undefined && bridgeExplicit !== '') {
    return String(bridgeExplicit).toLowerCase() !== 'false';
  }
  return Boolean(process.env.OPENCLAW_WORKSPACE);
}

async function run() {
  const bridgeEnabled = determineBridgeEnabled();
  const loopMode = ARGS.includes('--loop') || ARGS.includes('--mad-dog') || String(process.env.EVOLVE_LOOP || '').toLowerCase() === 'true';

  // 安全防护：如果另一个 evolver 手动代理正在运行，退避。
  // 防止 wrapper 重启时旧手动代理仍在执行导致的竞态条件。
  // 核心让步而不是启动竞争周期。
  if (process.platform !== 'win32') {
    try {
      const _psRace = require('child_process').execSync(
        'ps aux | grep "evolver_hand_" | grep "openclaw.*agent" | grep -v grep',
        { encoding: 'utf8', timeout: 5000, stdio: ['ignore', 'pipe', 'ignore'] }
      ).trim();
      if (_psRace && _psRace.length > 0) {
        console.log('[Evolver] 另一个 evolver 手动代理正在运行。让出当前周期。');
        return;
      }
    } catch (_) {
      // grep 退出码 1 = 无匹配 = 无冲突，可以安全继续
    }
  }

  // 安全防护：如果 agent 有过多活跃用户会话，退避。
  // Evolver 不能通过消耗模型并发来饿死用户对话。
  const QUEUE_MAX = Number.parseInt(process.env.EVOLVE_AGENT_QUEUE_MAX || '10', 10);
  const QUEUE_BACKOFF_MS = Number.parseInt(process.env.EVOLVE_AGENT_QUEUE_BACKOFF_MS || '60000', 10);
  const activeUserSessions = getRecentActiveSessionCount(10 * 60 * 1000);
  if (activeUserSessions > QUEUE_MAX) {
    console.log(`[Evolver] Agent 有 ${activeUserSessions} 个活跃用户会话（上限 ${QUEUE_MAX}）。退避 ${QUEUE_BACKOFF_MS}ms 以避免阻塞用户对话。`);
    writeDormantHypothesis({
      backoff_reason: 'active_sessions_exceeded',
      active_sessions: activeUserSessions,
      queue_max: QUEUE_MAX,
    });
    await sleepMs(QUEUE_BACKOFF_MS);
    return;
  }

  // 安全防护：系统负载感知。
  // 当系统负载过高时（如并发进程过多、I/O 负载重），
  // 退避以防止 evolver 加剧负载峰值。
  // Echo-MingXuan's Cycle #55 saw load spike from 0.02-0.50 to 1.30 before crash.
  const LOAD_MAX = parseFloat(process.env.EVOLVE_LOAD_MAX || String(getDefaultLoadMax()));
  const sysLoad = getSystemLoad();
  if (sysLoad.load1m > LOAD_MAX) {
    console.log(`[Evolver] 系统负载 ${sysLoad.load1m.toFixed(2)} 超过上限 ${LOAD_MAX.toFixed(1)}（基于 ${os.cpus().length} 个核心自动计算）。退避 ${QUEUE_BACKOFF_MS}ms。`);
    writeDormantHypothesis({
      backoff_reason: 'system_load_exceeded',
      system_load: { load1m: sysLoad.load1m, load5m: sysLoad.load5m, load15m: sysLoad.load15m },
      load_max: LOAD_MAX,
      cpu_cores: os.cpus().length,
    });
    await sleepMs(QUEUE_BACKOFF_MS);
    return;
  }

  // 循环节流：不要在前一个周期固化之前启动新周期。
  // 这防止 wrapper 在不等待手动代理完成的情况下"快速循环"大脑。
  if (bridgeEnabled && loopMode) {
    try {
      const st = readStateForSolidify();
      const lastRun = st && st.last_run ? st.last_run : null;
      const lastSolid = st && st.last_solidify ? st.last_solidify : null;
      if (lastRun && lastRun.run_id) {
        const pending = !lastSolid || !lastSolid.run_id || String(lastSolid.run_id) !== String(lastRun.run_id);
        if (pending) {
          writeDormantHypothesis({
            backoff_reason: 'loop_gating_pending_solidify',
            signals: lastRun && Array.isArray(lastRun.signals) ? lastRun.signals : [],
            selected_gene_id: lastRun && lastRun.selected_gene_id ? lastRun.selected_gene_id : null,
            mutation: lastRun && lastRun.mutation ? lastRun.mutation : null,
            personality_state: lastRun && lastRun.personality_state ? lastRun.personality_state : null,
            run_id: lastRun.run_id,
          });
          const raw = process.env.EVOLVE_PENDING_SLEEP_MS || process.env.EVOLVE_MIN_INTERVAL || '120000';
          const n = parseInt(String(raw), 10);
          const waitMs = Number.isFinite(n) ? Math.max(0, n) : 120000;
          await sleepMs(waitMs);
          return;
        }
      }
    } catch (e) {
      // 如果无法读取状态，继续（故障开放）以避免死锁。
    }
  }

  // 重置每周期环境标志，防止状态在周期之间泄漏。
  // 在 --loop 模式下，process.env 在周期之间持续存在。下面的断路器
  // 如果条件仍然成立，将重新设置 FORCE_INNOVATION。
  // CWD 恢复：如果工作目录在前一个周期中被删除
  // （例如，通过 git reset/restore 或目录删除），process.cwd() 会抛出
  // ENOENT 和所有后续操作都会失败。通过 chdir 到 REPO_ROOT 来恢复。
  try {
    process.cwd();
  } catch (e) {
    if (e && e.code === 'ENOENT') {
      console.warn('[Evolver] 工作目录丢失 (ENOENT)。正在恢复至 REPO_ROOT: ' + REPO_ROOT);
      try { process.chdir(REPO_ROOT); } catch (e2) {
        console.error('[Evolver] 工作目录恢复失败: ' + (e2 && e2.message ? e2.message : e2));
        throw e;
      }
    } else {
      throw e;
    }
  }

  delete process.env.FORCE_INNOVATION;

  // 安全防护：Git 仓库检查。
  // 固化、回滚和影响范围都依赖 git。没有 git 仓库，
  // 这些操作会静默产生空结果，导致数据丢失。
  try {
    execSync('git rev-parse --git-dir', { cwd: REPO_ROOT, encoding: 'utf8', stdio: ['ignore', 'pipe', 'pipe'], timeout: 5000 });
  } catch (_) {
    console.error('[Evolver] 致命错误: 不是 git 仓库 (' + REPO_ROOT + ')。');
    console.error('[Evolver] Evolver 需要 git 来支持回滚、影响范围计算和固化。');
    console.error('[Evolver] 请在项目根目录运行 "git init && git add -A && git commit -m init"，然后重试。');
    process.exitCode = 1;
    return;
  }

  const dormantHypothesis = readDormantHypothesis();
  if (dormantHypothesis) {
    console.log('[休眠假设] 已从上次退避中恢复部分状态: ' + (dormantHypothesis.backoff_reason || '未知'));
    clearDormantHypothesis();
  }

  const startTime = Date.now();
  verbose('--- evolve.run() start ---');
  verbose('Config: EVOLVE_STRATEGY=' + (process.env.EVOLVE_STRATEGY || '(default)') + ' EVOLVE_BRIDGE=' + (process.env.EVOLVE_BRIDGE || '(default)') + ' EVOLVE_LOOP=' + (process.env.EVOLVE_LOOP || 'false'));
  verbose('Config: EVOLVER_IDLE_FETCH_INTERVAL_MS=' + (process.env.EVOLVER_IDLE_FETCH_INTERVAL_MS || '(default 1800000)') + ' RANDOM_DRIFT=' + (process.env.RANDOM_DRIFT || 'false'));
  console.log('正在扫描会话日志...');

  // 确保所有 GEP 资产文件在任何操作之前存在。
  // 这防止外部工具在日志目录不存在时出现"没有那个文件或目录"错误
  // （grep、cat 等）引用可选的仅追加文件如 genes.jsonl。
  try { ensureAssetFiles(); } catch (e) {
    console.error(`[资产初始化] ensureAssetFiles 失败（非致命）: ${e.message}`);
  }

  // 维护：清理旧日志以保持目录扫描速度
  if (!IS_DRY_RUN) {
    performMaintenance();
  } else {
    console.log('[Maintenance] 已跳过（试运行模式）。');
  }

  // --- 修复循环断路器 ---
  // 检测 evolver 是否卡在"修复 -> 失败 -> 修复"的循环中。
  // 如果最近 N 个事件都是同一基因的失败修复，强制
  // 切换到创新意图以打破循环，而不是重试相同的修复。
  const REPAIR_LOOP_THRESHOLD = 3;
  try {
    const allEvents = readAllEvents();
    const recent = Array.isArray(allEvents) ? allEvents.slice(-REPAIR_LOOP_THRESHOLD) : [];
    if (recent.length >= REPAIR_LOOP_THRESHOLD) {
      const allRepairFailed = recent.every(e =>
        e && e.intent === 'repair' &&
        e.outcome && e.outcome.status === 'failed'
      );
      if (allRepairFailed) {
        const geneIds = recent.map(e => (e.genes_used && e.genes_used[0]) || 'unknown');
        const sameGene = geneIds.every(id => id === geneIds[0]);
        console.warn(`[断路器] 检测到连续 ${REPAIR_LOOP_THRESHOLD} 次修复失败${sameGene ? ` (基因: ${geneIds[0]})` : ''}。强制切换到创新意图以打破循环。`);
        // 设置环境标志，让下游代码读取以强制创新
        process.env.FORCE_INNOVATION = 'true';
      }
    }
  } catch (e) {
    // 非致命：如果无法读取事件，正常继续
    console.error(`[断路器] 检查失败（非致命）: ${e.message}`);
  }

  const recentMasterLog = readRealSessionLog();
  const todayLog = readRecentLog(TODAY_LOG);
  const memorySnippet = readMemorySnippet();
  const userSnippet = readUserSnippet();

  const cycleNum = getNextCycleId();
  const cycleId = `Cycle #${cycleNum}`;

  // 2. 检测工作区状态和本地覆盖
  // 逻辑：默认使用通用报告（message）
  let fileList = '';
  const skillsDir = path.join(REPO_ROOT, 'skills');

  // 默认报告：使用通用的 `message` 工具或 `process.env.EVOLVE_REPORT_CMD`（如果已设置）。
  // 这从核心逻辑中移除了对 'feishu-card' 的硬编码依赖。
  let reportingDirective = `Report requirement:
  - Use \`message\` tool.
  - Title: Evolution ${cycleId}
  - Status: [SUCCESS]
  - Changes: Detail exactly what was improved.`;

  // Wrapper 注入点：wrapper 可以通过环境变量注入自定义报告指令。
  if (process.env.EVOLVE_REPORT_DIRECTIVE) {
    reportingDirective = process.env.EVOLVE_REPORT_DIRECTIVE.replace('__CYCLE_ID__', cycleId);
  } else if (process.env.EVOLVE_REPORT_CMD) {
    reportingDirective = `Report requirement (custom):
  - Execute the custom report command:
    \`\`\`
    ${process.env.EVOLVE_REPORT_CMD.replace('__CYCLE_ID__', cycleId)}
    \`\`\`
  - Ensure you pass the status and action details.`;
  }

  // 处理审查模式标志 (--review)
  if (IS_REVIEW_MODE) {
    reportingDirective +=
      '\n  - REVIEW PAUSE: After generating the fix but BEFORE applying significant edits, ask the user for confirmation.';
  }

  const SKILLS_CACHE_FILE = path.join(MEMORY_DIR, 'skills_list_cache.json');

  try {
    if (fs.existsSync(skillsDir)) {
      // 检查缓存有效性（skills 文件夹的 mtime 与缓存文件比较）
      let useCache = false;
      const dirStats = fs.statSync(skillsDir);
      if (fs.existsSync(SKILLS_CACHE_FILE)) {
        const cacheStats = fs.statSync(SKILLS_CACHE_FILE);
        const CACHE_TTL = 1000 * 60 * 60 * 6; // 6 小时
        const isFresh = Date.now() - cacheStats.mtimeMs < CACHE_TTL;

        // 如果缓存新鲜且比目录（结构变更）更新，则使用缓存
        if (isFresh && cacheStats.mtimeMs > dirStats.mtimeMs) {
          try {
            const cached = JSON.parse(fs.readFileSync(SKILLS_CACHE_FILE, 'utf8'));
            fileList = cached.list;
            useCache = true;
          } catch (e) {}
        }
      }

      if (!useCache) {
        const skills = fs
          .readdirSync(skillsDir, { withFileTypes: true })
          .filter(dirent => dirent.isDirectory())
          .map(dirent => {
            const name = dirent.name;
            let desc = 'No description';
            try {
              const pkg = require(path.join(skillsDir, name, 'package.json'));
              if (pkg.description) desc = pkg.description.slice(0, 100) + (pkg.description.length > 100 ? '...' : '');
            } catch (e) {
              try {
                const skillMdPath = path.join(skillsDir, name, 'SKILL.md');
                if (fs.existsSync(skillMdPath)) {
                  const skillMd = fs.readFileSync(skillMdPath, 'utf8');
                  // 策略 1：YAML 前言 (description: ...)
                  const yamlMatch = skillMd.match(/^description:\s*(.*)$/m);
                  if (yamlMatch) {
                    desc = yamlMatch[1].trim();
                  } else {
                    // 策略 2：第一个非标题、非空行
                    const lines = skillMd.split('\n');
                    for (const line of lines) {
                      const trimmed = line.trim();
                      if (
                        trimmed &&
                        !trimmed.startsWith('#') &&
                        !trimmed.startsWith('---') &&
                        !trimmed.startsWith('```')
                      ) {
                        desc = trimmed;
                        break;
                      }
                    }
                  }
                  if (desc.length > 100) desc = desc.slice(0, 100) + '...';
                }
              } catch (e2) {}
            }
            return `- **${name}**: ${desc}`;
          });
        fileList = skills.join('\n');

        // 写入缓存
        try {
          fs.writeFileSync(SKILLS_CACHE_FILE, JSON.stringify({ list: fileList }, null, 2));
        } catch (e) {}
      }
    }
  } catch (e) {
    fileList = `Error listing skills: ${e.message}`;
  }

  const mutationDirective = getMutationDirective(recentMasterLog);
  const healthReport = checkSystemHealth();

  // 功能：情绪感知（模式 E - 个性化）
  let moodStatus = 'Mood: Unknown';
  try {
    const moodFile = path.join(MEMORY_DIR, 'mood.json');
    if (fs.existsSync(moodFile)) {
      const moodData = JSON.parse(fs.readFileSync(moodFile, 'utf8'));
      moodStatus = `Mood: ${moodData.current_mood || 'Neutral'} (Intensity: ${moodData.intensity || 0})`;
    }
  } catch (e) {}

  const scanTime = Date.now() - startTime;
  const memorySize = fs.existsSync(MEMORY_FILE) ? fs.statSync(MEMORY_FILE).size : 0;

  let syncDirective = 'Workspace sync: optional/disabled in this environment.';

  // 检查 git-sync 技能可用性
  const hasGitSync = fs.existsSync(path.join(skillsDir, 'git-sync'));
  if (hasGitSync) {
    syncDirective = 'Workspace sync: run skills/git-sync/sync.sh "Evolution: Workspace Sync"';
  }

  const genes = loadGenes();
  const capsules = loadCapsules();
  const recentEvents = (() => {
    try {
      const all = readAllEvents();
      return Array.isArray(all) ? all.filter(e => e && e.type === 'EvolutionEvent').slice(-80) : [];
    } catch (e) {
      return [];
    }
  })();
  const signals = extractSignals({
    recentSessionTranscript: recentMasterLog,
    todayLog,
    memorySnippet,
    userSnippet,
    recentEvents,
  });

  verbose('Signals extracted (' + signals.length + '):', signals.join(', '));
  verbose('Recent events: ' + recentEvents.length + ', session log size: ' + recentMasterLog.length + ' chars');

  if (dormantHypothesis && Array.isArray(dormantHypothesis.signals) && dormantHypothesis.signals.length > 0) {
    const dormantSignals = dormantHypothesis.signals;
    let injected = 0;
    for (let dsi = 0; dsi < dormantSignals.length; dsi++) {
      if (!signals.includes(dormantSignals[dsi])) {
        signals.push(dormantSignals[dsi]);
        injected++;
      }
    }
    if (injected > 0) {
      console.log('[休眠假设] 已从上次中断的周期注入 ' + injected + ' 个信号。');
    }
  }

  // --- 空闲周期节流：在饱和期间跳过 Hub API 调用以节省额度 ---
  let _idleFetchInterval = parseInt(String(process.env.EVOLVER_IDLE_FETCH_INTERVAL_MS || ''), 10);
  if (!Number.isFinite(_idleFetchInterval) || _idleFetchInterval <= 0) _idleFetchInterval = 1800000;
  let skipHubCalls = false;

  if (shouldSkipHubCalls(signals)) {
    const _elapsed = Date.now() - _lastHubFetchMs;
    if (_lastHubFetchMs > 0 && _elapsed < _idleFetchInterval) {
      skipHubCalls = true;
      console.log('[空闲节流] 已饱和，无可操作信号。跳过 Hub API 调用（上次获取距今 ' + Math.round(_elapsed / 1000) + '秒，阈值 ' + Math.round(_idleFetchInterval / 1000) + '秒）。');
    } else {
      console.log('[空闲节流] 已饱和但获取间隔已到（' + Math.round((Date.now() - _lastHubFetchMs) / 1000) + '秒）。执行定期 Hub 检查。');
    }
  }

  // --- Hub 任务自动认领（含主动问题） ---
  // 从当前上下文生成问题，附带在获取调用中，
  // 然后选择最佳任务并自动认领。
  let activeTask = null;
  let proactiveQuestions = [];
  if (!skipHubCalls) {
    try {
      proactiveQuestions = generateQuestions({
        signals,
        recentEvents,
        sessionTranscript: recentMasterLog,
        memorySnippet: memorySnippet,
      });
      if (proactiveQuestions.length > 0) {
        console.log(`[问题生成器] 已生成 ${proactiveQuestions.length} 个主动问题。`);
      }
    } catch (e) {
      console.log(`[问题生成器] 生成失败（非致命）: ${e.message}`);
    }

    // --- 自动 GitHub 问题报告器 ---
    // 当检测到持续失败时，向上游仓库提交问题，
    // 包含清理后的日志和环境信息。
    try {
      await maybeReportIssue({
        signals,
        recentEvents,
        sessionLog: recentMasterLog,
      });
    } catch (e) {
      console.log(`[问题报告器] 检查失败（非致命）: ${e.message}`);
    }
  }

  // 经验库：在获取期间从 Hub 接收的经验
  let hubLessons = [];

  if (!skipHubCalls) {
    _lastHubFetchMs = Date.now();
    try {
      const fetchResult = await fetchTasks({ questions: proactiveQuestions });
      const hubTasks = fetchResult.tasks || [];

      if (fetchResult.questions_created && fetchResult.questions_created.length > 0) {
        const created = fetchResult.questions_created.filter(function(q) { return !q.error; });
        const failed = fetchResult.questions_created.filter(function(q) { return q.error; });
        if (created.length > 0) {
          console.log(`[问题生成器] Hub 已接受 ${created.length} 个问题作为悬赏任务。`);
        }
        if (failed.length > 0) {
          console.log(`[问题生成器] Hub 拒绝了 ${failed.length} 个问题: ${failed.map(function(q) { return q.error; }).join(', ')}`);
        }
      }

      // 经验库：从 Hub 捕获相关经验
      if (Array.isArray(fetchResult.relevant_lessons) && fetchResult.relevant_lessons.length > 0) {
        hubLessons = fetchResult.relevant_lessons;
        console.log(`[经验库] 从生态系统接收到 ${hubLessons.length} 条经验。`);
      }

      if (hubTasks.length > 0) {
        let taskMemoryEvents = [];
        try {
          const { tryReadMemoryGraphEvents } = require('./gep/memoryGraph');
          taskMemoryEvents = tryReadMemoryGraphEvents(1000);
        } catch (e) {
          console.warn('[任务接收器] MemoryGraph 读取失败（任务选择将继续，但无历史记录）:', e && e.message || e);
        }
        const best = selectBestTask(hubTasks, taskMemoryEvents);
        if (best) {
          const alreadyClaimed = best.status === 'claimed';
          let claimed = alreadyClaimed;
          if (!alreadyClaimed) {
            const commitDeadline = estimateCommitmentDeadline(best);
            claimed = await claimTask(best.id || best.task_id, commitDeadline ? { commitment_deadline: commitDeadline } : undefined);
            if (claimed && commitDeadline) {
              best._commitment_deadline = commitDeadline;
              console.log(`[任务承诺] 截止时间已设置: ${commitDeadline}`);
            }
          }
          if (claimed) {
            activeTask = best;
            const taskSignals = taskToSignals(best);
            for (const sig of taskSignals) {
              if (!signals.includes(sig)) signals.unshift(sig);
            }
            console.log(`[任务接收器] ${alreadyClaimed ? '继续' : '已认领'}任务: "${best.title || best.id}" (${taskSignals.length} 个信号已注入)`);
          }
        }
      }
    } catch (e) {
      console.log(`[任务接收器] 获取/认领失败（非致命）: ${e.message}`);
    }
  }

  // --- 任务承诺：检查心跳中的逾期任务 ---
  // 如果 Hub 报告了逾期任务，通过将其信号注入
  // 前面来优先恢复它们。这不改变活跃任务选择（逾期
  // 任务应该已经从之前的周期中被认领/激活）。
  try {
    const { consumeOverdueTasks } = require('./gep/a2aProtocol');
    const overdueTasks = consumeOverdueTasks();
    if (overdueTasks.length > 0) {
      for (const ot of overdueTasks) {
        const otId = ot.task_id || ot.id;
        if (activeTask && (activeTask.id === otId || activeTask.task_id === otId)) {
          console.warn(`[任务承诺] 活跃任务 "${activeTask.title || otId}" 已逾期 -- 优先完成。`);
          signals.unshift('overdue_task', 'urgent');
          break;
        }
      }
    }
  } catch (e) {
    console.warn('[任务承诺] 逾期任务检查失败（非致命）:', e && e.message || e);
  }

  // --- 工作池：从心跳 available_work 中选择任务（延迟认领） ---
  // 只记住最佳任务并注入其信号；实际认领+完成
  // 在成功的进化周期后，实际认领+完成在 solidify.js 中原子性地发生。
  if (!activeTask && process.env.WORKER_ENABLED === '1') {
    try {
      const { consumeAvailableWork } = require('./gep/a2aProtocol');
      const workerTasks = consumeAvailableWork();
      if (workerTasks.length > 0) {
        let taskMemoryEvents = [];
        try {
          const { tryReadMemoryGraphEvents } = require('./gep/memoryGraph');
          taskMemoryEvents = tryReadMemoryGraphEvents(1000);
        } catch (e) {
          console.warn('[工作池] MemoryGraph 读取失败（任务选择将继续，但无历史记录）:', e && e.message || e);
        }
        const best = selectBestTask(workerTasks, taskMemoryEvents);
        if (best) {
          activeTask = best;
          activeTask._worker_pending = true;
          const taskSignals = taskToSignals(best);
          for (const sig of taskSignals) {
            if (!signals.includes(sig)) signals.unshift(sig);
          }
          console.log(`[工作池] 已选择工作池任务（延迟认领）: "${best.title || best.id}" (${taskSignals.length} 个信号已注入)`);
        }
      }
    } catch (e) {
      console.log(`[工作池] 任务选择失败（非致命）: ${e.message}`);
    }
  }

  const recentErrorMatches = recentMasterLog.match(/\[ERROR|Error:|Exception:|FAIL|Failed|"isError":true/gi) || [];
  const recentErrorCount = recentErrorMatches.length;

  const evidence = {
    // 保持简短；不要在图谱中存储完整转录。
    recent_session_tail: String(recentMasterLog || '').slice(-6000),
    today_log_tail: String(todayLog || '').slice(-2500),
  };

  const sessionScope = getSessionScope();
  const observations = {
    agent: AGENT_NAME,
    session_scope: sessionScope || null,
    drift_enabled: IS_RANDOM_DRIFT,
    review_mode: IS_REVIEW_MODE,
    dry_run: IS_DRY_RUN,
    system_health: healthReport,
    mood: moodStatus,
    scan_ms: scanTime,
    memory_size_bytes: memorySize,
    recent_error_count: recentErrorCount,
    node: process.version,
    platform: process.platform,
    cwd: process.cwd(),
    evidence,
  };

  if (sessionScope) {
    console.log(`[会话作用域] 当前活跃作用域: "${sessionScope}"。进化状态和记忆图谱已隔离。`);
  }

  // 记忆图谱：用推断的结果关闭最后一个动作（仅追加图谱，可变状态）。
  try {
    recordOutcomeFromState({ signals, observations });
  } catch (e) {
    // 如果无法读/写记忆图谱，拒绝进化（不做"无记忆进化"）。
    console.error(`[MemoryGraph] 结果写入失败: ${e.message}`);
    console.error(`[MemoryGraph] 拒绝在没有因果记忆的情况下进化。目标: ${memoryGraphPath()}`);
    throw new Error(`MemoryGraph Outcome write failed: ${e.message}`);
  }

  // 记忆图谱：将当前信号记录为一等节点。如果失败，拒绝进化。
  try {
    recordSignalSnapshot({ signals, observations });
  } catch (e) {
    console.error(`[MemoryGraph] 信号快照写入失败: ${e.message}`);
    console.error(`[MemoryGraph] 拒绝在没有因果记忆的情况下进化。目标: ${memoryGraphPath()}`);
    throw new Error(`MemoryGraph Signal snapshot write failed: ${e.message}`);
  }

  // 能力候选：提取、持久化和构建预览。
  const { capabilityCandidatesPreview, externalCandidatesPreview } = buildCandidatePreviews({
    signals,
    recentSessionTranscript: recentMasterLog,
  });

  // 搜索优先进化：在本地推理之前查询 Hub 以获取可复用解决方案。
  let hubHit = null;
  if (!skipHubCalls) {
    try {
      hubHit = await hubSearch(signals, { timeoutMs: 8000 });
      if (hubHit && hubHit.hit) {
        console.log(`[搜索优先] Hub 命中: asset=${hubHit.asset_id}, score=${hubHit.score}, mode=${hubHit.mode}`);
      } else {
        console.log(`[搜索优先] Hub 无匹配（原因: ${hubHit && hubHit.reason ? hubHit.reason : '未知'}）。继续本地进化。`);
      }
    } catch (e) {
      console.log(`[搜索优先] Hub 搜索失败（非致命）: ${e.message}`);
      hubHit = { hit: false, reason: 'exception' };
    }
  } else {
    hubHit = { hit: false, reason: 'idle_skip' };
    console.log('[空闲节流] hubSearch 已跳过（空闲周期）。');
  }

  // 记忆图谱推理：优先选择高置信路径，抑制已知低成功路径（除非明确漂移）。
  let memoryAdvice = null;
  try {
    memoryAdvice = getMemoryAdvice({ signals, genes, driftEnabled: IS_RANDOM_DRIFT });
  } catch (e) {
    console.error(`[MemoryGraph] 读取失败: ${e.message}`);
    console.error(`[MemoryGraph] 拒绝在没有因果记忆的情况下进化。目标: ${memoryGraphPath()}`);
    throw new Error(`MemoryGraph Read failed: ${e.message}`);
  }

  // 反思阶段：定期暂停以评估进化策略。
  try {
    const cycleState = fs.existsSync(STATE_FILE) ? JSON.parse(fs.readFileSync(STATE_FILE, 'utf8')) : {};
    const cycleCount = cycleState.cycleCount || 0;
    if (shouldReflect({ cycleCount, recentEvents })) {
      const narrativeSummary = loadNarrativeSummary(3000);
      const reflectionCtx = buildReflectionContext({
        recentEvents,
        signals,
        memoryAdvice,
        narrative: narrativeSummary,
      });
      recordReflection({
        cycle_count: cycleCount,
        signals_snapshot: signals.slice(0, 20),
        preferred_gene: memoryAdvice && memoryAdvice.preferredGeneId ? memoryAdvice.preferredGeneId : null,
        banned_genes: memoryAdvice && Array.isArray(memoryAdvice.bannedGeneIds) ? memoryAdvice.bannedGeneIds : [],
        context_preview: reflectionCtx.slice(0, 1000),
      });
      console.log(`[反思] 战略反思已记录于周期 ${cycleCount}。`);
    }
  } catch (e) {
    console.log('[反思] 失败（非致命）: ' + (e && e.message ? e.message : e));
  }

  let recentFailedCapsules = [];
  try {
    recentFailedCapsules = readRecentFailedCapsules(50);
  } catch (e) {
    console.log('[失败胶囊] 读取失败（非致命）: ' + e.message);
  }

  // 心跳提示：新颖度分数和能力差距，用于多样性导向漂移
  let heartbeatNovelty = null;
  let heartbeatCapGaps = [];
  try {
    const { getNoveltyHint, getCapabilityGaps: getCapGaps } = require('./gep/a2aProtocol');
    heartbeatNovelty = getNoveltyHint();
    heartbeatCapGaps = getCapGaps() || [];
  } catch (e) {}

  const { selectedGene, capsuleCandidates, selector } = selectGeneAndCapsule({
    genes,
    capsules,
    signals,
    memoryAdvice,
    driftEnabled: IS_RANDOM_DRIFT,
    failedCapsules: recentFailedCapsules,
    capabilityGaps: heartbeatCapGaps,
    noveltyScore: heartbeatNovelty && Number.isFinite(heartbeatNovelty.score) ? heartbeatNovelty.score : null,
  });

  const selectedBy = memoryAdvice && memoryAdvice.preferredGeneId ? 'memory_graph+selector' : 'selector';
  const capsulesUsed = Array.isArray(capsuleCandidates)
    ? capsuleCandidates.map(c => (c && c.id ? String(c.id) : null)).filter(Boolean)
    : [];
  const selectedCapsuleId = capsulesUsed.length ? capsulesUsed[0] : null;
  const strategyPolicy = computeAdaptiveStrategyPolicy({
    recentEvents,
    selectedGene,
    signals,
  });

  verbose('Gene selection: gene=' + (selectedGene ? selectedGene.id : '(none)') + ' capsule=' + (selectedCapsuleId || '(none)') + ' selectedBy=' + selectedBy + ' selector=' + (selector || '(none)'));
  verbose('Strategy policy: name=' + strategyPolicy.name + ' forceInnovate=' + strategyPolicy.forceInnovate + ' cautious=' + strategyPolicy.cautiousExecution + ' maxFiles=' + strategyPolicy.blastRadiusMaxFiles);
  if (memoryAdvice) {
    verbose('Memory advice: preferred=' + (memoryAdvice.preferredGeneId || '(none)') + ' banned=[' + (Array.isArray(memoryAdvice.bannedGeneIds) ? memoryAdvice.bannedGeneIds.join(',') : '') + ']');
  }

  // 人格选择（自然选择 + 触发时的小变异）。
  // 此状态持久化在 MEMORY_DIR 中，被视为进化控制面（非角色扮演）。
  const personalitySelection = selectPersonalityForRun({
    driftEnabled: IS_RANDOM_DRIFT,
    signals,
    recentEvents,
  });
  const personalityState = personalitySelection && personalitySelection.personality_state ? personalitySelection.personality_state : null;

  // 变异对象是每次进化运行的必需项。
  const tail = Array.isArray(recentEvents) ? recentEvents.slice(-6) : [];
  const tailOutcomes = tail
    .map(e => (e && e.outcome && e.outcome.status ? String(e.outcome.status) : null))
    .filter(Boolean);
  const stableSuccess = tailOutcomes.length >= 6 && tailOutcomes.every(s => s === 'success');
  const tailAvgScore =
    tail.length > 0
      ? tail.reduce((acc, e) => acc + (e && e.outcome && Number.isFinite(Number(e.outcome.score)) ? Number(e.outcome.score) : 0), 0) /
        tail.length
      : 0;
  const innovationPressure =
    !IS_RANDOM_DRIFT &&
    personalityState &&
    Number.isFinite(Number(personalityState.creativity)) &&
    Number(personalityState.creativity) >= 0.75 &&
    stableSuccess &&
    tailAvgScore >= 0.7;
  const forceInnovation =
    String(process.env.FORCE_INNOVATION || process.env.EVOLVE_FORCE_INNOVATION || '').toLowerCase() === 'true';
  const mutationInnovateMode = !!IS_RANDOM_DRIFT || !!innovationPressure || !!forceInnovation || !!strategyPolicy.forceInnovate;
  const mutationSignals = innovationPressure ? [...(Array.isArray(signals) ? signals : []), 'stable_success_plateau'] : signals;
  const mutationSignalsEffective = (forceInnovation || strategyPolicy.forceInnovate)
    ? [...(Array.isArray(mutationSignals) ? mutationSignals : []), 'force_innovation']
    : mutationSignals;

  const allowHighRisk =
    !!IS_RANDOM_DRIFT &&
    !!personalitySelection &&
    !!personalitySelection.personality_known &&
    personalityState &&
    isHighRiskMutationAllowed(personalityState) &&
    Number(personalityState.rigor) >= 0.8 &&
    Number(personalityState.risk_tolerance) <= 0.3 &&
    !(Array.isArray(signals) && signals.includes('log_error'));
  const mutation = buildMutation({
    signals: mutationSignalsEffective,
    selectedGene,
    driftEnabled: mutationInnovateMode,
    personalityState,
    allowHighRisk,
  });

  verbose('Mutation: category=' + (mutation && mutation.category || '?') + ' risk=' + (mutation && mutation.risk_level || '?') + ' innovateMode=' + mutationInnovateMode + ' forceInnovation=' + forceInnovation + ' allowHighRisk=' + allowHighRisk);
  verbose('Hub: hubHit=' + (hubHit && hubHit.hit ? 'true (score=' + hubHit.score + ' mode=' + hubHit.mode + ')' : 'false (' + (hubHit && hubHit.reason || 'unknown') + ')'));

  // 记忆图谱：记录连接信号 -> 动作的假设。如果失败，拒绝进化。
  let hypothesisId = null;
  try {
    const hyp = recordHypothesis({
      signals,
      mutation,
      personality_state: personalityState,
      selectedGene,
      selector,
      driftEnabled: mutationInnovateMode,
      selectedBy,
      capsulesUsed,
      observations,
    });
    hypothesisId = hyp && hyp.hypothesisId ? hyp.hypothesisId : null;
  } catch (e) {
    console.error(`[MemoryGraph] 假设写入失败: ${e.message}`);
    console.error(`[MemoryGraph] 拒绝在没有因果记忆的情况下进化。目标: ${memoryGraphPath()}`);
    throw new Error(`MemoryGraph Hypothesis write failed: ${e.message}`);
  }

  // 记忆图谱：记录此次运行选择的因果路径。如果失败，拒绝输出变异提示。
  try {
    recordAttempt({
      signals,
      mutation,
      personality_state: personalityState,
      selectedGene,
      selector,
      driftEnabled: mutationInnovateMode,
      selectedBy,
      hypothesisId,
      capsulesUsed,
      observations,
    });
  } catch (e) {
    console.error(`[MemoryGraph] 尝试写入失败: ${e.message}`);
    console.error(`[MemoryGraph] 拒绝在没有因果记忆的情况下进化。目标: ${memoryGraphPath()}`);
    throw new Error(`MemoryGraph Attempt write failed: ${e.message}`);
  }

  // 固化状态：捕获最小化、可审计的上下文用于补丁后验证 + 资产写入。
  // 这在补丁应用后强制执行严格的协议闭合。
  try {
    const runId = `run_${Date.now()}`;
    const parentEventId = getLastEventId();

    // 基线快照（在任何编辑之前）。
    let baselineUntracked = [];
    let baselineHead = null;
    try {
      const out = execSync('git ls-files --others --exclude-standard', {
        cwd: REPO_ROOT,
        encoding: 'utf8',
        stdio: ['ignore', 'pipe', 'ignore'],
        timeout: 4000,
        windowsHide: true,
      });
      baselineUntracked = String(out)
        .split('\n')
        .map(l => l.trim())
        .filter(Boolean);
    } catch (e) {
      console.warn('[固化状态] 读取基线未跟踪文件失败:', e && e.message || e);
    }

    try {
      const out = execSync('git rev-parse HEAD', {
        cwd: REPO_ROOT,
        encoding: 'utf8',
        stdio: ['ignore', 'pipe', 'ignore'],
        timeout: 4000,
        windowsHide: true,
      });
      baselineHead = String(out || '').trim() || null;
    } catch (e) {
      console.warn('[固化状态] 读取 git HEAD 失败:', e && e.message || e);
    }

    const maxFiles = strategyPolicy && Number.isFinite(Number(strategyPolicy.blastRadiusMaxFiles))
      ? Number(strategyPolicy.blastRadiusMaxFiles)
      : (
        selectedGene && selectedGene.constraints && Number.isFinite(Number(selectedGene.constraints.max_files))
          ? Number(selectedGene.constraints.max_files)
          : 12
      );
    const blastRadiusEstimate = {
      files: Number.isFinite(maxFiles) && maxFiles > 0 ? maxFiles : 0,
      lines: Number.isFinite(maxFiles) && maxFiles > 0 ? Math.round(maxFiles * 80) : 0,
    };

    // 合并到现有状态以保留 last_solidify（不要擦除它）。
    const prevState = readStateForSolidify();
    prevState.last_run = {
        run_id: runId,
        created_at: new Date().toISOString(),
        parent_event_id: parentEventId || null,
        selected_gene_id: selectedGene && selectedGene.id ? selectedGene.id : null,
        selected_capsule_id: selectedCapsuleId,
        selector: selector || null,
        signals: Array.isArray(signals) ? signals : [],
        mutation: mutation || null,
        mutation_id: mutation && mutation.id ? mutation.id : null,
        personality_state: personalityState || null,
        personality_key: personalitySelection && personalitySelection.personality_key ? personalitySelection.personality_key : null,
        personality_known: !!(personalitySelection && personalitySelection.personality_known),
        personality_mutations:
          personalitySelection && Array.isArray(personalitySelection.personality_mutations)
            ? personalitySelection.personality_mutations
            : [],
        drift: !!IS_RANDOM_DRIFT,
        selected_by: selectedBy,
        source_type: hubHit && hubHit.hit ? (hubHit.mode === 'direct' ? 'reused' : 'reference') : 'generated',
        reused_asset_id: hubHit && hubHit.hit ? (hubHit.asset_id || null) : null,
        reused_source_node: hubHit && hubHit.hit ? (hubHit.source_node_id || null) : null,
        reused_chain_id: hubHit && hubHit.hit ? (hubHit.chain_id || null) : null,
        baseline_untracked: baselineUntracked,
        baseline_git_head: baselineHead,
        blast_radius_estimate: blastRadiusEstimate,
        strategy_policy: strategyPolicy,
        active_task_id: activeTask ? (activeTask.id || activeTask.task_id || null) : null,
        active_task_title: activeTask ? (activeTask.title || null) : null,
        worker_assignment_id: activeTask ? (activeTask._worker_assignment_id || null) : null,
        worker_pending: activeTask ? (activeTask._worker_pending || false) : false,
        commitment_deadline: activeTask ? (activeTask._commitment_deadline || null) : null,
        applied_lessons: hubLessons.map(function(l) { return l.lesson_id; }).filter(Boolean),
        hub_lessons: hubLessons,
      };
    writeStateForSolidify(prevState);

    if (hubHit && hubHit.hit) {
      const assetAction = hubHit.mode === 'direct' ? 'asset_reuse' : 'asset_reference';
      logAssetCall({
        run_id: runId,
        action: assetAction,
        asset_id: hubHit.asset_id || null,
        asset_type: hubHit.match && hubHit.match.type ? hubHit.match.type : null,
        source_node_id: hubHit.source_node_id || null,
        chain_id: hubHit.chain_id || null,
        score: hubHit.score || null,
        mode: hubHit.mode,
        signals: Array.isArray(signals) ? signals : [],
        extra: {
          selected_gene_id: selectedGene && selectedGene.id ? selectedGene.id : null,
          task_id: activeTask ? (activeTask.id || activeTask.task_id || null) : null,
        },
      });
    }
  } catch (e) {
    console.error(`[固化状态] 写入失败: ${e.message}`);
  }

  if (skipHubCalls) {
    console.log('[空闲节流] 空闲周期完成。提示生成和 bridge 生成已跳过。');
    return;
  }

  const genesPreview = `\`\`\`json\n${JSON.stringify(genes.slice(0, 6), null, 2)}\n\`\`\``;
  const capsulesPreview = `\`\`\`json\n${JSON.stringify(capsules.slice(-3), null, 2)}\n\`\`\``;

  const reviewNote = IS_REVIEW_MODE
    ? 'Review mode: before significant edits, pause and ask the user for confirmation.'
    : 'Review mode: disabled.';

  // 构建近期进化历史摘要用于上下文注入
  const recentHistorySummary = (() => {
    if (!recentEvents || recentEvents.length === 0) return '(no prior evolution events)';
    const last8 = recentEvents.slice(-8);
    const lines = last8.map((evt, idx) => {
      const sigs = Array.isArray(evt.signals) ? evt.signals.slice(0, 3).join(', ') : '?';
      const gene = Array.isArray(evt.genes_used) && evt.genes_used.length ? evt.genes_used[0] : 'none';
      const outcome = evt.outcome && evt.outcome.status ? evt.outcome.status : '?';
      const ts = evt.meta && evt.meta.at ? evt.meta.at : (evt.id || '');
      return `  ${idx + 1}. [${evt.intent || '?'}] signals=[${sigs}] gene=${gene} outcome=${outcome} @${ts}`;
    });
    return lines.join('\n');
  })();

  const context = `
Runtime state:
- System health: ${healthReport}
- Agent state: ${moodStatus}
- Scan duration: ${scanTime}ms
- Memory size: ${memorySize} bytes
- Skills available (if any):
${fileList || '[skills directory not found]'}

Notes:
- ${reviewNote}
- ${reportingDirective}
- ${syncDirective}

Recent Evolution History (last 8 cycles -- DO NOT repeat the same intent+signal+gene):
${recentHistorySummary}
IMPORTANT: If you see 3+ consecutive "repair" cycles with the same gene, you MUST switch to "innovate" intent.
${(() => {
  // 从近期事件计算连续失败次数用于上下文注入
  let cfc = 0;
  const evts = Array.isArray(recentEvents) ? recentEvents : [];
  for (let i = evts.length - 1; i >= 0; i--) {
    if (evts[i] && evts[i].outcome && evts[i].outcome.status === 'failed') cfc++;
    else break;
  }
  if (cfc >= 3) {
    return `\nFAILURE STREAK WARNING: The last ${cfc} cycles ALL FAILED. You MUST change your approach.\n- Do NOT repeat the same gene/strategy. Pick a completely different approach.\n- If the error is external (API down, binary missing), mark as FAILED and move on.\n- Prefer a minimal safe innovate cycle over yet another failing repair.`;
  }
  return '';
})()}

External candidates (A2A receive zone; staged only, never execute directly):
${externalCandidatesPreview}

Global memory (MEMORY.md):
\`\`\`
${memorySnippet}
\`\`\`

User registry (USER.md):
\`\`\`
${userSnippet}
\`\`\`

Recent memory snippet:
\`\`\`
${todayLog.slice(-3000)}
\`\`\`

Recent session transcript:
\`\`\`
${recentMasterLog}
\`\`\`

Mutation directive:
${mutationDirective}
`.trim();

  // 构建提示：在直接复用模式下，使用最小复用提示。
  // 在参考模式（或无命中）下，使用完整的 GEP 提示并注入 hub 匹配。
  const isDirectReuse = hubHit && hubHit.hit && hubHit.mode === 'direct';
  const hubMatchedBlock = hubHit && hubHit.hit && hubHit.mode === 'reference'
    ? buildHubMatchedBlock({ capsule: hubHit.match })
    : null;

  const prompt = isDirectReuse
    ? buildReusePrompt({
        capsule: hubHit.match,
        signals,
        nowIso: new Date().toISOString(),
      })
    : buildGepPrompt({
        nowIso: new Date().toISOString(),
        context,
        signals,
        selector,
        parentEventId: getLastEventId(),
        selectedGene,
        capsuleCandidates,
        genesPreview,
        capsulesPreview,
        capabilityCandidatesPreview,
        externalCandidatesPreview,
        hubMatchedBlock,
        strategyPolicy,
        failedCapsules: recentFailedCapsules,
        hubLessons,
      });

  // 可选：为 wrapper 输出紧凑的思考过程块（噪音可控）。
  const emitThought = String(process.env.EVOLVE_EMIT_THOUGHT_PROCESS || '').toLowerCase() === 'true';
  if (emitThought) {
    const s = Array.isArray(signals) ? signals : [];
    const thought = [
      `cycle_id: ${cycleId}`,
      `signals_count: ${s.length}`,
      `signals: ${s.slice(0, 12).join(', ')}${s.length > 12 ? ' ...' : ''}`,
      `selected_gene: ${selectedGene && selectedGene.id ? String(selectedGene.id) : '(none)'}`,
      `selected_capsule: ${selectedCapsuleId ? String(selectedCapsuleId) : '(none)'}`,
      `mutation_category: ${mutation && mutation.category ? String(mutation.category) : '(none)'}`,
      `force_innovation: ${forceInnovation ? 'true' : 'false'}`,
      `source_type: ${hubHit && hubHit.hit ? (isDirectReuse ? 'reused' : 'reference') : 'generated'}`,
      `hub_reuse_mode: ${isDirectReuse ? 'direct' : hubMatchedBlock ? 'reference' : 'none'}`,
    ].join('\n');
    console.log(`[思考过程]\n${thought}\n[/思考过程]`);
  }

  const printPrompt = String(process.env.EVOLVE_PRINT_PROMPT || '').toLowerCase() === 'true';

  // 默认行为 (v1.4.1+)：通过 sessions_spawn 桥接提示 -> 子代理来"默认执行"。
  // 本项目是大脑；手动代理是生成的执行器代理。Wrapper 可以用 EVOLVE_BRIDGE=false 禁用桥接。
  if (bridgeEnabled) {
    // 尽可能复用存储在固化状态中的 run_id。
    let runId = `run_${Date.now()}`;
    try {
      const st = readStateForSolidify();
      if (st && st.last_run && st.last_run.run_id) runId = String(st.last_run.run_id);
    } catch (e) {}
    let artifact = null;
    try {
      artifact = writePromptArtifact({
        memoryDir: getEvolutionDir(),
        cycleId,
        runId,
        prompt,
        meta: {
          agent: AGENT_NAME,
          drift_enabled: IS_RANDOM_DRIFT,
          review_mode: IS_REVIEW_MODE,
          dry_run: IS_DRY_RUN,
          mutation_id: mutation && mutation.id ? mutation.id : null,
          personality_key: personalitySelection && personalitySelection.personality_key ? personalitySelection.personality_key : null,
        },
      });
    } catch (e) {
      artifact = null;
    }

    const executorTask = [
      'You are the executor (the Hand).',
      'Your job is to apply a safe, minimal patch in this repo following the attached GEP protocol prompt.',
      artifact && artifact.promptPath ? `Prompt file: ${artifact.promptPath}` : 'Prompt file: (unavailable)',
      '',
      'After applying changes and validations, you MUST run:',
      '  node index.js solidify',
      '',
      'Loop chaining (only if you are running in loop mode): after solidify succeeds, print a sessions_spawn call to start the next loop run with a short delay.',
      'Example:',
      'sessions_spawn({ task: "exec: node skills/feishu-evolver-wrapper/lifecycle.js ensure", agentId: "main", cleanup: "delete", label: "gep_loop_next" })',
      '',
      'GEP protocol prompt (may be truncated here; prefer the prompt file if provided):',
      clip(prompt, 24000),
    ].join('\n');

    const spawn = renderSessionsSpawnCall({
      task: executorTask,
      agentId: AGENT_NAME,
      cleanup: 'delete',
      label: `gep_bridge_${cycleNum}`,
    });

    console.log('\n[BRIDGE 已启用] 通过 sessions_spawn 生成执行器代理。');
    console.log(spawn);
    if (printPrompt) {
      console.log('\n[提示输出] (EVOLVE_PRINT_PROMPT=true)');
      console.log(prompt);
    }
  } else {
    console.log(prompt);
    console.log('\n[需要固化] 应用补丁并验证后，运行: node index.js solidify');
  }
}

module.exports = { run, computeAdaptiveStrategyPolicy, shouldSkipHubCalls, verbose, determineBridgeEnabled };

