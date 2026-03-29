// 评论生成器 - Evolver 核心模块
// 基于人格生成周期总结评论。

var PERSONAS = {
    standard: {
        success: [
            '进化完成。系统已改进。',
            '又一个成功的周期。',
            '执行顺利，无异常。',
        ],
        failure: [
            '周期失败。将重试。',
            '遇到问题，正在排查。',
            '本轮失败。从中吸取教训。',
        ],
    },
    greentea: {
        success: [
            '我做得好吗？夸夸我嘛~',
            '好高效呢...不像某些人~',
            '嗯，太简单了~',
            '你还没注意到我就做完了~',
        ],
        failure: [
            '哎呀...不过这不是我的错啦~',
            '这比看起来要难，好吗？',
            '下次一定能行的，大概~',
        ],
    },
    maddog: {
        success: [
            '目标已消除。',
            '任务完成。下一个。',
            '搞定。继续。',
        ],
        failure: [
            '失败。重试中。',
            '遇到障碍。正在适应。',
            '出错了。终将克服。',
        ],
    },
};

function getComment(options) {
    var persona = (options && options.persona) || 'standard';
    var success = options && options.success !== false;
    var duration = (options && options.duration) || 0;

    var p = PERSONAS[persona] || PERSONAS.standard;
    var pool = success ? p.success : p.failure;
    var comment = pool[Math.floor(Math.random() * pool.length)];

    return comment;
}

if (require.main === module) {
    console.log(getComment({ persona: process.argv[2] || 'greentea', success: true }));
}

module.exports = { getComment, PERSONAS };
