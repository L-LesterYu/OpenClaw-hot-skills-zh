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
