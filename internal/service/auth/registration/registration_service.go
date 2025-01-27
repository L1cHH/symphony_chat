package registration

import (
	"context"
	"errors"
	publicDto "symphony_chat/internal/application/dto"
	tx "symphony_chat/internal/application/transaction"
	"symphony_chat/internal/domain/users"
	authdto "symphony_chat/internal/dto/auth"
	jwtService "symphony_chat/internal/service/jwt"
	utils "symphony_chat/utils/service"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var (
	//Correct format does not include numbers or punctuation symbols. Also there must be 8 characters in len
	ErrUncorrectFormatLogin = errors.New("uncorrect format login or password")
	//This login is already owned by someone else
	ErrLoginAlreadyExists = errors.New("user with this login already exists")
	//Problem with database
	ErrDatabaseProblem = errors.New("problem with database")
	//Problem with hashing password
	ErrHashingPassword = errors.New("error occurs while hashing")
	//Problem with creating JWT token
	ErrProblemWithJWT = errors.New("jwt error")
	//Uncoverable
	ErrUnimplementedError = errors.New("this error is uncovered")
)

type RegistrationService struct {
	authUserRepo users.AuthUserRepository
	jwtService   *jwtService.JWTtokenService
	transactionManager tx.TransactionManager
}

type RegistrationConfiguration func(*RegistrationService) error

func NewRegistrationService(configs ...RegistrationConfiguration) (*RegistrationService, error) {
	rs := &RegistrationService{}

	for _, cfgFunc := range configs {
		error := cfgFunc(rs)
		if error != nil {
			return nil, error
		}
	}

	return rs, nil
}

func WithAuthUserRepository(au users.AuthUserRepository) RegistrationConfiguration {

	return func(rs *RegistrationService) error {
		rs.authUserRepo = au
		return nil
	}
}

func WithJWTtokenService(jwtService *jwtService.JWTtokenService) RegistrationConfiguration {
	return func(rs *RegistrationService) error {
		rs.jwtService = jwtService
		return nil
	}
}

func WithTransactionManager(tm tx.TransactionManager) RegistrationConfiguration {
	return func(rs *RegistrationService) error {
		rs.transactionManager = tm
		return nil
	}
}

func (rs *RegistrationService) SignUpUser(ctx context.Context, userInput publicDto.LoginCredentials) (authdto.AuthTokens, error) {

	authTokens := authdto.AuthTokens{}

	err := rs.transactionManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		
		//Validation user input
		exists, err := rs.authUserRepo.IsUserExists(txCtx,userInput.Login)
		if err != nil {
			return &users.AuthError{
				Code: "REGISTRATION_ERROR",
				Message: "failed to check user with this login existense",
				Err: err,
			}
		}

		if exists {
			return users.ErrLoginAlreadyExists
		}

		//Hashing password
		hashedPassword, err := utils.HashPassword(userInput.Password)
		if err != nil {
			return &users.AuthError{
				Code: "PASSWORD_HASHING_ERROR",
				Message: "failed to hash password",
				Err: err,
			}
		}

		//Creating AuthUser
		authUser, err := rs.CreateAuthUser(txCtx, userInput.Login, hashedPassword)
		if err != nil {
			return &users.AuthError{
				Code: "REGISTRATION_ERROR",
				Message: "failed to create new auth user",
				Err: err,
			}
		}

		//Creating pair of jwt tokens(access and refresh)
		authTokens, err = rs.jwtService.GetCreatedPairTokens(txCtx, authUser.GetID())
		if err != nil {
			return &users.AuthError{
				Code: "REGISTRATION_ERROR",
				Message: "failed to create jwt tokens",
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

func (rs *RegistrationService) CreateAuthUser(ctx context.Context, login string, password string) (users.AuthUser, error) {
	authUser := users.NewAuthUser(uuid.New(), login, password, time.Now())
	//Adding AuthUser to database
	err := rs.authUserRepo.AddAuthUser(ctx, authUser)
	if err != nil {
		return users.AuthUser{}, err
	}
	return authUser, nil
}


func (rs *RegistrationService) SetRefreshTokenInHTTPCookie(c *gin.Context, refreshToken string) {
	c.SetCookie("refresh_token", refreshToken, int(rs.jwtService.GetRefreshTokenTTL()), "/", "localhost", false, true)
}

func (rs *RegistrationService) SetAccessTokenInHTTPCookie(c *gin.Context, accessToken string) {
	c.SetCookie("access_token", accessToken, int(rs.jwtService.GetAccessTokenTTL()), "/", "localhost", false, true)
}

