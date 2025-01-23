package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTtoken struct {
	auth_user_id uuid.UUID
	token        string
}

func (jt JWTtoken) GetAuthUserID() uuid.UUID {
	return jt.auth_user_id
}

func (jt JWTtoken) GetToken() string {
	return jt.token
}

// /Function for generating new JWT token
func NewJWT(userID uuid.UUID, minutesTTL uint, daysTTL uint, secretKey []byte) (JWTtoken, error) {
	token := jwt.NewWithClaims(&jwt.SigningMethodHMAC{},
		jwt.MapClaims{
			"sub": userID,
			"iat": time.Now(),
			"exp": time.Now().Add(time.Duration(minutesTTL)*time.Minute + time.Duration(daysTTL)*time.Hour*24),
		},
	)

	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		return JWTtoken{}, err
	}

	return JWTtoken{
		auth_user_id: userID,
		token:        signedToken,
	}, nil
}

// /Function that converts JWTtoken from database format to domain format
func FromDB(authUserID uuid.UUID, token string) JWTtoken {
	return JWTtoken{
		auth_user_id: authUserID,
		token:        token,
	}
}

type JwtRepository interface {
	AddJWTtoken(token JWTtoken) error
	GetJWTtoken(authUserID uuid.UUID) (JWTtoken, error)
	UpdateJWTtoken(authUserID uuid.UUID, newToken string) error
	DeleteJWTtoken(authUserID uuid.UUID) error
}
