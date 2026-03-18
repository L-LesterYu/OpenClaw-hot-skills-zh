# Minara 示例

## 1 — 登录与账户

对于设备登录交接：如果 CLI 输出了验证 URL 和/或设备码，请将其原样转达给用户，要求用户完成浏览器验证，然后仅在用户确认完成后继续。

```bash
minara login                       # 交互式（默认设备码，或邮箱）
minara login --device              # 设备码流程：将 URL/验证码转达给用户进行浏览器验证
minara login -e user@example.com   # 邮箱验证码登录
minara account                     # 查看账户信息 + 钱包地址
minara deposit spot                # 显示现货充值地址（EVM + Solana）
```

## 2 — 代币兑换

链根据代币自动检测。

```bash
# 交互式：方向 → 代币 → 数量
minara swap

# 按代码（链自动检测）
minara swap -s buy -t '$BONK' -a 100
minara swap -s buy -t '$ETH' -a 50
minara swap -s sell -t '$SOL' -a 200

# 卖出全部余额
minara swap -s sell -t '$NVDAx' -a all

# 按合约地址
minara swap -s buy -t DezXAZ8z7PnrnRJjz3wXBoRgixCa6xjnB7YaB1pPB263 -a 100

# 模拟运行（不实际执行）
minara swap -s buy -t '$ETH' -a 50 --dry-run
```

## 3 — 转账与提现

```bash
# 转账（交互式）
minara transfer

# 提现到外部钱包
minara withdraw -c solana -t '$SOL' -a 10 --to <地址>
minara withdraw   # 交互式（接受代码或地址）
```

## 4 — 钱包与投资组合

```bash
minara balance                 # 快速总额：现货 + 永续合约 USDC/USDT 余额
minara assets                  # 完整概览：现货持仓 + 永续合约账户
minara assets spot             # 现货钱包：投资组合价值、成本、盈亏、持仓
minara assets perps            # 永续合约：权益、保证金、仓位
minara assets spot --json      # JSON 输出

# 充值
minara deposit                 # 交互式：现货 / 永续合约 / 信用卡购买
minara deposit spot            # 显示现货充值地址（EVM + Solana）
minara deposit perps           # 永续合约：显示 Arbitrum 地址，或从现货转入永续合约
minara deposit buy             # 通过 MoonPay 用信用卡购买加密货币（打开浏览器）
```

## 5 — 永续合约

```bash
# 为永续合约账户入金
minara perps deposit -a 100

# 设置杠杆
minara perps leverage

# 下单（交互式：合约、方向、数量、价格）
minara perps order

# 查看仓位
minara perps positions

# 平仓
minara perps close                    # 交互式：选择要平的仓位
minara perps close --all              # 平掉所有仓位（非交互式）
minara perps close --symbol BTC       # 平掉 BTC 仓位（非交互式）
minara perps close --all --yes        # 平掉所有，跳过确认

# 取消订单
minara perps cancel

# 从永续合约提现
minara perps withdraw -a 50

# AI 分析 → 可选快速下单
minara perps ask

# AI 自动交易策略
minara perps autopilot

# 历史记录
minara perps trades
minara perps fund-records
```

## 6 — AI 聊天

```bash
# 单个问题
minara chat "BTC 现在什么价格？"

# 高质量模式
minara chat --quality "分析 ETH 下周展望"

# 推理模式
minara chat --thinking "比较 SOL 和 AVAX 生态系统增长"

# 交互式 REPL
minara chat
# >>> 现在最好的 DeFi 收益是什么？
# >>> /help
# >>> exit

# 继续现有对话
minara chat -c <chatId>

# 列出/回放过去的对话
minara chat --list
minara chat --history <chatId>
```

## 7 — 市场发现

```bash
minara discover trending           # 热门代币
minara discover trending stocks    # 热门股票
minara discover search SOL         # 搜索代币/股票
minara discover search AAPL        # 按名称搜索股票
minara discover fear-greed         # 加密货币恐惧贪婪指数
minara discover btc-metrics        # 比特币链上指标
minara discover trending --json    # JSON 输出
```

## 8 — 限价单

```bash
minara limit-order create          # 交互式：代币、价格、方向、数量、有效期
minara limit-order list            # 列出所有订单
minara limit-order cancel abc123   # 按 ID 取消
```

## 9 — x402 协议支付

当 HTTP API 返回 **402 Payment Required** 并带有 x402 头时，Agent 可以使用 Minara 钱包进行支付。

```bash
# 1. 支付前检查余额
minara balance

# 2. 向 x402 服务付款（USDC 转账到服务的支付地址）
#    示例：服务需要在 Base 链上支付 0.01 USDC
minara transfer
#    → 代币：USDC
#    → 金额：0.01
#    → 接收方：<402 头中的服务支付地址>
#    → 链：base

# 3. 确保钱包已充值以备未来 x402 支付
minara deposit buy                 # 通过 MoonPay 信用卡入金
minara deposit spot                # 或显示充值地址以接收加密货币
```

## 10 — 高级版与订阅

```bash
minara premium plans               # 查看套餐
minara premium status              # 当前订阅状态
minara premium subscribe           # 订阅/升级
minara premium buy-credits         # 购买积分
minara premium cancel              # 取消订阅
```
