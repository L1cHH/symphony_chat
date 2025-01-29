package service

import (
	"context"
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
		return jwt.JWTtoken{}, &jwt.TokenError{
			Code: "ACCESS_TOKEN_NOT_CREATED",
			Message: "access token cant be created",
			Err: err,
		}
	}

	return accessToken, nil
}

// /Function for getting new refresh token
func (js *JWTtokenService) GetUpdatedRefreshToken(userID uuid.UUID) (jwt.JWTtoken, error) {
	refreshToken, err := jwt.NewJWT(userID, 0, js.jwtConfig.RefreshTTLinDays, []byte(js.jwtConfig.SecretKey))
	if err != nil {
		return jwt.JWTtoken{}, &jwt.TokenError{
			Code: "REFRESH_TOKEN_NOT_CREATED",
			Message: "refresh token cant be created",
			Err: err,
		}
	}

	return refreshToken, nil
}

// /Function that used when user again write login and password (when refresh token expires)
func (js *JWTtokenService) GetUpdatedPairTokens(txCtx context.Context, userID uuid.UUID) (authdto.AuthTokens, error) {
	accessToken, err := js.GetUpdatedAccessToken(userID)
	if err != nil {
		return authdto.AuthTokens{}, &jwt.TokenError{
			Code: "ACCESS_TOKEN_CREATION_FAILED",
			Message: "access token cant be created",
			Err: err,
		}
	}

	refreshToken, err := js.GetUpdatedRefreshToken(userID)
	if err != nil {
		return authdto.AuthTokens{}, &jwt.TokenError{
			Code: "REFRESH_TOKEN_CREATION_FAILED",
			Message: "refresh token cant be created",
			Err: err,
		}
	}

	err = js.jwtRepo.UpdateJWTtoken(txCtx, userID, refreshToken.GetToken())
	if err != nil {
		return authdto.AuthTokens{}, &jwt.TokenError{
			Code: "REFRESH_TOKEN_UPDATE_FAILED",
			Message: "refresh token cant be updated",
			Err: err,
		}
	}

	return authdto.AuthTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// /Function that used when user first time write login and password 
func (js *JWTtokenService) GetCreatedPairTokens(txCtx context.Context, userID uuid.UUID) (authdto.AuthTokens, error) {
	accessToken, err := js.GetUpdatedAccessToken(userID)
	if err != nil {
		return authdto.AuthTokens{}, &jwt.TokenError{
			Code: "ACCESS_TOKEN_CREATION_FAILED",
			Message: "access token cant be created",
			Err: err,
		}
	}

	refreshToken, err := js.GetUpdatedRefreshToken(userID)
	if err != nil {
		return authdto.AuthTokens{}, &jwt.TokenError{
			Code: "REFRESH_TOKEN_CREATION_FAILED",
			Message: "refresh token cant be created",
			Err: err,
		}
	}

	err = js.jwtRepo.AddJWTtoken(txCtx, refreshToken)
	if err != nil {
		return authdto.AuthTokens{}, &jwt.TokenError{
			Code: "REFRESH_TOKEN_CREATION_FAILED",
			Message: "refresh token cant be created",
			Err: err,
		}
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
			return nil, &jwt.TokenError{
				Code: "INVALID_TOKEN_SIGNING_METHOD",
				Message: "invalid token signing method",
				Err: errors.New("invalid token signing method"),
			}
		}
		return []byte(js.jwtConfig.SecretKey), nil
	})

	if err != nil {
		return uuid.Nil, &jwt.TokenError{
			Code: "TOKEN_PARSING_FAILED",
			Message: "token cant be parsed",
			Err: err,
		}
	}

	if !token.Valid {
		return uuid.Nil, jwt.ErrTokenExpired
	}

	claims, ok := token.Claims.(JWT.MapClaims)
	if !ok {
		return uuid.Nil, &jwt.TokenError{
			Code: "INVALID_TOKEN_CLAIMS_FORMAT",
			Message: "invalid token claims format",
			Err: errors.New("invalid token claims format"),
		}
	}

	userIDStr, ok := claims["sub"].(string)
	if !ok {
		return uuid.Nil, &jwt.TokenError{
			Code: "SUB_CLAIM_WAS_NOT_PROVIDED_IN_TOKEN_CLAIMS",
			Message: "sub claim was not provided in token claims",
			Err: errors.New("sub claim was not provided in token claims"),
		}
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return uuid.Nil, &jwt.TokenError{
			Code: "SUB_CLAIM_CANT_BE_PARSED_TO_UUID",
			Message: "sub claim cant be parsed to uuid",
			Err: err,
		}
	}

	return userID, nil
}

func (js *JWTtokenService) InvalidateRefreshToken(txCtx context.Context, userID uuid.UUID) error {
	err := js.jwtRepo.DeleteJWTtoken(txCtx, userID)
	if err != nil {
		return &jwt.TokenError{
			Code: "REFRESH_TOKEN_INVALDIATION_FAILED",
			Message: "refresh token cant be deleted",
			Err: err,
		}
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
