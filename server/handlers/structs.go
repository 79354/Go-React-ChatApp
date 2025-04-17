package handlers

import (
	"github.com/gorilla/websocket"
	"time"
)

type UserDetails struct {
	ID       string `bson:"_id,omitempty"`
	Username string `json:"username" binding:"required" bson:"username"`
	Password string `json:"-" bson:"password"`
	Online   string `json:"online" bson:"online"`
	SocketID  string    `json:"socketId,omitempty" bson:"socketId,omitempty"`
	CreatedAt time.Time `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
}

type Message struct {
	ID         string    `json:"id" bson:"_id,omitempty"`
	Message    string    `json:"message" binding:"required" bson:"message"`
	ToUserID   string    `json:"toUserID" binding:"required" bson:"toUserID"`
	FromUserID string    `json:"fromUserID" binding:"required" bson:"fromUserID"`
	CreatedAt  time.Time `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
}

// Registration data and login credentials
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
type RegistrationRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// user data returned to clients
type UserResponse struct {
	Username string `json:"username"`
	UserID   string `json:"userID"`
	Online   string `json:"online"`
}

type SocketEvent struct {
	EventName	 string		 `json:"eventname"`
	EventPayload interface{} `json:"eventpayload"`
}

type Client struct {
	Lobby   *Lobby
	Conn    *websocket.Conn
	Send    chan SocketEvent
	UserID  string
}

type MessagePayload struct {
	FromUserID string `json:"from_user_id" binding:"required"`
	ToUserID   string `json:"to_user_id" binding:"required"`
	Message    string `json:"message" binding:"required"`
}

type APIResponse struct{
	Code 	 int          `json:"code"`
	Status   string       `json:"status"`
	Message  string       `json:"message"`
	Response interface{}  `json:"response"`
}