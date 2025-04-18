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
		var requestPayload RegistrationRequest

		if err := c.ShouldBindJSON(&requestPayload); err != nil{
			c.JSON(http.StatusBadRequest, APIResponse{
				Code: http.StatusBadRequest,
				Status: http.StatusText(http.StatusBadRequest),
				Message: constants.ServerFailedResponse,
				Response: nil,
			})
			return
		}

		if requestPayload.Username == "" {
			c.JSON(http.StatusBadRequest, APIResponse{
				Code:     http.StatusBadRequest,
				Status:   http.StatusText(http.StatusBadRequest),
				Message:  constants.UsernameCantBeEmpty,
				Response: nil,
			})
			return
		}

		if requestPayload.Password == "" {
			c.JSON(http.StatusBadRequest, APIResponse{
				Code:     http.StatusBadRequest,
				Status:   http.StatusText(http.StatusBadRequest),
				Message:  constants.PasswordCantBeEmpty,
				Response: nil,
			})
			return
		}

		userObjectID, registrationErr := RegisterQueryHandler(requestPayload)
		if registrationErr != nil{
			c.JSON(http.StatusInternalServerError, APIResponse{
				Code:     http.StatusInternalServerError,
				Status:   http.StatusText(http.StatusInternalServerError),
				Message:  constants.ServerFailedResponse,
				Response: nil,
			})
			return
		}

		c.JSON(http.StatusOK, APIResponse{
			Code:     http.StatusOK,
			Status:   http.StatusText(http.StatusOK),
			Message:  constants.UserRegistrationCompleted,
			Response: UserResponse{
				Username: requestPayload.Username,
				UserID: userObjectID,
			},
		})
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