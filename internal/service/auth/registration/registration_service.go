package registration

import (
	"errors"
	publicDto "symphony_chat/internal/application/dto"
	"symphony_chat/internal/domain/users"
	authdto "symphony_chat/internal/dto/auth"
	jwtService "symphony_chat/internal/service/jwt"
	utils "symphony_chat/utils/service"
	"time"

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

func (rs *RegistrationService) SignUpUser(userInput publicDto.LoginCredentials) (authdto.AuthTokens, error) {

	//Validation user input
	exists, err := rs.authUserRepo.IsUserExists(userInput.Login)
	if err != nil {
		return authdto.AuthTokens{}, errors.New(ErrDatabaseProblem.Error() + ": " + err.Error())
	}

	if exists {
		return authdto.AuthTokens{}, ErrLoginAlreadyExists
	}

	//Hashing password
	hashedPassword, err := utils.HashPassword(userInput.Password)
	if err != nil {
		return authdto.AuthTokens{}, errors.New(ErrHashingPassword.Error() + ": " + err.Error())
	}

	//Creating AuthUser
	authUser, err := rs.CreateAuthUser(userInput.Login, hashedPassword)
	if err != nil {
		return authdto.AuthTokens{}, err
	}

	//Creating pair of jwt tokens(access and refresh)
	jwtTokens, err := rs.jwtService.GetCreatedPairTokens(authUser.GetID())
	if err != nil {
		return authdto.AuthTokens{}, errors.New(ErrProblemWithJWT.Error() + ": " + err.Error())
	}

	return jwtTokens, nil
}

func (rs *RegistrationService) CreateAuthUser(login string, password string) (users.AuthUser, error) {
	authUser := users.NewAuthUser(uuid.New(), login, password, time.Now())
	//Adding AuthUser to database
	err := rs.authUserRepo.AddAuthUser(authUser)
	if err != nil {
		return users.AuthUser{}, errors.New(ErrDatabaseProblem.Error() + ": " + err.Error())
	}
	return authUser, nil
}
