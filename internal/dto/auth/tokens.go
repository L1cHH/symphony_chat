package authdto

import "symphony_chat/internal/domain/jwt"

type AuthTokens struct {
	AccessToken  jwt.JWTtoken `json:"accessToken"`
	RefreshToken jwt.JWTtoken `json:"refreshToken"`
}
