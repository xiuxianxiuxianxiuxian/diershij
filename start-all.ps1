# Cultivation World - Start Server Services
# Uses compiled binaries for instant startup

$ErrorActionPreference = "Continue"
$Green = "Green"; $Yellow = "Yellow"; $Red = "Red"; $Cyan = "Cyan"

Write-Host "========================================" -ForegroundColor $Cyan
Write-Host "    Cultivation World - 启动服务端      " -ForegroundColor $Cyan
Write-Host "========================================" -ForegroundColor $Cyan
Write-Host ""

# ── Database / Redis connection config ──
$env:DB_PASSWORD = "123456"
$env:DB_HOST = "localhost"
$env:DB_PORT = "5432"
$env:DB_USER = "postgres"
$env:DB_NAME = "cultivation"
$env:REDIS_HOST = "localhost"
$env:REDIS_PORT = "6379"
$env:SERVER_PORT = "8081"

$RootDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$ServerDir = Join-Path $RootDir "server"
$LogDir = Join-Path $RootDir "logs"
if (-not (Test-Path $LogDir)) { New-Item -ItemType Directory -Path $LogDir -Force | Out-Null }

$Global:ServiceProcs = @{}   # name → Process

function Log($msg, $color = $Green) {
    Write-Host "[$(Get-Date -Format 'HH:mm:ss')] $msg" -ForegroundColor $color
}

function Free-Port($port) {
    try {
        $conn = Get-NetTCPConnection -LocalPort $port -ErrorAction SilentlyContinue
        if ($conn) {
            $pid = $conn.OwningProcess
            $proc = Get-Process -Id $pid -ErrorAction SilentlyContinue
            if ($proc) {
                Log "释放端口 $port (PID: $pid, $($proc.ProcessName))..." $Yellow
                Stop-Process -Id $pid -Force -ErrorAction SilentlyContinue
                Start-Sleep -Milliseconds 500
                if (Get-NetTCPConnection -LocalPort $port -ErrorAction SilentlyContinue) {
                    taskkill /F /PID $pid 2>$null
                    Start-Sleep -Milliseconds 500
                }
                Log "端口 $port 已释放" $Green
            }
        }
    } catch { Log "释放端口 $port 异常: $_" $Red }
}

function Test-Port($port) {
    try { return (Test-NetConnection localhost -Port $port -WarningAction SilentlyContinue).TcpTestSucceeded }
    catch { return $false }
}

function Start-Svc($name, $exe, $workDir, $port) {
    $log = Join-Path $LogDir "$($name -replace ' ','_').log"
    $errLog = $log -replace '\.log$','_err.log'

    Log "启动 $name (端口 $port)..." $Yellow

    # Free the port first (important after restart)
    Free-Port $port

    # Build full path to exe
    $exePath = Join-Path $workDir $exe

    # Auto-build if binary missing
    if (-not (Test-Path $exePath)) {
        Log "未找到 $exe，正在编译..." $Yellow
        Push-Location $workDir
        go build -o $exe ./cmd 2>&1 | ForEach-Object { Write-Host "  $_" }
        Pop-Location
        if (-not (Test-Path $exePath)) { Log "编译 $name 失败" $Red; return $false }
    }

    # Start in background, capture output to log files
    $p = Start-Process -FilePath $exePath -WorkingDirectory $workDir `
        -NoNewWindow -PassThru `
        -RedirectStandardOutput $log -RedirectStandardError $errLog

    $Global:ServiceProcs[$name] = $p

    # Wait for port to open (up to 15s)
    for ($i = 0; $i -lt 15; $i++) {
        Start-Sleep -Seconds 1
        if ($p.HasExited) {
            Log "$name 进程已退出 (exit code: $($p.ExitCode))" $Red
            if (Test-Path $errLog) { Get-Content $errLog -Tail 5 | ForEach-Object { Write-Host "  $_" -ForegroundColor $Red } }
            return $false
        }
        try {
            if ((Test-NetConnection -ComputerName "localhost" -Port $port -WarningAction SilentlyContinue).TcpTestSucceeded) {
                Log "$name 启动成功 (PID: $($p.Id))" $Green; return $true
            }
        } catch {}
    }
    Log "$name 端口 $port 未在 15s 内就绪" $Yellow
    return $true  # might still come up
}

