package handlers

import (
	"context"
	"time"
	"os"

	"chat-app/config"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func UpdateUserOnlineByUserID(userId, status string) error{

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