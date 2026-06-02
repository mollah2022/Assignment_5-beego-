package models

import (
	"os"
	"path/filepath"
	"testing"
)

// TestMain sets up isolated temp files for all model tests.
func TestMain(m *testing.M) {
	tmp, err := os.MkdirTemp("", "expense-models-*")
	if err != nil {
		panic("failed to create temp dir: " + err.Error())
	}

	UserFile = filepath.Join(tmp, "users.csv")
	ExpenseFile = filepath.Join(tmp, "expenses.csv")

	code := m.Run()

	os.RemoveAll(tmp)
	os.Exit(code)
}

// makeUser is a shared helper that creates a test user.
func makeUser(t *testing.T, email string) *User {
	t.Helper()
	u := &User{Name: "Test", Email: email, Password: "pass123"}
	if err := CreateUser(u); err != nil {
		t.Fatalf("makeUser: %v", err)
	}
	return u
}