function Stop-All {
    Log "正在停止所有服务..." $Yellow
    $Global:ServiceProcs.Values | Where-Object { -not $_.HasExited } | ForEach-Object {
        Stop-Process -Id $_.Id -Force -ErrorAction SilentlyContinue
    }
    # Also kill orphaned processes
    @("game-server","gateway","heavenly-dao","ai-scheduler","world-engine") | ForEach-Object {
        Get-Process -Name $_ -ErrorAction SilentlyContinue | Stop-Process -Force
    }
    Log "所有服务已停止" $Green
}

function Show-Status {
    Write-Host ""
    Write-Host "========================================" -ForegroundColor $Cyan
    Write-Host "         服务状态                       " -ForegroundColor $Cyan
    Write-Host "========================================" -ForegroundColor $Cyan
    @( @{N="Gateway";P=8081}, @{N="Game Server";P=50051}, @{N="AI Scheduler";P=50052},
       @{N="Heavenly Dao";P=50053}, @{N="World Engine";P=50054},
       @{N="PostgreSQL";P=5432}, @{N="Redis";P=6379} ) | ForEach-Object {
        $ok = Test-Port $_.P
        if ($ok) { Write-Host "  [OK]  $($_.N) (port: $($_.P))" -ForegroundColor $Green }
        else     { Write-Host "  [--]  $($_.N) (port: $($_.P))" -ForegroundColor $Red }
    }
    Write-Host ""
    Write-Host "API Gateway: http://localhost:8081/" -ForegroundColor $Cyan
    Write-Host "WebSocket:   ws://localhost:8081/ws" -ForegroundColor $Cyan
    Write-Host "日志目录:    $LogDir" -ForegroundColor $Cyan
}

# ── Prerequisites ──
Log "检查 PostgreSQL..." $Yellow
if (-not (Test-Port 5432)) {
    Log "PostgreSQL 未运行，请先启动 PostgreSQL" $Red; pause; exit 1
}
Log "PostgreSQL OK" $Green

Log "检查 Redis..." $Yellow
if (-not (Test-Port 6379)) {
    Log "Redis 未运行（可选）" $Yellow
} else { Log "Redis OK" $Green }
Write-Host ""

# ── Start services in order ──
$svcs = @(
    @{N="World Engine";  E="world-engine.exe";  D=(Join-Path $ServerDir "world-engine");  P=50054}
    @{N="Heavenly Dao";  E="heavenly-dao.exe";  D=(Join-Path $ServerDir "heavenly-dao");  P=50053}
    @{N="AI Scheduler"; E="ai-scheduler.exe";  D=(Join-Path $ServerDir "ai-scheduler");   P=50052}
    @{N="Game Server";  E="game-server.exe";   D=(Join-Path $ServerDir "game-server");    P=50051}
    @{N="Gateway";      E="gateway.exe";       D=(Join-Path $ServerDir "gateway");        P=8081}
)

# Optional: rebuild all if -rebuild flag passed
if ($args -contains "-rebuild") {
    Log "重新编译所有服务..." $Yellow
    $svcs | ForEach-Object {
        Log "编译 $($_.N)..." $Yellow
        Push-Location $_.D; go build -o $_.E ./cmd; Pop-Location
    }
    Log "编译完成" $Green; Write-Host ""
}

$ok = $true
$svcs | ForEach-Object {
    if (-not (Start-Svc $_.N $_.E $_.D $_.P)) { $ok = $false }
    Start-Sleep -Seconds 1
}

if (-not $ok) {
    Log "部分服务启动失败，停止所有服务..." $Red; Stop-All; pause; exit 1
}

# Show final status
Start-Sleep -Seconds 2; Show-Status

Write-Host ""
Write-Host "========================================" -ForegroundColor $Cyan
Write-Host "    所有服务已启动!                      " -ForegroundColor $Cyan
Write-Host "========================================" -ForegroundColor $Cyan
Write-Host ""
Write-Host "按 Ctrl+C 停止所有服务" -ForegroundColor $Yellow
Write-Host "提示: 使用 停止游戏.bat 可干净关闭服务" -ForegroundColor $Yellow
Write-Host ""

# Keep alive & monitor
while ($true) {
    Start-Sleep -Seconds 2
    $dead = $Global:ServiceProcs.Keys | Where-Object { $Global:ServiceProcs[$_].HasExited }
    if ($dead) { Log "已停止: $($dead -join ', ')  (查看日志: $LogDir)" $Red }

    if ([Console]::KeyAvailable) {
        $key = [Console]::ReadKey($true)
        if ($key.Key -eq "C" -and $key.Modifiers -eq "Control") { Stop-All; break }
    }
}
