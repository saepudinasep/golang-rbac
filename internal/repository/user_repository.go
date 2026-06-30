package repository

import (
	"context"
	"database/sql"
	"golang-rbac/internal/domain"
)

type mysqlUserRepository struct {
	db *sql.DB
}

func NewMysqlUserRepository(db *sql.DB) domain.UserRepository {
	return &mysqlUserRepository{db: db}
}

func (m *mysqlUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
		SELECT u.id, u.name, u.email, u.password, u.role_id, r.name 
		FROM users u 
		LEFT JOIN roles r ON u.role_id = r.id 
		WHERE u.email = ?`

	user := &domain.User{}
	err := m.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID, &user.Name, &user.Email, &user.Password, &user.RoleID, &user.RoleName,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (m *mysqlUserRepository) Create(ctx context.Context, user *domain.User) error {
	query := "INSERT INTO users (name, email, password, role_id) VALUES (?, ?, ?, ?)"
	_, err := m.db.ExecContext(ctx, query, user.Name, user.Email, user.Password, user.RoleID)
	return err
}
