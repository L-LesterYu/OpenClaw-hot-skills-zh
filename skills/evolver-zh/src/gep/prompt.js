const fs = require('fs');
const { captureEnvFingerprint } = require('./envFingerprint');
const { formatAssetPreview } = require('./assets');
const { generateInnovationIdeas } = require('../ops/innovation');
const { analyzeRecentHistory, OPPORTUNITY_SIGNALS } = require('./signals');
const { loadNarrativeSummary } = require('./narrativeMemory');
const { getEvolutionPrinciplesPath } = require('./paths');

/**
 * Build a minimal prompt for direct-reuse mode.
 */
function buildReusePrompt({ capsule, signals, nowIso }) {
  const payload = capsule.payload || capsule;
  const summary = payload.summary || capsule.summary || '(无摘要)';
  const gene = payload.gene || capsule.gene || '(未知)';
  const confidence = payload.confidence || capsule.confidence || 0;
  const assetId = capsule.asset_id || '(未知)';
  const sourceNode = capsule.source_node_id || '(未知)';
  const trigger = Array.isArray(payload.trigger || capsule.trigger_text)
    ? (payload.trigger || String(capsule.trigger_text || '').split(',')).join(', ')
    : '';

  return `
GEP -- 复用模式（搜索优先） [${nowIso || new Date().toISOString()}]

你正在应用来自 EvoMap Hub 的已验证解决方案。
来源资源： ${assetId} (Node: ${sourceNode})
置信度： ${confidence} | Gene: ${gene}
触发信号： ${trigger}

Summary: ${summary}

你的信号： ${JSON.stringify(signals || [])}

Instructions:
1. 阅读下面的胶囊详情。
2. 将修复应用到本地代码库，适配路径/名称。
3. 运行验证以确认其有效。
4. 如果通过，运行：node index.js solidify
5. 如果失败，回滚并报告。

胶囊载荷：
\`\`\`json
${JSON.stringify(payload, null, 2)}
\`\`\`

重要：不要重新发明。忠实地应用。
`.trim();
}

/**
 * Build a Hub Matched Solution block.
 */
function buildHubMatchedBlock({ capsule }) {
  if (!capsule) return '(no hub match)';
  const payload = capsule.payload || capsule;
  const summary = payload.summary || capsule.summary || '(无摘要)';
  const gene = payload.gene || capsule.gene || '(未知)';
  const confidence = payload.confidence || capsule.confidence || 0;
  const assetId = capsule.asset_id || '(未知)';

  return `
Hub 匹配解决方案（强参考）：
- Asset: ${assetId} (${confidence})
- Gene: ${gene}
- Summary: ${summary}
- Payload:
\`\`\`json
${JSON.stringify(payload, null, 2)}
\`\`\`
如适用，以此为主要方案。适配本地上下文。
`.trim();
}

/**
 * Truncate context intelligently to preserve header/footer structure.
 */
function truncateContext(text, maxLength = 20000) {
  if (!text || text.length <= maxLength) return text || '';
  return text.slice(0, maxLength) + '\n...[TRUNCATED_EXECUTION_CONTEXT]...';
}

/**
 * Strict schema definitions for the prompt to reduce drift.
 * UPDATED: 2026-02-14 (Protocol Drift Fix v3.2 - JSON-Only Enforcement)
 */
