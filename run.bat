@echo off
REM ==============================================================================
REM NeuroFlow - Lanceur Windows optimisé UTF-8
REM ==============================================================================

chcp 65001 > nul
cls

if exist "neuroflow.exe" (
    neuroflow.exe %*
) else (
    echo [NeuroFlow] Compilation en cours...
    go build -ldflags="-s -w" -o neuroflow.exe .\cmd\neuroflow
    if %errorlevel% neq 0 (
        echo [Erreur] Échec de la compilation Go sous Windows.
        pause
        exit /b 1
    )
    neuroflow.exe %*
)
