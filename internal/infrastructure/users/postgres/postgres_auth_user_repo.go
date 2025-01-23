package postgres

import (
	"database/sql"
	"fmt"
	"symphony_chat/internal/domain/users"
	"time"

	"github.com/google/uuid"
)

type PostgresAuthUserRepo struct {
	db *sql.DB
}

func NewPostgresAuthUserRepo(db *sql.DB) *PostgresAuthUserRepo {
	return &PostgresAuthUserRepo{
		db: db,
	}
}

func (pr *PostgresAuthUserRepo) GetAuthUserById(user_id uuid.UUID) (users.AuthUser, error) {
	var id uuid.UUID
	var login, password string
	var registrationAt time.Time
	err := pr.db.QueryRow(
		"SELECT id, login, password, registration_at FROM auth_user WHERE id = $1", user_id,
	).Scan(&id, &login, &password, &registrationAt)

	if err != nil {
		return users.AuthUser{}, fmt.Errorf("failed to get auth_user by id: %w", err)
	}

	return users.NewAuthUser(id, login, password, registrationAt), nil
}

func (pr *PostgresAuthUserRepo) GetAuthUserByLogin(user_login string) (users.AuthUser, error) {
	var id uuid.UUID
	var login, password string
	var registrationAt time.Time
	err := pr.db.QueryRow(
		"SELECT id, login, password, registration_at FROM auth_user WHERE login = $1", user_login,
	).Scan(&id, &login, &password, &registrationAt)

	if err != nil {
		return users.AuthUser{}, fmt.Errorf("failed to get auth_user by login: %w", err)
	}

	return users.NewAuthUser(id, login, password, registrationAt), nil
}

func (pr *PostgresAuthUserRepo) IsUserExists(user_login string) (bool, error) {
	var id uuid.UUID

	err := pr.db.QueryRow(
		"SELECT id FROM auth_user WHERE login = $1", user_login,
	).Scan(&id)

	if err != nil && err != sql.ErrNoRows {
		return false, fmt.Errorf("failed to check if user exists: %w", err)
	}

	if err == sql.ErrNoRows {
		return false, nil
	}

	return true, nil
}

func (pr *PostgresAuthUserRepo) AddAuthUser(user users.AuthUser) error {
	_, err := pr.db.Exec(
		"INSERT INTO auth_user (id, login, password, registration_at) VALUES ($1, $2, $3, $4)",
		user.GetID(), user.GetLogin(), user.GetPassword(), user.GetRegistrationAt(),
	)
	if err != nil {
		return fmt.Errorf("failed to add auth_user: %w", err)
	}

	return nil
}

func (pr *PostgresAuthUserRepo) UpdateLogin(user_id uuid.UUID, new_login string) error {
	_, err := pr.db.Exec(
		"UPDATE auth_user SET login = $1 WHERE id = $2",
		new_login, user_id,
	)
	if err != nil {
		return fmt.Errorf("failed to update login: %w", err)
	}

	return nil
}

func (pr *PostgresAuthUserRepo) UpdatePassword(user_id uuid.UUID, new_password string) error {
	_, err := pr.db.Exec(
		"UPDATE auth_user SET password = $1 WHERE id = $2",
		new_password, user_id,
	)
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	return nil
}

func (pr *PostgresAuthUserRepo) DeleteAuthUser(user_id uuid.UUID) error {
	_, err := pr.db.Exec(
		"DELETE FROM auth_user WHERE id = $1",
		user_id,
	)
	if err != nil {
		return fmt.Errorf("failed to delete auth_user: %w", err)
	}

	return nil
}
