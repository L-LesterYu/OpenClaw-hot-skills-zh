// 进化策略预设 (v1.1)
// 控制修复、优化和创新意图之间的平衡。
//
// 用法：设置 EVOLVE_STRATEGY 环境变量为以下之一：balanced、innovate、harden、repair-only、
//        early-stabilize、steady-state，或 "auto" 进行自适应选择。
// 默认：balanced（或根据周期数/饱和信号自动检测）
//
// 每个策略定义：
//   repair/optimize/innovate  - 目标分配比例（影响 LLM prompt）
//   repairLoopThreshold       - 最近 8 个周期中触发强制创新的修复比例
//   label                     - 注入到 GEP prompt 中的人类可读名称

var fs = require('fs');
var path = require('path');

var STRATEGIES = {
  'balanced': {
    repair: 0.20,
    optimize: 0.30,
    innovate: 0.50,
    repairLoopThreshold: 0.50,
    label: '均衡模式',
    description: '正常运行。稳定增长，兼顾稳定性。',
  },
  'innovate': {
    repair: 0.05,
    optimize: 0.15,
    innovate: 0.80,
    repairLoopThreshold: 0.30,
    label: '创新优先',
    description: '系统稳定。最大化新功能和能力。',
  },
  'harden': {
    repair: 0.40,
    optimize: 0.40,
    innovate: 0.20,
    repairLoopThreshold: 0.70,
    label: '加固模式',
    description: '重大变更后。专注于稳定性和健壮性。',
  },
  'repair-only': {
    repair: 0.80,
    optimize: 0.20,
    innovate: 0.00,
    repairLoopThreshold: 1.00,
    label: '仅修复',
    description: '紧急模式。先修复所有问题再做其他事情。',
  },
  'early-stabilize': {
    repair: 0.60,
    optimize: 0.25,
    innovate: 0.15,
    repairLoopThreshold: 0.80,
    label: '早期稳定',
    description: '初始周期。优先修复现有问题再进行创新。',
  },
  'steady-state': {
    repair: 0.60,
    optimize: 0.30,
    innovate: 0.10,
    repairLoopThreshold: 0.90,
    label: '稳态维持',
    description: '进化饱和。维持现有能力，最小化创新。',
  },
};

// Read evolution_state.json to get the current cycle count for auto-detection.
function _readCycleCount() {
  try {
    // evolver/memory/evolution_state.json (local to the skill)
    var localPath = path.resolve(__dirname, '..', '..', 'memory', 'evolution_state.json');
    // workspace/memory/evolution/evolution_state.json (canonical path used by evolve.js)
    var workspacePath = path.resolve(__dirname, '..', '..', '..', '..', 'memory', 'evolution', 'evolution_state.json');
    var candidates = [localPath, workspacePath];
    for (var i = 0; i < candidates.length; i++) {
      if (fs.existsSync(candidates[i])) {
        var data = JSON.parse(fs.readFileSync(candidates[i], 'utf8'));
        return data && Number.isFinite(data.cycleCount) ? data.cycleCount : 0;
      }
    }
  } catch (e) {}
  return 0;
}

function resolveStrategy(opts) {
  var signals = (opts && Array.isArray(opts.signals)) ? opts.signals : [];
  var name = String(process.env.EVOLVE_STRATEGY || 'balanced').toLowerCase().trim();

  // Backward compatibility: FORCE_INNOVATION=true maps to 'innovate'
  if (!process.env.EVOLVE_STRATEGY) {
    var fi = String(process.env.FORCE_INNOVATION || process.env.EVOLVE_FORCE_INNOVATION || '').toLowerCase();
    if (fi === 'true') name = 'innovate';
  }

  // Auto-detection: when no explicit strategy is set (defaults to 'balanced'),
  // apply heuristics inspired by Echo-MingXuan's "fix first, innovate later" pattern.
  var isDefault = !process.env.EVOLVE_STRATEGY || name === 'balanced' || name === 'auto';

  if (isDefault) {
    // Early-stabilize: first 5 cycles should focus on fixing existing issues.
    var cycleCount = _readCycleCount();
    if (cycleCount > 0 && cycleCount <= 5) {
      name = 'early-stabilize';
    }

    // Saturation detection: if saturation signals are present, switch to steady-state.
    if (signals.indexOf('force_steady_state') !== -1) {
      name = 'steady-state';
    } else if (signals.indexOf('evolution_saturation') !== -1) {
      name = 'steady-state';
    }
  }

  // Explicit "auto" maps to whatever was auto-detected above (or balanced if no heuristic fired).
  if (name === 'auto') name = 'balanced';

  var strategy = STRATEGIES[name] || STRATEGIES['balanced'];
  strategy.name = name;
  return strategy;
}

function getStrategyNames() {
  return Object.keys(STRATEGIES);
}

module.exports = { resolveStrategy, getStrategyNames, STRATEGIES };
