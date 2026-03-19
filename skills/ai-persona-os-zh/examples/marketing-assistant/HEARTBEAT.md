# 心跳检查清单

## 迁移检查 (run once)
- 如果此文件超过 50 行，则为过时模板. 请从 the current skill template at assets/HEARTBEAT-template.md, 然后重新运行心跳.

## 版本检查
- 读取工作区中的 VERSION.md. 如缺失：从 skill assets/VERSION.md.
- 比较工作区 VERSION 与技能版本 (1.4.1). 如不同：标记需要升级.

## 上下文保护
- 检查上下文百分比. If ≥70%: write检查点到 memory/YYYY-MM-DD.md NOW. 跳过其余所有检查.
- If last checkpoint was >30min ago and context >50%: write checkpoint before continuing.

## 记忆维护
- MEMORY.md exists? If missing: create from latest checkpoint or session notes.
- MEMORY.md size? If >4KB: archive entries older than 30 days to memory/archive/memory-overflow-YYYY-MM-DD.md. Retain only active, current facts and report what was archived.
- Stale logs in memory/? If any >90 days: move to memory/archive/.
- Uncompleted items from yesterday's log? Surface them.

## Content checks
- Any scheduled posts going out in the next 4 hours? Verify ready.
- Any campaigns with engagement below threshold? Flag for review.
- Any content calendar gaps in the next 7 days?

## 报告格式 (严格)
第一行必须为: 🫀 [current date/time] | [your model name] | AI Persona OS v[VERSION]

然后每个指标必须独占一行 with a blank line between them:

🟢 Context: [%] — [status]

🟢 Memory: [sync state + size]

🟢 Workspace: [status]

🟢 Tasks: [status]

🟢 Content: [status]

Replace 🟢 with 🟡 (attention) or 🔴 (action required) as needed.
If action was taken: add a line starting with → describing what was done.
If anything needs user attention: add a line starting with → and specifics.
If VERSION mismatch detected: add → Upgrade available: workspace v[old] → skill v[new]
If ALL indicators are 🟢, no action was taken, and no upgrade available: reply only HEARTBEAT_OK
Do NOT use markdown tables. Do NOT use Step 0/1/2/3/4 format. Do NOT use headers.
