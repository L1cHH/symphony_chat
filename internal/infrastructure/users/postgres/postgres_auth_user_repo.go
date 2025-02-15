package postgres

import (
	"context"
	"database/sql"
	"errors"
	"symphony_chat/internal/application/transaction"
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

func (pr *PostgresAuthUserRepo) GetAuthUserById(ctx context.Context, user_id uuid.UUID) (users.AuthUser, error) {
	var id uuid.UUID
	var login, password string
	var registrationAt time.Time

	tx := pr.GetTransaction(ctx)

	err := tx.QueryRowContext(
		ctx,
		"SELECT id, login, password, registration_at FROM auth_user WHERE id = $1", user_id,
	).Scan(&id, &login, &password, &registrationAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return users.AuthUser{}, users.ErrAuthUserNotFound
		}

		return users.AuthUser{}, &users.AuthError{
			Code: "DATABASE_ERROR",
			Message: "failed to get auth_user by id",
			Err: err,
		}
	}

	return users.NewAuthUser(id, login, password, registrationAt), nil
}

func (pr *PostgresAuthUserRepo) GetAuthUserByLogin(ctx context.Context, user_login string) (users.AuthUser, error) {
	var id uuid.UUID
	var login, password string
	var registrationAt time.Time

	tx := pr.GetTransaction(ctx)

	err := tx.QueryRowContext(
		ctx,
		"SELECT id, login, password, registration_at FROM auth_user WHERE login = $1", user_login,
	).Scan(&id, &login, &password, &registrationAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return users.AuthUser{}, users.ErrAuthUserNotFound
		}

		return users.AuthUser{}, &users.AuthError{
			Code: "DATABASE_ERROR",
			Message: "failed to get auth_user by login",
			Err: err,
		}
	}

	return users.NewAuthUser(id, login, password, registrationAt), nil
}

func (pr *PostgresAuthUserRepo) IsUserExists(ctx context.Context, user_login string) (bool, error) {
	var id uuid.UUID

	tx := pr.GetTransaction(ctx)

	err := tx.QueryRowContext(
		ctx,
		"SELECT id FROM auth_user WHERE login = $1", user_login,
	).Scan(&id)

	if err != nil && err != sql.ErrNoRows {
		return false, &users.AuthError{
			Code: "DATABASE_ERROR",
			Message: "failed to check if user exists",
			Err: err,
		}
	}

	if err == sql.ErrNoRows {
		return false, nil
	}

	return true, nil
}

func (pr *PostgresAuthUserRepo) AddAuthUser(ctx context.Context, user users.AuthUser) error {
	tx := pr.GetTransaction(ctx)

	_, err := tx.ExecContext(
		ctx,
		"INSERT INTO auth_user (id, login, password, registration_at) VALUES ($1, $2, $3, $4)",
		user.GetID(), user.GetLogin(), user.GetPassword(), user.GetRegistrationAt(),
	)
	if err != nil {
		return &users.AuthError{
			Code: "DATABASE_ERROR",
			Message: "failed to add auth_user",
			Err: err,
		}
	}

	return nil
}

func (pr *PostgresAuthUserRepo) UpdateLogin(ctx context.Context, user_id uuid.UUID, new_login string) error {
	tx := pr.GetTransaction(ctx)

	_, err := tx.ExecContext(
		ctx,
		"UPDATE auth_user SET login = $1 WHERE id = $2",
		new_login, user_id,
	)
	if err != nil {
		return &users.AuthError{
			Code: "DATABASE_ERROR",
			Message: "failed to update login",
			Err: err,
		}
	}

	return nil
}

func (pr *PostgresAuthUserRepo) UpdatePassword(ctx context.Context, user_id uuid.UUID, new_password string) error {
	
	tx := pr.GetTransaction(ctx)
	_, err := tx.ExecContext(
		ctx,
		"UPDATE auth_user SET password = $1 WHERE id = $2",
		new_password, user_id,
	)
	if err != nil {
		return &users.AuthError{
			Code: "DATABASE_ERROR",
			Message: "failed to update password",
			Err: err,
		}
	}

	return nil
}

func (pr *PostgresAuthUserRepo) DeleteAuthUser(ctx context.Context, user_id uuid.UUID) error {
	
	tx := pr.GetTransaction(ctx)
	
	_, err := tx.ExecContext(
		ctx,
		"DELETE FROM auth_user WHERE id = $1",
		user_id,
	)
	if err != nil {
		return &users.AuthError{
			Code: "DATABASE_ERROR",
			Message: "failed to delete auth_user",
			Err: err,
		}
	}

	return nil
}

//Function that gets transaction from context
//If there is no transaction in context, it returns pr.db
func (pr *PostgresAuthUserRepo) GetTransaction(ctx context.Context) transaction.DBTX {
	if tx := transaction.IsTransaction(ctx); tx != nil {
		return tx
	}

	return pr.db
}