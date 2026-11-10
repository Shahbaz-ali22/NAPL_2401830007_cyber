@echo off
title Smart Expense Tracker (Go) - ENSP465 Capstone
cd /d "%~dp0"
if exist expense-tracker.exe (
    expense-tracker.exe
) else (
    echo Compiling Smart Expense Tracker...
    go build -o expense-tracker.exe main.go
    expense-tracker.exe
)
pause

