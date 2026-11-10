package validator

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"smart-expense-tracker/internal/model"
)

const (
	DateFormat       = "2006-01-02"
	MaxDescriptionLen = 255
	MaxCategoryLen    = 50
	MaxAmountCents    = 100_000_000_00 // 100 million
)

// ValidationError represents an input validation failure.
type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("validation error on '%s': %s", e.Field, e.Message)
}

// ValidationErrors collects multiple validation issues.
type ValidationErrors []ValidationError

func (ve ValidationErrors) Error() string {
	var msgs []string
	for _, err := range ve {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// ValidateTransaction validates a transaction before creation or updating.
func ValidateTransaction(t *model.Transaction) error {
	var errs ValidationErrors

	// 1. Validate Transaction Type
	if t.Type != model.TypeIncome && t.Type != model.TypeExpense {
		errs = append(errs, ValidationError{
			Field:   "type",
			Message: "must be either 'income' or 'expense'",
		})
	}

	// 2. Validate Amount
	if t.Amount <= 0 {
		errs = append(errs, ValidationError{
			Field:   "amount",
			Message: "must be greater than 0.00",
		})
	} else if t.Amount > MaxAmountCents {
		errs = append(errs, ValidationError{
			Field:   "amount",
			Message: "exceeds maximum allowable transaction limit",
		})
	}

	// 3. Validate Category
	trimmedCategory := strings.TrimSpace(t.Category)
	if trimmedCategory == "" {
		errs = append(errs, ValidationError{
			Field:   "category",
			Message: "cannot be empty",
		})
	} else if len(trimmedCategory) > MaxCategoryLen {
		errs = append(errs, ValidationError{
			Field:   "category",
			Message: fmt.Sprintf("cannot exceed %d characters", MaxCategoryLen),
		})
	} else {
		t.Category = trimmedCategory
	}

	// 4. Validate Date
	trimmedDate := strings.TrimSpace(t.Date)
	if trimmedDate == "" {
		// Default to current date if omitted
		t.Date = time.Now().Format(DateFormat)
	} else {
		parsedDate, err := time.Parse(DateFormat, trimmedDate)
		if err != nil {
			errs = append(errs, ValidationError{
				Field:   "date",
				Message: "invalid date format; expected YYYY-MM-DD (e.g. 2026-11-10)",
			})
		} else {
			// Date cannot be unreasonably far in the future (e.g. > 10 years)
			if parsedDate.Year() < 1900 || parsedDate.Year() > time.Now().Year()+10 {
				errs = append(errs, ValidationError{
					Field:   "date",
					Message: "date must be between 1900 and 10 years in the future",
				})
			}
			t.Date = parsedDate.Format(DateFormat)
		}
	}

	// 5. Validate Description (optional, but bounded)
	trimmedDesc := strings.TrimSpace(t.Description)
	if len(trimmedDesc) > MaxDescriptionLen {
		errs = append(errs, ValidationError{
			Field:   "description",
			Message: fmt.Sprintf("cannot exceed %d characters", MaxDescriptionLen),
		})
	} else {
		t.Description = trimmedDesc
	}

	if len(errs) > 0 {
		return errs
	}
	return nil
}

// ValidateID checks if a database ID is valid.
func ValidateID(id int64) error {
	if id <= 0 {
		return errors.New("transaction ID must be a positive integer greater than 0")
	}
	return nil
}

