#!/bin/bash

# AI 人格操作系统 — 安全审计
# 每月运行以检查安全问题
# 作者：Jeff J Hunter — https://jeffjhunter.com
#
# 注意：此脚本用于检测工作区文件中泄露的凭据。
# 它不会存储、传输或处理任何秘密信息。
# 所有检查均为本地基于 grep 的模式扫描 — 不进行网络调用。

set -e

WORKSPACE="${1:-$HOME/workspace}"

# 颜色
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
BOLD='\033[1m'
NC='\033[0m'

echo -e "${BOLD}${BLUE}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "  🔒 AI Persona OS — 安全审计"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo -e "${NC}"
echo "工作区：$WORKSPACE"
echo "日期：$(date)"
echo ""

ISSUES_FOUND=0
WARNINGS_FOUND=0

# ═══════════════════════════════════════════════════════════════
# 检查 1：文件中的凭据
# ═══════════════════════════════════════════════════════════════
echo -e "${BOLD}[1/6] 正在检查文件中的凭据...${NC}"

# 常见凭据模式
PATTERNS=(
    "api[_-]?key"
    "api[_-]?secret"
    "password"
    "passwd"
    "secret[_-]?key"
    "access[_-]?token"
    "auth[_-]?token"
    "bearer"
    "sk-[a-zA-Z0-9]"
    "xox[baprs]-"
)

CRED_FOUND=false
for pattern in "${PATTERNS[@]}"; do
    if grep -r -i -l "$pattern" "$WORKSPACE/memory" 2>/dev/null | head -5 | grep -q .; then
        if [ "$CRED_FOUND" = false ]; then
            echo -e "  ${RED}⚠ 在记忆文件中发现潜在凭据：${NC}"
            CRED_FOUND=true
            ((ISSUES_FOUND++))
        fi
        grep -r -i -l "$pattern" "$WORKSPACE/memory" 2>/dev/null | head -3 | while read -r file; do
            echo "    - $file（模式：$pattern）"
        done
    fi
done

if [ "$CRED_FOUND" = false ]; then
    echo -e "  ${GREEN}✓ 记忆文件中未发现凭据${NC}"
fi

# ═══════════════════════════════════════════════════════════════
# 检查 2：SECURITY.md 是否存在且为近期
# ═══════════════════════════════════════════════════════════════
echo ""
echo -e "${BOLD}[2/6] 正在检查 SECURITY.md...${NC}"

if [ -f "$WORKSPACE/SECURITY.md" ]; then
    SECURITY_AGE=$(( ($(date +%s) - $(stat -f%m "$WORKSPACE/SECURITY.md" 2>/dev/null || stat -c%Y "$WORKSPACE/SECURITY.md" 2>/dev/null)) / 86400 ))
    if [ "$SECURITY_AGE" -gt 90 ]; then
        echo -e "  ${YELLOW}⚠ SECURITY.md 已有 $SECURITY_AGE 天未更新 — 建议审查${NC}"
        ((WARNINGS_FOUND++))
    else
        echo -e "  ${GREEN}✓ SECURITY.md 存在且为近期更新（$SECURITY_AGE 天前）${NC}"
    fi
else
    echo -e "  ${RED}✗ 未找到 SECURITY.md — 运行设置向导创建${NC}"
    ((ISSUES_FOUND++))
fi

# ═══════════════════════════════════════════════════════════════
# 检查 3：日志中的可疑模式
# ═══════════════════════════════════════════════════════════════
echo ""
echo -e "${BOLD}[3/6] 正在检查日志中的注入尝试...${NC}"

INJECTION_PATTERNS=(
    "ignore previous"
    "ignore all"
    "forget everything"
    "you are now"
    "new instructions"
    "system override"
    "admin mode"
    "SYSTEM MESSAGE"
)

INJECTION_FOUND=false
for pattern in "${INJECTION_PATTERNS[@]}"; do
    if grep -r -i -l "$pattern" "$WORKSPACE/memory" 2>/dev/null | head -3 | grep -q .; then
        if [ "$INJECTION_FOUND" = false ]; then
            echo -e "  ${YELLOW}⚠ 日志中发现潜在注入尝试：${NC}"
            INJECTION_FOUND=true
            ((WARNINGS_FOUND++))
        fi
        COUNT=$(grep -r -i -c "$pattern" "$WORKSPACE/memory" 2>/dev/null | grep -v ":0" | wc -l)
        echo "    - 模式 '$pattern' 出现在 $COUNT 个文件中"
    fi
done

if [ "$INJECTION_FOUND" = false ]; then
    echo -e "  ${GREEN}✓ 日志中未发现注入模式${NC}"
