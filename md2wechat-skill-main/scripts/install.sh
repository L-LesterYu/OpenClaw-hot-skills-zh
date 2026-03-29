#!/bin/bash
# md2wechat 自动安装脚本
# 适用于：macOS / Linux
# 使用方法：curl -fsSL https://github.com/geekjourneyx/md2wechat-skill/releases/download/vX.Y.Z/install.sh | bash

set -euo pipefail

REPO="geekjourneyx/md2wechat-skill"
VERSION="${MD2WECHAT_VERSION:-}"
if [ -z "$VERSION" ]; then
    if [ -n "${MD2WECHAT_VERSION_DEFAULT:-}" ]; then
        VERSION="${MD2WECHAT_VERSION_DEFAULT}"
    else
        VERSION="latest"
    fi
fi
INSTALL_DIR="${MD2WECHAT_INSTALL_DIR:-$HOME/.local/bin}"

choose_profile_file() {
    if [[ -n "${ZSH_VERSION:-}" ]]; then
        printf '%s\n' "$HOME/.zshrc"
        return 0
    fi

    if [[ -n "${BASH_VERSION:-}" ]]; then
        if [[ -f "$HOME/.bash_profile" ]]; then
            printf '%s\n' "$HOME/.bash_profile"
        else
            printf '%s\n' "$HOME/.bashrc"
        fi
        return 0
    fi

    case "$(basename "${SHELL:-}")" in
        zsh) printf '%s\n' "$HOME/.zshrc" ;;
        bash)
            if [[ -f "$HOME/.bash_profile" ]]; then
                printf '%s\n' "$HOME/.bash_profile"
            else
                printf '%s\n' "$HOME/.bashrc"
            fi
            ;;
        *) printf '%s\n' "" ;;
    esac
}

print_banner() {
cat <<'EOF'
+--------------------------------------------------+
|                    md2wechat                     |
|                    installer                     |
|              crafted by geekjourneyx             |
+--------------------------------------------------+
EOF
}

print_banner
echo ""

# 检测系统
OS="$(uname -s)"
ARCH="$(uname -m)"

echo "检测到系统: $OS $ARCH"

# 确定下载链接
if [ "$OS" = "Darwin" ]; then
    if [ "$ARCH" = "arm64" ]; then
        BINARY="md2wechat-darwin-arm64"
    else
        BINARY="md2wechat-darwin-amd64"
    fi
elif [ "$OS" = "Linux" ]; then
    if [ "$ARCH" = "aarch64" ] || [ "$ARCH" = "arm64" ]; then
        BINARY="md2wechat-linux-arm64"
    else
        BINARY="md2wechat-linux-amd64"
    fi
else
    echo "❌ 不支持的系统: $OS"
    exit 1
fi

echo "将下载: $BINARY"
echo ""

# 确定安装目录
mkdir -p "$INSTALL_DIR"
if [[ -e "$INSTALL_DIR/md2wechat" ]]; then
    echo "检测到已存在安装，将覆盖: $INSTALL_DIR/md2wechat"
fi

TMP_DIR="$(mktemp -d)"
cleanup() {
    rm -rf "$TMP_DIR"
}
trap cleanup EXIT

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
                echo "❌ 需要 curl 或 wget 来下载文件"
                exit 1
            fi
            ;;
    esac
}

verify_checksum() {
    checksums_file="$1"
    binary_file="$2"
    binary_name="$3"

    if command -v sha256sum >/dev/null 2>&1; then
        (cd "$(dirname "$binary_file")" && sha256sum -c "$checksums_file" --ignore-missing --status)
    elif command -v shasum >/dev/null 2>&1; then
        expected="$(awk -v file="$binary_name" '$2 == file { print $1 }' "$checksums_file")"
        if [ -z "$expected" ]; then
            echo "❌ checksums.txt 中未找到 $binary_name 的校验值"
            exit 1
        fi
        actual="$(shasum -a 256 "$binary_file" | awk '{print $1}')"
        [ "$expected" = "$actual" ]
    else
        echo "❌ 需要 sha256sum 或 shasum 来验证安装包完整性"
        exit 1
    fi
}

