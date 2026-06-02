package controllers

import (
	"encoding/json"
	"expense-tracker-api/models"
	"testing"
)

func TestRegister(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		wantStatus int
		wantMsg    string
	}{
		{
			name:       "valid registration",
			body:       `{"name":"Alice","email":"alice@ctrl.com","password":"pass123"}`,
			wantStatus: 201,
			wantMsg:    "User registered successfully",
		},
		{
			name:       "missing name",
			body:       `{"email":"noname@ctrl.com","password":"pass123"}`,
			wantStatus: 400,
			wantMsg:    "Name is required",
		},
		{
			name:       "missing email",
			body:       `{"name":"Bob","password":"pass123"}`,
			wantStatus: 400,
			wantMsg:    "Email is required",
		},
		{
			name:       "invalid email format",
			body:       `{"name":"Bob","email":"not-valid","password":"pass123"}`,
			wantStatus: 400,
			wantMsg:    "Invalid email format",
		},
		{
			name:       "short password",
			body:       `{"name":"Bob","email":"bob@ctrl.com","password":"123"}`,
			wantStatus: 400,
			wantMsg:    "Password must be at least 6 characters",
		},
		{
			name:       "duplicate email",
			body:       `{"name":"Alice2","email":"alice@ctrl.com","password":"pass123"}`,
			wantStatus: 409,
			wantMsg:    "Email already exists",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := do("POST", "/api/v1/auth/register", tt.body, 0)
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

func TestLogin(t *testing.T) {
	_ = models.CreateUser(&models.User{Name: "Login User", Email: "login@ctrl.com", Password: "pass123"})

	tests := []struct {
		name       string
		body       string
		wantStatus int
		wantMsg    string
	}{
		{
			name:       "valid login",
			body:       `{"email":"login@ctrl.com","password":"pass123"}`,
			wantStatus: 200,
			wantMsg:    "Login successful",
		},
		{
			name:       "wrong password",
			body:       `{"email":"login@ctrl.com","password":"wrong"}`,
			wantStatus: 401,
			wantMsg:    "Invalid email or password",
		},
		{
			name:       "unknown email",
			body:       `{"email":"ghost@ctrl.com","password":"pass123"}`,
			wantStatus: 401,
			wantMsg:    "Invalid email or password",
		},
		{
			name:       "missing credentials",
			body:       `{}`,
			wantStatus: 400,
			wantMsg:    "Email and password are required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := do("POST", "/api/v1/auth/login", tt.body, 0)
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
