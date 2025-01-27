package users

type AuthError struct {
	Code string
	Message string
	Err error
}

func (e *AuthError) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}

var (
	ErrAuthUserNotFound = &AuthError {
		Code: "AUTH_USER_NOT_FOUND",
		Message: "auth user not found in storage",
	}

	ErrWrongPassword = &AuthError {
		Code: "WRONG_PASSWORD",
		Message: "wrong password for this user",
	}

	ErrLoginAlreadyExists = &AuthError {
		Code: "LOGIN_ALREADY_EXISTS",
		Message: "login already exists",
	}
)
