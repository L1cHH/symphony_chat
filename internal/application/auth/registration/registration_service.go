package registration

import (
	"errors"
	"symphony_chat/internal/domain/users"
	authdto "symphony_chat/internal/dto/auth"
	jwtService "symphony_chat/internal/service/jwt"
	utils "symphony_chat/utils/service"
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
	jwtService   jwtService.JWTtokenService
}

type RegistrationConfiguration func(*RegistrationService) error

func NewResgistretionService(configs ...RegistrationConfiguration) (*RegistrationService, error) {
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

func WithJWTtokenService(jwtService jwtService.JWTtokenService) RegistrationConfiguration {
	return func(rs *RegistrationService) error {
		rs.jwtService = jwtService
		return nil
	}
}

func (rs *RegistrationService) SignUpUser(userInput authdto.LoginCredentials) (authdto.AuthTokens, error) {

	//Validation user input

	if !utils.IsCorrectFormat(userInput.Login) || !utils.IsCorrectFormat(userInput.Password) {
		return authdto.AuthTokens{}, ErrUncorrectFormatLogin
	}

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
	authUser := rs.CreateAuthUser(userInput.Login, hashedPassword)

	//Adding AuthUser to database
	err = rs.authUserRepo.AddAuthUser(authUser)
	if err != nil {
		return authdto.AuthTokens{}, errors.New(ErrDatabaseProblem.Error() + ": " + err.Error())
	}

	//Creating pair of jwt tokens(access and refresh)
	jwtTokens, err := rs.jwtService.GetNewPairTokens(authUser.GetID())
	if err != nil {
		return authdto.AuthTokens{}, errors.New(ErrProblemWithJWT.Error() + ": " + err.Error())
	}

	return jwtTokens, nil
}

func (rs *RegistrationService) CreateAuthUser(login string, password string) users.AuthUser {
	authUser := users.NewAuthUser(login, password)
	return authUser
}
