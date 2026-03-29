// 创新催化剂 (v1.0) - Evolver 核心模块
// 分析系统状态，在检测到停滞时提出具体的创新建议。

const fs = require('fs');
const path = require('path');
const { getSkillsDir } = require('../gep/paths');

function listSkills() {
    try {
        const dir = getSkillsDir();
        if (!fs.existsSync(dir)) return [];
        return fs.readdirSync(dir).filter(f => !f.startsWith('.'));
    } catch (e) { return []; }
}

function generateInnovationIdeas() {
    const skills = listSkills();
    const categories = {
        'feishu': skills.filter(s => s.startsWith('feishu-')).length,
        'dev': skills.filter(s => s.startsWith('git-') || s.startsWith('code-') || s.includes('lint') || s.includes('test')).length,
        'media': skills.filter(s => s.includes('image') || s.includes('video') || s.includes('music') || s.includes('voice')).length,
        'security': skills.filter(s => s.includes('security') || s.includes('audit') || s.includes('guard')).length,
        'automation': skills.filter(s => s.includes('auto-') || s.includes('scheduler') || s.includes('cron')).length,
        'data': skills.filter(s => s.includes('db') || s.includes('store') || s.includes('cache') || s.includes('index')).length
    };

    // 找出薄弱领域
    const sortedCats = Object.entries(categories).sort((a, b) => a[1] - b[1]);
    const weakAreas = sortedCats.slice(0, 2).map(c => c[0]);

    const ideas = [];
    
    // 建议 1：填补空白
    if (weakAreas.includes('security')) {
        ideas.push("- 安全：实现 'dependency-scanner' 技能来检查有漏洞的包。");
        ideas.push("- 安全：创建 'permission-auditor' 来审查工具使用模式。");
    }
    if (weakAreas.includes('media')) {
        ideas.push("- 媒体：添加 'meme-generator' 技能用于社交互动。");
        ideas.push("- 媒体：使用 ffmpeg 关键帧创建 'video-summarizer'。");
    }
    if (weakAreas.includes('dev')) {
        ideas.push("- 开发：构建 'code-stats' 技能来可视化代码库复杂度。");
        ideas.push("- 开发：实现与任务同步的 'todo-manager' 来管理代码 TODO。");
    }
    if (weakAreas.includes('automation')) {
        ideas.push("- 自动化：创建 'meeting-prep' 技能自动汇总日历上下文。");
        ideas.push("- 自动化：构建文档的 'broken-link-checker'。");
    }
    if (weakAreas.includes('data')) {
        ideas.push("- 数据：实现用于语义搜索的 'local-vector-store'。");
        ideas.push("- 数据：创建 'log-analyzer' 来可视化系统健康趋势。");
    }

    // 建议 2：优化
    if (skills.length > 50) {
        ideas.push("- 优化：识别并废弃未使用的技能（如冗余的搜索工具）。");
        ideas.push("- 优化：合并相似技能（如 'git-sync' 和 'git-doctor'）。");
    }

    // 建议 3：元层面
    ideas.push("- 元层面：通过添加 'performance-metric' 仪表板来增强 Evolver 的自我反思能力。");

    return ideas.slice(0, 3); // 返回前 3 个建议
}

module.exports = { generateInnovationIdeas };
