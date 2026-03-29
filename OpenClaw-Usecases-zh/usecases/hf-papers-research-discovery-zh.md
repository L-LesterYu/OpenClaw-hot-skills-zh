# Hugging Face 论文研究发现

追踪机器学习研究进展意味着每天刷新 Hugging Face Papers 页面，浏览几十个标题，逐个点进去看摘要，还要手动交叉引用 GitHub 仓库。你想要一种对话式的方式来发现、筛选和深度阅读热门论文，而不必离开你的工作区。

这个工作流结合了两个技能，构建了一条完整的研究发现流水线：

- 浏览 Hugging Face 上当天的热门论文——按点赞数或日期排序
- 按关键词搜索论文，找到与任何主题相关的工作
- 获取完整论文元数据：摘要、作者、GitHub 仓库、社区点赞数、AI 生成的摘要
- 阅读任何论文的社区讨论和评论
- 通过 arXiv ID 深度阅读论文的完整 LaTeX 源码（使用 arxiv-source）

## 所需技能

- [hf-papers](https://github.com/openclaw/skills/tree/main/skills/willamhou/hf-papers) 技能（4 个工具：`hf_daily_papers`、`hf_search_papers`、`hf_paper_detail`、`hf_paper_comments`）
- [arxiv-source](https://github.com/openclaw/skills/tree/main/skills/willamhou/arxiv-source) 技能（3 个工具：`arxiv_fetch`、`arxiv_sections`、`arxiv_abstract`）——用于获取完整论文文本

无需 Docker 或身份验证——两个技能都使用公共 API 并支持本地缓存。

## 设置方法

1. 安装两个技能：
```bash
clawhub install hf-papers
clawhub install arxiv-source
```

2. 向 OpenClaw 发送你的研究工作流提示：

```text
我想持续追踪 ML 研究进展。以下是我的每日工作流：

1. 每天早上，给我展示 Hugging Face 上排名前 10 的热门论文（按点赞数排序）
   - 每篇论文显示：标题、点赞数、GitHub 仓库（如有）、一句话 AI 摘要

2. 当我说"搜索 [主题]"时：
   - 搜索 HF Papers 并展示最相关的结果
   - 高亮显示有关联 GitHub 仓库或高点赞数的论文

3. 当我选择一篇论文（通过 ID）时：
   - 展示完整摘要、作者和关联资源
   - 展示社区评论（如有）
   - 询问我是否需要深度阅读

4. 深度阅读时：
   - 通过 arxiv-source 获取完整论文
   - 总结核心贡献、方法和结果
   - 标注我应该查看的关联代码仓库

维护一个我今天已审阅论文的清单，附上一句话心得。
```

3. 试试看："今天 Hugging Face Papers 上有什么热门论文？"

---

**原文链接**：[English Version](https://github.com/hesamsheikh/awesome-openclaw-usecases/blob/main/usecases/hf-papers-research-discovery.md)
