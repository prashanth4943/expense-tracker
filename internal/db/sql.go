package db

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"github.com/prashanth4943/expense-tracker/internal/models"
	"github.com/prashanth4943/expense-tracker/internal/utils"
)

type Store struct {
	db *sql.DB
}

// New opens (or creates) the SQLite database and runs migrations.
func New(dsn string) (*Store, error) {
	db, err := sql.Open("sqlite3", dsn+"?_journal_mode=WAL&_foreign_keys=on")
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}
	// SQLite performs best with a single writer connection.
	db.SetMaxOpenConns(1)

	s := &Store{db: db}
	if err := s.migrate(); err != nil {
		return nil, fmt.Errorf("migrate: %w", err)
	}
	return s, nil
}

func (s *Store) migrate() error {
	_, err := s.db.Exec(`
		CREATE TABLE IF NOT EXISTS expenses (
			id               TEXT PRIMARY KEY,
			amount_paise     INTEGER NOT NULL CHECK(amount_paise > 0),
			category         TEXT    NOT NULL,
			description      TEXT    NOT NULL DEFAULT '',
			date             TEXT    NOT NULL,
			created_at       TEXT    NOT NULL,
			idempotency_key  TEXT    NOT NULL UNIQUE
		);
		CREATE INDEX IF NOT EXISTS idx_expenses_category ON expenses(category);
		CREATE INDEX IF NOT EXISTS idx_expenses_date     ON expenses(date DESC);
		CREATE UNIQUE INDEX IF NOT EXISTS idx_idempotency_key ON expenses(idempotency_key); 
		`)
	return err
}

// CreateExpense inserts a new expense.
// If idempotency_key already exists it returns the existing row — safe to retry.
func (s *Store) CreateExpense(req models.CreateExpenseRequest, amountPaise int64) (*models.Expense, error) {
	// Check idempotency first (existing key → return stored record, do not insert again).
	existing, err := s.findByIdempotencyKey(req.IdempotencyKey)
	if err == nil && existing != nil {
		return existing, nil
	}

	id := uuid.New().String()
	now := time.Now().UTC().Format(time.RFC3339)

	_, err = s.db.Exec(`
		INSERT INTO expenses (id, amount_paise, category, description, date, created_at, idempotency_key)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		id, amountPaise, req.Category, req.Description, req.Date, now, req.IdempotencyKey,
	)
	if err != nil {
		// Another concurrent insert with the same key won the race — fetch and return it.
		if isUniqueConstraintErr(err) {
			return s.findByIdempotencyKey(req.IdempotencyKey)
		}
		return nil, fmt.Errorf("insert expense: %w", err)
	}

	return s.findByID(id)
}

type ListFilter struct {
	Category string
	Sort     string // "date_desc" (default) or ""
}

// ListExpenses returns expenses, optionally filtered and sorted.
func (s *Store) ListExpenses(f ListFilter) ([]models.Expense, error) {
	query := `SELECT id, amount_paise, category, description, date, created_at FROM expenses WHERE 1=1`
	args := []any{}

	if f.Category != "" {
		query += ` AND category = ?`
		args = append(args, f.Category)
	}

	// 🔥 Safe sort mapping
	orderBy := "date DESC, created_at DESC" // default

	switch f.Sort {
	case "date_desc":
		orderBy = "date DESC, created_at DESC"
	case "date_asc":
		orderBy = "date ASC, created_at ASC"
	case "amount_desc":
		orderBy = "amount_paise DESC"
	case "amount_asc":
		orderBy = "amount_paise ASC"
	}

	query += " ORDER BY " + orderBy

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("list expenses: %w", err)
	}
	defer rows.Close()

	var expenses []models.Expense
	for rows.Next() {
		var e models.Expense
		if err := rows.Scan(&e.ID, &e.AmountPaise, &e.Category, &e.Description, &e.Date, &e.CreatedAt); err != nil {
			return nil, err
		}
		e.AmountDisplay = utils.FormatAmount(e.AmountPaise)
		expenses = append(expenses, e)
	}
	return expenses, rows.Err()
}

// Categories returns all distinct categories (for filter dropdown).
func (s *Store) Categories() ([]string, error) {
	rows, err := s.db.Query(`SELECT DISTINCT category FROM expenses ORDER BY category`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var cats []string
	for rows.Next() {
		var c string
		if err := rows.Scan(&c); err != nil {
			return nil, err
		}
		cats = append(cats, c)
	}
	return cats, rows.Err()
}

// --- helpers ---

func (s *Store) findByIdempotencyKey(key string) (*models.Expense, error) {
	row := s.db.QueryRow(
		`SELECT id, amount_paise, category, description, date, created_at FROM expenses WHERE idempotency_key = ?`, key)
	return scanExpense(row)
}

func (s *Store) findByID(id string) (*models.Expense, error) {
	row := s.db.QueryRow(
		`SELECT id, amount_paise, category, description, date, created_at FROM expenses WHERE id = ?`, id)
	return scanExpense(row)
}

func scanExpense(row *sql.Row) (*models.Expense, error) {
	var e models.Expense
	if err := row.Scan(&e.ID, &e.AmountPaise, &e.Category, &e.Description, &e.Date, &e.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	e.AmountDisplay = utils.FormatAmount(e.AmountPaise)
	return &e, nil
}

func isUniqueConstraintErr(err error) bool {
	return err != nil && (contains(err.Error(), "UNIQUE constraint failed") ||
		contains(err.Error(), "unique constraint"))
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && containsStr(s, sub))
}

func containsStr(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
