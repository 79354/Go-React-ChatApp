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
	"golang.org/x/net/idna"
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

func GetUserByUserID(userID string) error{

}

func IsUsernameAvailableQueryHandler(username string) error{

}

func LoginQueryHandler(userDetails LoginRequest) (UserResponse, error){

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