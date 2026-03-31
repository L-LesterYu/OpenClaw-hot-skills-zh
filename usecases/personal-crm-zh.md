# 个人 CRM（客户关系管理）与自动联系人发现

手动记录你见过的人、见面时间以及讨论内容几乎不可能。重要的跟进事项不知不觉就被遗漏，在重要会议前你也忘记了相关背景。

这个工作流自动构建和维护个人 CRM：

- 每日定时任务（cron job）扫描邮件和日历，发现新联系人和互动记录
- 将联系人存储在结构化数据库中，附带关系上下文
- 自然语言查询："我对 [某人] 了解多少？"、"谁需要跟进？"、"我上次和 [某人] 交流是什么时候？"
- 每日会议准备简报：在当天会议前，通过 CRM 和邮件历史研究外部与会者，并提供简报

## 痛点

维护个人关系网络是一项无形的但至关重要的工作。大多数人：
- 忘记重要联系人的详细信息
- 错失跟进时机
- 在重要会议前缺乏背景信息准备
- 人脉关系分散在邮件、日历、聊天记录中，难以统一管理

这个智能体解决方案将你的所有互动数据整合到一个统一的个人 CRM 中，让你始终掌握人脉动态。

## 功能介绍

**自动联系人发现**
- 每日定时扫描 Gmail 和 Google Calendar，自动发现新的联系人和互动记录
- 智能识别会议参与者和邮件收件人，提取关键信息
- 自动更新现有联系人的互动历史和关系状态

**结构化数据存储**
- 将联系人信息存储在 SQLite 数据库中，包含姓名、邮箱、首次接触时间、最后联系时间、互动次数和备注
- 保持数据的一致性和完整性，避免重复记录
- 支持灵活的查询和筛选功能

**智能查询功能**
- 通过自然语言查询联系人信息："我对张三了解多少？"
- 识别需要跟进的联系人："谁需要我本周内联系？"
- 查询特定联系人的互动历史："我上次和李四交流是什么时候？"

**每日会议准备**
- 在每天早上检查当日会议安排
- 为每个外部参会者准备背景资料，包括：
  - 联系人基本信息
  - 上次交流时间和内容摘要
  - 待办事项和跟进提醒
- 通过 Telegram 推送会议准备简报

## 所需技能

- `gog` CLI（用于 Gmail 和 Google Calendar）
- 自定义 CRM 数据库（SQLite 或类似方案），或使用 [crm-query](https://clawhub.ai) 技能（如果可用）
- Telegram 话题用于 CRM 查询

## 如何设置

### 1. 创建 CRM 数据库

```sql
CREATE TABLE contacts (
  id INTEGER PRIMARY KEY,
  name TEXT,
  email TEXT,
  first_seen TEXT,
  last_contact TEXT,
  interaction_count INTEGER,
  notes TEXT
);
```

### 2. 设置 Telegram 话题

在 Telegram 中创建一个名为 "personal-crm" 的话题，用于 CRM 查询。

### 3. 配置 OpenClaw 智能体

以下提示词配置自动联系人发现和每日会议简报：

```text
Run a daily cron job at 6 AM to:
1. Scan my Gmail and Calendar for the past 24 hours
2. Extract new contacts and update existing ones
3. Log interactions (meetings, emails) with timestamps and context

Also, every morning at 7 AM:
1. Check my calendar for today's meetings
2. For each external attendee, search my CRM and email history
3. Deliver a briefing to Telegram with: who they are, when we last spoke, what we discussed, and any follow-up items

When I ask about a contact in the personal-crm topic, search the database and give me everything you know.
```

### 4. 进阶配置选项

**自定义联系人识别规则：**
```text
When extracting contacts from emails and calendar events:
- Prioritize contacts with business email domains (@company.com, @startup.io)
- Group contacts by company/organization when possible
- Flag high-value contacts for special attention
```

**智能跟进提醒：**
```text
For contact follow-up:
- If no contact in 30 days → mark for follow-up
- If promised call/email but not delivered → urgent follow-up
- If important meeting coming up → prepare background briefing
```

### 5. 使用示例

**基础查询：**
```
# 在 personal-crm 话题中
我对王五了解多少？
```

**跟进管理：**
```
# 在 personal-crm 话题中
谁需要我本周内联系？
```

**会议准备：**
```
# 在 personal-crm 话题中
为明天的会议准备李四的背景资料
```

## 关键洞察

**数据源多样性是优势**
- 邮件提供正式沟通记录
- 日历显示会议安排和互动频率
- 聊天记录提供非正式互动上下文
- 整合多个数据源才能获得完整的联系人画像

**自动化需要人工校准**
- 自动提取的联系人信息可能需要人工修正
- 重要的业务关系需要额外标注和跟踪
- 定期清理和更新数据保持 CRM 的准确性

**隐私保护很重要**
- 确保数据库存储符合当地隐私法规
- 对敏感联系人信息进行适当标记
- 定期备份和清理不再需要的数据

## 中国用户适配

### 本地化邮箱和日历支持

对于国内用户，可以使用相应的本地服务替代 Gmail 和 Google Calendar：

| 国际服务 | 国内替代 | 说明 |
|---------|---------|------|
| Gmail | 网易邮箱/QQ邮箱 | 支持 IMAP 协议接入 |
| Google Calendar | 网易邮箱日历/企业微信日历 | 需要通过 API 或第三方工具接入 |
| Telegram | 微信企业号/钉钉 | 用于内部通知和查询 |

### 数据存储建议

对于国内用户，建议：
- 使用阿里云/腾讯云的数据库服务
- 考虑数据本地化存储要求
- 确保符合《个人信息保护法》规定

### 集成方案

**飞书企业版用户：**
- 使用飞书日历和邮箱服务
- 通过飞书机器人进行查询和通知
- 将数据存储在企业内部服务器

**钉钉企业用户：**
- 使用钉钉日历和企业邮箱
- 集成钉钉机器人进行日常操作
- 利用钉钉的审批流程管理重要联系人

## 相关链接

- [gog CLI GitHub](https://github.com/your-org/gog) — Gmail 和 Google Calendar CLI 工具
- [SQLite 官方文档](https://www.sqlite.org/docs.html) — 数据库操作指南
- [Telegram Bot API](https://core.telegram.org/bots/api) — 机器人开发文档
- [个人信息保护法](https://www.gov.cn/zhengce/2021-08/20/content_5631682.htm) — 国内数据合规要求
- [企业微信开放平台](https://open.work.weixin.qq.com/) — 国内企业协作工具

---

**原文链接**：[English Version](https://github.com/AlexAnys/awesome-openclaw-usecases/blob/main/usecases/personal-crm.md)