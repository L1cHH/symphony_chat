package postgres

import (
	"context"
	"database/sql"
	"errors"
	"symphony_chat/internal/application/transaction"
	"symphony_chat/internal/domain/roles"

	"github.com/google/uuid"
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

func (pr *PostgresChatRoleRepo) GetChatRoles(ctx context.Context) ([]roles.ChatRole, error) {
	
	tx := pr.GetTransaction(ctx)

	rows, err := tx.QueryContext(
		ctx,
		`SELECT id, name, permissions FROM chat_role
		INNER JOIN chat_role_permission ON chat_role.id = chat_role_permission.role_id`,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []roles.ChatRole{}, roles.ErrChatRoleNotFound
		}

		return []roles.ChatRole{}, &roles.ChatRoleError{
			Code: "DATABASE_ERROR",
			Message: "failed to get chat roles from storage",
			Err: err,
		}
	}

	defer rows.Close()

	foundRoles := make(map[uuid.UUID]struct {
		name string
		permissions []roles.Permission
	})

	for rows.Next() {
		var id uuid.UUID
		var name string
		var permission roles.Permission

		if err := rows.Scan(&id, &name, &permission); err != nil {
			return []roles.ChatRole{}, &roles.ChatRoleError{
				Code: "DATABASE_ERROR",
				Message: "failed to scan chat role",
				Err: err,
			}
		}

		role, exists := foundRoles[id]
		if !exists {
			role = struct{
				name string; 
				permissions []roles.Permission
			}{
				name: name,
				permissions: make([]roles.Permission, 0),
			}

			role.permissions = append(role.permissions, permission)

			foundRoles[id] = role

		} else {
			role.permissions = append(role.permissions, permission)
			foundRoles[id] = role
		}

	}

	chatRoles := make([]roles.ChatRole, 0, len(foundRoles))

	for id, role := range foundRoles {
		chatRoles = append(chatRoles, roles.ChatRoleFromDB(id, role.name, role.permissions))
	}

	return chatRoles, nil
}

func (pr *PostgresChatRoleRepo) GetTransaction(ctx context.Context) transaction.DBTX {
	if tx := transaction.IsTransaction(ctx); tx != nil {
		return tx
	}

	return pr.db
}