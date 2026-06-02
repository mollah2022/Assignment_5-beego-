package controllers

import (
	"encoding/json"
	"expense-tracker-api/models"
	"fmt"
	"testing"
)

func TestCreateExpense(t *testing.T) {
	u := newTestUser(t)

	tests := []struct {
		name       string
		body       string
		wantStatus int
		wantMsg    string
	}{
		{
			name:       "valid expense",
			body:       `{"title":"Lunch","amount":350.50,"category":"Food","expense_date":"2025-06-10"}`,
			wantStatus: 201,
			wantMsg:    "Expense created successfully",
		},
		{
			name:       "missing title",
			body:       `{"amount":100,"category":"Food","expense_date":"2025-06-10"}`,
			wantStatus: 400,
			wantMsg:    "Title is required",
		},
		{
			name:       "negative amount",
			body:       `{"title":"Lunch","amount":-10,"category":"Food","expense_date":"2025-06-10"}`,
			wantStatus: 400,
			wantMsg:    "Amount must be positive",
		},
		{
			name:       "missing date",
			body:       `{"title":"Lunch","amount":100,"category":"Food"}`,
			wantStatus: 400,
			wantMsg:    "Expense date is required",
		},
		{
			name:       "invalid category",
			body:       `{"title":"Lunch","amount":100,"category":"INVALID","expense_date":"2025-06-10"}`,
			wantStatus: 400,
			wantMsg:    "Invalid category",
		},
		{
			name:       "no auth header",
			body:       `{"title":"Lunch","amount":100,"category":"Food","expense_date":"2025-06-10"}`,
			wantStatus: 401,
			wantMsg:    "Unauthorized",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uid := u.ID
			if tt.name == "no auth header" {
				uid = 0
			}
			w := do("POST", "/api/v1/expenses", tt.body, uid)
			if w.Code != tt.wantStatus {
				t.Errorf("status: want %d, got %d | body: %s", tt.wantStatus, w.Code, w.Body.String())
			}
			var resp map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &resp)
			if resp["message"] != tt.wantMsg {
				t.Errorf("message: want %q, got %q", tt.wantMsg, resp["message"])
			}
		})
	}
}

func TestListExpenses(t *testing.T) {
	u := newTestUser(t)
	_ = models.CreateExpense(&models.Expense{UserID: u.ID, Title: "F1", Amount: 200, Category: "Food", ExpenseDate: "2025-06-01"})
	_ = models.CreateExpense(&models.Expense{UserID: u.ID, Title: "T1", Amount: 50, Category: "Transport", ExpenseDate: "2025-06-05"})

	tests := []struct {
		name       string
		url        string
		wantStatus int
	}{
		{name: "list all", url: "/api/v1/expenses", wantStatus: 200},
		{name: "filter by category", url: "/api/v1/expenses?category=Food", wantStatus: 200},
		{name: "sort by amount", url: "/api/v1/expenses?sort_by=amount&sort_order=asc", wantStatus: 200},
		{name: "with limit", url: "/api/v1/expenses?limit=1", wantStatus: 200},
		{name: "invalid category", url: "/api/v1/expenses?category=INVALID", wantStatus: 400},
		{name: "invalid sort_by", url: "/api/v1/expenses?sort_by=invalid", wantStatus: 400},
		{name: "invalid sort_order", url: "/api/v1/expenses?sort_order=random", wantStatus: 400},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := do("GET", tt.url, "", u.ID)
			if w.Code != tt.wantStatus {
				t.Errorf("status: want %d, got %d", tt.wantStatus, w.Code)
			}
		})
	}
}

func TestGetOneExpense(t *testing.T) {
	u := newTestUser(t)
	e := &models.Expense{UserID: u.ID, Title: "Movie", Amount: 120, Category: "Entertainment", ExpenseDate: "2025-06-10"}
	_ = models.CreateExpense(e)

	tests := []struct {
		name       string
		id         int
		wantStatus int
	}{
		{name: "found", id: e.ID, wantStatus: 200},
		{name: "not found", id: 99999, wantStatus: 404},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := do("GET", fmt.Sprintf("/api/v1/expenses/%d", tt.id), "", u.ID)
			if w.Code != tt.wantStatus {
				t.Errorf("status: want %d, got %d", tt.wantStatus, w.Code)
			}
		})
	}
}

