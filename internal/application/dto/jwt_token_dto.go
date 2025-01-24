package publicdto

import (
	"symphony_chat/internal/domain/jwt"

	"github.com/google/uuid"
)

type JWTRefreshTokenDTO struct {
	AuthUserID uuid.UUID `json:"auth_user_id"`
	Token      string    `json:"refresh_token"`
}

type JWTAccessTokenDTO struct {
	AuthUserID uuid.UUID `json:"auth_user_id"`
	Token      string    `json:"access_token"`
}

func ToJWTRefreshTokenDTO(token jwt.JWTtoken) JWTRefreshTokenDTO {
	return JWTRefreshTokenDTO{
		AuthUserID: token.GetAuthUserID(),
		Token:      token.GetToken(),
	}
}

func ToJWTAccessTokenDTO(token jwt.JWTtoken) JWTAccessTokenDTO {
	return JWTAccessTokenDTO{
		AuthUserID: token.GetAuthUserID(),
		Token:      token.GetToken(),
	}
}
