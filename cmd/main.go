package main

import (
	"log"
	"os"
	"strconv"
	config "symphony_chat/internal/infrastructure/configs"
	"symphony_chat/internal/infrastructure/database"
	transaction"symphony_chat/internal/infrastructure/transaction/postgres"

	jwtPostgresRepo "symphony_chat/internal/infrastructure/jwt/postgres"
	authUserPostgresRepo "symphony_chat/internal/infrastructure/users/postgres"

	authentication "symphony_chat/internal/service/auth/authentication"
	registration "symphony_chat/internal/service/auth/registration"
	jwtService "symphony_chat/internal/service/jwt"

	authHandlerHTTP "symphony_chat/internal/application/auth/http"

	middleware "symphony_chat/internal/application/middleware"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	// Loading .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Creating configs

	// Database config
	postgresConfig := database.PostgresConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	}

	// JWTtoken config
	accessTTLinMinutes, err := strconv.ParseUint(os.Getenv("ACCESS_TTL_IN_MINUTES"), 10, 32)
	if err != nil {
		log.Fatal("Failed to parse JWT_ACCESS_TTL_IN_MINUTES:", err)
	}

	refreshTTLinDays, err := strconv.ParseUint(os.Getenv("REFRESH_TTL_IN_DAYS"), 10, 32)
	if err != nil {
		log.Fatal("Failed to parse JWT_REFRESH_TTL_IN_DAYS:", err)
	}

	jwtConfig := config.NewJWTConfig(
		os.Getenv("JWT_SECRET"),
		uint(accessTTLinMinutes),
		uint(refreshTTLinDays),
	)

	// Creating database connection
	db, err := database.NewPostgresConnection(postgresConfig)
	if err != nil {
		log.Fatal("Failed to create database connection:", err)
	}

	// Creating Repositories
	authUserRepo := authUserPostgresRepo.NewPostgresAuthUserRepo(db)
	jwtRepo := jwtPostgresRepo.NewPostgresJWTtokenRepo(db)

	// Creating services

	//Transaction manager
	transactionManager := transaction.NewPostgresTransactionManager(db)

	// JWTtoken service
	jwtService, err := jwtService.NewJWTtokenService(
		jwtService.WithJWTConfig(jwtConfig),
		jwtService.WithJWTtokenRepository(jwtRepo),
	)
	if err != nil {
		log.Fatal("Failed to create JWT service:", err)
	}

	// Registration service
	registrationService, err := registration.NewRegistrationService(
		registration.WithAuthUserRepository(authUserRepo),
		registration.WithJWTtokenService(jwtService),
		registration.WithTransactionManager(transactionManager),
	)
	if err != nil {
		log.Fatal("Failed to create registration service:", err)
	}

	// Authentication service
	authenticationService, err := authentication.NewAuthenticationService(
		authentication.WithAuthUserRepository(authUserRepo),
		authentication.WithJWTtokenService(jwtService),
		authentication.WithTransactionManager(transactionManager),
	)
	if err != nil {
		log.Fatal("Failed to create authentication service:", err)
	}

	// Creating handlers

	// Auth handler
	authHandler := authHandlerHTTP.NewAuthHandler(registrationService, authenticationService)

	// Создаем роутер
	r := gin.Default()

	// Базовый маршрут для проверки
	r.POST("/signup", authHandler.SignUp)
	r.POST("/login", authHandler.LogIn)
	r.POST("/logout", middleware.AuthMiddleware(jwtService), authHandler.LogOut)

	// Запускаем сервер
	log.Println("Starting server at :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
