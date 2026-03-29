#!/usr/bin/env bash
set -euo pipefail

usage() {
  cat <<'EOF'
Usage:
  generate-homebrew-formula.sh \
    --version <semver> \
    --base-url <release base url> \
    --darwin-arm64-sha <sha256> \
    --darwin-x86-64-sha <sha256> \
    --linux-arm64-sha <sha256> \
    --linux-x86-64-sha <sha256> \
    [--output <path>]
EOF
}

version=""
base_url=""
darwin_arm64_sha=""
darwin_x86_64_sha=""
linux_arm64_sha=""
linux_x86_64_sha=""
output=""

while [[ $# -gt 0 ]]; do
  case "$1" in
    --version)
      version="${2:-}"
      shift 2
      ;;
    --base-url)
      base_url="${2:-}"
      shift 2
      ;;
    --darwin-arm64-sha)
      darwin_arm64_sha="${2:-}"
      shift 2
      ;;
    --darwin-x86-64-sha)
      darwin_x86_64_sha="${2:-}"
      shift 2
      ;;
    --linux-arm64-sha)
      linux_arm64_sha="${2:-}"
      shift 2
      ;;
    --linux-x86-64-sha)
      linux_x86_64_sha="${2:-}"
      shift 2
      ;;
    --output)
      output="${2:-}"
      shift 2
      ;;
    -h|--help)
      usage
      exit 0
      ;;
    *)
      echo "unknown argument: $1" >&2
      usage >&2
      exit 1
      ;;
  esac
done

[[ -n "$version" ]] || { echo "--version is required" >&2; exit 1; }
[[ -n "$base_url" ]] || { echo "--base-url is required" >&2; exit 1; }
[[ -n "$darwin_arm64_sha" ]] || { echo "--darwin-arm64-sha is required" >&2; exit 1; }
[[ -n "$darwin_x86_64_sha" ]] || { echo "--darwin-x86-64-sha is required" >&2; exit 1; }
[[ -n "$linux_arm64_sha" ]] || { echo "--linux-arm64-sha is required" >&2; exit 1; }
[[ -n "$linux_x86_64_sha" ]] || { echo "--linux-x86-64-sha is required" >&2; exit 1; }

formula=$(cat <<EOF
class Md2wechat < Formula
  desc "Convert Markdown to WeChat Official Account HTML"
  homepage "https://github.com/geekjourneyx/md2wechat-skill"
  license "MIT"
  version "${version}"

  on_macos do
    on_arm do
      url "${base_url}/md2wechat_Darwin_arm64.tar.gz"
      sha256 "${darwin_arm64_sha}"
    end
    on_intel do
      url "${base_url}/md2wechat_Darwin_x86_64.tar.gz"
      sha256 "${darwin_x86_64_sha}"
    end
  end

  on_linux do
    on_arm do
      url "${base_url}/md2wechat_Linux_arm64.tar.gz"
      sha256 "${linux_arm64_sha}"
    end
    on_intel do
      url "${base_url}/md2wechat_Linux_x86_64.tar.gz"
      sha256 "${linux_x86_64_sha}"
    end
  end

  def install
    bin.install "md2wechat"
  end

  test do
    require "json"

    payload = JSON.parse(shell_output("#{bin}/md2wechat version --json"))
    assert_equal version.to_s, payload.fetch("data").fetch("version")
  end
end
EOF
)

if [[ -n "$output" ]]; then
  printf '%s\n' "$formula" >"$output"
else
  printf '%s\n' "$formula"
fi
