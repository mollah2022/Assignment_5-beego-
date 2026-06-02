package models

import (
	"fmt"
	"testing"
	"time"
)

// uniqueEmail generates a unique email per test to avoid collisions.
func uniqueEmail() string {
	return fmt.Sprintf("user%d@example.com", time.Now().UnixNano())
}

func TestCreateExpense(t *testing.T) {
	u := makeUser(t, uniqueEmail())

	tests := []struct {
		name    string
		expense Expense
		wantErr bool
	}{
		{
			name:    "valid expense",
			expense: Expense{UserID: u.ID, Title: "Lunch", Amount: 350.50, Category: "Food", ExpenseDate: "2025-06-10"},
			wantErr: false,
		},
		{
			name:    "another valid expense",
			expense: Expense{UserID: u.ID, Title: "Bus", Amount: 50.00, Category: "Transport", ExpenseDate: "2025-06-11"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CreateExpense(&tt.expense)
			if (err != nil) != tt.wantErr {
				t.Fatalf("CreateExpense() error = %v, wantErr = %v", err, tt.wantErr)
			}
			if !tt.wantErr && tt.expense.ID == 0 {
				t.Error("expected ID to be set after CreateExpense")
			}
		})
	}
}

func TestGetExpensesByUserID(t *testing.T) {
	u := makeUser(t, uniqueEmail())
	_ = CreateExpense(&Expense{UserID: u.ID, Title: "Dinner", Amount: 200, Category: "Food", ExpenseDate: "2025-06-12"})

	expenses, err := GetExpensesByUserID(u.ID)
	if err != nil {
		t.Fatalf("GetExpensesByUserID() error: %v", err)
	}
	if len(expenses) == 0 {
		t.Fatal("expected at least one expense")
	}
	for _, e := range expenses {
		if e.UserID != u.ID {
			t.Errorf("expected userID %d, got %d", u.ID, e.UserID)
		}
	}
}

func TestGetExpenseByID(t *testing.T) {
	u := makeUser(t, uniqueEmail())
	e := &Expense{UserID: u.ID, Title: "Movie", Amount: 120, Category: "Entertainment", ExpenseDate: "2025-06-13"}
	_ = CreateExpense(e)

	tests := []struct {
		name      string
		id        int
		wantFound bool
	}{
		{name: "existing expense", id: e.ID, wantFound: true},
		{name: "unknown expense", id: 99999, wantFound: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetExpenseByID(tt.id, u.ID)
			if tt.wantFound {
				if err != nil {
					t.Fatalf("expected expense, got error: %v", err)
				}
				if got.ID != tt.id {
					t.Errorf("expected ID %d, got %d", tt.id, got.ID)
				}
			} else {
				if err == nil {
					t.Error("expected error for unknown expense, got nil")
				}
			}
		})
	}
}

func TestUpdateExpense(t *testing.T) {
	u := makeUser(t, uniqueEmail())
	e := &Expense{UserID: u.ID, Title: "Old Title", Amount: 100, Category: "Food", ExpenseDate: "2025-06-14"}
	_ = CreateExpense(e)

	e.Title = "New Title"
	e.Amount = 999

	if err := UpdateExpense(e); err != nil {
		t.Fatalf("UpdateExpense() error: %v", err)
	}
	got, _ := GetExpenseByID(e.ID, u.ID)
	if got.Title != "New Title" {
		t.Errorf("expected title 'New Title', got '%s'", got.Title)
	}
	if got.Amount != 999 {
		t.Errorf("expected amount 999, got %f", got.Amount)
	}
}

func TestDeleteExpense(t *testing.T) {
	u := makeUser(t, uniqueEmail())
	e := &Expense{UserID: u.ID, Title: "Delete Me", Amount: 50, Category: "Other", ExpenseDate: "2025-06-15"}
	_ = CreateExpense(e)

	if err := DeleteExpense(e.ID, u.ID); err != nil {
		t.Fatalf("DeleteExpense() first call error: %v", err)
	}
	if _, err := GetExpenseByID(e.ID, u.ID); err == nil {
		t.Error("expected error after deletion, got nil")
	}
	if err := DeleteExpense(e.ID, u.ID); err == nil {
		t.Error("expected error on second delete, got nil")
	}
}

func TestFilterExpenses(t *testing.T) {
	u := makeUser(t, uniqueEmail())
	_ = CreateExpense(&Expense{UserID: u.ID, Title: "F1", Amount: 200, Category: "Food", ExpenseDate: "2025-06-01"})
	_ = CreateExpense(&Expense{UserID: u.ID, Title: "T1", Amount: 50, Category: "Transport", ExpenseDate: "2025-06-05"})
	_ = CreateExpense(&Expense{UserID: u.ID, Title: "F2", Amount: 400, Category: "Food", ExpenseDate: "2025-06-10"})

	tests := []struct {
		name      string
		opts      FilterOptions
		wantCount int
	}{
		{name: "filter by Food", opts: FilterOptions{Category: "Food"}, wantCount: 2},
		{name: "filter by date range", opts: FilterOptions{DateFrom: "2025-06-04", DateTo: "2025-06-06"}, wantCount: 1},
		{name: "limit to 1", opts: FilterOptions{Limit: 1}, wantCount: 1},
		{name: "sort amount asc", opts: FilterOptions{SortBy: "amount", SortOrder: "asc"}, wantCount: 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FilterExpenses(u.ID, tt.opts)
			if err != nil {
				t.Fatalf("FilterExpenses() error: %v", err)
			}
			if len(got) != tt.wantCount {
				t.Errorf("expected %d, got %d", tt.wantCount, len(got))
			}
		})
	}
}

func TestGetSummary(t *testing.T) {
	u := makeUser(t, uniqueEmail())
	_ = CreateExpense(&Expense{UserID: u.ID, Title: "G1", Amount: 300, Category: "Food", ExpenseDate: "2025-06-05"})
	_ = CreateExpense(&Expense{UserID: u.ID, Title: "G2", Amount: 100, Category: "Transport", ExpenseDate: "2025-06-08"})

	s, err := GetSummary(u.ID, "2025-06-01", "2025-06-30")
	if err != nil {
		t.Fatalf("GetSummary() error: %v", err)
	}
	if s.TotalCount != 2 {
		t.Errorf("expected count 2, got %d", s.TotalCount)
	}
	if s.TotalAmount != 400 {
		t.Errorf("expected amount 400, got %f", s.TotalAmount)
	}
}