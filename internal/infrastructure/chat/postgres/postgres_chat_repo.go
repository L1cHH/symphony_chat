package postgres

import (
	"context"
	"database/sql"
	"errors"
	"symphony_chat/internal/application/transaction"
	"symphony_chat/internal/domain/chat"
	"time"

	"github.com/google/uuid"
)

type PostgresChatRepo struct {
	db *sql.DB
}

func NewPostgresChatRepo(db *sql.DB) *PostgresChatRepo {
	return &PostgresChatRepo{
		db: db,
	}
}

func (pr *PostgresChatRepo) GetChatByID(ctx context.Context, chat_id uuid.UUID) (chat.Chat, error) {
	tx := pr.GetTransaction(ctx)

	var id uuid.UUID
	var name string
	var createdAt time.Time
	var updatedAt time.Time

	err := tx.QueryRowContext(
		ctx,
		"SELECT id, name, created_at, updated_at FROM chats WHERE id = $1",
		chat_id,
	).Scan(&id, &name, &createdAt, &updatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return chat.Chat{}, chat.ErrChatNotFound
		}

		return chat.Chat{}, &chat.ChatError{
			Code:    "DATABASE_ERROR",
			Message: "failed to get chat",
			Err:     err,
		}
	}

	return chat.ChatFromDB(id, name, createdAt, updatedAt), nil
}

func (pr *PostgresChatRepo) AddChat(ctx context.Context, chatDB chat.Chat) error {
	tx := pr.GetTransaction(ctx)

	_, err := tx.ExecContext(
		ctx,
		"INSERT INTO chats (id, name, created_at, updated_at) VALUES ($1, $2, $3, $4)",
		chatDB.GetID(), chatDB.GetName(), chatDB.GetCreatedAt(), chatDB.GetUpdatedAt(),
	)

	if err != nil {
		return &chat.ChatError {
			Code:    "DATABASE_ERROR",
			Message: "failed to add chat",
			Err:     err,
		}
	}

	return nil
}

func (pr *PostgresChatRepo) UpdateChatName(ctx context.Context, chatId uuid.UUID, name string) error {
	tx := pr.GetTransaction(ctx)

	result, err := tx.ExecContext(
		ctx,
		"UPDATE chat SET name = $1 WHERE id = $2",
		name, chatId,
	)

	if err != nil {
		return &chat.ChatError {
			Code: "DATABASE_ERROR",
			Message: "failed to update chat name",
			Err:     err,
		}
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return &chat.ChatError {
			Code: "DATABASE_ERROR",
			Message: "failed to get affected rows after updated chat name",
			Err:     err,
		}
	}

	if rowsAffected == 0 {
		return chat.ErrChatNotFound
	}

	return nil
}

func (pr *PostgresChatRepo) UpdateChatUpdatedAt(ctx context.Context, chatID uuid.UUID, updatedAt time.Time) error {
	tx := pr.GetTransaction(ctx)

	result, err := tx.ExecContext(
		ctx,
		"UPDATE chat SET updated_at = $1 WHERE id = $2",
		updatedAt,
		chatID,
	)

	if err != nil {
		return &chat.ChatError{
			Code: "DATABASE_ERROR",
			Message: "failed to update chat updated_at",
			Err:     err,
		}
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return &chat.ChatError {
			Code: "DATABASE_ERROR",
			Message: "failed to update chat updated_at",
			Err:     err,
		}
	}

	if rowsAffected == 0 {
		return chat.ErrChatNotFound
	}

	return nil
}

func (pr *PostgresChatRepo) DeleteChat(ctx context.Context, chatId uuid.UUID) error {
	tx := pr.GetTransaction(ctx)

	result, err := tx.ExecContext(
		ctx,
		"DELETE FROM chat WHERE id = $1",
		chatId,
	)

	if err != nil {
		return &chat.ChatError {
			Code: "DATABASE_ERROR",
			Message: "failed to delete chat",
			Err:     err,
		}
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return &chat.ChatError{
			Code: "DATABASE_ERROR",
			Message: "failed to delete chat",
			Err:     err,
		}
	}
	
	if rowsAffected == 0 {
		return chat.ErrChatNotFound
	}

	return nil
}

func (pr *PostgresChatRepo) GetTransaction(ctx context.Context) transaction.DBTX {
	if tx := transaction.IsTransaction(ctx); tx != nil {
		return tx
	}
	return pr.db
}
