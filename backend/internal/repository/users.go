package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"pharmasense/backend/internal/models"
)

func (r *Repository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, pharmacy_id, email, password_hash, full_name, role, is_active, last_login_at, created_at
		FROM users WHERE email = $1 AND is_active = true`, email)

	var u models.User
	err := row.Scan(&u.ID, &u.PharmacyID, &u.Email, &u.PasswordHash, &u.FullName,
		&u.Role, &u.IsActive, &u.LastLoginAt, &u.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("get user by email: %w", err)
	}
	return &u, nil
}

func (r *Repository) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, pharmacy_id, email, password_hash, full_name, role, is_active, last_login_at, created_at
		FROM users WHERE id = $1`, id)

	var u models.User
	err := row.Scan(&u.ID, &u.PharmacyID, &u.Email, &u.PasswordHash, &u.FullName,
		&u.Role, &u.IsActive, &u.LastLoginAt, &u.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("get user by id: %w", err)
	}
	return &u, nil
}

func (r *Repository) UpdateLastLogin(ctx context.Context, userID uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `UPDATE users SET last_login_at = NOW() WHERE id = $1`, userID)
	return err
}

func (r *Repository) ListUsers(ctx context.Context, pharmacyID uuid.UUID) ([]models.User, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, pharmacy_id, email, password_hash, full_name, role, is_active, last_login_at, created_at
		FROM users WHERE pharmacy_id = $1 ORDER BY created_at`, pharmacyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.PharmacyID, &u.Email, &u.PasswordHash, &u.FullName,
			&u.Role, &u.IsActive, &u.LastLoginAt, &u.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}
