# 动态仪表盘与子代理生成

静态仪表盘显示的是过时数据，需要不断手动更新。您希望实时监控多个数据源，无需构建自定义前端或遇到 API 速率限制。

这个工作流程创建一个实时仪表盘，通过生成子代理并行获取和处理数据：

• 同时监控多个数据源（API、数据库、GitHub、社交媒体）
• 为每个数据源生成子代理，避免阻塞并分散 API 负载
• 将结果聚合到统一仪表盘中（文本、HTML 或 Canvas）
• 每 N 分钟更新一次最新数据
• 当指标超过阈值时发送警报
• 在数据库中维护历史趋势以供可视化

## 痛点

构建自定义仪表盘需要数周时间。等仪表盘完成时，需求已经改变。顺序轮询多个 API 速度慢且会达到速率限制。您需要的是即时洞察，而不是经过周末编程后的结果。

## 功能说明

您通过对话式定义要监控的内容："跟踪 GitHub 星标、Twitter 提及、Polymarket 交易量和系统健康状态。" OpenClaw 生成子代理并行获取每个数据源，聚合结果，并将格式化的仪表盘发送到 Discord 或作为 HTML 文件。更新在定时任务调度上自动运行。

示例仪表盘部分：
- **GitHub**：星标、分支、开放问题、最近提交
- **社交媒体**：Twitter 提及、Reddit 讨论、Discord 活动
- **市场**：Polymarket 交易量、预测趋势
- **系统健康**：CPU、内存、磁盘使用率、服务状态

## 所需技能

- 子代理生成用于并行执行
- `github`（gh CLI）用于 GitHub 指标
- `bird`（Twitter）用于社交媒体数据
- `web_search` 或 `web_fetch` 用于外部 API
- `postgres` 用于存储历史指标
- Discord 或 Canvas 用于渲染仪表盘
- Cron 任务用于定时更新

## 如何设置

1. 设置指标数据库：
```sql
CREATE TABLE metrics (
  id SERIAL PRIMARY KEY,
  source TEXT, -- 例如："github"、"twitter"、"polymarket"
  metric_name TEXT,
  metric_value NUMERIC,
  timestamp TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE alerts (
  id SERIAL PRIMARY KEY,
  source TEXT,
  condition TEXT,
  threshold NUMERIC,
  last_triggered TIMESTAMPTZ
);
```

2. 为仪表盘更新创建一个 Discord 频道（例如 #dashboard）。

3. 提示 OpenClaw：
```text
您是我的动态仪表盘管理器。每 15 分钟，运行一个定时任务：

1. 并行生成子代理从以下数据源获取数据：
   - GitHub：星标、分支、开放问题、提交（过去 24 小时）
   - Twitter：提及 "@username" 的内容、情感分析
   - Polymarket：追踪市场的交易量
   - 系统：通过 shell 命令获取 CPU、内存、磁盘使用率

2. 每个子代理将结果写入指标数据库。

3. 聚合所有结果并格式化仪表盘：

📊 **仪表盘更新** — [时间戳]

**GitHub**
- ⭐ 星标：[数量]（+[变化]）
- 🍴 分支：[数量]
- 🐛 开放问题：[数量]
- 💻 提交（24h）：[数量]

**社交媒体**
- 🐦 Twitter 提及：[数量]
- 📈 情感倾向：[正面/负面/中性]

**市场**
- 📊 Polymarket 交易量：$[金额]
- 🔥 热门市场：[市场名称]

**系统健康**
- 💻 CPU：[使用率]%
- 🧠 内存：[使用率]%
- 💾 磁盘：[使用率]%

4. 发布到 Discord #dashboard 频道。

5. 检查警报条件：
   - 如果 GitHub 星标在 1 小时内变化 > 50 → 通知我
   - 如果系统 CPU > 90% → 发出警报
   - 如果 Twitter 出现负面情感激增 → 通知我

将所有指标存储在数据库中，供历史分析使用。
```

4. 可选：使用 Canvas 渲染带有图表和图形的 HTML 仪表盘。

5. 查询历史数据："显示过去 30 天的 GitHub 星标增长情况。"

## 相关链接

- [子代理并行处理](https://docs.openclaw.ai/subagents)
- [仪表盘设计原则](https://www.nngroup.com/articles/dashboard-design/)