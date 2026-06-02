package controllers

import (
	"bytes"
	"expense-tracker-api/models"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	beego "github.com/beego/beego/v2/server/web"
)

// TestMain sets up temp files and registers routes for all controller tests.
func TestMain(m *testing.M) {
	tmp, err := os.MkdirTemp("", "expense-ctrl-*")
	if err != nil {
		panic("failed to create temp dir: " + err.Error())
	}

	models.UserFile = filepath.Join(tmp, "users.csv")
	models.ExpenseFile = filepath.Join(tmp, "expenses.csv")

	beego.Router("/api/v1/auth/register", &AuthController{}, "post:Register")
	beego.Router("/api/v1/auth/login", &AuthController{}, "post:Login")
	beego.Router("/api/v1/expenses/summary", &ExpenseController{}, "get:Summary")
	beego.Router("/api/v1/expenses", &ExpenseController{}, "post:Create;get:List")
	beego.Router("/api/v1/expenses/:id", &ExpenseController{}, "get:GetOne;put:Update;delete:Delete")

	// beego.TestBeegoInit("../")
	// এই দুই line দিয়ে replace করো
	appPath, _ := filepath.Abs("../")
	beego.TestBeegoInit(appPath)

	code := m.Run()
	os.RemoveAll(tmp)
	os.Exit(code)
}

// do sends an HTTP request and returns the response recorder.
func do(method, url, body string, userID int) *httptest.ResponseRecorder {
	var req *http.Request
	if body != "" {
		req, _ = http.NewRequest(method, url, bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, _ = http.NewRequest(method, url, nil)
	}
	if userID > 0 {
		req.Header.Set("X-User-ID", strconv.Itoa(userID))
	}
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, req)
	return w
}

// newTestUser creates a unique user for controller tests.
func newTestUser(t *testing.T) *models.User {
	t.Helper()
	u := &models.User{
		Name:     "Test User",
		Email:    fmt.Sprintf("u%d@example.com", time.Now().UnixNano()),
		Password: "pass123",
	}
	if err := models.CreateUser(u); err != nil {
		t.Fatalf("newTestUser: %v", err)
	}
	return u
}
