# 智能会议纪要生成器

你花费数小时参加会议、记笔记，试图捕捉所有重要内容。等会议结束时，你已经记了满篇笔记，还需要整理、总结并提取行动项。而下一场会议30分钟后就要开始了。

本用例将 OpenClaw 变成一个智能会议助手，自动参加你的会议、智能记录笔记，并生成可操作的摘要。

## 痛点

传统会议记录存在以下问题：
- 手动记录笔记让人分心，无法专注参与讨论
- 笔记往往不完整或杂乱无章
- 提取行动项和决策耗时费力
- 会议洞察被淹没在大量讨论中
- 跟进行动项变成了额外的工作负担

## 功能介绍

- **自动记录笔记**：OpenClaw 参加你的会议（线上或线下），捕捉关键要点、决策和行动项
- **智能摘要**：生成简洁的摘要，聚焦于成果和关键结论
- **行动项提取**：自动识别并分类行动项，标注负责人和截止日期
- **会议分析**：提供会议效率、发言时间分布和情感分析等洞察
- **跟进自动化**：创建日历事件并发送行动项提醒邮件

## 所需技能

- 日历集成（Google Calendar、Outlook 等）
- 邮件集成（用于跟进）
- `web_fetch`：用于在需要时获取会议记录
- `exec`：用于本地会议录音工具
- 可选：`browser`：用于虚拟会议平台

## 如何设置

### 1. 基础配置

配置 OpenClaw 的日历和邮件集成。创建 AGENTS.md 配置：

```text
## Meeting Notes Assistant

You are my intelligent meeting assistant. Your job is to attend meetings, take smart notes, and generate actionable summaries.

Access:
- My calendar (read access)
- My email (send access for follow-ups)
- Meeting recordings/transcripts when available

Rules:
- Focus on capturing decisions made and action items assigned
- Note who is responsible for each action item and when it's due
- Summarize key takeaways in 3-5 bullet points
- Flag any unresolved issues or follow-up needed
- Send follow-up emails within 1 hour of meeting end
```

> 以上为会议纪要助手的系统提示词：设定了智能会议助手的角色、访问权限和行为规则。

### 2. 会议出席设置

针对虚拟会议，设置会议平台的访问权限：

```text
For virtual meetings (Zoom, Google Meet, Microsoft Teams):
- When I join a meeting, automatically start recording if allowed
- Capture audio transcript using your available tools
- For Zoom meetings, use the Zoom API for transcript access
- For Google Meet, use the official meeting transcript feature
```

> 以上为虚拟会议出席的配置提示词：定义了在不同平台上的自动录音和转写行为。

### 3. 本地会议录音设置

针对线下会议或平台录音不可用的情况：

```text
For local meeting recordings:
- Use `exec` to access local audio recording tools
- Convert audio to text using Whisper or similar transcription service
- Process the transcript to extract key information
- Generate notes and action items from the processed transcript
```

> 以上为本地会议录音的配置提示词：定义了使用本地工具录音、转写并生成纪要的流程。

### 4. 每日会议摘要

设置每日定时任务来回顾会议：

```text
Every day at 5:00 PM:
- Review all meetings from the past 24 hours
- Generate a summary of action items due tomorrow
- Send a brief report of meeting outcomes and next steps
- Flag any meetings that went over time or had low participation
```

> 以上为每日会议摘要的定时任务配置：每天下午5点自动回顾会议并生成摘要报告。

### 5. 会议质量洞察

```text
For each meeting attended, provide:
- Meeting duration vs scheduled time
- Key decisions made (with time stamps)
- Action items with owners and deadlines
- Sentiment analysis of the discussion
- Talking time distribution (who spoke most/least)
- Suggestions for more effective meetings next time
```

> 以上为会议质量洞察的配置提示词：定义了每次会议需提供的分析维度和改进建议。

## 关键洞察

- **AI不会分心**：与人类不同，AI能关注每一个细节，不会遗漏重要内容
- **自动化跟进**：最大的价值在于建立问责机制——AI自动发送跟进邮件和日历事件
- **会议智能**：对会议模式进行长期分析，揭示团队动态和会议效率方面的洞察
- **实时价值**：AI可以在会议进行中提供实时笔记，供参会者参考

## 示例提示词

**开始会议记录：**
> "I'm joining a meeting about Q2 planning with the sales team. Please attend and take notes on decisions and action items."

> 中文说明：即将参加与销售团队的Q2规划会议，请出席并记录决策和行动项。

**处理会议录音：**
> "I have a meeting recording at /Users/username/Documents/meeting-recording.wav. Please transcribe it and generate meeting notes with action items."

> 中文说明：提供一个会议录音文件路径，请转写并生成包含行动项的会议纪要。

**生成周报摘要：**
> "Create a weekly summary of all meetings, showing progress on action items and key outcomes from each discussion."

> 中文说明：创建所有会议的周报摘要，展示行动项进展和每次讨论的关键成果。

**会议效率分析：**
> "Analyze my meeting patterns from this week - which meetings were most productive, and what patterns should I change?"

> 中文说明：分析本周会议模式——哪些会议最高效，哪些模式需要改变。

## 实用建议

- **先用示例会议测试**：在重要会议之前，先用几场会议测试系统
- **自定义行动项分类**：根据团队工作流程调整分类方式
- **正确配置权限**：确保日历和邮件权限已正确设置
- **使用一致的会议格式**：标准化的会议结构更适合AI分析
- **定期审查**：定期审查AI生成的纪要准确性，并根据需要调整提示词

## 相关链接

- [OpenClaw 文档](https://docs.openclaw.ai)
- [日历集成指南](https://docs.openclaw.ai/calendar-integration)
- [邮件设置指南](https://docs.openclaw.ai/email-setup)
- [会议最佳实践](https://hbr.org/2023/05/how-to-run-more-effective-meetings)

---

**注意**：本用例需要正确配置日历和邮件访问权限。请始终确保遵守公司关于会议录音和数据隐私的相关政策。
