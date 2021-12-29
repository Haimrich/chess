package websocket

import (
	"fmt"

	"github.com/gorilla/websocket"
)

type Hub struct {
	// Client connessi con WS, la chiave è l'_id dell'utente.
	clients map[string]*Client

	// Canale d'ingresso dei messaggi dai client
	channel chan []byte

	// Canale per richieste di registrazione e cancellazione dei client online
	register   chan *Client
	unregister chan *Client
}

const logPrefix = "      WebSocket - "

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[string]*Client),
		channel:    make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) Listen() {
	for {
		select {
		case client := <-h.register:
			h.clients[client.uid] = client
		case client := <-h.unregister:
			if _, ok := h.clients[client.uid]; ok {
				delete(h.clients, client.uid)
				close(client.send)
			}
		case message := <-h.channel:
			fmt.Println(string(message))
		}
	}
}

func (h *Hub) Register(uid string, conn *websocket.Conn) {
	client := &Client{uid: uid, hub: h, conn: conn, send: make(chan []byte, 256)}

	// Segnala il nuovo client alla goroutine Listen() dell'hub
	h.register <- client

	// Goroutine che gestiscono scrittura e lettura sul canale del nuovo client
	go client.Writer()
	go client.Reader()

	fmt.Printf(logPrefix+"Nuovo client registrato con uid %s\n", uid)
}

func (h *Hub) Unregister(c *Client) {
	h.unregister <- c
	c.conn.Close()
}
