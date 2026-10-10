param([switch]$SkipCore)
$ErrorActionPreference = 'Stop'
$appRoot = (Resolve-Path (Join-Path $PSScriptRoot '..')).Path
function Native([string]$Executable,[string[]]$Arguments) {
 & $Executable @Arguments
 if ($LASTEXITCODE -ne 0) { throw "Build step failed: $Executable (exit $LASTEXITCODE)" }
}
Push-Location $appRoot
try {
 if (-not $SkipCore) { & ./scripts/core.ps1 }
 Push-Location frontend
 try { Native 'npm.cmd' @('ci'); Native 'npm.cmd' @('run','build') } finally { Pop-Location }
 Native 'go' @('test','./internal/core')
 Native 'go' @('run','github.com/wailsapp/wails/v3/cmd/wails3@v3.0.0-alpha.74','generate','syso','-arch','amd64','-icon','build/windows/icon.ico','-manifest','build/windows/wails.exe.manifest','-info','build/windows/info.json','-out','wails_windows_amd64.syso')
 New-Item -ItemType Directory -Force build/bin | Out-Null
 $env:CGO_ENABLED = '0'
 Native 'go' @('build','-tags','production','-trimpath','-buildvcs=false','-ldflags=-w -s -H windowsgui','-o','build/bin/Sudoku.exe','.')
 Write-Host 'Ready: build/bin/Sudoku.exe'
} finally { Pop-Location }
