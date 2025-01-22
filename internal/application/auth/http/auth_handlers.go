package http

import (
	auth "symphony_chat/internal/dto/auth"
	as "symphony_chat/internal/service/auth/authentication"
	rs "symphony_chat/internal/service/auth/registration"
	utils "symphony_chat/utils/service"

	"github.com/gin-gonic/gin"
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
		"auth_token":    tokens.AccessToken,
		"refresh_token": tokens.RefreshToken,
	})

}
