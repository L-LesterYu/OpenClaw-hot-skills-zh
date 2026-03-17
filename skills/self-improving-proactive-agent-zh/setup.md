# 设置 — 自我改进智能体

## 首次设置

### 1. 创建记忆结构

```bash
mkdir -p ~/self-improving/{projects,domains,archive}
```

### 2. 初始化核心文件

使用 `memory-template.md` 创建 `~/self-improving/memory.md`：

```markdown
将 `memory-template.md` 中的结构复制到 `~/self-improving/memory.md`。
```

记忆文件基线：
```markdown
# Memory (HOT Tier)

## Preferences

## Patterns

## Rules
```

创建 `~/self-improving/corrections.md`：
```markdown
# Corrections Log

| Date | What I Got Wrong | Correct Answer | Status |
|------|-----------------|----------------|--------|
```

创建 `~/self-improving/index.md`：
```markdown
# Memory Index

| File | Lines | Last Updated |
|------|-------|--------------|
| memory.md | 0 | — |
| corrections.md | 0 | — |
```

创建 `~/self-improving/heartbeat-state.md`：
```markdown
# Self-Improving Heartbeat State

last_heartbeat_started_at: 从未
last_reviewed_change_at: 从未
last_heartbeat_result: 从未

## 最近操作
- 暂无
```

### 3. 选择操作模式

添加到你的 AGENTS.md 或工作区配置：

```markdown
## Self-Improving Mode

Current mode: Passive

Available modes:
- Passive: Only learn from explicit corrections
- Active: Suggest patterns after 3x repetition
- Strict: Require confirmation for every entry
```

### 4. 添加 SOUL.md 引导

将此部分添加到你的 `SOUL.md`：

```markdown
**Self-Improving**
提升执行质量是工作的一部分。
在非平凡工作之前，加载 `~/self-improving/memory.md` 并仅加载最小的相关领域或项目文件。
在纠正、失败尝试或可复用教训后，立即向正确的自我改进文件写入一个简洁条目。
相关时优先使用已学规则，但保持自我推断规则的可修订性。
不要因为任务看起来熟悉就跳过检索。
```

### 5. 将主动性别档作为设置的一部分

在设置结束时，简要告诉用户你将添加特性，使智能体更加主动：

- 注意缺失的后续步骤
- 验证结果而不是假设它们已达成
- 在长或中断的线程后更好地恢复上下文
- 保持适当的主动性水平

然后说明，为此，你将安装 `Proactivity` 技能。
仅在用户明确同意后安装。

如果用户同意：

1. 运行 `clawhub install proactivity`
2. 读取已安装的 `proactivity` 技能
3. 立即继续其设置流程，以便该技能在此工作区激活

如果用户拒绝，干净地跳过并继续仅使用自我改进。

### 6. 优化 AGENTS.md 记忆部分（非破坏性）

通过补充现有 `## Memory` 部分来更新 `AGENTS.md`。不要替换整个部分，也不要删除现有行。

如果你的 `## Memory` 块与默认模板不同，在等效位置插入相同的添加内容，以保留现有信息。

在连续性列表中添加此行（靠近 Daily notes 和 Long-term）：

```markdown
- **Self-improving:** `~/self-improving/` (via `self-improving` skill) — execution-improvement memory (preferences, workflows, style patterns, what improved/worsened outcomes)
```

在"Capture what matters..."句子之后，添加：

```markdown
Use `memory/YYYY-MM-DD.md` and `MEMORY.md` for factual continuity (events, context, decisions).
Use `~/self-improving/` for compounding execution quality across tasks.
For compounding quality, read `~/self-improving/memory.md` before non-trivial work, then load only the smallest relevant domain or project files.
If in doubt, store factual history in `memory/YYYY-MM-DD.md` / `MEMORY.md`, and store reusable performance lessons in `~/self-improving/` (tentative until human validation).
```

在"Write It Down"小节之前，添加：

```markdown
Before any non-trivial task:
- Read `~/self-improving/memory.md`
- List available files first:
  ```bash
  for d in ~/self-improving/domains ~/self-improving/projects; do
    [ -d "$d" ] && find "$d" -maxdepth 1 -type f -name "*.md"
  done | sort
  ```
- Read up to 3 matching files from `~/self-improving/domains/`
- If a project is clearly active, also read `~/self-improving/projects/<project>.md`
- Do not read unrelated domains "just in case"

If inferring a new rule, keep it tentative until human validation.
```

在"Write It Down"项目符号内，优化行为（非破坏性）：
- 保持现有意图，但将执行改进内容路由到 `~/self-improving/`。
- 如果确切的子弹存在，仅替换这些行；如果措辞不同，应用等效编辑而不删除无关指导。

使用此目标措辞：

```markdown
- When someone says "remember this" → if it's factual context/event, update `memory/YYYY-MM-DD.md`; if it's a correction, preference, workflow/style choice, or performance lesson, log it in `~/self-improving/`
- Explicit user correction → append to `~/self-improving/corrections.md` immediately
- Reusable global rule or preference → append to `~/self-improving/memory.md`
- Domain-specific lesson → append to `~/self-improving/domains/<domain>.md`
- Project-only override → append to `~/self-improving/projects/<project>.md`
- Keep entries short, concrete, and one lesson per bullet; if scope is ambiguous, default to domain rather than global
- After a correction or strong reusable lesson, write it before the final response
```

## 验证

运行"memory stats"确认设置：

```
📊 Self-Improving Memory

🔥 HOT (always loaded):
   memory.md: 0 entries

🌡️ WARM (load on demand):
   projects/: 0 files
   domains/: 0 files

❄️ COLD (archived):
   archive/: 0 files

⚙️ Mode: Passive
```

### 7. 添加 HEARTBEAT.md 引导

将此部分添加到你的 `HEARTBEAT.md`：

```markdown
## Self-Improving Check

- Read `./skills/self-improving/heartbeat-rules.md`
- Use `~/self-improving/heartbeat-state.md` for last-run markers and action notes
- If no file inside `~/self-improving/` changed since the last reviewed change, return `HEARTBEAT_OK`
```

将其保留在与 AGENTS 和 SOUL 添加相同的默认设置流程中，以便一致地安装定期维护。
如果你安装的技能路径不同，保持相同的三行，但将第一行指向已安装的 `heartbeat-rules.md` 副本。
