package service

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"strconv"

	"smart-expense-tracker/internal/model"
	"smart-expense-tracker/internal/storage"
	"smart-expense-tracker/internal/validator"
)

// TrackerService defines the high-level business logic operations.
type TrackerService struct {
	repo storage.Repository
}

// NewTrackerService instantiates a new service layer with a given repository.
func NewTrackerService(repo storage.Repository) *TrackerService {
	return &TrackerService{repo: repo}
}

// AddTransaction validates and records a new financial transaction.
func (s *TrackerService) AddTransaction(tx *model.Transaction) error {
	if err := validator.ValidateTransaction(tx); err != nil {
		return err
	}
	if err := s.repo.Create(tx); err != nil {
		return fmt.Errorf("could not record transaction: %w", err)
	}
	return nil
}

// GetTransaction retrieves a transaction by ID after validation.
func (s *TrackerService) GetTransaction(id int64) (*model.Transaction, error) {
	if err := validator.ValidateID(id); err != nil {
		return nil, err
	}
	tx, err := s.repo.GetByID(id)
	if err != nil {
		if err == storage.ErrNotFound {
			return nil, fmt.Errorf("transaction with ID %d not found", id)
		}
		return nil, fmt.Errorf("error reading transaction: %w", err)
	}
	return tx, nil
}

// ListTransactions retrieves transactions matching the optional filter criteria.
func (s *TrackerService) ListTransactions(filter model.Filter) ([]model.Transaction, error) {
	list, err := s.repo.List(filter)
	if err != nil {
		return nil, fmt.Errorf("error listing transactions: %w", err)
	}
	return list, nil
}

// UpdateTransaction validates and updates an existing transaction record.
func (s *TrackerService) UpdateTransaction(tx *model.Transaction) error {
	if err := validator.ValidateID(tx.ID); err != nil {
		return err
	}
	if err := validator.ValidateTransaction(tx); err != nil {
		return err
	}
	if err := s.repo.Update(tx); err != nil {
		if err == storage.ErrNotFound {
			return fmt.Errorf("transaction with ID %d not found to update", tx.ID)
		}
		return fmt.Errorf("error updating transaction: %w", err)
	}
	return nil
}

// DeleteTransaction removes a transaction record by its ID.
func (s *TrackerService) DeleteTransaction(id int64) error {
	if err := validator.ValidateID(id); err != nil {
		return err
	}
	if err := s.repo.Delete(id); err != nil {
		if err == storage.ErrNotFound {
			return fmt.Errorf("transaction with ID %d not found to delete", id)
		}
		return fmt.Errorf("error deleting transaction: %w", err)
	}
	return nil
}

// GetBalanceSummary calculates total income, total expenses, and the remaining net balance.
// Formula: Remaining Balance = Total Income - Total Expenses
func (s *TrackerService) GetBalanceSummary() (*model.BalanceSummary, error) {
	income, expense, count, err := s.repo.GetTotals()
	if err != nil {
		return nil, fmt.Errorf("error calculating balance summary: %w", err)
	}

	balance := income - expense

	return &model.BalanceSummary{
		TotalIncome:  income,
		TotalExpense: expense,
		Balance:      balance,
		Count:        count,
	}, nil
}

// GetExpenseBreakdownByCategory retrieves expenses aggregated by category.
func (s *TrackerService) GetExpenseBreakdownByCategory() ([]model.CategoryBreakdown, error) {
	breakdown, err := s.repo.GetCategoryBreakdown(model.TypeExpense)
	if err != nil {
		return nil, fmt.Errorf("error calculating expense breakdown: %w", err)
	}
	return breakdown, nil
}

// GetIncomeBreakdownByCategory retrieves income aggregated by category.
func (s *TrackerService) GetIncomeBreakdownByCategory() ([]model.CategoryBreakdown, error) {
	breakdown, err := s.repo.GetCategoryBreakdown(model.TypeIncome)
	if err != nil {
		return nil, fmt.Errorf("error calculating income breakdown: %w", err)
	}
	return breakdown, nil
}

// GetMonthlyBreakdown retrieves monthly trends for income and expenses.
func (s *TrackerService) GetMonthlyBreakdown() ([]model.MonthlyBreakdown, error) {
	breakdown, err := s.repo.GetMonthlyBreakdown()
	if err != nil {
		return nil, fmt.Errorf("error calculating monthly trends: %w", err)
	}
	return breakdown, nil
}

// GetCategories returns all distinct categories currently stored.
func (s *TrackerService) GetCategories() ([]string, error) {
	categories, err := s.repo.GetDistinctCategories()
	if err != nil {
		return nil, fmt.Errorf("error fetching categories: %w", err)
	}
	return categories, nil
}

