package postgres

import (
	"context"
	"database/sql"
	"errors"
	"symphony_chat/internal/application/transaction"
	"symphony_chat/internal/domain/messages"
	"time"

	"github.com/google/uuid"
)

type PostgresChatMessageRepo struct {
	db *sql.DB
}

func NewPostgresChatMessageRepo(db *sql.DB) PostgresChatMessageRepo {
	return PostgresChatMessageRepo{
		db: db,
	}
}

func (pr *PostgresChatMessageRepo) GetChatMessageById(ctx context.Context, messageID uuid.UUID) (messages.ChatMessage, error) {
	tx := pr.GetTransaction(ctx)

	var id uuid.UUID
	var chatID uuid.UUID
	var senderID uuid.UUID
	var content string
	var createdAt time.Time
	var status messages.MessageStatus
	

	err := tx.QueryRowContext(
		ctx,
		`SELECT id, chat_id, sender_id, content, created_at, status
		FROM chat_message WHERE id = $1`,
		messageID,
	).Scan(
		&id,
		&chatID,
		&senderID,
		&content,
		&createdAt,
		&status,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return messages.ChatMessage{}, messages.ErrChatMessageNotFound
		}

		return messages.ChatMessage{}, &messages.ChatMessageError {
			Code: "DATABASE_ERROR",
			Message: "failed to get chat message by id",
			Err: err,
		}
	}

	return messages.ChatMessageFromDB(id, chatID, senderID, content, createdAt, status), nil
}

func (pr *PostgresChatMessageRepo) GetChatMessagesByChatID(ctx context.Context, chatID uuid.UUID) ([]messages.ChatMessage, error) {
	tx := pr.GetTransaction(ctx)

	rows, err := tx.QueryContext(
		ctx,
		`SELECT id, sender_id, content, created_at, status
		FROM chat_message WHERE chat_id = $1`,
		chatID,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []messages.ChatMessage{}, messages.ErrChatMessageNotFound
		}

		return []messages.ChatMessage{}, &messages.ChatMessageError {
			Code: "DATABASE_ERROR",
			Message: "failed to get chat messages by chat id",
			Err: err,
		}
	}

	defer rows.Close()

	var chatMessages = make([]messages.ChatMessage, 0)

	for rows.Next() {
		var id uuid.UUID
		var senderID uuid.UUID
		var content string
		var createdAt time.Time
		var status messages.MessageStatus

		if err := rows.Scan(&id, &senderID, &content, &createdAt, &status); err != nil {
			return nil, &messages.ChatMessageError {
				Code: "DATABASE_ERROR",
				Message: "failed to scan chat message",
				Err: err,
			}
		}

		chatMessages = append(chatMessages, messages.ChatMessageFromDB(id, chatID, senderID, content, createdAt, status))
	}

	return chatMessages, nil
}

func (pr *PostgresChatMessageRepo) GetChatMessagesByContentAndChatID(ctx context.Context, content string, chatID uuid.UUID) ([]messages.ChatMessage, error) {
	tx := pr.GetTransaction(ctx)

	rows, err := tx.QueryContext(
		ctx,
		`SELECT id, sender_id, created_at, status
		FROM chat_message WHERE content = $1 AND chat_id = $2`,
		content,
		chatID,
	)

	if err != nil {
		return nil, &messages.ChatMessageError {
			Code: "DATABASE_ERROR",
			Message: "failed to get chat messages by content and chat id",
			Err: err,
		}
	}

	defer rows.Close()

	var foundChatMessages = make([]messages.ChatMessage, 0)

	for rows.Next() {
		var id uuid.UUID
		var senderID uuid.UUID
		var createdAt time.Time
		var status messages.MessageStatus

		if err := rows.Scan(&id, &senderID, &createdAt, &status); err != nil {
			return nil, &messages.ChatMessageError {
				Code: "DATABASE_ERROR",
				Message: "failed to scan chat message",
				Err: err,
			}
		}

		foundChatMessages = append(foundChatMessages, messages.ChatMessageFromDB(id, chatID, senderID, content, createdAt, status))
	}

	return foundChatMessages, nil
}

