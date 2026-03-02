---
name: tavily-search-zh
description: 使用 Tavily API 进行网页搜索。当用户需要搜索互联网信息、新闻、研究资料或任何在线内容时使用此技能。适用于 Brave API 不可用时的替代搜索方案。
metadata:
  {
    "openclaw":
      {
        "requires": { "bins": ["python3"] },
        "emoji": "🔍",
      },
  }
---

# Tavily 网页搜索

使用此技能通过 Tavily API 执行网页搜索。

## 使用方法

```bash
python3 ~/.openclaw/skills/tavily-search-zh/scripts/search.py "你的搜索查询"
```

## API 密钥配置

此技能从 OpenClaw 配置文件（`~/.openclaw/openclaw.json`）中读取 Tavily API 密钥，路径为 `skills.entries.web-search.apiKey`。

### 获取 Tavily API 密钥

1. 访问 [Tavily 官网](https://tavily.com) 注册账号
2. 在控制台获取 API 密钥（免费套餐无需信用卡）

### 配置步骤

打开 `~/.openclaw/openclaw.json`，添加或修改：

```json
{
  "skills": {
    "entries": {
      "web-search": {
        "apiKey": "你的_Tavily_API_密钥"
      }
    }
  }
}
```

## 输出格式

返回 JSON 格式的搜索结果，包含：
- **title**: 网页标题
- **url**: 网页链接
- **content**: 内容摘要

## 使用场景

当用户提出以下需求时，使用此技能：

- "帮我搜索一下最新的 AI 新闻"
- "查找关于 Python 异步编程的资料"
- "今天有什么热点新闻？"
- "搜索 OpenClaw 的使用教程"

## 为什么选择 Tavily？

OpenClaw 内置的 `web_search` 工具依赖 Brave Search API，申请需要信用卡。Tavily 提供**免费套餐**，注册即可使用，**无需信用卡**。

## 手动测试

```bash
python3 ~/.openclaw/skills/tavily-search-zh/scripts/search.py "今天的科技新闻"
```

## 技术说明

- 使用 Tavily Search API 进行搜索
- 支持中文搜索查询
- 默认返回 5 条结果
- 搜索深度：基础模式（basic）
