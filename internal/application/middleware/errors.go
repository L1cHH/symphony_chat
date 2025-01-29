package middleware

type AuthMiddlewareErr struct {
	Code string
	Message string
	Err error
}

func (e *AuthMiddlewareErr) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	} else {
		return e.Message
	}
}

var (
	ErrAccessTokenWasNotProvided = &AuthMiddlewareErr {
		Code: "ACCESS_TOKEN_WAS_NOT_PROVIDED",
		Message: "access token was not provided in authorization header",
	}

	ErrInvalidAccessTokenFormat = &AuthMiddlewareErr {
		Code: "INVALID_ACCESS_TOKEN_FORMAT",
		Message: "invalid access token format",
	}

	ErrRefreshTokenWasNotSetInCookie = &AuthMiddlewareErr {
		Code: "REFRESH_TOKEN_WAS_NOT_SET_IN_COOKIE",
		Message: "refresh token was not set in cookie",
	}

	ErrRefreshTokenInCookieWasExpired = &AuthMiddlewareErr {
		Code: "REFRESH_TOKEN_IN_COOKIE_WAS_EXPIRED",
		Message: "refresh token in cookie was expired",
	}

	ErrCreatingAccessToken = &AuthMiddlewareErr {
		Code: "ACCESS_TOKEN_CREATION_FAILED",
		Message: "access token cant be created",
	}
)