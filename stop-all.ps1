# Cultivation World - Stop All Services
$ErrorActionPreference = "SilentlyContinue"

$Green = "Green"; $Yellow = "Yellow"; $Red = "Red"; $Cyan = "Cyan"

Write-Host "========================================" -ForegroundColor $Cyan
Write-Host "    Cultivation World - 停止服务端      " -ForegroundColor $Cyan
Write-Host "========================================" -ForegroundColor $Cyan
Write-Host ""

# Service executable names (without .exe — use process name)
$svcExes = @("game-server", "gateway", "heavenly-dao", "ai-scheduler", "world-engine")
$svcPorts = @(8081, 50051, 50052, 50053, 50054)
$allPorts = @(8081, 50051, 50052, 50053, 50054, 5432, 6379, 1420)

# 1. Kill service processes by name
Write-Host "[1/3] 终止服务进程..." -ForegroundColor $Yellow
$killed = 0
foreach ($name in $svcExes) {
    $procs = Get-Process -Name $name -ErrorAction SilentlyContinue
    foreach ($p in $procs) {
        Write-Host "  关闭 $name (PID: $($p.Id))" -ForegroundColor $Yellow
        Stop-Process -Id $p.Id -Force
        $killed++
    }
}
if ($killed -eq 0) { Write-Host "  未发现运行中的服务进程" -ForegroundColor $Green }
else { Write-Host "  已关闭 $killed 个进程" -ForegroundColor $Green }

# 2. Kill any orphaned Go build processes (from `go run`)
$goProcs = Get-Process -Name "go" -ErrorAction SilentlyContinue | Where-Object {
    $cmd = (Get-CimInstance Win32_Process -Filter "ProcessId = $($_.Id)" -ErrorAction SilentlyContinue).CommandLine
    $cmd -match "diershij|go run"
}
if ($goProcs) {
    Write-Host "  关闭 Go 编译进程..." -ForegroundColor $Yellow
    $goProcs | Stop-Process -Force
}

# 3. Force-release service ports
Write-Host "[2/3] 释放服务端口..." -ForegroundColor $Yellow
$released = 0
foreach ($port in $svcPorts) {
    $conn = Get-NetTCPConnection -LocalPort $port -ErrorAction SilentlyContinue
    if ($conn) {
        $pid = $conn.OwningProcess
        $proc = Get-Process -Id $pid -ErrorAction SilentlyContinue
        if ($proc) {
            Write-Host "  释放端口 $port (PID: $pid, $($proc.ProcessName))" -ForegroundColor $Yellow
            Stop-Process -Id $pid -Force
            $released++
        }
    }
}
if ($released -eq 0) { Write-Host "  所有端口已释放" -ForegroundColor $Green }
else { Write-Host "  已释放 $released 个端口" -ForegroundColor $Green }

# 4. Final port status check
Write-Host "[3/3] 端口状态检查..." -ForegroundColor $Yellow
Write-Host ""
Write-Host "  Port Status:" -ForegroundColor $Cyan
foreach ($port in $allPorts) {
    $inUse = $false
    try { $inUse = (Test-NetConnection localhost -Port $port -WarningAction SilentlyContinue).TcpTestSucceeded } catch {}
    if ($inUse) {
        $owner = Get-NetTCPConnection -LocalPort $port -ErrorAction SilentlyContinue
        $ownerName = if ($owner) { "(PID: $($owner.OwningProcess))" } else { "" }
        Write-Host "    Port $port - 仍在使用 $ownerName" -ForegroundColor $Red
    } else {
        Write-Host "    Port $port - 已释放" -ForegroundColor $Green
    }
}

Write-Host ""
Write-Host "========================================" -ForegroundColor $Green
Write-Host "    所有服务已停止!                      " -ForegroundColor $Green
Write-Host "========================================" -ForegroundColor $Green
Write-Host ""
