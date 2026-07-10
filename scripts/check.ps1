Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"

$Root = Split-Path -Parent $PSScriptRoot
Push-Location $Root

try {
    Write-Host "==> gofmt check changed Go files"
    $ChangedGoFiles = git -c core.excludesfile= status --short --untracked-files=all |
        ForEach-Object {
            $Path = $_.Substring(3).Trim()
            if ($Path -like "*.go" -and (Test-Path $Path) -and $Path -notmatch "^\.tmp_") {
                Resolve-Path $Path
            }
        }

    $Files = @()
    if ($ChangedGoFiles) {
        $Files = gofmt -l $ChangedGoFiles
    }

    if ($Files) {
        $Files | ForEach-Object { Write-Host $_ }
        throw "gofmt check failed"
    }

    Write-Host "==> full test"
    & (Join-Path $PSScriptRoot "test.ps1")

    Write-Host "==> go vet"
    go vet ./...

    Write-Host "==> git status"
    git -c core.excludesfile= status --short
}
finally {
    Pop-Location
}
