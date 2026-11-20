@echo off
set "PATH=%LOCALAPPDATA%\Microsoft\WinGet\Packages\Babashka.Babashka_Microsoft.Winget.Source_8wekyb3d8bbwe;%PATH%"
bb "%~dp0Q1.clj"
pause

