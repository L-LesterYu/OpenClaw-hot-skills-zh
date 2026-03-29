# md2wechat Windows 自动安装脚本
# 使用方法：在 PowerShell 中运行
# Set-ExecutionPolicy Bypass -Scope Process -Force; [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072; iex ((New-Object System.Net.WebClient).DownloadString('https://github.com/geekjourneyx/md2wechat-skill/releases/download/vX.Y.Z/install.ps1'))

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "   md2wechat 安装向导" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

$repo = "geekjourneyx/md2wechat-skill"
$version = if ($env:MD2WECHAT_VERSION) {
    $env:MD2WECHAT_VERSION
} elseif ($env:MD2WECHAT_VERSION_DEFAULT) {
    $env:MD2WECHAT_VERSION_DEFAULT
} else {
    "latest"
}
$releaseBaseUrl = $env:MD2WECHAT_RELEASE_BASE_URL
if (-not $releaseBaseUrl) {
    if ($version -eq "latest") {
        $releaseBaseUrl = "https://github.com/$repo/releases/latest/download"
    } else {
        $releaseBaseUrl = "https://github.com/$repo/releases/download/v$version"
    }
}
$installDirOverride = $env:MD2WECHAT_INSTALL_DIR
$nonInteractive = [bool]($env:MD2WECHAT_NONINTERACTIVE -or $env:CI)
$skipPathUpdate = [bool]($env:MD2WECHAT_NO_PATH_UPDATE -or $env:CI)

# 检测是否是管理员
$isAdmin = ([Security.Principal.WindowsPrincipal] [Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)

# 确定安装目录
if ($installDirOverride) {
    $installDir = $installDirOverride
} elseif ($isAdmin) {
    $installDir = "C:\Program Files\md2wechat"
} else {
    $installDir = "$env:USERPROFILE\AppData\Local\md2wechat"
}

Write-Host "安装目录: $installDir" -ForegroundColor Yellow
Write-Host ""

# 创建目录
New-Item -ItemType Directory -Force -Path $installDir | Out-Null
if (Test-Path (Join-Path $installDir "md2wechat.exe")) {
    Write-Host "检测到已存在安装，将覆盖: $(Join-Path $installDir 'md2wechat.exe')" -ForegroundColor Yellow
}

# 下载
Write-Host "正在下载..." -ForegroundColor Green
$binaryName = "md2wechat-windows-amd64.exe"
$downloadUrl = "$releaseBaseUrl/$binaryName"
$checksumsUrl = "$releaseBaseUrl/checksums.txt"
$outputFile = "$installDir\md2wechat.exe"
$tempDir = Join-Path ([System.IO.Path]::GetTempPath()) ("md2wechat-install-" + [System.Guid]::NewGuid().ToString("N"))
$downloadedBinary = Join-Path $tempDir $binaryName
$checksumsFile = Join-Path $tempDir "checksums.txt"

New-Item -ItemType Directory -Force -Path $tempDir | Out-Null

try {
    function Download-File {
        param(
            [Parameter(Mandatory = $true)][string]$Uri,
            [Parameter(Mandatory = $true)][string]$OutFile
        )

        if ($Uri.StartsWith("file://")) {
            $sourcePath = ([Uri]$Uri).LocalPath
            Copy-Item -Force $sourcePath $OutFile
            return
        }

        Invoke-WebRequest -Uri $Uri -OutFile $OutFile -UseBasicParsing
    }

    [Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12
    Download-File -Uri $downloadUrl -OutFile $downloadedBinary
    Download-File -Uri $checksumsUrl -OutFile $checksumsFile

    Write-Host "正在验证 SHA-256 校验值..." -ForegroundColor Yellow
    $expectedLine = Select-String -Path $checksumsFile -Pattern (" " + [regex]::Escape($binaryName) + "$") | Select-Object -First 1
    if (-not $expectedLine) {
        throw "checksums.txt 中未找到 $binaryName 的校验值"
    }
    $expectedHash = ($expectedLine.Line -split '\s+')[0].ToLowerInvariant()
    $actualHash = (Get-FileHash -Path $downloadedBinary -Algorithm SHA256).Hash.ToLowerInvariant()
    if ($expectedHash -ne $actualHash) {
        throw "SHA-256 校验失败"
    }

    Move-Item -Force $downloadedBinary $outputFile
    Write-Host "✅ 下载完成！" -ForegroundColor Green
} catch {
    Write-Host "❌ 下载失败: $_" -ForegroundColor Red
    if (Test-Path $downloadedBinary) { Remove-Item -Force $downloadedBinary }
    if (-not $nonInteractive) {
        Read-Host "按回车键退出"
    }
    exit 1
} finally {
    if (Test-Path $checksumsFile) { Remove-Item -Force $checksumsFile }
    if (Test-Path $tempDir) { Remove-Item -Force -Recurse $tempDir }
}

Write-Host ""

$verifyCmd = "& `"$outputFile`" version --json"
$initCmd = "& `"$outputFile`" config init"

# 添加到 PATH
if ($skipPathUpdate) {
    Write-Host "ℹ️  跳过 PATH 更新（CI / non-interactive 模式）" -ForegroundColor Yellow
} else {
    $currentPath = [Environment]::GetEnvironmentVariable("Path", "User")
    if ($currentPath -notlike "*$installDir*") {
        Write-Host "添加到系统 PATH..." -ForegroundColor Yellow
        [Environment]::SetEnvironmentVariable("Path", "$currentPath;$installDir", "User")
        if ($env:Path -notlike "*$installDir*") {
            $env:Path = "$installDir;$env:Path"
        }
        Write-Host "✅ 已添加到 PATH" -ForegroundColor Green
        Write-Host ""
        Write-Host "当前 PowerShell 会话已可直接使用 md2wechat" -ForegroundColor Green
        Write-Host "新的终端或命令提示符也会继承该 PATH 更新" -ForegroundColor Yellow
    } else {
        if ($env:Path -notlike "*$installDir*") {
            $env:Path = "$installDir;$env:Path"
        }
        Write-Host "✅ 已在 PATH 中" -ForegroundColor Green
    }
}

Write-Host ""
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "   安装完成！" -ForegroundColor Green
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "下一步：" -ForegroundColor Yellow
Write-Host "  1. 运行: md2wechat version --json" -ForegroundColor White
Write-Host "  2. 运行: md2wechat config init" -ForegroundColor White
Write-Host "  3. 编辑生成的配置文件" -ForegroundColor White
Write-Host "  4. 运行: md2wechat inspect 文章.md" -ForegroundColor White
Write-Host "  5. 运行: md2wechat preview 文章.md" -ForegroundColor White
Write-Host ""
Write-Host "如果当前会话里仍然找不到命令，直接运行：" -ForegroundColor Yellow
Write-Host "  $verifyCmd" -ForegroundColor White
Write-Host "  $initCmd" -ForegroundColor White
Write-Host ""
Write-Host "查看帮助: & `"$outputFile`" --help" -ForegroundColor White
Write-Host ""

if (-not $nonInteractive) {
    Read-Host "按回车键退出"
}
