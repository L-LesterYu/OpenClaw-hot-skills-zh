# 心跳规则

使用心跳保持 `~/self-improving/` 有序，而不制造混乱或丢失数据。

## 真实来源

保持工作区 `HEARTBEAT.md` 片段最小化。
将此文件视为自我改进心跳行为的稳定契约。
仅在 `~/self-improving/heartbeat-state.md` 中存储可变运行状态。

## 每次心跳开始时

1. 确保 `~/self-improving/heartbeat-state.md` 存在。
2. 立即以 ISO 8601 格式写入 `last_heartbeat_started_at`。
3. 读取之前的 `last_reviewed_change_at`。
4. 扫描 `~/self-improving/` 查找在该时刻之后更改的文件，排除 `heartbeat-state.md` 本身。

## 如果没有变化

- 设置 `last_heartbeat_result: HEARTBEAT_OK`
- 如果你保留操作日志，附加简短的"无实质性变化"笔记
- 返回 `HEARTBEAT_OK`

## 如果有变化

仅进行保守的组织：

- 如果计数或文件引用偏离，刷新 `index.md`
- 通过合并重复项或总结重复条目来压缩过大的文件
- 仅当目标明确时，将明显错放的笔记移动到正确的命名空间
- 完全保留已确认的规则和明确的纠正
- 仅在审查干净完成后更新 `last_reviewed_change_at`

## 安全规则

- 大多数心跳运行应该什么都不做
- 优先追加、总结或索引修复，而非大规模重写
- 永不删除数据、清空文件或覆盖不确定的文本
- 永不重新组织 `~/self-improving/` 之外的文件
- 如果范围模糊，保持文件不变，改为记录建议的后续操作

## 状态字段

保持 `~/self-improving/heartbeat-state.md` 简单：

- `last_heartbeat_started_at`
- `last_reviewed_change_at`
- `last_heartbeat_result`
- `last_actions`

## 行为标准

心跳的存在是为了保持记忆系统整洁和可信。
如果没有明确违反规则，什么都不做。
