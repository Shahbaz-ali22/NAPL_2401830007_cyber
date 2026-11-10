package ui

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"smart-expense-tracker/internal/model"
	"smart-expense-tracker/internal/service"
)

// InteractiveCLI runs the interactive terminal interface.
type InteractiveCLI struct {
	svc    *service.TrackerService
	reader *bufio.Reader
}

// NewInteractiveCLI creates a new interactive shell session.
func NewInteractiveCLI(svc *service.TrackerService) *InteractiveCLI {
	return &InteractiveCLI{
		svc:    svc,
		reader: bufio.NewReader(os.Stdin),
	}
}

// Run starts the main loop.
func (ui *InteractiveCLI) Run() {
	for {
		ui.printMainMenu()
		choice := ui.readLine("Select an option (0-11): ")

		switch choice {
		case "1":
			ui.handleAddTransaction()
		case "2":
			ui.handleViewTransactions()
		case "3":
			ui.handleFilterTransactions()
		case "4":
			ui.handleCalculateTotals()
		case "5":
			ui.handleReviewSpending()
		case "6":
			ui.handleUpdateTransaction()
		case "7":
			ui.handleDeleteTransaction()
		case "8":
			ui.handleListCategories()
		case "9":
			ui.handleExport()
		case "10":
			ui.handleSeed()
		case "11":
			ui.handleGenerateReport()
		case "0", "exit", "quit":
			fmt.Println("\nThank you for using Smart Expense Tracker. Goodbye!")
			return
		default:
			fmt.Println("\n[!] Invalid option. Please enter a number between 0 and 11.")
		}

		fmt.Println("\nPress Enter to return to main menu...")
		_, _ = ui.reader.ReadString('\n')
	}
}

func (ui *InteractiveCLI) printMainMenu() {
	fmt.Println("==================================================")
	fmt.Println("         SMART EXPENSE TRACKER (Go CLI)           ")
	fmt.Println("     K.R. Mangalam University - Capstone Project  ")
	fmt.Println("==================================================")
	fmt.Println(" [1] Add Transaction (Income / Expense)")
	fmt.Println(" [2] View All Transactions")
	fmt.Println(" [3] Filter Transactions (Type / Category / Date)")
	fmt.Println(" [4] Calculate Totals & Remaining Balance")
	fmt.Println(" [5] Review Spending (Category & Monthly Summaries)")
	fmt.Println(" [6] Update an Existing Transaction")
	fmt.Println(" [7] Delete a Transaction")
	fmt.Println(" [8] View Registered Categories")
	fmt.Println(" [9] Export Transactions (CSV / JSON)")
	fmt.Println(" [10] Load Sample Demo Data")
	fmt.Println(" [11] Open Interactive Web Dashboard (Browser)")
	fmt.Println(" [0] Exit Application")
	fmt.Println("--------------------------------------------------")
}

func (ui *InteractiveCLI) handleAddTransaction() {
	fmt.Println("\n--- [ Add New Transaction ] ---")

	// 1. Transaction Type
	var txType model.TransactionType
	for {
		input := ui.readLine("Transaction Type (1: Expense, 2: Income) [default: 1]: ")
		if input == "" || input == "1" || strings.ToLower(input) == "expense" {
			txType = model.TypeExpense
			break
		} else if input == "2" || strings.ToLower(input) == "income" {
			txType = model.TypeIncome
			break
		}
		fmt.Println("  [!] Invalid choice. Enter 1 for Expense or 2 for Income.")
	}

	// 2. Amount
	var amount model.Money
	for {
		input := ui.readLine("Amount ($): ")
		val, err := strconv.ParseFloat(strings.TrimSpace(input), 64)
		if err != nil || val <= 0 {
			fmt.Println("  [!] Please enter a valid positive numerical amount (e.g., 25.50).")
			continue
		}
		amount = model.NewMoneyFromFloat(val)
		break
	}

	// 3. Category
	fmt.Println("Common categories: Food, Transport, Education, Shopping, Utilities, Salary, Other")
	category := ui.readLine("Category: ")
	if category == "" {
		category = "General"
	}

	// 4. Date
	today := time.Now().Format("2006-01-02")
	datePrompt := fmt.Sprintf("Date (YYYY-MM-DD) [default %s]: ", today)
	date := ui.readLine(datePrompt)
	if date == "" {
		date = today
	}

	// 5. Description
	desc := ui.readLine("Description (optional): ")

	tx := &model.Transaction{
		Type:        txType,
		Amount:      amount,
		Category:    category,
		Date:        date,
		Description: desc,
	}

	err := ui.svc.AddTransaction(tx)
	if err != nil {
		fmt.Printf("\n[X] Error saving transaction: %v\n", err)
		return
	}

	fmt.Printf("\n[✓] Transaction successfully recorded! Assigned ID: #%d\n", tx.ID)
}

func (ui *InteractiveCLI) handleViewTransactions() {
	fmt.Println("\n--- [ All Transactions ] ---")
	txs, err := ui.svc.ListTransactions(model.Filter{})
	if err != nil {
		fmt.Printf("[X] Error fetching transactions: %v\n", err)
		return
	}
	PrintTransactionsTable(txs)
}

