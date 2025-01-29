package http

import (
	"errors"
	"net/http"
	publicDto "symphony_chat/internal/application/dto"
	"symphony_chat/internal/domain/jwt"
	"symphony_chat/internal/domain/users"
	as "symphony_chat/internal/service/auth/authentication"
	rs "symphony_chat/internal/service/auth/registration"
	utils "symphony_chat/utils/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthHandler struct {
	registrationService   *rs.RegistrationService
	authenticationService *as.AuthenticationService
}

func NewAuthHandler(registrationService *rs.RegistrationService, authenticationService *as.AuthenticationService) *AuthHandler {
	return &AuthHandler{
		registrationService:   registrationService,
		authenticationService: authenticationService,
	}
}

func (ah *AuthHandler) SignUp(c *gin.Context) {
	var loginCredentials publicDto.LoginCredentials
	if err := c.ShouldBindJSON(&loginCredentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": "INVALID_INPUT",
			"message": "Invalid input format",
			"details": err.Error(),
		})
		return
	}

	//Validation user input

	if !utils.IsCorrectLoginFormat(loginCredentials.Login) || !utils.IsCorrectPasswordFormat(loginCredentials.Password) {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": "INVALID_LOGIN_OR_PASSWORD_FORMAT",
			"message": "login or password format is not valid",
			"details": "login must be between 6 and 30 characters and password must be between 10 and 25 characters in length",
		})
		return
	}

	tokens, err := ah.registrationService.SignUpUser(c.Request.Context(), loginCredentials)
	if err != nil {
		var authErr *users.AuthError
		var tokenErr *jwt.TokenError

		switch {
		case errors.As(err, &authErr):
			//Auth errors

			switch authErr.Code {
			case "CHECK_USER_EXISTENSE_ERROR", "CREATE_AUTH_USER_ERROR":
				c.JSON(http.StatusInternalServerError, gin.H{
					"code": "DATABASE_ERROR",
					"message": "internal server error, please try again later",
				})

			case "PASSWORD_HASHING_ERROR":
				c.JSON(http.StatusInternalServerError, gin.H{
					"code": "PASSWORD_HASHING_ERROR",
					"message": "internal server error, please try again later",
				})
			
			case "LOGIN_ALREADY_EXISTS":
				c.JSON(http.StatusConflict, gin.H{
					"code": "LOGIN_ALREADY_EXISTS",
					"message": "user with this login already exists",
				})

			default:
				//Unexpected error
				c.JSON(http.StatusInternalServerError, gin.H{
					"code": "INTERNAL_SERVER_ERROR",
					"message": "internal server error, please try again later",
				})
			}
		case errors.As(err, &tokenErr):
			//Token errors

			switch tokenErr.Code {

			case "CREATE_JWT_TOKENS_ERROR":
				c.JSON(http.StatusInternalServerError, gin.H{
					"code": "TOKEN_GENERATION_FAILED",
					"message": "failed to generate tokens, please try again later",
				})

			default:
				//Unexpected error
				c.JSON(http.StatusInternalServerError, gin.H{
					"code": "INTERNAL_SERVER_ERROR",
					"message": "internal server error, please try again later",
				})
			}
		default:
			//Unexpected error
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": "INTERNAL_SERVER_ERROR",
				"message": "internal server error, please try again later",
			})
		}

		return
	}

	//Setting only refresh token in cookies
	ah.registrationService.SetRefreshTokenInHTTPCookie(c, tokens.RefreshToken.GetToken())

	c.JSON(http.StatusOK, gin.H{
		"access_token":  publicDto.ToJWTTokenDTO(tokens.AccessToken),
		"refresh_token": publicDto.ToJWTTokenDTO(tokens.RefreshToken),
	})

}

func (ah *AuthHandler) LogIn(c *gin.Context) {
	var loginCredentials publicDto.LoginCredentials
	if err := c.ShouldBindJSON(&loginCredentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": "INVALID_INPUT",
			"message": "Invalid input format",
			"details": err.Error(),
		})
		return
	}

	tokens, err := ah.authenticationService.LogIn(c.Request.Context(),loginCredentials)
	if err != nil {
		var authErr *users.AuthError
		var tokenErr *jwt.TokenError

		switch {
		case errors.As(err, &authErr):
			switch authErr.Code {
			case "AUTH_USER_NOT_FOUND":
				c.JSON(http.StatusNotFound, gin.H {
					"code": "AUTH_USER_WITH_THIS_LOGIN_NOT_FOUND",
					"message": "auth user with this login not found",
				})

			case "WRONG_PASSWORD":
				c.JSON(http.StatusUnauthorized, gin.H {
					"code": "WRONG_PASSWORD",
					"message": "wrong password for this user",
				})

			default:
				c.JSON(http.StatusInternalServerError, gin.H {
					"code": "DATABASE_ERROR",
					"message": "internal server error, please try again later",
				})
			}
		
		case errors.As(err, &tokenErr):
			switch tokenErr.Code {
			case "CREATE_JWT_TOKENS_ERROR":
				c.JSON(http.StatusInternalServerError, gin.H {
					"code": "TOKEN_GENERATION_FAILED",
					"message": "failed to generate tokens, please try again later",
				})

			default:
				c.JSON(http.StatusInternalServerError, gin.H {
					"code": "INTERNAL_SERVER_ERROR",
					"message": "internal server error, please try again later",
				})
			}
		default:
			//Unexpected error
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": "INTERNAL_SERVER_ERROR",
				"message": "internal server error, please try again later",
			})
		}
		return
	}

	//Updating only refresh token in cookies
	ah.authenticationService.UpdateRefreshTokenInHTTPCookie(c, tokens.RefreshToken.GetToken())
	
	c.JSON(http.StatusOK, gin.H{
		"access_token":  publicDto.ToJWTTokenDTO(tokens.AccessToken),
		"refresh_token": publicDto.ToJWTTokenDTO(tokens.RefreshToken),
	})
}

func (ah *AuthHandler) LogOut(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"logout_error": "User id was not provided"})
		return
	}

	err := ah.authenticationService.LogOut(c.Request.Context(), userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"logout_error": err.Error()})
		return
	}

	//Clearing refresh token cookie
	ah.authenticationService.ClearRefreshTokenCookie(c)

	c.JSON(200, gin.H{
		"message": "Logout successful",
		"action":  "clear_tokens",
	})
}
