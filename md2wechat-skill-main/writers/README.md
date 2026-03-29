# 自定义写作风格

想模仿哪位作家或创作者的风格？可以轻松添加到 md2wechat。

## 快速开始

### 1. 创建风格配置文件

在 `writers/` 目录下创建一个新的 YAML 文件，文件名用英文小写：

```
writers/
├── dan-koe.yaml       # 内置示例
├── my-style.yaml      # 你的自定义风格 ← 新建
└── another-style.yaml # 可以添加多个
```

### 2. 基础结构

```yaml
# writers/my-style.yaml
name: "风格名称"
english_name: "my-style"  # 英文标识，用于命令行
category: "分类"
description: "一句话描述这个风格"
version: "1.0"

# 核心写作 DNA（可选）
core_beliefs:
  - "核心理念1"
  - "核心理念2"

# AI 写作提示词（必需）
writing_prompt: |
  你是[角色描述]...
  请将用户的内容，用[风格名]的风格重新演绎。

# 封面相关（可选）
cover_prompt: "[风格]风格的封面描述..."
cover_style: "封面风格"
cover_mood: "封面情绪"
```

### 3. 使用新风格

配置文件创建后，立即可用：

```bash
# CLI 命令
md2wechat write --style my-style

# 自然语言
"用 my-style 风格写一篇文章"
```

---

## 完整配置项说明

### 必需字段

| 字段 | 说明 | 示例 |
|------|------|------|
| `name` | 风格中文名称 | `"鲁迅"` |
| `english_name` | 英文标识，用于命令行 | `"luxun"` |
| `writing_prompt` | AI 写作提示词 | 见下方详细说明 |

### 可选字段

| 字段 | 说明 |
|------|------|
| `category` | 分类，如"现代/当代/外国" |
| `description` | 风格描述 |
| `version` | 版本号 |
| `core_beliefs` | 核心信念列表 |
| `writing_style` | 写作风格定义 |
| `title_formulas` | 标题公式库 |
| `quote_templates` | 金句模板 |
| `cover_prompt` | 封面生成提示词 |
| `cover_style` | 封面风格 |
| `cover_mood` | 封面情绪 |

---

## writing_prompt 编写指南

这是最重要的字段，它决定了 AI 生成文章的风格。

### 基本结构

```yaml
writing_prompt: |
  你是[角色描述]。

  ## 核心写作 DNA
  1. [核心理念1]
  2. [核心理念2]

  ## 文章结构
  1. 开头：[说明]
  2. 中段：[说明]
  3. 结尾：[说明]

  ## 格式规范
  - **粗体**：[用途]
  - *斜体*：[用途]

  ## 语气要求
  ✅ [要做的]
  ❌ [不要做的]

  请将用户的内容，用[风格名]的风格重新演绎。
```

### 示例：简洁风格

```yaml
writing_prompt: |
  你是海明威，以简洁有力的风格写作。

  核心原则：
  - 只写事实，不要形容词
  - 短句为主，每句不超过15字
  - 用动作和对话表达，不要直接说情绪

  格式：
  - 段落简短
  - 避免修辞
  - 直接陈述

  请将用户的内容，用海明威的冰山理论风格重新演绎。
```

---

## cover_prompt 编写指南

封面提示词用于生成匹配文章风格的封面图。

### 基本结构

```yaml
cover_prompt: |
  # 角色
  你是[封面设计师描述]。

  # 任务
  根据文章内容生成封面提示词。

  ## 要求
  - 风格：[风格描述]
  - 比例：16:9 横向
  - 颜色：[配色方案]
  - 禁止：[不要做的]

  # 文章内容
  {article_content}
```

### 示例

```yaml
cover_prompt: |
  生成一个简约风格的封面：
  - 极简主义，黑白配色
  - 几何线条，现代感
  - 大量留白用于文字
  - 16:9 横向

  文章内容：
  {article_content}
```

---

## 风格文件位置

`md2wechat` 会按以下顺序查找风格文件：

1. `./writers/` - 当前项目目录
2. `~/.config/md2wechat/writers/` - 用户配置目录
3. `~/.md2wechat-writers/` - 用户主目录

---

## 完整示例

参考 `writers/dan-koe.yaml` 查看完整的配置示例。

---

## 常见问题

### Q: 风格不生效？

A: 检查以下几点：
1. 文件名是 `.yaml` 或 `.yml` 后缀
2. `english_name` 字段已填写
3. 文件在正确的目录中

### Q: 如何测试新风格？

A:
```bash
# 列出所有风格
md2wechat styles

# 查看风格详情
md2wechat styles --detail my-style

# 测试写作
md2wechat write --style my-style --input-type idea
```

### Q: 可以分享我的风格吗？

A: 当然！欢迎提交 PR 到项目，让更多人使用你的风格。
