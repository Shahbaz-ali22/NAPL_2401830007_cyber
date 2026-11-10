package ui

import (
	"fmt"
	"strings"

	"smart-expense-tracker/internal/model"
)

// PrintTransactionsTable prints a formatted table of transactions.
func PrintTransactionsTable(txs []model.Transaction) {
	if len(txs) == 0 {
		fmt.Println("\n  [i] No transactions found matching criteria.")
		return
	}

	fmt.Println()
	fmt.Println("+" + strings.Repeat("-", 7) + "+" + strings.Repeat("-", 10) + "+" + strings.Repeat("-", 12) + "+" + strings.Repeat("-", 16) + "+" + strings.Repeat("-", 13) + "+" + strings.Repeat("-", 32) + "+")
	fmt.Printf("| %-5s | %-8s | %10s | %-14s | %-11s | %-30s |\n", "ID", "TYPE", "AMOUNT", "CATEGORY", "DATE", "DESCRIPTION")
	fmt.Println("+" + strings.Repeat("-", 7) + "+" + strings.Repeat("-", 10) + "+" + strings.Repeat("-", 12) + "+" + strings.Repeat("-", 16) + "+" + strings.Repeat("-", 13) + "+" + strings.Repeat("-", 32) + "+")

	for _, tx := range txs {
		desc := tx.Description
		if len(desc) > 30 {
			desc = desc[:27] + "..."
		}
		typeFormatted := strings.ToUpper(string(tx.Type))
		fmt.Printf("| %-5d | %-8s | %10s | %-14s | %-11s | %-30s |\n",
			tx.ID,
			typeFormatted,
			tx.Amount.String(),
			tx.Category,
			tx.Date,
			desc,
		)
	}
	fmt.Println("+" + strings.Repeat("-", 7) + "+" + strings.Repeat("-", 10) + "+" + strings.Repeat("-", 12) + "+" + strings.Repeat("-", 16) + "+" + strings.Repeat("-", 13) + "+" + strings.Repeat("-", 32) + "+")
	fmt.Printf("  Total Records: %d\n\n", len(txs))
}

// PrintBalanceSummary prints the overall balance sheet.
func PrintBalanceSummary(summary *model.BalanceSummary) {
	fmt.Println()
	fmt.Println("==================================================")
	fmt.Println("             FINANCIAL BALANCE SUMMARY            ")
	fmt.Println("==================================================")
	fmt.Printf(" Total Transactions: %d\n", summary.Count)
	fmt.Printf(" Total Income:       $%s\n", summary.TotalIncome.String())
	fmt.Printf(" Total Expenses:     $%s\n", summary.TotalExpense.String())
	fmt.Println("--------------------------------------------------")
	if summary.Balance >= 0 {
		fmt.Printf(" Net Balance:        +$%s (Surplus)\n", summary.Balance.String())
	} else {
		fmt.Printf(" Net Balance:        -$%s (Deficit)\n", (-summary.Balance).String())
	}
	fmt.Println("==================================================")
	fmt.Println(" Formula: Remaining Balance = Total Income - Total Expenses")
	fmt.Println()
}

// PrintCategoryBreakdown prints category-wise distribution with visual bar representation.
func PrintCategoryBreakdown(title string, list []model.CategoryBreakdown) {
	fmt.Println()
	fmt.Printf("================= %s =================\n", title)
	if len(list) == 0 {
		fmt.Println("  No data available for category analysis.")
		return
	}

	fmt.Println("+" + strings.Repeat("-", 18) + "+" + strings.Repeat("-", 12) + "+" + strings.Repeat("-", 9) + "+" + strings.Repeat("-", 26) + "+")
	fmt.Printf("| %-16s | %10s | %7s | %-24s |\n", "CATEGORY", "TOTAL", "SHARE", "DISTRIBUTION")
	fmt.Println("+" + strings.Repeat("-", 18) + "+" + strings.Repeat("-", 12) + "+" + strings.Repeat("-", 9) + "+" + strings.Repeat("-", 26) + "+")

	for _, item := range list {
		barLen := int((item.Percentage / 100.0) * 20.0)
		if barLen < 1 && item.Percentage > 0 {
			barLen = 1
		}
		bar := strings.Repeat("#", barLen)

		fmt.Printf("| %-16s | %10s | %6.1f%% | %-24s |\n",
			item.Category,
			item.Total.String(),
			item.Percentage,
			bar,
		)
	}
	fmt.Println("+" + strings.Repeat("-", 18) + "+" + strings.Repeat("-", 12) + "+" + strings.Repeat("-", 9) + "+" + strings.Repeat("-", 26) + "+")
	fmt.Println()
}

// PrintMonthlyBreakdown prints monthly breakdown report.
func PrintMonthlyBreakdown(list []model.MonthlyBreakdown) {
	fmt.Println()
	fmt.Println("================== MONTHLY FINANCIAL TRENDS ==================")
	if len(list) == 0 {
		fmt.Println("  No historical monthly records found.")
		return
	}

	fmt.Println("+" + strings.Repeat("-", 10) + "+" + strings.Repeat("-", 14) + "+" + strings.Repeat("-", 14) + "+" + strings.Repeat("-", 16) + "+")
	fmt.Printf("| %-8s | %12s | %12s | %14s |\n", "MONTH", "INCOME", "EXPENSES", "NET BALANCE")
	fmt.Println("+" + strings.Repeat("-", 10) + "+" + strings.Repeat("-", 14) + "+" + strings.Repeat("-", 14) + "+" + strings.Repeat("-", 16) + "+")

	for _, item := range list {
		balPrefix := ""
		if item.Balance > 0 {
			balPrefix = "+"
		}
		fmt.Printf("| %-8s | %12s | %12s | %14s |\n",
			item.Month,
			item.Income.String(),
			item.Expense.String(),
			balPrefix+item.Balance.String(),
		)
	}
	fmt.Println("+" + strings.Repeat("-", 10) + "+" + strings.Repeat("-", 14) + "+" + strings.Repeat("-", 14) + "+" + strings.Repeat("-", 16) + "+")
	fmt.Println()
}

