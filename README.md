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

### 29. skill-vetter-zh - 技能审查器

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

### 30. automation-workflows-zh - 自动化工作流程

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

### 31. humanize-ai-text-zh - AI 文本人性化工具

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

### 32. markdown-converter-zh - Markdown 转换器

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
