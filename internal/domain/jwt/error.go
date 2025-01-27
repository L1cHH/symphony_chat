package jwt


type TokenError struct {
	Code string
	Message string
	Err error
}

func (e *TokenError) Error() string {
	if e.Err != nil {
		return e.Message + ":" + e.Err.Error()
	}
	return e.Message
}

var (
	//Repo errors
	ErrTokenNotFound = &TokenError {
		Code: "TOKEN_NOT_FOUND",
		Message: "token not found in storage",
	}

	//Service errors
	ErrTokenExpired = &TokenError {
		Code: "TOKEN_EXPIRED",
		Message: "token is expired",
	}

	ErrTokenNotValid = &TokenError {
		Code: "TOKEN_NOT_VALID",
		Message: "token is expired",
	}
)