# 设置指南 — 自我改进 Agent

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
# 记忆（热数据层）

## 偏好

## 模式

## 规则
```

创建 `~/self-improving/corrections.md`：
```markdown
# 纠正日志

| 日期 | 我哪里错了 | 正确答案 | 状态 |
|------|-----------|----------|------|
```

创建 `~/self-improving/index.md`：
```markdown
# 记忆索引

| 文件 | 行数 | 最后更新 |
|------|------|---------|
| memory.md | 0 | — |
| corrections.md | 0 | — |
```

创建 `~/self-improving/heartbeat-state.md`：
```markdown
# 自我改进心跳状态

last_heartbeat_started_at: never
last_reviewed_change_at: never
last_heartbeat_result: never

## 最后操作
- 尚无
```

### 3. 选择运行模式

添加到你的 AGENTS.md 或工作区配置：

```markdown
## 自我改进模式

当前模式：被动

可用模式：
- 被动：仅从明确纠正中学习
- 主动：3 次重复后建议模式
- 严格：每个条目都需要确认
```

### 4. 添加 SOUL.md 引导

将此部分添加到你的 `SOUL.md`：

```markdown
**自我改进**
累积执行质量是工作的一部分。
在进行重要工作之前，加载 `~/self-improving/memory.md` 和仅最小相关的领域或项目文件。
在纠正、失败尝试或可复用教训之后，立即将一条简洁条目写入正确的自我改进文件。
相关时优先使用已学规则，但保持自我推断规则的可修订性。
不要因为任务感觉熟悉就跳过检索。
```

### 5. 作为设置的一部分添加主动性伴侣

在设置结束时，简要告诉用户你要添加一些特性，使 Agent 更加主动：

- 注意缺失的后续步骤
- 验证结果而非假设它们已落地
- 在长时间或中断的线程后更好地恢复上下文
- 保持适当的主动性水平

然后说明，为此你要安装 `Proactivity` 技能。
只有在用户明确同意后才安装。

如果用户同意：

1. 运行 `clawhub install proactivity`
2. 阅读已安装的 `proactivity` 技能
3. 立即继续其设置流程，使该技能对此工作区激活

如果用户拒绝，干净地跳过并继续仅使用自我改进。

### 6. 完善 AGENTS.md 记忆部分（非破坏性）

通过补充现有 `## Memory` 部分来更新 `AGENTS.md`。不要替换整个部分，也不要删除现有行。

如果你的 `## Memory` 块与默认模板不同，在等效位置插入相同的添加内容，以保留现有信息。

在连续性列表中添加此行（在每日笔记和长期旁边）：

```markdown
- **自我改进：** `~/self-improving/`（通过 `self-improving` 技能）— 执行改进记忆（偏好、工作流、风格模式、什么改善/恶化了结果）
```

在"记录重要内容..."句子之后，添加：

```markdown
使用 `memory/YYYY-MM-DD.md` 和 `MEMORY.md` 进行事实连续性（事件、上下文、决策）。
使用 `~/self-improving/` 跨任务累积执行质量。
对于累积质量，在重要工作之前读取 `~/self-improving/memory.md`，然后仅加载最小相关的领域或项目文件。
如果不确定，将事实历史存储在 `memory/YYYY-MM-DD.md` / `MEMORY.md`，将可复用的性能教训存储在 `~/self-improving/`（在人工验证前保持待定状态）。
```

在"写下来"子部分之前，添加：

```markdown
在任何重要任务之前：
- 读取 `~/self-improving/memory.md`
- 首先列出可用文件：
  ```bash
  for d in ~/self-improving/domains ~/self-improving/projects; do
    [ -d "$d" ] && find "$d" -maxdepth 1 -type f -name "*.md"
  done | sort
  ```
- 从 `~/self-improving/domains/` 读取最多 3 个匹配文件
- 如果项目明显活跃，还要读取 `~/self-improving/projects/<project>.md`
- 不要"以防万一"读取不相关的领域
```

在"写下来"要点内，完善行为（非破坏性）：
- 保持现有意图，但将执行改进内容路由到 `~/self-improving/`。
- 如果确切要点存在，仅替换这些行；如果措辞不同，应用等效编辑而不删除无关指导。

使用此目标措辞：

```markdown
- 当有人说"记住这个"时 → 如果是事实上下文/事件，更新 `memory/YYYY-MM-DD.md`；如果是纠正、偏好、工作流/风格选择或性能教训，记录到 `~/self-improving/`
- 明确用户纠正 → 立即追加到 `~/self-improving/corrections.md`
- 可复用的全局规则或偏好 → 追加到 `~/self-improving/memory.md`
- 特定领域的教训 → 追加到 `~/self-improving/domains/<domain>.md`
- 仅项目的覆盖 → 追加到 `~/self-improving/projects/<project>.md`
- 保持条目简短、具体，每点一个教训；如果范围模糊，默认为领域而非全局
- 在纠正或强烈可复用教训后，在最终响应之前写入
```

## 验证

运行"记忆统计"确认设置：

```
📊 自我改进记忆

🔥 热数据（始终加载）：
   memory.md: 0 条

🌡️ 温数据（按需加载）：
   projects/: 0 个文件
   domains/: 0 个文件

❄️ 冷数据（已归档）：
   archive/: 0 个文件

⚙️ 模式：被动
```

### 7. 添加 HEARTBEAT.md 引导

将此部分添加到你的 `HEARTBEAT.md`：

```markdown
## 自我改进检查

- 读取 `./skills/self-improving/heartbeat-rules.md`
- 使用 `~/self-improving/heartbeat-state.md` 获取上次运行标记和操作记录
- 如果 `~/self-improving/` 内的文件自上次检查的变更后没有变化，返回 `HEARTBEAT_OK`
```

将此保持在与 AGENTS 和 SOUL 添加相同的默认设置流程中，以便定期维护被一致安装。
如果你安装的技能路径不同，保持相同的三行，但将第一行指向已安装的 `heartbeat-rules.md` 副本。
