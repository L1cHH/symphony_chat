package http

import (
	"os"
	"strconv"
	authHandlers "symphony_chat/internal/application/auth/http"
	"symphony_chat/internal/application/middleware"
	config "symphony_chat/internal/infrastructure/configs"
	jwtRepo "symphony_chat/internal/infrastructure/jwt/postgres"
	tx "symphony_chat/internal/infrastructure/transaction/postgres"
	"symphony_chat/internal/infrastructure/users/postgres"
	"symphony_chat/internal/service/auth/authentication"
	"symphony_chat/internal/service/auth/registration"
	jwtService "symphony_chat/internal/service/jwt"
	"symphony_chat/tests/integration/setup"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func SetupRouter(t *testing.T, db *setup.TestDB) *gin.Engine {
	// Отключаем логи gin в тестах
    gin.SetMode(gin.TestMode)

    // JWT конфигурация
    accessTTL, err := strconv.ParseUint(os.Getenv("ACCESS_TTL_IN_MINUTES"), 10, 32)
    require.NoError(t, err)
    refreshTTL, err := strconv.ParseUint(os.Getenv("REFRESH_TTL_IN_DAYS"), 10, 32)
    require.NoError(t, err)

    jwtConfig := config.NewJWTConfig(
        os.Getenv("JWT_SECRET"),
        uint(accessTTL),
        uint(refreshTTL),
    )

    // Создаем сервисы
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

    authenticationService, err := authentication.NewAuthenticationService(
        authentication.WithAuthUserRepository(postgres.NewPostgresAuthUserRepo(db.DB)),
        authentication.WithTransactionManager(tx.NewPostgresTransactionManager(db.DB)),
        authentication.WithJWTtokenService(jwtService),
    )
    require.NoError(t, err)

    // Создаем роутер
    router := gin.New()
    authHandler := authHandlers.NewAuthHandler(registrationService, authenticationService)

    // Регистрируем маршруты
    router.POST("/auth/signup", authHandler.SignUp)
    router.POST("/auth/login", authHandler.LogIn)
    router.POST("/auth/logout", middleware.AuthMiddleware(jwtService), authHandler.LogOut)

    return router
}