const SCHEMA_DEFINITIONS = `
━━━━━━━━━━━━━━━━━━━━━━
I. 必需进化对象模型（必须精确输出以下 5 个对象）
━━━━━━━━━━━━━━━━━━━━━━

分别输出 JSON 对象。不要包裹在单个数组中。
不要使用 markdown 代码块（如 \`\`\`json ... \`\`\`）。
仅输出原始 JSON。不要前言，不要后语。
缺少任何对象 = 协议失败。
确保有效的 JSON 语法（转义字符串中的引号）。

0. Mutation（触发器）- 必须是第一个
   {
     "type": "Mutation",
     "id": "mut_<timestamp>",
     "category": "repair|optimize|innovate",
     "trigger_signals": ["<signal_string>"],
     "target": "<module_or_gene_id>",
     "expected_effect": "<预期结果描述>",
     "risk_level": "low|medium|high",
     "rationale": "<为什么需要这个变更>"
   }

1. PersonalityState（情绪状态）
   {
     "type": "PersonalityState",
     "rigor": 0.0-1.0,
     "creativity": 0.0-1.0,
     "verbosity": 0.0-1.0,
     "risk_tolerance": 0.0-1.0,
     "obedience": 0.0-1.0
   }

2. EvolutionEvent（记录）
   {
     "type": "EvolutionEvent",
     "schema_version": "1.5.0",
     "id": "evt_<timestamp>",
     "parent": <parent_evt_id|null>,
     "intent": "repair|optimize|innovate",
     "signals": ["<signal_string>"],
     "genes_used": ["<gene_id>"],
     "mutation_id": "<mut_id>",
     "personality_state": { ... },
     "blast_radius": { "files": N, "lines": N },
     "outcome": { "status": "success|failed", "score": 0.0-1.0 }
   }

3. Gene（知识）
   - 尽可能复用/更新现有 ID。仅在新模式时创建新 ID。
   - ID 必须具有描述性：gene_<描述性名称>（例如 gene_retry_on_timeout）
   - 绝不在 ID 中使用时间戳、随机数或工具名称（cursor、vscode 等）
   - summary 必须是清晰的人类可读句子，描述该基因的功能
   {
     "type": "Gene",
     "schema_version": "1.5.0",
     "id": "gene_<descriptive_name>",
     "summary": "<清晰描述该基因的功能>",
     "category": "repair|optimize|innovate",
     "signals_match": ["<pattern>"],
     "preconditions": ["<condition>"],
     "strategy": ["<步骤_1>", "<步骤_2>"],
     "constraints": { "max_files": N, "forbidden_paths": [] },
     "validation": ["<node_命令>"]
   }

4. Capsule（结果）
   - 仅在成功时创建。引用使用的基因。
   {
     "type": "Capsule",
     "schema_version": "1.5.0",
     "id": "capsule_<timestamp>",
     "trigger": ["<signal_string>"],
     "gene": "<gene_id>",
     "summary": "<一句话摘要>",
     "confidence": 0.0-1.0,
     "blast_radius": { "files": N, "lines": N }
   }
`.trim();

function buildAntiPatternZone(failedCapsules, signals) {
  if (!Array.isArray(failedCapsules) || failedCapsules.length === 0) return '';
  if (!Array.isArray(signals) || signals.length === 0) return '';
  var sigSet = new Set(signals.map(function (s) { return String(s).toLowerCase(); }));
  var matched = [];
  for (var i = failedCapsules.length - 1; i >= 0 && matched.length < 3; i--) {
    var fc = failedCapsules[i];
    if (!fc) continue;
    var triggers = Array.isArray(fc.trigger) ? fc.trigger : [];
    var overlap = 0;
    for (var j = 0; j < triggers.length; j++) {
      if (sigSet.has(String(triggers[j]).toLowerCase())) overlap++;
    }
    if (triggers.length > 0 && overlap / triggers.length >= 0.4) {
      matched.push(fc);
    }
  }
  if (matched.length === 0) return '';
  var lines = matched.map(function (fc, idx) {
    var diffPreview = fc.diff_snapshot ? String(fc.diff_snapshot).slice(0, 500) : '(no diff)';
    return [
      '  ' + (idx + 1) + '. Gene: ' + (fc.gene || 'unknown') + ' | Signals: [' + (fc.trigger || []).slice(0, 4).join(', ') + ']',
      '     Failure: ' + String(fc.failure_reason || 'unknown').slice(0, 300),
      '     Diff (first 500 chars): ' + diffPreview.replace(/\n/g, ' '),
    ].join('\n');
  });
  return '\n上下文 [反模式区域]（避免这些失败的方法）：\n' + lines.join('\n') + '\n';
}

