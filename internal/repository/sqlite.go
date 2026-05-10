package repository

import (
	"database/sql"
	"expense_tracker_api/internal/models"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type ExpenseRepository struct {
	db *sql.DB
}

func NewExpenseRepository(dbPath string) (*ExpenseRepository, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	repo := &ExpenseRepository{db: db}
	if err := repo.initTable(); err != nil {
		return nil, err
	}

	return repo, nil
}

func (r *ExpenseRepository) initTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS expenses (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		amount REAL NOT NULL,
		category TEXT NOT NULL,
		note TEXT,
		spent_on TEXT NOT NULL,
		created_at DATETIME NOT NULL
	);`
	_, err := r.db.Exec(query)
	return err
}

func (r *ExpenseRepository) Create(e *models.Expense) error {
	e.CreatedAt = time.Now().UTC()
	query := `INSERT INTO expenses (amount, category, note, spent_on, created_at) VALUES (?, ?, ?, ?, ?)`
	res, err := r.db.Exec(query, e.Amount, e.Category, e.Note, e.SpentOn, e.CreatedAt)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err == nil {
		e.ID = int(id)
	}
	return err
}

func (r *ExpenseRepository) GetAll() ([]models.Expense, error) {
	query := `SELECT id, amount, category, note, spent_on, created_at FROM expenses ORDER BY created_at DESC`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var expenses []models.Expense
	for rows.Next() {
		var e models.Expense
		err := rows.Scan(&e.ID, &e.Amount, &e.Category, &e.Note, &e.SpentOn, &e.CreatedAt)
		if err != nil {
			return nil, err
		}
		expenses = append(expenses, e)
	}
	return expenses, nil
}

func (r *ExpenseRepository) GetByID(id int) (*models.Expense, error) {
	query := `SELECT id, amount, category, note, spent_on, created_at FROM expenses WHERE id = ?`
	row := r.db.QueryRow(query, id)

	var e models.Expense
	err := row.Scan(&e.ID, &e.Amount, &e.Category, &e.Note, &e.SpentOn, &e.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &e, err
}

func (r *ExpenseRepository) Update(id int, e *models.Expense) error {
	query := `UPDATE expenses SET amount = ?, category = ?, note = ?, spent_on = ? WHERE id = ?`
	_, err := r.db.Exec(query, e.Amount, e.Category, e.Note, e.SpentOn, id)
	return err
}

func (r *ExpenseRepository) Delete(id int) (bool, error) {
	query := `DELETE FROM expenses WHERE id = ?`
	res, err := r.db.Exec(query, id)
	if err != nil {
		return false, err
	}
	rows, _ := res.RowsAffected()
	return rows > 0, nil
}

func (r *ExpenseRepository) GetSummary() (*models.Summary, error) {
	summary := &models.Summary{
		ByCategory: make(map[string]float64),
	}

	// Total amount
	row := r.db.QueryRow(`SELECT COALESCE(SUM(amount), 0) FROM expenses`)
	if err := row.Scan(&summary.TotalAmount); err != nil {
		return nil, err
	}

	// Group by category
	rows, err := r.db.Query(`SELECT category, SUM(amount) FROM expenses GROUP BY category`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var cat string
		var sum float64
		if err := rows.Scan(&cat, &sum); err != nil {
			return nil, err
		}
		summary.ByCategory[cat] = sum
	}

	return summary, nil
}