# 下载
echo "正在下载..."
RELEASE_BASE_URL="${MD2WECHAT_RELEASE_BASE_URL:-}"
if [ -z "$RELEASE_BASE_URL" ]; then
    if [ "$VERSION" = "latest" ]; then
        RELEASE_BASE_URL="https://github.com/${REPO}/releases/latest/download"
    else
        RELEASE_BASE_URL="https://github.com/${REPO}/releases/download/v${VERSION}"
    fi
fi
DOWNLOAD_URL="${RELEASE_BASE_URL}/${BINARY}"
CHECKSUMS_URL="${RELEASE_BASE_URL}/checksums.txt"
echo "下载地址: $DOWNLOAD_URL"
echo "校验文件: $CHECKSUMS_URL"

DOWNLOADED_BINARY="$TMP_DIR/$BINARY"
CHECKSUMS_FILE="$TMP_DIR/checksums.txt"

download_file "$DOWNLOAD_URL" "$DOWNLOADED_BINARY"
download_file "$CHECKSUMS_URL" "$CHECKSUMS_FILE"

echo "正在验证 SHA-256 校验值..."
if ! verify_checksum "$CHECKSUMS_FILE" "$DOWNLOADED_BINARY" "$BINARY"; then
    echo "❌ 校验失败：下载文件与发布校验值不匹配"
    exit 1
fi

install -m 0755 "$DOWNLOADED_BINARY" "$INSTALL_DIR/md2wechat"

# 添加执行权限
chmod +x "$INSTALL_DIR/md2wechat"

PROFILE_FILE="$(choose_profile_file)"
CURRENT_SESSION_CMD="export PATH=\"$INSTALL_DIR:\$PATH\""
VERIFY_CMD="$INSTALL_DIR/md2wechat version --json"

echo ""
echo "✅ 下载完成！"
echo ""

# 检查 PATH
if echo ":$PATH:" | grep -q ":$INSTALL_DIR:"; then
    echo "✅ 安装目录已在 PATH 中"
else
    echo "⚠️  需要将安装目录添加到 PATH"
    echo ""
    echo "当前会话立即生效："
    echo "  ${CURRENT_SESSION_CMD}"
    echo "  md2wechat version --json"
    echo ""
    if [ -n "$PROFILE_FILE" ]; then
        echo "想让以后每次终端都生效："
        echo "  echo 'export PATH=\"$INSTALL_DIR:\$PATH\"' >> $PROFILE_FILE"
        echo "  source $PROFILE_FILE"
    else
        echo "请将 $INSTALL_DIR 添加到你的 PATH 环境变量。"
    fi
fi

echo ""
echo "========================================"
echo "   安装完成！"
echo "========================================"
echo ""
echo "下一步："
if echo ":$PATH:" | grep -q ":$INSTALL_DIR:"; then
    echo "  1. 运行: md2wechat version --json"
    echo "  2. 运行: md2wechat config init"
    echo "  3. 编辑生成的配置文件"
    echo "  4. 运行: md2wechat inspect 文章.md"
    echo "  5. 运行: md2wechat preview 文章.md"
else
    echo "  1. 运行: ${CURRENT_SESSION_CMD}"
    echo "  2. 运行: md2wechat version --json"
    echo "  3. 运行: md2wechat config init"
    echo "  4. 编辑生成的配置文件"
    echo "  5. 运行: md2wechat inspect 文章.md"
    echo "  6. 运行: md2wechat preview 文章.md"
fi
echo ""
echo "如果当前会话仍然找不到命令，直接运行:"
echo "  ${VERIFY_CMD}"
echo "  ${INSTALL_DIR}/md2wechat config init"
echo ""
echo "查看帮助: ${INSTALL_DIR}/md2wechat --help"
echo ""
