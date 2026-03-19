#!/bin/bash
# AI 人格操作系统 — 配置验证器
# 只读：此脚本仅读取配置文件来验证设置。
# 它不会修改任何文件、进行网络调用或安装任何东西。
# 检查内容：工作区结构、openclaw 配置、心跳设置。
#
# 检查所有必需的 OpenClaw/Clawdbot 设置并标记缺失项
# 运行：./scripts/config-validator.sh

echo "🔍 AI Persona OS 配置验证器"
echo "================================="
echo ""

ISSUES=0
WARNINGS=0

# 查找配置文件
CONFIG_FILE=""
for f in openclaw.json clawdbot-mac.json clawdbot.json; do
  if [ -f "$HOME/.openclaw/$f" ]; then
    CONFIG_FILE="$HOME/.openclaw/$f"
    break
  fi
  if [ -f "$f" ]; then
    CONFIG_FILE="$f"
    break
  fi
done

if [ -z "$CONFIG_FILE" ]; then
  echo "❌ 未找到配置文件（已检查 openclaw.json、clawdbot-mac.json、clawdbot.json）"
  echo "   检查位置：~/.openclaw/ 和当前目录"
  ISSUES=$((ISSUES + 1))
else
  echo "✅ 配置文件：$CONFIG_FILE"
  
  # 检查心跳配置
  echo ""
  echo "## 心跳配置"
  
  if grep -q '"heartbeat"' "$CONFIG_FILE" 2>/dev/null; then
    echo "✅ 心跳配置段存在"
    
    if grep -q '"every"' "$CONFIG_FILE" 2>/dev/null; then
      echo "✅ heartbeat.every 已设置"
    else
      echo "❌ heartbeat.every 缺失 — 心跳不会自动触发"
      echo "   修复：在心跳配置中添加 \"every\": \"30m\""
      ISSUES=$((ISSUES + 1))
    fi
    
    if grep -q '"target"' "$CONFIG_FILE" 2>/dev/null; then
      echo "✅ heartbeat.target 已设置"
    else
      echo "❌ heartbeat.target 缺失 — 心跳不会正确路由"
      echo "   修复：在心跳配置中添加 \"target\": \"last\""
      ISSUES=$((ISSUES + 1))
    fi
    
    if grep -q '"ackMaxChars"' "$CONFIG_FILE" 2>/dev/null; then
      echo "✅ heartbeat.ackMaxChars 已设置（HEARTBEAT_OK 抑制）"
    else
      echo "⚠️  heartbeat.ackMaxChars 缺失 — HEARTBEAT_OK 将在聊天中显示"
      echo "   修复：在心跳配置中添加 \"ackMaxChars\": 20"
      WARNINGS=$((WARNINGS + 1))
    fi
    
    if grep -q '"prompt"' "$CONFIG_FILE" 2>/dev/null; then
      echo "✅ heartbeat.prompt 覆盖已设置"
    else
      echo "⚠️  heartbeat.prompt 缺失 — 智能体可能使用旧格式"
      echo "   修复：添加 prompt 覆盖（参见 references/heartbeat-automation.md）"
      WARNINGS=$((WARNINGS + 1))
    fi
  else
    echo "❌ 未找到心跳配置段"
    echo "   修复：在 agents.defaults 中添加心跳配置（参见 references/heartbeat-automation.md）"
    ISSUES=$((ISSUES + 1))
  fi
  
  # 检查 Discord requireMention
  echo ""
  echo "## Discord 配置"
  
  if grep -q '"guilds"' "$CONFIG_FILE" 2>/dev/null; then
    GUILD_COUNT=$(grep -c '"requireMention"' "$CONFIG_FILE" 2>/dev/null || echo "0")
    if [ "$GUILD_COUNT" -gt 0 ]; then
      echo "✅ 在 $GUILD_COUNT 个服务器中找到 requireMention"
      
      # 检查是否有服务器缺少 requireMention
      TOTAL_GUILDS=$(grep -c '"[0-9]\{17,\}"' "$CONFIG_FILE" 2>/dev/null || echo "0")
      if [ "$GUILD_COUNT" -lt "$TOTAL_GUILDS" ]; then
        echo "⚠️  部分服务器可能缺少 requireMention: true"
        echo "   修复：为所有 Discord 服务器设置 requireMention: true"
        WARNINGS=$((WARNINGS + 1))
      fi
    else
      echo "⚠️  未找到 requireMention 设置 — 智能体可能未受邀请就响应"
      echo "   修复：为所有 Discord 服务器添加 requireMention: true"
      WARNINGS=$((WARNINGS + 1))
    fi
  else
    echo "ℹ️  未配置 Discord 服务器（不使用 Discord 则正常）"
  fi
