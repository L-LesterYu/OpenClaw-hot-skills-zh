# Claude Skills 最佳实践指南

本文档总结 Claude 官方 Skills 最佳实践，用于指导开发符合标准的 Skills。

---

## 核心原则

### 1. 简洁至上

Context window 是公共资源，每个 token 都要物有所值。

**默认假设**: Claude 已经很聪明，只添加它不知道的信息。

挑战每一段落：
- "Claude 真的需要这个解释吗？"
- "可以假设 Claude 知道这个吗？"
- "这段内容值得消耗这些 token 吗？"

**简洁示例** (~50 tokens):
```
## Extract PDF text

Use pdfplumber for text extraction:

```python
import pdfplumber

with pdfplumber.open("file.pdf") as pdf:
    text = pdf.pages[0].extract_text()
```
```

### 2. 设置适当的自由度

根据任务脆弱性和可变性决定指令的具体程度：

| 自由度 | 适用场景 | 示例 |
|--------|----------|------|
| **高自由度** | 多种有效方法、依赖上下文判断 | 代码审查流程 |
| **中自由度** | 有首选模式、允许变化 | 生成报告（带模板） |
| **低自由度** | 易错、一致性关键、必须按序 | 数据库迁移 |

**比喻**: 把 Claude 想象成探索路径的机器人
- 窄桥悬崖 → 提供具体护栏和精确指令（低自由度）
- 开阔平原 → 给出大致方向，信任 Claude 找路（高自由度）

### 3. 在所有计划使用的模型上测试

| 模型 | 测试重点 |
|------|----------|
| **Haiku** | 是否提供足够指导？ |
| **Sonnet** | 是否清晰高效？ |
| **Opus** | 是否过度解释？ |

---

## Skill 结构

### YAML Frontmatter 要求

```yaml
---
name: max-64-chars, only-lowercase-numbers-hyphens
description: 非空，max-1024字符，说明做什么+何时使用
---
```

**name 字段规则**:
- 最多 64 字符
- 只能包含小写字母、数字、连字符
- 不能包含 XML 标签
- 禁用保留词: "anthropic", "claude"

**description 字段规则**:
- 必须非空
- 最多 1024 字符
- 不能包含 XML 标签
- **必须用第三人称**（会注入 system prompt）

### 命名约定

推荐使用 **动名词形式 (gerund form)**:

| 推荐 | 可接受 | 避免 |
|------|--------|------|
| `processing-pdfs` | `pdf-processing` | `helper`, `utils` |
| `analyzing-spreadsheets` | `spreadsheet-analysis` | `documents`, `data` |
| `managing-databases` | `process-pdfs` | `tools`, `files` |

### 描述写作技巧

**始终用第三人称**:

- ✅ `"Processes Excel files and generates reports"`
- ❌ `"I can help you process Excel files"`
- ❌ `"You can use this to process Excel files"`

**具体且包含关键术语**:

```
description: Extract text and tables from PDF files, fill forms, merge documents.
Use when working with PDF files or when the user mentions PDFs, forms, or document extraction.
```

---

## 渐进式披露模式

SKILL.md 作为目录，按需加载其他内容。

### 目录结构示例

```
skill/
├── SKILL.md              # 主指令（触发时加载）
├── FORMS.md              # 表单指南（按需）
├── REFERENCE.md          # API 参考（按需）
├── EXAMPLES.md           # 使用示例（按需）
└── scripts/
    ├── analyze_form.py   # 执行脚本，不加载内容
    ├── fill_form.py      # 表单填充脚本
    └── validate.py       # 验证脚本
```

### 关键规则

1. **SKILL.md 主体保持在 500 行以内**
2. **引用只保持一层深度**，避免嵌套引用

**反模式 - 嵌套过深**:
```
# SKILL.md
See [advanced.md](advanced.md)...

# advanced.md
See [details.md](details.md)...

# details.md
Here's the actual information...
```

**正确模式 - 一层引用**:
```
# SKILL.md

**Basic usage**: [instructions in SKILL.md]
**Advanced features**: See [advanced.md](advanced.md)
**API reference**: See [reference.md](reference.md)
**Examples**: See [examples.md](examples.md)
```

3. **长文件（>100 行）顶部添加目录**

```markdown
# API Reference

## Contents
- Authentication and setup
- Core methods (create, read, update, delete)
- Advanced features (batch operations, webhooks)
- Error handling patterns
- Code examples

## Authentication and setup
...
```

### 组织模式

**模式 1: 高级指南 + 引用**

