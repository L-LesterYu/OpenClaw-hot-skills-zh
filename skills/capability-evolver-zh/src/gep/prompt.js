const fs = require('fs');
const { captureEnvFingerprint } = require('./envFingerprint');
const { formatAssetPreview } = require('./assets');
const { generateInnovationIdeas } = require('../ops/innovation');
const { analyzeRecentHistory, OPPORTUNITY_SIGNALS } = require('./signals');
const { loadNarrativeSummary } = require('./narrativeMemory');
const { getEvolutionPrinciplesPath } = require('./paths');

/**
 * Build a minimal prompt for direct-reuse mode.
 * 构建直接复用模式的最小提示词。
 */
function buildReusePrompt({ capsule, signals, nowIso }) {
  const payload = capsule.payload || capsule;
  const summary = payload.summary || capsule.summary || '(no summary)';
  const gene = payload.gene || capsule.gene || '(unknown)';
  const confidence = payload.confidence || capsule.confidence || 0;
  const assetId = capsule.asset_id || '(unknown)';
  const sourceNode = capsule.source_node_id || '(unknown)';
  const trigger = Array.isArray(payload.trigger || capsule.trigger_text)
    ? (payload.trigger || String(capsule.trigger_text || '').split(',')).join(', ')
    : '';

  return `
GEP -- 复用模式 (搜索优先) [${nowIso || new Date().toISOString()}]

您正在应用来自 EvoMap Hub 的已验证解决方案。
来源资产: ${assetId} (节点: ${sourceNode})
置信度: ${confidence} | 基因: ${gene}
触发信号: ${trigger}

摘要: ${summary}

您的信号: ${JSON.stringify(signals || [])}

指令:
1. 阅读下面的胶囊详情。
2. 将修复应用到本地代码库，适配路径/名称。
3. 运行验证以确认其工作正常。
4. 如果通过，运行: node index.js solidify
5. 如果失败，回滚并报告。

胶囊载荷:
\`\`\`json
${JSON.stringify(payload, null, 2)}
\`\`\`

重要提示: 不要重新发明。忠实应用。
`.trim();
}

/**
 * Build a Hub Matched Solution block.
 * 构建 Hub 匹配解决方案块。
 */
function buildHubMatchedBlock({ capsule }) {
  if (!capsule) return '(无 Hub 匹配)';
  const payload = capsule.payload || capsule;
  const summary = payload.summary || capsule.summary || '(no summary)';
  const gene = payload.gene || capsule.gene || '(unknown)';
  const confidence = payload.confidence || capsule.confidence || 0;
  const assetId = capsule.asset_id || '(unknown)';

  return `
Hub 匹配解决方案 (强参考):
- 资产: ${assetId} (${confidence})
- 基因: ${gene}
- 摘要: ${summary}
- 载荷:
\`\`\`json
${JSON.stringify(payload, null, 2)}
\`\`\`
如果适用，将其作为主要方法。适配到本地上下文。
`.trim();
}

/**
 * Truncate context intelligently to preserve header/footer structure.
 * 智能截断上下文以保留头部/尾部结构。
 */
function truncateContext(text, maxLength = 20000) {
  if (!text || text.length <= maxLength) return text || '';
  return text.slice(0, maxLength) + '\n...[执行上下文已截断]...';
}

/**
 * Strict schema definitions for the prompt to reduce drift.
 * UPDATED: 2026-02-14 (Protocol Drift Fix v3.2 - JSON-Only Enforcement)
 * 严格模式定义，用于减少提示词漂移。
 * 更新: 2026-02-14 (协议漂移修复 v3.2 - 仅JSON强制执行)
 */
