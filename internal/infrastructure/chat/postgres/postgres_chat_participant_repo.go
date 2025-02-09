package postgres

import (
	"context"
	"database/sql"
	"errors"
	"symphony_chat/internal/application/transaction"
	"symphony_chat/internal/domain/chat_participant"
	"time"

	"github.com/google/uuid"
)

type PostgresChatParticipantRepo struct {
	db *sql.DB
}

func (pr *PostgresChatParticipantRepo) GetChatParticipantByIDs(ctx context.Context, chatID uuid.UUID, userID uuid.UUID) (chatparticipant.ChatParticipant, error) {
	tx := pr.GetTransaction(ctx)

	var foundChatID uuid.UUID
	var foundUserID uuid.UUID
	var roleID uuid.UUID
	var joinedAt time.Time

	err := tx.QueryRowContext(
		ctx,
		`SELECT chat_id, user_id, role_id, joined_at
		FROM chat_participant WHERE chat_id = $1 AND user_id = $2`,
		chatID,
		userID,
	).Scan(&foundChatID, &foundUserID, &roleID, &joinedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return chatparticipant.ChatParticipant{}, chatparticipant.ErrChatParticipantNotFound
		}

		return chatparticipant.ChatParticipant{}, &chatparticipant.ChatParticipantError{
			Code:    "DATABASE_ERROR",
			Message: "failed to get chat participant",
			Err:     err,
		}
	}

	return chatparticipant.ChatParticipantFromDB(foundChatID, foundUserID, roleID, joinedAt), nil
}

func (pr *PostgresChatParticipantRepo) GetAllChatParticipantsByChatID(ctx context.Context, chatID uuid.UUID) ([]chatparticipant.ChatParticipant, error) {
	tx := pr.GetTransaction(ctx)

	rows, err := tx.QueryContext(
		ctx,
		`SELECT chat_id, user_id, role_id, joined_at
		FROM chat_participant WHERE chat_id = $1`,
		chatID,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []chatparticipant.ChatParticipant{}, nil
		}

		return nil, &chatparticipant.ChatParticipantError{
			Code:    "DATABASE_ERROR",
			Message: "failed to get chat participants",
			Err:     err,
		}
	}

	defer rows.Close()

	foundChatParticipants := make([]chatparticipant.ChatParticipant, 0)

	for rows.Next() {
		var chatID uuid.UUID
		var userID uuid.UUID
		var roleID uuid.UUID
		var joinedAt time.Time

		if err := rows.Scan(&chatID, &userID, &roleID, &joinedAt); err != nil {
			return nil, &chatparticipant.ChatParticipantError{
				Code:    "DATABASE_ERROR",
				Message: "failed to scan chat participant",
				Err:     err,
			}
		}

		foundChatParticipants = append(foundChatParticipants, chatparticipant.ChatParticipantFromDB(chatID, userID, roleID, joinedAt))
	}

	return foundChatParticipants, nil
}

func (pr *PostgresChatParticipantRepo) GetAllChatsByUserID(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error) {
	tx := pr.GetTransaction(ctx)

	rows, err := tx.QueryContext(
		ctx,
		`SELECT chat_id
		FROM chat_participant
		WHERE user_id = $1`,
		userID,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []uuid.UUID{}, chatparticipant.ErrChatParticipantChatsByUserNotFound
		}
		return []uuid.UUID{}, &chatparticipant.ChatParticipantError{
			Code: "DATABASE_ERROR",
			Message: "failed to get all chats ids by user id",
			Err: err,
		}
	}

	defer rows.Close()

	foundChatIDs := make([]uuid.UUID, 0)

	for rows.Next() {
		var chatID uuid.UUID

		if err := rows.Scan(&chatID); err != nil {
			return nil, &chatparticipant.ChatParticipantError{
				Code:    "DATABASE_ERROR",
				Message: "failed to scan chat chatIDs",
				Err:     err,
			}
		}
		foundChatIDs = append(foundChatIDs, chatID)
	}

	return foundChatIDs, nil
}