```markdown
---
name: pdf-processing
description: Extracts text and tables from PDF files...
---

# PDF Processing

## Quick start

Extract text with pdfplumber:
```python
import pdfplumber
with pdfplumber.open("file.pdf") as pdf:
    text = pdf.pages[0].extract_text()
```

## Advanced features

**Form filling**: See [FORMS.md](FORMS.md) for complete guide
**API reference**: See [REFERENCE.md](REFERENCE.md) for all methods
**Examples**: See [EXAMPLES.md](EXAMPLES.md) for common patterns
```

**模式 2: 按领域组织**

```
bigquery-skill/
├── SKILL.md (overview and navigation)
└── reference/
    ├── finance.md (revenue, billing metrics)
    ├── sales.md (opportunities, pipeline)
    ├── product.md (API usage, features)
    └── marketing.md (campaigns, attribution)
```

```markdown
# BigQuery Data Analysis

## Available datasets

**Finance**: Revenue, ARR, billing → See [reference/finance.md](reference/finance.md)
**Sales**: Opportunities, pipeline, accounts → See [reference/sales.md](reference/sales.md)
**Product**: API usage, features, adoption → See [reference/product.md](reference/product.md)
**Marketing**: Campaigns, attribution, email → See [reference/marketing.md](reference/marketing.md)

## Quick search

Find specific metrics using grep:

```bash
grep -i "revenue" reference/finance.md
grep -i "pipeline" reference/sales.md
grep -i "api usage" reference/product.md
```
```

**模式 3: 条件细节**

```markdown
# DOCX Processing

## Creating documents

Use docx-js for new documents. See [DOCX-JS.md](DOCX-JS.md).

## Editing documents

For simple edits, modify the XML directly.

**For tracked changes**: See [REDLINING.md](REDLINING.md)
**For OOXML details**: See [OOXML.md](OOXML.md)
```

---

## 工作流和反馈循环

### 复杂任务使用工作流

提供可复制的 checklist:

```markdown
## Research synthesis workflow

Copy this checklist and track your progress:

```
Research Progress:
- [ ] Step 1: Read all source documents
- [ ] Step 2: Identify key themes
- [ ] Step 3: Cross-reference claims
- [ ] Step 4: Create structured summary
- [ ] Step 5: Verify citations
```

**Step 1: Read all source documents**

Review each document in the `sources/` directory. Note the main arguments and supporting evidence.

**Step 2: Identify key themes**

Look for patterns across sources...
```

### 实现反馈循环

**验证器模式**: 运行验证器 → 修复错误 → 重复

```markdown
## Document editing process

1. Make your edits to `word/document.xml`
2. **Validate immediately**: `python ooxml/scripts/validate.py unpacked_dir/`
3. If validation fails:
   - Review the error message carefully
   - Fix the issues in the XML
   - Run validation again
4. **Only proceed when validation passes**
5. Rebuild: `python ooxml/scripts/pack.py unpacked_dir/ output.docx`
6. Test the output document
```

---

## 内容指南

### 避免时效性信息

**反模式** (会过时):
```
If you're doing this before August 2025, use the old API.
After August 2025, use the new API.
```

**正确模式** (使用 `<details>`):
```
## Current method

Use the v2 API endpoint: `api.example.com/v2/messages`

## Old patterns

<details>
<summary>Legacy v1 API (deprecated 2025-08)</summary>

The v1 API used: `api.example.com/v1/messages`

This endpoint is no longer supported.
</details>
```

### 统一术语

| 好的做法 | 坏的做法 |
|----------|----------|
| 始终 "API endpoint" | 混用 "API endpoint", "URL", "route", "path" |
| 始终 "field" | 混用 "field", "box", "element", "control" |
| 始终 "extract" | 混用 "extract", "pull", "get", "retrieve" |

### 模板模式

**严格需求** (使用 ALWAYS):
```markdown
## Report structure

ALWAYS use this exact template structure:

```markdown
# [Analysis Title]

## Executive summary
[One-paragraph overview of key findings]

## Key findings
- Finding 1 with supporting data
- Finding 2 with supporting data
```
```

**灵活需求** (使用 use your judgment):
```markdown
## Report structure

Here is a sensible default format, but use your best judgment based on the analysis:

```markdown
# [Analysis Title]

## Executive summary
[Overview]

## Key findings
[Adapt sections based on what you discover]
```
```

### 示例模式

提供 input/output 对:

```markdown
## Commit message format

Generate commit messages following these examples:

**Example 1:**
Input: Added user authentication with JWT tokens
Output:
```
feat(auth): implement JWT-based authentication