func TestUpdateExpense(t *testing.T) {
	u := newTestUser(t)
	e := &models.Expense{UserID: u.ID, Title: "Old", Amount: 100, Category: "Food", ExpenseDate: "2025-06-10"}
	_ = models.CreateExpense(e)

	tests := []struct {
		name       string
		id         int
		body       string
		wantStatus int
		wantMsg    string
	}{
		{
			name: "valid update", id: e.ID,
			body:       `{"title":"New","amount":500}`,
			wantStatus: 200, wantMsg: "Expense updated successfully",
		},
		{
			name: "invalid category", id: e.ID,
			body:       `{"category":"INVALID"}`,
			wantStatus: 400, wantMsg: "Invalid category",
		},
		{
			name: "not found", id: 99999,
			body:       `{"title":"X"}`,
			wantStatus: 404, wantMsg: "Expense not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := do("PUT", fmt.Sprintf("/api/v1/expenses/%d", tt.id), tt.body, u.ID)
			if w.Code != tt.wantStatus {
				t.Errorf("status: want %d, got %d", tt.wantStatus, w.Code)
			}
			var resp map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &resp)
			if resp["message"] != tt.wantMsg {
				t.Errorf("message: want %q, got %q", tt.wantMsg, resp["message"])
			}
		})
	}
}

func TestDeleteExpense(t *testing.T) {
	u := newTestUser(t)
	e := &models.Expense{UserID: u.ID, Title: "Del", Amount: 50, Category: "Other", ExpenseDate: "2025-06-10"}
	_ = models.CreateExpense(e)

	tests := []struct {
		name       string
		id         int
		wantStatus int
		wantMsg    string
	}{
		{name: "valid delete", id: e.ID, wantStatus: 200, wantMsg: "Expense deleted successfully"},
		{name: "already deleted", id: e.ID, wantStatus: 404, wantMsg: "Expense not found"},
		{name: "not found", id: 99999, wantStatus: 404, wantMsg: "Expense not found"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := do("DELETE", fmt.Sprintf("/api/v1/expenses/%d", tt.id), "", u.ID)
			if w.Code != tt.wantStatus {
				t.Errorf("status: want %d, got %d", tt.wantStatus, w.Code)
			}
			var resp map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &resp)
			if resp["message"] != tt.wantMsg {
				t.Errorf("message: want %q, got %q", tt.wantMsg, resp["message"])
			}
		})
	}
}

func TestSummary(t *testing.T) {
	u := newTestUser(t)
	_ = models.CreateExpense(&models.Expense{UserID: u.ID, Title: "G1", Amount: 300, Category: "Food", ExpenseDate: "2025-06-05"})
	_ = models.CreateExpense(&models.Expense{UserID: u.ID, Title: "G2", Amount: 100, Category: "Transport", ExpenseDate: "2025-06-08"})

	tests := []struct {
		name       string
		url        string
		wantStatus int
		wantMsg    string
	}{
		{
			name:       "valid summary",
			url:        "/api/v1/expenses/summary?date_from=2025-06-01&date_to=2025-06-30",
			wantStatus: 200, wantMsg: "Summary generated",
		},
		{
			name:       "missing date_from",
			url:        "/api/v1/expenses/summary?date_to=2025-06-30",
			wantStatus: 400, wantMsg: "date_from is required",
		},
		{
			name:       "missing date_to",
			url:        "/api/v1/expenses/summary?date_from=2025-06-01",
			wantStatus: 400, wantMsg: "date_to is required",
		},
		{
			name:       "date_from after date_to",
			url:        "/api/v1/expenses/summary?date_from=2025-07-01&date_to=2025-06-01",
			wantStatus: 400, wantMsg: "date_from cannot be after date_to",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := do("GET", tt.url, "", u.ID)
			if w.Code != tt.wantStatus {
				t.Errorf("status: want %d, got %d", tt.wantStatus, w.Code)
			}
			var resp map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &resp)
			if resp["message"] != tt.wantMsg {
				t.Errorf("message: want %q, got %q", tt.wantMsg, resp["message"])
			}
		})
	}
}
