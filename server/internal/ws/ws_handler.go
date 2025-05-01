package ws

import (
	"fmt"
	"log"
	"net/http"
	"server/enums"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Handler struct {
	hub *Hub
}

func NewHandler(h *Hub) *Handler {
	return &Handler{hub: h}
}

type CreateRoomReq struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (h *Handler) CreateRoom(c *gin.Context) {
	var req CreateRoomReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.hub.Rooms[req.ID] = &Room{
		ID:      req.ID,
		Name:    req.Name,
		Clients: make(map[string]*Client),
	}

	c.JSON(http.StatusOK, req)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
		// origin := r.Header.Get("Origin")
		// return origin == "http://localhost:3000" // replace with your own origin
	},
}

func (h *Handler) JoinRoom(c *gin.Context) {
	log.Println("Attempting to upgrade connection to WebSocket")
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("WebSocket upgrade failed:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	log.Println("WebSocket upgrade successful")

	// *** /ws/JoinRoom/:roomId?userId=1&username=user
	// Retrieve roomId, userId, and username from the request
	roomId := c.Param("roomId")
	userId := c.Query("userId")
	username := c.Query("username")

	// ! DEBUGGING!!!
	// if len(h.hub.Rooms) == 0 {
	// 	log.Println("No rooms available")
	// } else {
	// 	log.Println("Existing rooms:")
	// 	for id := range h.hub.Rooms {
	// 		log.Printf("- Room ID: %s\n", id)
	// 	}
	// }

	room := h.hub.Rooms[roomId]
	if room == nil {
		log.Println("Room not found, closing WebSocket connection")
		conn.WriteJSON(gin.H{"error": "Room not found"}) // Send an error message over WebSocket instead
		conn.Close()
		return
	}

	cl := &Client{
		Conn:     conn,
		Message:  make(chan *Message, 10),
		ID:       userId,
		RoomID:   roomId,
		Username: username,
	}

	m := &Message{
		Content:  fmt.Sprintf("%s has joined the room", username),
		RoomID:   roomId,
		UserID:   userId,
		Username: username,
		Action:   enums.Join,
	}

	// ? Register a new client through the register channel
	h.hub.Register <- cl
	// ? Broadcast that message
	h.hub.Broadcast <- m

	// ? writeMessage()
	go cl.writeMessage()

	// ? readMessage()
	cl.readMessage(h.hub)
}

// ! ROOMS
type RoomRes struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (h *Handler) GetRooms(c *gin.Context) {
	rooms := make([]RoomRes, 0)

	for _, r := range h.hub.Rooms {
		rooms = append(rooms, RoomRes{
			ID:   r.ID,
			Name: r.Name,
		})
	}

	c.JSON(http.StatusOK, rooms)
}

// ! CLIENTS
type ClientRes struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

func (h *Handler) GetClients(c *gin.Context) {
	roomId := c.Param("roomId")

	clients := make([]ClientRes, 0)

	if _, ok := h.hub.Rooms[roomId]; !ok {
		c.JSON(http.StatusOK, clients)
	}

	room := h.hub.Rooms[roomId]

	for _, cl := range room.Clients {
		clients = append(clients, ClientRes{
			ID:       cl.ID,
			Username: cl.Username,
		})
	}

	c.JSON(http.StatusOK, clients)
}