fi

# 检查工作区文件
echo ""
echo "## 工作区文件"

for f in SOUL.md USER.md MEMORY.md AGENTS.md HEARTBEAT.md SECURITY.md; do
  if [ -f "$f" ]; then
    SIZE=$(wc -c < "$f")
    if [ "$f" = "MEMORY.md" ] && [ "$SIZE" -gt 4096 ]; then
      echo "⚠️  $f 存在（$SIZE 字节 — 超过 4KB 限制）"
      echo "   修复：清理 MEMORY.md — 将旧事实归档到 memory/"
      WARNINGS=$((WARNINGS + 1))
    else
      echo "✅ $f 存在（$SIZE 字节）"
    fi
  else
    echo "❌ $f 缺失"
    echo "   修复：从技能 assets/ 复制或运行设置向导"
    ISSUES=$((ISSUES + 1))
  fi
done

# 检查 VERSION 文件
echo ""
echo "## 版本跟踪"

if [ -f "VERSION" ]; then
  VERSION=$(cat VERSION | tr -d '[:space:]')
  echo "✅ VERSION 文件：v$VERSION"
else
  echo "⚠️  工作区中无 VERSION 文件"
  echo "   修复：从技能 assets/VERSION 复制"
  WARNINGS=$((WARNINGS + 1))
fi

# 检查 HEARTBEAT.md 大小（旧模板检测）
if [ -f "HEARTBEAT.md" ]; then
  LINES=$(wc -l < "HEARTBEAT.md")
  if [ "$LINES" -gt 50 ]; then
    echo ""
    echo "⚠️  HEARTBEAT.md 有 $LINES 行 — 可能是旧模板"
    echo "   修复：替换为 assets/HEARTBEAT-template.md"
    WARNINGS=$((WARNINGS + 1))
  fi
fi

# 检查记忆目录
echo ""
echo "## 记忆系统"

if [ -d "memory" ]; then
  LOG_COUNT=$(ls memory/*.md 2>/dev/null | wc -l)
  echo "✅ memory/ 目录存在（$LOG_COUNT 个日志文件）"
  
  if [ -d "memory/archive" ]; then
    echo "✅ memory/archive/ 存在"
  else
    echo "⚠️  memory/archive/ 缺失"
    echo "   修复：mkdir -p memory/archive"
    WARNINGS=$((WARNINGS + 1))
  fi
  
  # 检查过期日志
  OLD_LOGS=$(find memory/ -maxdepth 1 -name "*.md" -mtime +90 2>/dev/null | wc -l)
  if [ "$OLD_LOGS" -gt 0 ]; then
    echo "⚠️  有 $OLD_LOGS 份超过 90 天的日志 — 应归档"
    WARNINGS=$((WARNINGS + 1))
  fi
else
  echo "❌ memory/ 目录缺失"
  echo "   修复：mkdir -p memory/archive"
  ISSUES=$((ISSUES + 1))
fi

# 检查学习目录
if [ -d ".learnings" ]; then
  echo "✅ .learnings/ 目录存在"
else
  echo "⚠️  .learnings/ 目录缺失"
  echo "   修复：mkdir -p .learnings"
  WARNINGS=$((WARNINGS + 1))
fi

# 检查 ESCALATION.md
echo ""
echo "## 升级协议"

if [ -f "ESCALATION.md" ]; then
  echo "✅ ESCALATION.md 存在"
else
  echo "⚠️  ESCALATION.md 缺失 — 智能体不会遵循结构化交接"
  echo "   修复：从技能 assets/ESCALATION-template.md 复制"
  WARNINGS=$((WARNINGS + 1))
fi

# 摘要
echo ""
echo "================================="
echo "摘要：$ISSUES 个问题，$WARNINGS 个警告"
if [ "$ISSUES" -eq 0 ] && [ "$WARNINGS" -eq 0 ]; then
  echo "🟢 所有检查通过 — 系统健康"
elif [ "$ISSUES" -eq 0 ]; then
  echo "🟡 无关键问题 — $WARNINGS 项可改进"
else
  echo "🔴 $ISSUES 个关键问题需要修复"
fi
echo ""