Add login endpoint and token validation middleware
```

**Example 2:**
Input: Fixed bug where dates displayed incorrectly in reports
Output:
```
fix(reports): correct date formatting in timezone conversion

Use UTC timestamps consistently across report generation
```

Follow this style: type(scope): brief description, then detailed explanation.
```

### 条件工作流模式

```markdown
## Document modification workflow

1. Determine the modification type:

   **Creating new content?** → Follow "Creation workflow" below
   **Editing existing content?** → Follow "Editing workflow" below

2. Creation workflow:
   - Use docx-js library
   - Build document from scratch
   - Export to .docx format

3. Editing workflow:
   - Unpack existing document
   - Modify XML directly
   - Validate after each change
   - Repack when complete
```

---

## 反模式

| 反模式 | 正确做法 |
|--------|----------|
| Windows 路径 `scripts\\helper.py` | Unix 路径 `scripts/helper.py` |
| 提供太多选项 "你可以用 A 或 B 或 C..." | 提供默认选项 + 逃生舱 |
| 脚本只是让 Claude 去猜 | 脚本显式处理错误 |

---

## 可执行代码的 Skill

### 解决问题而非推诿

**好的示例 - 显式处理错误**:
```python
def process_file(path):
    """Process a file, creating it if it doesn't exist."""
    try:
        with open(path) as f:
            return f.read()
    except FileNotFoundError:
        # Create file with default content instead of failing
        print(f"File {path} not found, creating default")
        with open(path, 'w') as f:
            f.write('')
        return ''
    except PermissionError:
        # Provide alternative instead of failing
        print(f"Cannot access {path}, using default")
        return ''
```

**坏的示例 - 推给 Claude**:
```python
def process_file(path):
    # Just fail and let Claude figure it out
    return open(path).read()
```

### 配置参数要有文档

避免 "voodoo constants":

**好的示例**:
```python
# HTTP requests typically complete within 30 seconds
# Longer timeout accounts for slow connections
REQUEST_TIMEOUT = 30

# Three retries balances reliability vs speed
# Most intermittent failures resolve by the second retry
MAX_RETRIES = 3
```

**坏的示例**:
```python
TIMEOUT = 47  # Why 47?
RETRIES = 5   # Why 5?
```

### 提供实用脚本

优势：
- 比生成代码更可靠
- 节省 token（无需加载代码内容）
- 节省时间（无需代码生成）
- 确保一致性

**重要区别**: 指令中要明确：
- **执行脚本**: "Run `analyze_form.py` to extract fields"
- **作为参考阅读**: "See `analyze_form.py` for the extraction algorithm"

### 创建可验证的中间输出

plan-validate-execute 模式:
1. 分析 → 创建计划文件 → 验证计划 → 执行 → 验证

**适用场景**: 批量操作、破坏性更改、复杂验证规则、高风险操作。

### 依赖管理

**不要假设包已安装**:

```
**Bad example: Assumes installation**:
"Use the pdf library to process the file."

**Good example: Explicit about dependencies**:
"Install required package: `pip install pypdf`

Then use it:
```python
from pypdf import PdfReader
reader = PdfReader("file.pdf")
```"
```

---

## MCP 工具引用

使用完全限定工具名称: `ServerName:tool_name`

```
Use the BigQuery:bigquery_schema tool to retrieve table schemas.
Use the GitHub:create_issue tool to create issues.
```

---

## 有效 Skill 检查清单

### 核心质量

- [ ] 描述具体且包含关键术语
- [ ] 描述包含做什么 + 何时使用
- [ ] 始终用第三人称描述
- [ ] SKILL.md 主体 < 500 行
- [ ] 无时效性信息（或在 "old patterns" 部分）
- [ ] 术语一致
- [ ] 引用仅一层深
- [ ] 使用渐进式披露

### 代码和脚本

- [ ] 脚本解决问题而非推给 Claude
- [ ] 显式错误处理
- [ ] 无魔法数字（所有值都有说明）
- [ ] 列出所需包
- [ ] 无 Windows 风格路径
- [ ] 关键操作有验证步骤

### 测试

- [ ] 至少 3 个评估场景
- [ ] 在 Haiku/Sonnet/Opus 上测试
- [ ] 真实使用场景测试

---

## 开发流程建议

1. **先不用 Skill 完成任务** → 注意反复提供的信息
2. **用 Claude A 创建 Skill** → 让它帮忙生成结构
3. **用 Claude B 测试** → 新实例测试真实任务
4. **观察并迭代** → 基于实际行为而非假设改进

