@echo off
chcp 65001 >nul
echo ========================================
echo    Cultivation World - Starter
echo ========================================
echo.

powershell -ExecutionPolicy Bypass -Command "& '%~dp0start-all.ps1'"

pause
