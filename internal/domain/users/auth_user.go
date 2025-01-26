package users

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type AuthUser struct {
	id             uuid.UUID
	login          string
	password       string
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

func NewAuthUser(id uuid.UUID, login string, password string, registrationAt time.Time) AuthUser {
	return AuthUser{
		id:             id,
		login:          login,
		password:       password,
		registrationAt: registrationAt,
	}
}

//AuthUserRepository for managing AuthUser aggregate
type AuthUserRepository interface {
	GetAuthUserById(ctx context.Context, id uuid.UUID) (AuthUser, error)
	GetAuthUserByLogin(ctx context.Context, login string) (AuthUser, error)
	IsUserExists(ctx context.Context, login string) (bool, error)
	AddAuthUser(ctx context.Context, authUser AuthUser) error
	UpdateLogin(ctx context.Context, id uuid.UUID, newLogin string) error
	UpdatePassword(ctx context.Context, id uuid.UUID, newPassword string) error
	DeleteAuthUser(ctx context.Context, id uuid.UUID) error
}
