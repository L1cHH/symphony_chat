package users

import "github.com/google/uuid"

type AuthUser struct {
	id uuid.UUID
	nickname string
	login string
	password string
}

func (au AuthUser) GetID() uuid.UUID {
	return au.id
}