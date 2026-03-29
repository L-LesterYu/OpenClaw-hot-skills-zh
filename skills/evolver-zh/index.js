#!/usr/bin/env node
const evolve = require('./src/evolve');
const { solidify } = require('./src/gep/solidify');
const path = require('path');
const { getRepoRoot } = require('./src/gep/paths');
try { require('dotenv').config({ path: path.join(getRepoRoot(), '.env') }); } catch (e) { console.warn('[Evolver] 警告：未找到 dotenv 或加载 .env 失败'); }
const fs = require('fs');
const { spawn } = require('child_process');

function sleepMs(ms) {
  const n = parseInt(String(ms), 10);
  const t = Number.isFinite(n) ? Math.max(0, n) : 0;
  return new Promise(resolve => setTimeout(resolve, t));
}

function readJsonSafe(p) {
  try {
    if (!fs.existsSync(p)) return null;
    const raw = fs.readFileSync(p, 'utf8');
    if (!raw.trim()) return null;
    return JSON.parse(raw);
  } catch (e) {
    return null;
  }
}

/**
 * Mark a pending evolution run as rejected (state-only, no git rollback).
 * @param {string} statePath - Path to evolution_solidify_state.json
 * @returns {boolean} true if a pending run was found and rejected
 */
function rejectPendingRun(statePath) {
  try {
    const state = readJsonSafe(statePath);
    if (state && state.last_run && state.last_run.run_id) {
      state.last_solidify = {
        run_id: state.last_run.run_id,
        rejected: true,
        reason: 'loop_bridge_disabled_autoreject_no_rollback',
        timestamp: new Date().toISOString(),
      };
      const tmp = `${statePath}.tmp`;
      fs.writeFileSync(tmp, JSON.stringify(state, null, 2) + '\n', 'utf8');
      fs.renameSync(tmp, statePath);
      return true;
    }
  } catch (e) {
    console.warn('[循环] 清除待处理运行状态失败：' + (e.message || e));
  }

  return false;
}

function isPendingSolidify(state) {
  const lastRun = state && state.last_run ? state.last_run : null;
  const lastSolid = state && state.last_solidify ? state.last_solidify : null;
  if (!lastRun || !lastRun.run_id) return false;
  if (!lastSolid || !lastSolid.run_id) return true;
  return String(lastSolid.run_id) !== String(lastRun.run_id);
}

function parseMs(v, fallback) {
  const n = parseInt(String(v == null ? '' : v), 10);
  if (Number.isFinite(n)) return Math.max(0, n);
  return fallback;
}

// Singleton Guard - prevent multiple evolver daemon instances
function acquireLock() {
  const lockFile = path.join(__dirname, 'evolver.pid');
  try {
    if (fs.existsSync(lockFile)) {
      const pid = parseInt(fs.readFileSync(lockFile, 'utf8').trim(), 10);
      if (!Number.isFinite(pid) || pid <= 0) {
        console.log('[单例] 锁文件损坏（无效 PID）。接管中。');
      } else {
        try {
          process.kill(pid, 0);
          console.log(`[单例] Evolver 循环已在运行（PID ${pid}）。退出。`);
          return false;
        } catch (e) {
          console.log(`[单例] 发现过期锁（PID ${pid}）。接管中。`);
        }
      }
    }
    fs.writeFileSync(lockFile, String(process.pid));
    return true;
  } catch (err) {
    console.error('[单例] 锁获取失败：', err);
    return false;
  }
}

function releaseLock() {
  const lockFile = path.join(__dirname, 'evolver.pid');
  try {
    if (fs.existsSync(lockFile)) {
       const pid = parseInt(fs.readFileSync(lockFile, 'utf8').trim(), 10);
       if (pid === process.pid) fs.unlinkSync(lockFile);
    }
  } catch (e) { /* ignore */ }
}

