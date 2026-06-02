package controllers

import (
	"encoding/json"
	"expense-tracker-api/models"
	"strconv"

	"github.com/beego/beego/v2/core/logs"
)

// ExpenseController handles all expense operations.
type ExpenseController struct {
	BaseController
}

type ExpenseCreateRequest struct {
	Title       string  `json:"title"`
	Amount      float64 `json:"amount"`
	Category    string  `json:"category"`
	Note        string  `json:"note"`
	ExpenseDate string  `json:"expense_date"`
}

type ExpenseUpdateRequest struct {
	Title       string  `json:"title"`
	Amount      float64 `json:"amount"`
	Category    string  `json:"category"`
	Note        string  `json:"note"`
	ExpenseDate string  `json:"expense_date"`
}

func (c *ExpenseController) authenticatedUserID() (int, bool) {
	header := c.Ctx.Input.Header("X-User-ID")
	if header == "" {
		c.SendError(401, "Unauthorized")
		return 0, false
	}
	id, err := strconv.Atoi(header)
	if err != nil {
		c.SendError(401, "Unauthorized")
		return 0, false
	}
	if _, err := models.GetUserByID(id); err != nil {
		c.SendError(401, "Unauthorized")
		return 0, false
	}
	return id, true
}

// Create adds a new expense.
// @Title Create Expense
// @Description Add a new expense for the authenticated user
// @Param	X-User-ID	header	int		true	"User ID from login"
// @Param	body	body	controllers.ExpenseCreateRequest	true	"Expense payload"
// @Success 201 {string} success
// @Failure 400 {string} error
// @Failure 401 {string} error
// @router / [post]
func (c *ExpenseController) Create() {
	userID, ok := c.authenticatedUserID()
	if !ok {
		return
	}
	var input struct {
		Title       string  `json:"title"`
		Amount      float64 `json:"amount"`
		Category    string  `json:"category"`
		Note        string  `json:"note"`
		ExpenseDate string  `json:"expense_date"`
	}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &input); err != nil {
		c.SendError(400, "Invalid request body")
		return
	}
	switch {
	case input.Title == "":
		c.SendError(400, "Title is required")
		return
	case input.Amount <= 0:
		c.SendError(400, "Amount must be positive")
		return
	case input.ExpenseDate == "":
		c.SendError(400, "Expense date is required")
		return
	case !models.IsValidCategory(input.Category):
		c.SendError(400, "Invalid category")
		return
	}
	expense := &models.Expense{
		UserID:      userID,
		Title:       input.Title,
		Amount:      input.Amount,
		Category:    input.Category,
		Note:        input.Note,
		ExpenseDate: input.ExpenseDate,
	}
	if err := models.CreateExpense(expense); err != nil {
		logs.Error("Create expense:", err)
		c.SendError(500, "Failed to create expense")
		return
	}
	c.SendSuccess(201, "Expense created successfully", expense)
}

// List returns filtered expenses.
// @Title List Expenses
// @Description List all expenses with optional filters and sorting
// @Param	X-User-ID	header	int		true	"User ID from login"
// @Param	category	query	string	false	"Filter by category"
// @Param	date_from	query	string	false	"From date (YYYY-MM-DD)"
// @Param	date_to		query	string	false	"To date (YYYY-MM-DD)"
// @Param	sort_by		query	string	false	"Sort by: amount or expense_date"
// @Param	sort_order	query	string	false	"Sort order: asc or desc"
// @Param	limit		query	int		false	"Max results (default 10)"
// @Success 200 {string} success
// @Failure 401 {string} error
// @router / [get]
func (c *ExpenseController) List() {
	userID, ok := c.authenticatedUserID()
	if !ok {
		return
	}
	category := c.GetString("category")
	dateFrom := c.GetString("date_from")
	dateTo := c.GetString("date_to")
	sortBy := c.GetString("sort_by")
	sortOrder := c.GetString("sort_order", "desc")
	limitStr := c.GetString("limit", "10")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		c.SendError(400, "limit must be a positive integer")
		return
	}
	if category != "" && !models.IsValidCategory(category) {
		c.SendError(400, "Invalid category")
		return
	}
	if sortBy != "" && sortBy != "amount" && sortBy != "expense_date" {
		c.SendError(400, "sort_by must be 'amount' or 'expense_date'")
		return
	}
	if sortOrder != "asc" && sortOrder != "desc" {
		c.SendError(400, "sort_order must be 'asc' or 'desc'")
		return
	}
	expenses, err := models.FilterExpenses(userID, models.FilterOptions{
		Category:  category,
		DateFrom:  dateFrom,
		DateTo:    dateTo,
		SortBy:    sortBy,
		SortOrder: sortOrder,
		Limit:     limit,
	})
	if err != nil {
		logs.Error("List expenses:", err)
		c.SendError(500, "Failed to retrieve expenses")
		return
	}
	c.SendSuccess(200, "Expenses retrieved", expenses)
}

