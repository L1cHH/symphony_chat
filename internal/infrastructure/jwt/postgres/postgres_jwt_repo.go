package postgres

import (
	"database/sql"
	"fmt"
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

func (pr *PostgresJWTtokenRepo) AddJWTtoken(token jwt.JWTtoken) error {
	_, err := pr.db.Exec(
		"INSERT INTO jwt_token (auth_user_id, token) VALUES ($1, $2)",
		token.GetAuthUserID(), token.GetToken(),
	)
	if err != nil {
		return fmt.Errorf("failed to add jwt_token: %w", err)
	}

	return nil
}

func (pr *PostgresJWTtokenRepo) GetJWTtoken(userID uuid.UUID) (jwt.JWTtoken, error) {
	var authUserID uuid.UUID
	var token string
	err := pr.db.QueryRow(
		"SELECT auth_user_id, token FROM jwt_token WHERE auth_user_id = $1", userID,
	).Scan(&authUserID, &token)

	if err != nil {
		return jwt.JWTtoken{}, fmt.Errorf("failed to get jwt_token: %w", err)
	}

	return jwt.FromDB(authUserID, token), nil
}

func (pr *PostgresJWTtokenRepo) UpdateJWTtoken(authUserID uuid.UUID, newToken string) error {
	_, err := pr.db.Exec(
		"UPDATE jwt_token SET token = $1 WHERE auth_user_id = $2",
		newToken, authUserID,
	)
	if err != nil {
		return fmt.Errorf("failed to update jwt_token: %w", err)
	}

	return nil
}

func (pr *PostgresJWTtokenRepo) DeleteJWTtoken(authUserID uuid.UUID) error {
	_, err := pr.db.Exec(
		"DELETE FROM jwt_token WHERE auth_user_id = $1",
		authUserID,
	)
	if err != nil {
		return fmt.Errorf("failed to delete jwt_token: %w", err)
	}

	return nil
}
