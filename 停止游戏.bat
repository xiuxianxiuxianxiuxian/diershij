@echo off
chcp 65001 >nul
title Cultivation World - Stop Server
cd /d "%~dp0"

echo ========================================
echo    Cultivation World - Stop Server
echo ========================================
echo.

echo [INFO] Stopping all services...
powershell -ExecutionPolicy Bypass -File "%~dp0stop-all.ps1"
echo.
echo [INFO] Done! Press any key to exit.
pause
