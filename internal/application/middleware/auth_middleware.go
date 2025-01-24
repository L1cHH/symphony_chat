package middleware

import (
	"fmt"
	"net/http"
	"strings"
	jwtService "symphony_chat/internal/service/jwt"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(js *jwtService.JWTtokenService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString := ctx.GetHeader("Authorization")
		if tokenString == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
			return 
		}

		parts := strings.Split(tokenString, " ") 
		if len(parts) != 2 || parts[0] != "Bearer" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			return 
		}

		authID, err := js.ValidateAccessToken(parts[1])
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": fmt.Errorf("invalid token: %w", err)})
			return
		}

		ctx.Set("user_id", authID)
		ctx.Next()
	}
}