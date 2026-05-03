@echo off
chcp 65001 >nul
echo ========================================
echo    Cultivation World - Stopper
echo ========================================
echo.

powershell -ExecutionPolicy Bypass -Command "& '%~dp0stop-all.ps1'"

pause
