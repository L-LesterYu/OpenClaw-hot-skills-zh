// Git 自修复 - Evolver 核心模块
// Git 同步失败的紧急修复：中止 rebase/merge，移除过期锁文件。

const fs = require('fs');
const path = require('path');
const { execSync } = require('child_process');
const { getWorkspaceRoot } = require('../gep/paths');

var LOCK_MAX_AGE_MS = 10 * 60 * 1000; // 10 分钟

function repair(gitRoot) {
    var root = gitRoot || getWorkspaceRoot();
    var repaired = [];

    // 1. 中止待处理的 rebase
    try {
        execSync('git rebase --abort', { cwd: root, stdio: 'ignore' });
        repaired.push('rebase_aborted');
        console.log('[SelfRepair] 已中止待处理的 rebase。');
    } catch (e) {}

    // 2. 中止待处理的 merge
    try {
        execSync('git merge --abort', { cwd: root, stdio: 'ignore' });
        repaired.push('merge_aborted');
        console.log('[SelfRepair] 已中止待处理的 merge。');
    } catch (e) {}

    // 3. 移除过期的 index.lock
    var lockFile = path.join(root, '.git', 'index.lock');
    if (fs.existsSync(lockFile)) {
        try {
            var stat = fs.statSync(lockFile);
            var age = Date.now() - stat.mtimeMs;
            if (age > LOCK_MAX_AGE_MS) {
                fs.unlinkSync(lockFile);
                repaired.push('stale_lock_removed');
                console.log('[SelfRepair] 已移除过期的 index.lock（已存在 ' + Math.round(age / 60000) + ' 分钟）。');
            }
        } catch (e) {}
    }

    // 4. 如果本地损坏则重置到远程 main（最后手段 - 受标志保护）
    // 仅在显式调用 --force-reset 或 EVOLVE_GIT_RESET=true 时启用
    if (process.env.EVOLVE_GIT_RESET === 'true') {
        try {
            console.log('[SelfRepair] 正在将本地分支重置为 origin/main（硬重置）...');
            execSync('git fetch origin main', { cwd: root, stdio: 'ignore' });
            execSync('git reset --hard origin/main', { cwd: root, stdio: 'ignore' });
            repaired.push('hard_reset_to_origin');
        } catch (e) {
            console.warn('[SelfRepair] 硬重置失败: ' + e.message);
        }
    } else {
        // 安全 fetch
        try {
            execSync('git fetch origin', { cwd: root, stdio: 'ignore', timeout: 30000 });
            repaired.push('fetch_ok');
        } catch (e) {
            console.warn('[SelfRepair] git fetch 失败: ' + e.message);
        }
    }

    return repaired;
}

if (require.main === module) {
    var result = repair();
    console.log('[SelfRepair] 结果:', result.length > 0 ? result.join(', ') : '无需修复');
}

module.exports = { repair };
