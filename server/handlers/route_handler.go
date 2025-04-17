package handlers

import (
	"http"
	"net/http"
	"regexp"

	"chat-app/constants"

	"github.com/gin-gonic/gin"
)

func RenderHome() gin.HandlerFunc{
	return func(c *gin.Context){
		c.JSON(http.StatusOK, APIResponse{
			Code: http.StatusOK,
			Status: http.StatusText(http.StatusOK),
			Message: constants.APIWelcomeMessage,
			Response: nil,
		})
	}
}

func IsUsernameAvailable() gin.HandlerFunc{
	return func(c *gin.Context){
		username:= c.Param("username")
		isAlphaNumeric := regexp.MustCompile(`^[A-Za-z0-9]([A-Za-z0-9_-]*[A-Za-z0-9])?$`).MatchString

		if !isAlphaNumeric(username){
			c.JSON(http.StatusBadRequest, APIResponse{
				Code: http.StatusBadRequest,
				Status: http.StatusText(http.StatusBadRequest),
				Message: constants.APIWelcomeMessage,
				Response: nil,
			})
		}
	}
}

func Login() gin.HandlerFunc{
	return func(c *gin.Context){
		var userDetails UserDetails

		if err := c.ShouldBindJSON(&userDetails); err != nil{
			c.JSON(http.StatusBadRequest, APIResponse{
				Code:     http.StatusBadRequest,
				Status:   http.StatusText(http.StatusBadRequest),
				Message:  constants.UsernameAndPasswordCantBeEmpty,
				Response: nil,
			})
			return
		}

		// succesfil login
		c.JSON(http.StatusOK, APIResponse{
			Code: http.StatusOK,
			Status: http.StatusText(http.StatusOK),
			Message: constants.UserLoginCompleted,
			Response: userDetails,
		})
	}
}

func Registration() gin.HandlerFunc{
	return func(c *gin.Context){
		
	}
}

func UserSessionCheck() gin.HandlerFunc{
	return func(c *gin.Context){
		
	}
}

func GetMessagesHandler() gin.HandlerFunc{
	return func(c *gin.Context){
		
	}
}