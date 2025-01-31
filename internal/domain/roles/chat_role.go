package roles

import (
	"context"

	"github.com/google/uuid"
)

type ChatRole struct {
	id		    uuid.UUID
	name	    string
	permissions []Permission
}

type Permission string 

const (
	PermissionAddMember Permission = "ADD_MEMBER_TO_CHAT"
	PermissionRemoveMember Permission = "REMOVE_MEMBER_FROM_CHAT"
	PermissionManageRoles Permission = "MANAGE_ROLES_OF_CHAT"
	PermissionDeleteChat Permission = "DELETE_CHAT"
	PermissionUpdateChatName Permission = "UPDATE_CHAT_NAME"
	PermissionAddMessage Permission = "ADD_MESSAGE_TO_CHAT"
)

func (c ChatRole) GetID() uuid.UUID {
	return c.id
}

func (c ChatRole) GetName() string {
	return c.name
}

func (c ChatRole) GetPermissions() []Permission {
	return c.permissions
}

func NewChatRole(id uuid.UUID, name string, permissions []Permission) (ChatRole, error) {
	
	if len(name) == 0 {
		return ChatRole{}, ErrWrongChatRoleName
	}
	
	return ChatRole{
		id: id,
		name: name,
		permissions: permissions,
	}, nil
}

func ChatRoleFromDB(id uuid.UUID, name string, permissions []Permission) ChatRole {
	return ChatRole{
		id: id,
		name: name,
		permissions: permissions,
	}
}


type ChatRoleRepository interface {
	GetChatRoleByID(ctx context.Context, id uuid.UUID) (ChatRole, error)
	GetChatRoleByName(ctx context.Context, name string) (ChatRole, error)
	AddChatRole(ctx context.Context, chatRole ChatRole) error
	DeleteChatRoleByID(ctx context.Context, id uuid.UUID) error
	UpdateChatRolePermissions(ctx context.Context, id uuid.UUID, newPermissions []Permission) error
}

