package models

import "testing"

func TestCreateUser(t *testing.T) {
	tests := []struct {
		name    string
		user    User
		wantErr bool
	}{
		{
			name:    "valid user",
			user:    User{Name: "Alice", Email: "alice@example.com", Password: "pass123"},
			wantErr: false,
		},
		{
			name:    "second valid user",
			user:    User{Name: "Bob", Email: "bob@example.com", Password: "pass456"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CreateUser(&tt.user)
			if (err != nil) != tt.wantErr {
				t.Fatalf("CreateUser() error = %v, wantErr = %v", err, tt.wantErr)
			}
			if !tt.wantErr && tt.user.ID == 0 {
				t.Error("expected ID to be set after CreateUser")
			}
		})
	}
}

func TestGetUserByEmail(t *testing.T) {
	seed := makeUser(t, "charlie@example.com")

	tests := []struct {
		name      string
		email     string
		wantFound bool
	}{
		{name: "existing email", email: seed.Email, wantFound: true},
		{name: "unknown email", email: "ghost@example.com", wantFound: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u, err := GetUserByEmail(tt.email)
			if tt.wantFound {
				if err != nil {
					t.Fatalf("expected user, got error: %v", err)
				}
				if u.Email != tt.email {
					t.Errorf("expected email %s, got %s", tt.email, u.Email)
				}
			} else {
				if err == nil {
					t.Error("expected error for unknown email, got nil")
				}
			}
		})
	}
}

func TestGetUserByID(t *testing.T) {
	seed := makeUser(t, "dave@example.com")

	tests := []struct {
		name      string
		id        int
		wantFound bool
	}{
		{name: "existing ID", id: seed.ID, wantFound: true},
		{name: "unknown ID", id: 99999, wantFound: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u, err := GetUserByID(tt.id)
			if tt.wantFound {
				if err != nil {
					t.Fatalf("expected user, got error: %v", err)
				}
				if u.ID != tt.id {
					t.Errorf("expected ID %d, got %d", tt.id, u.ID)
				}
			} else {
				if err == nil {
					t.Error("expected error for unknown ID, got nil")
				}
			}
		})
	}
}