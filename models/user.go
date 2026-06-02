package models

import (
	"encoding/csv"
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

// UserFile is the CSV path for users. Override in tests.
var UserFile = "data/users.csv"

// User represents a registered account.
type User struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	CreatedAt string `json:"created_at"`
}

// ensureUserFile creates the file and parent directory if missing.
func ensureUserFile() error {
	if err := os.MkdirAll(filepath.Dir(UserFile), 0755); err != nil {
		return err
	}
	if _, err := os.Stat(UserFile); os.IsNotExist(err) {
		f, err := os.Create(UserFile)
		if err != nil {
			return err
		}
		defer f.Close()
		w := csv.NewWriter(f)
		if err := w.Write([]string{"id", "name", "email", "password", "created_at"}); err != nil {
			return err
		}
		w.Flush()
		return w.Error()
	}
	return nil
}

// GetAllUsers reads every user from the CSV.
func GetAllUsers() ([]User, error) {
	if err := ensureUserFile(); err != nil {
		return nil, err
	}
	f, err := os.Open(UserFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	rows, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return nil, err
	}

	users := make([]User, 0, len(rows))
	for _, row := range rows[1:] { // skip header
		id, _ := strconv.Atoi(row[0])
		users = append(users, User{
			ID:        id,
			Name:      row[1],
			Email:     row[2],
			Password:  row[3],
			CreatedAt: row[4],
		})
	}
	return users, nil
}

// GetUserByEmail finds a user by email address.
func GetUserByEmail(email string) (*User, error) {
	users, err := GetAllUsers()
	if err != nil {
		return nil, err
	}
	for i := range users {
		if users[i].Email == email {
			return &users[i], nil
		}
	}
	return nil, errors.New("user not found")
}

// GetUserByID finds a user by their numeric ID.
func GetUserByID(id int) (*User, error) {
	users, err := GetAllUsers()
	if err != nil {
		return nil, err
	}
	for i := range users {
		if users[i].ID == id {
			return &users[i], nil
		}
	}
	return nil, errors.New("user not found")
}

// nextUserID returns max existing ID + 1.
func nextUserID() (int, error) {
	users, err := GetAllUsers()
	if err != nil {
		return 0, err
	}
	max := 0
	for _, u := range users {
		if u.ID > max {
			max = u.ID
		}
	}
	return max + 1, nil
}

// CreateUser appends a new user row and sets ID and CreatedAt.
func CreateUser(user *User) error {
	if err := ensureUserFile(); err != nil {
		return err
	}
	id, err := nextUserID()
	if err != nil {
		return err
	}
	user.ID = id
	user.CreatedAt = time.Now().UTC().Format(time.RFC3339)

	f, err := os.OpenFile(UserFile, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	w := csv.NewWriter(f)
	if err := w.Write([]string{
		strconv.Itoa(user.ID),
		user.Name,
		user.Email,
		user.Password,
		user.CreatedAt,
	}); err != nil {
		return err
	}
	w.Flush()
	return w.Error()
}