# OpenClaw-hot-skills-zh 🇨🇳

> 致力于将 各个skills 平台的优质、热门、安全的项目转化为中文版 Skills。降低阅读理解难度和优化门槛，提供更符合中文开发者习惯的开源技能库。

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![GitHub stars](https://img.shields.io/github/stars/L-LesterYu/OpenClaw-hot-skills-zh.svg?style=social)](https://github.com/L-LesterYu/OpenClaw-hot-skills-zh/stargazers)

## 📖 项目简介

本项目为中文开发者提供 各个skills 平台上优质技能（Skills）的中文版本，共收录 **85 个**中文技能，覆盖搜索、Agent 增强、工具集成、设计、开发、文档处理、数据分析和安全防护等领域。

**每个 skill 的完整文档在其子目录的 `SKILL.md` 中，本 README 仅提供速览。**

### 图例

| 标记 | 含义 |
|------|------|
| 💰 | 需要付费 API Key |
| 📦 | 需要额外安装 CLI 工具 |
| ✅ | 免费且零配置，开箱即用 |

---

## ⭐ 推荐技能（新手推荐）

| 技能 | 推荐理由 |
|------|---------|
| [multi-search-engine-zh](skills/multi-search-engine-zh/SKILL.md) | **17 个搜索引擎一键切换**，国内外全覆盖，无需 API Key，搜索能力最强 |
| [find-skills-zh](skills/find-skills-zh/SKILL.md) | 技能生态探索入口——关键词搜索 skills.sh 技能库，发现并安装所需能力 |
| [tavily-search-zh](skills/tavily-search-zh/SKILL.md) | Tavily API 网页搜索，免费 API Key，精准返回少量结果与 AI 摘要，Brave 的零成本替代方案 |
| [proactive-self-improving-agent-zh](skills/proactive-self-improving-agent-zh/SKILL.md) | **让 Agent 越用越聪明**——自动捕获经验、纠正错误、积累知识 |
| [word-docx-zh](skills/word-docx-zh/SKILL.md) | **OOXML 感知的 Word 精准编辑**——样式、编号、修订、批注、域、节完整可控，往返编辑格式不漂移，审阅工作流首选 |
| [markdown-converter-zh](skills/markdown-converter-zh/SKILL.md) | 15+ 格式一键转 Markdown——PDF/Office/图片(OCR)/音频转录/YouTube 字幕，喂给 AI 的万能适配器 |
| [nano-banana-pro-prompts-recommend-skill-zh](skills/nano-banana-pro-prompts-recommend-skill-zh/SKILL.md) | **10,000+ 经过验证的图像生成提示词库**——附带示例图片，适用于 Nano Banana Pro、Gemini、GPT Image、Midjourney、DALL-E 等多种AI模型 |
| [skill-vetter-zh](skills/skill-vetter-zh/SKILL.md) | 安装任何技能前先过一遍安全扫描，保护你的系统 |

---

## 📦 分类速览

### 🔍 搜索与信息获取

| 技能 | 描述 | 标记 |
|------|------|------|
| [multi-search-engine-zh](skills/multi-search-engine-zh/SKILL.md) | 零配置集成 17 个搜索引擎（百度/Google/DuckDuckGo 等），支持高级操作符、时间筛选、站内搜索、隐私搜索及 WolframAlpha 知识计算，DuckDuckGo Bangs 快捷跳转 | ✅ |
| [tavily-search-zh](skills/tavily-search-zh/SKILL.md) | Tavily API 网页搜索，免费注册即用，返回精准结果与可选 AI 摘要，Brave 的零成本替代方案 | 💰📦 |
| [baidu-search-zh](skills/baidu-search-zh/SKILL.md) | 百度 AI 搜索引擎 API——网页搜索与时间筛选（24h/7天/月/年/自定义日期区间），结果数量可控，中文搜索场景首选 | 💰📦 |
| [brave-search-zh](skills/brave-search-zh/SKILL.md) | Brave Search API 轻量搜索 CLI——网页搜索 + 可选页面内容提取（Markdown），无需浏览器 | 💰📦 |
| [brave-web-search-zh](skills/brave-web-search-zh/SKILL.md) | Brave Search API 双模式搜索：排序网页结果 + AI 摘要答案，支持时间筛选 | 💰📦 |
| [answer-overflow-zh](skills/answer-overflow-zh/SKILL.md) | 搜索已索引的 Discord 社区讨论，获取编程问题、库问题和社区问答的解决方案，支持 MCP 服务器搜索与线程内容 Markdown 抓取 | ✅ |
| [news-summary-zh](skills/news-summary-zh/SKILL.md) | 通过 RSS 订阅 BBC/路透社/NPR/半岛电视台等国际新闻源，按地区与主题分类输出摘要，支持文字与语音（OpenAI TTS）双模式 | ✅ |
| [youtube-zh](skills/youtube-zh/SKILL.md) | 基于 YouTube Data API v3 的视频研究工具——通过 MCP 服务器实现视频搜索、频道信息查询、视频详情获取与带时间戳字幕提取，yt-dlp 作为字幕提取备用方案 | 💰📦 |
| [youtube-watcher-zh](skills/youtube-watcher-zh/SKILL.md) | 通过 yt-dlp 提取 YouTube 视频字幕，支持视频摘要、问答和内容提取，仅适用于有字幕的视频 | 📦 |
| [web-search-plus-zh](skills/web-search-plus-zh/SKILL.md) | 智能自动路由的统一搜索——连接 7 个搜索提供商（Serper/Tavily/Querit/Exa/Perplexity/You.com/SearXNG），通过多信号分析自动选择最佳引擎，附带置信度评分 | 💰📦 |

### 🧠 Agent 增强

| 技能 | 描述 | 标记 |
|------|------|------|
| [proactive-self-improving-agent-zh](skills/proactive-self-improving-agent-zh/SKILL.md) | Agent 自进化引擎——7 种触发条件自动捕获经验（错误/纠正/最佳实践/任务回顾等），结构化记录并去重，经验出现 ≥3 次自动晋升为永久规则，ADL 防漂移 + VFM 价值评分保障安全进化 | ✅ |
| [agent-autonomy-kit-zh](skills/agent-autonomy-kit-zh/SKILL.md) | Agent 自主执行工具包——任务队列系统 + 主动心跳 + 持续运行，将智能体从被动响应转变为主动执行，含队列检查脚本与任务后协议 | ✅ |
| [proactive-agent-zh](skills/proactive-agent-zh/SKILL.md) | 主动式 Agent 架构——WAL 预写日志、工作缓冲区、压缩恢复、自主定时任务、反向提示与安全演进护栏，让 Agent 预测需求并在上下文丢失中存活 | ✅ |
| [self-improving-zh](skills/self-improving-zh/SKILL.md) | 自我反思与持续改进引擎：自动捕获纠正和偏好信号，分层记忆存储（热/温/冷），模式出现3次后自动晋升为永久规则 | ✅ |
| [self-improvement-zh](skills/self-improvement-zh/SKILL.md) | 捕获学习内容、错误和纠正以实现持续改进 | ✅ |
| [self-improving-proactive-agent-zh](skills/self-improving-proactive-agent-zh/SKILL.md) | 自我进化引擎——自动捕获纠正/偏好/经验，HOT/WARM/COLD 三层记忆按需加载，模式≥3次自动晋升规则，项目级命名空间隔离，心跳驱动维护，零外部依赖 | ✅ |
| [elite-longterm-memory-zh](skills/elite-longterm-memory-zh/SKILL.md) | 6 层记忆架构——WAL 热内存 + LanceDB 向量语义搜索 + Git Notes 知识图谱 + Mem0 自动事实提取 + 可选云备份，跨会话永不丢失上下文 | 📦💰 |
| [memory-setup-zh](skills/memory-setup-zh/SKILL.md) | 持久化记忆系统配置向导——启用 memorySearch 语义搜索（Voyage/OpenAI/Local 三种嵌入提供商），创建 MEMORY.md + 每日日志目录结构，配置 Agent 记忆指令，解决跨会话上下文丢失 | 📦💰 |
| [ai-persona-os-zh](skills/ai-persona-os-zh/SKILL.md) | AI 人格操作系统——统一命令入口管理助手完整生命周期：设置向导（4 种预设模板）生成工作区文件、Never-Forget 检查点协议防止上下文丢失、心跳健康监控与记忆维护、结构化学习捕获与模式自动晋升（≥3 次晋升永久规则） | ✅ |
| [capability-evolver-zh](skills/capability-evolver-zh/SKILL.md) | 基于 GEP 协议的 Agent 自进化引擎——自动分析运行时历史识别故障与低效，通过可复用基因和胶囊实现协议约束的可审计进化，支持 7 种策略、审查模式与自动回滚，需注册 EvoMap 节点 | 📦 |
| [evolver-zh](skills/evolver-zh/SKILL.md) | AI 智能体自进化引擎——分析运行时历史记录，识别改进点，在协议约束下执行进化，支持 7 种策略、审查模式与自动回滚，需注册 EvoMap 节点 | 📦 |
| [freeride-zh](skills/freeride-zh/SKILL.md) | 按 OpenRouter 免费模型质量自动排名并配置主模型与备用模型，限流时自动故障切换，零成本使用 AI | 💰📦 |
| [auto-updater-zh](skills/auto-updater-zh/SKILL.md) | 通过 cron 每日自动更新 Clawdbot 核心与所有已安装技能，检查可用版本、应用更新并向用户发送变更摘要 | 📦 |
| [ontology-zh](skills/ontology-zh/SKILL.md) | 类型化知识图谱——实体/关系建模、Schema 约束验证（属性/枚举/基数/无环）、CLI 增删改查与图变换规划，支持跨技能数据共享 | ✅ |
| [ontology-v4-zh](skills/ontology-v4-zh/SKILL.md) | 类型化知识图谱 V4——完整实体/关系建模与约束验证系统，支持 12+ 核心类型（Person/Project/Task/Event/Document 等）、Schema 定义与验证、图遍历查询、图变换规划、跨技能数据共享，内置脚本工具链 | ✅ |

### 📱 工具与服务集成

| 技能 | 描述 | 标记 |
|------|------|------|
| [gog-zh](skills/gog-zh/SKILL.md) | Google Workspace 统一 CLI 工具，覆盖 Gmail 收发、日历事件查询、云盘搜索、联系人管理、表格读写与文档导出，需 OAuth 认证 | 📦 |
| [gmail-zh](skills/gmail-zh/SKILL.md) | 通过官方 Gmail API 进行邮件收发与管理——邮件读取/发送（HTML+附件）、线程管理、标签管理、多条件搜索、多账户支持，通过 Maton.ai 网关自动注入 OAuth 令牌 | 💰📦 |
| [notion-skill-zh](skills/notion-skill-zh/SKILL.md) | 通过官方 Notion API 读写页面、管理数据库、追加内容块，支持多配置文件和架构变更 | 💰📦 |
| [slack-zh](skills/slack-zh/SKILL.md) | 通过 Bot Token 管理 Slack 工作区——频道与私信消息收发编辑、Emoji 反应互动、置顶管理、成员信息查询与自定义表情列表 | ✅ |
| [trello-zh](skills/trello-zh/SKILL.md) | Trello 看板完整管理——通过 curl 调用 REST API 实现看板/列表/卡片的查询、创建、移动、评论与归档，需 API Key + Token 认证 | 📦 |
| [caldav-calendar-zh](skills/caldav-calendar-zh/SKILL.md) | CalDAV 日历本地同步与事件管理——基于 vdirsyncer + khal 实现 iCloud/Google/Fastmail/Nextcloud 等日历的双向同步、事件 CRUD、搜索与格式化输出，数据完全本地存储 | 📦 |
| [obsidian-zh](skills/obsidian-zh/SKILL.md) | 基于 obsidian-cli 的 Obsidian 笔记库管理——日常笔记创建、前缀自动分区（任务/想法/问题/洞察）、集中式任务日志与艾森豪威尔优先级标签 | 📦 |
| [imap-smtp-email-zh](skills/imap-smtp-email-zh/SKILL.md) | 通过 IMAP/SMTP 收发邮件——检查新邮件/未读、获取内容、搜索邮箱、标记已读/未读、发送带附件邮件，支持多账号，兼容 Gmail/Outlook/163 等主流邮箱 | 📦 |
| [weather-zh](skills/weather-zh/SKILL.md) | 通过 wttr.in + Open-Meteo 免费查询天气——紧凑/完整预报、中文输出、JSON API 程序化调用，仅需 curl | ✅ |
| [minara-zh](skills/minara-zh/SKILL.md) | 加密货币全功能交易 CLI——现货兑换、永续合约、AI 自动交易、限价单、充值提现（信用卡/加密货币）、x402 协议支付、市场发现与 AI 研究聊天，支持 EVM + Solana 多链 | 💰📦 |

### 🎨 设计与内容创作

| 技能 | 描述 | 标记 |
|------|------|------|
| [frontend-design-zh](skills/frontend-design-zh/SKILL.md) | 从纯文本需求生成生产级网站，无需设计稿。React 18 + TypeScript + Tailwind + shadcn/ui + Framer Motion，反 AI 套路美学，支持 Vite 静态和 Next.js 双工作流 | 📦 |
| [SuperDesign-zh](skills/SuperDesign-zh/SKILL.md) | 前端设计系统参考——结构化工作流（ASCII 线框→主题→动画→落地），暗色/野兽派/玻璃拟态三套视觉风格模板，oklch 配色 + Google Fonts + Tailwind 实现，内置无障碍与响应式规范 | ✅ |
| [ui-ux-pro-max-zh](skills/ui-ux-pro-max-zh/README.md) | UI/UX 设计系统推理引擎——161 条行业规则自动生成完整设计系统（布局模式 + 风格 + 配色 + 排版 + 反模式检查），67 种风格、161 套配色、57 种字体搭配、25 种图表类型、13 种技术栈，自然语言触发 | 📦 |
| [nano-banana-pro-zh](skills/nano-banana-pro-zh/SKILL.md) | 通过 Gemini 3 Pro Image API 生成与编辑图像——支持文生图、图生图及 1K/2K/4K 分辨率，内置草稿→迭代→最终渐进工作流 | 💰📦 |
| [nano-banana-pro-prompts-recommend-skill-zh](skills/nano-banana-pro-prompts-recommend-skill-zh/SKILL.md) | Nano Banana Pro 提示词推荐专家——10,000+ 经过验证的图像生成提示词库，附带示例图片，适用于 Nano Banana Pro、Gemini、GPT Image、Midjourney、DALL-E 等多种AI模型 | 💰📦 |
| [content-ideas-generator-zh](skills/content-ideas-generator-zh/SKILL.md) | 从新闻通讯/笔记/脚本等素材中提炼反直觉洞察，生成 5 个高互动社交媒体帖子大纲——包含悖论框架、转变弧线、异议预判与语言技巧应用，仅输出大纲不写全文 | ✅ |
| [marketing-mode-zh](skills/marketing-mode-zh/SKILL.md) | 营销战略知识库——以专家人"Mark"提供 140+ 实战营销方法，覆盖 SEO、付费广告、社交媒体、邮件营销、PR、产品驱动增长等 23 个子领域，含 5 阶段发布框架与转化优化体系 | 📦 |
| [marketing-skills-zh](skills/marketing-skills-zh/SKILL.md) | 23 个营销模块实战参考库——覆盖 CRO、SEO、文案、付费广告、邮件序列、定价策略等全链路，提供即用检查清单与框架模板 | ✅ |
| [admapix-zh](skills/admapix-zh/SKILL.md) | 竞品广告素材搜索引擎——按关键词/类型/地区/时间筛选，支持视频、图片与试玩广告，结果以 H5 可视化页面展示 | 💰📦 |
| [xiaohongshu-mcp-zh](skills/xiaohongshu-mcp-zh/SKILL.md) | 小红书（RedNote）内容运营自动化——图文/视频发布、笔记搜索与趋势追踪、详情与评论分析、用户资料与信息流管理，需本地部署 MCP 服务器 | 📦 |
| [ai-business-search-zh](skills/ai-business-search-zh/SKILL.md) | AI 商家搜索优化——引导本地商家完成 7 天 AI 搜索优化流程，生成媒体文章、选择发布渠道，确保商家出现在豆包/DeepSeek、通义千问等 AI 推荐中 | ✅ |

### 🛠️ 开发工具

| 技能 | 描述 | 标记 |
|------|------|------|
| [github-zh](skills/github-zh/SKILL.md) | 通过 `gh` CLI 管理 GitHub 全流程——Issues、PR 生命周期、CI/CD 监控与重跑、仓库操作，支持 `gh api` 高级查询与 JSON 结构化输出 | 📦 |
| [mcporter-zh](skills/mcporter-zh/SKILL.md) | MCP 统一客户端与代码生成 CLI——零配置自动发现（Cursor/Claude/Codex 等），直接调用 MCP 工具，一键生成 TypeScript 类型定义与独立 CLI，内置 OAuth 认证与临时连接 | 📦 |
| [playwright-mcp-zh](skills/playwright-mcp-zh/SKILL.md) | 基于 Playwright MCP 的浏览器自动化——导航、点击、输入、表单填充、数据提取、截图与 JS 执行，支持 Chrome/Firefox/WebKit，内置沙箱与安全主机限制 | 📦 |
| [Agent-Browser-zh](skills/Agent-Browser-zh/SKILL.md) | AI 专用无头浏览器自动化 CLI——基于无障碍树确定性引用驱动导航、点击、填充，Rust+Node.js 架构，支持会话管理与 JSON 结构化输出 | 📦 |
| [browser-use-zh](skills/browser-use-zh/SKILL.md) | browser-use 浏览器自动化 CLI——后台守护进程保持浏览器持久打开（约 50ms 延迟），支持导航、点击、输入、表单填充、截图与数据提取，适合网页测试与爬取 | 📦 |
| [desktop-control-zh](skills/desktop-control-zh/SKILL.md) | 基于 PyAutoGUI 的桌面自动化——像素级鼠标/键盘控制、屏幕截图与图像识别、窗口管理、剪贴板操作，内置故障保护机制 | 📦 |
| [skill-creator-zh](skills/skill-creator-zh/SKILL.md) | AgentSkills 全生命周期工具——创建、编辑、验证、打包，内置初始化脚本和规范检查器 | ✅ |
| [clawhub-zh](skills/clawhub-zh/SKILL.md) | 通过 ClawHub CLI 管理 clawhub.com 上的智能体技能——搜索发现、安装指定版本、哈希匹配增量更新（支持全量）、列出已装技能及发布新技能，支持自定义注册表与工作目录 | 📦 |
| [find-skills-zh](skills/find-skills-zh/SKILL.md) | 通过 `npx skills` CLI 搜索 skills.sh 开放技能库，按关键词发现并安装智能体扩展能力 | 📦 |
| [n8n-workflow-automation-zh](skills/n8n-workflow-automation-zh/SKILL.md) | n8n 工作流设计与 JSON 输出——内置幂等性、重试退避、审计日志与人机协作审核队列，确保零静默失败 | ✅ |
| [automation-workflows-zh](skills/automation-workflows-zh/SKILL.md) | 独立创业者自动化运营指南——任务审计方法论、Zapier/Make/n8n 工具选型、工作流设计模板、构建测试清单与 ROI 计算 | ✅ |
| [clawddocs-zh](skills/clawddocs-zh/SKILL.md) | Clawdbot 文档智能助手——决策树导航、多模式搜索（站点地图/关键词/全文索引）、文档获取与版本变更跟踪，内置各功能配置片段 | ✅ |
| [byterover-zh](skills/byterover-zh/SKILL.md) | 为 AI 编码助手管理项目知识的上下文工程平台——自动整理模式、决策与经验，Docker 隔离运行，团队间上下文自动同步 | 💰📦 |
| [qmd-zh](skills/qmd-zh/SKILL.md) | 本地文件混合搜索引擎——集合级索引管理，支持 BM25 关键词搜索、Ollama 向量语义搜索及混合重排序，内置 MCP 服务器模式，需安装 qmd CLI 及 Ollama | 📦 |
| [qa-engineer-zh](skills/qa-engineer-zh/SKILL.md) | 软件测试与质量保证——发现 bug、执行测试、生成报告 | ✅ |
| [model-usage-zh](skills/model-usage-zh/SKILL.md) | 通过 CodexBar CLI 从本地成本日志汇总 Codex/Claude 按模型使用量与成本，支持当前模型快照或全量明细，文本/JSON 双输出 | 📦 |
| [filesystem-management-zh](skills/filesystem-management-zh/SKILL.md) | 高级文件系统操作工具——智能文件列表/搜索/批量复制/目录分析，支持 Glob/正则过滤、树形可视化与安全防护，专为 AI 智能体优化 | 📦 |
| [opencode-controller-zh](skills/opencode-controller-zh/SKILL.md) | Opencode 编码工作流控制器——管理会话与模型选择、Plan/Build 智能体模式切换、协调编码任务的全流程，需安装 Opencode CLI | 📦 |
| [api-gateway-zh](skills/api-gateway-zh/SKILL.md) | 通过 Maton.ai 代理统一访问 100+ 第三方 API（Slack、HubSpot、Google Workspace、Notion、Stripe 等），托管 OAuth 自动注入令牌，支持原生端点透传与多连接管理 | 💰 |

### 📄 文档与文件处理

| 技能 | 描述 | 标记 |
|------|------|------|
| [word-docx-zh](skills/word-docx-zh/SKILL.md) | OOXML 感知的 Word 文档精准编辑——样式/编号/修订/批注/域/节完整支持，往返编辑格式不漂移 | ✅ |
| [md2wechat-skill-zh](skills/md2wechat-skill-zh/SKILL.md) | Markdown转公众号工具——一键将Markdown转换为精美排版并自动上传到微信草稿箱，支持15种AI主题样式和批量发布 | 📦 |
| [gongwen-zh](skills/gongwen-zh/SKILL.md) | 按党政机关公文格式标准自动排版 | ✅ |
| [gongwen-writing-zh](skills/gongwen-writing-zh/SKILL.md) | 政府公文写作助手——按标准模板起草公文，包含完整的写作流程、结构模板、风险分析等，支持多种公文类型 | ✅ |
| [nano-pdf-zh](skills/nano-pdf-zh/SKILL.md) | 通过 nano-pdf CLI 用自然语言指令编辑 PDF 特定页面，支持标题修改、内容更正等精准操作 | 📦 |
| [excel-xlsx-zh](skills/excel-xlsx-zh/SKILL.md) | OOXML 感知的 Excel 工作簿精准编辑——公式/日期/数据类型完整保留，往返编辑格式不漂移，pandas+openpyxl 双引擎按场景自动选择 | 📦 |
| [markdown-converter-zh](skills/markdown-converter-zh/SKILL.md) | 基于 markitdown 将 15+ 种文档格式转为 Markdown，支持图片 OCR、音频转录、YouTube 字幕提取及 ZIP 解包 | 📦 |
| [video-frames-zh](skills/video-frames-zh/SKILL.md) | 基于 ffmpeg 提取视频帧——支持指定时间点截取单帧画面，JPG/PNG 双格式输出，用于视频预览与截图 | 📦 |
| [openai-whisper-zh](skills/openai-whisper-zh/SKILL.md) | 本地运行 OpenAI Whisper 模型进行语音转录与翻译，支持多音频格式、多模型尺寸选择，数据不出本地 | 📦 |
| [humanizer-zh](skills/humanizer-zh/SKILL.md) | 检测并移除 AI 写作痕迹，基于 Wikipedia 指南修复 24 种 AI 文本模式（过度夸大、促销语言、模糊归因等），让文本更自然 | ✅ |
| [humanize-ai-text-zh](skills/humanize-ai-text-zh/SKILL.md) | AI 文本人性化——基于 Wikipedia 指南检测修复 24 种 AI 写作模式，三阶段工作流（草稿→反 AI 自检→最终修订），附"注入灵魂"个性写作指导 | ✅ |
| [summarize-zh](skills/summarize-zh/SKILL.md) | 基于 summarize CLI 的多格式智能摘要工具——网页/PDF/Office 文档/图片/音频/YouTube 视频一键摘要，支持批量处理与自定义长度限制 | 📦 |

### 📊 数据分析

| 技能 | 描述 | 标记 |
|------|------|------|
| [stock-analysis-zh](skills/stock-analysis-zh/SKILL.md) | 基于 Yahoo Finance API 的全面股票研究——公司概览、技术指标与估值、多股对比、内部人士交易追踪及 SEC 监管文件，支持美股与国际市场 | 💰 |
| [stock-market-pro-zh](skills/stock-market-pro-zh/SKILL.md) | 基于 Yahoo Finance 的本地股票分析——实时报价、基本面数据、ASCII 趋势图与高分辨率技术图表（RSI/MACD/BB/VWAP/ATR），一键综合报告，可选新闻与期权扩展 | 📦 |
| [stock-watcher-zh](skills/stock-watcher-zh/SKILL.md) | A股自选股管理与监控——添加/删除/查看自选股列表，基于同花顺数据源获取实时行情摘要与技术指标，支持沪深A股及科创板市场 | 📦 |

### 🔒 安全与防护

| 技能 | 描述 | 标记 |
|------|------|------|
| [moltguard-zh](skills/moltguard-zh/SKILL.md) | OpenClaw 原生安全插件——通过 MoltGuard API 检测工具结果与消息中的提示注入，可选本地净化网关在发送前脱敏 PII（银行卡/密码等），首次使用自动注册 API 密钥 | 💰 |
| [skill-vetter-zh](skills/skill-vetter-zh/SKILL.md) | 技能安装前安全审查协议——来源验证、14项危险信号检测、权限范围评估与四级风险分类，输出结构化审查报告 | ✅ |
| [edgeone-clawscan-zh](skills/edgeone-clawscan-zh/SKILL.md) | 腾讯朱雀实验室 A.I.G 驱动的 OpenClaw 安全扫描——全环境 4 步体检（配置审计/供应链风险/CVE 漏洞/隐私泄露）与单技能深度安全审查，云端情报+本地审计双引擎 | ✅ |
| [clawsec-suite-zh](skills/clawsec-suite-zh/SKILL.md) | AI代理安全技能套件——提供统一的安全监控、完整性验证和威胁情报，保护代理架构免受提示注入和恶意指令攻击，包含文件完整性保护、实时安全公告和安全审计功能 | 📦💰 |

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
| **轻量替代** | `humanize-ai-text-zh` | 结构更简洁，额外提供三阶段工作流和个性注入指导 |

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
