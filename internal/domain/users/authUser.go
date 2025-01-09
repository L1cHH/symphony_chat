package users

import (
	"time"
	"github.com/google/uuid"
)

type AuthUser struct {
	id uuid.UUID
	login string
	password string
	registrationAt time.Time
}

func (au AuthUser) GetID() uuid.UUID {
	return au.id
}

func (au AuthUser) GetLogin() string {
	return au.login
}

func (au AuthUser) GetPassword() string {
	return au.password
}

func (au AuthUser) GetRegistrationAt() time.Time {
	return au.registrationAt
}

//AuthUserRepository for managing AuthUser aggregate

type AuthUserRepository interface {

}