package ws

import (
	"github.com/coffee/support/internal/controller/http/v1/dto"
)

type Room struct {
	clients map[uint64]*Client
}

type Hub struct {
	rooms      map[uint64]*Room
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan *dto.WsMessage
}

func NewHub() *Hub {
	return &Hub{
		rooms:      make(map[uint64]*Room),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan *dto.WsMessage, 5),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			if _, ok := h.rooms[client.RoomId]; !ok {
				h.rooms[client.RoomId] = &Room{
					clients: make(map[uint64]*Client),
				}
			}

			room := h.rooms[client.RoomId]

			if _, ok := room.clients[client.UserId]; !ok {
				room.clients[client.UserId] = client
			}
		case client := <-h.Unregister:
			if _, ok := h.rooms[client.RoomId]; ok {
				room := h.rooms[client.RoomId]

				if _, ok := room.clients[client.UserId]; ok {
					delete(room.clients, client.UserId)
					close(client.Send)
				}

				if len(room.clients) == 0 {
					delete(h.rooms, client.RoomId)
				}
			}
		case msg := <-h.Broadcast:
			if _, ok := h.rooms[msg.ChatId]; ok {
				for _, client := range h.rooms[msg.ChatId].clients {
					client.Send <- msg
				}
			}
		}
	}
}
