package domain

import "context"

type Role struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"-"` // Sembunyikan saat di-marshal ke JSON
	RoleID   int    `json:"role_id"`
	RoleName string `json:"role_name,omitempty"`
}

// Kontrak untuk Data Layer
type UserRepository interface {
	GetByEmail(ctx context.Context, email string) (*User, error)
	Create(ctx context.Context, user *User) error
}

// Kontrak untuk Business Logic Layer
type UserUsecase interface {
	Login(ctx context.Context, email, password string) (string, *User, error) // Mengembalikan token & user info
	Register(ctx context.Context, user *User) error
}
