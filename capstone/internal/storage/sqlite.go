package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	_ "modernc.org/sqlite"

	"smart-expense-tracker/internal/model"
)

// SQLiteRepository implements Repository using SQLite.
type SQLiteRepository struct {
	db *sql.DB
}

// NewSQLiteRepository creates or opens a SQLite database at dbPath and initializes the schema.
func NewSQLiteRepository(dbPath string) (*SQLiteRepository, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Verify database connection
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	repo := &SQLiteRepository{db: db}
	if err := repo.initSchema(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to initialize database schema: %w", err)
	}

	return repo, nil
}

// initSchema creates required tables and indexes if they do not exist.
func (r *SQLiteRepository) initSchema() error {
	query := `
	CREATE TABLE IF NOT EXISTS transactions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		type TEXT NOT NULL CHECK(type IN ('income', 'expense')),
		amount INTEGER NOT NULL CHECK(amount > 0),
		category TEXT NOT NULL,
		date TEXT NOT NULL,
		description TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_transactions_date ON transactions(date);
	CREATE INDEX IF NOT EXISTS idx_transactions_category ON transactions(category);
	CREATE INDEX IF NOT EXISTS idx_transactions_type ON transactions(type);
	`
	_, err := r.db.Exec(query)
	return err
}

// Create inserts a new transaction record using parameterized queries.
func (r *SQLiteRepository) Create(t *model.Transaction) error {
	query := `
		INSERT INTO transactions (type, amount, category, date, description)
		VALUES (?, ?, ?, ?, ?)
	`
	res, err := r.db.Exec(query, string(t.Type), int64(t.Amount), t.Category, t.Date, t.Description)
	if err != nil {
		return fmt.Errorf("database insert failed: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("retrieving insert ID failed: %w", err)
	}
	t.ID = id
	return nil
}

// GetByID retrieves a transaction record by its unique identifier.
func (r *SQLiteRepository) GetByID(id int64) (*model.Transaction, error) {
	query := `
		SELECT id, type, amount, category, date, description, created_at
		FROM transactions
		WHERE id = ?
	`
	row := r.db.QueryRow(query, id)

	var t model.Transaction
	var txType string
	var amount int64

	err := row.Scan(&t.ID, &txType, &amount, &t.Category, &t.Date, &t.Description, &t.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("database query failed: %w", err)
	}

	t.Type = model.TransactionType(txType)
	t.Amount = model.Money(amount)
	return &t, nil
}

// List returns a list of transactions matching the provided filter criteria.
func (r *SQLiteRepository) List(filter model.Filter) ([]model.Transaction, error) {
	var conditions []string
	var args []interface{}

	if filter.Type != "" {
		conditions = append(conditions, "type = ?")
		args = append(args, string(filter.Type))
	}
	if filter.Category != "" {
		conditions = append(conditions, "LOWER(category) = LOWER(?)")
		args = append(args, filter.Category)
	}
	if filter.StartDate != "" {
		conditions = append(conditions, "date >= ?")
		args = append(args, filter.StartDate)
	}
	if filter.EndDate != "" {
		conditions = append(conditions, "date <= ?")
		args = append(args, filter.EndDate)
	}

	query := `
		SELECT id, type, amount, category, date, description, created_at
		FROM transactions
	`
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}
	query += " ORDER BY date DESC, id DESC"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("database query failed: %w", err)
	}
	defer rows.Close()

	var transactions []model.Transaction
	for rows.Next() {
		var t model.Transaction
		var txType string
		var amount int64

		if err := rows.Scan(&t.ID, &txType, &amount, &t.Category, &t.Date, &t.Description, &t.CreatedAt); err != nil {
			return nil, fmt.Errorf("reading transaction row failed: %w", err)
		}
		t.Type = model.TransactionType(txType)
		t.Amount = model.Money(amount)
		transactions = append(transactions, t)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("transaction iteration error: %w", err)
	}

	return transactions, nil
}

