package service

import "symphony_chat/internal/domain/users"

type RegistrationService struct {
	authUserRepo users.AuthUserRepository
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