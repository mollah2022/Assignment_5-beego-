package models

import (
	"encoding/csv"
	"errors"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"time"
)

// ExpenseFile is the CSV path for expenses. Override in tests.
var ExpenseFile = "data/expenses.csv"

// AllowedCategories lists every valid expense category.
var AllowedCategories = []string{
	"Food", "Transport", "Housing", "Entertainment",
	"Shopping", "Healthcare", "Education", "Utilities", "Other",
}

// Expense represents a single spending record.
type Expense struct {
	ID          int     `json:"id"`
	UserID      int     `json:"user_id"`
	Title       string  `json:"title"`
	Amount      float64 `json:"amount"`
	Category    string  `json:"category"`
	Note        string  `json:"note"`
	ExpenseDate string  `json:"expense_date"`
	CreatedAt   string  `json:"created_at"`
}

// FilterOptions holds query parameters for filtering and sorting.
type FilterOptions struct {
	Category  string
	DateFrom  string
	DateTo    string
	SortBy    string
	SortOrder string
	Limit     int
}

// Summary holds the result of a date-range spending query.
type Summary struct {
	DateFrom    string  `json:"date_from"`
	DateTo      string  `json:"date_to"`
	TotalAmount float64 `json:"total_amount"`
	TotalCount  int     `json:"total_count"`
}

// IsValidCategory reports whether cat is an allowed category.
func IsValidCategory(cat string) bool {
	for _, c := range AllowedCategories {
		if c == cat {
			return true
		}
	}
	return false
}

// ensureExpenseFile creates the file and parent directory if missing.
func ensureExpenseFile() error {
	if err := os.MkdirAll(filepath.Dir(ExpenseFile), 0755); err != nil {
		return err
	}
	if _, err := os.Stat(ExpenseFile); os.IsNotExist(err) {
		f, err := os.Create(ExpenseFile)
		if err != nil {
			return err
		}
		defer f.Close()
		w := csv.NewWriter(f)
		if err := w.Write([]string{
			"id", "user_id", "title", "amount",
			"category", "note", "expense_date", "created_at",
		}); err != nil {
			return err
		}
		w.Flush()
		return w.Error()
	}
	return nil
}

// getAllExpenses reads every expense row from the CSV.
func getAllExpenses() ([]Expense, error) {
	if err := ensureExpenseFile(); err != nil {
		return nil, err
	}
	f, err := os.Open(ExpenseFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	rows, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return nil, err
	}

	expenses := make([]Expense, 0, len(rows))
	for _, row := range rows[1:] { // skip header
		id, _ := strconv.Atoi(row[0])
		userID, _ := strconv.Atoi(row[1])
		amount, _ := strconv.ParseFloat(row[3], 64)
		expenses = append(expenses, Expense{
			ID:          id,
			UserID:      userID,
			Title:       row[2],
			Amount:      amount,
			Category:    row[4],
			Note:        row[5],
			ExpenseDate: row[6],
			CreatedAt:   row[7],
		})
	}
	return expenses, nil
}

// expenseToRow converts an Expense to a CSV string slice.
func expenseToRow(e Expense) []string {
	return []string{
		strconv.Itoa(e.ID),
		strconv.Itoa(e.UserID),
		e.Title,
		strconv.FormatFloat(e.Amount, 'f', 2, 64),
		e.Category,
		e.Note,
		e.ExpenseDate,
		e.CreatedAt,
	}
}

// writeAllExpenses overwrites the CSV with the provided slice.
func writeAllExpenses(expenses []Expense) error {
	f, err := os.Create(ExpenseFile)
	if err != nil {
		return err
	}
	defer f.Close()

	w := csv.NewWriter(f)
	if err := w.Write([]string{
		"id", "user_id", "title", "amount",
		"category", "note", "expense_date", "created_at",
	}); err != nil {
		return err
	}
	for _, e := range expenses {
		if err := w.Write(expenseToRow(e)); err != nil {
			return err
		}
	}
	w.Flush()
	return w.Error()
}

// nextExpenseID returns max existing ID + 1.
func nextExpenseID() (int, error) {
	expenses, err := getAllExpenses()
	if err != nil {
		return 0, err
	}
	max := 0
	for _, e := range expenses {
		if e.ID > max {
			max = e.ID
		}
	}
	return max + 1, nil
}

