# CodexBar CLI 快速参考（使用量 + 成本）

## 安装
- 应用程序：偏好设置 -> 高级 -> 安装 CLI
- 仓库：./bin/install-codexbar-cli.sh

## 命令
- 使用快照（web/cli 来源）：
  - codexbar usage --format json --pretty
  - codexbar --provider all --format json
- 本地成本使用（仅 Codex + Claude）：
  - codexbar cost --format json --pretty
  - codexbar cost --provider codex|claude --format json

## 成本 JSON 字段
负载是一个数组（每个 provider 一个）。
- provider, source, updatedAt
- sessionTokens, sessionCostUSD
- last30DaysTokens, last30DaysCostUSD
- daily[]: date, inputTokens, outputTokens, cacheReadTokens, cacheCreationTokens, totalTokens, totalCost, modelsUsed, modelBreakdowns[]
- modelBreakdowns[]: modelName, cost
- totals: totalInputTokens, totalOutputTokens, cacheReadTokens, cacheCreationTokens, totalTokens, totalCost

## 注意事项
- 成本使用仅限本地。它读取以下位置的 JSONL 日志：
  - Codex: ~/.codex/sessions/**/*.jsonl
  - Claude: ~/.config/claude/projects/**/*.jsonl 或 ~/.claude/projects/**/*.jsonl
- 如果需要 web 使用（非本地），请使用 codexbar usage（而不是 cost）。
