package ui

import (
	"flag"
	"fmt"
	"os"

	"smart-expense-tracker/internal/model"
	"smart-expense-tracker/internal/service"
)

// RunCLI handles non-interactive command-line subcommands.
func RunCLI(svc *service.TrackerService, args []string) {
	if len(args) == 0 {
		PrintUsage()
		return
	}

	subcommand := args[0]
	cmdArgs := args[1:]

	switch subcommand {
	case "add":
		handleAddCmd(svc, cmdArgs)
	case "list":
		handleListCmd(svc, cmdArgs)
	case "balance":
		handleBalanceCmd(svc)
	case "summary":
		handleSummaryCmd(svc)
	case "update":
		handleUpdateCmd(svc, cmdArgs)
	case "delete":
		handleDeleteCmd(svc, cmdArgs)
	case "categories":
		handleCategoriesCmd(svc)
	case "export":
		handleExportCmd(svc, cmdArgs)
	case "seed":
		handleSeedCmd(svc)
	case "help", "-h", "--help":
		PrintUsage()
	default:
		fmt.Printf("Unknown command '%s'. Run 'expense-tracker help' for usage.\n", subcommand)
	}
}

func handleAddCmd(svc *service.TrackerService, args []string) {
	fs := flag.NewFlagSet("add", flag.ExitOnError)
	txType := fs.String("type", "expense", "Transaction type: 'income' or 'expense'")
	amount := fs.Float64("amount", 0.0, "Transaction amount (e.g. 45.50)")
	category := fs.String("category", "General", "Category name (e.g. Food, Transport)")
	date := fs.String("date", "", "Date in YYYY-MM-DD format (default: today)")
	desc := fs.String("desc", "", "Optional description")

	_ = fs.Parse(args)

	parsedType, err := model.ParseTransactionType(*txType)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	tx := &model.Transaction{
		Type:        parsedType,
		Amount:      model.NewMoneyFromFloat(*amount),
		Category:    *category,
		Date:        *date,
		Description: *desc,
	}

	if err := svc.AddTransaction(tx); err != nil {
		fmt.Printf("Error adding transaction: %v\n", err)
		return
	}

	fmt.Printf("[✓] Added %s: $%.2f (%s) on %s [ID: #%d]\n", tx.Type, tx.Amount.ToFloat(), tx.Category, tx.Date, tx.ID)
}

func handleListCmd(svc *service.TrackerService, args []string) {
	fs := flag.NewFlagSet("list", flag.ExitOnError)
	txType := fs.String("type", "", "Filter by type: 'income' or 'expense'")
	category := fs.String("category", "", "Filter by category")
	startDate := fs.String("from", "", "Filter starting from date (YYYY-MM-DD)")
	endDate := fs.String("to", "", "Filter up to date (YYYY-MM-DD)")

	_ = fs.Parse(args)

	var filter model.Filter
	if *txType != "" {
		pt, err := model.ParseTransactionType(*txType)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		filter.Type = pt
	}
	filter.Category = *category
	filter.StartDate = *startDate
	filter.EndDate = *endDate

	txs, err := svc.ListTransactions(filter)
	if err != nil {
		fmt.Printf("Error listing transactions: %v\n", err)
		return
	}
	PrintTransactionsTable(txs)
}

func handleBalanceCmd(svc *service.TrackerService) {
	summary, err := svc.GetBalanceSummary()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	PrintBalanceSummary(summary)
}

func handleSummaryCmd(svc *service.TrackerService) {
	summary, err := svc.GetBalanceSummary()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	PrintBalanceSummary(summary)

	catExp, _ := svc.GetExpenseBreakdownByCategory()
	PrintCategoryBreakdown("EXPENSE BREAKDOWN BY CATEGORY", catExp)

	catInc, _ := svc.GetIncomeBreakdownByCategory()
	if len(catInc) > 0 {
		PrintCategoryBreakdown("INCOME BREAKDOWN BY CATEGORY", catInc)
	}

	monthly, _ := svc.GetMonthlyBreakdown()
	PrintMonthlyBreakdown(monthly)
}