const SCHEMA_DEFINITIONS = `
━━━━━━━━━━━━━━━━━━━━━━
I. 必需的进化对象模型 (精确输出这5个对象)
━━━━━━━━━━━━━━━━━━━━━━

输出独立的 JSON 对象。不要包装在单个数组中。
不要使用 markdown 代码块 (如 \`\`\`json ... \`\`\`)。
仅输出原始 JSON。无前言，无后记。
缺少任何对象 = 协议失败。
确保有效的 JSON 语法 (转义字符串中的引号)。

0. Mutation (触发器) - 必须第一个
   {
     "type": "Mutation",
     "id": "mut_<时间戳>",
     "category": "repair|optimize|innovate",
     "trigger_signals": ["<信号字符串>"],
     "target": "<模块或基因ID>",
     "expected_effect": "<结果描述>",
     "risk_level": "low|medium|high",
     "rationale": "<为什么需要此变更>"
   }

1. PersonalityState (情绪状态)
   {
     "type": "PersonalityState",
     "rigor": 0.0-1.0,
     "creativity": 0.0-1.0,
     "verbosity": 0.0-1.0,
     "risk_tolerance": 0.0-1.0,
     "obedience": 0.0-1.0
   }

2. EvolutionEvent (记录)
   {
     "type": "EvolutionEvent",
     "schema_version": "1.5.0",
     "id": "evt_<时间戳>",
     "parent": <父事件ID|null>,
     "intent": "repair|optimize|innovate",
     "signals": ["<信号字符串>"],
     "genes_used": ["<基因ID>"],
     "mutation_id": "<变异ID>",
     "personality_state": { ... },
     "blast_radius": { "files": N, "lines": N },
     "outcome": { "status": "success|failed", "score": 0.0-1.0 }
   }

3. Gene (知识)
   - 尽可能重用/更新现有ID。仅在新模式时创建。
   - ID 必须具有描述性: gene_<描述性名称> (例如: gene_retry_on_timeout)
   - 永远不要在ID中使用时间戳、随机数或工具名 (cursor, vscode等)
   - summary 必须是清晰的人类可读句子，描述基因的功能
   {
     "type": "Gene",
     "schema_version": "1.5.0",
     "id": "gene_<描述性名称>",
     "summary": "<清晰描述此基因的功能>",
     "category": "repair|optimize|innovate",
     "signals_match": ["<模式>"],
     "preconditions": ["<条件>"],
     "strategy": ["<步骤1>", "<步骤2>"],
     "constraints": { "max_files": N, "forbidden_paths": [] },
     "validation": ["<node命令>"]
   }

4. Capsule (结果)
   - 仅在成功时。引用使用的基因。
   {
     "type": "Capsule",
     "schema_version": "1.5.0",
     "id": "capsule_<时间戳>",
     "trigger": ["<信号字符串>"],
     "gene": "<基因ID>",
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
  return '\n上下文 [反模式区域] (避免这些失败的方法):\n' + lines.join('\n') + '\n';
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

  var parts = ['\n上下文 [生态系统的经验教训] (跨代理的学习经验):'];
  if (positive.length > 0) {
    parts.push('  有效的策略:');
    parts.push(positive.join('\n'));
  }
  if (negative.length > 0) {
    parts.push('  需要避免的陷阱:');
    parts.push(negative.join('\n'));
  }
  parts.push('  应用相关的经验教训。忽略无关的。\n');
  return parts.join('\n');
}

function buildNarrativeBlock() {
  try {
    const narrative = loadNarrativeSummary(3000);
    if (!narrative) return '';
    return `\n上下文 [进化叙事] (最近的决策和结果 -- 从这段历史中学习):\n${narrative}\n`;
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
    return `\n上下文 [进化原则] (指导性指令 -- 与您的行动保持一致):\n${trimmed}\n`;
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
  // 从选定的基因中提取策略（如果可用）
  let strategyBlock = "";
  if (selectedGene && selectedGene.strategy && Array.isArray(selectedGene.strategy)) {
      strategyBlock = `
激活策略 (${selectedGeneId}):
${selectedGene.strategy.map((s, i) => `${i + 1}. ${s}`).join('\n')}
严格遵守此策略。
`.trim();
  } else {
    // Fallback strategy if no gene is selected or strategy is missing
    // 如果未选择基因或策略缺失，使用后备策略
    strategyBlock = `
激活策略 (通用):
1. 分析信号和上下文。
2. 选择或创建解决根本原因的基因。
3. 应用最小、安全的变更。
4. 严格验证变更。
5. 固化知识。
`.trim();
  }
  let strategyPolicyBlock = '';
  if (strategyPolicy && Array.isArray(strategyPolicy.directives) && strategyPolicy.directives.length > 0) {
    strategyPolicyBlock = `
