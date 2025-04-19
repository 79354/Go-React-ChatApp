package handlers

import (
	"log"
	"net/http"
	"regexp"
	"strconv"

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
		type usernameAvailable struct{
			IsUsernameAvailable bool `json:"isUsernameAvailable"`
		}
		
		username:= c.Param("username")
		isAlphaNumeric := regexp.MustCompile(`^[A-Za-z0-9]([A-Za-z0-9_-]*[A-Za-z0-9])?$`).MatchString

		if !isAlphaNumeric(username){
			c.JSON(http.StatusBadRequest, APIResponse{
				Code: http.StatusBadRequest,
				Status: http.StatusText(http.StatusBadRequest),
				Message: constants.UsernameCantBeEmpty,
				Response: nil,
			})
			return
		}

		isUsernameAvailable := IsUsernameAvailableQueryHandler(username)
		if isUsernameAvailable{
			c.JSON(http.StatusOK, APIResponse{
				Code: http.StatusOK,
				Status: http.StatusText(http.StatusOK),
				Message: constants.UsernameIsAvailable,
				Response: usernameAvailable{
					IsUsernameAvailable: true,
				},
			})
		}else{
			c.JSON(http.StatusOK, APIResponse{
				Code:    http.StatusOK,
				Status:  http.StatusText(http.StatusOK),
				Message: constants.UsernameIsNotAvailable,
				Response: usernameAvailable{
					IsUsernameAvailable: false,
				},
			})
		}
	}
}

func Login() gin.HandlerFunc{
	return func(c *gin.Context){
		var userDetails LoginRequest

		if err := c.ShouldBindJSON(&userDetails); err != nil{
			c.JSON(http.StatusBadRequest, APIResponse{
				Code:     http.StatusBadRequest,
				Status:   http.StatusText(http.StatusBadRequest),
				Message:  constants.UsernameAndPasswordCantBeEmpty,
				Response: nil,
			})
			return
		}

		if userDetails.Username == "" {
			c.JSON(http.StatusBadRequest, APIResponse{
				Code:     http.StatusBadRequest,
				Status:   http.StatusText(http.StatusBadRequest),
				Message:  constants.UsernameCantBeEmpty,
				Response: nil,
			})
		}

		if userDetails.Password == "" {
			c.JSON(http.StatusInternalServerError, APIResponse{
				Code:     http.StatusInternalServerError,
				Status:   http.StatusText(http.StatusInternalServerError),
				Message:  constants.PasswordCantBeEmpty,
				Response: nil,
			})
	
		}

		userDetailsResponse, loginErrorMessage :=  LoginQueryHandler(userDetails)

		if loginErrorMessage != nil {
			c.JSON(http.StatusNotFound, APIResponse{
				Code:     http.StatusNotFound,
				Status:   http.StatusText(http.StatusNotFound),
				Message:  loginErrorMessage.Error(),
				Response: nil,
			})
			return
		}

		// succesfil login
		c.JSON(http.StatusOK, APIResponse{
			Code: http.StatusOK,
			Status: http.StatusText(http.StatusOK),
			Message: constants.UserLoginCompleted,
			Response: userDetailsResponse,
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
		var IsAlphaNumeric = regexp.MustCompile(`^[A-Za-z0-9]([A-Za-z0-9_-]*[A-Za-z0-9])?$`).MatchString
		uid := c.Param("userID")

		if !IsAlphaNumeric(uid){
			c.JSON(http.StatusBadRequest, APIResponse{
				Code:     http.StatusBadRequest,
				Status:   http.StatusText(http.StatusBadRequest),
				Message:  constants.UsernameCantBeEmpty,
				Response: nil,
			})
			return
		}

		userDetails := GetUserByUserID(uid)
		if userDetails == (UserDetails{}){
			c.JSON(http.StatusOK, APIResponse{
				Code:     http.StatusOK,
				Status:   http.StatusText(http.StatusOK),
				Message:  constants.YouAreNotLoggedIN,
				Response: false,
			})
			return
		}

		c.JSON(http.StatusOK, APIResponse{
			Code: http.StatusOK,
			Status: http.StatusText(http.StatusOK),
			Message: constants.YouAreLoggedIN,
			Response: userDetails.Online == "Y",
		})
	}
}

func GetMessagesHandler() gin.HandlerFunc{
	return func(c *gin.Context){
		var IsAlphaNumeric = regexp.MustCompile(`^[A-Za-z0-9]([A-Za-z0-9_-]*[A-Za-z0-9])?$`).MatchString
		toUserID := c.Param("toUserID")
		fromUserID := c.Param("fromUserID")

		if !IsAlphaNumeric(toUserID) {
			c.JSON(http.StatusBadRequest, APIResponse{
				Code:     http.StatusBadRequest,
				Status:   http.StatusText(http.StatusBadRequest),
				Message:  constants.UsernameCantBeEmpty,
				Response: nil,
			})
			return
		} else if !IsAlphaNumeric(fromUserID) {
			c.JSON(http.StatusBadRequest, APIResponse{
				Code:     http.StatusBadRequest,
				Status:   http.StatusText(http.StatusBadRequest),
				Message:  constants.UsernameCantBeEmpty,
				Response: nil,
			})
			return
		}

		page, err := strconv.Atoi(c.Query("page"))
		if err != nil{
			log.Panic(err)
		}
		var pagee int64 = int64(page)

		conversations := GetConversationBetweenTwoUsers(toUserID, fromUserID, pagee)
		c.JSON(http.StatusOK, APIResponse{
			Code:     http.StatusOK,
			Status:   http.StatusText(http.StatusOK),
			Message:  constants.SuccessfulResponse,
			Response: conversations,
		})
	}
}