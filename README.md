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
