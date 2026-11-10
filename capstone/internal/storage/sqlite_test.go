package storage

import (
	"path/filepath"
	"testing"

	"smart-expense-tracker/internal/model"
)

func setupTestDB(t *testing.T) *SQLiteRepository {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test_expenses.db")
	repo, err := NewSQLiteRepository(dbPath)
	if err != nil {
		t.Fatalf("failed to create test repo: %v", err)
	}
	t.Cleanup(func() {
		repo.Close()
	})
	return repo
}

func TestSQLiteRepository_CreateAndGet(t *testing.T) {
	repo := setupTestDB(t)

	tx := &model.Transaction{
		Type:        model.TypeExpense,
		Amount:      model.NewMoneyFromFloat(25.50),
		Category:    "Food",
		Date:        "2026-11-10",
		Description: "Burgers",
	}

	err := repo.Create(tx)
	if err != nil {
		t.Fatalf("failed to create transaction: %v", err)
	}
	if tx.ID <= 0 {
		t.Fatalf("expected positive ID, got %d", tx.ID)
	}

	fetched, err := repo.GetByID(tx.ID)
	if err != nil {
		t.Fatalf("failed to fetch transaction: %v", err)
	}

	if fetched.Type != tx.Type || fetched.Amount != tx.Amount || fetched.Category != tx.Category {
		t.Errorf("fetched transaction mismatch: got %+v, expected %+v", fetched, tx)
	}
}

func TestSQLiteRepository_UpdateAndDelete(t *testing.T) {
	repo := setupTestDB(t)

	tx := &model.Transaction{
		Type:        model.TypeExpense,
		Amount:      model.NewMoneyFromFloat(100.00),
		Category:    "Books",
		Date:        "2026-11-01",
		Description: "Textbook",
	}

	if err := repo.Create(tx); err != nil {
		t.Fatalf("create failed: %v", err)
	}

	// Update
	tx.Amount = model.NewMoneyFromFloat(120.00)
	tx.Category = "Education"
	if err := repo.Update(tx); err != nil {
		t.Fatalf("update failed: %v", err)
	}

	updated, err := repo.GetByID(tx.ID)
	if err != nil {
		t.Fatalf("get updated failed: %v", err)
	}
	if updated.Amount != model.NewMoneyFromFloat(120.00) || updated.Category != "Education" {
		t.Errorf("update was not saved properly: %+v", updated)
	}

	// Delete
	if err := repo.Delete(tx.ID); err != nil {
		t.Fatalf("delete failed: %v", err)
	}

	_, err = repo.GetByID(tx.ID)
	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound after deletion, got: %v", err)
	}
}

func TestSQLiteRepository_TotalsAndBreakdowns(t *testing.T) {
	repo := setupTestDB(t)

	// Add income
	income := &model.Transaction{
		Type:     model.TypeIncome,
		Amount:   model.NewMoneyFromFloat(1000.00),
		Category: "Allowance",
		Date:     "2026-11-01",
	}
	_ = repo.Create(income)

	// Add expenses
	exp1 := &model.Transaction{
		Type:     model.TypeExpense,
		Amount:   model.NewMoneyFromFloat(200.00),
		Category: "Food",
		Date:     "2026-11-02",
	}
	exp2 := &model.Transaction{
		Type:     model.TypeExpense,
		Amount:   model.NewMoneyFromFloat(300.00),
		Category: "Transport",
		Date:     "2026-11-03",
	}
	_ = repo.Create(exp1)
	_ = repo.Create(exp2)

	totIncome, totExpense, count, err := repo.GetTotals()
	if err != nil {
		t.Fatalf("failed to get totals: %v", err)
	}

	if count != 3 {
		t.Errorf("expected count 3, got %d", count)
	}
	if totIncome != model.NewMoneyFromFloat(1000.00) {
		t.Errorf("expected income 1000.00, got %v", totIncome)
	}
	if totExpense != model.NewMoneyFromFloat(500.00) {
		t.Errorf("expected expense 500.00, got %v", totExpense)
	}

	// Category breakdown for expenses
	catBreakdown, err := repo.GetCategoryBreakdown(model.TypeExpense)
	if err != nil {
		t.Fatalf("failed to get category breakdown: %v", err)
	}
	if len(catBreakdown) != 2 {
		t.Fatalf("expected 2 categories, got %d", len(catBreakdown))
	}
	// Transport is 300 (60%), Food is 200 (40%)
	if catBreakdown[0].Category != "Transport" || catBreakdown[0].Total != model.NewMoneyFromFloat(300.00) {
		t.Errorf("unexpected top category: %+v", catBreakdown[0])
	}
	if catBreakdown[0].Percentage != 60.0 {
		t.Errorf("expected 60%%, got %.2f%%", catBreakdown[0].Percentage)
	}

	// Monthly breakdown
	monthly, err := repo.GetMonthlyBreakdown()
	if err != nil {
		t.Fatalf("failed to get monthly breakdown: %v", err)
	}
	if len(monthly) != 1 {
		t.Fatalf("expected 1 month record, got %d", len(monthly))
	}
	if monthly[0].Month != "2026-11" || monthly[0].Balance != model.NewMoneyFromFloat(500.00) {
		t.Errorf("unexpected monthly breakdown: %+v", monthly[0])
	}
}

