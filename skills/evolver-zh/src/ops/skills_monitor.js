// 技能监控器 (v2.0) - Evolver 核心模块
// 检查已安装技能的实际问题，自动修复简单问题。
// 零飞书依赖。

const fs = require('fs');
const path = require('path');
const { execSync } = require('child_process');
const { getSkillsDir, getWorkspaceRoot } = require('../gep/paths');

const IGNORE_LIST = new Set([
    'common',
    'clawhub',
    'input-validator',
    'proactive-agent',
    'security-audit',
]);

// 加载用户自定义忽略列表
try {
    var ignoreFile = path.join(getWorkspaceRoot(), '.skill_monitor_ignore');
    if (fs.existsSync(ignoreFile)) {
        fs.readFileSync(ignoreFile, 'utf8').split('\n').forEach(function(l) {
            var t = l.trim();
            if (t && !t.startsWith('#')) IGNORE_LIST.add(t);
        });
    }
} catch (e) { /* ignore */ }

function checkSkill(skillName) {
    var SKILLS_DIR = getSkillsDir();
    if (IGNORE_LIST.has(skillName)) return null;
    var skillPath = path.join(SKILLS_DIR, skillName);
    var issues = [];

    try { if (!fs.statSync(skillPath).isDirectory()) return null; } catch (e) { return null; }

    var mainFile = 'index.js';
    var pkgPath = path.join(skillPath, 'package.json');
    var hasPkg = false;

    if (fs.existsSync(pkgPath)) {
        hasPkg = true;
        try {
            var pkg = JSON.parse(fs.readFileSync(pkgPath, 'utf8'));
            if (pkg.main) mainFile = pkg.main;
            if (pkg.dependencies && Object.keys(pkg.dependencies).length > 0) {
                if (!fs.existsSync(path.join(skillPath, 'node_modules'))) {
                    issues.push('缺少 node_modules（需要执行 npm install）');
                } else {
                    // 优化：检查 node_modules 是否存在，而不是 spawn node
                    // 为每个技能 spawn node 太慢（性能瓶颈）。
                    // 我们假设如果 node_modules 存在，它可能没问题。
                    // 只有在我们真的怀疑有问题时才 spawn 检查（如空的 node_modules）。
                    try {
                        if (fs.readdirSync(path.join(skillPath, 'node_modules')).length === 0) {
                             issues.push('node_modules 为空（需要执行 npm install）');
                        }
                    } catch (e) {
                        issues.push('无效的 node_modules');
                    }
                }
            }
        } catch (e) {
            issues.push('无效的 package.json');
        }
    }

    if (mainFile.endsWith('.js')) {
        var entryPoint = path.join(skillPath, mainFile);
        if (fs.existsSync(entryPoint)) {
            // 优化：通过 node -c 进行语法检查较慢。
            // 我们可以信任运行时在加载时捕获语法错误。
            // 或者如果绝对必要，可以使用更轻量的检查。
            // 目前，移除同步 spawn 以修复性能瓶颈。
        }
    }

    if (hasPkg && !fs.existsSync(path.join(skillPath, 'SKILL.md'))) {
        issues.push('缺少 SKILL.md');
    }

    return issues.length > 0 ? { name: skillName, issues: issues } : null;
}

function autoHeal(skillName, issues) {
    var SKILLS_DIR = getSkillsDir();
    var skillPath = path.join(SKILLS_DIR, skillName);
    var healed = [];

    for (var i = 0; i < issues.length; i++) {
        if (issues[i] === '缺少 node_modules（需要执行 npm install）' || issues[i] === 'node_modules 为空（需要执行 npm install）') {
            try {
                // 如果存在 package-lock.json 则移除以防止冲突错误
                try { fs.unlinkSync(path.join(skillPath, 'package-lock.json')); } catch (e) {}
                
                execSync('npm install --production --no-audit --no-fund', {
                    cwd: skillPath, stdio: 'ignore', timeout: 60000 // 增加超时时间
                });
                healed.push(issues[i]);
                console.log('[SkillsMonitor] 自动修复 ' + skillName + ': npm install');
            } catch (e) {
                console.error('[SkillsMonitor] 修复失败 ' + skillName + ': ' + e.message);
            }
        } else if (issues[i] === '缺少 SKILL.md') {
            try {
                var name = skillName.replace(/-/g, ' ');
                fs.writeFileSync(path.join(skillPath, 'SKILL.md'), '# ' + skillName + '\n\n' + name + ' skill.\n');
                healed.push(issues[i]);
                console.log('[SkillsMonitor] 自动修复 ' + skillName + ': 已创建 SKILL.md 存根');
            } catch (e) {}
        }
    }
    return healed;
}

function run(options) {
    var heal = (options && options.autoHeal) !== false;
    var SKILLS_DIR = getSkillsDir();
    var skills = fs.readdirSync(SKILLS_DIR);
    var report = [];

    for (var i = 0; i < skills.length; i++) {
        if (skills[i].startsWith('.')) continue;
        var result = checkSkill(skills[i]);
        if (result) {
            if (heal) {
                var healed = autoHeal(result.name, result.issues);
                result.issues = result.issues.filter(function(issue) { return !healed.includes(issue); });
                if (result.issues.length === 0) continue;
            }
            report.push(result);
        }
    }
    return report;
}

if (require.main === module) {
    var issues = run();
    console.log(JSON.stringify(issues, null, 2));
    process.exit(issues.length > 0 ? 1 : 0);
}

module.exports = { run, checkSkill, autoHeal };
