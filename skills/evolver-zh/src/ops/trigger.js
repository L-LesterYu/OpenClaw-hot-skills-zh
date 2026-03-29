// Evolver 唤醒触发器 - Evolver 核心模块
// 写入信号文件，供 wrapper 轮询以立即唤醒。

const fs = require('fs');
const path = require('path');
const { getWorkspaceRoot } = require('../gep/paths');

var WAKE_FILE = path.join(getWorkspaceRoot(), 'memory', 'evolver_wake.signal');

function send() {
    try {
        fs.writeFileSync(WAKE_FILE, 'WAKE');
        console.log('[Trigger] 唤醒信号已发送至 ' + WAKE_FILE);
        return true;
    } catch (e) {
        console.error('[Trigger] 发送失败: ' + e.message);
        return false;
    }
}

function clear() {
    try { if (fs.existsSync(WAKE_FILE)) fs.unlinkSync(WAKE_FILE); } catch (e) {}
}

function isPending() {
    return fs.existsSync(WAKE_FILE);
}

if (require.main === module) {
    send();
}

module.exports = { send, clear, isPending };
