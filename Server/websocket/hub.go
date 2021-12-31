package websocket

import (
	"fmt"

	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/mongo"
)

type Hub struct {
	// Client connessi con WS, la chiave è l'_id dell'utente.
	clients map[string]*Client

	// Canale d'uscita dei messaggi verso i client
	channel chan Message

	// Canale per richieste di registrazione e cancellazione dei client online
	register   chan *Client
	unregister chan *Client

	// DB
	db *mongo.Database
}

const logPrefix = "      WebSocket - "

func NewHub(dbc *mongo.Database) *Hub {
	return &Hub{
		clients:    make(map[string]*Client),
		channel:    make(chan Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		db:         dbc,
	}
}

func (h *Hub) Listen() {
	for {
		select {
		case client := <-h.register:
			// Disconnetti altra sessione se esiste
			if old_client, ok := h.clients[client.uid]; ok {
				delete(h.clients, old_client.uid)
				close(old_client.send)
			}
			h.clients[client.uid] = client

		case client := <-h.unregister:
			if c, ok := h.clients[client.uid]; ok && c == client {
				delete(h.clients, client.uid)
				close(client.send)
			}
		case message := <-h.channel:
			h.clients[message.Destination].send <- message.Content
		}
	}
}

func (h *Hub) Register(uid string, conn *websocket.Conn) {

	client := &Client{uid: uid, hub: h, conn: conn, send: make(chan MessageContent, 256)}

	// Segnala il nuovo client alla goroutine Listen() dell'hub
	h.register <- client

	// Goroutine che gestiscono scrittura e lettura sul canale del nuovo client
	go client.Writer()
	go client.Reader()

	fmt.Printf(logPrefix+"Nuovo client registrato con uid %s\n", uid)

	h.SendWelcome(client.uid)
	h.updateUserStatus(client.uid, "online")
}

func (h *Hub) Unregister(c *Client) {
	h.unregister <- c
	c.conn.Close()

	h.updateUserStatus(c.uid, "offline")
}
