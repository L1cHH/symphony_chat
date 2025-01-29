package auth_test

import (
	"context"
	"os"
	"strconv"
	publicdto "symphony_chat/internal/application/dto"
	"symphony_chat/internal/domain/users"
	config "symphony_chat/internal/infrastructure/configs"
	jwtRepo "symphony_chat/internal/infrastructure/jwt/postgres"
	tx "symphony_chat/internal/infrastructure/transaction/postgres"
	"symphony_chat/internal/infrastructure/users/postgres"
	"symphony_chat/internal/service/auth/registration"
	jwtService "symphony_chat/internal/service/jwt"
	"symphony_chat/tests/integration/setup"
	"testing"

	_"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSignUp(t *testing.T) {
	db, err := setup.NewTestDB()
	if err != nil {
		t.Fatalf("failed to create test database connection: %v", err)
	}
	defer db.Close()

	testCases := []struct {
		name         string
		credentials  publicdto.LoginCredentials
		wantErr      bool
		errMessage   string
	}{ 
		{
		name: "Success registration",
		credentials: publicdto.LoginCredentials{
			Login: "Andrei.Karpukh2000@gmail.com",
			Password: "fhigbgiwgwwhnwihwgwb",
		},
		wantErr: false,
	},
	{
		name: "Duplicate login",
		credentials: publicdto.LoginCredentials{
			Login: "Kolomin.Andrey@gmail.com",
			Password: "kolomin.andrey2005",
		},
		wantErr: true,
		errMessage: users.ErrLoginAlreadyExists.Error(),
	},

	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err = db.TruncateAllTables()
			require.NoError(t, err)

			// JWTtoken config
			accessTTLinMinutes, err := strconv.ParseUint(os.Getenv("ACCESS_TTL_IN_MINUTES"), 10, 32)
			if err != nil {
				t.Fatal("Failed to parse JWT_ACCESS_TTL_IN_MINUTES:", err)
			}

			refreshTTLinDays, err := strconv.ParseUint(os.Getenv("REFRESH_TTL_IN_DAYS"), 10, 32)
			if err != nil {
				t.Fatal("Failed to parse JWT_REFRESH_TTL_IN_DAYS:", err)
			}

			jwtConfig := config.NewJWTConfig(
				os.Getenv("JWT_SECRET"),
				uint(accessTTLinMinutes),
				uint(refreshTTLinDays),
			)

			jwtService, err := jwtService.NewJWTtokenService(
				jwtService.WithJWTtokenRepository(jwtRepo.NewPostgresJWTtokenRepo(db.DB)),
				jwtService.WithJWTConfig(jwtConfig),
			)

			require.NoError(t, err)

			registrationService, err := registration.NewRegistrationService(
				registration.WithAuthUserRepository(postgres.NewPostgresAuthUserRepo(db.DB)),
				registration.WithTransactionManager(tx.NewPostgresTransactionManager(db.DB)),
				registration.WithJWTtokenService(jwtService),
			)

			require.NoError(t, err)

			tokens, err := registrationService.SignUpUser(context.Background(), tc.credentials)

			if tc.wantErr && tc.errMessage == users.ErrLoginAlreadyExists.Error() {
				tokens, err := registrationService.SignUpUser(context.Background(), tc.credentials)
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.errMessage)
				require.Empty(t, tokens)
				return
			}

			require.NoError(t, err)
			require.NotEmpty(t, tokens.AccessToken)
			require.NotEmpty(t, tokens.RefreshToken)
		})
	}
}