func (pr *PostgresChatParticipantRepo) AddChatParticipant(ctx context.Context, chatParticipant chatparticipant.ChatParticipant) error {
	tx := pr.GetTransaction(ctx)

	result, err := tx.ExecContext(
		ctx,
		`INSERT INTO chat_participant (chat_id, user_id, role_id, joined_at)
		VALUES ($1, $2, $3, $4)`,
		chatParticipant.GetChatID(),
		chatParticipant.GetUserID(),
		chatParticipant.GetRoleID(),
		chatParticipant.GetJoinedAt(),
	)

	if err != nil {
		return &chatparticipant.ChatParticipantError{
			Code: "DATABASE_ERROR",
			Message: "failed to add chat participant",
			Err:     err,
		}
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return &chatparticipant.ChatParticipantError{
			Code: "DATABASE_ERROR",
			Message: "failed to get affected rows after added chat participant",
			Err:     err,
		}
	}

	if rowsAffected == 0 {
		return &chatparticipant.ChatParticipantError{
			Code: "UNEXPECTED_ERROR",
			Message: "chat participant with that id already exists",
			Err:     err,
		}
	}

	return nil
}

func (pr *PostgresChatParticipantRepo) UpdateChatParticipantRole(ctx context.Context, chatID uuid.UUID, userID uuid.UUID, newRoleID uuid.UUID) error {
	tx := pr.GetTransaction(ctx)

	result, err := tx.ExecContext(
		ctx,
		`UPDATE chat_participant SET role_id = $1
		WHERE chat_id = $2 AND user_id = $3`,
		newRoleID,
		chatID,
		userID,
	)

	if err != nil {
		return &chatparticipant.ChatParticipantError{
			Code: "DATABASE_ERROR",
			Message: "failed to update chat participant role",
			Err:     err,
		}
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return &chatparticipant.ChatParticipantError{
			Code: "DATABASE_ERROR",
			Message: "failed to get affected rows after updated chat participant role",
			Err:     err,
		}
	}

	if rowsAffected == 0 {
		return chatparticipant.ErrChatParticipantNotFound
	}

	return nil
}

func (pr *PostgresChatParticipantRepo) DeleteChatParticipant(ctx context.Context, chatID uuid.UUID, userID uuid.UUID) error {
	tx := pr.GetTransaction(ctx)

	result, err := tx.ExecContext(
		ctx,
		`DELETE FROM chat_participant
		WHERE chat_id = $1 AND user_id = $2`,
		chatID,
		userID,
	)

	if err != nil {
		return &chatparticipant.ChatParticipantError{
			Code: "DATABASE_ERROR",
			Message: "failed to delete chat participant",
			Err:     err,
		}
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return &chatparticipant.ChatParticipantError{
			Code: "DATABASE_ERROR",
			Message: "failed to get affected rows after deleted chat participant",
			Err:     err,
		}
	}

	if rowsAffected == 0 {
		return chatparticipant.ErrChatParticipantNotFound
	}

	return nil
}

func (pr *PostgresChatParticipantRepo) DeleteAllChatParticipants(ctx context.Context, chatID uuid.UUID) error {
	tx := pr.GetTransaction(ctx)

	_, err := tx.ExecContext(
		ctx,
		`DELETE FROM chat_participant
		WHERE chat_id = $1`,
		chatID,
	)

	if err != nil {
		return &chatparticipant.ChatParticipantError{
			Code: "DATABASE_ERROR",
			Message: "failed to delete all chat participants",
			Err:     err,
		}
	}


	return nil
}

func NewPostgresChatParticipantRepo(db *sql.DB) *PostgresChatParticipantRepo {
	return &PostgresChatParticipantRepo{
		db: db,
	}
}

func (pr *PostgresChatParticipantRepo) GetTransaction(ctx context.Context) transaction.DBTX {
	if tx := transaction.IsTransaction(ctx); tx != nil {
		return tx
	}
	return pr.db
}