func (ui *InteractiveCLI) handleFilterTransactions() {
	fmt.Println("\n--- [ Filter Transactions ] ---")
	var filter model.Filter

	typeInput := ui.readLine("Filter by Type (1: Expense, 2: Income, Enter to skip): ")
	if typeInput == "1" || strings.ToLower(typeInput) == "expense" {
		filter.Type = model.TypeExpense
	} else if typeInput == "2" || strings.ToLower(typeInput) == "income" {
		filter.Type = model.TypeIncome
	}

	catInput := ui.readLine("Filter by Category (Enter to skip): ")
	if strings.TrimSpace(catInput) != "" {
		filter.Category = strings.TrimSpace(catInput)
	}

	startInput := ui.readLine("From Date (YYYY-MM-DD, Enter to skip): ")
	if strings.TrimSpace(startInput) != "" {
		filter.StartDate = strings.TrimSpace(startInput)
	}

	endInput := ui.readLine("To Date (YYYY-MM-DD, Enter to skip): ")
	if strings.TrimSpace(endInput) != "" {
		filter.EndDate = strings.TrimSpace(endInput)
	}

	txs, err := ui.svc.ListTransactions(filter)
	if err != nil {
		fmt.Printf("[X] Error filtering transactions: %v\n", err)
		return
	}
	PrintTransactionsTable(txs)
}

func (ui *InteractiveCLI) handleCalculateTotals() {
	summary, err := ui.svc.GetBalanceSummary()
	if err != nil {
		fmt.Printf("[X] Error: %v\n", err)
		return
	}
	PrintBalanceSummary(summary)
}

func (ui *InteractiveCLI) handleReviewSpending() {
	fmt.Println("\n--- [ Financial Summaries & Spending Review ] ---")

	// 1. Overall Balance Summary
	summary, err := ui.svc.GetBalanceSummary()
	if err != nil {
		fmt.Printf("[X] Error: %v\n", err)
		return
	}
	PrintBalanceSummary(summary)

	// 2. Category Spending Breakdown
	catExpenses, err := ui.svc.GetExpenseBreakdownByCategory()
	if err == nil {
		PrintCategoryBreakdown("EXPENSE BREAKDOWN BY CATEGORY", catExpenses)
	}

	// 3. Category Income Breakdown
	catIncome, err := ui.svc.GetIncomeBreakdownByCategory()
	if err == nil && len(catIncome) > 0 {
		PrintCategoryBreakdown("INCOME BREAKDOWN BY CATEGORY", catIncome)
	}

	// 4. Monthly Trend
	monthly, err := ui.svc.GetMonthlyBreakdown()
	if err == nil {
		PrintMonthlyBreakdown(monthly)
	}
}

func (ui *InteractiveCLI) handleUpdateTransaction() {
	fmt.Println("\n--- [ Update Transaction ] ---")
	idStr := ui.readLine("Enter Transaction ID to update: ")
	id, err := strconv.ParseInt(strings.TrimSpace(idStr), 10, 64)
	if err != nil || id <= 0 {
		fmt.Println("[!] Invalid transaction ID.")
		return
	}

	existing, err := ui.svc.GetTransaction(id)
	if err != nil {
		fmt.Printf("[X] %v\n", err)
		return
	}

	fmt.Printf("\nCurrent Record: #%d | Type: %s | Amount: $%s | Cat: %s | Date: %s | Desc: %s\n",
		existing.ID, existing.Type, existing.Amount.String(), existing.Category, existing.Date, existing.Description)
	fmt.Println("(Press Enter on any prompt to retain the current value)")

	// Update Type
	typeInput := ui.readLine(fmt.Sprintf("New Type (1: Expense, 2: Income) [current: %s]: ", existing.Type))
	if typeInput == "1" || strings.ToLower(typeInput) == "expense" {
		existing.Type = model.TypeExpense
	} else if typeInput == "2" || strings.ToLower(typeInput) == "income" {
		existing.Type = model.TypeIncome
	}

	// Update Amount
	amountInput := ui.readLine(fmt.Sprintf("New Amount [current: %s]: ", existing.Amount.String()))
	if strings.TrimSpace(amountInput) != "" {
		val, err := strconv.ParseFloat(strings.TrimSpace(amountInput), 64)
		if err != nil || val <= 0 {
			fmt.Println("[!] Invalid amount format. Skipping amount update.")
		} else {
			existing.Amount = model.NewMoneyFromFloat(val)
		}
	}

	// Update Category
	catInput := ui.readLine(fmt.Sprintf("New Category [current: %s]: ", existing.Category))
	if strings.TrimSpace(catInput) != "" {
		existing.Category = strings.TrimSpace(catInput)
	}

	// Update Date
	dateInput := ui.readLine(fmt.Sprintf("New Date (YYYY-MM-DD) [current: %s]: ", existing.Date))
	if strings.TrimSpace(dateInput) != "" {
		existing.Date = strings.TrimSpace(dateInput)
	}

	// Update Description
	descInput := ui.readLine(fmt.Sprintf("New Description [current: %s]: ", existing.Description))
	if strings.TrimSpace(descInput) != "" {
		existing.Description = strings.TrimSpace(descInput)
	}

	if err := ui.svc.UpdateTransaction(existing); err != nil {
		fmt.Printf("[X] Failed to update: %v\n", err)
		return
	}

	fmt.Printf("[✓] Transaction #%d successfully updated!\n", existing.ID)
}

