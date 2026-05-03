# Cultivation World - Stop All Services
# PowerShell Script

$ErrorActionPreference = "SilentlyContinue"

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "    Cultivation World - Stopper         " -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

# Stop Go processes
$goProcesses = Get-Process | Where-Object { $_.ProcessName -eq "go" }
if ($goProcesses) {
    Write-Host "Stopping Go services..." -ForegroundColor Yellow
    foreach ($proc in $goProcesses) {
        Stop-Process -Id $proc.Id -Force
        Write-Host "  Stopped Go PID: $($proc.Id)" -ForegroundColor Green
    }
}

# Stop Node processes
$nodeProcesses = Get-Process | Where-Object { $_.ProcessName -eq "node" }
if ($nodeProcesses) {
    Write-Host "Stopping Node processes..." -ForegroundColor Yellow
    foreach ($proc in $nodeProcesses) {
        Stop-Process -Id $proc.Id -Force
        Write-Host "  Stopped Node PID: $($proc.Id)" -ForegroundColor Green
    }
}

# Stop PowerShell child processes
$psProcesses = Get-CimInstance Win32_Process | Where-Object { 
    $_.Name -eq "powershell.exe" -and $_.CommandLine -like "*diershij*" 
}
if ($psProcesses) {
    Write-Host "Stopping PowerShell processes..." -ForegroundColor Yellow
    foreach ($proc in $psProcesses) {
        Stop-Process -Id $proc.ProcessId -Force
        Write-Host "  Stopped PowerShell PID: $($proc.ProcessId)" -ForegroundColor Green
    }
}

Write-Host ""
Write-Host "========================================" -ForegroundColor Green
Write-Host "    All services stopped!               " -ForegroundColor Green
Write-Host "========================================" -ForegroundColor Green
Write-Host ""

# Check port status
Write-Host "Port Status:" -ForegroundColor Cyan
$ports = @(8081, 50051, 50052, 50053, 50054, 1420)
foreach ($port in $ports) {
    $connection = Test-NetConnection -ComputerName "localhost" -Port $port -WarningAction SilentlyContinue
    if ($connection.TcpTestSucceeded) {
        Write-Host "  Port $port - Still in use" -ForegroundColor Red
    } else {
        Write-Host "  Port $port - Released" -ForegroundColor Green
    }
}

Write-Host ""
