package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"smart-expense-tracker/internal/service"
	"smart-expense-tracker/internal/storage"
	"smart-expense-tracker/internal/ui"
)

func main() {
	// Database file location
	defaultDB := "expenses.db"
	if envDB := os.Getenv("EXPENSE_DB_PATH"); envDB != "" {
		defaultDB = envDB
	}

	// Optional global flag for custom database file
	dbPathFlag := flag.String("db", defaultDB, "Path to SQLite database file")
	
	// Only parse top-level flags if not a known subcommand
	args := os.Args[1:]
	if len(args) > 0 && (args[0] == "-db" || args[0] == "--db") {
		flag.CommandLine.Parse(args)
		args = flag.Args()
	}

	absDBPath, err := filepath.Abs(*dbPathFlag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error resolving database path: %v\n", err)
		os.Exit(1)
	}

	// Initialize Storage Layer (SQLite with parameterized queries)
	repo, err := storage.NewSQLiteRepository(absDBPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing database (%s): %v\n", absDBPath, err)
		os.Exit(1)
	}
	defer repo.Close()

	// Initialize Application Logic / Service Layer
	svc := service.NewTrackerService(repo)

	// Route based on command-line arguments
	if len(args) > 0 {
		ui.RunCLI(svc, args)
	} else {
		// Launch interactive mode
		cli := ui.NewInteractiveCLI(svc)
		cli.Run()
	}
}

