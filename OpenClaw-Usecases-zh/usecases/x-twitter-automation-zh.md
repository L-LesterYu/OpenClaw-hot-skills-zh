# X/Twitter 自动化操作

通过自然语言实现完整的 X/Twitter 自动化——发推文、回复、点赞、转推、关注、私信、搜索、提取数据、运行抽奖活动和监控账号，全部在你的 OpenClaw 聊天中完成。

## 痛点

管理 X/Twitter 账号需要在应用、第三方仪表盘和分析工具之间来回切换。运行抽奖活动意味着手动挑选获奖者。提取粉丝、点赞者或转推者需要爬虫脚本。没有一个统一的界面能让你用对话的方式完成所有这些操作。

## 功能介绍

TweetClaw 是一个 OpenClaw 插件，将你的智能体连接到 X/Twitter API。你完全通过聊天进行交互：

- **发布与互动**——撰写推文、回复帖子、点赞、转推、关注/取关、发送私信
- **搜索与提取**——搜索推文和用户，提取粉丝、点赞者、转推者、引用推文者和列表成员
- **抽奖活动**——从推文互动中随机挑选获奖者，支持可配置的筛选条件（最低粉丝数、账号年龄、关键词要求等）
- **监控**——监控账号的新推文或粉丝变化并在有更新时通知你

所有操作通过托管 API 进行——无需浏览器 Cookie、无需爬虫、不会暴露凭据。

## 提示词

**安装插件：**
```text
openclaw plugins install @xquik/tweetclaw
```

**发布推文：**
```text
发布推文："刚刚发布了一个新功能——来试试吧！"
```

**运行抽奖活动：**
```text
从这条推文的转推者中随机抽取 3 名获奖者：https://x.com/username/status/123456789。排除粉丝数少于 50 的账号。
```

**提取数据：**
```text
提取所有点赞了这条推文的用户并导出为 CSV：https://x.com/username/status/123456789
```

**监控账号：**
```text
监控 @elonmusk，他每次发新推文时通知我。
```

## 所需技能

- [@xquik/tweetclaw](https://www.npmjs.com/package/@xquik/tweetclaw) —— 通过 `openclaw plugins install @xquik/tweetclaw` 安装

## 相关链接

- [GitHub 仓库](https://github.com/Xquik-dev/tweetclaw)
- [npm 包](https://www.npmjs.com/package/@xquik/tweetclaw)

---

**原文链接**：[English Version](https://github.com/hesamsheikh/awesome-openclaw-usecases/blob/main/usecases/x-twitter-automation.md)
