// GEP 制品清理 - Evolver 核心模块
// 从 evolution 目录中删除旧的 gep_prompt_*.json/txt 文件。
// 无论文件年龄，至少保留最近的 10 个文件。

const fs = require('fs');
const path = require('path');
const { getEvolutionDir } = require('../gep/paths');

var MAX_AGE_MS = 24 * 60 * 60 * 1000; // 24 小时
var MIN_KEEP = 10;

function safeBatchDelete(batch) {
    var deleted = 0;
    for (var i = 0; i < batch.length; i++) {
        try { fs.unlinkSync(batch[i]); deleted++; } catch (_) {}
    }
    return deleted;
}

function run() {
    var evoDir = getEvolutionDir();
    if (!fs.existsSync(evoDir)) return;

    var files = fs.readdirSync(evoDir)
        .filter(function(f) { return /^gep_prompt_.*\.(json|txt)$/.test(f); })
        .map(function(f) {
            var full = path.join(evoDir, f);
            var stat = fs.statSync(full);
            return { name: f, path: full, mtime: stat.mtimeMs };
        })
        .sort(function(a, b) { return b.mtime - a.mtime; }); // 最新的在前

    var now = Date.now();
    var deleted = 0;

    // 阶段 1：基于年龄的清理（至少保留 MIN_KEEP 个）
    var filesToDelete = [];
    for (var i = MIN_KEEP; i < files.length; i++) {
        if (now - files[i].mtime > MAX_AGE_MS) {
            filesToDelete.push(files[i].path);
        }
    }

    if (filesToDelete.length > 0) {
        deleted += safeBatchDelete(filesToDelete);
    }

    // 阶段 2：基于大小的安全上限（最多保留 10 个文件）
    try {
        var remainingFiles = fs.readdirSync(evoDir)
            .filter(function(f) { return /^gep_prompt_.*\.(json|txt)$/.test(f); })
            .map(function(f) {
                var full = path.join(evoDir, f);
                var stat = fs.statSync(full);
                return { name: f, path: full, mtime: stat.mtimeMs };
            })
            .sort(function(a, b) { return b.mtime - a.mtime; }); // 最新的在前

        var MAX_FILES = 10;
        if (remainingFiles.length > MAX_FILES) {
            var toDelete = remainingFiles.slice(MAX_FILES).map(function(f) { return f.path; });
            deleted += safeBatchDelete(toDelete);
        }
    } catch (e) {
        console.warn('[Cleanup] 第二阶段失败:', e.message);
    }

    if (deleted > 0) {
        console.log('[Cleanup] 已删除 ' + deleted + ' 个旧 GEP 制品。');
    }
    return deleted;
}

if (require.main === module) {
    console.log('[Cleanup] 正在扫描旧制品...');
    var count = run();
    console.log('[Cleanup] ' + (count > 0 ? '已删除 ' + count + ' 个文件。' : '没有需要删除的文件。'));
}

module.exports = { run };