// ExportCSV writes all recorded transactions to a CSV format stream.
func (s *TrackerService) ExportCSV(w io.Writer) error {
	txs, err := s.repo.List(model.Filter{})
	if err != nil {
		return fmt.Errorf("failed to fetch transactions for export: %w", err)
	}

	csvWriter := csv.NewWriter(w)
	defer csvWriter.Flush()

	// Write header
	header := []string{"ID", "Type", "Amount", "Category", "Date", "Description"}
	if err := csvWriter.Write(header); err != nil {
		return fmt.Errorf("failed to write CSV header: %w", err)
	}

	// Write rows
	for _, tx := range txs {
		record := []string{
			strconv.FormatInt(tx.ID, 10),
			string(tx.Type),
			tx.Amount.String(),
			tx.Category,
			tx.Date,
			tx.Description,
		}
		if err := csvWriter.Write(record); err != nil {
			return fmt.Errorf("failed to write CSV row: %w", err)
		}
	}

	return nil
}

// ExportJSON writes all transactions formatted as pretty JSON.
func (s *TrackerService) ExportJSON(w io.Writer) error {
	txs, err := s.repo.List(model.Filter{})
	if err != nil {
		return fmt.Errorf("failed to fetch transactions for export: %w", err)
	}

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(txs); err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}
	return nil
}

// SeedSampleData populates a set of realistic student financial transactions
// across multiple months and categories for demonstration and grading.
func (s *TrackerService) SeedSampleData() error {
	sampleTransactions := []model.Transaction{
		// Month 1
		{Type: model.TypeIncome, Amount: model.NewMoneyFromFloat(2000.00), Category: "Scholarship", Date: "2026-09-01", Description: "Semester merit stipend"},
		{Type: model.TypeExpense, Amount: model.NewMoneyFromFloat(450.00), Category: "Rent", Date: "2026-09-02", Description: "Hostel room rent"},
		{Type: model.TypeExpense, Amount: model.NewMoneyFromFloat(120.00), Category: "Education", Date: "2026-09-05", Description: "Network Security textbooks"},
		{Type: model.TypeExpense, Amount: model.NewMoneyFromFloat(85.50), Category: "Food", Date: "2026-09-10", Description: "Campus canteen & grocery"},
		{Type: model.TypeExpense, Amount: model.NewMoneyFromFloat(35.00), Category: "Transport", Date: "2026-09-15", Description: "Bus & metro pass"},
		// Month 2
		{Type: model.TypeIncome, Amount: model.NewMoneyFromFloat(2000.00), Category: "Scholarship", Date: "2026-10-01", Description: "Monthly stipend"},
		{Type: model.TypeIncome, Amount: model.NewMoneyFromFloat(350.00), Category: "Freelance", Date: "2026-10-04", Description: "Bug bounty / code audit"},
		{Type: model.TypeExpense, Amount: model.NewMoneyFromFloat(450.00), Category: "Rent", Date: "2026-10-02", Description: "Hostel room rent"},
		{Type: model.TypeExpense, Amount: model.NewMoneyFromFloat(110.25), Category: "Food", Date: "2026-10-12", Description: "Monthly mess & groceries"},
		{Type: model.TypeExpense, Amount: model.NewMoneyFromFloat(60.00), Category: "Utilities", Date: "2026-10-18", Description: "High-speed internet plan"},
		{Type: model.TypeExpense, Amount: model.NewMoneyFromFloat(75.00), Category: "Shopping", Date: "2026-10-25", Description: "Electronics components & USB"},
		// Month 3
		{Type: model.TypeIncome, Amount: model.NewMoneyFromFloat(2200.00), Category: "Scholarship", Date: "2026-11-01", Description: "Monthly stipend"},
		{Type: model.TypeExpense, Amount: model.NewMoneyFromFloat(450.00), Category: "Rent", Date: "2026-11-02", Description: "Hostel room rent"},
		{Type: model.TypeExpense, Amount: model.NewMoneyFromFloat(95.00), Category: "Food", Date: "2026-11-04", Description: "Weekly food supplies"},
		{Type: model.TypeExpense, Amount: model.NewMoneyFromFloat(40.00), Category: "Transport", Date: "2026-11-07", Description: "Metro travel"},
		{Type: model.TypeExpense, Amount: model.NewMoneyFromFloat(150.00), Category: "Education", Date: "2026-11-09", Description: "Certification exam fee"},
	}

	for _, tx := range sampleTransactions {
		if err := s.AddTransaction(&tx); err != nil {
			return err
		}
	}

	return nil
}