// Update updates an existing transaction record by ID.
func (r *SQLiteRepository) Update(t *model.Transaction) error {
	query := `
		UPDATE transactions
		SET type = ?, amount = ?, category = ?, date = ?, description = ?
		WHERE id = ?
	`
	res, err := r.db.Exec(query, string(t.Type), int64(t.Amount), t.Category, t.Date, t.Description, t.ID)
	if err != nil {
		return fmt.Errorf("database update failed: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check affected rows: %w", err)
	}
	if rowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

// Delete deletes a transaction record by ID.
func (r *SQLiteRepository) Delete(id int64) error {
	query := `DELETE FROM transactions WHERE id = ?`
	res, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("database delete failed: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check affected rows: %w", err)
	}
	if rowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

// GetTotals computes total income, total expense, and transaction count.
func (r *SQLiteRepository) GetTotals() (model.Money, model.Money, int, error) {
	query := `
		SELECT 
			COALESCE(SUM(CASE WHEN type = 'income' THEN amount ELSE 0 END), 0) AS total_income,
			COALESCE(SUM(CASE WHEN type = 'expense' THEN amount ELSE 0 END), 0) AS total_expense,
			COUNT(*) AS total_count
		FROM transactions
	`
	row := r.db.QueryRow(query)

	var income, expense int64
	var count int
	if err := row.Scan(&income, &expense, &count); err != nil {
		return 0, 0, 0, fmt.Errorf("failed to calculate totals: %w", err)
	}

	return model.Money(income), model.Money(expense), count, nil
}

// GetCategoryBreakdown computes spending or income aggregated by category.
func (r *SQLiteRepository) GetCategoryBreakdown(txType model.TransactionType) ([]model.CategoryBreakdown, error) {
	query := `
		SELECT category, SUM(amount) AS cat_total, COUNT(*) AS cat_count
		FROM transactions
		WHERE type = ?
		GROUP BY category
		ORDER BY cat_total DESC
	`
	rows, err := r.db.Query(query, string(txType))
	if err != nil {
		return nil, fmt.Errorf("failed to get category summary: %w", err)
	}
	defer rows.Close()

	var totalSum int64
	var items []struct {
		category string
		total    int64
		count    int
	}

	for rows.Next() {
		var cat string
		var total int64
		var count int
		if err := rows.Scan(&cat, &total, &count); err != nil {
			return nil, fmt.Errorf("failed to scan category breakdown: %w", err)
		}
		items = append(items, struct {
			category string
			total    int64
			count    int
		}{category: cat, total: total, count: count})
		totalSum += total
	}

	results := make([]model.CategoryBreakdown, 0, len(items))
	for _, item := range items {
		pct := 0.0
		if totalSum > 0 {
			pct = (float64(item.total) / float64(totalSum)) * 100.0
		}
		results = append(results, model.CategoryBreakdown{
			Category:   item.category,
			Total:      model.Money(item.total),
			Percentage: pct,
			Count:      item.count,
		})
	}

	return results, nil
}

// GetMonthlyBreakdown computes income and expense sums grouped by month (YYYY-MM).
func (r *SQLiteRepository) GetMonthlyBreakdown() ([]model.MonthlyBreakdown, error) {
	query := `
		SELECT 
			strftime('%Y-%m', date) AS month,
			COALESCE(SUM(CASE WHEN type = 'income' THEN amount ELSE 0 END), 0) AS inc,
			COALESCE(SUM(CASE WHEN type = 'expense' THEN amount ELSE 0 END), 0) AS exp
		FROM transactions
		GROUP BY strftime('%Y-%m', date)
		ORDER BY month DESC
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get monthly breakdown: %w", err)
	}
	defer rows.Close()

	var list []model.MonthlyBreakdown
	for rows.Next() {
		var month string
		var inc, exp int64
		if err := rows.Scan(&month, &inc, &exp); err != nil {
			return nil, fmt.Errorf("failed to scan monthly row: %w", err)
		}
		list = append(list, model.MonthlyBreakdown{
			Month:   month,
			Income:  model.Money(inc),
			Expense: model.Money(exp),
			Balance: model.Money(inc - exp),
		})
	}

	return list, nil
}

// GetDistinctCategories returns all distinct categories present in the database.
func (r *SQLiteRepository) GetDistinctCategories() ([]string, error) {
	query := `SELECT DISTINCT category FROM transactions ORDER BY category ASC`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch categories: %w", err)
	}
	defer rows.Close()

	var categories []string
	for rows.Next() {
		var cat string
		if err := rows.Scan(&cat); err != nil {
			return nil, fmt.Errorf("failed to scan category: %w", err)
		}
		categories = append(categories, cat)
	}
	return categories, nil
}

// Close closes the underlying SQLite database connection.
func (r *SQLiteRepository) Close() error {
	if r.db != nil {
		return r.db.Close()
	}
	return nil
}
