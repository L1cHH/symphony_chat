package service

import (
	"errors"
	"symphony_chat/internal/domain/jwt"
	authdto "symphony_chat/internal/dto/auth"
	config "symphony_chat/internal/infrastructure/configs"

	JWT "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTtokenService struct {
	jwtRepo   jwt.JwtRepository
	jwtConfig config.JWTConfig
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

func WithJWTConfig(jwtConfig config.JWTConfig) JWTtokenConfiguration {
	return func(js *JWTtokenService) error {
		js.jwtConfig = jwtConfig
		return nil
	}
}

func WithJWTtokenRepository(jt jwt.JwtRepository) JWTtokenConfiguration {
	return func(js *JWTtokenService) error {
		js.jwtRepo = jt
		return nil
	}
}

// /Function for getting new access token
func (js *JWTtokenService) GetUpdatedAccessToken(userID uuid.UUID) (jwt.JWTtoken, error) {
	accessToken, err := jwt.NewJWT(userID, js.jwtConfig.AccessTTLinMinutes, 0, []byte(js.jwtConfig.SecretKey))
	if err != nil {
		return jwt.JWTtoken{}, err
	}

	return accessToken, nil
}

// /Function for getting new refresh token
func (js *JWTtokenService) GetUpdatedRefreshToken(userID uuid.UUID) (jwt.JWTtoken, error) {
	refreshToken, err := jwt.NewJWT(userID, 0, js.jwtConfig.RefreshTTLinDays, []byte(js.jwtConfig.SecretKey))
	if err != nil {
		return jwt.JWTtoken{}, err
	}

	return refreshToken, nil
}

// /Function that used when user again write login and password (when refresh token expires)
func (js *JWTtokenService) GetUpdatedPairTokens(userID uuid.UUID) (authdto.AuthTokens, error) {
	accessToken, err := js.GetUpdatedAccessToken(userID)
	if err != nil {
		return authdto.AuthTokens{}, err
	}

	refreshToken, err := js.GetUpdatedRefreshToken(userID)
	if err != nil {
		return authdto.AuthTokens{}, err
	}

	err = js.jwtRepo.UpdateJWTtoken(userID, refreshToken.GetToken())
	if err != nil {
		return authdto.AuthTokens{}, err
	}

	return authdto.AuthTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// /Function that used when user first time write login and password 
func (js *JWTtokenService) GetCreatedPairTokens(userID uuid.UUID) (authdto.AuthTokens, error) {
	accessToken, err := js.GetUpdatedAccessToken(userID)
	if err != nil {
		return authdto.AuthTokens{}, err
	}

	refreshToken, err := js.GetUpdatedRefreshToken(userID)
	if err != nil {
		return authdto.AuthTokens{}, err
	}

	err = js.jwtRepo.AddJWTtoken(refreshToken)
	if err != nil {
		return authdto.AuthTokens{}, err
	}

	return authdto.AuthTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// /Functon that used for validating access token
func (js *JWTtokenService) ValidateToken(tokenString string) (uuid.UUID, error) {
	token, err := JWT.Parse(tokenString, func(t *JWT.Token) (interface{}, error) {
		if _, ok := t.Method.(*JWT.SigningMethodHMAC); !ok {
			return nil, errors.New("wrong alg method in jwt token. So token cant be parsed")
		}
		return []byte(js.jwtConfig.SecretKey), nil
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

func (js *JWTtokenService) InvalidateRefreshToken(userID uuid.UUID) error {
	err := js.jwtRepo.DeleteJWTtoken(userID)
	if err != nil {
		return errors.New("problem with deleting refresh token")
	}
	return nil
}

///Function for getting refresh token TTL in seconds
func (js *JWTtokenService) GetRefreshTokenTTL() uint {
	return js.jwtConfig.RefreshTTLinDays * 24 * 3600
}

///Function for getting access token TTL in seconds
func (js *JWTtokenService) GetAccessTokenTTL() uint {
	return js.jwtConfig.AccessTTLinMinutes * 60
}