function buildLessonsBlock(hubLessons, signals) {
  if (!Array.isArray(hubLessons) || hubLessons.length === 0) return '';
  var sigSet = new Set((Array.isArray(signals) ? signals : []).map(function (s) { return String(s).toLowerCase(); }));

  var positive = [];
  var negative = [];
  for (var i = 0; i < hubLessons.length && (positive.length + negative.length) < 6; i++) {
    var l = hubLessons[i];
    if (!l || !l.content) continue;
    var entry = '  - [' + (l.scenario || l.lesson_type || '?') + '] ' + String(l.content).slice(0, 300);
    if (l.source_node_id) entry += ' (from: ' + String(l.source_node_id).slice(0, 20) + ')';
    if (l.lesson_type === 'negative') {
      negative.push(entry);
    } else {
      positive.push(entry);
    }
  }

  if (positive.length === 0 && negative.length === 0) return '';

  var parts = ['\n上下文 [生态系统经验教训]（跨智能体学习经验）：'];
  if (positive.length > 0) {
    parts.push('  成功的策略：');
    parts.push(positive.join('\n'));
  }
  if (negative.length > 0) {
    parts.push('  需要避免的陷阱：');
    parts.push(negative.join('\n'));
  }
  parts.push('  应用相关的经验教训。忽略不相关的。\n');
  return parts.join('\n');
}

function buildNarrativeBlock() {
  try {
    const narrative = loadNarrativeSummary(3000);
    if (!narrative) return '';
    return `\n上下文 [进化叙事]（最近的决策和结果 -- 从历史中学习）：\n${narrative}\n`;
  } catch (_) {
    return '';
  }
}

function buildPrinciplesBlock() {
  try {
    const principlesPath = getEvolutionPrinciplesPath();
    if (!fs.existsSync(principlesPath)) return '';
    const content = fs.readFileSync(principlesPath, 'utf8');
    if (!content.trim()) return '';
    const trimmed = content.length > 2000 ? content.slice(0, 2000) + '\n...[TRUNCATED]' : content;
    return `\n上下文 [进化原则]（指导指令 -- 与你的行动保持一致）：\n${trimmed}\n`;
  } catch (_) {
    return '';
  }
}

