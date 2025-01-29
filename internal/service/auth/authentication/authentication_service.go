package authentication

import (
	"context"
	publicDto "symphony_chat/internal/application/dto"
	tx "symphony_chat/internal/application/transaction"
	"symphony_chat/internal/domain/jwt"
	"symphony_chat/internal/domain/users"
	authdto "symphony_chat/internal/dto/auth"
	jwtService "symphony_chat/internal/service/jwt"
	utils "symphony_chat/utils/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)


type AuthenticationService struct {
	jwtService *jwtService.JWTtokenService
	userRepo   users.AuthUserRepository
	transactionManager tx.TransactionManager
}

type AuthenticationConfiguration func(*AuthenticationService) error

func WithJWTtokenService(jwtService *jwtService.JWTtokenService) AuthenticationConfiguration {
	return func(as *AuthenticationService) error {
		as.jwtService = jwtService
		return nil
	}
}

func WithAuthUserRepository(au users.AuthUserRepository) AuthenticationConfiguration {
	return func(as *AuthenticationService) error {
		as.userRepo = au
		return nil
	}
}

func WithTransactionManager(tm tx.TransactionManager) AuthenticationConfiguration {
	return func(as *AuthenticationService) error {
		as.transactionManager = tm
		return nil
	}
}


func NewAuthenticationService(configs ...AuthenticationConfiguration) (*AuthenticationService, error) {
	as := &AuthenticationService{}

	for _, cfg := range configs {
		err := cfg(as)
		if err != nil {
			return nil, err
		}
	}

	return as, nil
}

func (as *AuthenticationService) LogIn(ctx context.Context,userInput publicDto.LoginCredentials) (authdto.AuthTokens, error) {

	authTokens := authdto.AuthTokens{}

	err := as.transactionManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		authUser, err := as.userRepo.GetAuthUserByLogin(txCtx,userInput.Login)
		if err != nil {
			if err == users.ErrAuthUserNotFound {
				return err
			}
			return &users.AuthError{
				Code: "GET_AUTH_USER_ERROR",
				Message: "failed to get auth_user by login from storage",
				Err: err,
			}
		}

		if !utils.CheckPassword(userInput.Password, authUser.GetPassword()) {
			return users.ErrWrongPassword
		}

		authTokens, err = as.jwtService.GetUpdatedPairTokens(txCtx, authUser.GetID())
		if err != nil {
			return &jwt.TokenError{
				Code: "CREATE_JWT_TOKENS_ERROR",
				Message: "failed to generate new jwt tokens",
				Err: err,
			}
		}

		return nil
	})

	if err != nil {
		return authdto.AuthTokens{}, err
	}

	return authTokens, nil
}

func (as *AuthenticationService) LogOut(ctx context.Context, userID uuid.UUID) error {

	err := as.transactionManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		err := as.jwtService.InvalidateRefreshToken(txCtx, userID)
		if err != nil {
			return &users.AuthError{
				Code: "LOGOUT_ERROR",
				Message: "failed to invalidate refresh token",
				Err: err,
			}
		}

		return nil
	})

	return err
}

func (as *AuthenticationService) UpdateRefreshTokenInHTTPCookie(c *gin.Context, refreshToken string) {
	c.SetCookie("refresh_token", refreshToken, int(as.jwtService.GetRefreshTokenTTL()), "/", "localhost", false, true)
}

func (as *AuthenticationService) UpdateAccessTokenInHTTPCookie(c *gin.Context, accessToken string) {
	c.SetCookie("access_token", accessToken, int(as.jwtService.GetAccessTokenTTL()), "/", "localhost", false, true)
}

func (as *AuthenticationService) ClearRefreshTokenCookie(c *gin.Context) {
	c.SetCookie("refresh_token", "", 0, "/", "localhost", false, true)
}
