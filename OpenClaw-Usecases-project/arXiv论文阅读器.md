# arXiv论文阅读器

阅读arXiv论文意味着下载PDF文件，在不同论文间切换时失去上下文，还要费力解析密集的LaTeX符号。你希望能在工作空间中通过对话方式来阅读、分析和比较论文，而无需离开当前环境。

这个工作流程将你的智能体变成一个研究阅读助手：

- 通过ID获取任何arXiv论文并获得清晰易读的文本（LaTeX自动扁平化处理）
- 先浏览论文结构——列出章节，在决定阅读全文前有选择性地阅读
- 快速扫描多篇论文的摘要，为阅读清单排优先级
- 要求智能体总结、比较或评论特定章节
- 结果会被本地缓存——重新访问论文时即时加载

## 你需要的技能

- [arxiv-reader](https://github.com/Prismer-AI/Prismer/tree/main/skills/arxiv-reader) 技能（3个工具：`arxiv_fetch`、`arxiv_sections`、`arxiv_abstract`）

不需要Docker或Python——该技能使用Node.js内置功能独立运行。它直接从arXiv下载，解压LaTeX源代码，并自动处理包含文件。

## 如何设置

1. 从[Prismer仓库](https://github.com/Prismer-AI/Prismer/tree/main/skills/arxiv-reader)安装`arxiv-reader`技能——将`skills/arxiv-reader/`目录复制到你的OpenClaw技能文件夹中。

2. 技能已准备就绪。提示OpenClaw：
```text
我正在研究[主题]。这是我的工作流程：

1. 当我给你一个arXiv ID（如2301.00001）时：
   - 首先获取摘要，让我判断是否相关
   - 如果我说"阅读全文"，则获取完整论文（默认移除附录）
   - 总结主要贡献、方法和结果

2. 当我提供多个ID时：
   - 获取所有摘要并给我一个比较表格
   - 根据我的研究主题按相关性排序

3. 当我询问特定章节时：
   - 先列出论文的章节
   - 然后获取并详细解释相关章节

保持我已读论文和主要收获的运行列表。
```

3. 尝试使用："阅读2401.04088——主要贡献是什么？"