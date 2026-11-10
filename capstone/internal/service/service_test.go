package service

import (
	"path/filepath"
	"strings"
	"testing"

	"smart-expense-tracker/internal/model"
	"smart-expense-tracker/internal/storage"
)

func setupTestService(t *testing.T) *TrackerService {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "service_test.db")
	repo, err := storage.NewSQLiteRepository(dbPath)
	if err != nil {
		t.Fatalf("failed to init repo: %v", err)
	}
	t.Cleanup(func() {
		repo.Close()
	})
	return NewTrackerService(repo)
}

func TestTrackerService_AddAndGet(t *testing.T) {
	svc := setupTestService(t)

	tx := &model.Transaction{
		Type:        model.TypeExpense,
		Amount:      model.NewMoneyFromFloat(50.00),
		Category:    "Food",
		Date:        "2026-11-10",
		Description: "Groceries",
	}

	err := svc.AddTransaction(tx)
	if err != nil {
		t.Fatalf("AddTransaction failed: %v", err)
	}

	fetched, err := svc.GetTransaction(tx.ID)
	if err != nil {
		t.Fatalf("GetTransaction failed: %v", err)
	}

	if fetched.Amount != tx.Amount || fetched.Category != "Food" {
		t.Errorf("mismatched transaction data: %+v", fetched)
	}
}

func TestTrackerService_InvalidInput(t *testing.T) {
	svc := setupTestService(t)

	// Zero amount
	tx := &model.Transaction{
		Type:     model.TypeExpense,
		Amount:   0,
		Category: "Food",
		Date:     "2026-11-10",
	}

	err := svc.AddTransaction(tx)
	if err == nil {
		t.Fatalf("expected validation error on zero amount, got nil")
	}

	// Invalid ID lookup
	_, err = svc.GetTransaction(-1)
	if err == nil {
		t.Fatalf("expected error on negative ID, got nil")
	}
}

func TestTrackerService_BalanceCalculation(t *testing.T) {
	svc := setupTestService(t)

	// Income: 2000.00
	_ = svc.AddTransaction(&model.Transaction{
		Type:     model.TypeIncome,
		Amount:   model.NewMoneyFromFloat(2000.00),
		Category: "Salary",
		Date:     "2026-11-01",
	})

	// Expense 1: 500.00 (Rent)
	_ = svc.AddTransaction(&model.Transaction{
		Type:     model.TypeExpense,
		Amount:   model.NewMoneyFromFloat(500.00),
		Category: "Rent",
		Date:     "2026-11-02",
	})

	// Expense 2: 150.00 (Food)
	_ = svc.AddTransaction(&model.Transaction{
		Type:     model.TypeExpense,
		Amount:   model.NewMoneyFromFloat(150.00),
		Category: "Food",
		Date:     "2026-11-03",
	})

	summary, err := svc.GetBalanceSummary()
	if err != nil {
		t.Fatalf("GetBalanceSummary failed: %v", err)
	}

	// Formula: Remaining Balance = Total Income - Total Expenses
	// Balance = 2000 - 650 = 1350
	if summary.TotalIncome != model.NewMoneyFromFloat(2000.00) {
		t.Errorf("expected TotalIncome 2000.00, got %v", summary.TotalIncome)
	}
	if summary.TotalExpense != model.NewMoneyFromFloat(650.00) {
		t.Errorf("expected TotalExpense 650.00, got %v", summary.TotalExpense)
	}
	if summary.Balance != model.NewMoneyFromFloat(1350.00) {
		t.Errorf("expected Balance 1350.00, got %v", summary.Balance)
	}
	if summary.Count != 3 {
		t.Errorf("expected Count 3, got %d", summary.Count)
	}
}

func TestTrackerService_UpdateAndDelete(t *testing.T) {
	svc := setupTestService(t)

	tx := &model.Transaction{
		Type:     model.TypeExpense,
		Amount:   model.NewMoneyFromFloat(100.00),
		Category: "Books",
		Date:     "2026-11-01",
	}
	_ = svc.AddTransaction(tx)

	// Update
	tx.Amount = model.NewMoneyFromFloat(150.00)
	if err := svc.UpdateTransaction(tx); err != nil {
		t.Fatalf("UpdateTransaction failed: %v", err)
	}

	updated, _ := svc.GetTransaction(tx.ID)
	if updated.Amount != model.NewMoneyFromFloat(150.00) {
		t.Errorf("expected updated amount 150.00, got %v", updated.Amount)
	}

	// Delete
	if err := svc.DeleteTransaction(tx.ID); err != nil {
		t.Fatalf("DeleteTransaction failed: %v", err)
	}

	_, err := svc.GetTransaction(tx.ID)
	if err == nil {
		t.Fatalf("expected error getting deleted transaction, got nil")
	}
}

func TestTrackerService_ExportAndSeed(t *testing.T) {
	svc := setupTestService(t)

	// Seed demo data
	if err := svc.SeedSampleData(); err != nil {
		t.Fatalf("SeedSampleData failed: %v", err)
	}

	txs, err := svc.ListTransactions(model.Filter{})
	if err != nil {
		t.Fatalf("ListTransactions failed: %v", err)
	}
	if len(txs) == 0 {
		t.Fatalf("expected seeded transactions, got 0")
	}

	// Test CSV Export
	var csvBuf strings.Builder
	if err := svc.ExportCSV(&csvBuf); err != nil {
		t.Fatalf("ExportCSV failed: %v", err)
	}
	if !strings.Contains(csvBuf.String(), "Scholarship") {
		t.Errorf("expected CSV to contain 'Scholarship'")
	}

	// Test JSON Export
	var jsonBuf strings.Builder
	if err := svc.ExportJSON(&jsonBuf); err != nil {
		t.Fatalf("ExportJSON failed: %v", err)
	}
	if !strings.Contains(jsonBuf.String(), "Scholarship") {
		t.Errorf("expected JSON to contain 'Scholarship'")
	}
}

