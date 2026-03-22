# OpenClaw-hot-skills-zh 🇨🇳

> 致力于将 ClawHub 平台的优质安全项目转化为中文版 Skills，降低阅读理解门槛，提供更符合中文开发者习惯的开源技能库。

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![GitHub stars](https://img.shields.io/github/stars/L-LesterYu/OpenClaw-hot-skills-zh.svg?style=social)](https://github.com/L-LesterYu/OpenClaw-hot-skills-zh/stargazers)

## 📖 项目简介

本项目为中文开发者提供 ClawHub 平台上优质技能（Skills）的中文版本，共收录 **69 个**中文技能，覆盖搜索、Agent 增强、工具集成、设计、开发、文档处理、数据分析和安全防护等领域。

**每个 skill 的完整文档在其子目录的 `SKILL.md` 中，本 README 仅提供速览。**

### 图例

| 标记 | 含义 |
|------|------|
| 💰 | 需要付费 API Key |
| 📦 | 需要额外安装 CLI 工具 |
| ✅ | 免费且零配置，开箱即用 |

---

## ⭐ 推荐技能（新手必装）

| 技能 | 推荐理由 |
|------|---------|
| [multi-search-engine-zh](skills/multi-search-engine-zh/SKILL.md) | **17 个搜索引擎一键切换**，国内外全覆盖，无需 API Key，搜索能力最强 |
| [find-skills-zh](skills/find-skills-zh/SKILL.md) | 技能生态探索入口——关键词搜索 skills.sh 技能库，发现并安装所需能力 |
| [tavily-search-zh](skills/tavily-search-zh/SKILL.md) | Tavily API 网页搜索，免费 API Key，精准返回少量结果与 AI 摘要，Brave 的零成本替代方案 |
| [proactive-self-improving-agent-zh](skills/proactive-self-improving-agent-zh/SKILL.md) | **让 Agent 越用越聪明**——自动捕获经验、纠正错误、积累知识 |
| [frontend-design-zh](skills/frontend-design-zh/SKILL.md) | 纯文本描述即可生成生产级 React 网站，设计师的最强辅助 |
| [word-docx-zh](skills/word-docx-zh/SKILL.md) | 一句话创建/编辑 Word 文档，支持样式、表格、修订，日常办公必备 |
| [markdown-converter-zh](skills/markdown-converter-zh/SKILL.md) | PDF/Word/PPT/Excel → Markdown 一键转换，喂给 AI 的万能适配器 |
| [skill-vetter-zh](skills/skill-vetter-zh/SKILL.md) | 安装任何技能前先过一遍安全扫描，保护你的系统 |

---

## 📦 分类速览

### 🔍 搜索与信息获取

| 技能 | 描述 | 标记 |
|------|------|------|
| [multi-search-engine-zh](skills/multi-search-engine-zh/SKILL.md) | 零配置集成 17 个搜索引擎（百度/Google/DuckDuckGo 等），支持高级操作符、时间筛选、站内搜索、隐私搜索及 WolframAlpha 知识计算，DuckDuckGo Bangs 快捷跳转 | ✅ |
| [tavily-search-zh](skills/tavily-search-zh/SKILL.md) | Tavily API 网页搜索，免费注册即用，返回精准结果与可选 AI 摘要，Brave 的零成本替代方案 | 💰📦 |
| [baidu-search-zh](skills/baidu-search-zh/SKILL.md) | 百度 AI 搜索引擎，适合中文内容查询 | 💰📦 |
| [brave-search-zh](skills/brave-search-zh/SKILL.md) | Brave Search API 轻量搜索 CLI——网页搜索 + 可选页面内容提取（Markdown），无需浏览器 | 💰📦 |
| [brave-web-search-zh](skills/brave-web-search-zh/SKILL.md) | Brave Search API 双模式搜索：排序网页结果 + AI 摘要答案，支持时间筛选 | 💰📦 |
| [answer-overflow-zh](skills/answer-overflow-zh/SKILL.md) | 搜索已索引的 Discord 社区讨论，查找编程问题解决方案 | ✅ |
| [news-summary-zh](skills/news-summary-zh/SKILL.md) | 新闻摘要聚合，快速了解最新资讯 | ✅ |
| [youtube-zh](skills/youtube-zh/SKILL.md) | YouTube 视频/频道搜索、详情获取和字幕提取 | 💰📦 |
| [youtube-watcher-zh](skills/youtube-watcher-zh/SKILL.md) | 通过 yt-dlp 提取 YouTube 视频字幕，支持视频摘要、问答和内容提取，仅适用于有字幕的视频 | 📦 |
| [web-search-plus-zh](skills/web-search-plus-zh/SKILL.md) | 智能自动路由的统一搜索——连接 7 个搜索提供商（Serper/Tavily/Querit/Exa/Perplexity/You.com/SearXNG），通过多信号分析自动选择最佳引擎，附带置信度评分 | 💰📦 |

### 🧠 Agent 增强

| 技能 | 描述 | 标记 |
|------|------|------|
| [proactive-self-improving-agent-zh](skills/proactive-self-improving-agent-zh/SKILL.md) | 自动捕获经验、纠正错误、积累学习，经验反复出现时自动晋升为永久规则 | ✅ |
| [proactive-agent-zh](skills/proactive-agent-zh/SKILL.md) | 主动式 Agent 架构——WAL 预写日志、工作缓冲区、压缩恢复、自主定时任务、反向提示与安全演进护栏，让 Agent 预测需求并在上下文丢失中存活 | ✅ |
| [self-improving-zh](skills/self-improving-zh/SKILL.md) | 自我反思与持续改进引擎：自动捕获纠正和偏好信号，分层记忆存储（热/温/冷），模式出现3次后自动晋升为永久规则 | ✅ |
| [self-improvement-zh](skills/self-improvement-zh/SKILL.md) | 捕获学习内容、错误和纠正以实现持续改进 | ✅ |
| [self-improving-proactive-agent-zh](skills/self-improving-proactive-agent-zh/SKILL.md) | 自我改进 + 主动式 Agent 的组合版本 | ✅ |
| [elite-longterm-memory-zh](skills/elite-longterm-memory-zh/SKILL.md) | 终极记忆系统——WAL 协议 + 向量搜索 + Git Notes + 云备份 | 📦 |
| [memory-setup-zh](skills/memory-setup-zh/SKILL.md) | 启用并配置持久化记忆搜索功能，解决"金鱼记忆"问题 | ✅ |
| [ai-persona-os-zh](skills/ai-persona-os-zh/SKILL.md) | AI 人格操作系统，统一管理人格设置、聊天指令和监控 | ✅ |
| [capability-evolver-zh](skills/capability-evolver-zh/SKILL.md) | 面向 Agent 的自进化引擎，分析运行时历史发现改进机会 | 📦 |
| [freeride-zh](skills/freeride-zh/SKILL.md) | 按 OpenRouter 免费模型质量自动排名并配置主模型与备用模型，限流时自动故障切换，零成本使用 AI | 💰📦 |
| [auto-updater-zh](skills/auto-updater-zh/SKILL.md) | 通过 cron 每日自动更新 Clawdbot 核心与所有已安装技能，检查可用版本、应用更新并向用户发送变更摘要 | 📦 |
| [ontology-zh](skills/ontology-zh/SKILL.md) | 类型化知识图谱——实体/关系建模、Schema 约束验证（属性/枚举/基数/无环）、CLI 增删改查与图变换规划，支持跨技能数据共享 | ✅ |

### 📱 工具与服务集成

| 技能 | 描述 | 标记 |
|------|------|------|
| [gog-zh](skills/gog-zh/SKILL.md) | Google Workspace 统一 CLI 工具，覆盖 Gmail 收发、日历事件查询、云盘搜索、联系人管理、表格读写与文档导出，需 OAuth 认证 | 📦 |
| [notion-skill-zh](skills/notion-skill-zh/SKILL.md) | 通过官方 Notion API 读写页面、管理数据库、追加内容块，支持多配置文件和架构变更 | 💰📦 |
| [slack-zh](skills/slack-zh/SKILL.md) | 控制 Slack 频道消息发送、置顶等操作 | ✅ |
| [trello-zh](skills/trello-zh/SKILL.md) | 通过 REST API 管理 Trello 看板、列表和卡片 | 💰📦 |
| [caldav-calendar-zh](skills/caldav-calendar-zh/SKILL.md) | CalDAV 日历本地同步与事件管理——基于 vdirsyncer + khal 实现 iCloud/Google/Fastmail/Nextcloud 等日历的双向同步、事件 CRUD、搜索与格式化输出，数据完全本地存储 | 📦 |
| [obsidian-zh](skills/obsidian-zh/SKILL.md) | 基于 obsidian-cli 的 Obsidian 笔记库管理——日常笔记创建、前缀自动分区（任务/想法/问题/洞察）、集中式任务日志与艾森豪威尔优先级标签 | 📦 |
| [weather-zh](skills/weather-zh/SKILL.md) | 通过 wttr.in + Open-Meteo 免费查询天气——紧凑/完整预报、中文输出、JSON API 程序化调用，仅需 curl | ✅ |
| [minara-zh](skills/minara-zh/SKILL.md) | 加密货币交易 CLI——兑换、永续合约、支付，支持 EVM + Solana | 💰📦 |

### 🎨 设计与内容创作

| 技能 | 描述 | 标记 |
|------|------|------|
| [frontend-design-zh](skills/frontend-design-zh/SKILL.md) | 从纯文本需求生成生产级网站，无需设计稿。React 18 + TypeScript + Tailwind + shadcn/ui + Framer Motion，反 AI 套路美学，支持 Vite 静态和 Next.js 双工作流 | 📦 |
| [SuperDesign-zh](skills/SuperDesign-zh/SKILL.md) | 前端设计指南，专精创建美观现代的 UI 界面 | ✅ |
| [ui-ux-pro-max-zh](skills/ui-ux-pro-max-zh/README.md) | AI 设计智能工具包——UI 风格、配色、字体、UX 指南数据库 | 📦 |
| [nano-banana-pro-zh](skills/nano-banana-pro-zh/SKILL.md) | 通过 Gemini 3 Pro Image API 生成与编辑图像——支持文生图、图生图及 1K/2K/4K 分辨率，内置草稿→迭代→最终渐进工作流 | 💰📦 |
| [content-ideas-generator-zh](skills/content-ideas-generator-zh/SKILL.md) | 从新闻通讯/笔记/脚本等素材中提炼反直觉洞察，生成 5 个高互动社交媒体帖子大纲——包含悖论框架、转变弧线、异议预判与语言技巧应用，仅输出大纲不写全文 | ✅ |
| [marketing-mode-zh](skills/marketing-mode-zh/SKILL.md) | 23 个营销技能整合——战略、心理学、SEO、转化优化、付费增长 | ✅ |
| [marketing-skills-zh](skills/marketing-skills-zh/SKILL.md) | 23 个营销实战手册——CRO、SEO、文案、分析、实验、定价等 | ✅ |
| [admapix-zh](skills/admapix-zh/SKILL.md) | 广告创意素材搜索，查找竞品广告视频和创意素材 | 💰 |
| [xiaohongshu-mcp-zh](skills/xiaohongshu-mcp-zh/SKILL.md) | 小红书内容运营自动化——发布、搜索、趋势分析 | 📦 |

### 🛠️ 开发工具

| 技能 | 描述 | 标记 |
|------|------|------|
| [github-zh](skills/github-zh/SKILL.md) | 通过 `gh` CLI 管理 GitHub 全流程——Issues、PR 生命周期、CI/CD 监控与重跑、仓库操作，支持 `gh api` 高级查询与 JSON 结构化输出 | 📦 |
| [mcporter-zh](skills/mcporter-zh/SKILL.md) | MCP 服务器管理——列出、配置、认证和调用 MCP 工具 | 📦 |
| [playwright-mcp-zh](skills/playwright-mcp-zh/SKILL.md) | 通过 Playwright MCP 实现浏览器自动化 | 📦 |
| [Agent-Browser-zh](skills/Agent-Browser-zh/SKILL.md) | AI 专用无头浏览器自动化 CLI——基于无障碍树确定性引用驱动导航、点击、填充，Rust+Node.js 架构，支持会话管理与 JSON 结构化输出 | 📦 |
| [desktop-control-zh](skills/desktop-control-zh/SKILL.md) | 基于 PyAutoGUI 的桌面自动化——像素级鼠标/键盘控制、屏幕截图与图像识别、窗口管理、剪贴板操作，内置故障保护机制 | 📦 |
| [skill-creator-zh](skills/skill-creator-zh/SKILL.md) | AgentSkills 全生命周期工具——创建、编辑、验证、打包，内置初始化脚本和规范检查器 | ✅ |
| [clawhub-zh](skills/clawhub-zh/SKILL.md) | 通过 ClawHub CLI 管理 clawhub.com 上的智能体技能——搜索发现、安装指定版本、哈希匹配增量更新（支持全量）、列出已装技能及发布新技能，支持自定义注册表与工作目录 | 📦 |
| [find-skills-zh](skills/find-skills-zh/SKILL.md) | 通过 `npx skills` CLI 搜索 skills.sh 开放技能库，按关键词发现并安装智能体扩展能力 | 📦 |
| [n8n-workflow-automation-zh](skills/n8n-workflow-automation-zh/SKILL.md) | 设计健壮的 n8n 工作流 JSON，含重试、日志和审核队列 | ✅ |
| [automation-workflows-zh](skills/automation-workflows-zh/SKILL.md) | 自动化工作流设计——识别重复任务、选择工具、构建流程 | ✅ |
| [clawddocs-zh](skills/clawddocs-zh/SKILL.md) | Clawdbot 文档智能助手——决策树导航、多模式搜索（站点地图/关键词/全文索引）、文档获取与版本变更跟踪，内置各功能配置片段 | ✅ |
| [byterover-zh](skills/byterover-zh/SKILL.md) | ByteRover CLI 上下文工程平台——管理项目知识和决策 | 💰📦 |
| [qmd-zh](skills/qmd-zh/SKILL.md) | 本地搜索/索引 CLI（BM25 + 向量搜索 + 重排序），附带 MCP 模式 | 📦 |
| [qa-engineer-zh](skills/qa-engineer-zh/SKILL.md) | 软件测试与质量保证——发现 bug、执行测试、生成报告 | ✅ |
| [model-usage-zh](skills/model-usage-zh/SKILL.md) | 汇总 Codex/Claude 的按模型使用情况和成本数据 | 📦 |
| [api-gateway-zh](skills/api-gateway-zh/SKILL.md) | 通过 Maton.ai 代理统一访问 100+ 第三方 API（Slack、HubSpot、Google Workspace、Notion、Stripe 等），托管 OAuth 自动注入令牌，支持原生端点透传与多连接管理 | 💰 |

### 📄 文档与文件处理

| 技能 | 描述 | 标记 |
|------|------|------|
| [word-docx-zh](skills/word-docx-zh/SKILL.md) | 创建/编辑 Word 文档，支持样式、表格、修订跟踪 | ✅ |
| [gongwen-zh](skills/gongwen-zh/SKILL.md) | 按党政机关公文格式标准自动排版 | ✅ |
| [nano-pdf-zh](skills/nano-pdf-zh/SKILL.md) | 用自然语言指令编辑 PDF 文件 | 📦 |
| [excel-xlsx-zh](skills/excel-xlsx-zh/SKILL.md) | 创建、检查和编辑 Excel 工作簿，支持公式、日期、格式化和模板保留 | 📦 |
| [markdown-converter-zh](skills/markdown-converter-zh/SKILL.md) | PDF/Word/PPT/Excel/HTML 等转 Markdown，支持 OCR 和音频转录 | 📦 |
| [video-frames-zh](skills/video-frames-zh/SKILL.md) | 使用 ffmpeg 从视频中提取帧或短片段 | 📦 |
| [openai-whisper-zh](skills/openai-whisper-zh/SKILL.md) | 本地运行 OpenAI Whisper 模型进行语音转录与翻译，支持多音频格式、多模型尺寸选择，数据不出本地 | 📦 |
| [humanizer-zh](skills/humanizer-zh/SKILL.md) | 检测并移除 AI 写作痕迹，基于 Wikipedia 指南修复 24 种 AI 文本模式（过度夸大、促销语言、模糊归因等），让文本更自然 | ✅ |

### 📊 数据分析

| 技能 | 描述 | 标记 |
|------|------|------|
| [stock-analysis-zh](skills/stock-analysis-zh/SKILL.md) | 基于 Yahoo Finance API 的全面股票研究——公司概览、技术指标与估值、多股对比、内部人士交易追踪及 SEC 监管文件，支持美股与国际市场 | 💰付费API |
| [stock-market-pro-zh](skills/stock-market-pro-zh/SKILL.md) | 专业股市分析工具 | 📦 |

### 🔒 安全与防护

| 技能 | 描述 | 标记 |
|------|------|------|
| [moltguard-zh](skills/moltguard-zh/openclaw.plugin.json) | 本地提示词净化与提示注入检测，保护 Agent 安全 | ✅ |
| [skill-vetter-zh](skills/skill-vetter-zh/SKILL.md) | 技能安装前安全审查协议——来源验证、14项危险信号检测、权限范围评估与四级风险分类，输出结构化审查报告 | ✅ |
| [edgeone-clawscan-zh](skills/edgeone-clawscan-zh/SKILL.md) | 由腾讯朱雀实验室 AI-Infra-Guard 驱动的全面安全扫描，支持 OpenClaw 环境体检和技能安全审查 | 📦 |

---

## 🤔 选择指南

面对功能相似的技能时，可以根据以下建议选择：

### 🔍 搜索引擎怎么选？

| 场景 | 推荐 | 理由 |
|------|------|------|
| **全能搜索** | `multi-search-engine-zh` | 17 个引擎全覆盖，无需 API Key |
| **中文内容** | `baidu-search-zh` | 百度生态，中文结果更精准 |
| **精准搜索** | `tavily-search-zh` | 免费 API Key，精准结果 + AI 摘要，Brave 替代 |
| **Brave 用户** | `brave-web-search-zh` | 带 AI 摘要，体验更好；`brave-search-zh` 更轻量 |

### 🧠 记忆/自我改进怎么选？

| 场景 | 推荐 | 理由 |
|------|------|------|
| **推荐首选** | `proactive-self-improving-agent-zh` | 功能最完整，自动捕获经验并晋升为永久规则 |
| **主动行为** | `proactive-agent-zh` | 侧重让 Agent 主动预测需求和自主执行任务 |
| **轻量改进** | `self-improving-zh` | 简洁的自我反思和学习循环 |
| **长期记忆** | `elite-longterm-memory-zh` | WAL 协议 + 向量搜索，适合大量上下文持久化 |

> 💡 `self-improvement-zh` 和 `self-improving-zh` `self-improving-proactive-agent-zh`功能高度重叠，选其一即可。`self-improving-proactive-agent-zh` 是改进 + 主动的组合版。

### 📈 股票分析怎么选？

| 场景 | 推荐 | 理由 |
|------|------|------|
| **全面研究** | `stock-analysis-zh` | 公司简介、技术分析、SEC 文件、内部人士数据一应俱全 |
| **本地运行** | `stock-market-pro-zh` | 使用 `uv run` 本地运行，不依赖外部 API |

### 🎬 YouTube 怎么选？

| 场景 | 推荐 | 理由 |
|------|------|------|
| **搜索 + 字幕** | `youtube-zh` | 功能更全，支持频道信息、搜索和字幕提取 |
| **纯字幕读取** | `youtube-watcher-zh` | 更轻量，专注获取字幕文本 |

### ✍️ 文本人性化怎么选？

| 场景 | 推荐 | 理由 |
|------|------|------|
| **首选** | `humanizer-zh` | 基于 Wikipedia 指南，覆盖 24 种 AI 写作模式，修复更全面 |
| **轻量替代** | `humanize-ai-text-zh` | 功能相同，结构更简洁，任选其一即可 |

---

## 🚀 快速开始

### 前置要求

- OpenClaw 已安装并配置
- 基本的命令行操作能力

### 安装方式

**方式一：克隆后选择性安装**
```bash
git clone https://github.com/L-LesterYu/OpenClaw-hot-skills-zh.git
cd OpenClaw-hot-skills-zh/
cp -r skills/<skill-name> ~/.openclaw/skills/
```

**方式二：通过 ClawHub 安装**（安装英文原版，部分技能支持）
```bash
openclaw skill install <skill-name>
```

### 验证安装
```bash
ls ~/.openclaw/skills/<skill-name>/
```

---

## 🤝 贡献指南

欢迎所有形式的贡献！

1. **Fork** 本仓库
2. **创建特性分支** (`git checkout -b feature/AmazingSkill`)
3. **提交更改** (`git commit -m '添加某个中文版 Skill'`)
4. **推送** (`git push origin feature/AmazingSkill`)
5. **创建 Pull Request**

### 贡献规范

- ✅ 确保翻译准确、流畅
- ✅ 保持代码示例和技术术语不变
- ✅ 添加详细的使用说明和示例
- ✅ 遵循现有的文件结构

---

## 🔗 相关链接

- [OpenClaw 官网](https://openclaw.ai)
- [ClawHub Skills 市场](https://clawhub.ai)
- [OpenClaw 文档](https://docs.openclaw.ai)
- [Agent Skills 规范](https://agentskills.io/specification)

---

## 📄 许可证

本项目采用 MIT 许可证 - 详见 [LICENSE](LICENSE) 文件

## 🙏 致谢

- [OpenClaw](https://github.com/openclaw/openclaw) - AI 代理框架
- [ClawHub](https://clawhub.ai) - Skills 分享平台
- 所有原始 Skills 的作者们
- 为本项目贡献的开发者们

## 📮 联系方式

- **项目地址**: https://github.com/L-LesterYu/OpenClaw-hot-skills-zh
- **问题反馈**: [GitHub Issues](https://github.com/L-LesterYu/OpenClaw-hot-skills-zh/issues)

---

**如果这个项目对您有帮助，请给一个 ⭐️ Star 支持一下！**

Made with ❤️ by [L-LesterYu](https://github.com/L-LesterYu)