// GetOne returns a single expense by ID.
// @Title Get One Expense
// @Description Get a single expense by its ID
// @Param	X-User-ID	header	int	true	"User ID from login"
// @Param	id			path	int	true	"Expense ID"
// @Success 200 {string} success
// @Failure 404 {string} error
// @Failure 401 {string} error
// @router /:id [get]
func (c *ExpenseController) GetOne() {
	userID, ok := c.authenticatedUserID()
	if !ok {
		return
	}
	id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err != nil {
		c.SendError(400, "Invalid expense ID")
		return
	}
	expense, err := models.GetExpenseByID(id, userID)
	if err != nil {
		c.SendError(404, "Expense not found")
		return
	}
	c.SendSuccess(200, "Expense retrieved", expense)
}

// Update modifies an existing expense.
// @Title Update Expense
// @Description Update an existing expense by ID
// @Param	X-User-ID	header	int		true	"User ID from login"
// @Param	id			path	int		true	"Expense ID"
// @Param	body	body	controllers.ExpenseUpdateRequest	true	"Expense payload"
// @Success 200 {string} success
// @Failure 404 {string} error
// @Failure 401 {string} error
// @router /:id [put]
func (c *ExpenseController) Update() {
	userID, ok := c.authenticatedUserID()
	if !ok {
		return
	}
	id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err != nil {
		c.SendError(400, "Invalid expense ID")
		return
	}
	expense, err := models.GetExpenseByID(id, userID)
	if err != nil {
		c.SendError(404, "Expense not found")
		return
	}
	var input struct {
		Title       string  `json:"title"`
		Amount      float64 `json:"amount"`
		Category    string  `json:"category"`
		Note        string  `json:"note"`
		ExpenseDate string  `json:"expense_date"`
	}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &input); err != nil {
		c.SendError(400, "Invalid request body")
		return
	}
	if input.Title != "" {
		expense.Title = input.Title
	}
	if input.Amount > 0 {
		expense.Amount = input.Amount
	}
	if input.Category != "" {
		if !models.IsValidCategory(input.Category) {
			c.SendError(400, "Invalid category")
			return
		}
		expense.Category = input.Category
	}
	if input.Note != "" {
		expense.Note = input.Note
	}
	if input.ExpenseDate != "" {
		expense.ExpenseDate = input.ExpenseDate
	}
	if err := models.UpdateExpense(expense); err != nil {
		logs.Error("Update expense:", err)
		c.SendError(500, "Failed to update expense")
		return
	}
	c.SendSuccess(200, "Expense updated successfully", expense)
}

// Delete removes an expense.
// @Title Delete Expense
// @Description Delete an expense by ID
// @Param	X-User-ID	header	int	true	"User ID from login"
// @Param	id			path	int	true	"Expense ID"
// @Success 200 {string} success
// @Failure 404 {string} error
// @Failure 401 {string} error
// @router /:id [delete]
func (c *ExpenseController) Delete() {
	userID, ok := c.authenticatedUserID()
	if !ok {
		return
	}
	id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err != nil {
		c.SendError(400, "Invalid expense ID")
		return
	}
	if err := models.DeleteExpense(id, userID); err != nil {
		c.SendError(404, "Expense not found")
		return
	}
	c.SendSuccess(200, "Expense deleted successfully", nil)
}

// Summary returns total spending within a date range.
// @Title Expense Summary
// @Description Get total spending summary for a date range
// @Param	X-User-ID	header	int		true	"User ID from login"
// @Param	date_from	query	string	true	"Start date (YYYY-MM-DD)"
// @Param	date_to		query	string	true	"End date (YYYY-MM-DD)"
// @Success 200 {string} success
// @Failure 400 {string} error
// @Failure 401 {string} error
// @router /summary [get]
func (c *ExpenseController) Summary() {
	userID, ok := c.authenticatedUserID()
	if !ok {
		return
	}
	dateFrom := c.GetString("date_from")
	dateTo := c.GetString("date_to")
	if dateFrom == "" {
		c.SendError(400, "date_from is required")
		return
	}
	if dateTo == "" {
		c.SendError(400, "date_to is required")
		return
	}
	if dateFrom > dateTo {
		c.SendError(400, "date_from cannot be after date_to")
		return
	}
	summary, err := models.GetSummary(userID, dateFrom, dateTo)
	if err != nil {
		logs.Error("Summary:", err)
		c.SendError(500, "Failed to generate summary")
		return
	}
	c.SendSuccess(200, "Summary generated", summary)
}
