# LaTeX 论文写作

搭建本地 LaTeX 环境很痛苦 —— 安装 TeX Live 需要数 GB 空间，调试编译错误很繁琐，在编辑器和 PDF 查看器之间切换会打断思路。你希望无需任何本地设置就能通过对话方式编写和编译 LaTeX 论文。

这个工作流程将你的智能体变成一个具有即时编译功能的 LaTeX 写作助手：

- 与智能体协作编写 LaTeX —— 描述你想要的，它会生成源码
- 使用 pdflatex、xelatex 或 lualatex 即时编译为 PDF（无需本地 TeX 安装）
- 无需切换到其他应用即可内联预览 PDF
- 使用起始模板（article、IEEE、beamer、中文文章）跳过模板代码
- 支持 BibTeX/BibLaTeX 文献管理 —— 只需粘贴你的 .bib 内容

## 所需技能

- [latex-compiler](https://github.com/Prismer-AI/Prismer/tree/main/skills/latex-compiler) 技能（4 个工具：`latex_compile`、`latex_preview`、`latex_templates`、`latex_get_template`）
- Prismer 工作区容器（在端口 8080 运行完整的 TeX Live LaTeX 服务器）

## 如何设置

1. 使用 Docker 克隆并部署 [Prismer](https://github.com/Prismer-AI/Prismer)（LaTeX 服务器会自动启动，包含完整的 TeX Live）：
```bash
git clone https://github.com/Prismer-AI/Prismer.git && cd Prismer
docker compose -f docker/docker-compose.dev.yml up
```

2. `latex-compiler` 技能是内置的 —— 无需安装。向 OpenClaw 发送提示词：
```text
帮助我用 LaTeX 写研究论文。这是我的工作流程：

1. 从 IEEE 模板开始（或根据需要选择 article/beamer）
2. 当我描述一个章节时，为它生成 LaTeX 源码
3. 每次主要编辑后，编译并预览 PDF 以便检查格式
4. 如果有编译错误，读取日志并自动修复
5. 当我提供 BibTeX 条目时，将它们添加到参考文献并重新编译

如果需要中文/CJK 支持，使用 xelatex，否则默认使用 pdflatex。
始终运行 2 遍以处理交叉引用。
```

3. 尝试一下："开始一篇题为 'LLM Agents 调查' 的新 IEEE 论文。给我模板并填写摘要和引言部分，然后编译它。"

---

**原文链接**：[English Version](https://github.com/AlexAnys/awesome-openclaw-usecases/blob/main/usecases/latex-paper-writing.md)