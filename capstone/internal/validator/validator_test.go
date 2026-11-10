package validator

import (
	"strings"
	"testing"

	"smart-expense-tracker/internal/model"
)

func TestValidateTransaction_Valid(t *testing.T) {
	tx := &model.Transaction{
		Type:        model.TypeExpense,
		Amount:      model.NewMoneyFromFloat(45.50),
		Category:    "Food",
		Date:        "2026-11-10",
		Description: "Lunch at cafe",
	}

	if err := ValidateTransaction(tx); err != nil {
		t.Fatalf("expected valid transaction, got error: %v", err)
	}
}

func TestValidateTransaction_DefaultDate(t *testing.T) {
	tx := &model.Transaction{
		Type:     model.TypeIncome,
		Amount:   model.NewMoneyFromFloat(1500.00),
		Category: "Scholarship",
		Date:     "",
	}

	if err := ValidateTransaction(tx); err != nil {
		t.Fatalf("expected date to default without error, got: %v", err)
	}

	if tx.Date == "" {
		t.Errorf("expected date to be auto-populated with today's date")
	}
}

func TestValidateTransaction_InvalidFields(t *testing.T) {
	tests := []struct {
		name        string
		tx          model.Transaction
		errContains string
	}{
		{
			name: "invalid type",
			tx: model.Transaction{
				Type:     "unknown",
				Amount:   model.NewMoneyFromFloat(10),
				Category: "Food",
				Date:     "2026-11-10",
			},
			errContains: "type",
		},
		{
			name: "zero amount",
			tx: model.Transaction{
				Type:     model.TypeExpense,
				Amount:   0,
				Category: "Food",
				Date:     "2026-11-10",
			},
			errContains: "amount",
		},
		{
			name: "negative amount",
			tx: model.Transaction{
				Type:     model.TypeExpense,
				Amount:   -500,
				Category: "Food",
				Date:     "2026-11-10",
			},
			errContains: "amount",
		},
		{
			name: "empty category",
			tx: model.Transaction{
				Type:     model.TypeExpense,
				Amount:   model.NewMoneyFromFloat(10),
				Category: "   ",
				Date:     "2026-11-10",
			},
			errContains: "category",
		},
		{
			name: "invalid date format",
			tx: model.Transaction{
				Type:     model.TypeExpense,
				Amount:   model.NewMoneyFromFloat(10),
				Category: "Food",
				Date:     "10/11/2026",
			},
			errContains: "date",
		},
		{
			name: "description too long",
			tx: model.Transaction{
				Type:        model.TypeExpense,
				Amount:      model.NewMoneyFromFloat(10),
				Category:    "Food",
				Date:        "2026-11-10",
				Description: strings.Repeat("A", 300),
			},
			errContains: "description",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTransaction(&tt.tx)
			if err == nil {
				t.Fatalf("expected error containing %q, got nil", tt.errContains)
			}
			if !strings.Contains(err.Error(), tt.errContains) {
				t.Errorf("expected error to contain %q, got %q", tt.errContains, err.Error())
			}
		})
	}
}

func TestValidateID(t *testing.T) {
	if err := ValidateID(1); err != nil {
		t.Errorf("expected id 1 to be valid, got: %v", err)
	}

	if err := ValidateID(0); err == nil {
		t.Errorf("expected id 0 to be invalid")
	}

	if err := ValidateID(-5); err == nil {
		t.Errorf("expected negative id to be invalid")
	}
}