### 观察 Claude 如何导航 Skill

- **意外探索路径**: 结构可能不够直观
- **遗漏连接**: 链接需要更明确
- **过度依赖某些部分**: 考虑移到主 SKILL.md
- **忽略内容**: 可能不需要或信号不明显

---

## skill-creator 工具对比

参考: `awesome-claude-skills/skill-creator` 是一个标准的 Skill 创建工具。

### ✅ 完全一致的部分

| 项目 | 官方 | skill-creator |
|------|------|---------------|
| **YAML frontmatter** | name + description 必需 | ✅ 相同 |
| **第三人称描述** | "Processes Excel files..." | ✅ 相同 |
| **渐进式披露** | 三级加载 | ✅ 相同 |
| **目录结构** | scripts/, references/ | ✅ 相同（+ assets/）|
| **迭代开发** | 观察→改进→测试 | ✅ 相同 |

### ⚠️ 差异点（无冲突）

| 项目 | 官方 | skill-creator | 说明 |
|------|------|---------------|------|
| **name 长度** | max 64 字符 | 提示 "max 40" | skill-creator 更保守 |
| **命名形式** | 推荐 gerund form | 只要求连字符格式 | skill-creator 更宽松 |
| **SKILL.md 大小** | < 500 行 | < 5k 字词 | 度量方式不同 |

### 🔍 skill-creator 额外验证规则

```python
# quick_validate.py 中的额外验证：
# 1. 不能以连字符开头/结尾
if name.startswith('-') or name.endswith('-'):
    return False, "Name cannot start/end with hyphen"

# 2. 不能有连续连字符
if '--' in name:
    return False, "Name cannot contain consecutive hyphens"

# 3. description 不能有尖括号
if '<' in description or '>' in description:
    return False, "Description cannot contain angle brackets"
```

### 📁 skill-creator 目录结构

```
skill-name/
├── SKILL.md              # 主文件（必需）
├── scripts/              # 可执行代码（可选）
│   └── example.py        # 模板脚本
├── references/           # 文档参考（可选）
│   └── api_reference.md  # 模板文档
└── assets/               # 输出资源（可选）
    └── example_asset.txt # 模板资源
```

### 🛠️ 创建新 Skill 命令

```bash
# 初始化
python scripts/init_skill.py <skill-name> --path <output-directory>

# 验证
python scripts/quick_validate.py <path/to/skill-folder>

# 打包（会自动验证）
python scripts/package_skill.py <path/to/skill-folder> [output-directory]
```

### 📋 skill-creator 6 步流程

| 步骤 | 操作 | 说明 |
|------|------|------|
| 1 | 理解用途 | 明确功能、场景、触发模式 |
| 2 | 规划内容 | 分析重复性工作、脚本化需求 |
| 3 | 初始化 | `init_skill.py` 生成目录结构 |
| 4 | 编辑内容 | 完善 SKILL.md、scripts、references |
| 5 | 打包 | `package_skill.py` 打包成 zip |
| 6 | 迭代优化 | 实战测试、收集反馈 |

### 🎯 skill-creator 独有的 4 种结构模式

init_skill.py 模板中提供：

| 模式 | 适用场景 | 结构 |
|------|----------|------|
| **Workflow-Based** | 顺序流程 | Overview → Decision Tree → Step 1 → Step 2 |
| **Task-Based** | 工具集合 | Overview → Quick Start → Task 1 → Task 2 |
| **Reference/Guidelines** | 标准规范 | Overview → Guidelines → Specifications |
| **Capabilities-Based** | 集成系统 | Overview → Core Capabilities → Feature 1, 2... |

### 📦 assets/ 目录说明

skill-creator 引入的 `assets/` 明确区分：

| 目录 | 用途 | 示例 |
|------|------|------|
| `scripts/` | 执行或参考的代码 | Python/Bash 脚本 |
| `references/` | 加载到上下文的文档 | API 文档、Schema |
| `assets/` | 输出中使用的文件 | 模板、图片、字体、样例代码 |

### 结论

- **无冲突**: skill-creator 完全符合官方规范
- **有增强**: 增加 assets/、更多验证、结构指南
- **可放心使用**: 官方最佳实践的优秀实现 + 实用补充

---

## skill-prompt-generator 优秀实践

参考: `skill-prompt-generator` 是一个高质量的领域专用 Skill 系统。

### 核心架构特点