func handleUpdateCmd(svc *service.TrackerService, args []string) {
	fs := flag.NewFlagSet("update", flag.ExitOnError)
	id := fs.Int64("id", 0, "Transaction ID to update (required)")
	txType := fs.String("type", "", "New type: 'income' or 'expense'")
	amount := fs.Float64("amount", -1.0, "New amount")
	category := fs.String("category", "", "New category")
	date := fs.String("date", "", "New date (YYYY-MM-DD)")
	desc := fs.String("desc", "", "New description")

	_ = fs.Parse(args)

	if *id <= 0 {
		fmt.Println("Error: --id flag is required and must be greater than 0")
		return
	}

	tx, err := svc.GetTransaction(*id)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if *txType != "" {
		pt, err := model.ParseTransactionType(*txType)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		tx.Type = pt
	}
	if *amount > 0 {
		tx.Amount = model.NewMoneyFromFloat(*amount)
	}
	if *category != "" {
		tx.Category = *category
	}
	if *date != "" {
		tx.Date = *date
	}
	if *desc != "" {
		tx.Description = *desc
	}

	if err := svc.UpdateTransaction(tx); err != nil {
		fmt.Printf("Error updating transaction: %v\n", err)
		return
	}

	fmt.Printf("[✓] Transaction #%d updated successfully.\n", tx.ID)
}

func handleDeleteCmd(svc *service.TrackerService, args []string) {
	fs := flag.NewFlagSet("delete", flag.ExitOnError)
	id := fs.Int64("id", 0, "Transaction ID to delete (required)")
	_ = fs.Parse(args)

	if *id <= 0 {
		fmt.Println("Error: --id flag is required and must be greater than 0")
		return
	}

	if err := svc.DeleteTransaction(*id); err != nil {
		fmt.Printf("Error deleting transaction: %v\n", err)
		return
	}

	fmt.Printf("[✓] Transaction #%d deleted successfully.\n", *id)
}

func handleCategoriesCmd(svc *service.TrackerService) {
	cats, err := svc.GetCategories()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Println("\nRegistered Categories:")
	for i, c := range cats {
		fmt.Printf("  %d. %s\n", i+1, c)
	}
	fmt.Println()
}

func handleExportCmd(svc *service.TrackerService, args []string) {
	fs := flag.NewFlagSet("export", flag.ExitOnError)
	format := fs.String("format", "csv", "Export format: 'csv' or 'json'")
	filePath := fs.String("file", "", "Target file name (default: expenses.csv or expenses.json)")
	_ = fs.Parse(args)

	targetFile := *filePath
	if targetFile == "" {
		if *format == "json" {
			targetFile = "expenses.json"
		} else {
			targetFile = "expenses.csv"
		}
	}

	f, err := os.Create(targetFile)
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		return
	}
	defer f.Close()

	if *format == "json" {
		err = svc.ExportJSON(f)
	} else {
		err = svc.ExportCSV(f)
	}

	if err != nil {
		fmt.Printf("Export failed: %v\n", err)
		return
	}

	fmt.Printf("[✓] Exported all transactions to '%s'\n", targetFile)
}

func handleSeedCmd(svc *service.TrackerService) {
	if err := svc.SeedSampleData(); err != nil {
		fmt.Printf("Error seeding demo data: %v\n", err)
		return
	}
	fmt.Println("[✓] Sample demonstration transactions loaded successfully. Run 'list' or 'summary' to inspect.")
}

// PrintUsage displays helpful CLI commands.
func PrintUsage() {
	fmt.Println(`Smart Expense Tracker - CLI Usage:

  expense-tracker [command] [flags]

Commands:
  (no command)        Launch interactive terminal mode
  add                 Add a new financial transaction
                      Flags: --type [income|expense] --amount [0.00] --category [name] --date [YYYY-MM-DD] --desc [text]
  list                List recorded transactions
                      Flags: --type [income|expense] --category [name] --from [YYYY-MM-DD] --to [YYYY-MM-DD]
  balance             Show total income, total expenses, and remaining balance
  summary             Show category spending breakdowns and monthly trends
  update              Update an existing transaction
                      Flags: --id [ID] [--type] [--amount] [--category] [--date] [--desc]
  delete              Delete an existing transaction
                      Flags: --id [ID]
  categories          Show all registered categories
  export              Export transactions to file
                      Flags: --format [csv|json] --file [filename]
  seed                Load sample student transactions across 3 months for demo
  help                Display this help screen`)
	os.Stdout.Sync()
}

