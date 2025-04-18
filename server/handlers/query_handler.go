package handlers

import (
	"context"
	"errors"
	"os"
	"time"

	"chat-app/config"
	"chat-app/constants"
	"chat-app/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func UpdateUserOnlineStatusByUserID(userId, status string) error{

}

func GetUserByUsername(username string) UserDetails{
	var userDetails UserDetails
	collection := config.Client.Database(os.Getenv("MONGODB_DATABASE")).Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_ = collection.FindOne(ctx, bson.M{"username": username,}).Decode(&userDetails)

	return userDetails
}

func GetUserByUserID(userID string) UserDetails{
	var userDetails UserDetails

	docID, err := primitive.ObjectIDFromHex(userID)
	if err != nil{
		return UserDetails{}
	}

	collection := config.Client.Database(os.Getenv("MONGODB_DATABASE")).Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	_ = collection.FindOne(ctx, bson.M{
		"_id": docID,
	}).Decode(&userDetails)
	
	cancel()
	return	userDetails
}

func IsUsernameAvailableQueryHandler(username string) bool{
	userDetails := GetUserByUsername(username)
	if userDetails == (UserDetails{}){
		return true
	}
	return false
}

func LoginQueryHandler(userDetailsRequest LoginRequest) (UserResponse, error){
	if userDetailsRequest.Username == "" {
		return UserResponse{}, errors.New(constants.UsernameCantBeEmpty)
	} else if userDetailsRequest.Password == "" {
		return UserResponse{}, errors.New(constants.PasswordCantBeEmpty)
	} else{
		userDetails := GetUserByUsername(userDetailsRequest.Username)
		if userDetails == (UserDetails{}){
			return UserResponse{}, errors.New(constants.UserIsNotRegisteredWithUs)
		}

		if passErr := utils.VerifyPassword(userDetails.Password, userDetailsRequest.Password); passErr != nil{
			return UserResponse{}, errors.New(constants.LoginPasswordIsInCorrect)
		}

		if onlineStatusErr := UpdateUserOnlineStatusByUserID(userDetails.ID, "Y"); onlineStatusErr != nil{
			return UserResponse{}, errors.New(constants.LoginPasswordIsInCorrect)
		}

		return	UserResponse{
			Username: userDetails.Username,
			UserID: userDetails.ID,
		}, nil
	}
}

// check the username from the database
func RegisterQueryHandler(userDetails RegistrationRequest) (string, error){

	if userDetails.Username == ""{
		return "", errors.New(constants.UsernameCantBeEmpty)
	}else if userDetails.Password == ""{
		return "", errors.New(constants.PasswordCantBeEmpty)
	}else{
		newPasswordHash, PassErr := utils.HashPassword(userDetails.Password)
		if PassErr != nil{
			return "", errors.New(constants.ServerFailedResponse)
		}

		collection := config.Client.Database(os.Getenv("MONGODB_DATABASE")).Collection("users")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		id := primitive.NewObjectID()
		uid := id.Hex()

		_, registrationErr := collection.InsertOne(ctx, bson.M{
			"username": userDetails.Username,
			"password": newPasswordHash,
			"online": "N",
			"createdAt": time.Now(),
			"id": id,
		})
		defer cancel()

		if registrationErr != nil{
			return "", errors.New(constants.ServerFailedResponse)
		}

		if onlineStatusError := UpdateUserOnlineStatusByUserID(uid, "Y"); onlineStatusError != nil {
			return " ", errors.New(constants.ServerFailedResponse)
		}
		return uid, nil
	}
}

func GetAllOnlineUsers(userID string) []UserDetails{
	var onlineUsers []UserDetails

	collection := config.Client.Database(os.Getenv("MONGODB_DATABASE")).Collection("users")
	context, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()


}

func StoreNewMessages(message MessagePayload) bool{
	collection := config.Client.Database(os.Getenv("MONGODB_DATABASE")).Collection("messages")

	context, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
}

func GetConversationBetweenTwoUsers(toUser, fromUser string) []Message{
	var conversation []Message
	collection := config.Client.Database(os.Getenv("MONGODB_DATABASE")).Collection("messages")

	context, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()


}