#!/usr/bin/env bash
#
# md2wechat OpenClaw Skill Installer
#
# Usage:
#   curl -fsSL https://github.com/geekjourneyx/md2wechat-skill/releases/download/vX.Y.Z/install-openclaw.sh | bash
#

set -euo pipefail

REPO="geekjourneyx/md2wechat-skill"
VERSION="${MD2WECHAT_VERSION:-}"
if [[ -z "$VERSION" ]]; then
    if [[ -n "${MD2WECHAT_VERSION_DEFAULT:-}" ]]; then
        VERSION="${MD2WECHAT_VERSION_DEFAULT}"
    fi
fi
RELEASE_BASE_URL="${MD2WECHAT_RELEASE_BASE_URL:-}"
SKILL_NAME="md2wechat"
SKILL_ARCHIVE="md2wechat-openclaw-skill.tar.gz"
INSTALL_DIR="${MD2WECHAT_OPENCLAW_INSTALL_DIR:-${HOME}/.openclaw/skills/${SKILL_NAME}}"
CLI_INSTALL_DIR="${MD2WECHAT_INSTALL_DIR:-${MD2WECHAT_OPENCLAW_RUNTIME_DIR:-${HOME}/.local/bin}}"
NON_INTERACTIVE="${MD2WECHAT_NONINTERACTIVE:-}"

if [[ -z "$RELEASE_BASE_URL" ]]; then
    [[ -n "$VERSION" ]] || error "必须提供 MD2WECHAT_VERSION 或使用固定版本 release 安装器 / MD2WECHAT_VERSION is required unless a fixed-version installer injected it"
    RELEASE_BASE_URL="https://github.com/${REPO}/releases/download/v${VERSION}"
fi

CHECKSUMS_URL="${RELEASE_BASE_URL}/checksums.txt"
ARCHIVE_URL="${RELEASE_BASE_URL}/${SKILL_ARCHIVE}"

detect_binary() {
    local os arch
    os="$(uname -s)"
    arch="$(uname -m)"

    if [[ "$os" == "Darwin" ]]; then
        if [[ "$arch" == "arm64" ]]; then
            printf 'md2wechat-darwin-arm64\n'
        else
            printf 'md2wechat-darwin-amd64\n'
        fi
        return 0
    fi

    if [[ "$os" == "Linux" ]]; then
        if [[ "$arch" == "aarch64" || "$arch" == "arm64" ]]; then
            printf 'md2wechat-linux-arm64\n'
        else
            printf 'md2wechat-linux-amd64\n'
        fi
        return 0
    fi

    error "不支持的系统 / Unsupported platform: ${os} ${arch}"
}

RUNTIME_BINARY="$(detect_binary)"
RUNTIME_URL="${RELEASE_BASE_URL}/${RUNTIME_BINARY}"

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
NC='\033[0m'

info()    { printf "${BLUE}ℹ${NC} %s\n" "$1"; }
success() { printf "${GREEN}✓${NC} %s\n" "$1"; }
warn()    { printf "${YELLOW}⚠${NC} %s\n" "$1"; }
error()   { printf "${RED}✗${NC} %s\n" "$1" >&2; exit 1; }

confirm_or_continue() {
    prompt="$1"
    if [[ -n "$NON_INTERACTIVE" || -n "${CI:-}" ]]; then
        return 0
    fi
    read -p "$prompt [y/N] " -n 1 -r
    printf "\n"
    [[ $REPLY =~ ^[Yy]$ ]]
}

download_file() {
    url="$1"
    output="$2"
    case "$url" in
        file://*)
            cp "${url#file://}" "$output"
            ;;
        *)
            if command -v curl >/dev/null 2>&1; then
                curl -fsSL "$url" -o "$output"
            elif command -v wget >/dev/null 2>&1; then
                wget -q "$url" -O "$output"
            else
                error "需要 curl 或 wget / Need curl or wget"
            fi
            ;;
    esac
}

verify_checksum() {
    checksums_file="$1"
    archive_file="$2"
    archive_name="$3"

    if command -v sha256sum >/dev/null 2>&1; then
        (cd "$(dirname "$archive_file")" && sha256sum -c "$checksums_file" --ignore-missing --status)
    elif command -v shasum >/dev/null 2>&1; then
        expected="$(awk -v file="$archive_name" '$2 == file { print $1 }' "$checksums_file")"
        [[ -n "$expected" ]] || error "checksums.txt 中未找到 ${archive_name} 的校验值"
        actual="$(shasum -a 256 "$archive_file" | awk '{print $1}')"
        [[ "$expected" == "$actual" ]]
    else
        error "需要 sha256sum 或 shasum 来验证安装包完整性"
    fi
}

print_banner() {
cat <<'EOF'
+--------------------------------------------------+
|               md2wechat OpenClaw                 |
|                    installer                     |
|              crafted by geekjourneyx             |
+--------------------------------------------------+
EOF
}

printf "\n"
print_banner
printf "\n"

if command -v clawhub >/dev/null 2>&1; then
    info "检测到 clawhub CLI / ClawHub CLI detected"
    printf "\n"
    printf "推荐使用 ClawHub 安装 / Recommend using ClawHub:\n"
    printf "  ${GREEN}clawhub install md2wechat${NC}\n"
    printf "\n"
    if ! confirm_or_continue "继续手动安装？/ Continue manual install?"; then
        exit 0
    fi
