# 新闻摘要

## 概述

从可信的国际新闻源通过 RSS 订阅获取并摘要新闻。

## 使用时机

当用户询问新闻动态、每日简报或世界正在发生什么时，使用此技能。

## RSS 新闻源

### BBC（主要来源）

# 世界新闻
curl -s "https://feeds.bbci.co.uk/news/world/rss.xml"

# 头条新闻
curl -s "https://feeds.bbci.co.uk/news/rss.xml"

# 财经新闻
curl -s "https://feeds.bbci.co.uk/news/business/rss.xml"

# 科技新闻
curl -s "https://feeds.bbci.co.uk/news/technology/rss.xml"

### 路透社

# 世界新闻
curl -s "https://www.reutersagency.com/feed/?best-regions=world&post_type=best"

### NPR（美国视角）

curl -s "https://feeds.npr.org/1001/rss.xml"

### 半岛电视台（全球南方视角）

curl -s "https://www.aljazeera.com/xml/rss/all.xml"

## 解析 RSS

提取标题和描述：

curl -s "https://feeds.bbci.co.uk/news/world/rss.xml" | \
grep -E "<title>|<description>" | \
sed 's/<[^>]*>//g' | \
sed 's/^[ \t]*//' | \
head -30

## 工作流程

### 文字摘要

- 获取 BBC 世界新闻头条
- 可选择补充路透社/NPR 内容
- 摘要关键新闻
- 按地区或主题分类

### 语音摘要

- 创建文字摘要
- 使用 OpenAI TTS 生成语音
- 以音频消息发送

curl -s https://api.openai.com/v1/audio/speech \
  -H "Authorization: Bearer $OPENAI_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "model": "tts-1-hd",
    "input": "<新闻摘要文本>",
    "voice": "onyx",
    "speed": 0.95
  }' \
  --output /tmp/news.mp3

## 输出格式示例

📰 新闻摘要 [日期]

🌍 国际
- [头条 1]
- [头条 2]

💼 财经
- [头条 1]

💻 科技
- [头条 1]

## 最佳实践

- 摘要保持简洁（5-8 条重点新闻）
- 优先报道突发新闻和重大事件
- 语音摘要控制在约 2 分钟以内
- 平衡不同视角（西方 + 全球南方）
- 被问及时标注来源
