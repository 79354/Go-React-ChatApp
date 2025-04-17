package handlers

// Lobby maintains list of active clients and broadcasts messages to client
type Lobby struct{
	// Registered Clients
	clients map[*Client]bool
	
	register chan *Client
	unregister chan *Client
}

func NewLobby() *Lobby{
	return &Lobby{
		clients: make(map[*Client]bool),
		register: make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (lobby *Lobby) Run(){
	for {
		select	{
		case client := <- lobby.register:
			
		case client := <- lobby.unregister:
			
		}
	}
}