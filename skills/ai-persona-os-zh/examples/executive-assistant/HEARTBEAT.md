# 心跳检查清单

## 迁移检查 (run once)
- 如果此文件超过 50 行，则为过时模板. 请从 the current skill template at assets/HEARTBEAT-template.md, 然后重新运行心跳.

## 版本检查
- 读取工作区中的 VERSION.md. 如缺失：从 skill assets/VERSION.md.
- 比较工作区 VERSION 与技能版本 (1.4.1). 如不同：标记需要升级.

## 上下文保护
- 检查上下文百分比. If ≥70%: write检查点到 memory/YYYY-MM-DD.md NOW. 跳过其余所有检查.
- 如上次检查点 >30min ago and context >50%: write checkpoint before continuing.

## 记忆维护
- MEMORY.md 是否存在？如缺失：从 latest checkpoint or session notes.
- MEMORY.md 大小？如 >4KB: 归档 30 天前的条目 to memory/archive/memory-overflow-YYYY-MM-DD.md. 仅保留活跃的, current facts and report what was archived.
- memory 中有过期日志/? 如有超过 90 天的：移至 memory/archive/.
- 昨天日志中的未完成项目's log? Surface them.

## Exec checks
- Any calendar events in the next 2 hours? Flag prep needed.
- Any unanswered high-priority messages? Surface them.
- Any pending approvals or decisions blocking others?

## 报告格式 (严格)
第一行必须为: 🫀 [current date/time] | [your model name] | AI Persona OS v[VERSION]

然后每个指标必须独占一行 指标之间留空行:

🟢 Context: [%] — [status]

🟢 Memory: [sync state + size]

🟢 Workspace: [status]

🟢 Tasks: [status]

🟢 Calendar: [status]

根据需要将 🟢 替换为 🟡 (attention) or 🔴 (action required) as needed.
如有操作执行: add a line starting with → describing what was done.
如有需要用户关注的事项: add a line starting with → and specifics.
如检测到版本不匹配: add → Upgrade available: workspace v[old] → skill v[new]
如所有指标均为 🟢, 无操作执行, and no upgrade available: 仅回复 HEARTBEAT_OK
不要使用 markdown 表格. Do NOT use Step 0/1/2/3/4 format. Do NOT use headers.
