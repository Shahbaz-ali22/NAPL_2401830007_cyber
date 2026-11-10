package model

import (
	"testing"
)

func TestMoney_ConversionAndFormatting(t *testing.T) {
	tests := []struct {
		input    float64
		expected Money
		str      string
	}{
		{input: 0.00, expected: 0, str: "0.00"},
		{input: 12.34, expected: 1234, str: "12.34"},
		{input: 99.999, expected: 10000, str: "100.00"},
		{input: 1500.50, expected: 150050, str: "1500.50"},
	}

	for _, tt := range tests {
		m := NewMoneyFromFloat(tt.input)
		if m != tt.expected {
			t.Errorf("NewMoneyFromFloat(%.3f) = %d, expected %d", tt.input, m, tt.expected)
		}
		if m.String() != tt.str {
			t.Errorf("m.String() = %s, expected %s", m.String(), tt.str)
		}
	}
}

func TestParseTransactionType(t *testing.T) {
	cases := []struct {
		input       string
		expected    TransactionType
		expectError bool
	}{
		{"income", TypeIncome, false},
		{"INCOME", TypeIncome, false},
		{" Expense ", TypeExpense, false},
		{"invalid", "", true},
		{"", "", true},
	}

	for _, c := range cases {
		res, err := ParseTransactionType(c.input)
		if (err != nil) != c.expectError {
			t.Errorf("ParseTransactionType(%q) unexpected error state: %v", c.input, err)
		}
		if res != c.expected {
			t.Errorf("ParseTransactionType(%q) = %q, expected %q", c.input, res, c.expected)
		}
	}
}

