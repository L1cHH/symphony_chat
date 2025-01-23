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

type JwtRepository interface {
	AddJWTtoken(userID uuid.UUID, token JWTtoken) error
	GetJWTtoken(userID uuid.UUID) (JWTtoken, error)
	UpdateJWTtoken(userID uuid.UUID) error
	DeleteJWTtoken(userID uuid.UUID) error
}