fi

if [[ ! -d "${HOME}/.openclaw" ]]; then
    warn "未检测到 OpenClaw 安装 / OpenClaw not detected"
    info "请先安装 OpenClaw: https://openclaw.ai/"
    info "Install OpenClaw first: https://openclaw.ai/"
    printf "\n"
    if ! confirm_or_continue "仍要继续安装技能？/ Continue installing skill anyway?"; then
        exit 0
    fi
fi

if [[ -d "$INSTALL_DIR" ]]; then
    warn "已存在安装 / Existing installation: $INSTALL_DIR"
    if ! confirm_or_continue "覆盖？/ Overwrite?"; then
        exit 0
    fi
    rm -rf "$INSTALL_DIR"
fi

if [[ -x "${CLI_INSTALL_DIR}/md2wechat" ]]; then
    warn "已存在 CLI / Existing CLI: ${CLI_INSTALL_DIR}/md2wechat"
    if ! confirm_or_continue "覆盖 CLI？/ Overwrite CLI?"; then
        exit 0
    fi
fi

TMP_DIR="$(mktemp -d)"
cleanup() {
    rm -rf "$TMP_DIR"
}
trap cleanup EXIT

ARCHIVE_PATH="${TMP_DIR}/${SKILL_ARCHIVE}"
RUNTIME_PATH="${TMP_DIR}/${RUNTIME_BINARY}"
CHECKSUMS_FILE="${TMP_DIR}/checksums.txt"

info "下载技能文件 / Downloading release assets..."
info "技能包 / Skill bundle: ${ARCHIVE_URL}"
info "CLI binary: ${RUNTIME_URL}"
info "校验文件 / Checksums: ${CHECKSUMS_URL}"

download_file "$ARCHIVE_URL" "$ARCHIVE_PATH"
download_file "$RUNTIME_URL" "$RUNTIME_PATH"
download_file "$CHECKSUMS_URL" "$CHECKSUMS_FILE"

info "正在验证 SHA-256 校验值 / Verifying SHA-256 checksum..."
if ! verify_checksum "$CHECKSUMS_FILE" "$ARCHIVE_PATH" "$SKILL_ARCHIVE"; then
    error "校验失败：下载文件与发布校验值不匹配 / Checksum verification failed"
fi
if ! verify_checksum "$CHECKSUMS_FILE" "$RUNTIME_PATH" "$RUNTIME_BINARY"; then
    error "CLI 校验失败：下载文件与发布校验值不匹配 / CLI checksum verification failed"
fi

mkdir -p "$INSTALL_DIR"
mkdir -p "$CLI_INSTALL_DIR"
tar -xzf "$ARCHIVE_PATH" -C "$TMP_DIR"

EXTRACTED_DIR="${TMP_DIR}/skills/md2wechat"
[[ -d "$EXTRACTED_DIR" ]] || error "技能包结构无效 / Invalid skill bundle layout"

cp -r "${EXTRACTED_DIR}/"* "$INSTALL_DIR/"
chmod +x "${INSTALL_DIR}/scripts/"*.sh 2>/dev/null || true
install -m 0755 "$RUNTIME_PATH" "${CLI_INSTALL_DIR}/md2wechat"

success "安装完成 / Installation complete!"

printf "\n"
printf "${BLUE}========================================${NC}\n"
printf "${BLUE}   下一步 / Next Steps${NC}\n"
printf "${BLUE}========================================${NC}\n"
printf "\n"
printf "${YELLOW}注意 / Note:${NC}\n"
printf "  • WECHAT_APPID/SECRET 仅草稿上传需要，预览转换可不配置\n"
printf "  • 图片生成需额外配置 IMAGE_API_KEY\n"
printf "  • OpenClaw 安装器会一并安装 skill 壳并把 md2wechat CLI 安装到正式环境路径\n"
printf "  • 推荐始终使用固定版本 release 资产，不要使用 main/raw 作为安装入口\n"
printf "\n"
printf "推荐执行 / Recommended commands:\n"
printf "  ${GREEN}%s/md2wechat config init${NC}\n" "${CLI_INSTALL_DIR}"
printf "  ${GREEN}%s/md2wechat config validate${NC}\n" "${CLI_INSTALL_DIR}"
printf "\n"
printf "默认配置文件 / Default config file:\n"
printf "  ${GREEN}~/.config/md2wechat/config.yaml${NC}\n"
printf "\n"
printf "安装路径 / Installed to: ${GREEN}%s${NC}\n" "$INSTALL_DIR"
printf "CLI installed to: ${GREEN}%s${NC}\n" "${CLI_INSTALL_DIR}/md2wechat"
if [[ ":$PATH:" != *":${CLI_INSTALL_DIR}:"* ]]; then
    printf "\n"
    printf "如果当前 shell 还找不到命令 / If the current shell still cannot find md2wechat:\n"
    printf "  ${GREEN}export PATH=\"%s:\$PATH\"${NC}\n" "${CLI_INSTALL_DIR}"
    printf "  ${GREEN}md2wechat version --json${NC}\n"
fi
printf "文档 / Documentation: https://github.com/${REPO}/blob/main/docs/OPENCLAW.md\n"
printf "OpenClaw 官网 / OpenClaw: https://openclaw.ai/\n"
printf "ClawHub 技能市场 / ClawHub: https://clawhub.ai/\n"
printf "\n"
