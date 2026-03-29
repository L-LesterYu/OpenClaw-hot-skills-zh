const fs = require('fs');
const os = require('os');
const path = require('path');
const { execSync } = require('child_process');

function getDiskUsage(mount) {
    try {
        // 使用 Node 18+ statfs（如果可用）
        if (fs.statfsSync) {
            const stats = fs.statfsSync(mount || '/');
            const total = stats.blocks * stats.bsize;
            const free = stats.bavail * stats.bsize; // 普通用户可用空间
            const used = total - free;
            return {
                pct: Math.round((used / total) * 100),
                freeMb: Math.round(free / 1024 / 1024)
            };
        }
        // 回退方案
        const out = execSync(`df -P "${mount || '/'}" | tail -1 | awk '{print $5, $4}'`).toString().trim().split(' ');
        return {
            pct: parseInt(out[0].replace('%', '')),
            freeMb: Math.round(parseInt(out[1]) / 1024) // df 通常返回 1k 块
        };
    } catch (e) {
        return { pct: 0, freeMb: 999999, error: e.message };
    }
}

function runHealthCheck() {
    const checks = [];
    let criticalErrors = 0;
    let warnings = 0;

    // 1. 密钥检查（对外部服务至关重要，但对 agent 本身运行可能不是必需的）
    const criticalSecrets = ['FEISHU_APP_ID', 'FEISHU_APP_SECRET'];
    criticalSecrets.forEach(key => {
        if (!process.env[key] || process.env[key].trim() === '') {
            checks.push({ name: `env:${key}`, ok: false, status: 'missing', severity: 'warning' }); // 降级为警告以防止重启循环
            warnings++;
        } else {
            checks.push({ name: `env:${key}`, ok: true, status: 'present' });
        }
    });

    const optionalSecrets = ['CLAWHUB_TOKEN', 'OPENAI_API_KEY'];
    optionalSecrets.forEach(key => {
        if (!process.env[key] || process.env[key].trim() === '') {
            checks.push({ name: `env:${key}`, ok: false, status: 'missing', severity: 'info' });
        } else {
            checks.push({ name: `env:${key}`, ok: true, status: 'present' });
        }
    });

    // 2. 磁盘空间检查
    const disk = getDiskUsage('/');
    if (disk.pct > 90) {
        checks.push({ name: 'disk_space', ok: false, status: `${disk.pct}% used`, severity: 'critical' });
        criticalErrors++;
    } else if (disk.pct > 80) {
        checks.push({ name: 'disk_space', ok: false, status: `${disk.pct}% used`, severity: 'warning' });
        warnings++;
    } else {
        checks.push({ name: 'disk_space', ok: true, status: `${disk.pct}% used` });
    }

    // 3. 内存检查
    const memFree = os.freemem();
    const memTotal = os.totalmem();
    const memPct = Math.round(((memTotal - memFree) / memTotal) * 100);
    if (memPct > 95) {
        checks.push({ name: 'memory', ok: false, status: `${memPct}% used`, severity: 'critical' });
        criticalErrors++;
    } else {
        checks.push({ name: 'memory', ok: true, status: `${memPct}% used` });
    }

    // 4. 进程数量（检查 fork 炸弹或泄漏）
    // 仅限 Linux
    if (process.platform === 'linux') {
        try {
            // 优化：readdirSync /proc 较重。使用更轻量的检查或在高频调用时跳过。
            // 但由于这是健康检查，保留它但提高阈值以减少噪音。
            const pids = fs.readdirSync('/proc').filter(f => /^\d+$/.test(f));
            if (pids.length > 2000) { // 阈值提高到 2000
                 checks.push({ name: 'process_count', ok: false, status: `${pids.length} procs`, severity: 'warning' });
                 warnings++;
            } else {
                 checks.push({ name: 'process_count', ok: true, status: `${pids.length} procs` });
            }
        } catch(e) {}
    }

    // 判断总体状态
    let status = 'ok';
    if (criticalErrors > 0) status = 'error';
    else if (warnings > 0) status = 'warning';

    return {
        status,
        timestamp: new Date().toISOString(),
        checks
    };
}

module.exports = { runHealthCheck };