```
.claude/
├── skills/           # 12个专业领域 Skills
│   ├── domain-classifier/       # 领域分类（智能路由）
│   ├── prompt-master/           # 主控调度
│   ├── intelligent-prompt-generator/  # 人像专家
│   ├── design-master/           # 设计专家
│   └── ...                      # 其他领域专家
└── CLAUDE.md                     # Skill 路由指南
```

| 特点 | 说明 |
|------|------|
| **Skills First** | Skills 作为主逻辑，Python 作为后端支持 |
| **智能路由** | 自动识别领域，调用对应专家 Skill |
| **分层设计** | 用户层 → Skills层 → 引擎层 → 数据层 |
| **专业化** | 每个 Skill 专注一个垂直领域 |

### SKILL.md 撰写高级技巧

#### 1. 结构化文档设计

```markdown
# [Skill 名称]

## 概述
[1-2 句话说明功能]

## 核心能力
[能力列表]

## 工作流程
1. **步骤一**: [具体操作]
2. **步骤二**: [具体操作]
...

## 完整示例
[5+ 个真实场景示例]

## 常见问题
[FAQ]
```

#### 2. 意图构造规范

```markdown
## Intent 结构

必须字段：
- lighting: 光影设置（必选！）
- [其他领域特定字段]

智能补全：
- 自动推导依赖关系
- 应用默认值
- 检测逻辑冲突
```

#### 3. 示例驱动设计

```markdown
## 示例场景

### 场景 1: 基础场景
用户输入: "xxx"
输出: "xxx"

### 场景 2: 复杂场景
用户输入: "xxx"
输出: "xxx"

### 场景 3: 边界情况
用户输入: "xxx"
输出: "xxx"
```

### 框架驱动设计模式

将复杂逻辑抽象为 YAML 配置：

```yaml
# prompt_framework.yaml 结构
framework:
  categories:
    - name: subject      # 主体
      required: true
    - name: facial       # 面部
    - name: styling      # 造型
    - name: lighting     # 光影（必选！）
    - name: scene        # 场景

  # 依赖规则
  dependencies:
    - when:
        scene.era: ancient
      then:
        styling.clothing: traditional_chinese

  # 验证规则
  validations:
    - rule: "古装不能使用现代妆容"
      check: ...
```

### 元素化思维

将知识拆分为可复用元素：

| 元素属性 | 说明 |
|----------|------|
| element_id | 唯一标识 |
| category | 分类（7大类别） |
| keywords | 搜索关键词 |
| reusability_score | 复用性评分（0-1） |
| conflicts_with | 冲突元素 |
| required_combinations | 必须组合 |

### 智能选择策略

**全局最优策略**（优于贪心）：

```python
score = (
    keyword_match * 0.60 +      # 关键词匹配度
    quality_score * 0.30 +      # 元素质量评分
    consistency_bonus * 0.10    # 语义一致性
)
```

### 双轨制系统

| 轨道 | 用途 | 示例 |
|------|------|------|
| **元素级** | 灵活组合 | 从 1140+ 元素中选择 |
| **模板级** | 完整方案 | Apple PPT 模板（12元素） |

### 知识库设计

内置常识和领域知识：

```python
knowledge = {
    # 生物学一致性
    'ethnicity_typical_eyes': {
        'East_Asian': ['black', 'dark brown'],
        'European': ['blue', 'green', 'brown'],
    },
    # 导演风格映射
    'director_styles': {
        'zhang_yimou': ['dramatic', 'shadow', 'rim'],
        'wong_kar_wai': ['moody', 'nostalgic', 'warm'],
    }
}
```

### 高级技巧清单

| 技巧 | 说明 | 适用场景 |
|------|------|----------|
| **智能路由** | 分类器自动选择 Skill | 多领域系统 |
| **依赖推导** | 自动补全相关字段 | 有固定规则的领域 |
| **冲突检测** | 修正逻辑矛盾 | 需要一致性的输出 |
| **语义理解** | 区分属性/风格/场景 | 复杂用户输入 |
| **全局最优** | 多维度评分排序 | 元素选择 |
| **模板系统** | 保存完整设计 | 重复性场景 |
| **学习机制** | 保存历史优化 | 持续改进 |

### 创建高质量 Skills 的关键经验

1. **单一职责** - 每个 Skill 专注一个领域
2. **示例丰富** - 提供 5+ 个完整示例
3. **流程清晰** - 6 步工作流模式
4. **错误友好** - 优雅降级和友好提示
5. **结构化数据** - 元素化 + 评分系统
6. **框架驱动** - YAML 配置分离逻辑
7. **持续迭代** - 基于实际使用优化