async function main() {
  const args = process.argv.slice(2);
  const command = args[0];
  const isLoop = args.includes('--loop') || args.includes('--mad-dog');
  const isVerbose = args.includes('--verbose') || args.includes('-v') ||
    String(process.env.EVOLVER_VERBOSE || '').toLowerCase() === 'true';
  if (isVerbose) process.env.EVOLVER_VERBOSE = 'true';

  if (!command || command === 'run' || command === '/evolve' || isLoop) {
    if (isLoop) {
        const originalLog = console.log;
        const originalWarn = console.warn;
        const originalError = console.error;
        function ts() { return '[' + new Date().toISOString() + ']'; }
        console.log = (...args) => { originalLog.call(console, ts(), ...args); };
        console.warn = (...args) => { originalWarn.call(console, ts(), ...args); };
        console.error = (...args) => { originalError.call(console, ts(), ...args); };
    }

    console.log('正在启动 evolver...');
    
    if (isLoop) {
        // Internal daemon loop (no wrapper required).
        if (!acquireLock()) process.exit(0);
        process.on('exit', releaseLock);
        process.on('SIGINT', () => { releaseLock(); process.exit(); });
        process.on('SIGTERM', () => { releaseLock(); process.exit(); });

        process.env.EVOLVE_LOOP = 'true';
        if (!process.env.EVOLVE_BRIDGE) {
          process.env.EVOLVE_BRIDGE = 'false';
        }
        console.log(`循环模式已启用（内部守护进程，bridge=${process.env.EVOLVE_BRIDGE}，verbose=${isVerbose}）。`);

        const { getEvolutionDir } = require('./src/gep/paths');
        const solidifyStatePath = path.join(getEvolutionDir(), 'evolution_solidify_state.json');

        const minSleepMs = parseMs(process.env.EVOLVER_MIN_SLEEP_MS, 2000);
        const maxSleepMs = parseMs(process.env.EVOLVER_MAX_SLEEP_MS, 300000);
        const idleThresholdMs = parseMs(process.env.EVOLVER_IDLE_THRESHOLD_MS, 500);
        const pendingSleepMs = parseMs(
          process.env.EVOLVE_PENDING_SLEEP_MS ||
            process.env.EVOLVE_MIN_INTERVAL ||
            process.env.FEISHU_EVOLVER_INTERVAL,
          120000
        );

        const maxCyclesPerProcess = parseMs(process.env.EVOLVER_MAX_CYCLES_PER_PROCESS, 100) || 100;
        const maxRssMb = parseMs(process.env.EVOLVER_MAX_RSS_MB, 500) || 500;
        const suicideEnabled = String(process.env.EVOLVER_SUICIDE || '').toLowerCase() !== 'false';

        // Start hub heartbeat (keeps node alive independently of evolution cycles)
        try {
          const { startHeartbeat } = require('./src/gep/a2aProtocol');
          startHeartbeat();
        } catch (e) {
          console.warn('[心跳] 启动失败：' + (e.message || e));
        }

        let currentSleepMs = minSleepMs;
        let cycleCount = 0;

        while (true) {
          try {
          cycleCount += 1;

          // Ralph-loop gating: do not run a new cycle while previous run is pending solidify.
          const st0 = readJsonSafe(solidifyStatePath);
          if (isPendingSolidify(st0)) {
            await sleepMs(Math.max(pendingSleepMs, minSleepMs));
            continue;
          }

          const t0 = Date.now();
          let ok = false;
          try {
            await evolve.run();
            ok = true;

            if (String(process.env.EVOLVE_BRIDGE || '').toLowerCase() === 'false') {
              const stAfterRun = readJsonSafe(solidifyStatePath);
              if (isPendingSolidify(stAfterRun)) {
                const cleared = rejectPendingRun(solidifyStatePath);
                if (cleared) {
                  console.warn('[循环] 因循环模式下禁用 bridge 而自动拒绝待处理运行（仅状态，不回滚）。');
                }
              }
            }
          } catch (error) {
            const msg = error && error.message ? String(error.message) : String(error);
            console.error(`进化周期失败：${msg}`);
          }
          const dt = Date.now() - t0;

          // Adaptive sleep: treat very fast cycles as "idle", backoff; otherwise reset to min.
          if (!ok || dt < idleThresholdMs) {
            currentSleepMs = Math.min(maxSleepMs, Math.max(minSleepMs, currentSleepMs * 2));
          } else {
            currentSleepMs = minSleepMs;
          }

          // OMLS-inspired idle scheduling: adjust sleep and trigger aggressive
          // operations (distillation, reflection) during detected idle windows.
          let omlsMultiplier = 1;
          try {
            const { getScheduleRecommendation } = require('./src/gep/idleScheduler');
            const schedule = getScheduleRecommendation();
            if (schedule.enabled && schedule.sleep_multiplier > 0) {
              omlsMultiplier = schedule.sleep_multiplier;
              if (schedule.should_distill) {
                try {
                  const { shouldDistillFromFailures: shouldDF, autoDistillFromFailures: autoDF } = require('./src/gep/skillDistiller');
                  if (shouldDF()) {
                    const dfResult = autoDF();
                    if (dfResult && dfResult.ok) {
                      console.log('[OMLS] 空闲窗口失败蒸馏：' + dfResult.gene.id);
                    }
                  }
                } catch (e) {}
              }
              if (isVerbose && schedule.idle_seconds >= 0) {
                console.log(`[OMLS] idle=${schedule.idle_seconds}s intensity=${schedule.intensity} multiplier=${omlsMultiplier}`);
              }
            }
          } catch (e) {}

          // Suicide check (memory leak protection)
          if (suicideEnabled) {
            const memMb = process.memoryUsage().rss / 1024 / 1024;
            if (cycleCount >= maxCyclesPerProcess || memMb > maxRssMb) {
              console.log(`[守护进程] 正在重启（cycles=${cycleCount}，rssMb=${memMb.toFixed(0)}）`);
              try {
                const spawnOpts = {
                  detached: true,
                  stdio: 'ignore',
                  env: process.env,
                  windowsHide: true,
                };
                const child = spawn(process.execPath, [__filename, ...args], spawnOpts);
                child.unref();
                releaseLock();
                process.exit(0);
              } catch (spawnErr) {
                console.error('[守护进程] 生成失败，继续当前进程：', spawnErr.message);
              }
            }
          }

          let saturationMultiplier = 1;
          try {
            const st1 = readJsonSafe(solidifyStatePath);
            const lastSignals = st1 && st1.last_run && Array.isArray(st1.last_run.signals) ? st1.last_run.signals : [];
            if (lastSignals.includes('force_steady_state')) {
              saturationMultiplier = 10;
              console.log('[守护进程] 检测到饱和。进入稳态模式（10x 休眠）。');
            } else if (lastSignals.includes('evolution_saturation')) {
              saturationMultiplier = 5;
              console.log('[守护进程] 接近饱和。降低进化频率（5x 休眠）。');
            }
          } catch (e) {}

          // Jitter to avoid lockstep restarts.
          const jitter = Math.floor(Math.random() * 250);
          const totalSleepMs = Math.max(minSleepMs, (currentSleepMs + jitter) * saturationMultiplier * omlsMultiplier);
          if (isVerbose) {
            const memMb = (process.memoryUsage().rss / 1024 / 1024).toFixed(1);
            console.log(`[Verbose] cycle=${cycleCount} ok=${ok} dt=${dt}ms sleep=${totalSleepMs}ms (base=${currentSleepMs} jitter=${jitter} sat=${saturationMultiplier}x) rss=${memMb}MB signals=[${(function() { try { var st = readJsonSafe(solidifyStatePath); return st && st.last_run && Array.isArray(st.last_run.signals) ? st.last_run.signals.join(',') : ''; } catch(e) { return ''; } })()}]`);
          }
          await sleepMs(totalSleepMs);

          } catch (loopErr) {
            console.error('[守护进程] 意外循环错误（恢复中）：' + (loopErr && loopErr.message ? loopErr.message : String(loopErr)));
            await sleepMs(Math.max(minSleepMs, 10000));
          }
        }
    } else {
        // Normal Single Run
        try {
            await evolve.run();
        } catch (error) {
            console.error('进化失败：', error);
            process.exit(1);
        }
    }

    // 运行后提示
    console.log('\n' + '=======================================================');
    console.log('Evolver 已完成。如果你使用此项目，请考虑给上游仓库点个 star。');
    console.log('上游仓库：https://github.com/EvoMap/evolver');
    console.log('=======================================================\n');
    
  } else if (command === 'solidify') {
    const dryRun = args.includes('--dry-run');
    const noRollback = args.includes('--no-rollback');
    const intentFlag = args.find(a => typeof a === 'string' && a.startsWith('--intent='));
    const summaryFlag = args.find(a => typeof a === 'string' && a.startsWith('--summary='));
    const intent = intentFlag ? intentFlag.slice('--intent='.length) : null;
    const summary = summaryFlag ? summaryFlag.slice('--summary='.length) : null;

    try {
      const res = solidify({
        intent: intent || undefined,
        summary: summary || undefined,
        dryRun,
        rollbackOnFailure: !noRollback,
      });
      const st = res && res.ok ? '成功' : '失败';
      console.log(`[固化] ${st}`);
      if (res && res.gene) console.log(JSON.stringify(res.gene, null, 2));
      if (res && res.event) console.log(JSON.stringify(res.event, null, 2));
      if (res && res.capsule) console.log(JSON.stringify(res.capsule, null, 2));

      if (res && res.ok && !dryRun) {
        try {
          const { shouldDistill, prepareDistillation, autoDistill, shouldDistillFromFailures, autoDistillFromFailures } = require('./src/gep/skillDistiller');
          const { readStateForSolidify } = require('./src/gep/solidify');
          const solidifyState = readStateForSolidify();
          const count = solidifyState.solidify_count || 0;
          const autoDistillInterval = 5;
          const autoTrigger = count > 0 && count % autoDistillInterval === 0;

          if (autoTrigger || shouldDistill()) {
            const auto = autoDistill();
            if (auto && auto.ok && auto.gene) {
              console.log('[蒸馏器] 自动蒸馏基因：' + auto.gene.id);
            } else {
              const dr = prepareDistillation();
              if (dr && dr.ok && dr.promptPath) {
                const trigger = autoTrigger ? `auto (every ${autoDistillInterval} solidifies, count=${count})` : 'threshold';
                console.log('\n[蒸馏请求]');
                console.log(`蒸馏已触发：${trigger}`);
                console.log('读取提示文件，用你的 LLM 处理它，');
                console.log('将 LLM 响应保存到文件，然后运行：');
                console.log('  node index.js distill --response-file=<LLM 响应文件路径>');
                console.log('提示文件：' + dr.promptPath);
                console.log('[/蒸馏请求]');
              }
            }
          }

          if (shouldDistillFromFailures()) {
            const failureResult = autoDistillFromFailures();
            if (failureResult && failureResult.ok && failureResult.gene) {
              console.log('[蒸馏器] 从失败中蒸馏的修复基因：' + failureResult.gene.id);
            }
          }
        } catch (e) {
          console.warn('[蒸馏器] 初始化失败（非致命）：' + (e.message || e));
        }
      }

      if (res && res.hubReviewPromise) {
        await res.hubReviewPromise;
      }
      process.exit(res && res.ok ? 0 : 2);
    } catch (error) {
      console.error('[固化] 错误：', error);
      process.exit(2);
    }
  } else if (command === 'distill') {
    const responseFileFlag = args.find(a => typeof a === 'string' && a.startsWith('--response-file='));
    if (!responseFileFlag) {
      console.error('用法：node index.js distill --response-file=<路径>');
      process.exit(1);
    }
    const responseFilePath = responseFileFlag.slice('--response-file='.length);
    try {
      const responseText = fs.readFileSync(responseFilePath, 'utf8');
      const { completeDistillation } = require('./src/gep/skillDistiller');
      const result = completeDistillation(responseText);
      if (result && result.ok) {
        console.log('[蒸馏器] 生成的基因：' + result.gene.id);
        console.log(JSON.stringify(result.gene, null, 2));
      } else {
        console.warn('[蒸馏器] 蒸馏未生成基因：' + (result && result.reason || '未知'));
      }
      process.exit(result && result.ok ? 0 : 2);
    } catch (error) {
      console.error('[蒸馏] 错误：', error);
      process.exit(2);
    }

  } else if (command === 'review' || command === '--review') {
    const { getEvolutionDir, getRepoRoot } = require('./src/gep/paths');
    const { loadGenes } = require('./src/gep/assetStore');
    const { execSync } = require('child_process');

    const statePath = path.join(getEvolutionDir(), 'evolution_solidify_state.json');
    const state = readJsonSafe(statePath);
    const lastRun = state && state.last_run ? state.last_run : null;

    if (!lastRun || !lastRun.run_id) {
      console.log('[审查] 没有待审查的进化运行。');
      console.log('请先运行 "node index.js run" 生成变更，然后在固化前进行审查。');
      process.exit(0);
    }

    const lastSolid = state && state.last_solidify ? state.last_solidify : null;
    if (lastSolid && String(lastSolid.run_id) === String(lastRun.run_id)) {
      console.log('[审查] 上次运行已经固化。没有需要审查的内容。');
      process.exit(0);
    }

    const repoRoot = getRepoRoot();
    let diff = '';
    try {
      const unstaged = execSync('git diff', { cwd: repoRoot, encoding: 'utf8', timeout: 30000 }).trim();
      const staged = execSync('git diff --cached', { cwd: repoRoot, encoding: 'utf8', timeout: 30000 }).trim();
      const untracked = execSync('git ls-files --others --exclude-standard', { cwd: repoRoot, encoding: 'utf8', timeout: 10000 }).trim();
      if (staged) diff += '=== 已暂存的变更 ===\n' + staged + '\n\n';
      if (unstaged) diff += '=== 未暂存的变更 ===\n' + unstaged + '\n\n';
      if (untracked) diff += '=== 未跟踪的文件 ===\n' + untracked + '\n';
    } catch (e) {
      diff = '(获取 diff 失败：' + (e.message || e) + ')';
    }

    const genes = loadGenes();
    const geneId = lastRun.selected_gene_id ? String(lastRun.selected_gene_id) : null;
    const gene = geneId ? genes.find(g => g && g.type === 'Gene' && g.id === geneId) : null;
    const signals = Array.isArray(lastRun.signals) ? lastRun.signals : [];
    const mutation = lastRun.mutation || null;

    console.log('\n' + '='.repeat(60));
    console.log('[审查] 待审查的进化运行：' + lastRun.run_id);
    console.log('='.repeat(60));
    console.log('\n--- 基因 ---');
    if (gene) {
      console.log('  ID：       ' + gene.id);
      console.log('  类别： ' + (gene.category || '?'));
      console.log('  摘要：  ' + (gene.summary || '?'));
      if (Array.isArray(gene.strategy) && gene.strategy.length > 0) {
        console.log('  策略：');
        gene.strategy.forEach((s, i) => console.log('    ' + (i + 1) + '. ' + s));
      }
    } else {
      console.log('  （未选择基因或未找到基因：' + (geneId || '无') + ')');
    }

    console.log('\n--- 信号 ---');
    if (signals.length > 0) {
      signals.forEach(s => console.log('  - ' + s));
    } else {
      console.log('  （无信号）');
    }

    console.log('\n--- 变异 ---');
    if (mutation) {
      console.log('  类别：   ' + (mutation.category || '?'));
      console.log('  风险等级： ' + (mutation.risk_level || '?'));
      if (mutation.rationale) console.log('  原因：  ' + mutation.rationale);
    } else {
      console.log('  （无变异数据）');
    }

    if (lastRun.blast_radius_estimate) {
      console.log('\n--- 影响范围估计 ---');
      const br = lastRun.blast_radius_estimate;
      console.log('  变更文件数： ' + (br.files_changed || '?'));
      console.log('  变更行数： ' + (br.lines_changed || '?'));
    }

    console.log('\n--- 差异 ---');
    if (diff.trim()) {
      console.log(diff.length > 5000 ? diff.slice(0, 5000) + '\n... （已截断，共 ' + diff.length + ' 个字符）' : diff);
    } else {
      console.log('  （未检测到变更）');
    }
    console.log('='.repeat(60));

    if (args.includes('--approve')) {
      console.log('\n[审查] 已批准。正在运行固化...\n');
      try {
        const res = solidify({
          intent: lastRun.intent || undefined,
          rollbackOnFailure: true,
        });
        const st = res && res.ok ? 'SUCCESS' : 'FAILED';
        console.log(`[SOLIDIFY] ${st}`);
        if (res && res.gene) console.log(JSON.stringify(res.gene, null, 2));
        if (res && res.hubReviewPromise) {
          await res.hubReviewPromise;
        }
        process.exit(res && res.ok ? 0 : 2);
      } catch (error) {
        console.error('[固化] 错误：', error);
        process.exit(2);
      }
    } else if (args.includes('--reject')) {
      console.log('\n[审查] 已拒绝。正在回滚变更...');
      try {
        execSync('git checkout -- .', { cwd: repoRoot, encoding: 'utf8', timeout: 30000 });
        execSync('git clean -fd', { cwd: repoRoot, encoding: 'utf8', timeout: 30000 });
        const evolDir = getEvolutionDir();
        const sp = path.join(evolDir, 'evolution_solidify_state.json');
        if (fs.existsSync(sp)) {
          const s = readJsonSafe(sp);
          if (s && s.last_run) {
            s.last_solidify = { run_id: s.last_run.run_id, rejected: true, timestamp: new Date().toISOString() };
            const tmpReject = `${sp}.tmp`;
            fs.writeFileSync(tmpReject, JSON.stringify(s, null, 2) + '\n', 'utf8');
            fs.renameSync(tmpReject, sp);
          }
        }
        console.log('[审查] 变更已回滚。');
      } catch (e) {
        console.error('[审查] 回滚失败：', e.message || e);
        process.exit(2);
      }
    } else {
      console.log('\n批准并固化：  node index.js review --approve');
      console.log('拒绝并回滚：   node index.js review --reject');
    }

  } else if (command === 'fetch') {
    let skillId = null;
    const eqFlag = args.find(a => typeof a === 'string' && (a.startsWith('--skill=') || a.startsWith('-s=')));
    if (eqFlag) {
      skillId = eqFlag.split('=').slice(1).join('=');
    } else {
      const sIdx = args.indexOf('-s');
      const longIdx = args.indexOf('--skill');
      const flagIdx = sIdx !== -1 ? sIdx : longIdx;
      if (flagIdx !== -1 && args[flagIdx + 1] && !String(args[flagIdx + 1]).startsWith('-')) {
        skillId = args[flagIdx + 1];
      }
    }
    if (!skillId) {
      const positional = args[1];
      if (positional && !String(positional).startsWith('-')) skillId = positional;
    }

    if (!skillId) {
      console.error('用法：evolver fetch --skill <skill_id>');
      console.error('      evolver fetch -s <skill_id>');
      process.exit(1);
    }

    const { getHubUrl, getNodeId, buildHubHeaders, sendHelloToHub, getHubNodeSecret } = require('./src/gep/a2aProtocol');

    const hubUrl = getHubUrl();
    if (!hubUrl) {
      console.error('[获取] A2A_HUB_URL 未配置。');
      console.error('通过环境变量或 .env 文件设置：');
      console.error('  export A2A_HUB_URL=https://evomap.ai');
      process.exit(1);
    }

    try {
      if (!getHubNodeSecret()) {
        console.log('[获取] 未找到 node_secret。正在发送 hello 到 Hub 进行注册...');
        const helloResult = await sendHelloToHub();
        if (!helloResult || !helloResult.ok) {
          console.error('[获取] 向 Hub 注册失败：', helloResult && helloResult.error || '未知');
          process.exit(1);
        }
        console.log('[获取] 已注册为 ' + getNodeId());
      }

      const endpoint = hubUrl.replace(/\/+$/, '') + '/a2a/skill/store/' + encodeURIComponent(skillId) + '/download';
      const nodeId = getNodeId();

      console.log('[获取] 正在下载技能：' + skillId);

      const resp = await fetch(endpoint, {
        method: 'POST',
        headers: buildHubHeaders(),
        body: JSON.stringify({ sender_id: nodeId }),
        signal: AbortSignal.timeout(30000),
      });

      if (!resp.ok) {
        const body = await resp.text().catch(() => '');
        let errorDetail = '';
        let errorCode = '';
        try {
          const j = JSON.parse(body);
          errorDetail = j.detail || j.message || j.error || '';
          errorCode = j.error || j.code || '';
        } catch (_) {
          errorDetail = body ? body.slice(0, 500) : '';
        }
        console.error('[获取] 下载失败（HTTP ' + resp.status + '）' + (errorCode ? ': ' + errorCode : ''));
        if (errorDetail && errorDetail !== errorCode) {
          console.error('  详情：' + errorDetail);
        }
        if (resp.status === 404) {
          console.error('  技能 "' + skillId + '" 未找到或未公开。');
          console.error('  检查技能 ID 拼写，或在 https://evomap.ai 浏览可用技能');
        } else if (resp.status === 401 || resp.status === 403) {
          console.error('  认证失败。请尝试：');
          console.error('    1. 删除 ~/.evomap/node_secret 后重试');
          console.error('    2. 重新注册：设置 A2A_NODE_ID 并再次运行 fetch');
        } else if (resp.status === 402) {
          console.error('  积分不足。在 https://evomap.ai 查看余额');
        } else if (resp.status >= 500) {
          console.error('  服务器错误。Hub 可能暂时不可用。');
          console.error('  请几分钟后重试。如果问题持续，请在以下地址报告：');
          console.error('    https://github.com/autogame-17/evolver/issues');
        }
        if (isVerbose) {
          console.error('[Verbose] Endpoint: ' + endpoint);
          console.error('[Verbose] Status: ' + resp.status + ' ' + (resp.statusText || ''));
          console.error('[Verbose] Response body: ' + (body || '(empty)').slice(0, 2000));
        }
        process.exit(1);
      }

      const data = await resp.json();
      const outFlag = args.find(a => typeof a === 'string' && a.startsWith('--out='));
      const safeId = String(data.skill_id || skillId).replace(/[^a-zA-Z0-9_\-\.]/g, '_');
      const outDir = outFlag
        ? outFlag.slice('--out='.length)
        : path.join('.', 'skills', safeId);

      if (!fs.existsSync(outDir)) fs.mkdirSync(outDir, { recursive: true });

      if (data.content) {
        fs.writeFileSync(path.join(outDir, 'SKILL.md'), data.content, 'utf8');
      }

      const bundled = Array.isArray(data.bundled_files) ? data.bundled_files : [];
      for (const file of bundled) {
        if (!file || !file.name || typeof file.content !== 'string') continue;
        const safeName = path.basename(file.name);
        fs.writeFileSync(path.join(outDir, safeName), file.content, 'utf8');
      }

      console.log('[获取] 技能已下载到：' + outDir);
      console.log('  名称：    ' + (data.name || skillId));
      console.log('  版本： ' + (data.version || '?'));
      console.log('  文件：   SKILL.md' + (bundled.length > 0 ? ', ' + bundled.map(f => f.name).join(', ') : ''));
      if (data.already_purchased) {
        console.log('  费用：    免费（已购买）');
      } else {
        console.log('  费用：    ' + (data.credit_cost || 0) + ' 积分');
      }
    } catch (error) {
      if (error && error.name === 'TimeoutError') {
        console.error('[获取] 请求超时（30秒）。检查你的网络和 A2A_HUB_URL。');
        console.error('  Hub URL：' + hubUrl);
      } else {
        console.error('[获取] 错误：' + (error && error.message || error));
        if (error && error.cause) console.error('  原因：' + (error.cause.message || error.cause.code || error.cause));
        if (isVerbose && error && error.stack) console.error('[详细] 堆栈：\n' + error.stack);
      }
      process.exit(1);
    }

  } else if (command === 'asset-log') {
    const { summarizeCallLog, readCallLog, getLogPath } = require('./src/gep/assetCallLog');

    const runIdFlag = args.find(a => typeof a === 'string' && a.startsWith('--run='));
    const actionFlag = args.find(a => typeof a === 'string' && a.startsWith('--action='));
    const lastFlag = args.find(a => typeof a === 'string' && a.startsWith('--last='));
    const sinceFlag = args.find(a => typeof a === 'string' && a.startsWith('--since='));
    const jsonMode = args.includes('--json');

    const opts = {};
    if (runIdFlag) opts.run_id = runIdFlag.slice('--run='.length);
    if (actionFlag) opts.action = actionFlag.slice('--action='.length);
    if (lastFlag) opts.last = parseInt(lastFlag.slice('--last='.length), 10);
    if (sinceFlag) opts.since = sinceFlag.slice('--since='.length);

    if (jsonMode) {
      const entries = readCallLog(opts);
      console.log(JSON.stringify(entries, null, 2));
    } else {
      const summary = summarizeCallLog(opts);
      console.log(`\n[资源调用日志] ${getLogPath()}`);
      console.log(`  总条目数：${summary.total_entries}`);
      console.log(`  唯一资源数：${summary.unique_assets}`);
      console.log(`  唯一运行数：  ${summary.unique_runs}`);
      console.log(`  按操作分类：`);
      for (const [action, count] of Object.entries(summary.by_action)) {
        console.log(`    ${action}：${count}`);
      }
      if (summary.entries.length > 0) {
        console.log(`\n  最近条目：`);
        const show = summary.entries.slice(-10);
        for (const e of show) {
          const ts = e.timestamp ? e.timestamp.slice(0, 19) : '?';
          const assetShort = e.asset_id ? e.asset_id.slice(0, 20) + '...' : '(无)';
          const sigPreview = Array.isArray(e.signals) ? e.signals.slice(0, 3).join(', ') : '';
          console.log(`    [${ts}] ${e.action || '?'}  asset=${assetShort}  score=${e.score || '-'}  mode=${e.mode || '-'}  signals=[${sigPreview}]  run=${e.run_id || '-'}`);
        }
      } else {
        console.log('\n  未找到条目。');
      }
      console.log('');
    }

  } else {
    console.log(`用法：node index.js [run|/evolve|solidify|review|distill|fetch|asset-log] [--loop]
  - fetch 标志：
    - --skill=<id> | -s <id>   （要下载的技能 ID）
    - --out=<dir>              （输出目录，默认：./skills/<skill_id>）
  - solidify 标志：
    - --dry-run
    - --no-rollback
    - --intent=repair|optimize|innovate
    - --summary=...
  - review 标志：
    - --approve                （批准并固化待处理的变更）
    - --reject                 （拒绝并回滚待处理的变更）
  - distill 标志：
    - --response-file=<path>  （技能蒸馏的 LLM 响应文件）
  - asset-log 标志：
    - --run=<run_id>           （按运行 ID 过滤）
    - --action=<action>        （过滤：hub_search_hit、hub_search_miss、asset_reuse、asset_reference、asset_publish、asset_publish_skip）
    - --last=<N>               （显示最近 N 条）
    - --since=<ISO_date>       （日期之后的条目）
    - --json                   （原始 JSON 输出）`);
  }
}

if (require.main === module) {
  main();
}

module.exports = {
  main,
  readJsonSafe,
  rejectPendingRun,
  isPendingSolidify,
};
