const path = require('path');
const fs = require('fs');

function getRepoRoot() {
  if (process.env.EVOLVER_REPO_ROOT) {
    return process.env.EVOLVER_REPO_ROOT;
  }

  const ownDir = path.resolve(__dirname, '..', '..');

  // 安全检查：先检查 evolver 自身目录，防止在包含 .git 的
  // 父级仓库中操作（git reset --hard 在错误范围内运行可能导致数据丢失）。
  if (fs.existsSync(path.join(ownDir, '.git'))) {
    return ownDir;
  }

  let dir = path.dirname(ownDir);
  while (dir !== '/' && dir !== '.') {
    if (fs.existsSync(path.join(dir, '.git'))) {
      if (process.env.EVOLVER_USE_PARENT_GIT === 'true') {
        console.warn('[evolver] 使用父级 git 仓库：', dir);
        return dir;
      }
      console.warn(
        '[evolver] 检测到父级目录中存在 .git', dir,
        '-- 已忽略。设置 EVOLVER_USE_PARENT_GIT=true 覆盖，',
        '或设置 EVOLVER_REPO_ROOT 显式指定目标目录。'
      );
      return ownDir;
    }
    dir = path.dirname(dir);
  }

  return ownDir;
}

function getWorkspaceRoot() {
  if (process.env.OPENCLAW_WORKSPACE) {
    return process.env.OPENCLAW_WORKSPACE;
  }

  const repoRoot = getRepoRoot();
  const workspaceDir = path.join(repoRoot, 'workspace');
  if (fs.existsSync(workspaceDir)) {
    return workspaceDir;
  }

  // 独立/Cursor/非 OpenClaw：使用仓库根目录作为工作区。
  // 旧的 4 级向上回退假设了 OpenClaw 的技能目录布局
  // (/workspace/skills/evolver/)，在其他环境中解析不正确。
  return repoRoot;
}

function getLogsDir() {
  return process.env.EVOLVER_LOGS_DIR || path.join(getWorkspaceRoot(), 'logs');
}

function getEvolverLogPath() {
  return path.join(getLogsDir(), 'evolver_loop.log');
}

function getMemoryDir() {
  return process.env.MEMORY_DIR || path.join(getWorkspaceRoot(), 'memory');
}

// --- 会话作用域隔离 ---
// 当设置了 EVOLVER_SESSION_SCOPE 时（例如 Discord 频道 ID 或项目名称），
// 进化状态、记忆图谱和资产会被隔离到每个作用域的子目录中。
// 这防止了跨频道/跨项目的记忆污染。
// 未设置时，一切照旧工作（全局作用域，向后兼容）。
function getSessionScope() {
  const raw = String(process.env.EVOLVER_SESSION_SCOPE || '').trim();
  if (!raw) return null;
  // 净化：仅允许字母数字、连字符、下划线、点（防止路径穿越）。
  const safe = raw.replace(/[^a-zA-Z0-9_\-\.]/g, '_').slice(0, 128);
  if (!safe || /^\.{1,2}$/.test(safe) || /\.\./.test(safe)) return null;
  return safe;
}

function getEvolutionDir() {
  const baseDir = process.env.EVOLUTION_DIR || path.join(getMemoryDir(), 'evolution');
  const scope = getSessionScope();
  if (scope) {
    return path.join(baseDir, 'scopes', scope);
  }
  return baseDir;
}

function getGepAssetsDir() {
  const repoRoot = getRepoRoot();
  const baseDir = process.env.GEP_ASSETS_DIR || path.join(repoRoot, 'assets', 'gep');
  const scope = getSessionScope();
  if (scope) {
    return path.join(baseDir, 'scopes', scope);
  }
  return baseDir;
}

function getSkillsDir() {
  return process.env.SKILLS_DIR || path.join(getWorkspaceRoot(), 'skills');
}

function getNarrativePath() {
  return path.join(getEvolutionDir(), 'evolution_narrative.md');
}

function getEvolutionPrinciplesPath() {
  const repoRoot = getRepoRoot();
  const custom = path.join(repoRoot, 'EVOLUTION_PRINCIPLES.md');
  if (fs.existsSync(custom)) return custom;
  return path.join(repoRoot, 'assets', 'gep', 'EVOLUTION_PRINCIPLES.md');
}

function getReflectionLogPath() {
  return path.join(getEvolutionDir(), 'reflection_log.jsonl');
}

module.exports = {
  getRepoRoot,
  getWorkspaceRoot,
  getLogsDir,
  getEvolverLogPath,
  getMemoryDir,
  getEvolutionDir,
  getGepAssetsDir,
  getSkillsDir,
  getSessionScope,
  getNarrativePath,
  getEvolutionPrinciplesPath,
  getReflectionLogPath,
};

