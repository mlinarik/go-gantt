# Quick Start Script for Windows PowerShell

Write-Host "==================================" -ForegroundColor Cyan
Write-Host "Go Gantt Chart Application Setup" -ForegroundColor Cyan
Write-Host "==================================" -ForegroundColor Cyan
Write-Host ""

# Check if Docker is available
$dockerAvailable = $false
try {
    docker --version | Out-Null
    $dockerAvailable = $true
    Write-Host "✓ Docker is installed" -ForegroundColor Green
} catch {
    Write-Host "✗ Docker is not installed" -ForegroundColor Yellow
}

Write-Host ""

# Check if Go is available
$goAvailable = $false
try {
    go version | Out-Null
    $goAvailable = $true
    Write-Host "✓ Go is installed" -ForegroundColor Green
} catch {
    Write-Host "✗ Go is not installed" -ForegroundColor Yellow
}

Write-Host ""
Write-Host "Choose how to run the application:" -ForegroundColor Cyan
Write-Host "1. Docker (recommended)" -ForegroundColor White
Write-Host "2. Local Go installation" -ForegroundColor White
Write-Host "3. Exit" -ForegroundColor White
Write-Host ""

$choice = Read-Host "Enter your choice (1-3)"

switch ($choice) {
    "1" {
        if (-not $dockerAvailable) {
            Write-Host "Docker is not available. Please install Docker first." -ForegroundColor Red
            exit 1
        }
        
        Write-Host ""
        Write-Host "Building Docker image..." -ForegroundColor Cyan
        docker build -t go-ghant .
        
        if ($LASTEXITCODE -eq 0) {
            Write-Host "✓ Docker image built successfully" -ForegroundColor Green
            Write-Host ""
            Write-Host "Starting container..." -ForegroundColor Cyan
            
            # Create data directory if it doesn't exist
            if (-not (Test-Path ".\data")) {
                New-Item -ItemType Directory -Path ".\data" | Out-Null
            }
            
            docker run -d -p 8080:8080 -v "${PWD}/data:/root" --name go-ghant-app go-ghant
            
            if ($LASTEXITCODE -eq 0) {
                Write-Host "✓ Container started successfully" -ForegroundColor Green
                Write-Host ""
                Write-Host "Application is running at: http://localhost:8080" -ForegroundColor Green
                Write-Host ""
                Write-Host "To stop the container, run: docker stop go-ghant-app" -ForegroundColor Yellow
                Write-Host "To remove the container, run: docker rm go-ghant-app" -ForegroundColor Yellow
                
                # Open browser
                Start-Process "http://localhost:8080"
            }
        }
    }
    "2" {
        if (-not $goAvailable) {
            Write-Host "Go is not available. Please install Go first." -ForegroundColor Red
            exit 1
        }
        
        Write-Host ""
        Write-Host "Downloading dependencies..." -ForegroundColor Cyan
        go mod download
        
        Write-Host "Starting application..." -ForegroundColor Cyan
        Write-Host ""
        Write-Host "Application will run at: http://localhost:8080" -ForegroundColor Green
        Write-Host "Press Ctrl+C to stop the application" -ForegroundColor Yellow
        Write-Host ""
        
        # Open browser after a short delay
        Start-Job -ScriptBlock {
            Start-Sleep -Seconds 2
            Start-Process "http://localhost:8080"
        } | Out-Null
        
        go run .
    }
    "3" {
        Write-Host "Exiting..." -ForegroundColor Yellow
        exit 0
    }
    default {
        Write-Host "Invalid choice. Exiting..." -ForegroundColor Red
        exit 1
    }
}
