package publicdto

import (
	"symphony_chat/internal/domain/jwt"

	"github.com/google/uuid"
)

type JWTTokenDTO struct {
	AuthUserID uuid.UUID `json:"auth_user_id"`
	Token      string    `json:"token"`
}


func ToJWTTokenDTO(token jwt.JWTtoken) JWTTokenDTO {
	return JWTTokenDTO{
		AuthUserID: token.GetAuthUserID(),
		Token:      token.GetToken(),
	}
}

