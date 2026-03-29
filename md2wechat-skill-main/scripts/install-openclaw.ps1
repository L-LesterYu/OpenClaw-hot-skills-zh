# md2wechat OpenClaw Windows installer
# Usage:
#   $env:MD2WECHAT_RELEASE_BASE_URL = "https://github.com/geekjourneyx/md2wechat-skill/releases/download/vX.Y.Z"
#   iex ((New-Object System.Net.WebClient).DownloadString("$env:MD2WECHAT_RELEASE_BASE_URL/install-openclaw.ps1"))

Write-Host "+--------------------------------------------------+" -ForegroundColor Cyan
Write-Host "|               md2wechat OpenClaw                 |" -ForegroundColor Cyan
Write-Host "|                    installer                     |" -ForegroundColor Cyan
Write-Host "|              crafted by geekjourneyx             |" -ForegroundColor Cyan
Write-Host "+--------------------------------------------------+" -ForegroundColor Cyan
Write-Host ""

$repo = "geekjourneyx/md2wechat-skill"
$version = if ($env:MD2WECHAT_VERSION) {
    $env:MD2WECHAT_VERSION
} elseif ($env:MD2WECHAT_VERSION_DEFAULT) {
    $env:MD2WECHAT_VERSION_DEFAULT
} else {
    throw "MD2WECHAT_VERSION is required unless injected by a fixed-version installer."
}

$releaseBaseUrl = $env:MD2WECHAT_RELEASE_BASE_URL
if (-not $releaseBaseUrl) {
    $releaseBaseUrl = "https://github.com/$repo/releases/download/v$version"
}

$skillDir = if ($env:MD2WECHAT_OPENCLAW_INSTALL_DIR) {
    $env:MD2WECHAT_OPENCLAW_INSTALL_DIR
} else {
    Join-Path $env:USERPROFILE ".openclaw\skills\md2wechat"
}

$installDir = if ($env:MD2WECHAT_INSTALL_DIR) {
    $env:MD2WECHAT_INSTALL_DIR
} elseif ([bool]($env:MD2WECHAT_NONINTERACTIVE -or $env:CI)) {
    Join-Path $env:USERPROFILE ".local\bin"
} else {
    Join-Path $env:USERPROFILE "AppData\Local\md2wechat"
}

$nonInteractive = [bool]($env:MD2WECHAT_NONINTERACTIVE -or $env:CI)
$skipPathUpdate = [bool]($env:MD2WECHAT_NO_PATH_UPDATE -or $env:CI)

$skillArchive = "md2wechat-openclaw-skill.tar.gz"
$binaryName = "md2wechat-windows-amd64.exe"
$archiveUrl = "$releaseBaseUrl/$skillArchive"
$binaryUrl = "$releaseBaseUrl/$binaryName"
$checksumsUrl = "$releaseBaseUrl/checksums.txt"

Write-Host "Skill dir: $skillDir" -ForegroundColor Yellow
Write-Host "CLI dir:   $installDir" -ForegroundColor Yellow
Write-Host ""

New-Item -ItemType Directory -Force -Path $skillDir | Out-Null
New-Item -ItemType Directory -Force -Path $installDir | Out-Null

$tempDir = Join-Path ([System.IO.Path]::GetTempPath()) ("md2wechat-openclaw-" + [System.Guid]::NewGuid().ToString("N"))
New-Item -ItemType Directory -Force -Path $tempDir | Out-Null

$archivePath = Join-Path $tempDir $skillArchive
$binaryPath = Join-Path $tempDir $binaryName
$checksumsPath = Join-Path $tempDir "checksums.txt"
$extractRoot = Join-Path $tempDir "extract"
$outputFile = Join-Path $installDir "md2wechat.exe"

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

    Write-Host "Downloading release assets..." -ForegroundColor Green
    Download-File -Uri $archiveUrl -OutFile $archivePath
    Download-File -Uri $binaryUrl -OutFile $binaryPath
    Download-File -Uri $checksumsUrl -OutFile $checksumsPath

    Write-Host "Verifying SHA-256 checksums..." -ForegroundColor Yellow
    foreach ($name in @($skillArchive, $binaryName)) {
        $expectedLine = Select-String -Path $checksumsPath -Pattern (" " + [regex]::Escape($name) + "$") | Select-Object -First 1
        if (-not $expectedLine) {
            throw "checksums.txt 中未找到 $name 的校验值"
        }
        $expectedHash = ($expectedLine.Line -split '\s+')[0].ToLowerInvariant()
        $actualPath = if ($name -eq $skillArchive) { $archivePath } else { $binaryPath }
        $actualHash = (Get-FileHash -Path $actualPath -Algorithm SHA256).Hash.ToLowerInvariant()
        if ($expectedHash -ne $actualHash) {
            throw "$name SHA-256 校验失败"
        }
    }

    if (Test-Path $skillDir) {
        Remove-Item -Force -Recurse $skillDir
    }
    New-Item -ItemType Directory -Force -Path $extractRoot | Out-Null
    tar -xzf $archivePath -C $extractRoot

    $extractedSkillDir = Join-Path $extractRoot "skills\md2wechat"
    if (-not (Test-Path $extractedSkillDir)) {
        throw "OpenClaw skill bundle layout is invalid"
    }

    New-Item -ItemType Directory -Force -Path (Split-Path -Parent $skillDir) | Out-Null
    Copy-Item -Force -Recurse $extractedSkillDir $skillDir
    Move-Item -Force $binaryPath $outputFile

    Write-Host "✅ OpenClaw skill and CLI installed." -ForegroundColor Green
} catch {
    Write-Host "❌ Installation failed: $_" -ForegroundColor Red
    if (-not $nonInteractive) {
        Read-Host "按回车键退出"
    }
    exit 1
} finally {
    if (Test-Path $checksumsPath) { Remove-Item -Force $checksumsPath }
    if (Test-Path $archivePath) { Remove-Item -Force $archivePath }
    if (Test-Path $tempDir) { Remove-Item -Force -Recurse $tempDir }
}

$verifyCmd = "& `"$outputFile`" version --json"
$initCmd = "& `"$outputFile`" config init"
$capCmd = "& `"$outputFile`" capabilities --json"

if ($skipPathUpdate) {
    Write-Host "ℹ️  跳过 PATH 更新（CI / non-interactive 模式）" -ForegroundColor Yellow
} else {
    $currentPath = [Environment]::GetEnvironmentVariable("Path", "User")
    if ($currentPath -notlike "*$installDir*") {
        [Environment]::SetEnvironmentVariable("Path", "$currentPath;$installDir", "User")
    }
    if ($env:Path -notlike "*$installDir*") {
        $env:Path = "$installDir;$env:Path"
    }
}

Write-Host ""
Write-Host "Next steps:" -ForegroundColor Yellow
Write-Host "  1. md2wechat version --json" -ForegroundColor White
Write-Host "  2. md2wechat config init" -ForegroundColor White
Write-Host "  3. md2wechat capabilities --json" -ForegroundColor White
Write-Host ""
Write-Host "Skill installed to: $skillDir" -ForegroundColor White
Write-Host "CLI installed to:   $outputFile" -ForegroundColor White
Write-Host ""
Write-Host "If the current shell still cannot find md2wechat, run:" -ForegroundColor Yellow
Write-Host "  $verifyCmd" -ForegroundColor White
Write-Host "  $initCmd" -ForegroundColor White
Write-Host "  $capCmd" -ForegroundColor White
Write-Host ""

if (-not $nonInteractive) {
    Read-Host "按回车键退出"
}
