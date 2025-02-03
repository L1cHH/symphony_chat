package roles

type ChatRoleError struct {
	Code    string
	Message string
	Err     error
}

func (e *ChatRoleError) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}

var (
	ErrWrongChatRoleName = &ChatRoleError {
		Code: "WRONG_CHAT_ROLE_NAME",
		Message: "wrong chat role name",
	}

	ErrChatRoleNotFound = &ChatRoleError {
		Code: "CHAT_ROLE_NOT_FOUND",
		Message: "chat role not found",
	}

	ErrInsufficientPermissions = &ChatRoleError {
		Code: "INSUFFICIENT_PERMISSIONS",
		Message: "user does not have required permissions for this action",
	}
)