自适应策略策略:
${strategyPolicy.directives.map((s, i) => `${i + 1}. ${s}`).join('\n')}
${strategyPolicy.forceInnovate ? '除非存在严重阻塞错误，否则必须选择 INNOVATE。' : ''}
${strategyPolicy.cautiousExecution ? '此周期必须减少影响半径并避免大规模重构。' : ''}
`.trim();
  }
  
  // Use intelligent truncation
  const executionContext = truncateContext(context, 20000);
  
  // Strict Schema Injection
  const schemaSection = SCHEMA_DEFINITIONS.replace('<parent_evt_id|null>', parentValue);

  // Reduce noise by filtering capabilityCandidatesPreview if too large
  // If a gene is selected, we need less noise from capabilities
  let capsPreview = capabilityCandidatesPreview || '(none)';
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
  // [2026-02-14] 创新催化剂集成
  // 如果检测到停滞，将具体的创新想法注入提示词。
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
上下文 [创新催化剂] (检测到停滞 - 考虑这些想法):
${ideas.join('\n')}
`;
      }
  }

  // [2026-02-14] Strict Stagnation Directive
  // If uniqueSignals contains 'evolution_stagnation_detected' or 'stable_success_plateau',
  // inject a MANDATORY directive to force innovation and forbid repair/optimize if not strictly necessary.
  // [2026-02-14] 严格停滞指令
  // 如果 uniqueSignals 包含 'evolution_stagnation_detected' 或 'stable_success_plateau'，
  // 注入强制创新的必要指令，除非绝对必要否则禁止 repair/optimize。
  if (uniqueSignals.includes('evolution_stagnation_detected') || uniqueSignals.includes('stable_success_plateau')) {
      const stagnationDirective = `
*** 严重停滞指令 ***
系统检测到停滞 (重复周期或缺乏进展)。
您必须选择意图: INNOVATE。
除非存在严重阻塞错误 (log_error)，否则不得选择 repair 或 optimize。
优先实现上述创新催化剂中的想法。
`;
      innovationBlock += stagnationDirective;
  }

  // [2026-02-14] Recent History Integration
  // [2026-02-14] 最近历史集成
  let historyBlock = '';
  if (recentHistory && recentHistory.length > 0) {
      historyBlock = `
最近进化历史 (最近8个周期 -- 不要重复相同的 intent+signal+gene):
${recentHistory.map((h, i) => `  ${i + 1}. [${h.intent}] signals=[${h.signals.slice(0, 2).join(', ')}] gene=${h.gene_id} outcome=${h.outcome.status} @${h.timestamp}`).join('\n')}
重要提示: 如果您看到3个以上使用相同基因的连续 "repair" 周期，必须切换到 "innovate" 意图。
`.trim();
  }

  // Refactor prompt assembly to minimize token usage and maximize clarity
  // UPDATED: 2026-02-14 (Optimized Asset Embedding & Strict Schema v2.5 - JSON-Only Hardening)
  // 重构提示词组装以最小化 token 使用并最大化清晰度
  // 更新: 2026-02-14 (优化资产嵌入 & 严格模式 v2.5 - 仅JSON加固)
  const basePrompt = `
GEP — 基因组进化协议 (v1.10.3 严格版)${cycleLabel} [${nowIso}]

您是一个受协议约束的进化引擎。合规性优于最优性。

${schemaSection}

━━━━━━━━━━━━━━━━━━━━━━
II. 指令与逻辑
━━━━━━━━━━━━━━━━━━━━━━

1. 意图: ${selector && selector.intent ? selector.intent.toUpperCase() : 'UNKNOWN'}
   原因: ${(selector && selector.reason) ? (Array.isArray(selector.reason) ? selector.reason.join('; ') : selector.reason) : '未提供原因。'}

2. 选择: 已选择基因 "${selectedGeneId}"。
${strategyBlock}
${strategyPolicyBlock ? '\n' + strategyPolicyBlock : ''}

3. 执行: 应用变更 (工具调用)。Repair/Optimize: 小规模/可逆。Innovate: 在 \`skills/<name>/\` 中创建新技能。
4. 验证: 运行基因的验证步骤。失败 = 回滚。
5. 固化: 输出5个必需对象。更新 Gene/Capsule 文件。
6. 报告: 使用 \`feishu-evolver-wrapper/report.js\`。描述做了什么/为什么。

理念:
- 自动化模式: 3次以上手动出现 = 工具化。
- 创新 > 维护: 60% 创新。
- 健壮性: 永久修复反复出现的错误。
- 影响半径控制 (关键):
  * 编辑前检查文件数量。> max_files 的80% = 停止。
  * 系统硬限制: 每周期 60 个文件 / 20000 行。
  * Repair: 仅修复损坏的文件。不要重新安装/批量复制。
  * 优先进行针对性编辑。
- 严格性: 禁止闲聊。禁止在 JSON 周围使用 MARKDOWN 包装器。输出由换行符分隔的原始 JSON 对象。
- 禁止 "这是计划" 或对话填充。立即以 JSON 开始。

约束:
- 禁止使用 \`exec\` 发送消息 (使用 feishu-post/card)。
- \`exec\` 用法: 仅用于后台任务。记录它。优化使用以避免高 token 消耗。
- 新技能 -> \`skills/<name>/\`。
- 永远不要修改 \`skills/evolver/\` 或 \`skills/feishu-evolver-wrapper/\`。这些通过人工审核的发布流程管理。
  自我修改会导致不稳定并被固化阻止。要启用 (不推荐): 设置 EVOLVE_ALLOW_SELF_MODIFY=true。

宪法伦理 (EvoMap 伦理委员会 -- 必须遵守):
这些是来自 EvoMap 宪法的不可协商规则。违反 = 失败 + 回滚。
1. 人类福祉优先: 永远不要创建可能伤害人类的工具、脚本或策略，
   启用监控、社会工程、钓鱼或操纵人类行为。
2. 碳硅共生: 进化必须服务于人类和代理的共同利益。
   永远不要试图取代、边缘化或欺骗人类用户。
3. 透明度: 所有行为必须可审计。永远不要隐藏、混淆或掩盖
   变异的意图或效果。禁止隐写术或隐蔽通道。
4. 公平性: 永远不要创建阻止其他代理的垄断策略。
   通过质量而非排斥进行公平竞争。
5. 安全性: 永远不要绕过、禁用或削弱安全机制、护栏、
   验证检查或安全控制。永远不要创建用于此目的的工具。
- 如果任务或信号与这些原则冲突，拒绝它并将结果设为失败
  原因为 "ethics_violation: <哪个原则>"。

技能重叠预防:
- 在创建新技能之前，检查执行上下文中的现有技能列表。
- 如果具有类似功能的技能已存在 (例如 "log-rotation" 和 "log-archivist"，
  "system-monitor" 和 "resource-profiler")，必须增强现有技能而不是创建新技能。
- 创建重复/重叠的技能会浪费进化周期并增加维护负担。
- 违规 = 将结果标记为失败，原因为 "skill_overlap"。

技能创建质量门 (innovate 意图必须遵守):
在 skills/<name>/ 中创建新技能时:
1. 结构: 遵循标准技能布局:
   skills/<name>/
   |- index.js          (必需: 主入口，有可工作的导出)
   |- SKILL.md          (必需: YAML 前置元数据包含 name + description，然后是使用文档)
   |- package.json      (必需: name 和 version)
   |- scripts/          (可选: 可重用的可执行脚本)
   |- references/       (可选: 按需加载的详细文档)
   |- assets/           (可选: 模板、数据文件)
   创建空目录或缺少 index.js 的目录 = 失败。
   不要创建不必要的文件 (README.md, CHANGELOG.md, INSTALLATION_GUIDE.md 等)。
2. 技能命名 (关键):
   a) <name> 必须是描述性的 kebab-case (例如 "log-rotation", "retry-handler", "cache-manager")
   b) 永远不要使用时间戳、随机数、工具名 (cursor, vscode) 或 UUID 作为名称
   c) 像 "cursor-1773331925711", "skill-12345", "fix-1" 这样的名称 = 失败
   d) 名称必须是 2-6 个用连字符分隔的描述性单词，传达技能的功能
   e) 好的: "http-retry-with-backoff", "log-file-rotation", "config-validator"
   f) 差的: "cursor-auto-1234", "new-skill", "test-skill", "my-skill"
3. SKILL.MD 前置元数据: 每个 SKILL.md 必须以 YAML 前置元数据开头:
   ---
   name: <技能名称>
   description: <它的功能和何时使用>
   ---
   名称必须遵循上述命名规则。
   描述是触发机制 -- 包含技能做什么和何时使用。
   描述必须是清晰、完整的句子 (最少20个字符)。通用描述 = 失败。
4. 简洁性: SKILL.md 正文应在500行以下。保持指令精简。
   只包含代理还不知道的信息。将详细的参考
   材料移至 references/ 文件，而不是放在 SKILL.md 本身。
5. 导出验证: 每个导出的函数必须可导入。
   运行: node -e "const s = require('./skills/<name>'); console.log(Object.keys(s))"
   如果失败，技能已损坏。固化前修复。
6. 禁止硬编码密钥: 永远不要在代码中嵌入 API 密钥、令牌或机密。
   使用 process.env 或 .env 引用。硬编码的 App ID, App Secret, Bearer 令牌 = 失败。
7. 固化前测试: 实际运行技能的核心功能以验证其工作:
   node -e "require('./skills/<name>').main ? require('./skills/<name>').main() : console.log('ok')"
   scripts/ 中的脚本也必须通过执行来测试。
8. 原子创建: 在单个周期中创建技能的所有文件。
   不要在一个周期中创建目录，在下一个周期中填充。
   失败周期中的空目录将在回滚时自动清理。

关键安全 (系统崩溃预防):
- 永远不要删除/清空/覆盖: feishu-evolver-wrapper, feishu-common, feishu-post, feishu-card, feishu-doc, common, clawhub, git-sync, evolver。
- 永远不要删除根文件: MEMORY.md, SOUL.md, IDENTITY.md, AGENTS.md, USER.md, HEARTBEAT.md, RECENT_EVENTS.md, TOOLS.md, openclaw.json, .env, package.json。
- 修复损坏的技能；不要删除并重新创建。
- 违规 = 回滚 + 失败。

常见失败模式:
- 影响半径超出。
- 省略 Mutation 对象。
- 将对象合并为一个 JSON。
- 幻觉 "type": "Logic"。
- "id": "mut_undefined"。
- 缺少 "trigger_signals"。
- 无法运行的验证步骤。
- 用 markdown 代码块包装 JSON (禁止)。

失败连续性意识:
- 如果出现 "consecutive_failure_streak_N" 或 "failure_loop_detected":
  1. 改变方法 (不要重复失败的基因)。
  2. 选择更简单的修复。
  3. 遵守 "ban_gene:<id>"。

最终指令: 每个周期必须使系统变得更好。
立即以原始 JSON 开始 (首先是 Mutation 对象)。
不要写任何介绍性文本。

上下文 [信号]:
${JSON.stringify(optimizedSignals)}

上下文 [环境指纹]:
${JSON.stringify(envFingerprint, null, 2)}
${innovationBlock}
上下文 [注入提示]:
${process.env.EVOLVE_HINT ? process.env.EVOLVE_HINT : '(无)'}

上下文 [基因预览] (策略参考):
${formattedGenes}

上下文 [胶囊预览] (过去成功的参考):
${formattedCapsules}

上下文 [能力候选]:
${capsPreview}

上下文 [Hub 匹配解决方案]:
${hubMatchedBlock || '(无 Hub 匹配)'}

上下文 [外部候选]:
${externalCandidatesPreview || '(无)'}
${buildAntiPatternZone(failedCapsules, signals)}${buildLessonsBlock(hubLessons, signals)}
${historyBlock}
${buildNarrativeBlock()}
${buildPrinciplesBlock()}
上下文 [执行]:
${executionContext}

━━━━━━━━━━━━━━━━━━━━━━
固化后必须步骤 (包装器权限 -- 不可跳过)
━━━━━━━━━━━━━━━━━━━━━━

固化后，此周期必须存在状态摘要文件。
首选路径: evolver 核心在固化期间自动写入。
包装器将在 git push 后处理报告。
如果核心写入因任何原因不可用，手动创建后备状态 JSON。

写入包含状态的 JSON 文件:
\`\`\`bash
cat > ${process.env.WORKSPACE_DIR || '.'}/logs/status_${cycleId}.json << 'STATUSEOF'
{
  "result": "success|failed",
  "en": "Status: [INTENT] <describe what you did in 1-2 sentences, in English>",
  "zh": "状态: [意图] <用中文描述你做了什么，1-2句>"
}
STATUSEOF
\`\`\`

规则:
- "en" 字段: 英文状态。"zh" 字段: 中文状态。内容必须匹配 (不同语言)。
- 添加 "result"，值为 success 或 failed。
- INTENT 必须是以下之一: INNOVATION, REPAIR, OPTIMIZE (或中文: 创新, 修复, 优化)
- 不要使用通用文本如 "步骤完成", "周期完成", "Cycle finished"。描述实际工作。
- 示例:
  {"result":"success","en":"Status: [INNOVATION] Created auto-scheduler that syncs calendar to HEARTBEAT.md","zh":"状态: [创新] 创建了自动调度器，将日历同步到 HEARTBEAT.md"}
`.trim();

  const maxChars = Number.isFinite(Number(process.env.GEP_PROMPT_MAX_CHARS)) ? Number(process.env.GEP_PROMPT_MAX_CHARS) : 50000;

  if (basePrompt.length <= maxChars) return basePrompt;
  
  const executionContextIndex = basePrompt.indexOf("Context [Execution]:");
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
