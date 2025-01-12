package authentication

import (
	"errors"
	"symphony_chat/internal/domain/users"
	authdto "symphony_chat/internal/dto/auth"
	jwtService "symphony_chat/internal/service/jwt"
	utils "symphony_chat/utils/service"
)

var (
	//User with that login not found
	ErrUserNotFound = errors.New("user with that login not found")
	//Wrong password for this user
	ErrWrongPassword = errors.New("wrong password for this user")
	//Problem with updating jwt tokens
	ErrProblemWithJWT = errors.New("problem with updating jwt tokens")
)

type AuthenticationService struct {
	jwtService jwtService.JWTtokenService
	userRepo   users.AuthUserRepository
}

type AuthenticationConfiguration func(*AuthenticationService) error

func WithJWTtokenService(jwtService jwtService.JWTtokenService) AuthenticationConfiguration {
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

func (as *AuthenticationService) LogIn(userInput authdto.LoginCredentials) (authdto.AuthTokens, error) {
	authUser, err := as.userRepo.GetAuthUserByLogin(userInput.Login)
	if err != nil {
		return authdto.AuthTokens{}, errors.New(ErrUserNotFound.Error() + ": " + err.Error())
	}

	if !utils.CheckPassword(userInput.Password, authUser.GetPassword()) {
		return authdto.AuthTokens{}, ErrWrongPassword
	}

	tokens, err := as.jwtService.GetNewPairTokens(authUser.GetID())
	if err != nil {
		return authdto.AuthTokens{}, errors.New(ErrProblemWithJWT.Error() + ": " + err.Error())
	}

	return tokens, nil
}
