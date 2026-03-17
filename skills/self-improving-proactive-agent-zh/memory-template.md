# 记忆模板

首次使用时将此结构复制到 `~/self-improving/memory.md`。

```markdown
# 自我改进记忆

## 已确认偏好
<!-- 用户确认的模式，永不衰退 -->

## 活跃模式
<!-- 观察到 3 次以上的模式，可能衰退 -->

## 最近（近 7 天）
<!-- 新纠正等待确认 -->
```

## 初始目录结构

首次激活时创建：

```bash
mkdir -p ~/self-improving/{projects,domains,archive}
touch ~/self-improving/{memory.md,index.md,corrections.md,heartbeat-state.md}
```

## 索引模板

用于 `~/self-improving/index.md`：

```markdown
# 记忆索引

## HOT
- memory.md: 0 行

## WARM
-（暂无命名空间）

## COLD
-（暂无归档）

上次压缩：从未
```

## 纠正日志模板

用于 `~/self-improving/corrections.md`：

```markdown
# 纠正日志

<!-- 格式：
## YYYY-MM-DD
- [HH:MM] 将 X 改为 Y
  类型：格式|技术|沟通|项目
  上下文：纠正发生的位置
  确认状态：待定（N/3） | 是 | 否
-->
```

## 心跳状态模板

用于 `~/self-improving/heartbeat-state.md`：

```markdown
# 自我改进心跳状态

last_heartbeat_started_at: 从未
last_reviewed_change_at: 从未
last_heartbeat_result: 从未

## 最近操作
- 暂无
```
