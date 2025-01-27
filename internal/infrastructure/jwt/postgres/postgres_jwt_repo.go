package postgres

import (
	"context"
	"database/sql"
	"errors"
	"symphony_chat/internal/application/transaction"
	"symphony_chat/internal/domain/jwt"

	"github.com/google/uuid"
)

type PostgresJWTtokenRepo struct {
	db *sql.DB
}

func NewPostgresJWTtokenRepo(db *sql.DB) *PostgresJWTtokenRepo {
	return &PostgresJWTtokenRepo{
		db: db,
	}
}

func (pr *PostgresJWTtokenRepo) AddJWTtoken(ctx context.Context, token jwt.JWTtoken) error {

	tx := pr.GetTransaction(ctx)

	_, err := tx.ExecContext(
		ctx,
		"INSERT INTO jwt_token (auth_user_id, token) VALUES ($1, $2)",
		token.GetAuthUserID(), token.GetToken(),
	)
	if err != nil {
		return &jwt.TokenError{
			Code: "DATABASE_ERROR",
			Message: "failed to add jwt_token",
			Err: err,
		}
	}

	return nil
}

func (pr *PostgresJWTtokenRepo) GetJWTtoken(ctx context.Context, userID uuid.UUID) (jwt.JWTtoken, error) {
	var authUserID uuid.UUID
	var token string

	tx := pr.GetTransaction(ctx)

	err := tx.QueryRowContext(
		ctx,
		"SELECT auth_user_id, token FROM jwt_token WHERE auth_user_id = $1", userID,
	).Scan(&authUserID, &token)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return jwt.JWTtoken{}, jwt.ErrTokenNotFound
		}
		
		return jwt.JWTtoken{}, &jwt.TokenError{
			Code: "DATABASE_ERROR",
			Message: "failed to get jwt_token",
			Err: err,
		}
	}

	return jwt.FromDB(authUserID, token), nil
}

func (pr *PostgresJWTtokenRepo) UpdateJWTtoken(ctx context.Context, authUserID uuid.UUID, newToken string) error {
	
	tx := pr.GetTransaction(ctx)

	_, err := tx.ExecContext(
		ctx,
		"UPDATE jwt_token SET token = $1 WHERE auth_user_id = $2",
		newToken, authUserID,
	)
	if err != nil {
		return &jwt.TokenError{
			Code: "DATABASE_ERROR",
			Message: "failed to update jwt_token",
			Err: err,
		}
	}

	return nil
}

func (pr *PostgresJWTtokenRepo) DeleteJWTtoken(ctx context.Context, authUserID uuid.UUID) error {
	
	tx := pr.GetTransaction(ctx)
	
	_, err := tx.ExecContext(
		ctx,
		"DELETE FROM jwt_token WHERE auth_user_id = $1",
		authUserID,
	)
	if err != nil {
		return &jwt.TokenError{
			Code: "DATABASE_ERROR",
			Message: "failed to delete jwt_token",
			Err: err,
		}
	}

	return nil
}

func (pr *PostgresJWTtokenRepo) GetTransaction(ctx context.Context) transaction.DBTX {
	if tx := transaction.IsTransaction(ctx); tx != nil {
		return tx
	}

	return pr.db
}
