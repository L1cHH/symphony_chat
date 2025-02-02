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

var (
	OwnerChatRole ChatRole = ChatRole {
		id: uuid.MustParse("11111111-1111-1111-1111-111111111111"),
		name: "OWNER",
		permissions: []Permission{
			PermissionDeleteChat,
			PermissionUpdateChatName,
			PermissionRemoveMember,
			PermissionManageRoles,
			PermissionAddMember,
			PermissionAddMessage,
			PermissionDeleteMessage,
			PermissionEditMessage,
		},
	}

	AdminChatRole ChatRole = ChatRole {
		id: uuid.MustParse("22222222-2222-2222-2222-222222222222"),
		name: "ADMIN",
		permissions: []Permission{
			PermissionAddMember,
			PermissionRemoveMember,
			PermissionUpdateChatName,
			PermissionAddMessage,
			PermissionDeleteMessage,
			PermissionEditMessage,
		},
	}

	MemberChatRole ChatRole = ChatRole {
		id: uuid.MustParse("33333333-3333-3333-3333-333333333333"),
		name: "MEMBER",
		permissions: []Permission {
			PermissionAddMember,
			PermissionUpdateChatName,
			PermissionAddMessage,
			PermissionDeleteMessage,
			PermissionEditMessage,
		},
	}
)

type Permission string 

const (
	PermissionAddMember Permission = "ADD_MEMBER_TO_CHAT"
	PermissionRemoveMember Permission = "REMOVE_MEMBER_FROM_CHAT"
	PermissionManageRoles Permission = "MANAGE_ROLES_OF_CHAT"
	PermissionDeleteChat Permission = "DELETE_CHAT"
	PermissionUpdateChatName Permission = "UPDATE_CHAT_NAME"
	PermissionAddMessage Permission = "ADD_MESSAGE_TO_CHAT"
	PermissionDeleteMessage Permission = "DELETE_MESSAGE_FROM_CHAT"
	PermissionEditMessage Permission = "EDIT_MESSAGE_IN_CHAT"
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
	GetChatRoles(ctx context.Context) ([]ChatRole, error)
}