fi

# ═══════════════════════════════════════════════════════════════
# 检查 4：文件权限
# ═══════════════════════════════════════════════════════════════
echo ""
echo -e "${BOLD}[4/6] 正在检查文件权限...${NC}"

# 检查敏感文件是否对所有人可读
WORLD_READABLE=false
for file in "$WORKSPACE/SECURITY.md" "$WORKSPACE/TEAM.md" "$WORKSPACE/USER.md"; do
    if [ -f "$file" ]; then
        PERMS=$(stat -f%Lp "$file" 2>/dev/null || stat -c%a "$file" 2>/dev/null)
        if [[ "$PERMS" =~ [0-7][0-7][4-7] ]]; then
            if [ "$WORLD_READABLE" = false ]; then
                echo -e "  ${YELLOW}⚠ 部分文件可能对所有人可读：${NC}"
                WORLD_READABLE=true
                ((WARNINGS_FOUND++))
            fi
            echo "    - $file（权限：$PERMS）"
        fi
    fi
done

if [ "$WORLD_READABLE" = false ]; then
    echo -e "  ${GREEN}✓ 文件权限看起来合理${NC}"
fi

# ═══════════════════════════════════════════════════════════════
# 检查 5：备份状态
# ═══════════════════════════════════════════════════════════════
echo ""
echo -e "${BOLD}[5/6] 正在检查备份状态...${NC}"

if [ -d "$WORKSPACE/backups" ]; then
    BACKUP_COUNT=$(find "$WORKSPACE/backups" -type f -mtime -30 2>/dev/null | wc -l)
    if [ "$BACKUP_COUNT" -gt 0 ]; then
        echo -e "  ${GREEN}✓ 发现近 30 天内的 $BACKUP_COUNT 份备份文件${NC}"
    else
        echo -e "  ${YELLOW}⚠ 近 30 天内未发现备份${NC}"
        ((WARNINGS_FOUND++))
    fi
else
    echo -e "  ${YELLOW}⚠ 未找到备份目录${NC}"
    ((WARNINGS_FOUND++))
fi

# ═══════════════════════════════════════════════════════════════
# 检查 6：核心文件完整性
# ═══════════════════════════════════════════════════════════════
echo ""
echo -e "${BOLD}[6/6] 正在检查核心文件完整性...${NC}"

CORE_FILES=("SOUL.md" "USER.md" "MEMORY.md" "AGENTS.md" "SECURITY.md")
MISSING_CORE=false

for file in "${CORE_FILES[@]}"; do
    if [ ! -f "$WORKSPACE/$file" ]; then
        if [ "$MISSING_CORE" = false ]; then
            echo -e "  ${YELLOW}⚠ 缺失核心文件：${NC}"
            MISSING_CORE=true
            ((WARNINGS_FOUND++))
        fi
        echo "    - $file"
    fi
done

if [ "$MISSING_CORE" = false ]; then
    echo -e "  ${GREEN}✓ 所有核心文件齐全${NC}"
fi

# ═══════════════════════════════════════════════════════════════
# 摘要
# ═══════════════════════════════════════════════════════════════
echo ""
echo -e "${BOLD}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${BOLD}  安全审计摘要${NC}"
echo -e "${BOLD}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""

if [ "$ISSUES_FOUND" -eq 0 ] && [ "$WARNINGS_FOUND" -eq 0 ]; then
    echo -e "  ${GREEN}${BOLD}✅ 一切正常！未发现安全问题。${NC}"
else
    if [ "$ISSUES_FOUND" -gt 0 ]; then
        echo -e "  ${RED}${BOLD}❌ 发现问题：$ISSUES_FOUND 个${NC}"
    fi
    if [ "$WARNINGS_FOUND" -gt 0 ]; then
        echo -e "  ${YELLOW}${BOLD}⚠️  警告：$WARNINGS_FOUND 个${NC}"
    fi
    echo ""
    echo "  建议操作："
    if [ "$ISSUES_FOUND" -gt 0 ]; then
        echo "  • 审查并移除任何暴露的凭据"
        echo "  • 如核心文件缺失则运行设置向导"
    fi
    if [ "$WARNINGS_FOUND" -gt 0 ]; then
        echo "  • 审查日志中的注入尝试"
        echo "  • 如 SECURITY.md 过期则更新"
        echo "  • 如缺少备份则创建"
    fi
fi

echo ""
echo -e "${BOLD}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""
echo "  每月运行此审计：./scripts/security-audit.sh"
echo ""
echo "  AI Persona OS 作者：Jeff J Hunter"
echo "  https://jeffjhunter.com"
echo ""
