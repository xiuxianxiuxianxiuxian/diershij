@echo off
chcp 65001 >nul
echo ========================================
echo    Building Cultivation World Desktop App
echo ========================================
echo.

cd /d "%~dp0client"

echo Installing dependencies...
call npm install

echo.
echo Building Tauri desktop app...
call npm run tauri build

echo.
echo ========================================
echo    Build Complete!
echo ========================================
echo.
echo Executable location:
echo   src-tauri\target\release\修仙世界.exe
echo.
pause
