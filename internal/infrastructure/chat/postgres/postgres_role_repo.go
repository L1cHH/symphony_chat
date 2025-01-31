package postgres

import (
	"context"
	"database/sql"
	"errors"
	"symphony_chat/internal/application/transaction"
	"symphony_chat/internal/domain/roles"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type PostgresChatRoleRepo struct {
	db *sql.DB
}

func NewPostgresChatRoleRepo(db *sql.DB) *PostgresChatRoleRepo {
	return &PostgresChatRoleRepo{
		db: db,
	}
}

func (pr *PostgresChatRoleRepo) GetChatRoleByID(ctx context.Context, role_id uuid.UUID) (roles.ChatRole, error) {
	
	tx := pr.GetTransaction(ctx)

	rows, err := tx.QueryContext(
		ctx,
		`SELECT id, name, permissions 
		FROM chat_role
		INNER JOIN chat_role_permission ON chat_role.id = chat_role_permission.role_id
		WHERE chat_role.id = $1`,
		role_id,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return roles.ChatRole{}, roles.ErrChatRoleNotFound
		}

		return roles.ChatRole{}, &roles.ChatRoleError{
			Code: "DATABASE_ERROR",
			Message: "failed to get chat role by id",
			Err: err,
		}
	}

	defer rows.Close()

	var id uuid.UUID
	var name string
	var permissions []roles.Permission

	for rows.Next() {
		var permission string

		if err := rows.Scan(&id, &name, &permission); err != nil {
			return roles.ChatRole{}, &roles.ChatRoleError {
				Code: "DATABASE_ERROR",
				Message: "failed to scan chat role",
				Err: err,
			}
		}

		permissions = append(permissions, roles.Permission(permission))
	}

	return roles.ChatRoleFromDB(id, name, permissions), nil
}

func (pr *PostgresChatRoleRepo) GetChatRoleByName(ctx context.Context, roleName string) (roles.ChatRole, error) {
	
	tx := pr.GetTransaction(ctx) 

	rows, err := tx.QueryContext(
		ctx,
		`SELECT id, name, permission
		FROM chat_role
		INNER JOIN chat_role_permission ON chat_role.id = chat_role_permission.role_id
		WHERE chat_role.name = $1`,
		roleName,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return roles.ChatRole{}, roles.ErrChatRoleNotFound
		}

		return roles.ChatRole{}, &roles.ChatRoleError{
			Code: "DATABASE_ERROR",
			Message: "failed to get chat role by name",
			Err: err,
		}
	}

	defer rows.Close()

	var id uuid.UUID
	var name string
	var permissions []roles.Permission

	for rows.Next() {
		var permission string

		if err := rows.Scan(&id, &name, &permission); err != nil {
			return roles.ChatRole{}, &roles.ChatRoleError {
				Code: "DATABASE_ERROR",
				Message: "failed to scan chat role",
				Err: err,
			}
		}

		permissions = append(permissions, roles.Permission(permission))
	}

	return roles.ChatRoleFromDB(id, name, permissions), nil

}

func (pr *PostgresChatRoleRepo) AddChatRole(ctx context.Context, chatRole roles.ChatRole) error {

	tx := pr.GetTransaction(ctx)

	permissions := make([]string, len(chatRole.GetPermissions()))

	for i, permission := range chatRole.GetPermissions() {
		permissions[i] = string(permission)
	}

	result, err := tx.ExecContext(
		ctx,
		`INSERT INTO chat_role (id, name) VALUES ($1, $2)
		INSERT INTO chat_role_permission (role_id, permission)
		SELECT $1, unnest($3::text[])`,
		chatRole.GetID(), chatRole.GetName(), pq.Array(permissions),
	)

	if err != nil {
		return &roles.ChatRoleError{
			Code: "DATABASE_ERROR",
			Message: "failed to add chat role",
			Err: err,
		}
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return &roles.ChatRoleError{
			Code: "DATABASE_ERROR",
			Message: "failed to get affected rows after added chat role",
			Err: err,
		}
	}

	//1 row from chat_role table + len(permissions) rows from chat_role_permission table
	expectedRowsAffected := int64(len(chatRole.GetPermissions()) + 1)

	if rowsAffected != expectedRowsAffected {
		return &roles.ChatRoleError{
			Code: "DATABASE_ERROR",
			Message: "failed to fully add chat role or chat role permissions",
		}
	}

	return nil
}

func (pr *PostgresChatRoleRepo) DeleteChatRoleByID(ctx context.Context, roleID uuid.UUID) error {
	tx := pr.GetTransaction(ctx)

	result, err := tx.ExecContext(
		ctx,
		`DELETE FROM chat_role_permission WHERE role_id = $1
		DELETE FROM chat_role WHERE id = $1`,
		roleID,
	)

	if err != nil {
		return &roles.ChatRoleError{
			Code: "DATABASE_ERROR",
			Message: "failed to delete chat role",
			Err: err,
		}
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return &roles.ChatRoleError{
			Code: "DATABASE_ERROR",
			Message: "failed to get affected rows after deleted chat role",
			Err: err,
		}
	}

	if rowsAffected == 0 {
		return roles.ErrChatRoleNotFound
	}

	return nil
}

//This method add new permisions to role without changing existing ones
func (pr *PostgresChatRoleRepo) UpdateChatRolePermissions(ctx context.Context, roleID uuid.UUID, newPermissions []roles.Permission) error {
	tx := pr.GetTransaction(ctx)

	permissions := make([]string, len(newPermissions))

	for i, p := range newPermissions {
		permissions[i] = string(p)
	}

	result, err := tx.ExecContext(
		ctx,
		`INSERT INTO chat_role_permission (role_id, permission)
		VALUES ($1, unnest($2::text[]))`,
		roleID, pq.Array(permissions),
	)

	if err != nil {
		return &roles.ChatRoleError{
			Code: "DATABASE_ERROR",
			Message: "failed to update chat role permissions",
			Err: err,
		}
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return &roles.ChatRoleError{
			Code: "DATABASE_ERROR",
			Message: "failed to get affected rows after updated chat role permissions",
			Err: err,
		}
	}

	expectedRowsAffected := int64(len(newPermissions))

	if rowsAffected != expectedRowsAffected {
		return &roles.ChatRoleError {
			Code: "DATABASE_ERROR",
			Message: "failed to fully update chat role permissions",
		}
	}

	return nil
}

func (pr *PostgresChatRoleRepo) GetTransaction(ctx context.Context) transaction.DBTX {
	if tx := transaction.IsTransaction(ctx); tx != nil {
		return tx
	}

	return pr.db
}