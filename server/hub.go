package server

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) Run() {

	for {
		select {
		case client := <-h.register:
			h.clients[client] = true

			logrus.Info(len(h.clients))

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}

		case message := <-h.broadcast:

			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}

func (h *Hub) Send(incidents string)  {

	logrus.Info("send ")

	inc, err := json.Marshal(incidents)

	if err != nil {
		logrus.Info("failed to marshal")
	}

		for client := range h.clients {
			select {
			case client.send <- inc:
				logrus.Info("succes send ")
			default:
				close(client.send)
				delete(h.clients, client)
			}


		}
}

func (h *Hub) GetCountClient() (int)  {

	return len(h.clients)
}