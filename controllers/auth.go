package controllers

import (
	"encoding/json"
	"expense-tracker-api/models"
	"regexp"

	"github.com/beego/beego/v2/core/logs"
)

// AuthController handles registration and login.
type AuthController struct {
	BaseController
}

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

// Register creates a new user account.
// @Title Register
// @Description Create a new user with name, email and password
// @Param	body	body	controllers.RegisterRequest	true	"Register payload"
// @Success 201 {string} success
// @Failure 400 {string} error
// @Failure 409 {string} error
// @router /register [post]
func (c *AuthController) Register() {
	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &input); err != nil {
		c.SendError(400, "Invalid request body")
		return
	}
	switch {
	case input.Name == "":
		c.SendError(400, "Name is required")
		return
	case input.Email == "":
		c.SendError(400, "Email is required")
		return
	case !emailRegex.MatchString(input.Email):
		c.SendError(400, "Invalid email format")
		return
	case len(input.Password) < 6:
		c.SendError(400, "Password must be at least 6 characters")
		return
	}
	if existing, _ := models.GetUserByEmail(input.Email); existing != nil {
		c.SendError(409, "Email already exists")
		return
	}
	user := &models.User{Name: input.Name, Email: input.Email, Password: input.Password}
	if err := models.CreateUser(user); err != nil {
		logs.Error("Register: failed to create user:", err)
		c.SendError(500, "Failed to register user")
		return
	}
	c.SendSuccess(201, "User registered successfully", nil)
}

// Login authenticates a user.
// @Title Login
// @Description Login with email and password, returns user_id
// @Param	body	body	controllers.LoginRequest	true	"Login payload"
// @Success 200 {string} success
// @Failure 400 {string} error
// @Failure 401 {string} error
// @router /login [post]
func (c *AuthController) Login() {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &input); err != nil {
		c.SendError(400, "Invalid request body")
		return
	}
	if input.Email == "" || input.Password == "" {
		c.SendError(400, "Email and password are required")
		return
	}
	user, err := models.GetUserByEmail(input.Email)
	if err != nil || user.Password != input.Password {
		c.SendError(401, "Invalid email or password")
		return
	}
	c.SendSuccess(200, "Login successful", map[string]interface{}{
		"user_id": user.ID,
		"name":    user.Name,
		"email":   user.Email,
	})
}
