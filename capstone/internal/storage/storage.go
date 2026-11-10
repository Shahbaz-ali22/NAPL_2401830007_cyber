package storage

import (
	"errors"

	"smart-expense-tracker/internal/model"
)

var (
	ErrNotFound = errors.New("record not found")
)

// Repository defines the interface for persisting and retrieving transactions.
type Repository interface {
	// Create inserts a new transaction.
	Create(t *model.Transaction) error

	// GetByID retrieves a single transaction by ID.
	GetByID(id int64) (*model.Transaction, error)

	// List returns transactions matching optional filter criteria.
	List(filter model.Filter) ([]model.Transaction, error)

	// Update modifies an existing transaction by ID.
	Update(t *model.Transaction) error

	// Delete removes a transaction by ID.
	Delete(id int64) error

	// GetTotals returns aggregate income, expense, and total record count.
	GetTotals() (income model.Money, expense model.Money, count int, err error)

	// GetCategoryBreakdown returns spending or income grouped by category.
	GetCategoryBreakdown(txType model.TransactionType) ([]model.CategoryBreakdown, error)

	// GetMonthlyBreakdown returns income/expense summary grouped by YYYY-MM.
	GetMonthlyBreakdown() ([]model.MonthlyBreakdown, error)

	// GetDistinctCategories returns all distinct categories recorded.
	GetDistinctCategories() ([]string, error)

	// Close cleanly closes the storage connection.
	Close() error
}

