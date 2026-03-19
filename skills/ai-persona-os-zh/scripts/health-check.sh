#!/bin/bash

# AI 人格操作系统 — 健康检查
# 验证工作区结构并识别问题
# 作者：Jeff J Hunter — https://jeffjhunter.com

set -e

WORKSPACE="${1:-$HOME/workspace}"

# 颜色
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[0;33m'
BOLD='\033[1m'
NC='\033[0m'

echo ""
echo -e "${BOLD}AI Persona OS — 健康检查${NC}"
echo "工作区：$WORKSPACE"
echo ""

issues=0
warnings=0

# 检查必需文件
echo "正在检查必需文件..."

required_files=("SOUL.md" "USER.md" "MEMORY.md" "AGENTS.md")
for file in "${required_files[@]}"; do
    if [ -f "$WORKSPACE/$file" ]; then
        echo -e "  ${GREEN}✓${NC} $file"
    else
        echo -e "  ${RED}✗${NC} $file — 缺失"
        ((issues++))
    fi
done

# 检查必需目录
echo ""
echo "正在检查目录..."

required_dirs=("memory" "memory/archive" ".learnings")
for dir in "${required_dirs[@]}"; do
    if [ -d "$WORKSPACE/$dir" ]; then
        echo -e "  ${GREEN}✓${NC} $dir/"
    else
        echo -e "  ${RED}✗${NC} $dir/ — 缺失"
        ((issues++))
    fi
done

# 检查 MEMORY.md 大小
echo ""
echo "正在检查文件大小..."

if [ -f "$WORKSPACE/MEMORY.md" ]; then
    size=$(wc -c < "$WORKSPACE/MEMORY.md")
    if [ "$size" -gt 4096 ]; then
        echo -e "  ${YELLOW}⚠${NC} MEMORY.md 为 $size 字节（应 < 4096）"
        ((warnings++))
    else
        echo -e "  ${GREEN}✓${NC} MEMORY.md 大小正常（$size 字节）"
    fi
fi

# 检查核心文件是否为空
echo ""
echo "正在检查文件内容..."

for file in "SOUL.md" "USER.md"; do
    if [ -f "$WORKSPACE/$file" ]; then
        words=$(wc -w < "$WORKSPACE/$file")
        if [ "$words" -lt 50 ]; then
            echo -e "  ${YELLOW}⚠${NC} $file 仅有 $words 个词（需要更多内容）"
            ((warnings++))
        else
            echo -e "  ${GREEN}✓${NC} $file 内容正常（$words 个词）"
        fi
    fi
done

# 检查过期日志
echo ""
echo "正在检查会话日志..."

if [ -d "$WORKSPACE/memory" ]; then
    stale=$(find "$WORKSPACE/memory" -maxdepth 1 -name "*.md" -type f -mtime +90 2>/dev/null | wc -l | tr -d ' ')
    if [ "$stale" -gt 0 ]; then
        echo -e "  ${YELLOW}⚠${NC} 有 $stale 份超过 90 天的日志（应归档）"
        ((warnings++))
    else
        echo -e "  ${GREEN}✓${NC} 无过期日志"
    fi
fi

# 摘要
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

if [ "$issues" -eq 0 ] && [ "$warnings" -eq 0 ]; then
    echo -e "${GREEN}${BOLD}✅ 所有检查通过！${NC}"
    exit 0
elif [ "$issues" -eq 0 ]; then
    echo -e "${YELLOW}${BOLD}⚠️  $warnings 个警告${NC}"
    exit 0
else
    echo -e "${RED}${BOLD}❌ $issues 个问题，$warnings 个警告${NC}"
    echo ""
    echo "运行 ./scripts/setup-wizard.sh 修复问题"
    exit 1
fi
