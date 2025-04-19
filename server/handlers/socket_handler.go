package handlers

import (
	"net/http"
	"bytes"
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait = 10 *time.Second		// prevents server hang, conn doesnt wait forever to send
	pongWait = 60 *time.Second		// keeps the server waiting too long, if client disconnects
	pingPeriod = (pongWait*9)/ 10   // sends regular pings to check if client is active
	maxMessageSize = 512			// prevents memory abuse
)

// Upgrader specifies parameters for upgrading an HTTP connection to a WebSocket connection
var Upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    // This CheckOrigin function allows connections from any origin
    // In production, you might want to be more restrictive
    CheckOrigin: func(r *http.Request) bool {
        return true
    },
}

func HandleSocketPayloadEvents(client *Client, socketEventPayload SocketEvent){
	type chatListResponse struct{
		Type 	 string 	 `json:"type"`
		Chatlist interface{} `json:"chatlist"`
	}

	switch socketEventPayload.EventName {
	case "join":
		userID := (socketEventPayload.EventPayload).(string)
		userDetails := GetUserByUserID(userID)

		if userDetails == (UserDetails{}){
			log.Println("An invalid user with userID " + userID + " tried to connect to Chat Server.")
		} else{
			if userDetails.Online == "N"{
				log.Println(userID + " tried to connect to Chat Server.")
			}else{
				// Announce in the chatlist arrival of new client
				newUserOnlinePayload := SocketEvent{
					EventName: "chatlist-response",
					EventPayload: chatListResponse{
						Type: "new-user-joined",
						Chatlist: UserResponse{
							Username: userDetails.Username,
							UserID: userDetails.ID,
							Online: userDetails.Online,
						},
					},
				}

				BroadcastToEveryoneExceptme(client.Lobby, newUserOnlinePayload, userID)

				// For the client to see everyone that's online
				allOnlineUsersPayload := SocketEvent{
					EventName: "chatlist-response",
					EventPayload: chatListResponse{
						Type: "my-chatlist",
						Chatlist: GetAllOnlineUsers(userDetails.ID),
					},
				}

				EmitToClient(client.Lobby, allOnlineUsersPayload, userDetails.ID)
			}
		}
	case "disconnect":
		if socketEventPayload.EventPayload != nil{
			userID := (socketEventPayload.EventPayload).(string)
			userDetails := GetUserByUserID(userID)
			UpdateUserOnlineStatusByUserID(userID, "N")

			BroadcastToEveryone(client.Lobby, SocketEvent{
				EventName: "chatlist-response",
				EventPayload: chatListResponse{
					Type: "user-disconnected",
					Chatlist: UserResponse{
						Online: "N",
						UserID: userDetails.ID,
						Username: userDetails.Username,
					},
				},
			})
		}
	case "message":
		//decoding JSON into Go types using the encoding/json package without a struct, Go uses this:
		//  map[string]interface{}

		message	   := (socketEventPayload.EventPayload.(map[string]interface{})["message"]).(string)
		toUserID   := (socketEventPayload.EventPayload.(map[string]interface{})["toUserID"]).(string)
		fromUserID := (socketEventPayload.EventPayload.(map[string]interface{})["fromUserID"]).(string)

		if message != "" && fromUserID != "" && toUserID != "" {
			messagePacket := MessagePayload{
				FromUserID: fromUserID,
				Message: message,
				ToUserID: toUserID,
			}
			StoreNewMessages(messagePacket)
			payload := SocketEvent{
				EventName: "message-response",
				EventPayload: messagePacket,
			}

			EmitToClient(client.Lobby, payload, toUserID)
		}
	}
}

func setSocketPayloadReadConfig(c *Client){
	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))	// deadline for pong reponse
	c.Conn.SetPongHandler(func (string) error { c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil }) // extends the pong deadline
}

// client stays in the lobby, but client side Conn is closed
func unRegisterAndCloseConn(c *Client){
	c.Lobby.unregister <- c
	c.Conn.Close()
}

func (c *Client) readPump(){
	var socketEvenPayload SocketEvent

	defer unRegisterAndCloseConn(c)

	setSocketPayloadReadConfig(c)

	for {
		_, payload, err := c.Conn.ReadMessage()

		decoder := json.NewDecoder(bytes.NewReader(payload))
		decoderErr := decoder.Decode(&socketEvenPayload)

		if decoderErr != nil{
			log.Printf("error: %v", decoderErr)
			break
		}

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error ===: %v", err)
			}
			break
		}

		HandleSocketPayloadEvents(c, socketEvenPayload)
	}
}

// sends mssg, from: server to client
// batches the messages
// sends periodic ping
// cleans up gracefully on errors or disconnects
func (c *Client) writePump(){
	// starts a ticker to trigger pings every pingPeriod(54 seconds)
	ticker := time.NewTicker(pingPeriod)
	defer func(){
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case payload, ok := <- c.Send:

			// struct Buffer implements the interface io.Writer{Write(p []byte) (n int, err error)}
			reqBodyBytes := new(bytes.Buffer)
			json.NewEncoder(reqBodyBytes).Encode(payload)
			finalPayload := reqBodyBytes.Bytes()	// returns the unread portion of Buffer

			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok{
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil{
				return
			}

			w.Write(finalPayload)

			n := len(c.Send)
			for i:= 0; i < n; i++{
				json.NewEncoder(reqBodyBytes).Encode(<-c.Send)
				w.Write(reqBodyBytes.Bytes())
			}

			if err := w.Close(); err != nil{
				return
			}

		// This sends a ping message every pingPeriod to check if the client is still connected.
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil{
				return
			}
		}
	}
}

func CreateClient(lobby *Lobby, connection *websocket.Conn, userID string){
	client := &Client{
		Lobby: lobby,
		Conn: connection,
		Send: make(chan SocketEvent),
		UserID: userID,
	}

	go client.writePump() // uses ping, mssg: server 
	go client.readPump() // uses pong

	client.Lobby.register <- client
}

// Join for new Socket Users
func HandleUserRegisterEvent(lobby *Lobby, client *Client){
	lobby.clients[client] = true
	HandleSocketPayloadEvents(client, SocketEvent{
		EventName: "join",
		EventPayload: client.UserID,
	})
}

// Disconnect for Socket Users
func HandleUserDisconnectEvent(lobby *Lobby, client *Client){
	_, ok := lobby.clients[client]
	if ok{
		// remove client from lobby and close the communication channel
		delete(lobby.clients, client)
		close(client.Send)

		// close the websocket connection
		HandleSocketPayloadEvents(client, SocketEvent{
			EventName: "disconnect",
			EventPayload: client.UserID,
		})
	}
}

func EmitToClient(lobby *Lobby, payload SocketEvent, userID string){

	for client := range lobby.clients{
		if client.UserID == userID{
			select {
			case client.Send <- payload:
			default:
				close(client.Send)
				delete(lobby.clients, client)
			}
		}
	}
}

func BroadcastToEveryone(lobby *Lobby, payload SocketEvent){
	for client := range lobby.clients{
		select{
		case client.Send <- payload:
		default:
			close(client.Send)
			delete(lobby.clients, client)
		}
	}
}

func BroadcastToEveryoneExceptme(lobby *Lobby, payload SocketEvent, myUserID string){
	for client := range lobby.clients{
		if client.UserID != myUserID{
			select {
			case client.Send <- payload:
			default:
				close(client.Send)
				delete(lobby.clients, client)
			}
		}
	}
}