func (ui *InteractiveCLI) handleDeleteTransaction() {
	fmt.Println("\n--- [ Delete Transaction ] ---")
	idStr := ui.readLine("Enter Transaction ID to delete: ")
	id, err := strconv.ParseInt(strings.TrimSpace(idStr), 10, 64)
	if err != nil || id <= 0 {
		fmt.Println("[!] Invalid transaction ID.")
		return
	}

	existing, err := ui.svc.GetTransaction(id)
	if err != nil {
		fmt.Printf("[X] %v\n", err)
		return
	}

	fmt.Printf("Record found: #%d | Type: %s | Amount: $%s | Cat: %s | Date: %s | Desc: %s\n",
		existing.ID, existing.Type, existing.Amount.String(), existing.Category, existing.Date, existing.Description)

	confirm := ui.readLine("Are you sure you want to permanently delete this transaction? (y/N): ")
	if strings.ToLower(strings.TrimSpace(confirm)) != "y" {
		fmt.Println("Deletion cancelled.")
		return
	}

	if err := ui.svc.DeleteTransaction(id); err != nil {
		fmt.Printf("[X] Failed to delete: %v\n", err)
		return
	}

	fmt.Printf("[✓] Transaction #%d deleted successfully.\n", id)
}

func (ui *InteractiveCLI) handleListCategories() {
	fmt.Println("\n--- [ Distinct Categories ] ---")
	categories, err := ui.svc.GetCategories()
	if err != nil {
		fmt.Printf("[X] Error: %v\n", err)
		return
	}
	if len(categories) == 0 {
		fmt.Println("  No recorded categories yet.")
		return
	}
	for i, c := range categories {
		fmt.Printf("  %d. %s\n", i+1, c)
	}
	fmt.Println()
}

func (ui *InteractiveCLI) handleExport() {
	fmt.Println("\n--- [ Export Transactions ] ---")
	fmt.Println(" [1] Export to CSV (e.g. expenses.csv)")
	fmt.Println(" [2] Export to JSON (e.g. expenses.json)")
	fmt.Println(" [0] Cancel")

	formatChoice := ui.readLine("Select format (1 or 2): ")
	var fileName string
	var isCSV bool

	switch formatChoice {
	case "1":
		isCSV = true
		fileName = ui.readLine("File name [default: expenses.csv]: ")
		if fileName == "" {
			fileName = "expenses.csv"
		}
	case "2":
		isCSV = false
		fileName = ui.readLine("File name [default: expenses.json]: ")
		if fileName == "" {
			fileName = "expenses.json"
		}
	default:
		fmt.Println("Export cancelled.")
		return
	}

	file, err := os.Create(fileName)
	if err != nil {
		fmt.Printf("[X] Could not create file %s: %v\n", fileName, err)
		return
	}
	defer file.Close()

	if isCSV {
		err = ui.svc.ExportCSV(file)
	} else {
		err = ui.svc.ExportJSON(file)
	}

	if err != nil {
		fmt.Printf("[X] Export failed: %v\n", err)
		return
	}

	fmt.Printf("[✓] Data successfully exported to '%s'!\n", fileName)
}

func (ui *InteractiveCLI) handleSeed() {
	fmt.Println("\n--- [ Load Demo Data ] ---")
	confirm := ui.readLine("This will insert sample student transactions across 3 months. Continue? (y/N): ")
	if strings.ToLower(strings.TrimSpace(confirm)) != "y" {
		fmt.Println("Cancelled.")
		return
	}

	if err := ui.svc.SeedSampleData(); err != nil {
		fmt.Printf("[X] Failed to load sample data: %v\n", err)
		return
	}

	fmt.Println("[✓] Sample demo data successfully loaded! Select option [2] to view records or [5] to see the charts.")
}

func (ui *InteractiveCLI) handleGenerateReport() {
	fmt.Println("\n--- [ Generate Web Dashboard Report ] ---")
	fileName := "expense_report.html"
	f, err := os.Create(fileName)
	if err != nil {
		fmt.Printf("[X] Could not create report file: %v\n", err)
		return
	}
	defer f.Close()

	if err := ui.svc.GenerateHTMLReport(f); err != nil {
		fmt.Printf("[X] Report generation failed: %v\n", err)
		return
	}

	fmt.Printf("[✓] Report generated successfully: %s\n", fileName)
	fmt.Println("Opening in your web browser...")
	service.OpenBrowser(fileName)
}

func (ui *InteractiveCLI) readLine(prompt string) string {
	fmt.Print(prompt)
	text, _ := ui.reader.ReadString('\n')
	return strings.TrimSpace(text)
}

