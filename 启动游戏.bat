@echo off
chcp 65001 >nul
title Cultivation World - Start Server
cd /d "%~dp0"

echo ========================================
echo    Cultivation World - Start Server
echo ========================================
echo.

echo [INFO] Starting all services...
powershell -NoExit -ExecutionPolicy Bypass -File "%~dp0start-all.ps1"
pause
