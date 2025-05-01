package ws

import (
	"fmt"
	"server/enums"
)

type Room struct {
	ID      string             `json:"id"`
	Name    string             `json:"name"`
	Clients map[string]*Client `json:"clients"`
}

type Hub struct {
	Rooms      map[string]*Room `json:"rooms"`
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan *Message
}

func NewHub() *Hub {
	return &Hub{
		Rooms:      make(map[string]*Room),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan *Message, 5),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case cl := <-h.Register:
			// h.Rooms[cl.RoomID].Clients[cl.UserID] = cl
			if _, ok := h.Rooms[cl.RoomID]; ok {
				r := h.Rooms[cl.RoomID]

				if _, ok := r.Clients[cl.ID]; !ok {
					r.Clients[cl.ID] = cl
				}
			}

		case cl := <-h.Unregister:
			if room, ok := h.Rooms[cl.RoomID]; ok {
				if _, ok := room.Clients[cl.ID]; ok {
					// ? broadcast a message saying that the client has left the room
					if len(room.Clients) != 0 {
						h.Broadcast <- &Message{
							Content:  fmt.Sprintf("%s has left the room", cl.Username),
							RoomID:   cl.RoomID,
							UserID:   cl.ID,
							Username: cl.Username,
							Action:   enums.Left,
						}
					}

					delete(room.Clients, cl.ID)
					close(cl.Message)
				}
			}

		case m := <-h.Broadcast:
			if room, ok := h.Rooms[m.RoomID]; ok {
				for _, client := range room.Clients {
					select {
					case client.Message <- m:
					default:
						close(client.Message)
						delete(room.Clients, client.ID)
					}
				}
			}

			// for _, client := range h.Rooms[message.RoomID].Clients {
			// 	select {
			// 	case client.Message <- message:
			// 	default:
			// 		close(client.Message)
			// 		delete(h.Rooms[message.RoomID].Clients, client.UserID)
			// 	}
			// }
		}
	}
}
