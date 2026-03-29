# 学习日志

此文件记录纠正、知识缺口和最佳实践。

## 使用方法

当发生以下情况时记录学习内容：
- 用户纠正你（"不对，应该是..."）
- 发现知识过时
- 找到更好的方法
- 发现重复模式

## 格式

```markdown
## [LRN-YYYYMMDD-XXX] category

**Logged**: ISO-8601 时间戳
**Priority**: low | medium | high | critical
**Status**: pending
**Area**: frontend | backend | infra | tests | docs | config

### Summary
学习内容的一行描述

### Details
完整上下文：发生了什么，哪里错了，什么是正确的

### Suggested Action
具体的修复或改进建议

### Metadata
- Source: conversation | error | user_feedback
- Related Files: path/to/file.ext
- Tags: tag1, tag2
- See Also: LRN-20250110-001 (如果与现有条目相关)

---
```

## 示例条目

## [LRN-20250224-001] best_practice

**Logged**: 2025-02-24T09:30:00Z
**Priority**: medium
**Status**: pending
**Area**: docs

### Summary
创建中文版自我改进技能以支持中文用户

### Details
原始技能是英文的，为了更好地服务中文用户，创建了完整的中文版本。包括翻译所有文档、说明和模板，同时保持代码示例和技术术语不变。

### Suggested Action
测试中文版本，确保所有功能正常工作

### Metadata
- Source: user_feedback
- Related Files: SKILL.md
- Tags: translation, localization, chinese

---

## [LRN-20250224-002] best_practice

**Logged**: 2025-02-24T17:36:00+08:00
**Priority**: high
**Status**: pending
**Area**: config

### Summary
使用 git clone 直接安装 ClawHub 技能

### Details
当 clawdhub CLI 工具不可用时，可以直接使用 git clone 将技能仓库克隆到 ~/.openclaw/skills/ 目录下。这是一个有效的替代方案。

**步骤：**
1. 导航到技能目录：`cd ~/.openclaw/skills/`
2. 克隆仓库：`git clone <repo-url> <skill-name>`
3. 验证安装：`ls -la ~/.openclaw/skills/<skill-name>/`

### Suggested Action
将此方法添加到 OpenClaw 文档中，作为技能安装的备用方法

### Metadata
- Source: user_feedback
- Related Files: ~/.openclaw/skills/self-improvement-zh/
- Tags: installation, git, skills, alternative-method
- See Also:

---

## [LRN-20250225-001] best_practice

**Logged**: 2025-02-25T11:27:00+08:00
**Priority**: medium
**Status**: pending
**Area**: docs

### Summary
测试 self-improvement-zh 技能时的学习记录演示

### Details
用户想测试 self-improvement-zh 技能。通过实际操作演示了：
1. 检查 `.learnings/` 目录结构
2. 查看现有的 LEARNINGS.md 和 ERRORS.md
3. 添加新的学习条目

这个技能的核心价值在于：
- 捕获非显而易见的问题解决方案
- 跨会话保持学习记忆
- 可将高频学习提升到 AGENTS.md / SOUL.md / TOOLS.md

### Suggested Action
继续在实际工作中使用此技能，当遇到以下情况时主动记录：
- 命令意外失败
- 用户纠正 AI
- 发现更好的方法

### Metadata
- Source: user_feedback
- Related Files: .learnings/LEARNINGS.md
- Tags: demo, self-improvement, learning
- See Also: LRN-20250224-001

---
