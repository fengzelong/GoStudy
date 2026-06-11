Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"

$Root = Split-Path -Parent $PSScriptRoot
$GoModCache = Join-Path $Root ".tmp_gomodcache"
$GoPath = Join-Path $Root ".tmp_gopath"

New-Item -ItemType Directory -Force -Path $GoModCache, $GoPath | Out-Null

$env:GOMODCACHE = $GoModCache
$env:GOPATH = $GoPath

function Invoke-GoTest {
    param(
        [Parameter(Mandatory = $true)]
        [string]$Path,

        [Parameter(Mandatory = $true)]
        [string]$Pattern
    )

    Push-Location $Path
    try {
        Write-Host "==> go test $Pattern ($Path)"
        go test $Pattern
        if ($LASTEXITCODE -ne 0) {
            throw "go test $Pattern failed in $Path"
        }
    }
    finally {
        Pop-Location
    }
}

Invoke-GoTest -Path $Root -Pattern "./..."

$SubModules = @(
    "gomysql",
    "goredis",
    "gorm",
    "gorabbitmq",
    "websocket"
)

foreach ($Module in $SubModules) {
    $ModulePath = Join-Path $Root $Module
    Invoke-GoTest -Path $ModulePath -Pattern "./..."
}

Write-Host "All tests passed"
