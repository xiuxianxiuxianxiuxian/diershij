# Cultivation World - Start Server Services Only
# PowerShell Script

$ErrorActionPreference = "Continue"

# Colors
$Green = "Green"
$Yellow = "Yellow"
$Red = "Red"
$Cyan = "Cyan"

Write-Host "========================================" -ForegroundColor $Cyan
Write-Host "    Cultivation World - Server Only     " -ForegroundColor $Cyan
Write-Host "========================================" -ForegroundColor $Cyan
Write-Host ""

# Environment Variables
$env:DB_PASSWORD = "123456"
$env:DB_HOST = "localhost"
$env:DB_PORT = "5432"
$env:DB_USER = "postgres"
$env:DB_NAME = "cultivation"
$env:REDIS_HOST = "localhost"
$env:REDIS_PORT = "6379"
$env:SERVER_PORT = "8081"

# Directories
$RootDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$ServerDir = Join-Path $RootDir "server"

# Store Process IDs
$Global:ProcessIds = @()

function Log-Info($msg, $color = $Green) {
    Write-Host "[$(Get-Date -Format 'HH:mm:ss')] $msg" -ForegroundColor $color
}

function Start-MyService($name, $dir, $cmd, $port) {
    Log-Info "Starting $name (port: $port)..." $Yellow
    try {
        # 先检查端口是否被占用
        $portInUse = Get-NetTCPConnection -LocalPort $port -ErrorAction SilentlyContinue
        if ($portInUse) {
            Log-Info "WARNING: Port $port is already in use" $Yellow
        }
        
        $process = Start-Process -FilePath "powershell.exe" -ArgumentList "-NoExit", "-Command", "cd '$dir'; $cmd" -PassThru -WindowStyle Minimized
        $Global:ProcessIds += $process.Id
        
        # 等待服务启动
        Log-Info "Waiting for $name to start..." $Yellow
        $maxAttempts = 30
        $attempt = 0
        $started = $false
        
        while ($attempt -lt $maxAttempts -and -not $started) {
            Start-Sleep -Seconds 1
            $attempt++
            
            try {
                $connection = Test-NetConnection -ComputerName "localhost" -Port $port -WarningAction SilentlyContinue
                if ($connection.TcpTestSucceeded) {
                    $started = $true
                    Log-Info "$name started successfully (PID: $($process.Id))" $Green
                }
            } catch {
                # 继续等待
            }
        }
        
        if (-not $started) {
            Log-Info "$name failed to start within 30 seconds" $Red
            return $false
        }
        
        return $true
    } catch {
        Log-Info "$name failed: $_" $Red
        return $false
    }
}

function Show-MyStatus() {
    Write-Host ""
    Write-Host "========================================" -ForegroundColor $Cyan
    Write-Host "         Service Status                 " -ForegroundColor $Cyan
    Write-Host "========================================" -ForegroundColor $Cyan
    Write-Host ""
    
    $services = @(
        @{ Name = "Gateway"; Port = 8081 },
        @{ Name = "Game Server"; Port = 50051 },
        @{ Name = "AI Scheduler"; Port = 50052 },
        @{ Name = "Heavenly Dao"; Port = 50053 },
        @{ Name = "World Engine"; Port = 50054 },
        @{ Name = "PostgreSQL"; Port = 5432 },
        @{ Name = "Redis"; Port = 6379 }
    )
    
    foreach ($svc in $services) {
        try {
            $connection = Test-NetConnection -ComputerName "localhost" -Port $svc.Port -WarningAction SilentlyContinue
            if ($connection.TcpTestSucceeded) {
                Write-Host "  OK  $($svc.Name) (port: $($svc.Port))" -ForegroundColor $Green
            } else {
                Write-Host "  FAIL $($svc.Name) (port: $($svc.Port))" -ForegroundColor $Red
            }
        } catch {
            Write-Host "  FAIL $($svc.Name) (port: $($svc.Port))" -ForegroundColor $Red
        }
    }
    
    Write-Host ""
    Write-Host "API Gateway: http://localhost:8081/" -ForegroundColor $Cyan
    Write-Host "WebSocket:   ws://localhost:8081/ws" -ForegroundColor $Cyan
    Write-Host ""
}

function Stop-MyServices() {
    Write-Host ""
    Log-Info "Stopping all services..." $Yellow
    foreach ($pid in $Global:ProcessIds) {
        try {
            Stop-Process -Id $pid -Force -ErrorAction SilentlyContinue
            Log-Info "Stopped PID: $pid" $Green
        } catch {}
    }
    Get-Process | Where-Object { $_.ProcessName -eq "go" -or $_.ProcessName -eq "node" } | Stop-Process -Force -ErrorAction SilentlyContinue
    Log-Info "All services stopped" $Green
}

# Check PostgreSQL
Log-Info "Checking PostgreSQL..." $Yellow
$pg = Test-NetConnection -ComputerName "localhost" -Port 5432 -WarningAction SilentlyContinue
if (-not $pg.TcpTestSucceeded) {
    Log-Info "ERROR: PostgreSQL not running (port 5432)" $Red
    exit 1
}
Log-Info "PostgreSQL OK" $Green

# Check Redis
Log-Info "Checking Redis..." $Yellow
$redis = Test-NetConnection -ComputerName "localhost" -Port 6379 -WarningAction SilentlyContinue
if (-not $redis.TcpTestSucceeded) {
    Log-Info "WARNING: Redis not running (port 6379)" $Yellow
} else {
    Log-Info "Redis OK" $Green
}

Write-Host ""

# Start Backend Services in order
$services = @(
    @{ Name = "World Engine"; Dir = Join-Path $ServerDir "world-engine"; Cmd = "go run ./cmd"; Port = 50054 },
    @{ Name = "Heavenly Dao"; Dir = Join-Path $ServerDir "heavenly-dao"; Cmd = "go run ./cmd"; Port = 50053 },
    @{ Name = "AI Scheduler"; Dir = Join-Path $ServerDir "ai-scheduler"; Cmd = "go run ./cmd"; Port = 50052 },
    @{ Name = "Game Server"; Dir = Join-Path $ServerDir "game-server"; Cmd = "go run ./cmd"; Port = 50051 },
    @{ Name = "Gateway"; Dir = Join-Path $ServerDir "gateway"; Cmd = "go run ./cmd"; Port = 8081 }
)

foreach ($svc in $services) {
    $success = Start-MyService -name $svc.Name -dir $svc.Dir -cmd $svc.Cmd -port $svc.Port
    if (-not $success) {
        Log-Info "Failed to start $($svc.Name), stopping all services..." $Red
        Stop-MyServices
        exit 1
    }
    Start-Sleep -Seconds 2
}

# Wait for all services to be ready
Write-Host ""
Log-Info "Waiting for all services to be ready (5s)..." $Yellow
Start-Sleep -Seconds 5

# Show Status
Show-MyStatus

Write-Host "========================================" -ForegroundColor $Cyan
Write-Host "    All server services started!        " -ForegroundColor $Cyan
Write-Host "========================================" -ForegroundColor $Cyan
Write-Host ""
Write-Host "Press Ctrl+C to stop all services" -ForegroundColor $Yellow
Write-Host ""

# Keep running
while ($true) {
    if ([Console]::KeyAvailable) {
        $key = [Console]::ReadKey($true)
        if ($key.Key -eq "C" -and $key.Modifiers -eq "Control") {
            Stop-MyServices
            break
        }
    }
    Start-Sleep -Milliseconds 100
}
