---
name: ontology-zh
description: 类型化知识图谱，用于结构化智能体记忆和可组合技能。适用于创建/查询实体（人物、项目、任务、事件、文档）、关联相关对象、执行约束检查、规划多步骤操作的图变换，或技能间需要共享状态的场景。触发词："记住"、"我知道关于X什么"、"将X关联到Y"、"显示依赖关系"、实体增删改查，或跨技能数据访问。
---

# 本体知识图谱（Ontology）

一个类型化词汇 + 约束系统，用于将知识表示为可验证的图谱。

## 核心概念

万物皆为**实体**，拥有**类型**、**属性**以及与其他实体的**关系**。每次变更都会在提交前进行类型约束验证。

```
实体: { id, type, properties, relations, created, updated }
关系: { from_id, relation_type, to_id, properties }
```

## 使用场景

| 触发词 | 操作 |
|---------|--------|
| "记住..." | 创建/更新实体 |
| "我知道关于 X 什么？" | 查询图谱 |
| "将 X 关联到 Y" | 创建关系 |
| "显示项目 Z 的所有任务" | 图谱遍历 |
| "什么依赖于 X？" | 依赖查询 |
| 规划多步骤工作 | 建模为图变换 |
| 技能需要共享状态 | 读写本体对象 |

## 核心类型

```yaml
# 智能体与人物
Person: { name, email?, phone?, notes? }
Organization: { name, type?, members[] }

# 工作管理
Project: { name, status, goals[], owner? }
Task: { title, status, due?, priority?, assignee?, blockers[] }
Goal: { description, target_date?, metrics[] }

# 时间与地点
Event: { title, start, end?, location?, attendees[], recurrence? }
Location: { name, address?, coordinates? }

# 信息
Document: { title, path?, url?, summary? }
Message: { content, sender, recipients[], thread? }
Thread: { subject, participants[], messages[] }
Note: { content, tags[], refs[] }

# 资源
Account: { service, username, credential_ref? }
Device: { name, type, identifiers[] }
Credential: { service, secret_ref }  # 切勿直接存储密钥

# 元数据
Action: { type, target, timestamp, outcome? }
Policy: { scope, rule, enforcement }
```

## 存储方式

默认路径：`memory/ontology/graph.jsonl`

```jsonl
{"op":"create","entity":{"id":"p_001","type":"Person","properties":{"name":"张三"}}}
{"op":"create","entity":{"id":"proj_001","type":"Project","properties":{"name":"网站重构","status":"active"}}}
{"op":"relate","from":"proj_001","rel":"has_owner","to":"p_001"}
```

通过脚本或直接文件操作进行查询。对于复杂图谱，可迁移到 SQLite。

### 仅追加规则

在操作现有本体数据或模式时，请**追加/合并**变更而非覆盖文件。这样可以保留历史记录并避免破坏已有定义。

## 工作流

### 创建实体

```bash
python3 scripts/ontology.py create --type Person --props '{"name":"张三","email":"zhangsan@example.com"}'
```

### 查询

```bash
python3 scripts/ontology.py query --type Task --where '{"status":"open"}'
python3 scripts/ontology.py get --id task_001
python3 scripts/ontology.py related --id proj_001 --rel has_task
```

### 关联实体

```bash
python3 scripts/ontology.py relate --from proj_001 --rel has_task --to task_001
```

### 验证

```bash
python3 scripts/ontology.py validate  # 检查所有约束
```

## 约束定义

在 `memory/ontology/schema.yaml` 中定义：

```yaml
types:
  Task:
    required: [title, status]
    status_enum: [open, in_progress, blocked, done]
  
  Event:
    required: [title, start]
    validate: "end >= start if end exists"

  Credential:
    required: [service, secret_ref]
    forbidden_properties: [password, secret, token]  # 强制间接引用

relations:
  has_owner:
    from_types: [Project, Task]
    to_types: [Person]
    cardinality: many_to_one
  
  blocks:
    from_types: [Task]
    to_types: [Task]
    acyclic: true  # 禁止循环依赖
```

## 技能契约

使用本体的技能应声明以下内容：

```yaml
# 在 SKILL.md 的 frontmatter 或头部中声明
ontology:
  reads: [Task, Project, Person]
  writes: [Task, Action]
  preconditions:
    - "Task.assignee 必须存在"
  postconditions:
    - "新创建的 Task 状态为 open"
```

## 将规划建模为图变换

将多步骤规划建模为一系列图操作：

```
规划: "安排团队会议并创建后续任务"

1. 创建 Event { title: "团队同步会", attendees: [p_001, p_002] }
2. 关联 Event -> has_project -> proj_001
3. 创建 Task { title: "准备议程", assignee: p_001 }
4. 关联 Task -> for_event -> event_001
5. 创建 Task { title: "发送会议纪要", assignee: p_001, blockers: [task_001] }
```

每个步骤在执行前都会进行验证。如果违反约束则回滚。

## 集成模式

### 与因果推理集成

将本体变更记录为因果操作日志：

```python
# 在创建/更新实体时，同时记录到因果操作日志
action = {
    "action": "create_entity",
    "domain": "ontology", 
    "context": {"type": "Task", "project": "proj_001"},
    "outcome": "created"
}
```

### 跨技能通信

```python
# 邮件技能创建承诺
commitment = ontology.create("Commitment", {
    "source_message": msg_id,
    "description": "周五前发送报告",
    "due": "2026-01-31"
})

# 任务技能接手处理
tasks = ontology.query("Commitment", {"status": "pending"})
for c in tasks:
    ontology.create("Task", {
        "title": c.description,
        "due": c.due,
        "source": c.id
    })
```

## 快速开始

```bash
# 初始化本体存储
mkdir -p memory/ontology
touch memory/ontology/graph.jsonl

# 创建模式（可选但推荐）
python3 scripts/ontology.py schema-append --data '{
  "types": {
    "Task": { "required": ["title", "status"] },
    "Project": { "required": ["name"] },
    "Person": { "required": ["name"] }
  }
}'

# 开始使用
python3 scripts/ontology.py create --type Person --props '{"name":"张三"}'
python3 scripts/ontology.py list --type Person
```

## 参考文档

- `references/schema.md` — 完整类型定义和约束模式
- `references/queries.md` — 查询语言和遍历示例

## 指令范围

运行时指令操作本地文件（`memory/ontology/graph.jsonl` 和 `memory/ontology/schema.yaml`）并提供 CLI 用法用于创建/查询/关联/验证，这属于正常使用范围。该技能读写工作区文件，首次使用时会自动创建 `memory/ontology` 目录。验证包括属性/枚举/禁止字段检查、关系类型/基数验证、标记为 `acyclic: true` 的关系无环检测以及 Event 的 `end >= start` 检查；其他更高级别的约束可能仅在文档中描述，除非在代码中实现。