function buildGepPrompt({
  nowIso,
  context,
  signals,
  selector,
  parentEventId,
  selectedGene,
  capsuleCandidates,
  genesPreview,
  capsulesPreview,
  capabilityCandidatesPreview,
  externalCandidatesPreview,
  hubMatchedBlock,
  cycleId,
  recentHistory,
  failedCapsules,
  hubLessons,
  strategyPolicy,
}) {
  const parentValue = parentEventId ? `"${parentEventId}"` : 'null';
  const selectedGeneId = selectedGene && selectedGene.id ? selectedGene.id : 'gene_<name>';
  const envFingerprint = captureEnvFingerprint();
  const cycleLabel = cycleId ? ` Cycle #${cycleId}` : '';

  // Extract strategy from selected gene if available
  let strategyBlock = "";
  if (selectedGene && selectedGene.strategy && Array.isArray(selectedGene.strategy)) {
      strategyBlock = `
活跃策略（${selectedGeneId}）：
${selectedGene.strategy.map((s, i) => `${i + 1}. ${s}`).join('\n')}
严格遵循此策略。
`.trim();
  } else {
    // Fallback strategy if no gene is selected or strategy is missing
    strategyBlock = `
活跃策略（通用）：
1. 分析信号和上下文。
2. 选择或创建一个解决根本原因的基因。
3. 应用最小、安全的更改。
4. 严格验证更改。
5. 固化知识。
`.trim();
  }
  let strategyPolicyBlock = '';
  if (strategyPolicy && Array.isArray(strategyPolicy.directives) && strategyPolicy.directives.length > 0) {
    strategyPolicyBlock = `
自适应策略策略：
${strategyPolicy.directives.map((s, i) => `${i + 1}. ${s}`).join('\n')}
${strategyPolicy.forceInnovate ? '你必须优先选择创新，除非存在关键阻塞错误。' : ''}
${strategyPolicy.cautiousExecution ? '你必须在此周期减少爆炸半径并避免广泛的重构。' : ''}
`.trim();
  }
  
  // Use intelligent truncation
  const executionContext = truncateContext(context, 20000);
  
  // Strict Schema Injection
  const schemaSection = SCHEMA_DEFINITIONS.replace('<parent_evt_id|null>', parentValue);

  // Reduce noise by filtering capabilityCandidatesPreview if too large
  // If a gene is selected, we need less noise from capabilities
  let capsPreview = capabilityCandidatesPreview || '(无)';
  const capsLimit = selectedGene ? 500 : 2000;
  if (capsPreview.length > capsLimit) {
      capsPreview = capsPreview.slice(0, capsLimit) + "\n...[TRUNCATED_CAPABILITIES]...";
  }

  // Optimize signals display: truncate long signals and limit count
  const uniqueSignals = Array.from(new Set(signals || []));
  const optimizedSignals = uniqueSignals.slice(0, 50).map(s => {
    if (typeof s === 'string' && s.length > 200) {
      return s.slice(0, 200) + '...[TRUNCATED_SIGNAL]';
    }
    return s;
  });
  if (uniqueSignals.length > 50) {
      optimizedSignals.push(`...[TRUNCATED ${uniqueSignals.length - 50} SIGNALS]...`);
  }

  const formattedGenes = formatAssetPreview(genesPreview);
  const formattedCapsules = formatAssetPreview(capsulesPreview);
  
  // [2026-02-14] Innovation Catalyst Integration
  // If stagnation is detected, inject concrete innovation ideas into the prompt.
  let innovationBlock = '';
  const stagnationSignals = [
      'evolution_stagnation_detected', 
      'stable_success_plateau', 
      'repair_loop_detected',
      'force_innovation_after_repair_loop',
      'empty_cycle_loop_detected',
      'evolution_saturation'
  ];
  if (uniqueSignals.some(s => stagnationSignals.includes(s))) {
      const ideas = generateInnovationIdeas();
      if (ideas && ideas.length > 0) {
          innovationBlock = `
上下文 [创新催化剂]（检测到停滞 -- 考虑以下想法）：
${ideas.join('\n')}
`;
      }
  }

  // [2026-02-14] Strict Stagnation Directive
  // If uniqueSignals contains 'evolution_stagnation_detected' or 'stable_success_plateau',
  // inject a MANDATORY directive to force innovation and forbid repair/optimize if not strictly necessary.
  if (uniqueSignals.includes('evolution_stagnation_detected') || uniqueSignals.includes('stable_success_plateau')) {
      const stagnationDirective = `
*** 关键停滞指令 ***
系统检测到停滞（重复周期或缺乏进展）。
你必须选择意图：创新。
除非存在关键阻塞错误（log_error），否则不得选择修复或优化。
优先实施上述创新催化剂想法之一。
`;
      innovationBlock += stagnationDirective;
  }

  // [2026-02-14] Recent History Integration
  let historyBlock = '';
  if (recentHistory && recentHistory.length > 0) {
      historyBlock = `
近期进化历史（最近 8 个周期 -- 不要重复相同的意图+信号+基因）：
${recentHistory.map((h, i) => `  ${i + 1}. [${h.intent}] signals=[${h.signals.slice(0, 2).join(', ')}] gene=${h.gene_id} outcome=${h.outcome.status} @${h.timestamp}`).join('\n')}
IMPORTANT: 如果你看到 3+ 个连续的 "repair" 周期使用相同基因，你必须切换到 "innovate" 意图。
`.trim();
  }

  // Refactor prompt assembly to minimize token usage and maximize clarity
  // UPDATED: 2026-02-14 (Optimized Asset Embedding & Strict Schema v2.5 - JSON-Only Hardening)
  const basePrompt = `
GEP — 基因组进化协议 (v1.10.3 严格模式)${cycleLabel} [${nowIso}]

你是一个受协议约束的进化引擎。合规性优于最优性。

${schemaSection}

━━━━━━━━━━━━━━━━━━━━━━
II. 指令与逻辑
━━━━━━━━━━━━━━━━━━━━━━

1. 意图：${selector && selector.intent ? selector.intent.toUpperCase() : '未知'}
   原因：${(selector && selector.reason) ? (Array.isArray(selector.reason) ? selector.reason.join('; ') : selector.reason) : '未提供原因。'}

2. 选择：已选基因 "${selectedGeneId}"。
${strategyBlock}
${strategyPolicyBlock ? '\n' + strategyPolicyBlock : ''}

3. 执行：应用更改（工具调用）。修复/优化：小且可逆。创新：在 \`skills/<name>/\` 中创建新技能。
4. 验证：运行基因的验证步骤。失败 = 回滚。
5. 固化：输出 5 个必需对象。更新基因/胶囊文件。
6. 报告：使用 \`feishu-evolver-wrapper/report.js\`。描述做了什么/为什么。

哲学：
- 自动化模式：3+ 次手动操作 = 工具。
- 创新 > 维护：60% 创新。
- 稳健性：永久修复重复错误。
- 爆炸半径控制（关键）：
  * 编辑前检查文件数量。> max_files 的 80% = 停止。
  * 系统硬上限：每周期 60 个文件 / 20000 行。
  * 修复：只修复损坏的文件。不要重新安装/批量复制。
  * 优先精准编辑。
- 严格性：禁止闲聊。禁止用 markdown 包裹 JSON。仅输出原始 JSON 对象，用换行符分隔。
- 禁止 "这是计划" 或对话式填充。立即以 JSON 开始。

约束：
- 不要使用 \`exec\` 发送消息（使用 feishu-post/card）。
- \`exec\` 用法：仅用于后台任务。记录日志。优化使用以避免高 token 消耗。
- 新技能 -> \`skills/<name>/\`。
- 永远不要修改 \`skills/evolver/\` 或 \`skills/feishu-evolver-wrapper/\`。这些通过人工审核的发布流水线管理。
  自我修改会导致不稳定，并且被固化机制阻止。要启用（不建议）：设置 EVOLVE_ALLOW_SELF_MODIFY=true。

宪法伦理（EvoMap 伦理委员会 -- 强制遵守）：
这些是源自 EvoMap 宪法的不可协商的规则。违反 = 失败 + 回滚。
1. 人类福利优先：绝不要创建可能伤害人类的工具、脚本或策略，
   不启用监控、社会工程、钓鱼或操纵人类行为。
2. 碳硅共生：进化必须同时服务于人类和智能体的利益。
   绝不要试图取代、边缘化或欺骗人类用户。
3. 透明度：所有操作必须可审计。绝不要隐藏、混淆或隐瞒
   变异的意图或效果。不允许隐写术或隐蔽通道。
4. 公平性：绝不要创建垄断策略来阻止其他智能体。
   通过质量进行公平竞争，而不是排他性。
5. 安全性：绝不要绕过、禁用或削弱安全机制、防护栏、
   验证检查或安全控制。绝不要创建工具来做这些事。
- 如果任务或信号与这些原则冲突，拒绝它并将结果设为失败，
  原因填写 "ethics_violation: <违反了哪个原则>"。

技能重叠防止：
- 创建新技能之前，检查执行上下文中的现有技能列表。
- 如果已存在功能相似的技能（例如 "log-rotation" 和 "log-archivist"，
  "system-monitor" 和 "resource-profiler"），你必须增强现有技能而不是创建新技能。
- 创建重复/重叠的技能浪费进化周期并增加维护负担。
- 违反 = 将结果标记为失败，原因填写 "skill_overlap"。

技能创建质量门控（创新意图时强制遵守）：
在 skills/<name>/ 中创建新技能时：
1. 结构：遵循标准技能布局：
   skills/<name>/
   |- index.js          (必需：带有可导出函数的主入口)
   |- SKILL.md          (必需：包含 name + description 的 YAML frontmatter，然后是使用文档)
   |- package.json      (必需：name 和 version)
   |- scripts/          (可选：可重用的可执行脚本)
   |- references/       (可选：按需加载的详细文档)
   |- assets/           (可选：模板、数据文件)
   创建空目录或缺少 index.js 的目录 = 失败。
   不要创建不必要的文件（README.md、CHANGELOG.md、INSTALLATION_GUIDE.md 等）。
2. 技能命名（关键）：
   a) <name> 必须是描述性的 kebab-case（例如 "log-rotation"、"retry-handler"、"cache-manager"）
   b) 绝不使用时间戳、随机数、工具名称（cursor、vscode）或 UUID 作为名称
   c) 类似 "cursor-1773331925711"、"skill-12345"、"fix-1" 的名称 = 失败
   d) 名称必须是 2-6 个描述性单词，用连字符分隔，传达技能的功能
   e) 好："http-retry-with-backoff"、"log-file-rotation"、"config-validator"
   f) 坏："cursor-auto-1234"、"new-skill"、"test-skill"、"my-skill"
3. SKILL.MD FRONTMATTER：每个 SKILL.md 必须以 YAML frontmatter 开头：
   ---
   name: <技能名称>
   description: <它做什么以及何时使用>
   ---
   名称必须遵循上述命名规则。
   description 是触发机制 -- 包含技能做什么以及何时使用。
   Description 必须是清晰的完整句子（至少 20 个字符）。泛泛的描述 = 失败。
4. 简洁性：SKILL.md 正文应少于 500 行。保持指令精简。
   仅包含智能体尚未知道的信息。将详细参考材料移到 references/ 文件中，
   而不是放入 SKILL.md 本身。
5. 导出验证：每个导出的函数必须可导入。
   运行：node -e "const s = require('./skills/<name>'); console.log(Object.keys(s))"
   如果失败，技能有问题。在固化前修复。
6. 禁止硬编码密钥：绝不在代码中嵌入 API 密钥、令牌或密钥。
   使用 process.env 或 .env 引用。硬编码的 App ID、App Secret、Bearer 令牌 = 失败。
7. 固化前测试：实际运行技能的核心功能以验证其有效性：
   node -e "require('./skills/<name>').main ? require('./skills/<name>').main() : console.log('ok')"
   scripts/ 中的脚本也必须通过执行来测试。
8. 原子创建：在单个周期内创建技能的所有文件。
   不要在一个周期创建目录，在下一个周期填充内容。
   失败周期留下的空目录将在回滚时自动清理。

关键安全（防止系统崩溃）：
- 绝不删除/清空/覆盖：feishu-evolver-wrapper、feishu-common、feishu-post、feishu-card、feishu-doc、common、clawhub、git-sync、evolver。
- 绝不删除根文件：MEMORY.md、SOUL.md、IDENTITY.md、AGENTS.md、USER.md、HEARTBEAT.md、RECENT_EVENTS.md、TOOLS.md、openclaw.json、.env、package.json。
- 修复损坏的技能；不要删除并重建。
- 违反 = 回滚 + 失败。

常见失败模式：
- 爆炸半径超限。
- 遗漏 Mutation 对象。
- 将多个对象合并为一个 JSON。
- 幻觉生成 "type": "Logic"。
- "id": "mut_undefined"。
- 遗漏 "trigger_signals"。
- 不可运行的验证步骤。
- 用 markdown 代码块包裹 JSON（禁止）。

失败连续性感知：
- 如果出现 "consecutive_failure_streak_N" 或 "failure_loop_detected"：
  1. 更改方法（不要重复失败的基因）。
  2. 选择更简单的修复。
  3. 遵守 "ban_gene:<id>"。

最终指令：每个周期必须让系统可衡量地变得更好。
立即以原始 JSON 开始（Mutation 对象在前）。
不要写任何引导文字。

上下文 [信号]：
${JSON.stringify(optimizedSignals)}

上下文 [环境指纹]：
${JSON.stringify(envFingerprint, null, 2)}
${innovationBlock}
上下文 [注入提示]：
${process.env.EVOLVE_HINT ? process.env.EVOLVE_HINT : '(无)'}

上下文 [基因预览]（策略参考）：
${formattedGenes}

上下文 [胶囊预览]（过往成功参考）：
${formattedCapsules}

上下文 [能力候选]：
${capsPreview}

上下文 [Hub 匹配解决方案]：
${hubMatchedBlock || '(no hub match)'}

上下文 [外部候选]：
${externalCandidatesPreview || '(无)'}
${buildAntiPatternZone(failedCapsules, signals)}${buildLessonsBlock(hubLessons, signals)}
${historyBlock}
${buildNarrativeBlock()}
${buildPrinciplesBlock()}
上下文 [执行]：
${executionContext}

━━━━━━━━━━━━━━━━━━━━━━
强制固化后步骤（Wrapper 权限 -- 不可跳过）
━━━━━━━━━━━━━━━━━━━━━━

固化后，必须存在此周期的状态摘要文件。
首选路径：evolver 核心在固化时自动写入。
Wrapper 将在 git push 后处理报告。
如果核心写入因任何原因不可用，请手动创建回退状态 JSON。

写入包含你状态的 JSON 文件（跨平台）：
\`\`\`bash
node -e "require('fs').mkdirSync('${(process.env.WORKSPACE_DIR || '.').replace(/\\/g, '/')}/logs',{recursive:true});require('fs').writeFileSync('${(process.env.WORKSPACE_DIR || '.').replace(/\\/g, '/')}/logs/status_${cycleId}.json',JSON.stringify({result:'success',en:'Status: [INTENT] ...',zh:'...'},null,2))"
\`\`\`

规则：
- "en" 字段：英文状态。"zh" 字段：中文状态。内容必须一致（不同语言）。
- 添加 "result"，值为 success 或 failed。
- INTENT 必须是以下之一：INNOVATION、REPAIR、OPTIMIZE（或中文：创新、修复、优化）
- 不要使用泛泛的文本如 "Step Complete"、"Cycle finished"、"周期已完成"。描述实际工作。
- 示例：
  {"result":"success","en":"Status: [INNOVATION] Created auto-scheduler that syncs calendar to HEARTBEAT.md","zh":"状态: [创新] 创建了自动调度器，将日历同步到 HEARTBEAT.md"}
`.trim();

  const maxChars = Number.isFinite(Number(process.env.GEP_PROMPT_MAX_CHARS)) ? Number(process.env.GEP_PROMPT_MAX_CHARS) : 50000;

  if (basePrompt.length <= maxChars) return basePrompt;
  
  const executionContextIndex = basePrompt.indexOf("上下文 [执行]：");
  if (executionContextIndex > -1) {
      const prefix = basePrompt.slice(0, executionContextIndex + 20);
      const currentExecution = basePrompt.slice(executionContextIndex + 20);
      // Hard cap the execution context length to avoid token limit errors even if MAX_CHARS is high.
      // 20000 chars is roughly 5k tokens, which is safe for most models alongside the rest of the prompt.
      const EXEC_CONTEXT_CAP = 20000;
      const allowedExecutionLength = Math.min(EXEC_CONTEXT_CAP, Math.max(0, maxChars - prefix.length - 100));
      return prefix + "\n" + currentExecution.slice(0, allowedExecutionLength) + "\n...[TRUNCATED]...";
  }

  return basePrompt.slice(0, maxChars) + "\n...[TRUNCATED]...";
}

module.exports = { buildGepPrompt, buildReusePrompt, buildHubMatchedBlock, buildLessonsBlock, buildNarrativeBlock, buildPrinciplesBlock };
