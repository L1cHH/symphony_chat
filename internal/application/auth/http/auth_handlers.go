package http

import (
	"net/http"
	auth "symphony_chat/internal/dto/auth"
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

func (ah *AuthHandler) SignUp(c *gin.Context) {
	var loginCredentials auth.LoginCredentials
	if err := c.ShouldBindJSON(&loginCredentials); err != nil {
		c.JSON(400, gin.H{"registration_error": err.Error()})
		return
	}

	//Validation user input

	if !utils.IsCorrectFormat(loginCredentials.Login) || !utils.IsCorrectFormat(loginCredentials.Password) {
		c.JSON(400, gin.H{"registration_error": "Uncorrect format login or password"})
		return
	}

	tokens, err := ah.registrationService.SignUpUser(loginCredentials)
	if err != nil {
		c.JSON(400, gin.H{"registration_error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"access_token":  tokens.AccessToken,
		"refresh_token": tokens.RefreshToken,
	})

}

func (ah *AuthHandler) LogIn(c *gin.Context) {
	var loginCredentials auth.LoginCredentials
	if err := c.ShouldBindJSON(&loginCredentials); err != nil {
		c.JSON(400, gin.H{"login_error": err.Error()})
		return
	}

	tokens, err := ah.authenticationService.LogIn(loginCredentials)
	if err != nil {
		c.JSON(400, gin.H{"login_error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"access_token":  tokens.AccessToken,
		"refresh_token": tokens.RefreshToken,
	})
}

func (ah *AuthHandler) LogOut(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"logout_error": "User id was not provided"})
		return
	}

	err := ah.authenticationService.LogOut(userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"logout_error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message": "Logout successful",
		"action":  "clear_tokens",
	})
}
