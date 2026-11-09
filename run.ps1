# ==============================================================================
# NeuroFlow - Lanceur PowerShell Windows (UTF-8 & Compatibilité)
# ==============================================================================

[Console]::OutputEncoding = [System.Text.Encoding]::UTF8
[Console]::InputEncoding  = [System.Text.Encoding]::UTF8

$exePath = Join-Path $PSScriptRoot "neuroflow.exe"

if (-not (Test-Path $exePath)) {
    Write-Host "[NeuroFlow] Compilation en cours pour Windows..." -ForegroundColor Cyan
    & go build -ldflags="-s -w" -o $exePath (Join-Path $PSScriptRoot "cmd\neuroflow")
    if ($LASTEXITCODE -ne 0) {
        Write-Error "[Erreur] Échec de la compilation Go."
        Exit 1
    }
}

& $exePath $args
