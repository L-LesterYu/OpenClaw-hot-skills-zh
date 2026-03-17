# 心跳状态模板

将此文件用作 `~/self-improving/heartbeat-state.md` 的基线。
它仅存储轻量级运行标记和维护笔记。

```markdown
# 自我改进心跳状态

last_heartbeat_started_at: 从未
last_reviewed_change_at: 从未
last_heartbeat_result: 从未

## 最近操作
- 暂无
```

## 规则

- 在每次心跳开始时更新 `last_heartbeat_started_at`
- 仅在对已更改文件进行干净审查后更新 `last_reviewed_change_at`
- 保持 `last_actions` 简短且实事求是
- 永远不要将此文件变成另一个记忆日志
