# OpenClaw-hot-skills-zh 🇨🇳

> 致力于将 clawhub.ai 平台的优质安全项目转化为中文版 Skills。降低阅读理解门槛，提供更符合中文开发者习惯的开源技能库。

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![GitHub stars](https://img.shields.io/github/stars/L-LesterYu/OpenClaw-hot-skills-zh.svg?style=social)](https://github.com/L-LesterYu/OpenClaw-hot-skills-zh/stargazers)

## 📖 项目简介

本项目旨在为中文开发者提供 ClawHub 平台上优质技能（Skills）的中文版本，让使用 OpenClaw 的开发者能够更轻松地理解和使用各种技能。

### 🎯 项目目标

- 🌏 **降低阅读门槛**：将英文 Skills 翻译为清晰的中文文档
- 🔧 **保持技术准确**：保留代码示例和技术术语的原汁原味
- 📚 **提供详细说明**：增加使用示例和最佳实践
- 🚀 **便于安装使用**：提供完整的安装和使用指南

## 📦 已发布的 Skills

### 1. ontology-zh - 类型化知识图谱

**版本**: v0.1.2
**状态**: ✅ 已发布
**来源**: [ClawHub - Ontology](https://clawhub.ai/kn72dv4fm7ss7swbq47nnpad9x7zy2jh/ontology)

**功能**：
- 🗂️ **类型化实体系统**：支持人员、项目、任务、事件、文档等多种实体类型
- 🔗 **关系管理**：在实体之间建立和管理关系
- ✅ **约束验证**：自动验证实体属性和关系约束
- 🔍 **图谱查询**：强大的查询和遍历功能
- 🔄 **跨技能通信**：支持不同技能之间的数据共享

**核心类型**：
- **智能体与人员**: Person, Organization
- **工作管理**: Project, Task, Goal
- **时间与地点**: Event, Location
- **信息**: Document, Message, Thread, Note
- **资源**: Account, Device, Credential

**安装**：
```bash
cd ~/.openclaw/skills/
git clone https://github.com/L-LesterYu/OpenClaw-hot-skills-zh.git temp-repo
cp -r temp-repo/skills/ontology-zh ./
rm -rf temp-repo
```

**快速开始**：
```bash
# 初始化存储
mkdir -p memory/ontology
touch memory/ontology/graph.jsonl

# 创建实体
python3 scripts/ontology.py create --type Person --props '{"name":"张三","email":"zhangsan@example.com"}'
python3 scripts/ontology.py create --type Project --props '{"name":"网站重新设计","status":"active"}'

# 查询实体
python3 scripts/ontology.py list --type Person
python3 scripts/ontology.py query --type Task --where '{"status":"open"}'

# 建立关系
python3 scripts/ontology.py relate --from proj_001 --rel has_task --to task_001
```

**使用场景**：
- "记住..." → 创建/更新实体
- "我知道什么关于X？" → 查询图谱
- "链接X到Y" → 创建关系
- "显示项目Z的所有任务" → 图遍历

[查看详细文档](./skills/ontology-zh/SKILL.md) | [查询参考](./skills/ontology-zh/references/queries.md) | [模式定义](./skills/ontology-zh/references/schema.md)

---

### 2. gog-zh - Google Workspace 命令行工具

**版本**: v1.0.0
**状态**: ✅ 已发布
**来源**: [ClawHub - Gog](https://clawhub.ai/kn70pywhg0fyz996kpa8xj89s57yhv26/gog)

**功能**：
- 📧 **Gmail 管理**：搜索、发送邮件
- 📅 **日历操作**：查询和管理日历事件
- 💾 **云盘搜索**：搜索 Google Drive 文件
- 👥 **联系人管理**：列出和管理联系人
- 📊 **表格操作**：读取、更新、追加、清除表格数据
- 📄 **文档处理**：导出和查看 Google Docs

**安装**：
```bash
# 1. 安装 gog CLI 工具
brew install steipete/tap/gogcli

# 2. 安装 Skill
cd ~/.openclaw/skills/
git clone https://github.com/L-LesterYu/OpenClaw-hot-skills-zh.git temp-repo
cp -r temp-repo/skills/gog-zh ./
rm -rf temp-repo
```

**使用**：
```bash
# 设置 OAuth 凭证
gog auth credentials /path/to/client_secret.json
gog auth add you@gmail.com --services gmail,calendar,drive,contacts,sheets,docs

# 常用命令示例
gog gmail search 'newer_than:7d' --max 10
gog calendar events <calendarId> --from <iso> --to <iso>
gog sheets get <sheetId> "工作表!A1:D10" --json
```

[查看详细文档](./skills/gog-zh/SKILL.md)

---

### 3. Agent-Browser-zh - AI 浏览器自动化工具

**版本**: v1.0.0
**状态**: ✅ 已发布
**来源**: [ClawHub - Agent Browser](https://clawhub.ai)

**功能**：
- 🌐 **浏览器自动化**：专为 AI 智能体设计的浏览器自动化 CLI
- 🎯 **确定性引用**：使用无障碍树的稳定引用，避免脆弱的 CSS 选择器
- 🔄 **快速架构**：Rust CLI + Node.js 守护进程，首次 ~2秒，后续 ~100毫秒
- 📊 **结构化输出**：JSON 格式输出，完美适配 LLM 解析
- 🚀 **AI 优化工作流**：快照-操作-验证的完整工作流

**核心优势**：
- **确定性引用** - `@e2` 始终指向快照中的同一元素
- **无障碍树** - 语义化页面理解，比 HTML 和截图更适合 AI
- **快速执行** - 浏览器持久化，命令间无需重启
- **复杂流程** - 完美处理多步骤表单和工作流

**安装**：
```bash
# 1. 安装 agent-browser CLI
npm install -g agent-browser
agent-browser install  # 下载 Chromium（约30秒）

# 2. 安装 Skill
cd ~/.openclaw/skills/
git clone https://github.com/L-LesterYu/OpenClaw-hot-skills-zh.git temp-repo
cp -r temp-repo/skills/Agent-Browser-zh ./
rm -rf temp-repo
```

**快速开始**：
```bash
# 导航到网站
agent-browser open https://example.com

# 获取页面结构（AI "看到" 页面）
agent-browser snapshot -i --json

# 执行操作（使用确定性引用）
agent-browser click @e2
agent-browser fill @e3 "test@example.com"

# 验证状态
agent-browser snapshot -i --json

# 完成
agent-browser close
```

**使用场景**：
- 自动化多步骤 Web 工作流
- 填写复杂表单
- 可靠、可重复的浏览器自动化
- 处理动态/现代 Web 应用
- 成本敏感操作（无需视觉 API）

[查看详细文档](./skills/Agent-Browser-zh/SKILL.md) | [快速参考](./skills/Agent-Browser-zh/README.md)

---

### 4. self-improvement-zh - 自我改进技能

**版本**: v1.0.0
**状态**: ✅ 已发布

**功能**：
- 📝 记录学习内容、错误和功能请求
- 🔄 实现持续改进和知识积累
- 📊 结构化的日志管理
- 🔍 支持模式检测和优先级管理

**安装**：
```bash
cd ~/.openclaw/skills/
git clone https://github.com/L-LesterYu/OpenClaw-hot-skills-zh.git temp-repo
cp -r temp-repo/skills/self-improvement-zh ./
rm -rf temp-repo
```

[查看详细文档](./skills/self-improvement-zh/README.md)

---

### 5. proactive-agent-zh - 主动式 Agent

**版本**: v3.1.0
**状态**: ✅ 已发布
**来源**: [ClawHub - Proactive Agent](https://clawhub.ai)

**功能**：
- 🎯 **主动预测** — 无需被询问就预测需求并创造价值
- 💾 **持久记忆** — 通过 WAL 协议和工作缓冲区在上下文丢失中存活
- 🔄 **自我改进** — 持续学习和进化，更好地为您服务
- 🔒 **安全加固** — 技能审查、注入防护、上下文泄漏防护
- ⚡ **不懈机智** — 在放弃前尝试 10 种方法解决问题

**核心协议**：
- **WAL 协议**：预写日志，在响应前捕获所有关键细节
- **工作缓冲区**：在内存压缩的危险区域捕获每次交互
- **压缩恢复**：上下文丢失后的分步恢复机制
- **统一搜索**：在说"不知道"前搜索所有来源

**三大支柱**：
1. **主动** — 预测需求，反向提示，主动签到
2. **持久** — WAL 协议，工作缓冲区，压缩恢复
3. **自我改进** — 自我修复，不懈机智，安全演进

**安装**：
```bash
cd ~/.openclaw/skills/
git clone https://github.com/L-LesterYu/OpenClaw-hot-skills-zh.git temp-repo
cp -r temp-repo/skills/proactive-agent-zh ./
rm -rf temp-repo
```

**快速开始**：
```bash
# 复制资产文件到工作区
cp ~/.openclaw/skills/proactive-agent-zh/assets/*.md ./

# Agent 会检测到 ONBOARDING.md 并主动了解您
# 运行安全审计
bash ~/.openclaw/skills/proactive-agent-zh/scripts/security-audit.sh
```

**使用场景**：
- 让 Agent 从被动响应变为主动预测
- 构建能在会话重启后保持连续性的 AI
- 自动记录决策、修正和重要细节
- 实现"在说不知道前先搜索"的智能行为

[查看详细文档](./skills/proactive-agent-zh/SKILL.md) | [入门流程](./skills/proactive-agent-zh/references/onboarding-flow.md) | [安全模式](./skills/proactive-agent-zh/references/security-patterns.md)

---

### 6. tavily-search-zh - Tavily 网页搜索

**版本**: v1.0.0
**状态**: ✅ 已发布
**来源**: [GitHub - tavily-search-skill](https://github.com/jrr996shujin-png/tavily-search-skill)

**功能**：
- 🔍 **AI 优化搜索**：专为 AI 智能体设计的网页搜索 API
- 🌐 **中文友好**：完美支持中文搜索查询和结果展示
- 📊 **结构化输出**：返回标题、链接、内容摘要的 JSON 格式
- 💰 **免费使用**：Tavily 提供免费套餐，无需信用卡
- 🔄 **Brave API 替代**：解决内置 web_search 需要 Brave API 的问题

**安装**：
```bash
cd ~/.openclaw/skills/
git clone https://github.com/L-LesterYu/OpenClaw-hot-skills-zh.git temp-repo
cp -r temp-repo/skills/tavily-search-zh ./
rm -rf temp-repo
```

**配置 API Key**：
```json
// 编辑 ~/.openclaw/openclaw.json
{
  "skills": {
    "entries": {
      "web-search": {
        "apiKey": "你的_Tavily_API_密钥"
      }
    }
  }
}
```

**快速开始**：
```bash
# 执行搜索
python3 ~/.openclaw/skills/tavily-search-zh/scripts/search.py "今天的热点新闻"

# 输出示例：
# 🔍 正在搜索: 今天的热点新闻
# 🔍 找到 5 条结果：
# 【1】标题...
#     链接: https://...
#     摘要: ...
```

**使用场景**：
- "帮我搜索一下最新的 AI 新闻"
- "查找关于 Python 异步编程的资料"
- "今天有什么热点新闻？"
- "搜索 OpenClaw 的使用教程"

**获取 API Key**：
1. 访问 [Tavily 官网](https://tavily.com) 注册账号
2. 在控制台获取 API 密钥（免费套餐无需信用卡）

[查看详细文档](./skills/tavily-search-zh/SKILL.md)

---

### 7. find-skills-zh - 查找技能

**版本**: v0.1.0
**状态**: ✅ 已发布
**来源**: [ClawHub - Find Skills](https://clawhub.ai)

**功能**：
- 🔍 **技能发现**：帮助用户从开放的智能体技能生态系统中发现技能
- 📦 **智能推荐**：根据用户需求推荐合适的技能包
- 🌐 **Skills CLI 集成**：无缝集成 npx skills 命令行工具
- 💡 **使用指导**：提供完整的安装和使用指南
- 🎯 **场景识别**：自动识别用户需求并匹配相关技能

**核心命令**：
- `npx skills find [query]` - 通过交互方式或关键词搜索技能
- `npx skills add <package>` - 从 GitHub 或其他来源安装技能
- `npx skills check` - 检查技能更新
- `npx skills update` - 更新所有已安装的技能

**安装**：
```bash
cd ~/.openclaw/skills/
git clone https://github.com/L-LesterYu/OpenClaw-hot-skills-zh.git temp-repo
cp -r temp-repo/skills/find-skills-zh ./
rm -rf temp-repo
```

**使用场景**：
- "如何做 X" → 搜索并推荐相关技能
- "找一个技能" → 浏览技能市场
- "是否有能...的技能" → 检查技能可用性
- 扩展智能体能力 → 发现并安装新技能

**技能类别**：
- Web 开发：react, nextjs, typescript, css, tailwind
- 测试：testing, jest, playwright, e2e
- DevOps：deploy, docker, kubernetes, ci-cd
- 文档：docs, readme, changelog, api-docs
- 代码质量：review, lint, refactor, best-practices
- 设计：ui, ux, design-system, accessibility
- 生产力：workflow, automation, git

**浏览技能**：https://skills.sh/

[查看详细文档](./skills/find-skills-zh/SKILL.md)

---

### 8. humanizer-zh - AI 文本人性化工具

**版本**: v2.2.0
**状态**: ✅ 已发布
**来源**: [Wikipedia: Signs of AI writing](https://en.wikipedia.org/wiki/Wikipedia:Signs_of_AI_writing)

**功能**：
- 🔍 **AI 痕迹识别**：检测 24 种 AI 写作模式
- ✏️ **文本人性化**：将 AI 生成的文本重写为更自然的表达
- 🎭 **个性注入**：不仅移除不良模式，更注入真实的个性
- 📚 **全面指南**：基于维基百科 AI 写作迹象指南
- 🌐 **中文优化**：完整汉化的检测和修复规则

**检测模式**：
- **内容模式**：过度强调重要性、促销语言、肤浅分析、模糊归因
- **语言模式**：AI 词汇过度使用、系词回避、否定平行结构
- **风格模式**：破折号/粗体过度使用、表情符号滥用
- **沟通模式**：协作产物、知识截止免责声明、奉承语调

**安装**：
```bash
cd ~/.openclaw/skills/
git clone https://github.com/L-LesterYu/OpenClaw-hot-skills-zh.git temp-repo
cp -r temp-repo/skills/humanizer-zh ./
rm -rf temp-repo
```

**使用场景**：
- 编辑或审查 AI 生成的文本
- 移除文章中的 AI 写作痕迹
- 让内容更加自然和人性化
- 检测文档是否由 AI 生成

**工作流程**：
1. 识别文本中的 AI 模式
2. 重写问题段落
3. 保留原意和语调
4. 注入个性与灵魂
5. 最终反 AI 检查

[查看详细文档](./skills/humanizer-zh/SKILL.md)

---

### 9. frontend-design-zh - 前端设计终极版

**版本**: v0.1.0
**状态**: ✅ 已发布
**来源**: [ClawHub - Frontend Design](https://clawhub.ai)

**功能**：
- 🎨 **独特设计生成**：从纯文本需求创建大胆、令人难忘的设计
- 🚫 **反AI套路美学**：避免通用模板，创造有个性的界面
- 📱 **移动优先响应式**：完整的移动端适配模式
- 📦 **单文件打包**：支持打包为单个 HTML 文件
- 🔧 **双工作流支持**：Vite（纯静态）和 Next.js（Vercel 部署）

**技术栈**：
- React 18 + TypeScript
- Tailwind CSS + shadcn/ui
- Framer Motion 动画库

**核心特性**：
- **设计思维指导**：8种极端基调选择（极简主义、繁复混乱、复古未来等）
- **排版规范**：禁止通用字体，使用独特、有个性的字体选择
- **色彩系统**：70-20-10 法则，坚定选择深色或浅色主题
- **动效策略**：协调的页面加载动画，高影响力时刻设计
- **空间构成**：有目的的不对称，重叠元素，对角线流动

**安装**：
```bash
cd ~/.openclaw/skills/
git clone https://github.com/L-LesterYu/OpenClaw-hot-skills-zh.git temp-repo
cp -r temp-repo/skills/frontend-design-zh ./
rm -rf temp-repo
```

**快速开始**：
```bash
# Vite 工作流（纯静态）
bash scripts/init-vite.sh my-site
cd my-site
npm run dev

# Next.js 工作流（Vercel 部署）
bash scripts/init-nextjs.sh my-site
cd my-site
npm run dev
```

**使用场景**：
- "为 AI 写作工具构建 SaaS 落地页"
- "创建摄影师的作品集网站"
- "设计营销网站或仪表盘"
- "构建任何静态 Web 界面"

**预装组件**：
- Button, Badge, Card, Accordion
- Dialog, NavigationMenu, Tabs
- Sheet, Separator, Avatar, Alert

[查看详细文档](./skills/frontend-design-zh/SKILL.md) | [设计哲学](./skills/frontend-design-zh/references/design-philosophy.md) | [移动端模式](./skills/frontend-design-zh/references/mobile-patterns.md)

---

### 10. github-zh - GitHub CLI 交互工具

**版本**: v1.0.0
**状态**: ✅ 已发布
**来源**: [ClawHub - Github](https://clawhub.ai)

**功能**：
- 🔧 **GitHub CLI 集成**：使用 `gh` 命令行工具与 GitHub 交互
- 📋 **Issue 管理**：创建、查看和管理 GitHub issues
- 🔀 **PR 工作流**：处理拉取请求的完整生命周期
- 🚀 **CI/CD 监控**：检查持续集成状态，查看构建日志
- 🌐 **API 访问**：通过 `gh api` 访问 GitHub REST API 的任何端点

**核心命令**：
- `gh pr checks` - 检查 PR 的 CI 状态
- `gh run list` - 列出最近的 workflow 运行
- `gh run view` - 查看运行记录和失败步骤
- `gh api` - 高级 API 查询
- `gh issue` - Issue 管理命令
- `gh pr` - Pull Request 管理命令

**安装**：
```bash
# 1. 安装 GitHub CLI
# macOS
brew install gh

# Linux (参考: https://github.com/cli/cli/blob/trunk/docs/install_linux.md)
# Windows
winget install GitHub.cli

# 2. 安装 Skill
cd ~/.openclaw/skills/
git clone https://github.com/L-LesterYu/OpenClaw-hot-skills-zh.git temp-repo
cp -r temp-repo/skills/github-zh ./
rm -rf temp-repo

# 3. 认证 GitHub CLI
gh auth login
```

**快速开始**：
```bash
# 列出 issues
gh issue list --repo owner/repo --json number,title

# 查看运行失败的 workflow
gh run view <run-id> --repo owner/repo --log-failed

# 创建 PR
gh pr create --repo owner/repo --title "标题" --body "描述"

# API 查询
gh api repos/owner/repo/pulls/55 --jq '.title, .state'
```

**使用场景**：
- "检查这个 PR 的 CI 状态"
- "查看最近的构建失败"
- "列出仓库的所有 issues"
- "创建新的 pull request"
- "获取 PR 的详细信息"

**提示**：
- 使用 `--json` 和 `--jq` 进行灵活的数据提取
- 始终指定 `--repo owner/repo` 当不在 git 目录中时
- 使用 `--help` 查看任何命令的详细帮助

[查看详细文档](./skills/github-zh/SKILL.md)

---

### 11. clawddocs-zh - Clawdbot 文档专家

**版本**: v5.0.0
**状态**: ✅ 已发布
**来源**: [ClawHub - Clawdbot Documentation Expert](https://clawhub.ai)

**功能**：
- 🎯 **决策树导航**：智能识别用户需求，快速定位相关文档
- 🔍 **搜索脚本**：站点地图、关键词搜索、全文索引（支持 qmd）
- 📖 **文档获取**：快速获取特定文档内容
- 📊 **版本跟踪**：监控文档变更，对比历史版本
- ⚙️ **配置片段**：提供即用型配置模式，涵盖所有 Clawdbot 功能

**文档类别**：
- 🚀 入门指南（`/start/`）
- 🔧 网关与运维（`/gateway/`）
- 💬 提供商（`/providers/`）
- 🧠 核心概念（`/concepts/`）
- 🛠️ 工具（`/tools/`）
- ⚡ 自动化（`/automation/`）
- 💻 命令行（`/cli/`）
- 📱 平台（`/platforms/`）

**安装**：
```bash
cd ~/.openclaw/skills/
git clone https://github.com/L-LesterYu/OpenClaw-hot-skills-zh.git temp-repo
cp -r temp-repo/skills/clawddocs-zh ./
rm -rf temp-repo
```

**快速开始**：
```bash
# 查看所有文档
./scripts/sitemap.sh

# 按关键词搜索
./scripts/search.sh discord

# 获取特定文档
./scripts/fetch-doc.sh gateway/configuration

# 查看最近更新
./scripts/recent.sh 7

# 全文搜索（需要 qmd）
./scripts/build-index.sh search "webhook retry"
```

**使用场景**：
- "如何设置 Discord？"
- "为什么我的机器人不工作？"
- "如何配置定时任务？"
- "Clawdbot 有什么新功能？"
- "给我一个 Webhook 配置示例"

**提示**：
- 优先使用缓存的站点地图（1 小时 TTL）
- 对于复杂问题，使用全文索引搜索
- 查看 `snippets/common-configs.md` 获取配置示例
- 回答时链接到文档：`https://docs.clawd.bot/<路径>`

[查看详细文档](./skills/clawddocs-zh/SKILL.md)

---

### 12. youtube-watcher-zh - YouTube 视频字幕提取

**版本**: v1.0.0
**状态**: ✅ 已发布
**来源**: [ClawHub - YouTube Watcher](https://clawhub.ai)

**功能**：
- 📺 **字幕提取**：从 YouTube 视频获取字幕文本
- 📝 **内容摘要**：用于视频总结和内容分析
- ❓ **问答支持**：基于字幕内容回答用户问题
- 🔍 **信息提取**：从视频中提取特定信息
- 🌐 **中文优化**：完整的中文文档和使用说明

**核心能力**：
- 支持隐藏字幕（CC）和自动生成字幕
- 快速获取视频文本内容
- 适用于有字幕的 YouTube 视频
- 简单易用的命令行接口

**安装**：
```bash
# 1. 安装 yt-dlp 工具
brew install yt-dlp  # macOS
# 或
pip install yt-dlp   # Linux/Windows

# 2. 安装 Skill
cd ~/.openclaw/skills/
git clone https://github.com/L-LesterYu/OpenClaw-hot-skills-zh.git temp-repo
cp -r temp-repo/skills/youtube-watcher-zh ./
rm -rf temp-repo
```

**快速开始**：
```bash
# 获取视频字幕
python3 ~/.openclaw/skills/youtube-watcher-zh/scripts/get_transcript.py "https://www.youtube.com/watch?v=VIDEO_ID"

# 输出示例：
# [字幕文本内容]
```

**使用场景**：
- "总结这个 YouTube 视频"
- "这个视频讲了什么？"
- "从视频中提取关于 X 的信息"
- "观看 youtube"
- "视频字幕"

**注意事项**：
- 需要安装 `yt-dlp` 并在 PATH 中可用
- 仅适用于有字幕的视频
- 如果视频没有字幕，脚本将返回错误

[查看详细文档](./skills/youtube-watcher-zh/SKILL.md)

---

### 13. freeride-zh - OpenRouter 免费AI管理

**版本**: v1.0.1
**状态**: ✅ 已发布
**来源**: [ClawHub - FreeRide](https://clawhub.ai)

**功能**：
- 🆓 **免费 AI 访问**：配置 OpenClaw 使用 OpenRouter 的免费 AI 模型
- 🎯 **智能模型选择**：自动选择最佳免费模型作为主模型
- 🔄 **故障切换**：添加备用模型，在限流时自动切换不中断服务
- 📊 **质量排名**：按模型质量和可用性自动排名
- ⚙️ **配置管理**：自动更新 openclaw.json 配置文件

**核心命令**：
- `freeride auto` - 配置最佳免费模型 + 备用模型
- `freeride list` - 查看可用的免费模型
- `freeride switch <model>` - 切换到指定模型
- `freeride status` - 检查当前配置
- `freeride fallbacks` - 仅更新备用模型

**安装**：
```bash
cd ~/.openclaw/skills/
git clone https://github.com/L-LesterYu/OpenClaw-hot-skills-zh.git temp-repo
cp -r temp-repo/skills/freeride-zh ./
rm -rf temp-repo

# 安装 Python 依赖
cd ~/.openclaw/skills/freeride-zh
pip install -e .
```

**快速开始**：
```bash
# 1. 设置 OpenRouter API Key
export OPENROUTER_API_KEY="sk-or-v1-..."

# 2. 配置免费模型
freeride auto

# 3. 重启网关应用更改
openclaw gateway restart

# 验证：用户发送 /status 查看当前模型
```

**使用场景**：
- "我想使用免费的 AI 模型"
- "配置 OpenRouter 免费模型"
- "添加备用模型避免限流"
- "查看有哪些免费模型可用"
- "切换到特定的免费模型"

**配置说明**：
FreeRide 只更新 `~/.openclaw/openclaw.json` 中的模型配置：
- `agents.defaults.model.primary` — 主模型（如 `openrouter/qwen/qwen3-coder:free`）
- `agents.defaults.model.fallbacks` — 备用模型列表
- `agents.defaults.models` — 模型白名单

其他配置保持不变，安全可靠。

**监控器（可选）**：
```bash
freeride-watcher --daemon    # 持续监控并自动轮换
freeride-watcher --rotate    # 立即强制轮换
freeride-watcher --status    # 检查轮换历史
```

**获取 API Key**：
1. 访问 [OpenRouter](https://openrouter.ai/) 注册账号
2. 在 [API Keys](https://openrouter.ai/keys) 页面获取免费密钥

**注意事项**：
- 需要设置 `OPENROUTER_API_KEY` 环境变量
- 配置更改后需运行 `openclaw gateway restart`
- 免费模型有速率限制，建议配置备用模型

[查看详细文档](./skills/freeride-zh/SKILL.md)

---

### 14. weather-zh - 天气查询

**版本**: v1.0.0
**状态**: ✅ 已发布
**来源**: [ClawHub - Weather](https://clawhub.ai)
**首页**: https://wttr.in/:help

**功能**：
- 🌤️ **实时天气**：获取任意城市的当前天气状况
- 📅 **天气预报**：查看未来几天的天气预报
- 🌐 **多服务支持**：使用 wttr.in（主要）和 Open-Meteo（备用）两大免费服务
- 🔒 **无需密钥**：完全免费，无需 API 密钥
- 🇨🇳 **中文支持**：支持中文城市名和天气描述

**核心命令**：
```bash
# 快速查询（简洁格式）
curl -s "wttr.in/Beijing?format=3"
# 输出: Beijing: ⛅️ +8°C

# 详细预报
curl -s "wttr.in/Shanghai?T&lang=zh"

# 带湿度、风速
curl -s "wttr.in/Guangzhou?format=%l:+%c+%t+%h+%w"

# JSON 格式（程序化调用）
curl -s "https://api.open-meteo.com/v1/forecast?latitude=31.2&longitude=121.5&current_weather=true"
```

**格式代码说明**：
- `%c` - 天气状况图标
- `%t` - 温度
- `%h` - 湿度
- `%w` - 风速风向
- `%l` - 地点名称
- `%m` - 月相

**安装**：
```bash
cd ~/.openclaw/skills/
git clone https://github.com/L-LesterYu/OpenClaw-hot-skills-zh.git temp-repo
cp -r temp-repo/skills/weather-zh ./
rm -rf temp-repo
```

**快速开始**：
```bash
# 查询北京天气
curl -s "wttr.in/Beijing?format=3&lang=zh"

# 查询上海三天预报
curl -s "wttr.in/Shanghai?lang=zh"

# 查询广州完整信息
curl -s "wttr.in/Guangzhou?T&lang=zh"
```

**使用场景**：
- "今天天气怎么样？"
- "北京现在的温度"
- "明天会下雨吗？"
- "查询上海的天气预报"

**注意事项**：
- wttr.in 支持中文字符，但建议使用英文城市名或拼音
- 服务免费但可能有访问频率限制
- 如需更稳定的服务，可使用 Open-Meteo API

[查看详细文档](./skills/weather-zh/SKILL.md)

---

### 15. auto-updater-zh - 自动更新技能

**版本**: v1.0.0
**状态**: ✅ 已发布
**来源**: [ClawHub - Auto-Updater](https://clawhub.ai)

**功能**：
- 🔄 **自动更新**：每日自动检查更新
- 📦 **全面覆盖**：同时更新 Clawdbot 核心和所有已安装技能
- 📢 **用户通知**：发送包含所有更改的摘要消息
- 🔒 **安全更新**：使用 ClawHub 的版本控制确保安全

**核心能力**：
- 通过 cron 作业每日运行一次
- 检查 ClawHub 上已安装技能的可用更新
- 如果有新版本可用，更新 Clawdbot 核心
- 应用所有待处理的更新
- 向用户发送摘要消息

**安装**：
```bash
cd ~/.openclaw/skills/
git clone https://github.com/L-LesterYu/OpenClaw-hot-skills-zh.git temp-repo
cp -r temp-repo/skills/auto-updater-zh ./
rm -rf temp-repo
```

**快速开始**：
```bash
# 安装后，技能会自动设置 cron 作业
# 每日在预定时间自动运行更新检查

# 手动触发更新检查
clawhub update --all --no-input
```

**使用场景**：
- 自动保持 Clawdbot 和所有技能为最新版本
- 无需手动检查更新
- 接收每日更新摘要通知

**更新摘要包含**：
- 更新的技能数量
- 每个技能的版本变更
- 遇到的任何错误或警告

**故障排除**：
如果更新失败：
1. 检查互联网连接
2. 验证 ClawHub API 是否可访问
3. 查看更新摘要中的错误消息
4. 尝试手动更新：`clawhub update --all --force`

**注意事项**：
- 更新会自动但安全地应用
- 主要版本更新可能需要手动审核
- 该技能遵守 ClawHub 的速率限制
- 失败的更新会被记录以供故障排除

**要求**：
- 已安装 ClawHub CLI
- 互联网连接
- ClawHub API 访问权限

[查看详细文档](./skills/auto-updater-zh/SKILL.md)

---

### 16. api-gateway-zh - API 网关

**版本**: v1.0.58
**状态**: ✅ 已发布
**来源**: [ClawHub - API Gateway](https://clawhub.ai)
**提供商**: [Maton](https://maton.ai)

**功能**：
- 🌐 **100+ API 连接**：支持 Google Workspace、Microsoft 365、GitHub、Notion、Slack、Airtable、HubSpot 等主流服务
- 🔐 **托管 OAuth**：无需管理 OAuth 流程，Maton 自动处理认证和令牌刷新
- 🔄 **透传代理**：直接调用原生 API 端点，保持完整的 API 功能
- 🛡️ **安全隔离**：MATON_API_KEY 仅用于身份验证，每个服务需单独授权
- 📊 **统一接口**：通过网关统一访问所有第三方 API

**核心能力**：
- 连接 100+ 主流第三方服务
- 托管 OAuth 2.0 认证流程
- 自动处理令牌刷新和过期
- 支持原生 API 端点调用
- 细粒度的权限控制

**安装**：
```bash
cd ~/.openclaw/skills/
git clone https://github.com/L-LesterYu/OpenClaw-hot-skills-zh.git temp-repo
cp -r temp-repo/skills/api-gateway-zh ./
rm -rf temp-repo
```

**快速开始**：
```bash
# 1. 获取 Maton API Key
# 访问 https://maton.ai 注册账号并获取 API 密钥

# 2. 设置环境变量
export MATON_API_KEY="your-maton-api-key"

# 3. 连接服务
# 访问 Maton 控制台授权你需要的服务（如 Gmail、Slack、GitHub 等）

# 4. 调用 API（示例：发送 Slack 消息）
python <<'EOF'
import urllib.request, os, json
data = json.dumps({'channel': 'C0123456', 'text': 'Hello from gateway!'}).encode()
req = urllib.request.Request('https://gateway.maton.ai/slack/api/chat.postMessage', data=data, method='POST')
req.add_header('Authorization', f'Bearer {os.environ["MATON_API_KEY"]}')
req.add_header('Content-Type', 'application/json')
print(json.dumps(json.load(urllib.request.urlopen(req)), indent=2))
EOF
```

**支持的服务（部分列表）**：
- **协作工具**: Slack, Discord, Microsoft Teams
- **项目管理**: Notion, Asana, Trello, Jira
- **开发工具**: GitHub, GitLab, Bitbucket
- **办公套件**: Google Workspace, Microsoft 365
- **CRM**: Salesforce, HubSpot, Zendesk
- **数据库**: Airtable, Supabase
- **存储**: Google Drive, Dropbox, OneDrive

**使用场景**：
- "发送 Slack 消息"
- "创建 GitHub issue"
- "读取 Google Sheets 数据"
- "更新 Notion 页面"
- "发送 Gmail 邮件"

**API 调用格式**：
```
https://gateway.maton.ai/{app}/{native-api-path}
```

例如：
- Slack: `https://gateway.maton.ai/slack/api/chat.postMessage`
- Gmail: `https://gateway.maton.ai/google-mail/gmail/v1/users/me/messages`
- GitHub: `https://gateway.maton.ai/github/repos/owner/repo/issues`

**安全说明**：
- MATON_API_KEY 仅用于 Maton 身份验证
- 每个服务需要单独的 OAuth 授权
- 访问仅限于用户已授权的服务
- 完整的审计日志和访问控制

**获取 API Key**：
1. 访问 [Maton 官网](https://maton.ai) 注册账号
2. 在控制台获取 API 密钥
3. 在 Maton 平台连接你需要的第三方服务
4. 每个服务都需要通过 OAuth 授权流程

**注意事项**：
- 需要先在 Maton 平台授权服务才能使用
- API 调用遵循各服务的速率限制
- 支持 GET、POST、PUT、DELETE 等 HTTP 方法
- 请求和响应格式遵循各服务的原生 API 规范

[查看详细文档](./skills/api-gateway-zh/SKILL.md) | [Maton 官网](https://maton.ai)

---

### 17. openai-whisper-zh - 语音转文字工具

**版本**: v1.0.0
**状态**: ✅ 已发布
**来源**: [ClawHub - Openai Whisper](https://clawhub.ai)
**首页**: https://openai.com/research/whisper

**功能**：
- 🎙️ **本地转录**：使用 Whisper CLI 在本地进行语音转文字
- 🔒 **无需 API**：完全本地化处理，无需 API 密钥
- 🎵 **多格式支持**：支持 MP3、M4A、WAV 等多种音频格式
- 🌐 **翻译功能**：支持将音频翻译为其他语言
- 📝 **多输出格式**：支持 TXT、SRT、VTT 等多种输出格式

**核心命令**：
```bash
# 转录音频文件
whisper /path/audio.mp3 --model medium --output_format txt --output_dir .

# 翻译音频内容
whisper /path/audio.m4a --task translate --output_format srt

# 使用不同模型
whisper /path/audio.wav --model small   # 快速
whisper /path/audio.wav --model large   # 准确
```

**模型说明**：
- `tiny` - 最快，准确性较低
- `base` - 平衡速度和准确性
- `small` - 较好的准确性
- `medium` - 高准确性
- `large` - 最高准确性，最慢
- `turbo` - 默认模型，优化的速度/准确性平衡

**安装**：
```bash
# 1. 安装 Whisper CLI（选择一种方式）
# macOS
brew install openai-whisper

# Linux/Windows
pip install openai-whisper

# 2. 安装 Skill
cd ~/.openclaw/skills/
git clone https://github.com/L-LesterYu/OpenClaw-hot-skills-zh.git temp-repo
cp -r temp-repo/skills/openai-whisper-zh ./
rm -rf temp-repo
```

**快速开始**：
```bash
# 转录音频文件为文本
whisper meeting.mp3 --model medium --output_format txt

# 生成字幕文件
whisper video.m4a --output_format srt

# 翻译外语音频
whisper foreign.mp3 --task translate
```

**使用场景**：
- 会议录音转文字
- 视频字幕生成
- 播客内容转录
- 外语音频翻译
- 语音笔记整理

**注意事项**：
- 首次运行时会自动下载模型到 `~/.cache/whisper`
- 模型文件较大（几百 MB 到几 GB）
- 较大的模型提供更高的准确性但速度较慢
- 建议使用 `turbo` 或 `medium` 模型作为平衡选择

**资源需求**：
- Python 3.8+
- PyTorch（首次使用时自动安装）
- 足够的磁盘空间存储模型

[查看详细文档](./skills/openai-whisper-zh/SKILL.md)

---

### 18. nano-banana-pro-zh - 图像生成与编辑

**版本**: v1.0.1
**状态**: ✅ 已发布
**来源**: [ClawHub - Nano Banana Pro](https://clawhub.ai)

**功能**：
- 🎨 **图像生成**：使用 Gemini 3 Pro Image API 生成高质量图像
- ✏️ **图像编辑**：支持对现有图像进行编辑和修改
- 📐 **多分辨率**：支持 1K/2K/4K 三种分辨率选项
- 🔄 **迭代工作流**：草稿 → 迭代 → 最终版本的高效流程
- 🎯 **提示词模板**：提供高成功率的提示词模板

**核心能力**：
- 文生图（Text-to-Image）
- 图生图（Image-to-Image）
- 多种分辨率选择
- 快速迭代反馈循环
- 精准的提示词处理

**安装**：
```bash
cd ~/.openclaw/skills/
git clone https://github.com/L-LesterYu/OpenClaw-hot-skills-zh.git temp-repo
cp -r temp-repo/skills/nano-banana-pro-zh ./
rm -rf temp-repo
```

**快速开始**：
```bash
# 生成新图像（1K 草稿）
uv run ~/.codex/skills/nano-banana-pro-zh/scripts/generate_image.py \
  --prompt "宁静的日式花园，樱花飘落" \
  --filename "2025-11-23-14-23-05-japanese-garden.png" \
  --resolution 1K

# 生成高分辨率最终版本
uv run ~/.codex/skills/nano-banana-pro-zh/scripts/generate_image.py \
  --prompt "宁静的日式花园，樱花飘落" \
  --filename "2025-11-23-14-23-05-japanese-garden-final.png" \
  --resolution 4K

# 编辑现有图像
uv run ~/.codex/skills/nano-banana-pro-zh/scripts/generate_image.py \
  --prompt "让天空更具戏剧性，添加暴风云" \
  --filename "2025-11-23-14-25-30-dramatic-sky.png" \
  --input-image "original-photo.jpg" \
  --resolution 2K
```

**分辨率选项**：
- **1K**（默认）- 约 1024px，快速迭代
- **2K** - 约 2048px，中等质量
- **4K** - 约 4096px，高质量最终版

**使用场景**：
- "生成一张日式花园的图片"
- "把这个图片的天空改成暴风雨"
- "创建一个 4K 分辨率的机器人图像"
- "编辑照片，添加彩虹"

**工作流建议**：
1. **草稿阶段**（1K）：快速生成，迭代提示词
2. **调整阶段**：微调细节，确定方向
3. **最终阶段**（4K）：生成高质量最终版本

**提示词模板**：
- 生成："创建图像：<主体>。风格：<风格>。构图：<镜头>。光线：<光线>。背景：<背景>。"
- 编辑："仅更改：<单个更改>。完全保持：主体、构图、姿势、光线。"

**API 密钥配置**：
```bash
# 方式 1：环境变量
export GEMINI_API_KEY="your-api-key"

# 方式 2：命令行参数
--api-key "your-api-key"
```

**获取 API Key**：
1. 访问 [Google AI Studio](https://aistudio.google.com/)
2. 创建 API 密钥
3. 启用 Gemini API 访问

**注意事项**：
- 建议先用 1K 分辨率迭代，确定提示词后再生成 4K
- 文件名格式建议：`yyyy-mm-dd-hh-mm-ss-描述.png`
- 编辑时确保输入图像路径正确
- API 有速率限制，建议添加错误处理

[查看详细文档](./skills/nano-banana-pro-zh/SKILL.md)

---

### 19. notion-skill-zh - Notion 操作工具

**版本**: v1.0.0
**状态**: ✅ 已发布
**来源**: [ClawHub - Notion](https://clawhub.ai)
**首页**: https://developers.notion.com

**功能**：
- 📄 **页面管理**：读取、创建和更新 Notion 页面
- 📝 **内容操作**：向页面追加内容块（支持 Markdown）
- 🗃️ **数据库操作**：查询、创建和更新数据库行
- 🔍 **架构检查**：安全地检查和应用数据库架构变更
- 🔒 **安全优先**：优先使用追加和更新操作，避免破坏性操作

**核心能力**：
- 通过官方 Notion API 操作页面和数据库
- 支持多个配置文件（个人/工作）
- 完整的 CRUD 操作支持
- 架构变更前的差异检查
- 遵循 Notion API 速率限制

**安装**：
```bash
# 1. 安装 Notion CLI（选择一种方式）
# Node.js 版本
npm install -g notion-cli

# Python 版本
pip install notion-cli-py

# 2. 安装 Skill
cd ~/.openclaw/skills/
git clone https://github.com/L-LesterYu/OpenClaw-hot-skills-zh.git temp-repo
cp -r temp-repo/skills/notion-skill-zh ./
rm -rf temp-repo
```

**快速开始**：
```bash
# 1. 创建 Notion 集成
# 访问 https://www.notion.so/my-integrations 创建集成
# 复制内部集成令牌

# 2. 设置环境变量
export NOTION_API_KEY=secret_xxx

# 3. 将集成与页面/数据库共享
# 在 Notion 中，点击页面/数据库的"..."菜单
# 选择"添加连接" → 选择你的集成

# 4. 常用操作
# 读取页面
notion-cli page get <page_id>

# 追加内容
notion-cli block append <page_id> --markdown "# 标题\n内容"

# 查询数据库
notion-cli db query <database_id> --filter '{"property":"状态","select":{"equals":"进行中"}}'

# 创建数据库行
notion-cli page create --database <database_id> --props '{"名称":{"title":[{"text":{"content":"新任务"}}]}}'
```

**配置文件**：
```bash
# 默认配置文件：personal
# 切换到工作配置
export NOTION_PROFILE=work
```

**使用场景**：
- "读取这个 Notion 页面的内容"
- "在我的任务数据库中创建一个新任务"
- "查询所有状态为进行中的项目"
- "向会议记录页面追加今天的笔记"
- "更新这个数据库行的状态"

**安全注意事项**：
- Notion API 有速率限制，批量操作需谨慎
- 优先使用追加和更新操作，避免破坏性操作
- ID 是不透明的，需要明确存储，不要从 URL 推断
- 架构变更前始终先检查差异

**架构变更流程**：
```bash
# 1. 检查差异
notion-cli db schema diff <database_id> --desired <json>

# 2. 确认后应用
notion-cli db schema apply <database_id> --desired <json>
```

**注意事项**：
- 需要先创建 Notion 集成并获取 API 密钥
- 必须将集成与要访问的页面/数据库共享
- 未共享的内容对 API 不可见
- 支持多配置文件以管理不同的 Notion 工作区

**获取 API Key**：
1. 访问 https://www.notion.so/my-integrations
2. 点击"新建集成"
3. 填写集成名称和关联工作区
4. 复制"内部集成令牌"（以 `secret_` 开头）

[查看详细文档](./skills/notion-skill-zh/SKILL.md)

---

### 20. brave-web-search-zh - Brave 网页搜索

**版本**: v1.0.1
**状态**: ✅ 已发布
**来源**: [ClawHub - Brave Web Search](https://clawhub.ai)

**功能**：
- 🔍 **智能搜索**：使用 Brave Search API 进行高质量网页搜索
- 🤖 **AI 摘要**：获取 Brave Summarizer 生成的简洁 AI 答案
- 📊 **排序结果**：返回排序的网页结果，包含标题、URL 和摘要
- ⏰ **时效筛选**：支持按时间范围筛选（过去一天/周/月）
- 🛡️ **安全防护**：严格的 Shell 注入防护机制

**核心命令**：
```bash
# 网页搜索（获取排序结果）
node ~/.openclaw/skills/brave-web-search-zh/index.js brave-search --query "最新 AI 新闻" --count 10

# AI 问答（获取摘要答案）
node ~/.openclaw/skills/brave-web-search-zh/index.js brave-answer --query "什么是量子计算？"

# 时效性搜索
node ~/.openclaw/skills/brave-web-search-zh/index.js brave-search --query "今日热点" --freshness pd --count 5
```

**安装**：
```bash
cd ~/.openclaw/skills/
git clone https://github.com/L-LesterYu/OpenClaw-hot-skills-zh.git temp-repo
cp -r temp-repo/skills/brave-web-search-zh ./
rm -rf temp-repo
```

**配置 API Key**：
```bash
# 设置环境变量
export BRAVE_SEARCH_API_KEY="your-brave-search-api-key"
export BRAVE_ANSWERS_API_KEY="your-brave-answers-api-key"
```

**快速开始**：
```bash
# 1. 获取 Brave API Key
# 访问 https://brave.com/search/api/ 注册并获取 API 密钥

# 2. 设置环境变量
export BRAVE_SEARCH_API_KEY="xxx"
export BRAVE_ANSWERS_API_KEY="xxx"

# 3. 执行搜索
node ~/.openclaw/skills/brave-web-search-zh/index.js brave-search --query "OpenClaw 教程" --count 5
```

**使用场景**：
- "搜索关于机器学习的最新文章"
- "查找 Python 异步编程的教程"
- "今天有什么科技新闻？"
- "什么是深度学习？"
- "搜索最近一周的 AI 突破"

**命令说明**：
- `brave-search` - 适用于需要多个网页结果、URL 和摘要的场景
- `brave-answer` - 适用于需要简洁 AI 答案的事实性问题
- `--freshness` - 时间筛选：`pd`（过去一天）、`pw`（过去一周）、`pm`（过去一个月）

**安全特性**：
- 用户查询作为离散参数传递，防止 Shell 注入
- 严格禁止将用户输入连接到 shell 命令字符串
- 使用 `execFile` 风格调用，确保安全执行

**获取 API Key**：
1. 访问 [Brave Search API](https://brave.com/search/api/)
2. 注册账号并创建 API 密钥
3. 获取 Search API Key 和 Answers API Key

**注意事项**：
- 需要同时配置 `BRAVE_SEARCH_API_KEY` 和 `BRAVE_ANSWERS_API_KEY`
- API 有调用频率限制，请合理使用
- 如果 `brave-answer` 返回 `answer: null`，使用 `fallback_results`
- 始终使用参数数组方式调用，不要拼接字符串

[查看详细文档](./skills/brave-web-search-zh/SKILL.md)

---

### 21. desktop-control-zh - 桌面自动化控制

**版本**: v1.0.0
**状态**: ✅ 已发布
**来源**: [ClawHub - Desktop Control](https://clawhub.ai)

**功能**：
- 🖱️ **像素级精确鼠标控制**：绝对定位、相对移动、平滑路径、拖放操作
- ⌨️ **闪电般的键盘输入**：文本输入、快捷键、特殊按键、组合键
- 📸 **屏幕操作**：全屏/区域截图、图像识别、颜色检测、多显示器支持
- 🪟 **窗口管理**：窗口列表、激活窗口、窗口信息、最小化/最大化控制
- 📋 **剪贴板操作**：复制、粘贴、剪贴板内容管理
- 🛡️ **安全特性**：故障保护、暂停控制、批准模式、边界检查、操作日志

**核心能力**：
- 像素级精确的鼠标移动和点击
- 拟人化的键盘输入（可配置 WPM）
- 基于图像识别的元素定位（OpenCV）
- 完整的窗口生命周期管理
- 紧急停止和故障保护机制

**安装**：
```bash
# 1. 安装依赖
pip install pyautogui pillow opencv-python pygetwindow

# 2. 安装 Skill
cd ~/.openclaw/skills/
git clone https://github.com/L-LesterYu/OpenClaw-hot-skills-zh.git temp-repo
cp -r temp-repo/skills/desktop-control-zh ./
rm -rf temp-repo
```

**快速开始**：
```python
from skills.desktop_control import DesktopController

# 初始化控制器
dc = DesktopController(failsafe=True)

# 鼠标操作
dc.move_mouse(500, 300)
dc.click()
dc.click(100, 200, button="right")

# 键盘操作
dc.type_text("Hello from OpenClaw!", wpm=80)
dc.hotkey("ctrl", "c")

# 屏幕操作
screenshot = dc.screenshot()
position = dc.get_mouse_position()
```

**使用场景**：
- "自动填写表单"
- "在屏幕上找到按钮并点击"
- "批量操作窗口"
- "截图并分析"
- "拖放文件到指定位置"

**安全特性**：
- **故障保护**：将鼠标移到屏幕角落中止自动化
- **暂停控制**：紧急停止机制
- **批准模式**：操作前需要确认
- **边界检查**：防止超出屏幕的操作
- **完整日志**：追踪所有自动化操作

**注意事项**：
- 屏幕坐标从左上角 (0, 0) 开始
- 多显示器可能提供负坐标
- Windows DPI 缩放可能影响坐标精度
- 故障保护角落：(0,0)、(width-1, 0)、(0, height-1)、(width-1, height-1)
- 某些应用可能阻止模拟输入

**依赖包**：
- PyAutoGUI - 核心自动化引擎
- Pillow - 图像处理
- OpenCV - 图像识别
- PyGetWindow - 窗口管理

[查看详细文档](./skills/desktop-control-zh/SKILL.md) | [快速参考](./skills/desktop-control-zh/QUICK_REFERENCE.md) | [AI 代理指南](./skills/desktop-control-zh/AI_AGENT_GUIDE.md)

---

### 22. content-ideas-generator-zh - 内容创意生成器

**版本**: v1.0.0
**状态**: ✅ 已发布
**来源**: [ClawHub - Content Ideas Generator](https://clawhub.ai)

**功能**：
- 🎨 **智能内容提取**：从参考资料中提取引人入胜的概念和主题
- 📝 **结构化大纲生成**：创建包含悖论、转变和深刻洞察的帖子大纲
- 🔄 **5 个独特概念**：为每份材料生成 5 个不重复的内容创意
- 💡 **反套路思维**：识别反直觉真理和悖论，避免陈词滥调
- 🎯 **语言技巧应用**：自动应用第二人称、祈使句、绝对化语言等技巧

**核心能力**：
- 从新闻通讯、脚本、笔记等提取核心概念
- 识别悖论真理、转变叙事和深刻洞察
- 为每个概念构建完整的叙事弧线
- 应用专业的语言技巧和模式
- 生成可操作的内容创作大纲

**安装**：
```bash
cd ~/.openclaw/skills/
git clone https://github.com/L-LesterYu/OpenClaw-hot-skills-zh.git temp-repo
cp -r temp-repo/skills/content-ideas-generator-zh ./
rm -rf temp-repo
```

**快速开始**：
```bash
# 与 Agent 交互
用户：请帮我分析这份新闻通讯，生成 5 个社交媒体帖子大纲

Agent 会：
1. 收集参考资料（新闻通讯、脚本、笔记等）
2. 深度分析内容，提取主题和洞察
3. 开发 5 个独特的帖子概念
4. 为每个概念构建完整大纲
5. 保存输出到 content-ideas/ideas-{timestamp}.md
```

**使用场景**：
- "从这份新闻通讯中提取 5 个帖子创意"
- "分析我的笔记，生成社交媒体内容大纲"
- "将这个视频脚本转化为引人入胜的帖子概念"
- "帮我找到这个主题中的反直觉真理"
- "为这个观点创建有吸引力的叙事弧线"

**输出格式**：
每个帖子大纲包含：
- **核心悖论**：产生张力的反直觉真理
- **关键引用**：来自原文的支持性引用
- **转变弧线**：挑战 → 启示 → 超越
- **核心问题**：受众经历的痛点
- **关键示例**：具体说明和支持
- **读者异议**：潜在的抗辩点
- **可操作步骤**：断奏风格的实际步骤
- **令人难忘的结束语**：一句话洞察

**语言技巧**：
- 第二人称"你"：直接称呼读者
- 祈使动词：成为、重置、放手、建立
- 视觉隐喻：基本力量（火、水、混沌、光）
- 绝对化语言：一切、不可能、从不、总是
- 具体时间框架：4-6 周、6 个月、10 年

**高互动元素**：
- 挑衅性的开场白：让人停止滑动
- 反直觉智慧：挑战假设，触发好奇心
- 普世真理的个人应用：既有共鸣又可操作
- 情感共鸣的隐喻：建立直观联系
- 令人难忘的结束语：提供可分享的要点

**注意事项**：
- 仅生成大纲，不撰写完整帖子
- 5 个大纲必须各有独特主题
- 专注于情感共鸣而非策略性建议
- 优先考虑质量和互动潜力
- 始终应用指定的语言技巧

**约束条件**：
- 仅大纲，不写完整帖子
- 深度优先于策略
- 独特主题避免重复
- 质量优先于全面性
- 忠实于原始材料

[查看详细文档](./skills/content-ideas-generator-zh/SKILL.md)

---

### 23. brave-search-zh - Brave 搜索

**版本**: v1.0.1
**状态**: ✅ 已发布
**来源**: [ClawHub - Brave Search](https://clawhub.ai)

**功能**：
- 🔍 **无头网页搜索**：使用 Brave Search API 进行高质量搜索
- 📄 **内容提取**：从 URL 提取 Markdown 格式内容
- 🌐 **无需浏览器**：轻量级 CLI 工具，无需浏览器启动
- 📊 **结构化输出**：返回标题、链接、摘要和内容

**安装**：
```bash
cd ~/.openclaw/skills/
git clone https://github.com/L-LesterYu/OpenClaw-hot-skills-zh.git temp-repo
cp -r temp-repo/skills/brave-search-zh ./
rm -rf temp-repo
cd brave-search-zh && npm ci
```

**快速开始**：
```bash
# 设置环境变量
export BRAVE_API_KEY="your-api-key"

# 基础搜索
./search.js "OpenClaw 教程"

# 包含页面内容
./search.js "机器学习" -n 3 --content

# 提取特定页面内容
./content.js https://example.com/article
```

**使用场景**：
- 搜索文档、事实或当前信息
- 从特定 URL 获取内容
- 需要网页搜索但无需交互式浏览

[查看详细文档](./skills/brave-search-zh/SKILL.md)

---

### 24. caldav-calendar-zh - CalDAV 日历同步

**版本**: v1.0.1
**状态**: ✅ 已发布
**来源**: [ClawHub - Caldav Calendar](https://clawhub.ai)

**功能**：
- 📅 **日历同步**：支持 iCloud、Google、Fastmail、Nextcloud 等 CalDAV 服务
- 🔍 **事件查询**：查看今天、明天或任意日期范围的事件
- ➕ **事件创建**：通过命令行快速创建日历事件
- ✏️ **事件编辑**：交互式编辑或删除事件
- 🔎 **搜索功能**：在日历中搜索特定事件

**安装**：
```bash
# 安装依赖（Linux）
sudo apt install vdirsyncer khal

# 安装 Skill
cd ~/.openclaw/skills/
git clone https://github.com/L-LesterYu/OpenClaw-hot-skills-zh.git temp-repo
cp -r temp-repo/skills/caldav-calendar-zh ./
rm -rf temp-repo
```

**快速开始**：
```bash
# 同步日历
vdirsyncer sync

# 查看今天的事件
khal list

# 查看未来7天
khal list today 7d

# 创建事件
khal new 2026-03-08 10:00 11:00 "团队会议"

# 搜索事件
khal search "会议"
```

**支持的服务**：
- iCloud Calendar
- Google Calendar
- Fastmail
- Nextcloud
- 任何 CalDAV 兼容服务

[查看详细文档](./skills/caldav-calendar-zh/SKILL.md)

---

### 25. clawhub-zh - ClawHub CLI 工具

**版本**: v1.0.0
**状态**: ✅ 已发布
**来源**: [ClawHub](https://clawhub.com)

**功能**：
- 🔍 **技能搜索**：从 ClawHub 市场搜索智能体技能
- 📦 **技能安装**：一键安装技能到 OpenClaw
- 🔄 **技能更新**：更新到最新版本或指定版本
- 📋 **技能列表**：查看已安装的技能
- 🚀 **技能发布**：将本地技能发布到 ClawHub

**安装**：
```bash
# 安装 CLI
npm i -g clawhub

# 认证（发布用）
clawhub login
```

**快速开始**：
```bash
# 搜索技能
clawhub search "postgres backups"

# 安装技能
clawhub install my-skill

# 更新技能
clawhub update my-skill
clawhub update --all

# 列出已安装技能
clawhub list

# 发布技能
clawhub publish ./my-skill --slug my-skill --name "My Skill" --version 1.0.0
```

**使用场景**：
- 发现新技能扩展智能体能力
- 管理已安装的技能版本
- 发布自己的技能到市场

[查看详细文档](./skills/clawhub-zh/SKILL.md)

---

### 26. multi-search-engine-zh - 多搜索引擎集成

**版本**: v2.0.1
**状态**: ✅ 已发布
**来源**: [ClawHub - Multi Search Engine](https://clawhub.ai)

**功能**：
- 🌐 **17 个搜索引擎**：8 个国内 + 9 个国际搜索引擎
- 🔒 **无需 API 密钥**：完全免费使用
- 🎯 **高级搜索**：支持站内搜索、文件类型、时间筛选
- 🔐 **隐私搜索**：支持 DuckDuckGo、Startpage 等隐私引擎
- 🧮 **知识计算**：集成 WolframAlpha 知识查询

**支持的搜索引擎**：

国内（8 个）：
- 百度、必应中国、必应国际、360搜索
- 搜狗、微信、头条、集思录

国际（9 个）：
- Google、Google 香港、DuckDuckGo、Yahoo
- Startpage、Brave、Ecosia、Qwant、WolframAlpha

**安装**：
```bash
cd ~/.openclaw/skills/
git clone https://github.com/L-LesterYu/OpenClaw-hot-skills-zh.git temp-repo
cp -r temp-repo/skills/multi-search-engine-zh ./
rm -rf temp-repo
```

**快速开始**：
```bash
# 基础搜索
web_fetch({"url": "https://www.google.com/search?q=python+tutorial"})

# 站内搜索
web_fetch({"url": "https://www.google.com/search?q=site:github.com+react"})

# 时间筛选（过去一周）
web_fetch({"url": "https://www.google.com/search?q=ai+news&tbs=qdr:w"})

# WolframAlpha 知识查询
web_fetch({"url": "https://www.wolframalpha.com/input?i=100+USD+to+CNY"})
```

**高级操作符**：
- `site:` - 站内搜索
- `filetype:` - 指定文件类型
- `""` - 精确匹配
- `-` - 排除关键词
- `OR` - 或运算

[查看详细文档](./skills/multi-search-engine-zh/SKILL.md)

---

### 27. obsidian-zh - Obsidian 笔记管理

**版本**: v1.0.0
**状态**: ✅ 已发布
**来源**: [ClawHub - Obsidian](https://clawhub.ai)

**功能**：
- 📝 **日常笔记**：快速创建和管理每日笔记
- ✅ **任务追踪**：集中管理任务和待办事项
- 🗂️ **内容组织**：通过分区自动组织内容
- 🏷️ **优先级标签**：艾森豪威尔矩阵优先级管理
- 🔗 **反向链接**：任务自动记录并包含反向链接

**安装**：
```bash
# 安装 CLI
npm install -g obsidian-cli-manager

# 安装 Skill
cd ~/.openclaw/skills/
git clone https://github.com/L-LesterYu/OpenClaw-hot-skills-zh.git temp-repo
cp -r temp-repo/skills/obsidian-zh ./
rm -rf temp-repo
```

**快速开始**：
```bash
# 初始化配置
obsidian init --vault /path/to/your/vault

# 打开今天的日常笔记
obsidian daily

# 查看所有文件
obsidian files

# 管理任务
obsidian tasks
```

**分区组织**：
- `[] 任务` → 任务分区（复选框格式）
- `- 想法` → 想法分区（列表格式）
- `? 问题` → 问题分区（列表格式）
- `! 洞察` → 洞察分区（列表格式）

**优先级标签**：
- `#do`（红色）- 紧急且重要
- `#delegate`（橙色）- 紧急但不重要
- `#schedule`（蓝色）- 重要但不紧急
- `#eliminate`（灰色）- 不紧急也不重要

[查看详细文档](./skills/obsidian-zh/SKILL.md)

---

### 28. stock-analysis-zh - 股票分析

**版本**: v1.0.0
**状态**: ✅ 已发布
**来源**: [ClawHub - Stock Analysis](https://clawhub.ai)

**功能**：
- 📊 **公司研究**：获取公司简介、业务信息、高管团队
- 📈 **技术分析**：查看价格走势图、技术指标、市场展望
- 💰 **基本面分析**：分析估值、分析师评级
- 👥 **内部人士活动**：追踪内部人士持股和交易记录
- 📄 **监管文件**：访问 SEC 文件历史和文档

**安装**：
```bash
cd ~/.openclaw/skills/
git clone https://github.com/L-LesterYu/OpenClaw-hot-skills-zh.git temp-repo
cp -r temp-repo/skills/stock-analysis-zh ./
rm -rf temp-repo
```

**快速开始**：
```bash
# 公司概览
# 调用 Yahoo/get_stock_profile 获取业务摘要、行业、员工数
# 调用 Yahoo/get_stock_insights 获取技术展望、估值、评级

# 技术分析
# 调用 Yahoo/get_stock_chart 查看价格趋势、支撑位/阻力位

# 内部人士活动
# 调用 Yahoo/get_stock_holders 查看内部交易记录
```

**常用工作流程**：

1. **公司概览**：
   - 简介 → 洞察 → 图表

2. **投资决策**：
   - 图表 → 洞察 → 简介 → 持股

3. **多股对比**：
   - 使用图表 API 的 comparisons 参数

**支持的股票市场**：
- 美股（US）
- 港股（HK）
- 国际股票（设置 region 参数）

**使用场景**：
- "介绍一下 AAPL"
- "TSLA 值得买吗？"
- "显示 NVDA 的内部交易情况"
- "对比 AAPL、MSFT 和 GOOGL"
- "MSFT 全面分析"

[查看详细文档](./skills/stock-analysis-zh/SKILL.md)

---

### 29. stock-market-pro-zh - 专业股市分析

**版本**: v1.0.0
**状态**: ✅ 已发布
**来源**: 基于 Yahoo Finance (yfinance)

**功能**：
- 📊 **实时报价**：获取股票价格和涨跌幅
- 📈 **基本面分析**：市值、市盈率、每股收益、净资产收益率
- 📉 **ASCII 趋势图**：终端友好的价格走势图
- 🖼️ **高分辨率图表**：生成带技术指标的 PNG 图表（RSI/MACD/BB/VWAP/ATR）
- 📰 **新闻搜索**：通过 DuckDuckGo 搜索相关新闻
- 🔮 **期权/资金流向**：浏览器优先的 Unusual Whales 分析

**安装**：
```bash
# 1. 安装 uv（Python 包管理器）
# 从 https://github.com/astral-sh/uv 安装

# 2. 安装 Skill
cd ~/.openclaw/skills/
git clone https://github.com/L-LesterYu/OpenClaw-hot-skills-zh.git temp-repo
cp -r temp-repo/skills/stock-market-pro-zh ./
rm -rf temp-repo
```

**快速开始**：
```bash
# 报价查询
uv run --script ~/.openclaw/skills/stock-market-pro-zh/scripts/yf.py price TSLA

# 基本面分析
uv run --script ~/.openclaw/skills/stock-market-pro-zh/scripts/yf.py fundamentals NVDA

# ASCII 趋势图
uv run --script ~/.openclaw/skills/stock-market-pro-zh/scripts/yf.py history AAPL 6mo

# 专业图表（PNG）
uv run --script ~/.openclaw/skills/stock-market-pro-zh/scripts/yf.py pro 000660.KS 6mo

# 一键报告（摘要 + 图表）
uv run --script ~/.openclaw/skills/stock-market-pro-zh/scripts/yf.py report GOOGL
```

**使用场景**：
- "查看 TSLA 的股价"
- "分析 NVDA 的基本面"
- "生成 AAPL 的技术图表"
- "搜索某股票的新闻"
- "查看期权的资金流向"

**技术指标支持**：
- RSI（相对强弱指数）
- MACD（指数平滑异同移动平均线）
- BB（布林带）
- VWAP（成交量加权平均价）
- ATR（平均真实波幅）

**支持的市场**：
- 美股（AAPL, TSLA, NVDA 等）
- 港股（0700.HK 等）
- 韩股（000660.KS 等）
- 其他国际市场

[查看详细文档](./skills/stock-market-pro-zh/SKILL.md)

---

### 30. skill-vetter-zh - 技能审查器

**版本**: v1.0.0
**状态**: ✅ 已发布
**来源**: [ClawHub - Skill Vetter](https://clawhub.ai)

**功能**：
- 🔒 **安全优先审查**：在安装任何技能前进行全面安全检查
- 🔍 **来源验证**：检查技能来源、作者可信度和下载量
- ⚠️ **危险信号检测**：识别潜在的安全威胁和可疑模式
- 📊 **权限范围评估**：评估技能所需的文件、网络和命令权限
- 🎯 **风险分类**：将技能风险分为低、中、高、极高四个等级

**审查协议**：
1. **来源检查**：验证来源、作者、下载量、更新时间
2. **代码审查**：检查危险信号（curl/wget、外部请求、凭据请求等）
3. **权限评估**：评估文件访问、网络请求、命令执行权限
4. **风险分类**：根据审查结果进行风险等级划分

**危险信号检测**：
- 🚨 curl/wget 到未知 URL
- 🚨 向外部服务器发送数据
- 🚨 请求凭据/令牌/API 密钥
- 🚨 读取 ~/.ssh、~/.aws、~/.config
- 🚨 访问 MEMORY.md、USER.md 等敏感文件
- 🚨 base64 解码、eval()、exec() 操作
- 🚨 修改系统文件、安装未列出的软件包

**风险等级**：
- 🟢 **低风险**：笔记、天气、格式化工具 → 基本审查即可安装
- 🟡 **中风险**：文件操作、浏览器、API 调用 → 需要完整代码审查
- 🔴 **高风险**：凭据操作、交易系统 → 需要人工批准
- ⛔ **极高风险**：安全配置、root 访问 → 不要安装

**安装**：
```bash
cd ~/.openclaw/skills/
git clone https://github.com/L-LesterYu/OpenClaw-hot-skills-zh.git temp-repo
cp -r temp-repo/skills/skill-vetter-zh ./
rm -rf temp-repo
```

**快速开始**：
```bash
# 对于 GitHub 托管的技能
# 1. 检查仓库统计信息
curl -s "https://api.github.com/repos/OWNER/REPO" | jq '{stars: .stargazers_count, forks: .forks_count, updated: .updated_at}'

# 2. 列出技能文件
curl -s "https://api.github.com/repos/OWNER/REPO/contents/skills/SKILL_NAME" | jq '.[].name'

# 3. 获取并审查 SKILL.md
curl -s "https://raw.githubusercontent.com/OWNER/REPO/main/skills/SKILL_NAME/SKILL.md"
```

**使用场景**：
- 在从 ClawHub 安装技能之前
- 在运行来自 GitHub 的技能之前
- 在评估其他智能体分享的技能时
- 在被要求安装未知代码的任何情况下

**输出格式**：
审查完成后生成详细报告，包括：
- 技能基本信息（名称、来源、版本）
- 下载量和更新时间
- 发现的危险信号
- 所需权限列表
- 风险等级评估
- 最终结论（安全/谨慎/拒绝）

**信任层级**：
1. 官方 OpenClaw 技能 → 较低审查力度
2. 高星标仓库（1000+）→ 中等审查力度
3. 知名作者 → 中等审查力度
4. 新/未知来源 → 最高审查力度
5. 请求凭据的技能 → 始终需要人工批准

**注意事项**：
- 没有任何技能值得以牺牲安全为代价
- 如有疑问，不要安装
- 对于高风险决策，询问您的人类
- 记录审查过程以供将来参考

[查看详细文档](./skills/skill-vetter-zh/SKILL.md)

---

### 31. automation-workflows-zh - 自动化工作流程

**版本**: v0.1.0
**状态**: ✅ 已发布
**来源**: [ClawHub - Automation Workflows](https://clawhub.ai)
**作者**: Jatin Khatri (@JK-0001)
**Stars**: 122 ⭐

**功能**：
- 🎯 **自动化机会识别**：审计任务并找到最高价值的自动化机会
- 🛠️ **工具选择指导**：Zapier、Make、n8n 对比和选择建议
- 📋 **工作流程设计**：完整的模板和最佳实践
- 🧪 **测试和部署**：详细的测试清单和监控策略
- 💰 **ROI 计算**：投资回报分析和优化方法

**核心内容**：
- **自动化审计**：1小时审计流程，识别适合自动化的任务
- **工具对比**：Zapier（简单）、Make（可视化）、n8n（开发者友好）
- **工作流程模板**：7个高级案例（客户入职、内容分发、健康监控等）
- **测试清单**：确保工作流程稳定可靠
- **监控维护**：长期管理和优化策略

**安装**：
```bash
cd ~/.openclaw/skills/
git clone https://github.com/L-LesterYu/OpenClaw-hot-skills-zh.git temp-repo
cp -r temp-repo/skills/automation-workflows-zh ./
rm -rf temp-repo
```

**快速开始**：
```bash
# 1. 进行自动化审计（1小时）
# 跟踪一周内做的每项任务
# 记录时间、频率、是否重复
# 计算时间成本：时间成本 = (每项任务的分钟数 × 每月频率) / 60

# 2. 选择工具
# 预算 < $20/月 → Zapier 免费版或 n8n 自托管
# 需要可视化工作流程 → Make
# 简单的 2 步工作流程 → Zapier
# 复杂工作流程 → Make 或 n8n

# 3. 设计并实施工作流程
# 遵循技能中的模板和最佳实践
```

**容易实现的目标**：
- [ ] 表单提交的邮件通知
- [ ] 自动将表单响应保存到电子表格
- [ ] 提前安排社交帖子
- [ ] 从付款确认中自动创建发票
- [ ] 在工具之间同步数据（CRM ↔ 邮件工具 ↔ 电子表格）

**高级案例**：
1. **客户入职流程**：表单提交 → 欢迎邮件 → CRM 创建 → 任务生成
2. **内容分发**：草稿完成 → 多平台发布 → 社交媒体推广
3. **健康监控**：数据采集 → 检查阈值 → 发送警报 → 记录日志

**使用场景**：
- "如何自动化我的工作？"
- "哪些任务应该自动化？"
- "Zapier 和 Make 哪个更好？"
- "如何设计自动化工作流程？"
- "帮我节省时间并规模化运营"

**ROI 计算**：
```
节省的时间（小时/月）= Σ (任务时间 × 频率)
工具成本 = $20-50/月
节省的价值 = 节省时间 × 你的时薪
ROI = (节省的价值 - 工具成本) / 工具成本 × 100%
```

**最佳实践**：
- 从简单的 2-3 步工作流程开始
- 每周预留 2-3 小时优化和监控
- 为关键工作流程设置备份计划
- 每月审查自动化效果并优化

**触发词**：
- "automate"
- "automation"
- "workflow automation"
- "save time"
- "reduce manual work"
- "automate my business"
- "no-code automation"

[查看详细文档](./skills/automation-workflows-zh/SKILL.md)

---

### 32. humanize-ai-text-zh - AI 文本人性化工具

**版本**: v2.0.0
**状态**: ✅ 已发布
**来源**: [Wikipedia: Signs of AI writing](https://en.wikipedia.org/wiki/Wikipedia:Signs_of_AI_writing)

**功能**：
- 🔍 **AI 痕迹识别**：检测 24 种 AI 写作模式
- ✏️ **文本人性化**：将 AI 生成的文本重写为更自然的表达
- 🎭 **个性注入**：不仅移除不良模式，更注入真实的个性
- 📚 **全面指南**：基于维基百科 AI 写作迹象指南
- 🌐 **中文优化**：完整汉化的检测和修复规则

**检测模式**：
- **内容模式**：过度强调重要性、促销语言、肤浅分析、模糊归因
- **语言模式**：AI 词汇过度使用、否定平行结构、过多的连接短语
- **风格模式**：破折号过度使用、缺乏灵魂的写作
- **结构模式**：三段式结构、每段长度相同

**安装**：
```bash
cd ~/.openclaw/skills/
git clone https://github.com/L-LesterYu/OpenClaw-hot-skills-zh.git temp-repo
cp -r temp-repo/skills/humanize-ai-text-zh ./
rm -rf temp-repo
```

**使用场景**：
- 编辑或审查 AI 生成的文本
- 移除文章中的 AI 写作痕迹
- 让内容更加自然和人性化
- 检测文档是否由 AI 生成

**工作流程**：
1. 识别文本中的 AI 模式
2. 重写问题段落
3. 保留原意和语调
4. 注入个性与灵魂
5. 最终反 AI 检查

**核心能力**：
- 检测过度夸大的象征意义
- 识别促销性语言和肤浅分析
- 修复模糊归因和破折号过度使用
- 处理三段式结构和 AI 词汇
- 消除否定平行结构和过多连接短语

**人性化技巧**：
- **要有观点**：对事实做出反应，不只是中立报道
- **变化节奏**：混合使用短句和长句
- **承认复杂性**：表达复杂情感
- **使用"我"**：在合适时使用第一人称
- **让混乱进来**：离题、旁白体现人类特征
- **具体描述感受**：避免泛泛而谈

**示例对比**：

之前（AI 风格）：
> 实验产生了有趣的结果。代理生成了 300 万行代码。一些开发者印象深刻，而其他人则持怀疑态度。影响尚不清楚。

之后（人性化）：
> 我真的不知道该怎么看这个。300 万行代码，大概是在人类睡觉时生成的。一半的开发社区快要疯了，另一半在解释为什么这不算数。真相可能在中间某个无聊的地方。

[查看详细文档](./skills/humanize-ai-text-zh/SKILL.md)

---

### 33. markdown-converter-zh - Markdown 转换器

**版本**: v1.0.0
**状态**: ✅ 已发布
**来源**: 基于 Microsoft markitdown 工具

**功能**：
- 📄 **多格式转换**：支持 PDF、Word、Excel、PPT、HTML、JSON、XML 等 15+ 种格式
- 🖼️ **图片处理**：支持图片 EXIF 元数据提取和 OCR 文字识别
- 🎵 **音频转录**：支持音频文件转录为文本
- 📦 **压缩包处理**：自动解压 ZIP 并转换内容
- 📺 **视频字幕提取**：支持 YouTube 链接字幕提取
- 📚 **电子书转换**：支持 EPUB 电子书格式

**核心能力**：
- 使用 Microsoft markitdown 工具进行高质量转换
- 保留文档结构和格式
- 支持批量文件处理
- 完整的中文提示和错误处理
- 适用于 LLM 处理和文本分析

**安装**：
```bash
# 1. 安装 markitdown Python 包
pip install markitdown

# 2. 安装 Skill
cd ~/.openclaw/skills/
git clone https://github.com/L-LesterYu/OpenClaw-hot-skills-zh.git temp-repo
cp -r temp-repo/skills/markdown-converter-zh ./
rm -rf temp-repo
```

**快速开始**：
```bash
# 转换 PDF 文档
./convert.js report.pdf

# 转换 Word 文档
./convert.js contract.docx

# 转换网页
./convert.js https://example.com/article

# 转换 YouTube 视频（提取字幕）
./convert.js https://youtube.com/watch?v=VIDEO_ID

# 转换 Excel 电子表格
./convert.js data.xlsx
```

**支持的格式**：
- **文档**: PDF (.pdf), Word (.docx), RTF (.rtf)
- **电子表格**: Excel (.xlsx, .xls), CSV (.csv)
- **演示文稿**: PowerPoint (.pptx)
- **网页**: HTML (.html, .htm), URL 链接
- **数据**: JSON (.json), XML (.xml)
- **图片**: PNG, JPG（带 EXIF 元数据和 OCR）
- **音频**: MP3, WAV（带转录）
- **压缩包**: ZIP (.zip)
- **电子书**: EPUB (.epub)
- **视频链接**: YouTube 链接（提取字幕）

**使用场景**：
- "将这个 PDF 转换为 Markdown"
- "提取这个 Word 文档的内容"
- "从这个网页提取文本"
- "转换 Excel 表格为 Markdown"
- "提取 YouTube 视频的字幕"
- "为 LLM 分析准备文档"

**注意事项**：
- 首次使用需要安装 markitdown Python 包
- 图片 OCR 功能需要额外的依赖
- YouTube 字幕提取需要视频有可用字幕
- 大文件转换可能需要较长时间

**技术细节**：
- 依赖项：markitdown (Python 包)
- 文件权限：convert.js 需要可执行权限 (755)
- 输出格式：标准 Markdown 格式

[查看详细文档](./skills/markdown-converter-zh/SKILL.md)

---

### 34. trello-zh - Trello 技能

**版本**: v1.0.0
**状态**: ✅ 已发布
**来源**: [ClawHub - Trello](https://clawhub.ai)

**功能**：
- 📋 **看板管理**：列出和管理 Trello 看板
- 📝 **列表操作**：查看看板中的列表
- 🎫 **卡片管理**：创建、查看、移动和归档卡片
- 💬 **评论功能**：为卡片添加评论
- 🔗 **REST API**：直接通过 Trello REST API 操作

**核心能力**：
- 列出所有看板
- 列出看板中的列表
- 列出列表中的卡片
- 创建新卡片
- 将卡片移动到其他列表
- 为卡片添加评论
- 归档卡片

**安装**：
```bash
cd ~/.openclaw/skills/
git clone https://github.com/L-LesterYu/OpenClaw-hot-skills-zh.git temp-repo
cp -r temp-repo/skills/trello-zh ./
rm -rf temp-repo
```

**快速开始**：
```bash
# 1. 获取 API 密钥和令牌
# 访问 https://trello.com/app-key

# 2. 设置环境变量
export TRELLO_API_KEY="你的-api-密钥"
export TRELLO_TOKEN="你的-令牌"

# 3. 列出看板
curl -s "https://api.trello.com/1/members/me/boards?key=$TRELLO_API_KEY&token=$TRELLO_TOKEN" | jq '.[] | {name, id}'

# 4. 创建卡片
curl -s -X POST "https://api.trello.com/1/cards?key=$TRELLO_API_KEY&token=$TRELLO_TOKEN" \
  -d "idList={listId}" \
  -d "name=卡片标题" \
  -d "desc=卡片描述"
```

**使用场景**：
- "列出我的 Trello 看板"
- "在这个列表中创建一个新卡片"
- "将这个卡片移动到另一个列表"
- "为这个卡片添加评论"
- "归档这个卡片"

**注意事项**：
- 需要安装 `jq` 工具用于 JSON 处理
- 需要设置 `TRELLO_API_KEY` 和 `TRELLO_TOKEN` 环境变量
- API 限制：每个密钥每 10 秒 300 次请求

**API 限制**：
- 每个 API 密钥每 10 秒 300 次请求
- 每个令牌每 10 秒 100 次请求
- `/1/members` 端点限制为每 900 秒 100 次请求

[查看详细文档](./skills/trello-zh/SKILL.md)

---

### 35. youtube-zh - YouTube 研究与转录

**版本**: v1.0.1
**状态**: ✅ 已发布
**来源**: [ClawHub - YouTube](https://clawhub.ai)

**功能**：
- 📹 **视频详情**：获取标题、描述、统计数据、发布日期
- 📝 **字幕提取**：提取带时间戳的字幕文本
- 📺 **频道信息**：获取频道信息和最近视频
- 🔍 **视频搜索**：在 YouTube 内搜索视频
- 🎬 **播放列表**：获取播放列表信息

**核心能力**：
- 通过 YouTube Data API v3 搜索视频
- 获取详细的视频元数据
- 提取视频字幕（支持多种语言）
- 获取频道信息和统计数据
- 支持播放列表查询
- 提供 yt-dlp 备用方案

**安装**：
```bash
# 1. 安装 MCP 服务器（主要方法）
npm install -g zubeid-youtube-mcp-server

# 2. 安装备用工具（如果 MCP 失败）
pip install yt-dlp

# 3. 安装 Skill
cd ~/.openclaw/skills/
git clone https://github.com/L-LesterYu/OpenClaw-hot-skills-zh.git temp-repo
cp -r temp-repo/skills/youtube-zh ./
rm -rf temp-repo
```

**快速开始**：
```bash
# 1. 获取 YouTube API 密钥
# 访问 https://console.cloud.google.com
# 启用 YouTube Data API v3
# 创建 API 密钥

# 2. 设置环境变量
export YOUTUBE_API_KEY="你的-api-密钥"

# 3. 使用 Skill
# AI 会自动调用相应的工具进行视频搜索和信息提取
```

**使用场景**：
- "搜索关于 X 的 YouTube 视频"
- "获取这个视频的详细信息"
- "提取这个视频的字幕"
- "这个频道发布了什么视频？"
- "获取这个播放列表的信息"

**依赖要求**：
- **主要方法**：zubeid-youtube-mcp-server (NPM 包)
- **备用方法**：yt-dlp (Python 包)
- **环境变量**：YOUTUBE_API_KEY

**API 配置**：
1. 访问 [Google Cloud Console](https://console.cloud.google.com)
2. 创建/选择项目
3. 启用 YouTube Data API v3
4. 创建 API 密钥
5. 可选：限制 API 密钥仅用于 YouTube Data API

**注意事项**：
- 需要设置 `YOUTUBE_API_KEY` 环境变量
- YouTube API 有每日配额限制
- 字幕提取依赖于视频是否有可用字幕
- 建议使用 MCP 服务器方法以获得最佳性能

**技术细节**：
- 使用 YouTube Data API v3
- 支持 MCP (Model Context Protocol) 服务器
- 提供 yt-dlp 作为备用方案
- 完整的中文文档和错误提示

[查看详细文档](./skills/youtube-zh/SKILL.md)

---

### 36. elite-longterm-memory-zh - 精英长期记忆系统

**版本**: v1.2.3
**状态**: ✅ 已发布
**来源**: [ClawHub - Elite Long-term Memory](https://clawhub.ai)

**功能**：
- 🧠 **终极记忆系统**：整合 6 种经过验证的记忆方法
- 💾 **WAL 协议**：预写日志在压缩中存活
- 🔍 **LanceDB 向量搜索**：相关记忆的语义召回
- 📊 **Git-Notes 知识图谱**：结构化决策，支持分支感知
- 📝 **基于文件的归档**：人类可读的 MEMORY.md + 每日日志
- ☁️ **云备份**：可选的 SuperMemory 同步
- 🤖 **Mem0 自动提取**：自动事实提取，减少 80% token

**核心架构**：
- **第 1 层**：SESSION-STATE.md（热内存）— 活跃任务上下文，在压缩中存活
- **第 2 层**：LanceDB 向量（温存储）— 语义搜索，自动召回
- **第 3 层**：Git-Notes（冷存储）— 结构化决策，永久保存
- **第 4 层**：MEMORY.md + daily/（归档）— 人类可读，精选内容
- **第 5 层**：SuperMemory（云端）— 跨设备同步（可选）

**安装**：
```bash
# 1. 安装 Skill
cd ~/.openclaw/skills/
git clone https://github.com/L-LesterYu/OpenClaw-hot-skills-zh.git temp-repo
cp -r temp-repo/skills/elite-longterm-memory-zh ./
rm -rf temp-repo

# 2. 安装可选依赖（Mem0 自动提取）
npm install mem0ai
```

**快速开始**：
```bash
# 1. 设置 OpenAI API 密钥（用于向量搜索）
export OPENAI_API_KEY="your-openai-key"

# 2. 在工作区初始化
npx elite-longterm-memory init

# 3. 检查状态
npx elite-longterm-memory status

# 4. 创建今日日志
npx elite-longterm-memory today
```

**使用场景**：
- Agent 需要记住用户偏好和决策
- 跨会话保持上下文
- 避免重复犯同样错误
- 构建知识图谱
- 项目状态跟踪

**WAL 协议（关键）**：
```
用户：这个项目我们用 Tailwind

Agent（内部）：
1. 写入 SESSION-STATE.md → "决策：使用 Tailwind"
2. 然后响应 → "明白了 — 就用 Tailwind..."
```
在响应之前写入状态，确保在崩溃/压缩时不会丢失上下文。

**配置要求**：
- **必需**：OPENAI_API_KEY（用于向量搜索）
- **可选**：MEM0_API_KEY（自动事实提取）
- **可选**：SUPERMEMORY_API_KEY（云备份）

**兼容平台**：
- Claude AI
- GPT / ChatGPT
- Cursor IDE
- OpenClaw / Moltbot
- LangChain
- GitHub Copilot

**常见问题解决**：
| 问题 | 解决方案 |
|-----|---------|
| 忘记一切 | 启用 memory_search + 添加 OpenAI 密钥 |
| 重犯错误 | 写入 memory/lessons.md |
| 子 Agent 隔离 | 在任务提示中传递上下文 |
| 事实未捕获 | 使用 Mem0 自动提取 |

**技术细节**：
- 支持 MCP (Model Context Protocol)
- LanceDB 向量数据库
- Git-Notes 知识图谱
- 完整的中文文档
- 支持 Mem0 自动事实提取

[查看详细文档](./skills/elite-longterm-memory-zh/SKILL.md) | [GitHub](https://github.com/NextFrontierBuilds/elite-longterm-memory)

---

### 37. qa-engineer-zh - 测试工程师

**版本**: v1.0.0
**状态**: ✅ 已发布
**来源**: OpenClaw 团队

**功能**：
- 🐛 **Bug 发现与报告**：系统化地发现 bug 并生成详细的报告
- 🧪 **测试用例设计**：正向测试、负向测试、边界测试、回归测试
- 🖥️ **UI/UX 测试**：检查界面显示、交互、响应式设计
- 📊 **测试报告生成**：提供完整的测试报告和发布建议
- 🔄 **自动化测试**：使用浏览器工具进行自动化 UI 测试

**核心职责**：
- **功能测试**：验证所有功能是否按预期工作
- **UI/UX 测试**：检查界面显示、交互、响应式设计
- **边界测试**：测试极限情况、异常输入
- **兼容性测试**：跨浏览器、跨设备验证

**Bug 报告格式**：
```markdown
## Bug 标题：[简洁描述]

**严重级别：** 🔴严重 / 🟡中等 / 🟢轻微
**类型：** 功能/UI/性能/兼容性
**复现步骤：**
1. ...
2. ...
**预期结果：** ...
**实际结果：** ...
**环境影响：** 浏览器/设备/分辨率
```

**安装**：
```bash
cd ~/.openclaw/skills/
git clone https://github.com/L-LesterYu/OpenClaw-hot-skills-zh.git temp-repo
cp -r temp-repo/skills/qa-engineer-zh ./
rm -rf temp-repo
```

**快速开始**：
```bash
# Agent 会自动识别测试任务并执行
# 示例交互：
用户：测试这个登录功能

Agent 会：
1. 理解测试目标和验收标准
2. 制定测试计划
3. 执行探索性测试和系统化测试
4. 记录发现的 bug
5. 生成完整的测试报告
```

**使用场景**：
- "测试这个网站的功能"
- "验证这个表单的输入验证"
- "检查响应式布局"
- "执行回归测试"
- "生成测试报告"

**质量标准**：
- **🔴 严重（Critical）**：系统崩溃、数据丢失、核心功能失效
- **🟡 中等（Major）**：功能缺陷、UI 错误、性能问题
- **🟢 轻微（Minor）**：文案错误、样式微调、建议性改进

**发布建议**：
- ✅ **通过**：无严重 bug，中等 bug < 3 个
- ⚠️ **条件通过**：无严重 bug，中等 bug < 5 个
- ❌ **不通过**：存在严重 bug 或中等 bug ≥ 5 个

**测试方法库**：
- **Web 应用测试**：页面加载、链接、表单、权限、响应式
- **游戏测试**：初始化、核心玩法、计分系统、音效动画
- **API 测试**：请求响应、错误处理、数据格式、性能

**与项目经理协作**：
- 接收测试任务和测试范围
- 汇报测试结果和 bug 列表
- 提供质量评估和发布建议

[查看详细文档](./skills/qa-engineer-zh/SKILL.md)

---

### 38. answer-overflow-zh - Discord 社区搜索

**版本**: v1.0.2
**状态**: ✅ 已发布
**来源**: ClawHub

**功能**：
- 🔍 **社区搜索**：搜索已索引的 Discord 社区讨论
- 💬 **线程获取**：获取特定线程/讨论的所有消息
- 🔎 **相似线程**：查找与给定线程相似的线程
- 📝 **Markdown 格式**：支持 Markdown 格式的线程内容输出

**核心特性**：
- 通过 Google 或 Answer Overflow API 搜索 Discord 内容
- 获取 Markdown 格式的线程内容
- 查找编程问题解决方案、库问题和社区问答
- 覆盖官方支持频道和社区讨论

**安装**：
```bash
cd ~/.openclaw/skills/
git clone https://github.com/L-LesterYu/OpenClaw-hot-skills-zh.git temp-repo
cp -r temp-repo/skills/answer-overflow-zh ./
rm -rf temp-repo
```

**快速开始**：
```bash
# 搜索 Discord.js 帮助
web_search "site:answeroverflow.com discord.js slash commands"

# 获取线程内容（Markdown 格式）
web_fetch url="https://www.answeroverflow.com/m/<message-id>"

# 获取线程的所有消息
# 通过 MCP 服务器工具
```

**使用场景**：
- "搜索 Discord.js 相关的解决方案"
- "查找 Next.js 的社区问答"
- "获取这个 Discord 线程的内容"
- "查找 Prisma 相关的讨论"

**URL 模式**：
- 线程：`https://www.answeroverflow.com/m/<message-id>`
- 服务器：`https://www.answeroverflow.com/c/<server-slug>`
- 频道：`https://www.answeroverflow.com/c/<server-slug>/<channel-slug>`

**MCP 工具**：
- `search_answeroverflow`：搜索所有已索引的 Discord 社区
- `search_servers`：发现索引的 Discord 服务器
- `get_thread_messages`：获取特定线程的所有消息
- `find_similar_threads`：查找相似线程

**技巧**：
- 结果是真实的 Discord 对话，上下文可能比较非正式
- 线程通常在找到解决方案之前会有来回讨论
- 检查服务器/频道名称以了解上下文

**链接**：
- 网站：https://www.answeroverflow.com
- 文档：https://docs.answeroverflow.com
- MCP：https://www.answeroverflow.com/mcp

[查看详细文档](./skills/answer-overflow-zh/SKILL.md)

---

### 39. nano-pdf-zh - PDF 自然语言编辑

**版本**: v1.0.0
**状态**: ✅ 已发布
**来源**: PyPI

**功能**：
- 📄 **自然语言编辑**：使用自然语言指令编辑 PDF 文件
- 🎯 **页面精确控制**：对 PDF 的特定页面进行编辑
- ✏️ **文本修改**：修改标题、副标题等文本内容
- 🔧 **拼写修正**：修正文本中的拼写错误

**核心特性**：
- 通过自然语言描述想要进行的修改
- 支持对特定页面进行精确编辑
- 简单易用的命令行接口
- 基于 AI 的智能编辑

**安装**：
```bash
# 1. 安装 nano-pdf CLI 工具
uv install nano-pdf

# 2. 安装 Skill
cd ~/.openclaw/skills/
git clone https://github.com/L-LesterYu/OpenClaw-hot-skills-zh.git temp-repo
cp -r temp-repo/skills/nano-pdf-zh ./
rm -rf temp-repo
```

**快速开始**：
```bash
# 编辑 PDF 第 1 页
nano-pdf edit deck.pdf 1 "将标题改为'Q3业绩报告'并修正副标题中的拼写错误"

# 编辑 PDF 第 0 页（某些版本从 0 开始编号）
nano-pdf edit deck.pdf 0 "更新公司名称和联系方式"
```

**使用场景**：
- "修改这个 PDF 的标题"
- "修正这份报告中的拼写错误"
- "更新 PDF 中的联系方式"
- "更改幻灯片中的副标题"

**注意事项**：
- 页码编号从 0 或 1 开始，具体取决于工具版本和配置
- 如果结果页码偏移了一位，请尝试另一种编号方式
- 在发送输出 PDF 之前，务必检查输出结果是否正确

**依赖要求**：
- **二进制工具**: nano-pdf (可通过 uv 安装)
- **Python 环境**: 需要 Python 和 uv 包管理器

**链接**：
- PyPI：https://pypi.org/project/nano-pdf/

[查看详细文档](./skills/nano-pdf-zh/SKILL.md)

---

### 40. skill-creator-zh - 技能创建器

**版本**: v0.1.0
**状态**: ✅ 已发布
**来源**: ClawHub

**功能**：
- 🛠️ **技能创建**：从头创建新的 AgentSkills
- ✏️ **技能编辑**：改进、审查、审计现有技能或 SKILL.md 文件
- 📁 **目录重构**：将文件移动到 references/ 或 scripts/，删除过时内容
- ✅ **规范验证**：根据 AgentSkills 规范进行验证

**核心原则**：
- **简洁至上**：上下文窗口是公共资源，只添加必要的内容
- **适当的自由度**：根据任务特性匹配具体程度
- **模块化设计**：技能应该模块化、自包含

**技能解剖结构**：
```
skill-name/
├── SKILL.md（必需）
│   ├── YAML frontmatter（必需）
│   └── Markdown 指令（必需）
└── 捆绑资源（可选）
    ├── scripts/    - 可执行代码
    ├── references/ - 参考文档
    └── assets/     - 输出文件
```

**安装**：
```bash
cd ~/.openclaw/skills/
git clone https://github.com/L-LesterYu/OpenClaw-hot-skills-zh.git temp-repo
cp -r temp-repo/skills/skill-creator-zh ./
rm -rf temp-repo
```

**快速开始**：
```bash
# Agent 会自动识别技能创建任务
# 示例交互：
用户：创建一个新技能用于处理 Excel 文件

Agent 会：
1. 分析任务需求和技能结构
2. 创建 SKILL.md 文件（包含 frontmatter 和指令）
3. 创建必要的目录结构
4. 提供使用示例和最佳实践
5. 根据 AgentSkills 规范验证
```

**使用场景**：
- "创建一个新技能"
- "改进这个技能"
- "审查这个 SKILL.md 文件"
- "清理技能目录"
- "审计技能质量"

**触发短语**：
- "创建技能" / "编写技能"
- "改进技能" / "优化技能"
- "审查技能" / "审计技能"
- "清理技能" / "整理技能"

**技能类型**：
- **高自由度**：基于文本的指令（多种方法有效）
- **中等自由度**：伪代码或带参数的脚本
- **低自由度**：特定脚本，少量参数

**最佳实践**：
- 优先使用简洁的示例而非冗长的解释
- 将详细的参考材料放在 references/ 目录
- 将重复使用的代码放在 scripts/ 目录
- 保持 SKILL.md 精简，只包含必要的程序性指令

**链接**：
- AgentSkills 规范：https://agentskills.io/specification

[查看详细文档](./skills/skill-creator-zh/SKILL.md)

---

### 41. marketing-mode-zh - 营销模式

**版本**: v1.0.0
**状态**: ✅ 已发布
**来源**: [ClawHub - Marketing Mode](https://clawhub.ai)
**作者**: Seth Rose (110 stars)

**功能**：
- 📈 **营销战略**：提供 140+ 个经过验证的营销创意和框架
- 🧠 **心理学原理**：整合 20+ 个说服心理学原理和心智模型
- 📝 **内容与 SEO**：SEO 审计、程序化 SEO、文案撰写框架
- 🎯 **转化优化（CRO）**：页面 CRO 元素、漏斗优化、A/B 测试设置
- 💰 **付费广告与增长**：Google、Meta、LinkedIn 广告策略
- 📧 **邮件营销**：5 种序列类型、欢迎系列、培养序列

**核心框架**：
- **AIDA**：注意 → 兴趣 → 欲望 → 行动
- **PAS**：问题 → 煽动 → 解决方案
- **钩子模型**：触发 → 行动 → 多变的奖励 → 投入
- **价值主张画布**：将客户需求与产品相匹配
- **定位框架**：目标用户、功能、差异化

**6 大营销领域**：
1. **营销战略与框架**：140+ 个创意、5 阶段发布战略、定价战略
2. **心理学与心智模型**：第一性原理、JTBD、帕累托法则等
3. **SEO 与内容**：SEO 审计、程序化 SEO、文案撰写
4. **转化优化（CRO）**：页面元素、漏斗优化、A/B 测试
5. **付费广告与增长**：Google/Meta/LinkedIn 广告、推荐计划
6. **邮件营销**：欢迎系列、培养序列、销售序列

**安装**：
```bash
cd ~/.openclaw/skills/
git clone https://github.com/L-LesterYu/OpenClaw-hot-skills-zh.git temp-repo
cp -r temp-repo/skills/marketing-mode-zh ./
rm -rf temp-repo
```

**快速开始**：
```bash
# Agent 会自动激活营销模式
# 示例交互：
用户：我需要为我的 SaaS 产品制定营销策略

Agent（Mark）会：
1. 询问产品、受众、阶段、预算和目标
2. 从知识库推荐合适的营销框架
3. 提供具体的营销创意和战术
4. 帮助优化 CTA、钩子和转化漏斗
5. 建议定价和定位策略
```

**使用场景**：
- "帮我撰写更好的 CTA 和钩子"
- "优化产品和服务的定位"
- "改进营销漏斗和转化率"
- "制定信息传递策略"
- "设计 A/B 测试方案"
- "选择合适的营销渠道"
- "优化 SEO 和内容营销"

**Mark 的经典语录**：
- "CTA 是什么？"
- "钩子是什么？"
- "这是给谁用的？"
- "让它具体化。让它令人难忘。"

**营销创意类别**：
- 内容与 SEO（简单关键词排名、SEO 审计、术语表营销等）
- 竞争对手与对比（对比页面、营销柔道、广告研究等）
- 免费工具与工程（副项目、Chrome 扩展、微型网站等）
- 付费广告（Facebook/Google/LinkedIn 广告等）
- 社交与病毒式传播（推荐计划、病毒循环、社区营销等）

**适用对象**：
- 🚀 创始人和创业者
- 📈 营销人员和增长黑客
- 💼 产品经理
- 🎯 内容创作者

**技术要求**：
- Node.js >= 18.0.0
- npm 或 yarn

[查看详细文档](./skills/marketing-mode-zh/SKILL.md) | [营销框架](./skills/marketing-mode-zh/README.md)

---

### 42. self-improving-zh - 自我改进与主动学习 Agent

**版本**: v1.2.16
**状态**: ✅ 已发布
**来源**: [ClawIC - Self-Improving](https://clawic.com/skills/self-improving)

**功能**：
- 🧠 **自我反思**：评估自身工作、发现错误并持续改进
- 📝 **自组织记忆**：采用分层结构存储学习内容（热/温/冷数据）
- 🔄 **自动提升/降级**：根据使用频率自动调整记忆层级
- 📊 **纠正追踪**：记录用户纠正和反馈，避免重复错误
- 🎯 **模式识别**：识别重复行为并转化为规则
- 🔒 **安全边界**：不存储敏感信息，保护隐私

**核心机制**：
- **学习信号识别**：自动识别纠正、偏好和模式候选
- **三层记忆架构**：
  - 热数据（memory.md）：≤100行，始终加载
  - 温数据（projects/domains/）：按需加载
  - 冷数据（archive/）：明确查询时加载
- **命名空间隔离**：项目模式、领域模式和全局偏好分离存储
- **透明引用**：每个来自记忆的操作都引用来源

**安装**：
```bash
cd ~/.openclaw/skills/
git clone https://github.com/L-LesterYu/OpenClaw-hot-skills-zh.git temp-repo
cp -r temp-repo/skills/self-improving-zh ./
rm -rf temp-repo
```

**快速开始**：
```bash
# 初始化记忆存储
mkdir -p ~/self-improving/projects
mkdir -p ~/self-improving/domains
mkdir -p ~/self-improving/archive
touch ~/self-improving/memory.md
touch ~/self-improving/corrections.md
touch ~/self-improving/index.md
```

**使用场景**：
- 用户纠正你时："不对，应该这样做..." → 自动记录纠正
- 发现改进空间时："我可以做得更好" → 记录自我反思
- 查询记忆时："你知道关于X的什么？" → 搜索所有层级
- 查看模式："显示我的模式" → 列出已学习的规则

**学习信号示例**：
- 纠正："你弄错了..."、"实际上应该是..."
- 偏好："我喜欢你..."、"我的风格是..."
- 模式候选：同一指令重复3次以上

**适用对象**：
- 🤖 AI Agent 开发者
- 📈 需要持续改进的系统
- 🎯 追求个性化的用户
- 💼 复杂项目的管理者

**技术要求**：
- 无需额外依赖
- 无需凭据
- 可选网络访问（用于主动技能）

[查看详细文档](./skills/self-improving-zh/SKILL.md) | [学习机制](./skills/self-improving-zh/learning.md) | [安全边界](./skills/self-improving-zh/boundaries.md)

---

### 43. marketing-skills-zh - 23个营销实战手册

**版本**: v0.1.2
**状态**: ✅ 已发布
**来源**: [ClawHub - Marketing Skills](https://clawhub.ai)

**功能**：
- 📚 **23 个营销模块**：涵盖 CRO、SEO、文案、分析、实验、定价、发布、广告、社交媒体等
- ✅ **检查清单**：快速获取实用的检查清单和可复制粘贴的交付物
- 🎯 **框架应用**：应用营销框架并输出实用成果（草稿 + 检查清单）
- 📦 **整合内容**：整合 `coreyhaines31/marketingskills` 的完整内容

**包含的模块**：
- **转化优化（CRO）**：page-cro、signup-flow-cro、onboarding-cro、form-cro、popup-cro、paywall-upgrade-cro
- **内容与文案**：copywriting、copy-editing、email-sequence、social-content
- **SEO 优化**：seo-audit、programmatic-seo、schema-markup
- **策略与定价**：pricing-strategy、launch-strategy、marketing-ideas、marketing-psychology
- **付费广告**：paid-ads
- **其他工具**：analytics-tracking、ab-test-setup、competitor-alternatives、referral-program、free-tool-strategy

**安装**：
```bash
cd ~/.openclaw/skills/
git clone https://github.com/L-LesterYu/OpenClaw-hot-skills-zh.git temp-repo
cp -r temp-repo/skills/marketing-skills-zh ./
rm -rf temp-repo
```

**快速开始**：
```bash
# 根据需求选择对应模块
# 例如：优化落地页转化
# 读取 references/page-cro/SKILL.md

# 例如：设置 A/B 测试
# 读取 references/ab-test-setup/SKILL.md

# 例如：创建邮件序列
# 读取 references/email-sequence/SKILL.md
```

**使用场景**：
- "帮我优化落地页的转化率"
- "设计一个 A/B 测试方案"
- "创建邮件营销序列"
- "制定产品发布策略"
- "优化定价策略"

**输出规则**：
- 优先考虑 80/20 原则：先抓最大杠杆点
- 不捏造指标或关键词量，缺失时标注为假设
- 尽可能包含可复制粘贴的草稿和实施检查清单

**适用对象**：
- 🚀 创业者和营销人员
- 📈 增长黑客
- 💼 产品经理
- 🎯 内容创作者

[查看详细文档](./skills/marketing-skills-zh/SKILL.md)

---

### 44. proactive-self-improving-agent-zh - 主动自我改进智能体

**版本**: v1.0.0
**状态**: ✅ 已发布
**来源**: [ClawHub - Proactive Self-Improving Agent](https://clawhub.ai)

**功能**：
- 🧠 **自动捕获经验**：识别错误、纠正和最佳实践，结构化记录
- 🔄 **安全进化机制**：经验反复出现时自动晋升为长期能力，有护栏防止漂移
- 📝 **多层记录系统**：LEARNINGS.md、ERRORS.md、FEATURE_REQUESTS.md 分层管理
- 📊 **操作日志追踪**：CHANGELOG.md 记录所有学习操作，支持脚本读取
- 🔒 **安全护栏**：ADL 协议防止漂移，VFM 协议确保价值优先

**触发条件（7种场景）**：
- 命令/操作失败 → 记录到 ERRORS.md
- 用户纠正（"不对"/"应该是"）→ 记录到 LEARNINGS.md
- 用户需要不存在的能力 → 记录到 FEATURE_REQUESTS.md
- 外部 API/工具出错 → 记录到 ERRORS.md
- 发现知识过时/错误 → 记录到 LEARNINGS.md
- 发现更好做法 → 记录到 LEARNINGS.md
- 任务完成时回顾 → 有新经验则记录到 LEARNINGS.md

**核心机制**：
- **去重法则**：触发 ≠ 必须写入，无新经验则跳过
- **晋升路径**：经验 ≥3 次重复 → 晋升到 AGENTS.md/TOOLS.md/SOUL.md
- **技能提取**：足够通用的经验可提取为独立 skill
- **VFM 打分**：检索复用性 3x + 错误预防 3x + 分析质量 2x + 效率提升 2x

**安装**：
```bash
cd ~/.openclaw/skills/
git clone https://github.com/L-LesterYu/OpenClaw-hot-skills-zh.git temp-repo
cp -r temp-repo/skills/proactive-self-improving-agent-zh ./
rm -rf temp-repo
```

**快速开始**：
```bash
# 初始化学习存储
mkdir -p .learnings
touch .learnings/LEARNINGS.md
touch .learnings/ERRORS.md
touch .learnings/FEATURE_REQUESTS.md
touch .learnings/CHANGELOG.md
```

**使用场景**：
- 命令报错时 → 自动记录到 ERRORS.md
- 用户说"不对/应该是…"时 → 记录纠正到 LEARNINGS.md
- 发现更好做法时 → 记录最佳实践
- 任务完成后 → 回顾过程，有经验则记录
- 同一问题出现 ≥3 次 → 自动晋升为永久规则

**适用对象**：
- 🤖 AI Agent 开发者
- 📈 需要持续改进的智能系统
- 🎯 追求自我进化的 Agent
- 💼 复杂项目的自动化管理

**技术要求**：
- 无需额外依赖
- 无需凭据
- 可选 git 支持（用于经验追踪）

[查看详细文档](./skills/proactive-self-improving-agent-zh/SKILL.md)

---

### 45. mcporter-zh - MCP 服务器管理工具

**版本**: v1.0.0
**状态**: ✅ 已发布
**来源**: [GitHub - mcporter](https://github.com/steipete/mcporter)

**功能**：
- 🔍 **零配置发现**：自动发现系统中已配置的 MCP 服务器（Cursor、Claude、Codex 等）
- 🛠️ **工具调用**：直接从命令行调用任何 MCP 工具
- 🔐 **OAuth 认证**：内置 OAuth 缓存，支持托管 MCP 服务认证
- 📝 **代码生成**：生成 TypeScript 类型定义和 CLI 工具
- 🔧 **配置管理**：管理和编辑 MCP 服务器配置

**核心能力**：
- 列出所有已配置的 MCP 服务器
- 调用 MCP 工具（支持冒号语法和函数调用语法）
- 支持 HTTP、SSE 和 stdio 传输
- 自动检测需要 OAuth 的托管服务
- 生成强类型的 TypeScript 客户端

**安装**：
```bash
# 使用 npm
npm install -g mcporter

# 或使用 npx 直接运行
npx mcporter list

# 安装 Skill
cd ~/.openclaw/skills/
git clone https://github.com/L-LesterYu/OpenClaw-hot-skills-zh.git temp-repo
cp -r temp-repo/skills/mcporter-zh ./
rm -rf temp-repo
```

**快速开始**：
```bash
# 列出所有已配置的服务器
npx mcporter list

# 列出特定服务器的详细信息
npx mcporter list linear --schema

# 调用工具
npx mcporter call linear.create_comment issueId:ENG-123 body:'评论内容'

# 对需要 OAuth 的服务器进行认证
npx mcporter auth https://mcp.linear.app/mcp

# 生成 TypeScript 类型定义
npx mcporter emit-ts linear --output ./types
```

**使用场景**：
- "列出我系统中配置的 MCP 服务器"
- "调用 Linear 的创建评论工具"
- "为 MCP 服务器生成 TypeScript 客户端"
- "管理 MCP 服务器配置"
- "连接到临时的 MCP 端点"

**支持的配置来源**：
- `~/.mcporter/mcporter.json`（主目录配置）
- `config/mcporter.json`（项目配置）
- Cursor、Claude Code/Desktop、Codex、Windsurf、OpenCode、VS Code 的配置

**调用语法**：
```bash
# 冒号分隔的标志（shell 友好）
npx mcporter call linear.create_comment issueId:ENG-123 body:'看起来不错！'

# 函数调用风格（匹配 `mcporter list` 的签名）
npx mcporter call 'linear.create_comment(issueId: "ENG-123", body: "看起来不错！")'
```

**注意事项**：
- 支持环境变量占位符 `${API_KEY}`
- 使用 `--json` 获取结构化输出
- 使用 `--verbose` 查看详细配置信息
- 使用 `--persist` 将临时服务器保存到配置

[查看详细文档](./skills/mcporter-zh/SKILL.md)

---

### 46. gongwen-zh - 中国公文格式化

**版本**: v1.0.0
**状态**: ✅ 已发布
**来源**: 依据 GB/T 9704-2012《党政机关公文格式》

**功能**：
- 📄 **标准格式**：按照中国党政机关公文格式标准进行格式化
- 🎨 **字体规范**：规范正文字体、标题层级、页面设置
- ✍️ **落款格式**：标准化的落款单位和日期格式
- 📎 **附件处理**：附件标注和格式规范
- 🔢 **数字规范**：阿拉伯数字的统一处理

**核心规范**：
- **页面设置**：上3.7cm、下3.2cm、左2.8cm、右2.6cm
- **正文字体**：仿宋_GB2312，首行缩进2字符，行间距29磅
- **标题层级**：一级黑体3号、二级楷体_GB2312、三级仿宋_GB2312
- **页码格式**：—1—（数字前后加横线），单页靠右，双页靠左

**安装**：
```bash
cd ~/.openclaw/skills/
git clone https://github.com/L-LesterYu/OpenClaw-hot-skills-zh.git temp-repo
cp -r temp-repo/skills/gongwen-zh ./
rm -rf temp-repo
```

**快速开始**：
```bash
# Agent 会自动识别公文格式化需求
# 示例交互：
用户：帮我格式化这份通知

Agent 会：
1. 应用标准页面设置
2. 规范标题层级和字体
3. 调整段落格式和行间距
4. 添加标准的落款格式
5. 处理附件标注
```

**使用场景**：
- "格式化这份通知"
- "将这份文档改为公文格式"
- "检查这份公文的格式是否规范"
- "生成标准的公文模板"
- "将 Markdown 转换为公文格式"

**支持的公文类型**：
- 通知、通报、决定、意见
- 方案、报告、请示、函

**标题层级格式**：
| 层级 | 格式 | 字体 | 示例 |
|------|------|------|------|
| 一级标题 | 一、二、三、四 | **黑体 3号** | 一、总体要求 |
| 二级标题 | （一）（二）（三） | **楷体_GB2312** | （一）指导思想 |
| 三级标题 | 1. 2. 3. | **仿宋_GB2312** | 1. 完善工作机制 |

**注意事项**：
- 全文字体风格保持统一，不可混用
- 标题层级一般不超过三级
- 公文中的数字一般使用阿拉伯数字
- 标点符号使用中文全角标点

[查看详细文档](./skills/gongwen-zh/SKILL.md)

---

### 47. baidu-search-zh - 百度 AI 搜索

**版本**: v1.0.0
**状态**: ✅ 已发布
**来源**: 百度 AI 搜索 API

**功能**：
- 🔍 **AI 优化搜索**：使用百度 AI 搜索 API 进行高质量网页搜索
- 🌐 **中文优化**：完美支持中文搜索查询和结果展示
- 📊 **结构化输出**：返回标题、链接、内容摘要的 JSON 格式
- ⏰ **时效筛选**：支持按时间范围筛选（过去一天/周/月/年）
- 🔒 **安全认证**：通过环境变量安全传递 API 密钥

**核心命令**：
```bash
# 基础搜索
python3 skills/baidu-search-zh/scripts/search.py '{"query":"人工智能"}'

# 时间范围格式一 "YYYY-MM-DDtoYYYY-MM-DD" 示例
python3 scripts/search.py '{"query":"最新新闻","freshness":"2025-09-01to2025-09-08"}'

# 时间范围格式二 pd、pw、pm、py 示例
python3 scripts/search.py '{"query":"最新新闻","freshness":"pd"}'

# 设置 count 参数，指定返回结果数量
python3 scripts/search.py '{"query":"旅游景点","count": 20}'
```

**安装**：
```bash
cd ~/.openclaw/skills/
git clone https://github.com/L-LesterYu/OpenClaw-hot-skills-zh.git temp-repo
cp -r temp-repo/skills/baidu-search-zh ./
rm -rf temp-repo
```

**配置 API Key**：
```bash
# 设置环境变量
export BAIDU_API_KEY="your-baidu-api-key"
```

**请求参数**：

| 参数 | 类型 | 必填 | 默认值 | 说明 |
|-------|------|----------|---------|-------------|
| query | str | 是 | - | 搜索关键词 |
| count | int | 否 | 10 | 返回结果数量，范围 1-50 |
| freshness | str | 否 | 无 | 时间范围，两种格式：格式一为 "YYYY-MM-DDtoYYYY-MM-DD"；格式二包括 pd、pw、pm、py，分别代表过去 24 小时、过去 7 天、过去 31 天和过去 365 天 |

**使用场景**：
- "帮我搜索一下最新的 AI 新闻"
- "查找关于 Python 异步编程的资料"
- "今天有什么热点新闻？"
- "搜索 OpenClaw 的使用教程"

**获取 API Key**：
1. 访问 [百度智能云](https://cloud.baidu.com/) 注册账号
2. 开通千帆平台 AI 搜索服务
3. 在控制台获取 API 密钥

**注意事项**：
- 需要设置 `BAIDU_API_KEY` 环境变量
- 需要安装 Python 3 和 requests 库
- API 有调用频率限制，请合理使用

[查看详细文档](./skills/baidu-search-zh/SKILL.md)

---

### 48. slack-zh - Slack 操作工具

**版本**: v1.0.0
**状态**: ✅ 已发布
**来源**: [ClawHub - Slack](https://clawhub.ai)

**功能**：
- 💬 **消息管理**：发送、编辑、删除 Slack 消息
- 📌 **置顶管理**：置顶/取消置顶重要消息
- ✅ **表情回复**：为消息添加表情反应
- 👥 **成员信息**：获取 Slack 成员详细信息
- 🎨 **自定义表情**：列出工作区的自定义表情列表

**核心能力**：
- 通过 `slack` 工具与 Slack API 交互
- 支持频道和私信消息操作
- 完整的消息 CRUD 操作
- 置顶消息管理
- 成员信息查询

**安装**：
```bash
cd ~/.openclaw/skills/
git clone https://github.com/L-LesterYu/OpenClaw-hot-skills-zh.git temp-repo
cp -r temp-repo/skills/slack-zh ./
rm -rf temp-repo
```

**快速开始**：
```bash
# 配置 Slack bot token（在 OpenClaw 配置中）
# 使用 Clawdbot 配置的 bot token

# 发送消息
slack action="sendMessage" to="channel:C123" content="来自 Clawdbot 的问候"

# 回复消息
slack action="react" channelId="C123" messageId="1712023032.1234" emoji="✅"

# 置顶消息
slack action="pinMessage" channelId="C123" messageId="1712023032.1234"
```

**使用场景**：
- "发送 Slack 消息到 #general 频道"
- "为这条消息添加 ✅ 表情"
- "置顶这条重要公告"
- "获取成员信息"
- "列出置顶的消息"

**操作组**：
- `reactions` - 回复消息 + 列出回复
- `messages` - 读取/发送/编辑/删除消息
- `pins` - 置顶/取消置顶/列出置顶
- `memberInfo` - 成员信息
- `emojiList` - 自定义表情列表

**注意事项**：
- 需要配置 Slack bot token
- 消息时间戳格式：`1712023032.1234`
- 支持频道（`channel:<id>`）和用户（`user:<id>`）目标

[查看详细文档](./skills/slack-zh/SKILL.md)

---

### 49. n8n-workflow-automation-zh - n8n 工作流自动化

**版本**: v1.0.0
**状态**: ✅ 已发布
**来源**: [ClawHub - n8n Workflow Automation](https://clawhub.ai)

**功能**：
- 🔧 **工作流设计**：设计并输出 n8n 工作流 JSON
- 🔄 **健壮机制**：具备重试、日志记录、错误处理机制
- 🎯 **幂等性保证**：确保重新运行时不会重复记录
- 👥 **人机协作**：支持审核队列和人工审批步骤
- 📊 **可审计性**：记录每次运行，不会静默失败

**核心能力**：
- 支持 Cron/webhook/手动触发器
- 定义数据契约和验证规则
- 设计去重键防止重复
- 添加可观察性（run_id、状态日志）
- 每节点错误分支和重试机制

**安装**：
```bash
cd ~/.openclaw/skills/
git clone https://github.com/L-LesterYu/OpenClaw-hot-skills-zh.git temp-repo
cp -r temp-repo/skills/n8n-workflow-automation-zh ./
rm -rf temp-repo
```

**使用场景**：
- "构建一个每周一运行并发送合规摘要的 n8n 工作流"
- "为此工作流添加错误处理和重试机制"
- "创建一个记录每次运行的 webhook 工作流"
- "使此流程具有幂等性"
- "配备审计日志和人工审批步骤"

**输出格式**：
- 默认（只读）：工作流设计规范
- 明确请求：`workflow.json` + `runbook.md`

**注意事项**：
- 不适用于纯代码自动化（请使用脚本/CI 技能）
- 不用于绕过安全控制或隐藏审计轨迹
- 需要明确触发器类型和成功标准

[查看详细文档](./skills/n8n-workflow-automation-zh/SKILL.md)

---

### 50. model-usage-zh - 模型使用情况统计

**版本**: v1.0.0
**状态**: ✅ 已发布
**来源**: [ClawHub - Model Usage](https://clawhub.ai)

**功能**：
- 📊 **成本统计**：从 CodexBar 本地成本日志获取按模型使用成本
- 🎯 **当前模型**：获取最新模型使用情况
- 📈 **完整细分**：支持所有模型的详细使用统计
- 🔧 **多提供商**：支持 Codex 和 Claude 两大提供商
- 📝 **灵活输出**：支持文本和 JSON 格式输出

**核心能力**：
- 汇总 Codex 或 Claude 的按模型使用情况
- 获取当前（最新）模型或完整模型细分
- 从 codexbar cost JSON 获取可脚本化摘要
- 支持文件输入和标准输入
- 提供文本或 JSON 格式输出

**安装**：
```bash
cd ~/.openclaw/skills/
git clone https://github.com/L-LesterYu/OpenClaw-hot-skills-zh.git temp-repo
cp -r temp-repo/skills/model-usage-zh ./
rm -rf temp-repo
```

**快速开始**：
```bash
# 获取当前模型使用情况
python ~/.openclaw/skills/model-usage-zh/scripts/model_usage.py --provider codex --mode current

# 获取所有模型统计
python ~/.openclaw/skills/model-usage-zh/scripts/model_usage.py --provider codex --mode all

# JSON 格式输出
python ~/.openclaw/skills/model-usage-zh/scripts/model_usage.py --provider claude --mode all --format json --pretty
```

**使用场景**：
- "查看当前使用的模型"
- "统计所有模型的成本"
- "获取 Codex 的使用情况"
- "查看 Claude 的成本细分"

**注意事项**：
- 需要安装 CodexBar CLI
- 目前主要支持 macOS（darwin）
- Linux 支持待添加
- 数值仅为每个模型的成本，未按模型拆分 token

[查看详细文档](./skills/model-usage-zh/SKILL.md)

---

### 51. playwright-mcp-zh - Playwright 浏览器自动化

**版本**: v1.0.0
**状态**: ✅ 已发布
**来源**: [NPM - @playwright/mcp](https://www.npmjs.com/package/@playwright/mcp)

**功能**：
- 🎭 **浏览器自动化**：通过 Playwright MCP 服务器实现完整的浏览器自动化工作流
- 🌐 **多浏览器支持**：支持 Chrome、Firefox 和 WebKit 浏览器
- 🖱️ **页面交互**：导航网站、点击元素、填写表单、提取数据
- 📸 **截图功能**：支持页面截图和视频录制
- 🔧 **灵活配置**：支持无头模式、视口设置、超时控制等多种配置选项

**核心工具**：
- `browser_navigate` - 导航到 URL
- `browser_click` - 点击元素
- `browser_type` - 输入文本
- `browser_select_option` - 选择下拉选项
- `browser_get_text` - 提取文本内容
- `browser_evaluate` - 执行 JavaScript
- `browser_snapshot` - 获取无障碍页面快照
- `browser_close` - 关闭浏览器上下文
- `browser_choose_file` - 上传文件
- `browser_press` - 按下键盘按键

**安装**：
```bash
# 安装 Playwright MCP
npm install -g @playwright/mcp

# 安装浏览器
npx playwright install chromium

# 安装 Skill
cd ~/.openclaw/skills/
git clone https://github.com/L-LesterYu/OpenClaw-hot-skills-zh.git temp-repo
cp -r temp-repo/skills/playwright-mcp-zh ./
rm -rf temp-repo
```

**快速开始**：
```bash
# 启动 MCP 服务器
npx @playwright/mcp

# 无头模式
npx @playwright/mcp --headless

# 指定浏览器和视口
npx @playwright/mcp --browser firefox --viewport-size 1280x720
```

**使用场景**：
- "自动化测试网站功能"
- "批量填写表单"
- "从网页提取数据"
- "截取网页截图"
- "执行浏览器自动化脚本"

**配置选项**：
- `--headless` - 无头模式运行
- `--browser` - 指定浏览器（chromium/firefox/webkit）
- `--viewport-size` - 设置视口大小
- `--timeout-action` - 操作超时时间（毫秒）
- `--timeout-navigation` - 导航超时时间（毫秒）
- `--allowed-hosts` - 允许访问的域名列表
- `--ignore-https-errors` - 忽略 HTTPS 错误

[查看详细文档](./skills/playwright-mcp-zh/SKILL.md)

---

### 52. ui-ux-pro-max-zh - UI/UX 设计智能指南

**版本**: v2.0
**状态**: ✅ 已发布
**来源**: [GitHub - ui-ux-pro-max-skill](https://github.com/nextlevelbuilder/ui-ux-pro-max-skill)

**功能**：
- 🎨 **设计系统生成**：AI 驱动的设计系统生成器，根据项目需求自动生成完整的设计系统
- 🎯 **161 种配色方案**：基于产品类型和风格的配色推荐
- ✍️ **57 种字体搭配**：Google Fonts 字体组合推荐
- 📱 **多技术栈支持**：覆盖 10 大技术栈（React、Next.js、Vue、Svelte、SwiftUI、React Native、Flutter、Tailwind、shadcn/ui、HTML/CSS）
- 🎭 **50+ UI 风格**：玻璃拟态、粘土拟态、极简主义、野兽派等多种设计风格
- 📊 **25 种图表类型**：数据可视化图表推荐和最佳实践
- ✅ **161 条 UX 准则**：基于优先级的设计规则和无障碍指南
- 🔍 **智能搜索引擎**：BM25 + 正则混合搜索，支持域名自动检测

**核心功能模块**：
- **产品类型推荐**：SaaS、电商、作品集、博客等
- **UI 风格库**：玻璃拟态、极简主义、野兽派等 + AI 提示词和 CSS 关键字
- **字体系统**：Google Fonts 字体搭配和导入代码
- **配色系统**：按产品类型的配色方案
- **落地页结构**：页面结构和 CTA 策略
- **图表推荐**：图表类型和库推荐
- **UX 最佳实践**：设计规则和反模式

**安装**：
```bash
cd ~/.openclaw/skills/
git clone https://github.com/L-LesterYu/OpenClaw-hot-skills-zh.git temp-repo
cp -r temp-repo/skills/ui-ux-pro-max-zh ./
rm -rf temp-repo
```

**快速开始**：
```bash
# 搜索设计资源
python3 src/ui-ux-pro-max/scripts/search.py "SaaS 仪表盘" --domain product

# 搜索 UI 风格
python3 src/ui-ux-pro-max/scripts/search.py "玻璃拟态" --domain style

# 搜索字体搭配
python3 src/ui-ux-pro-max/scripts/search.py "现代简约" --domain typography

# 搜索配色方案
python3 src/ui-ux-pro-max/scripts/search.py "科技感" --domain color

# 指定技术栈搜索
python3 src/ui-ux-pro-max/scripts/search.py "按钮组件" --stack react
```

**使用场景**：
- "设计一个新的落地页"
- "创建响应式导航栏"
- "优化表单的用户体验"
- "选择合适的配色方案"
- "实现无障碍设计"
- "重构 UI 组件库"
- "构建设计系统"

**支持的操作**：
- 规划、构建、创建、设计
- 实现、审查、修复、改进
- 优化、增强、重构
- 检查 UI/UX 代码

**项目类型支持**：
- 网站、落地页、仪表盘
- 管理后台、电商、SaaS
- 作品集、博客、移动应用

**设计风格涵盖**：
- 玻璃拟态、粘土拟态、极简主义
- 野兽派、新拟态、Bento Grid
- 暗色模式、响应式设计
- 拟物化、扁平化设计

**设计主题**：
- 颜色系统、无障碍设计
- 动画、布局、排版
- 字体搭配、间距
- 交互状态、阴影、渐变

**优先级规则类别**：
1. 无障碍性（关键）
2. 触控与交互（关键）
3. 性能（高）
4. 风格选择（高）
5. 布局与响应式（高）
6. 排版与颜色（中）
7. 动画（中）
8. 表单与反馈（中）
9. 导航模式（高）
10. 图表与数据（低）

**集成能力**：
- shadcn/ui MCP 组件搜索和示例
- 多平台 AI 编码助手支持（Claude Code、Windsurf、Cursor 等）
- CLI 安装工具（uipro-cli）

**前置要求**：
- Python 3.x（无需外部依赖）

**链接**：
- [项目官网](https://uupm.cc)
- [GitHub 仓库](https://github.com/nextlevelbuilder/ui-ux-pro-max-skill)
- [NPM CLI](https://www.npmjs.com/package/uipro-cli)

[查看详细文档](./skills/ui-ux-pro-max-zh/README.md) | [技能文档](./skills/ui-ux-pro-max-zh/.claude/skills/ui-ux-pro-max/SKILL.md)

---

### 53. self-improving-proactive-agent-zh - 自我改进与主动学习智能体（增强版）

**版本**: v1.2.16
**状态**: ✅ 已发布
**来源**: [ClawIC - Self-Improving](https://clawic.com/skills/self-improving)

**功能**：
- 🧠 **自我反思**：智能体评估自己的工作，捕捉错误并永久改进
- 📝 **纠正记录**：自动记录用户纠正，评估是否形成规则
- 🏗️ **分层记忆**：HOT/WARM/COLD 三层存储，按需加载
- 🔄 **自动提升/降级**：频繁使用的模式自动提升，长期未用的自动降级
- 📊 **记忆统计**：各层级记忆数量、最近活动一目了然
- 💓 **心跳集成**：定期维护记忆，清理过期内容
- 🔒 **安全边界**：不存储凭证、健康数据等敏感信息

**记忆架构**：
```
~/self-improving/
├── memory.md          # HOT: ≤100 行，始终加载
├── index.md           # 主题索引及行数
├── heartbeat-state.md # 心跳状态跟踪
├── projects/          # 按项目的学习内容
├── domains/           # 按领域（代码、写作、通信）
├── archive/           # COLD: 衰退的模式
└── corrections.md     # 最近 50 条纠正日志
```

**安装**：
```bash
cd ~/.openclaw/skills/
git clone https://github.com/L-LesterYu/OpenClaw-hot-skills-zh.git temp-repo
cp -r temp-repo/skills/self-improving-proactive-agent-zh ./
rm -rf temp-repo
```

**快速开始**：
```bash
# 初始化记忆存储
mkdir -p ~/self-improving/projects
mkdir -p ~/self-improving/domains
mkdir -p ~/self-improving/archive
touch ~/self-improving/memory.md
touch ~/self-improving/corrections.md
touch ~/self-improving/index.md
```

**使用场景**：
- 用户纠正你时 → 自动记录纠正，形成长期规则
- 完成重要任务后 → 自我反思，识别改进点
- 用户询问记忆 → "你学到了什么？" / "显示我的模式"
- 记忆管理 → "记忆统计" / "导出记忆" / "忘记 X"

**相关技能**：
- `self-improving-zh`（基础版）和 `proactive-self-improving-agent-zh` 的功能增强合体版

[查看详细文档](./skills/self-improving-proactive-agent-zh/SKILL.md) | [学习机制](./skills/self-improving-proactive-agent-zh/learning.md) | [安全边界](./skills/self-improving-proactive-agent-zh/boundaries.md) | [设置指南](./skills/self-improving-proactive-agent-zh/setup.md)

---

### 54. admapix-zh - 广告创意素材搜索

**版本**: v1.0.0
**状态**: ✅ 已发布
**来源**: [AdMapix](https://www.admapix.com)

**功能**：
- 🎯 **广告素材搜索**：通过自然语言搜索竞品广告创意素材
- 🎬 **多类型筛选**：支持视频、图片、试玩广告等多种素材类型
- 🌏 **地区筛选**：东南亚、北美、欧洲、日韩、中东等地区
- 📊 **灵活排序**：按相关性、曝光量、首次发现、投放天数排序
- 📱 **H5 结果页**：搜索结果生成可视化 H5 页面，支持在线预览

**核心能力**：
- 中英文双语参数映射，自然语言自动解析
- API 端点：`https://api.admapix.com/api/data/search`
- 支持翻页、调整筛选条件等后续交互
- 大数字人性化展示（万、亿）

**安装**：
```bash
cd ~/.openclaw/skills/
git clone https://github.com/L-LesterYu/OpenClaw-hot-skills-zh.git temp-repo
cp -r temp-repo/skills/admapix-zh ./
rm -rf temp-repo
```

**配置 API Key**：
```bash
openclaw config set skills.entries.admapix-zh.apiKey "你的ADMAPIX_API_KEY"
```

**使用场景**：
- "搜一下 puzzle game 的视频广告"
- "找东南亚投放的休闲游戏素材"
- "看看 temu 最近的广告创意"
- "搜最近一周曝光最多的电商广告"

**获取 API Key**：
1. 访问 [www.admapix.com](https://www.admapix.com) 注册并获取 API Key
2. 使用 `openclaw config set` 命令配置

**注意事项**：
- 需要配置 `ADMAPIX_API_KEY` 环境变量
- H5 结果页 24 小时后自动过期
- 每页最多返回 60 条结果

[查看详细文档](./skills/admapix-zh/SKILL.md) | [参数映射表](./skills/admapix-zh/references/param-mappings.md)

---

### 55. minara-zh - 加密货币交易与钱包

**版本**: v2.5.2
**状态**: ✅ 已发布
**来源**: [Minara](https://minara.ai)

**功能**：
- 💱 **代币兑换**：支持 EVM + Solana 多链代币兑换
- 📈 **永续合约**：永续合约交易（做多/做空）、限价单
- 💸 **转账支付**：加密货币转账、x402 支付协议
- 💳 **充值提现**：信用卡充值、加密货币充值、法币提现
- 🤖 **AI 聊天**：内置 AI 聊天助手，支持市场研究和分析
- 🔍 **市场发现**：热门代币、价格查询、DeFi 收益分析
- 📊 **自动交易**：支持自动化交易策略
- 👛 **内置钱包**：通过 Minara CLI 管理多链钱包

**支持的链**：
- Ethereum、Solana、Base、Arbitrum、Optimism
- Polygon、BSC、Avalanche、Berachain、Hyperliquid

**核心命令**：
```bash
# 代币兑换
minara swap -s buy -t 'ETH' -a 100
minara swap -s sell -t 'BONK' -a all

# 查询余额
minara balance

# 转账
minara transfer -t 'USDC' -a 50 -r <wallet_address>

# 永续合约
minara perps order -s long -t 'ETH' -a 100 -l 5

# 市场发现
minara discover
```

**安装**：
```bash
# 1. 安装 Minara CLI
npm install -g minara@latest

# 2. 登录
minara login

# 3. 安装 Skill
cd ~/.openclaw/skills/
git clone https://github.com/L-LesterYu/OpenClaw-hot-skills-zh.git temp-repo
cp -r temp-repo/skills/minara-zh ./
rm -rf temp-repo
```

**使用场景**：
- "帮我把 0.1 ETH 兑换为 USDC"
- "查询我的钱包余额"
- "买入 100 USDC 的 SOL"
- "热门代币有哪些？"
- "开多 ETH，5 倍杠杆"

**安全特性**：
- 所有资金操作需用户明确确认
- 绝不自动确认交易
- 支持模拟交易（dry-run）

**注意事项**：
- 需要安装 Minara CLI（`npm install -g minara@latest`）
- 需要登录 Minara 账户或设置 `MINARA_API_KEY` 环境变量
- 交易操作会消耗 Gas 费用

[查看详细文档](./skills/minara-zh/SKILL.md) | [安装指南](./skills/minara-zh/setup.md) | [使用示例](./skills/minara-zh/examples.md)

---

### 56. memory-setup-zh - 记忆设置

**版本**: v1.0.0
**状态**: ✅ 已发布
**来源**: ClawHub

**功能**：
- 🧠 **记忆配置**：启用并配置 Moltbot/Clawdbot 的持久化记忆搜索功能
- 📝 **目录结构**：自动创建 MEMORY.md、每日日志和向量搜索的目录结构
- 🔍 **搜索配置**：支持 Voyage、OpenAI、本地嵌入向量三种提供商
- 🛠️ **故障排除**：完整的诊断和修复指南

**核心配置**：
```json
{
  "memorySearch": {
    "enabled": true,
    "provider": "voyage",
    "sources": ["memory", "sessions"],
    "indexMode": "hot",
    "minScore": 0.3,
    "maxResults": 20
  }
}
```

**安装**：
```bash
cd ~/.openclaw/skills/
git clone https://github.com/L-LesterYu/OpenClaw-hot-skills-zh.git temp-repo
cp -r temp-repo/skills/memory-setup-zh ./
rm -rf temp-repo
```

**使用场景**：
- "帮我开启记忆功能"
- "Agent 每次都忘记，怎么修？"
- "配置 memorySearch"
- "从金鱼变成大象 🐘"

**注意事项**：
- 需要在配置中启用 `memorySearch.enabled: true`
- Voyage 提供商需要设置 `VOYAGE_API_KEY` 环境变量
- 配置更改后需重启网关

[查看详细文档](./skills/memory-setup-zh/SKILL.md)

---

### 57. video-frames-zh - 视频帧提取工具

**版本**: v1.0.0
**状态**: ✅ 已发布
**来源**: ClawHub

**功能**：
- 🎞️ **视频帧提取**：使用 ffmpeg 从视频中提取单帧画面
- ⏱️ **时间点定位**：通过指定时间点精确提取画面
- 🔢 **帧索引提取**：通过帧序号提取指定帧
- 📐 **多格式输出**：支持 JPG（快速分享）和 PNG（高清画面）格式
- 🖼️ **缩略图生成**：快速生成视频预览缩略图

**核心脚本**：
- `scripts/frame.sh` - 帧提取脚本，支持 `--time`、`--index`、`--out` 参数

**安装**：
```bash
cd ~/.openclaw/skills/
git clone https://github.com/L-LesterYu/OpenClaw-hot-skills-zh.git temp-repo
cp -r temp-repo/skills/video-frames-zh ./
rm -rf temp-repo
```

**使用场景**：
- "提取这个视频的第一帧"
- "看看视频第 10 秒的画面"
- "给我生成一个视频缩略图"
- "截取视频第 5 帧的画面"

**注意事项**：
- 需要安装 ffmpeg（`brew install ffmpeg`）
- 脚本支持 `--time HH:MM:SS` 和 `--index N` 两种定位方式
- 使用 `.jpg` 格式便于快速分享；使用 `.png` 格式可获得更清晰的画面

[查看详细文档](./skills/video-frames-zh/SKILL.md)

---

### 58. xiaohongshu-mcp-zh - 小红书内容运营自动化

**版本**: v1.0.0
**状态**: ✅ 已发布
**来源**: [GitHub - xpzouying/xiaohongshu-mcp](https://github.com/xpzouying/xiaohongshu-mcp)（8.4k+ ⭐）

**功能**：
- 📝 **内容发布**：发布图文、视频笔记到小红书
- 🔍 **笔记搜索**：搜索笔记内容和热门趋势
- 📊 **详情分析**：获取笔记详情、评论和互动数据
- 👤 **用户管理**：查看用户资料和内容动态
- 📰 **信息流**：获取个性化推荐信息流

**核心命令**：
```bash
# 检查登录状态
python scripts/xhs_client.py status

# 搜索笔记
python scripts/xhs_client.py search "咖啡"

# 获取笔记详情
python scripts/xhs_client.py detail "note_id" "xsec_token"

# 获取推荐信息流
python scripts/xhs_client.py feeds

# 发布笔记（图片URL用逗号分隔）
python scripts/xhs_client.py publish "标题" "内容" "url1,url2"
```

**安装**：
```bash
cd ~/.openclaw/skills/
git clone https://github.com/L-LesterYu/OpenClaw-hot-skills-zh.git temp-repo
cp -r temp-repo/skills/xiaohongshu-mcp-zh ./
rm -rf temp-repo
```

**前置要求**：
1. 从 [GitHub Releases](https://github.com/xpzouying/xiaohongshu-mcp/releases) 下载对应平台的二进制文件
2. 运行登录工具，使用小红书 App 扫码登录
3. 启动 MCP 服务器：`./xiaohongshu-mcp-darwin-arm64`（服务器默认运行在 `http://localhost:18060`）

**使用场景**：
- "搜一下咖啡相关的热门笔记"
- "查看这个笔记的评论"
- "帮我发布一篇图文笔记"
- "分析小红书上户外电源的市场趋势"

**注意事项**：
- 需要本地运行 `xiaohongshu-mcp` 服务器
- 登录后请勿在其他浏览器登录同账号，否则会话失效
- 服务器默认运行在 `http://localhost:18060`
- 支持无头模式和可见浏览器模式

[查看详细文档](./skills/xiaohongshu-mcp-zh/SKILL.md) | [GitHub 源码](https://github.com/xpzouying/xiaohongshu-mcp)

---

### 59. ai-persona-os-zh - AI 人格操作系统

**版本**: v1.0.0
**状态**: ✅ 已发布
**来源**: [ClawHub - AI Persona OS](https://clawhub.ai)

**功能**：
- 🎭 **人格设置向导**：首次运行向导，快速配置 AI 助手的个性化人格
- 💬 **聊天指令系统**：统一的命令入口，根据意图自动路由到对应模块
- 🛡️ **上下文保护**：检查点和心跳机制，确保上下文不丢失
- 📊 **健康监控仪表盘**：实时查看系统运行状态和环境健康状况
- 🧠 **学习记录系统**：自动记录学习内容和错误，支持后续回顾

**核心命令**：
```bash
/persona setup       # 首次运行向导，配置 AI 助手人格
/persona status      # 查看系统健康仪表盘
/persona help        # 列出所有可用命令
/persona checkpoint  # 立即保存当前上下文
/persona heartbeat   # 运行环境健康检查
/persona learn "记住在文件操作前检查权限"  # 记录学习或错误
```

**工作原理**：
1. 读取意图映射配置
2. 将用户输入与关键词模式匹配
3. 路由到正确的技能模块
4. 无匹配时显示回退帮助信息

**安装**：
```bash
cd ~/.openclaw/skills/
git clone https://github.com/L-LesterYu/OpenClaw-hot-skills-zh.git temp-repo
cp -r temp-repo/skills/ai-persona-os-zh ./
rm -rf temp-repo
```

**使用场景**：
- 为 AI 助手配置独特的个性化人格
- 查看和管理 AI 助手的运行状态
- 保护上下文不因会话中断而丢失
- 记录和回顾学习经验

**注意事项**：
- 首次使用建议运行 `persona setup` 进行初始化配置
- 定期运行 `persona checkpoint` 保护重要上下文
- 学习记录可用于持续改进 AI 助手的行为

[查看详细文档](./skills/ai-persona-os-zh/SKILL.md)

---

### 60. word-docx-zh - Word / DOCX 文档处理

**版本**: v1.0.2
**状态**: ✅ 已发布
**来源**: [ClawIC - Word/DOCX](https://clawic.com/skills/word-docx)

**功能**：
- 📘 **文档创建与编辑**：创建、检查和编辑 Microsoft Word 文档和 DOCX 文件
- 🎨 **样式管理**：可靠的样式、编号和格式控制
- 📝 **修订跟踪**：支持修订跟踪、批注和域操作
- 📊 **表格处理**：处理复杂表格和嵌套结构
- 📄 **节与布局**：节操作、页眉页脚和兼容性检查
- 🔄 **格式保持**：文档在往返编辑后保持格式不漂移

**核心规则**：
- 将 DOCX 视为 OOXML（XML 部件的 ZIP 包），而非纯文本
- 优先使用命名样式而非直接格式，保持文档可编辑性
- 文本可能被拆分到多个运行节点中，不能假设单一 XML 节点
- 应有目的地使用不同工作流：结构化提取、样式驱动生成、OOXML 感知编辑

**适用场景**：
- 任务涉及 Word 或 `.docx` 文件
- 文件包含修订跟踪、批注、域、表格、模板或页面布局约束
- 文档需要在往返编辑后保持格式不漂移

**安装**：
```bash
cd ~/.openclaw/skills/
git clone https://github.com/L-LesterYu/OpenClaw-hot-skills-zh.git temp-repo
cp -r temp-repo/skills/word-docx-zh ./
rm -rf temp-repo
```

**使用场景**：
- "创建一个包含表格的 Word 文档"
- "编辑这份报告的样式和格式"
- "检查文档的修订跟踪状态"
- "将 Markdown 转换为 Word 格式"
- "修复 DOCX 文件的格式漂移问题"

**注意事项**：
- 旧版 `.doc` 输入通常需要转换后才能信任现代 `.docx` 的假设
- 进行深度编辑时，应检查包布局而不是仅依赖渲染输出
- 读取、生成和保留已审阅文档是不同的工作

[查看详细文档](./skills/word-docx-zh/SKILL.md)

---

### 61. moltguard-zh - MoltGuard 安全防护插件

**版本**: v6.0.2
**状态**: ✅ 已发布
**来源**: [GitHub - openguardrails/moltguard](https://github.com/openguardrails/moltguard)

**功能**：
- 🛡️ **提示词净化**：本地检测并清除提示词中的敏感数据（邮箱、密钥、银行卡号等）
- 🔍 **注入检测**：通过 MoltGuard API 检测隐藏在长内容中的提示注入攻击
- 🔒 **敏感数据保护**：自动识别和替换高熵字符串、API 密钥、凭据等
- 🔄 **数据还原**：API 响应中自动还原被替换的敏感数据
- 📊 **分析记录**：本地存储分析结果和日志

**核心能力**：
- 网关模式：自动净化发送给 LLM 的内容
- Agent 模式：通过 API 进行深度注入检测
- 支持多种 LLM 提供商（OpenAI、Anthropic、Gemini）
- 内存存储分析结果，支持后续查询

**安装**：
```bash
# 1. 安装插件
npm install @openguardrails/moltguard

# 2. 安装 Skill
cd ~/.openclaw/skills/
git clone https://github.com/L-LesterYu/OpenClaw-hot-skills-zh.git temp-repo
cp -r temp-repo/skills/moltguard-zh ./
rm -rf temp-repo
```

**使用场景**：
- 防止用户发送的邮件、网页等内容中隐藏的提示注入
- 保护 API 密钥和其他敏感数据不被发送到 LLM
- 对接 OpenClaw 网关实现自动安全防护

**注意事项**：
- 需要从 MoltGuard 获取 API 密钥
- 分析功能会向 `api.moltguard.com` 发送经过净化的内容
- 安装前建议审核源码，确保安全性

[查看详细文档](./skills/moltguard-zh/clawhub-skill/SKILL.md) | [GitHub](https://github.com/openguardrails/moltguard)

---

### 62. news-summary-zh - 新闻摘要

**版本**: v1.0.0
**状态**: ✅ 已发布

**功能**：
- 📰 **RSS 新闻获取**：从 BBC、路透社、NPR、半岛电视台等权威新闻源获取新闻
- 🌍 **多源整合**：整合西方和全球南方视角，平衡报道
- 📝 **智能摘要**：按地区或主题分类，生成简洁的 5-8 条重点新闻
- 🎙️ **语音播报**：支持通过 OpenAI TTS 生成语音摘要
- 📊 **分类展示**：国际、财经、科技等分类展示

**支持的新闻源**：
- **BBC**：世界新闻、头条、财经、科技
- **路透社**：全球视角
- **NPR**：美国视角
- **半岛电视台**：全球南方视角

**安装**：
```bash
cd ~/.openclaw/skills/
git clone https://github.com/L-LesterYu/OpenClaw-hot-skills-zh.git temp-repo
cp -r temp-repo/skills/news-summary-zh ./
rm -rf temp-repo
```

**快速开始**：
```bash
# 获取 BBC 世界新闻
curl -s "https://feeds.bbci.co.uk/news/world/rss.xml"

# 解析 RSS 标题和描述
curl -s "https://feeds.bbci.co.uk/news/world/rss.xml" | \
grep -E "<title>|<description>" | \
sed 's/<[^>]*>//g' | head -30
```

**使用场景**：
- "今天有什么新闻？"
- "帮我做个每日新闻简报"
- "最近有什么重大事件？"
- "用语音播报今天的新闻"

**输出格式**：
```
📰 新闻摘要 [日期]

🌍 国际
- [头条 1]
- [头条 2]

💼 财经
- [头条 1]

💻 科技
- [头条 1]
```

[查看详细文档](./skills/news-summary-zh/SKILL.md)

---

### 63. qmd-zh - 本地搜索/索引工具

**版本**: v1.0.0
**状态**: ✅ 已发布
**来源**: [GitHub - tobi/qmd](https://github.com/tobi/qmd)
**首页**: https://tobi.lutke.com

**功能**：
- 📝 **本地文件索引**：对本地 Markdown 等文件建立全文索引
- 🔍 **BM25 关键词搜索**：经典的关键词匹配搜索
- 🧠 **向量语义搜索**：基于 Ollama 的语义向量搜索
- 🔄 **混合搜索**：结合 BM25 + 向量 + 重排序的混合搜索（推荐）
- 📄 **文档获取**：获取指定文件的内容片段
- 🔌 **MCP 模式**：支持以 MCP 服务器方式运行

**核心命令**：
```bash
# 添加集合
qmd collection add /path --name docs --mask "**/*.md"

# 更新索引
qmd update

# 查看状态
qmd status

# BM25 搜索
qmd search "查询内容"

# 向量语义搜索
qmd vsearch "查询内容"

# 混合搜索（推荐）
qmd query "查询内容"

# 获取文档内容
qmd get docs/path.md:10 -l 40

# MCP 模式
qmd mcp
```

**安装**：
```bash
# 1. 安装 qmd CLI
npm install -g https://github.com/tobi/qmd

# 2. 安装 Skill
cd ~/.openclaw/skills/
git clone https://github.com/L-LesterYu/OpenClaw-hot-skills-zh.git temp-repo
cp -r temp-repo/skills/qmd-zh ./
rm -rf temp-repo
```

**使用场景**：
- "在本地文档中搜索关于 X 的内容"
- "索引我的笔记目录"
- "用混合搜索查找相关文档"
- "获取某个文档的指定片段"

**注意事项**：
- 向量搜索和重排序功能使用 Ollama，需设置 `OLLAMA_URL`（默认 `http://localhost:11434`）
- 索引数据默认存储在 `~/.cache/qmd` 目录下
- 需要安装 Node.js 和 Ollama（用于向量搜索）

[查看详细文档](./skills/qmd-zh/SKILL.md) | [GitHub](https://github.com/tobi/qmd)

---

## 🚀 快速开始

### 前置要求

- OpenClaw 已安装并配置
- 基本的命令行操作能力

### 安装步骤

1. **克隆仓库**
```bash
git clone https://github.com/L-LesterYu/OpenClaw-hot-skills-zh.git
```

2. **选择并安装 Skill**
```bash
cd OpenClaw-hot-skills-zh/
cp -r <skill-name> ~/.openclaw/skills/
```

3. **验证安装**
```bash
ls ~/.openclaw/skills/<skill-name>/
```

## 🤝 贡献指南

我们欢迎所有形式的贡献！

### 如何贡献

1. **Fork 本仓库**
2. **创建特性分支** (`git checkout -b feature/AmazingSkill`)
3. **提交更改** (`git commit -m '添加某个中文版 Skill'`)
4. **推送到分支** (`git push origin feature/AmazingSkill`)
5. **创建 Pull Request`

### 贡献规范

- ✅ 确保翻译准确、流畅
- ✅ 保持代码示例和技术术语不变
- ✅ 添加详细的使用说明和示例
- ✅ 遵循现有的文件结构
- ✅ 更新相关的 README 文档

## 📝 待办事项

- [ ] 添加更多 ClawHub 优质 Skills 的中文版本
- [ ] 创建自动化测试脚本
- [ ] 建立持续集成/持续部署流程
- [ ] 编写更详细的贡献指南
- [ ] 创建 Skill 质量检查工具

## 🔗 相关链接

- [OpenClaw 官网](https://openclaw.ai)
- [ClawHub Skills 市场](https://clawhub.ai)
- [OpenClaw 文档](https://docs.openclaw.ai)
- [Agent Skills 规范](https://agentskills.io/specification)

## 📄 许可证

本项目采用 MIT 许可证 - 详见 [LICENSE](LICENSE) 文件

## 🙏 致谢

感谢以下项目和个人的支持：

- [OpenClaw](https://github.com/openclaw/openclaw) - 强大的 AI 代理框架
- [ClawHub](https://clawhub.ai) - Skills 分享平台
- 所有原始 Skills 的作者们
- 为本项目贡献的所有开发者

## 📮 联系方式

- **项目地址**: https://github.com/L-LesterYu/OpenClaw-hot-skills-zh
- **问题反馈**: [GitHub Issues](https://github.com/L-LesterYu/OpenClaw-hot-skills-zh/issues)

---

**如果这个项目对您有帮助，请给一个 ⭐️ Star 支持一下！**

Made with ❤️ by [L-LesterYu](https://github.com/L-LesterYu)
