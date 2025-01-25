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

		authID, err := js.ValidateToken(parts[1])
		if err == nil {
			ctx.Set("user_id", authID)
			ctx.Next()
			return
		}

		if err.Error() != "token is not valid" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": fmt.Errorf("validation error: %w", err).Error()})
			return 
		}

		cookieValue, err := ctx.Cookie("refresh_token")
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "refresh token is missing"})
			return 
		}

		authID, err = js.ValidateToken(cookieValue)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
                "error": "Authentication required",
                "code": "AUTH_REQUIRED",
            })
            return
		}

		newAccessToken, err := js.GetUpdatedAccessToken(authID)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
                "error": fmt.Errorf("problem with getting new access token: %w", err).Error(),
            })
			return 
		}

		ctx.Header("New-Access-Token", newAccessToken.GetToken())

		ctx.Set("user_id", authID)
		ctx.Next()
	}
}