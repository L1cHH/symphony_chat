package postgres

import (
	"database/sql"
	"fmt"
	"symphony_chat/internal/domain/users"
	"time"

	"github.com/google/uuid"
)

type PostgresChatUserRepo struct {
	db *sql.DB
}

func NewPostgresChatUserRepo(db *sql.DB) *PostgresChatUserRepo {
	return &PostgresChatUserRepo{
		db: db,
	}
}

func (pr *PostgresChatUserRepo) GetChatUserById(chat_user_id uuid.UUID) (users.ChatUser, error) {
	var id uuid.UUID
	var username string
	var status users.UserStatus
	var created_at time.Time
	var last_seen_at time.Time

	err := pr.db.QueryRow(
		"SELECT id, username, status, created_at, last_seen_at FROM chat_user WHERE id = $1",
		chat_user_id,
	).Scan(&id, &username, &status, &created_at, &last_seen_at)

	if err != nil {
		return users.ChatUser{}, fmt.Errorf("failed to get chat_user by id: %w", err)
	}

	return users.ChatUserFromDB(id, username, status, created_at, last_seen_at), nil
}

func (pr *PostgresChatUserRepo) GetChatUserByUsername(chat_user_username string) (users.ChatUser, error) {
	var id uuid.UUID
	var username string
	var status users.UserStatus
	var created_at time.Time
	var last_seen_at time.Time

	err := pr.db.QueryRow(
		"SELECT id, username, status, created_at, last_seen_at FROM chat_user WHERE username = $1",
		chat_user_username,
	).Scan(&id, &username, &status, &created_at, &last_seen_at)

	if err != nil {
		return users.ChatUser{}, fmt.Errorf("failed to get chat_user by username: %w", err)
	}

	return users.ChatUserFromDB(id, username, status, created_at, last_seen_at), nil
}

func (pr *PostgresChatUserRepo) AddChatUser(chat_user users.ChatUser) error {
	_, err := pr.db.Exec(
		"INSERT INTO chat_user (id, username, status, created_at, last_seen_at) VALUES ($1, $2, $3, $4, $5)",
		chat_user.GetID(), chat_user.GetUsername(), chat_user.GetStatus(), chat_user.GetCreatedAt(), chat_user.GetLastSeenAt(),
	)
	if err != nil {
		return fmt.Errorf("failed to add chat_user: %w", err)
	}

	return nil
}

func (pr *PostgresChatUserRepo) DeleteChatUserByID(chat_user_id uuid.UUID) error {
	_, err := pr.db.Exec(
		"DELETE FROM chat_user WHERE id = $1",
		chat_user_id,
	)
	if err != nil {
		return fmt.Errorf("failed to delete chat_user: %w", err)
	}

	return nil
}

func (pr *PostgresChatUserRepo) UpdateUsername(chat_user_id uuid.UUID, new_username string) error {
	_, err := pr.db.Exec(
		"UPDATE chat_user SET username = $1 WHERE id = $2",
		new_username, chat_user_id,
	)
	if err != nil {
		return fmt.Errorf("failed to update username for chat_user: %w", err)
	}

	return nil
}

func (pr *PostgresChatUserRepo) UpdateStatus(chat_user_id uuid.UUID, new_status users.UserStatus) error {
	_, err := pr.db.Exec(
		"UPDATE chat_user SET status = $1 WHERE id = $2",
		new_status, chat_user_id,
	)
	if err != nil {
		return fmt.Errorf("failed to update status for chat_user: %w", err)
	}

	return nil
}

func (pr *PostgresChatUserRepo) UpdateLastSeenAt(chat_user_id uuid.UUID, new_last_seen_at time.Time) error {
	_, err := pr.db.Exec(
		"UPDATE chat_user SET last_seen_at = $1 WHERE id = $2",
		new_last_seen_at, chat_user_id,
	)
	if err != nil {
		return fmt.Errorf("failed to update last_seen_at for chat_user: %w", err)
	}

	return nil
}
