package middleware

import (
	"errors"
	"net/http"
	"strings"
	jwtService "symphony_chat/internal/service/jwt"
	jwt "symphony_chat/internal/domain/jwt"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(js *jwtService.JWTtokenService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString := ctx.GetHeader("Authorization")
		if tokenString == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code": ErrAccessTokenWasNotProvided.Code,
				"message": ErrAccessTokenWasNotProvided.Message,
			})
			return 
		}

		parts := strings.Split(tokenString, " ") 
		if len(parts) != 2 || parts[0] != "Bearer" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code": ErrInvalidAccessTokenFormat.Code,
				"message": ErrInvalidAccessTokenFormat.Message,
			})
			return 
		}

		authID, err := js.ValidateToken(parts[1])
		if err == nil {
			ctx.Set("user_id", authID)
			ctx.Next()
			return
		}

		if !errors.Is(err, jwt.ErrTokenNotValid) {
			statusCode, response := mapValidateError(err)
			ctx.AbortWithStatusJSON(statusCode, response)
			return 
		}

		refreshToken, err := ctx.Cookie("refresh_token")
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code": ErrRefreshTokenWasNotSetInCookie.Code,
				"message": ErrRefreshTokenWasNotSetInCookie.Message,
			})
			return 
		}

		authID, err = js.ValidateToken(refreshToken)
		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"code": ErrRefreshTokenInCookieWasExpired.Code,
					"message": ErrRefreshTokenInCookieWasExpired.Message,
				})
                return
			} else {
				httpStatus, response := mapValidateError(err)
				ctx.AbortWithStatusJSON(httpStatus, response)
				return
			}
		}

		newAccessToken, err := js.GetUpdatedAccessToken(authID)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
                "code": ErrCreatingAccessToken.Code,
				"message": ErrCreatingAccessToken.Message,
            })
			return 
		}

		ctx.Header("New-Access-Token", newAccessToken.GetToken())

		ctx.Set("user_id", authID)
		ctx.Next()
	}
}

func mapValidateError(err error) (int, gin.H) {
	var jwtError *jwt.TokenError
	errors.As(err, &jwtError)

	switch  jwtError.Code{
		case "INVALID_TOKEN_SIGNING_METHOD",
		    "TOKEN_PARSING_FAILED",
			"INVALID_TOKEN_CLAIMS_FORMAT",
			"SUB_CLAIM_WAS_NOT_PROVIDED_IN_TOKEN_CLAIMS",
			"SUB_CLAIM_CANT_BE_PARSED_TO_UUID":
			return http.StatusUnauthorized, gin.H{
				"code": "INVALID_TOKEN_FORMAT",
				"message": "invalid token format",
				"details": jwtError.Err.Error(),
			}
		case "TOKEN_EXPIRED":
			return http.StatusUnauthorized, gin.H{
				"code": "TOKEN_EXPIRED",
				"message": "token is expired",
			}
		default:
			return http.StatusInternalServerError, gin.H{
				"code": "INTERNAL_SERVER_ERROR",
				"message": "internal server error",
				"details": jwtError.Err.Error(),
			}
	}
}