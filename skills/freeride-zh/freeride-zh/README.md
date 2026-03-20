# 🎢 FreeRide

### 停止为 AI 付费。开始免费畅享。

[![License: MIT](https://img.shields.io/badge/License-MIT-green.svg)](https://opensource.org/licenses/MIT)
[![OpenClaw Compatible](https://img.shields.io/badge/OpenClaw-Compatible-blue.svg)](https://github.com/openclaw/openclaw)
[![OpenRouter](https://img.shields.io/badge/OpenRouter-30%2B%20Free%20Models-orange.svg)](https://openrouter.ai)

---

**FreeRide** 通过自动管理 OpenRouter 的免费模型，为 [OpenClaw](https://github.com/openclaw/openclaw) 提供无限免费 AI。

```
你：*触发了限流*
FreeRide："我来搞定。" *切换到下一个最佳模型*
你：*继续编程*
```

## 问题所在

你在使用 OpenClaw。你喜欢它。但是：

- 💸 API 成本快速累积
- 🚫 免费模型有限流限制
- 😤 手动切换模型很烦人
- 🤷 你不知道哪个免费模型真正好用

## 解决方案

一条命令。免费 AI。永远。

```bash
freeride auto
```

就这么简单。FreeRide 会：

1. **发现** OpenRouter 上 30+ 个免费模型
2. **排名** 按质量排序（上下文长度、能力、速度）
3. **设置** 最佳模型为主模型
4. **配置** 智能备用模型，当你触发限流时使用
5. **保留** 你现有的 OpenClaw 配置

## 安装

```bash
npx clawhub@latest install free-ride
cd ~/.openclaw/skills/freeride-zh
pip install -e .
```

完成。`freeride` 和 `freeride-watcher` 现在可以作为全局命令使用。

## 快速开始

### 1. 获取免费 OpenRouter 密钥

访问 [openrouter.ai/keys](https://openrouter.ai/keys) → 创建账户 → 生成密钥

无需信用卡。无试用期。真正免费。

### 2. 设置你的密钥

```bash
export OPENROUTER_API_KEY="sk-or-v1-..."
```

或添加到你的 OpenClaw 配置：

```bash
openclaw config set env.OPENROUTER_API_KEY "sk-or-v1-..."
```

### 3. 运行 FreeRide

```bash
freeride auto
```

### 4. 重启 OpenClaw

```bash
openclaw gateway restart
```

### 5. 验证是否生效

在 WhatsApp/Telegram/Discord 或仪表板上给你的代理发消息：

```
你：    /status
代理：  （显示免费模型名称 + token 数量）
```

完成。你现在已经运行在带自动故障切换的免费 AI 上。

## 你将获得

```
主模型: openrouter/nvidia/nemotron-3-nano-30b-a3b:free (256K 上下文)

备用模型:
  1. openrouter/free          ← 智能路由器（自动选择最佳可用模型）
  2. qwen/qwen3-coder:free    ← 编程利器
  3. stepfun/step-3.5:free    ← 快速响应
  4. deepseek/deepseek:free   ← 强大推理
  5. mistral/mistral:free     ← 可靠备用
```

当你触发限流时，OpenClaw 会自动尝试下一个模型。你继续工作。不中断。

## 命令

| 命令 | 作用 |
|---------|------|
| `freeride auto` | 自动配置最佳模型 + 备用模型 |
| `freeride list` | 查看所有 30+ 免费模型排名 |
| `freeride switch <model>` | 使用指定模型 |
| `freeride status` | 检查当前设置 |
| `freeride fallbacks` | 仅更新备用模型 |
| `freeride refresh` | 强制刷新模型缓存 |

### 专业技巧

```bash
# 已有心仪的模型？只添加备用：
freeride auto -f

# 想要更多备用模型以获得最大在线时间？
freeride auto -c 10

# 编程？切换到最佳编程模型：
freeride switch qwen3-coder

# 看看有什么可用的：
freeride list -n 30

# 更改后始终重启 OpenClaw：
openclaw gateway restart
```

## 模型排名机制

FreeRide 根据以下因素为每个模型评分（0-1）：

| 因素 | 权重 | 原因 |
|--------|--------|-----|
| 上下文长度 | 40% | 越长 = 能处理更大的代码库 |
| 能力 | 30% | 视觉、工具、结构化输出 |
| 新旧程度 | 20% | 新模型 = 更好的性能 |
| 提供商信任度 | 10% | Google、Meta、NVIDIA 等 |

**智能备用** `openrouter/free` 始终排第一 - 它根据你的请求需求自动选择。

## 使用 OpenClaw 代理测试

运行 `freeride auto` 和 `openclaw gateway restart` 后：

```bash
# 检查 OpenClaw 是否看到模型
openclaw models list

# 验证配置
openclaw doctor --fix

# 打开仪表板并聊天
openclaw dashboard
# 或在 WhatsApp/Telegram/Discord 上给你的代理发消息
```

用于验证的有用代理命令：

| 命令 | 告诉你什么 |
|---------|----------|
| `/status` | 当前模型 + token 使用量 |
| `/model` | 可用模型（你的免费模型应该列出） |
| `/new` | 用新模型开始新会话 |

## 监控器（自动轮换）

FreeRide 包含一个监控守护进程，监视限流并自动轮换模型：

```bash
# 运行一次（检查 + 如需要则轮换）
freeride-watcher

# 作为守护进程运行（持续监控）
freeride-watcher --daemon

# 强制轮换到下一个模型
freeride-watcher --rotate

# 检查监控器状态
freeride-watcher --status

# 清除限流冷却
freeride-watcher --clear-cooldowns
```

## 常见问题

**真的免费吗？**

是的。OpenRouter 为许多模型提供免费层。你只需要一个账户（无需信用卡）。

**限流怎么办？**

这正是重点所在。FreeRide 配置多个备用模型。当一个模型限流时，OpenClaw 自动切换到下一个。

**会搞乱我的 OpenClaw 配置吗？**

不会。FreeRide 只修改 `agents.defaults.model` 和 `agents.defaults.models`。你的网关、通道、插件、工作区、自定义指令 - 全部保留。

**哪些模型是免费的？**

运行 `freeride list` 查看当前可用性。它会变化，这就是 FreeRide 存在的原因。

**更改后需要重启 OpenClaw 吗？**

是的。在更改配置的 FreeRide 命令后运行 `openclaw gateway restart`。

## 算笔账

| 场景 | 每月成本 |
|----------|--------------|
| GPT-4 API | $50-200+ |
| Claude API | $50-200+ |
| OpenClaw + FreeRide | **$0** |

不客气。

## 要求

- [OpenClaw](https://github.com/openclaw/openclaw) 已安装（Node ≥22）
- Python 3.8+
- 免费 OpenRouter 账户（[获取密钥](https://openrouter.ai/keys)）

## 架构

```
┌──────────────┐     ┌──────────────┐     ┌──────────────────┐
│  你          │ ──→ │  FreeRide    │ ──→ │  OpenRouter API  │
│  "freeride   │     │              │     │  (30+ free       │
│   auto"      │     │  • 获取      │     │   models)        │
└──────────────┘     │  • 排名      │     └──────────────────┘
                     │  • 配置      │
                     └──────┬───────┘
                            │
                            ▼
                     ┌──────────────┐
                     │ ~/.openclaw/ │
                     │ openclaw.json│
                     └──────┬───────┘
                            │
                     openclaw gateway restart
                            │
                            ▼
                     ┌──────────────┐
                     │  OpenClaw    │
                     │  (免费 AI!)  │
                     └──────────────┘
```

## 贡献

发现 bug？想要功能？欢迎 PR。

```bash
cd ~/.openclaw/skills/freeride-zh

# 测试命令
freeride list
freeride status
freeride auto --help
```

## 相关项目

- [OpenClaw](https://github.com/openclaw/openclaw) - AI 编程代理
- [OpenRouter](https://openrouter.ai) - 模型路由器
- [ClawHub](https://github.com/clawhub) - 技能市场

## 许可证

MIT - 随心所欲。

---

<p align="center">
  <b>停止付费。开始畅享。</b>
  <br>
  <br>
  <a href="https://github.com/Shaivpidadi/FreeRide">⭐ 在 GitHub 上给我们点星</a>
  ·
  <a href="https://openrouter.ai/keys">🔑 获取 OpenRouter 密钥</a>
  ·
  <a href="https://github.com/openclaw/openclaw">🦞 安装 OpenClaw</a>
</p>
