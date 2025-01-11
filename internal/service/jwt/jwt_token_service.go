package service

import (
	"errors"
	"symphony_chat/internal/domain/jwt"

	JWT "github.com/golang-jwt/jwt/v5"
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

func (js *JWTtokenService) ValidateAccessToken(tokenString string) (uuid.UUID, error) {
	token, err := JWT.Parse(tokenString, func(t *JWT.Token) (interface{}, error) {
		if _, ok := t.Method.(*JWT.SigningMethodHMAC); !ok {
			return nil, errors.New("wrong alg method in jwt token. So token cant be parsed")
		}
		return []byte("secretKey"), nil
	})

	if err != nil {
		return uuid.Nil, err
	}

	if !token.Valid {
		return uuid.Nil, errors.New("token is not valid")
	}

	claims, ok := token.Claims.(JWT.MapClaims)
	if !ok {
		return uuid.Nil, errors.New("invalid token claims format")
	}

	userIDStr, ok := claims["sub"].(string)
	if !ok {
		return uuid.Nil, errors.New("uncorrect format of sub claim, must be string")
	}

	userID, err := uuid.Parse(userIDStr)
    if err != nil {
        return uuid.Nil, errors.New("invalid user ID format in token, must be uuid")
    }

	return userID, nil
}