// GetExpensesByUserID returns all expenses owned by userID.
func GetExpensesByUserID(userID int) ([]Expense, error) {
	all, err := getAllExpenses()
	if err != nil {
		return nil, err
	}
	result := make([]Expense, 0)
	for _, e := range all {
		if e.UserID == userID {
			result = append(result, e)
		}
	}
	return result, nil
}

// GetExpenseByID returns the expense matching id and userID.
func GetExpenseByID(id, userID int) (*Expense, error) {
	all, err := getAllExpenses()
	if err != nil {
		return nil, err
	}
	for i := range all {
		if all[i].ID == id && all[i].UserID == userID {
			return &all[i], nil
		}
	}
	return nil, errors.New("expense not found")
}

// CreateExpense appends a new row and sets ID and CreatedAt.
func CreateExpense(expense *Expense) error {
	if err := ensureExpenseFile(); err != nil {
		return err
	}
	id, err := nextExpenseID()
	if err != nil {
		return err
	}
	expense.ID = id
	expense.CreatedAt = time.Now().UTC().Format(time.RFC3339)

	f, err := os.OpenFile(ExpenseFile, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	w := csv.NewWriter(f)
	if err := w.Write(expenseToRow(*expense)); err != nil {
		return err
	}
	w.Flush()
	return w.Error()
}

// UpdateExpense replaces the matching row with updated data.
func UpdateExpense(updated *Expense) error {
	all, err := getAllExpenses()
	if err != nil {
		return err
	}
	found := false
	for i := range all {
		if all[i].ID == updated.ID && all[i].UserID == updated.UserID {
			all[i] = *updated
			found = true
			break
		}
	}
	if !found {
		return errors.New("expense not found")
	}
	return writeAllExpenses(all)
}

// DeleteExpense removes the matching row from the CSV.
func DeleteExpense(id, userID int) error {
	all, err := getAllExpenses()
	if err != nil {
		return err
	}
	updated := make([]Expense, 0, len(all))
	found := false
	for _, e := range all {
		if e.ID == id && e.UserID == userID {
			found = true
			continue
		}
		updated = append(updated, e)
	}
	if !found {
		return errors.New("expense not found")
	}
	return writeAllExpenses(updated)
}

// FilterExpenses returns expenses for a user, filtered and sorted by opts.
func FilterExpenses(userID int, opts FilterOptions) ([]Expense, error) {
	result, err := GetExpensesByUserID(userID)
	if err != nil {
		return nil, err
	}

	if opts.Category != "" {
		var filtered []Expense
		for _, e := range result {
			if e.Category == opts.Category {
				filtered = append(filtered, e)
			}
		}
		result = filtered
	}

	if opts.DateFrom != "" {
		var filtered []Expense
		for _, e := range result {
			if e.ExpenseDate >= opts.DateFrom {
				filtered = append(filtered, e)
			}
		}
		result = filtered
	}

	if opts.DateTo != "" {
		var filtered []Expense
		for _, e := range result {
			if e.ExpenseDate <= opts.DateTo {
				filtered = append(filtered, e)
			}
		}
		result = filtered
	}

	if opts.SortBy != "" {
		sort.Slice(result, func(i, j int) bool {
			switch opts.SortBy {
			case "amount":
				if opts.SortOrder == "asc" {
					return result[i].Amount < result[j].Amount
				}
				return result[i].Amount > result[j].Amount
			case "expense_date":
				if opts.SortOrder == "asc" {
					return result[i].ExpenseDate < result[j].ExpenseDate
				}
				return result[i].ExpenseDate > result[j].ExpenseDate
			}
			return false
		})
	}

	if opts.Limit > 0 && opts.Limit < len(result) {
		result = result[:opts.Limit]
	}

	return result, nil
}

// GetSummary calculates total spend within the given date range.
func GetSummary(userID int, dateFrom, dateTo string) (*Summary, error) {
	expenses, err := GetExpensesByUserID(userID)
	if err != nil {
		return nil, err
	}

	summary := &Summary{DateFrom: dateFrom, DateTo: dateTo}
	for _, e := range expenses {
		if e.ExpenseDate >= dateFrom && e.ExpenseDate <= dateTo {
			summary.TotalAmount += e.Amount
			summary.TotalCount++
		}
	}
	return summary, nil
}