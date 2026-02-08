@echo off
echo ========================================
echo   IPO~Master BY dallefx
echo   Professional IPO Automation Tool
echo ========================================
echo.

REM Check if Go is installed
go version >nul 2>&1
if errorlevel 1 (
    echo ERROR: Go is not installed!
    echo Please download and install Go from: https://golang.org/dl/
    echo.
    pause
    exit /b 1
)

echo Starting IPO~Master...
echo.

REM Run the customer version
go run ipo_master_customer.go

echo.
echo IPO~Master has exited.
pause