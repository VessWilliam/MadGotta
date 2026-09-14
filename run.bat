@echo off

echo ========================================
echo Development Start
echo ========================================

echo.
echo Checking Tailwind...
go run setup.go

echo.
echo Closing previous Tailwind CMD...
taskkill /FI "WINDOWTITLE eq Tailwind CSS*" /T /F >nul 2>&1

echo.
echo Closing previous Templ CMD...
taskkill /FI "WINDOWTITLE eq Templ Watcher*" /T /F >nul 2>&1

echo.
echo Starting Templ watcher...
start "Templ Watcher" cmd /k "templ generate --watch"

echo.
echo Starting Tailwind watcher...
start "Tailwind CSS" cmd /k ".\web\static\tailwindcss.exe -i .\web\static\css\input.css -o .\web\static\css\output.css --watch"

echo.
echo Starting Air...
air