func (pr *PostgresChatMessageRepo) AddChatMessage(ctx context.Context, chatMessage messages.ChatMessage) error {
	tx := pr.GetTransaction(ctx)

	result, err := tx.ExecContext(
		ctx,
		`INSERT INTO chat_message (id, chat_id, sender_id, content, created_at, status)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		chatMessage.GetID(),
		chatMessage.GetChatID(),
		chatMessage.GetSenderID(),
		chatMessage.GetContent(),
		chatMessage.GetCreatedAt(),
		chatMessage.GetStatus(),
	)

	if err != nil {
		return &messages.ChatMessageError{
			Code: "DATABASE_ERROR",
			Message: "failed to add chat message",
			Err: err,
		}
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return &messages.ChatMessageError{
			Code: "DATABASE_ERROR",
			Message: "failed to get affected rows after added chat message",
			Err: err,
		}
	}

	if rowsAffected == 0 {
		return messages.ErrChatMessageNotFound
	}

	return nil
}

func (pr *PostgresChatMessageRepo) UpdateChatMessageContent(ctx context.Context, messageID uuid.UUID, content string) error {
	tx := pr.GetTransaction(ctx)

	result, err := tx.ExecContext(
		ctx,
		`UPDATE chat_message SET content = $1 WHERE id = $2`,
		content,
		messageID,
	)

	if err != nil {
		return &messages.ChatMessageError{
			Code: "DATABASE_ERROR",
			Message: "failed to update chat message content",
			Err: err,
		}
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return &messages.ChatMessageError{
			Code: "DATABASE_ERROR",
			Message: "failed to get affected rows after updated chat message content",
			Err: err,
		}
	}

	if rowsAffected == 0 {
		return messages.ErrChatMessageNotFound
	}

	return nil
}

func (pr *PostgresChatMessageRepo) UpdateChatMessageStatus(ctx context.Context, messageID uuid.UUID, status messages.MessageStatus) error {
	tx := pr.GetTransaction(ctx)

	result, err := tx.ExecContext (
		ctx,
		`UPDATE chat_message SET status = $1 WHERE id = $2`,
		status,
		messageID,
	)

	if err != nil {
		return &messages.ChatMessageError{
			Code: "DATABASE_ERROR",
			Message: "failed to update chat message status",
			Err: err,
		}
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return &messages.ChatMessageError{
			Code: "DATABASE_ERROR",
			Message: "failed to get affected rows after updated chat message status",
			Err: err,
		}
	}

	if rowsAffected == 0 {
		return messages.ErrChatMessageNotFound
	}

	return nil
}

func (pr *PostgresChatMessageRepo) DeleteChatMessage(ctx context.Context, messageID uuid.UUID) error {
	tx := pr.GetTransaction(ctx)

	result, err := tx.ExecContext(
		ctx,
		`DELETE FROM chat_message WHERE id = $1`,
		messageID,
	)

	if err != nil {
		return &messages.ChatMessageError{
			Code: "DATABASE_ERROR",
			Message: "failed to delete chat message",
			Err: err,
		}
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return &messages.ChatMessageError{
			Code: "DATABASE_ERROR",
			Message: "failed to get affected rows after deleted chat message",
			Err: err,
		}
	}

	if rowsAffected == 0 {
		return messages.ErrChatMessageNotFound
	}

	return nil
}

func (pr *PostgresChatMessageRepo) DeleteAllChatMessagesByChatID(ctx context.Context, chatID uuid.UUID) error {
	tx := pr.GetTransaction(ctx)

	result, err := tx.ExecContext(
		ctx,
		`DELETE FROM chat_message WHERE chat_id = $1`,
		chatID,
	)

	if err != nil {
		return &messages.ChatMessageError{
			Code: "DATABASE_ERROR",
			Message: "failed to delete all chat messages by chat ID",
			Err: err,
		}
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return &messages.ChatMessageError{
			Code: "DATABASE_ERROR",
			Message: "failed to get affected rows after deleted all chat messages by chat ID",
			Err: err,
		}
	}

	if rowsAffected == 0 {
		return messages.ErrChatMessageNotFound
	}

	return nil
}

func (pr *PostgresChatMessageRepo) GetTransaction(ctx context.Context) transaction.DBTX {
	if tx := transaction.IsTransaction(ctx); tx != nil {
		return tx
	}
	return pr.db
}