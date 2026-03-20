---
name: qmd
description: 本地搜索/索引 CLI 工具（支持 BM25 + 向量搜索 + 重排序），附带 MCP 模式。
homepage: https://tobi.lutke.com
metadata: {"clawdbot":{"emoji":"📝","requires":{"bins":["qmd"]},"install":[{"id":"node","kind":"node","package":"https://github.com/tobi/qmd","bins":["qmd"],"label":"安装 qmd (node)"}]}}
---

# qmd

使用 `qmd` 对本地文件建立索引并进行搜索。

## 索引操作
- 添加集合：`qmd collection add /path --name docs --mask "**/*.md"`
- 更新索引：`qmd update`
- 查看状态：`qmd status`

## 搜索方式
- BM25 关键词搜索：`qmd search "查询内容"`
- 向量语义搜索：`qmd vsearch "查询内容"`
- 混合搜索（推荐）：`qmd query "查询内容"`
- 获取文档内容：`qmd get docs/path.md:10 -l 40`

## 注意事项
- 嵌入向量和重排序功能使用 Ollama，地址由 `OLLAMA_URL` 指定（默认 `http://localhost:11434`）。
- 索引数据默认存储在 `~/.cache/qmd` 目录下。
- MCP 模式启动命令：`qmd mcp`
