package service

import (
	"symphony_chat/internal/domain/jwt"

	"github.com/google/uuid"
)

type JWTtokenService struct {
	jwtRepo jwt.JwtRepository
}

type JWTtokenConfiguration func(*JWTtokenService) error

func NewJWTtokenService(configs ...JWTtokenConfiguration) (*JWTtokenService, error) {
	js := &JWTtokenService{}

	for _, cfg := range configs {
		err := cfg(js)
		if err != nil {
			return nil, err
		}
	}

	return js, nil
}

func WithJWTtokenRepository(jt jwt.JwtRepository) JWTtokenConfiguration {
	return func(js *JWTtokenService) error {
		js.jwtRepo = jt
		return nil
	}
}

///Function for getting new access token
func (js *JWTtokenService) GetUpdatedAccessToken(userID uuid.UUID) (jwt.JWTtoken, error) {
	accessToken, err := jwt.NewJWT(userID, 15, 0, []byte("secretKey"))
	if err != nil {
		return jwt.JWTtoken{}, err
	}

	return accessToken, nil
}

///Function for getting new refresh token
func (js *JWTtokenService) GetUpdatedRefreshToken(userID uuid.UUID) (jwt.JWTtoken, error) {
	refreshToken, err := jwt.NewJWT(userID, 0, 30, []byte("secretKey"))
	if err != nil {
		return jwt.JWTtoken{}, err
	}

	err = js.jwtRepo.AddJWTtoken(userID, refreshToken)
	if err != nil {
		return jwt.JWTtoken{}, err
	}
	return refreshToken, nil
}

///Function that used when user again write login and password (when refresh token expires)
///First is AccessToken, Second is RefreshToken
func (js *JWTtokenService) GetNewPairTokens(userID uuid.UUID) ([2]jwt.JWTtoken, error) {
	accessToken, err := js.GetUpdatedAccessToken(userID)
	if err != nil {
		return [2]jwt.JWTtoken{}, err
	}

	refreshToken, err := js.GetUpdatedRefreshToken(userID)
	if err != nil {
		return [2]jwt.JWTtoken{}, err
	}

	return [2]jwt.JWTtoken {accessToken, refreshToken}, nil
}