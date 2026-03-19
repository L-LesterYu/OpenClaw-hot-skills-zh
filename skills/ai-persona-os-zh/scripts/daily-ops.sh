#!/bin/bash

# AI 人格操作系统 — 每日运维
# 自动化每日运维检查
# 作者：Jeff J Hunter — https://jeffjhunter.com

WORKSPACE="${1:-$HOME/workspace}"
TODAY=$(date +%Y-%m-%d)

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "  AI Persona OS — 每日运维"
echo "  $(date)"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# 步骤 1：检查今天的日志是否存在
echo "📅 正在检查今天的日志..."
if [ ! -f "$WORKSPACE/memory/$TODAY.md" ]; then
    echo "   正在创建 $TODAY.md..."
    cat > "$WORKSPACE/memory/$TODAY.md" << EOF
# $TODAY — $(date +%A)

## 心跳

**完成时间：** $(date +%H:%M)
**状态：** 🟢
**关注焦点：** [设置您的关注焦点]

## 会话笔记

[今天发生了什么]

## 已做出的决策

| 决策 | 原因 | 时间 |
|------|------|------|

## 待办事项

### 进行中
- [ ] 

## 学习记录

- 

## 检查点

---

*由 AI 人格操作系统每日运维生成*
EOF
    echo "   ✓ 已创建"
else
    echo "   ✓ 已存在"
fi

# 步骤 2：检查昨天的待办事项
echo ""
echo "📋 正在检查昨天的待办事项..."
YESTERDAY=$(date -d "yesterday" +%Y-%m-%d 2>/dev/null || date -v-1d +%Y-%m-%d 2>/dev/null)
if [ -f "$WORKSPACE/memory/$YESTERDAY.md" ]; then
    pending=$(grep -c "^\- \[ \]" "$WORKSPACE/memory/$YESTERDAY.md" 2>/dev/null || echo "0")
    if [ "$pending" -gt 0 ]; then
        echo "   ⚠️  昨天有 $pending 个未完成事项"
    else
        echo "   ✓ 无待办事项"
    fi
else
    echo "   ○ 昨天无日志"
fi

# 步骤 3：检查 MEMORY.md 大小
echo ""
echo "💾 正在检查 MEMORY.md 大小..."
if [ -f "$WORKSPACE/MEMORY.md" ]; then
    size=$(wc -c < "$WORKSPACE/MEMORY.md")
    if [ "$size" -gt 4096 ]; then
        echo "   ⚠️  $size 字节（应 < 4096）"
    else
        echo "   ✓ $size 字节（正常）"
    fi
fi

# 步骤 4：检查学习系统待审查内容
echo ""
echo "📈 正在检查学习系统..."
if [ -d "$WORKSPACE/.learnings" ]; then
    pending=$(grep -c "pending" "$WORKSPACE/.learnings/LEARNINGS.md" 2>/dev/null || echo "0")
    if [ "$pending" -gt 5 ]; then
        echo "   ⚠️  $pending 条待处理学习内容（需要审查）"
    else
        echo "   ✓ $pending 条待处理学习内容"
    fi
fi

# 步骤 5：检查过期日志
echo ""
echo "🗄️  正在检查过期日志..."
stale=$(find "$WORKSPACE/memory" -maxdepth 1 -name "*.md" -type f -mtime +90 2>/dev/null | wc -l | tr -d ' ')
if [ "$stale" -gt 0 ]; then
    echo "   ⚠️  有 $stale 份日志需要归档（超过 90 天）"
else
    echo "   ✓ 无过期日志"
fi

# 摘要
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "  每日运维完成"
echo "  今日日志：memory/$TODAY.md"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
