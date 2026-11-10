package model

import (
	"fmt"
	"math"
	"strings"
	"time"
)

// TransactionType defines whether a transaction is income or expense.
type TransactionType string

const (
	TypeIncome  TransactionType = "income"
	TypeExpense TransactionType = "expense"
)

// ParseTransactionType converts a string to TransactionType with case-insensitivity.
func ParseTransactionType(s string) (TransactionType, error) {
	clean := strings.ToLower(strings.TrimSpace(s))
	switch clean {
	case string(TypeIncome):
		return TypeIncome, nil
	case string(TypeExpense):
		return TypeExpense, nil
	default:
		return "", fmt.Errorf("invalid transaction type: '%s'. Allowed values: income, expense", s)
	}
}

// Money represents monetary amounts in the smallest currency unit (e.g. cents/paise)
// to prevent floating-point rounding inaccuracies as specified in the Phase 1 design.
type Money int64

// NewMoneyFromFloat converts a float value (e.g. 10.50) to Money (1050 cents).
func NewMoneyFromFloat(f float64) Money {
	// Round to nearest integer to avoid precision issues
	return Money(math.Round(f * 100))
}

// ToFloat converts Money back to float64 for display purposes.
func (m Money) ToFloat() float64 {
	return float64(m) / 100.0
}

// String formats the Money value into standard decimal notation (e.g. "10.50").
func (m Money) String() string {
	return fmt.Sprintf("%.2f", m.ToFloat())
}

// Transaction represents a financial record in the expense tracker.
type Transaction struct {
	ID          int64           `json:"id"`
	Type        TransactionType `json:"type"`
	Amount      Money           `json:"amount"` // Stored in minor currency units
	Category    string          `json:"category"`
	Date        string          `json:"date"` // Format: YYYY-MM-DD
	Description string          `json:"description"`
	CreatedAt   time.Time       `json:"created_at"`
}

// Filter allows filtering transactions by various criteria.
type Filter struct {
	Type      TransactionType
	Category  string
	StartDate string // YYYY-MM-DD
	EndDate   string // YYYY-MM-DD
}

// BalanceSummary contains total calculations.
type BalanceSummary struct {
	TotalIncome  Money
	TotalExpense Money
	Balance      Money // TotalIncome - TotalExpense
	Count        int
}

// CategoryBreakdown represents total spending in a category.
type CategoryBreakdown struct {
	Category   string
	Total      Money
	Percentage float64
	Count      int
}

// MonthlyBreakdown represents spending/income grouped by month.
type MonthlyBreakdown struct {
	Month   string // YYYY-MM
	Income  Money
	Expense Money
	Balance Money
}

