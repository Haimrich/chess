package websocket

import (
	"backend/db"
	"backend/models"
	"context"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Hub struct {
	// Client connessi con WS, la chiave è l'_id dell'utente.
	clients    map[string]*Client
	register   chan *Client
	unregister chan *Client

	// Canale d'uscita dei messaggi verso i client
	channel chan Message

	// Partite in corso, la chiave è l'id della partita.
	games      map[string]*Game
	addGame    chan *Game
	removeGame chan *Game

	// DB
	db *mongo.Database
}

const logPrefix = "      WebSocket - "

func NewHub(dbc *mongo.Database) *Hub {
	return &Hub{
		clients:    make(map[string]*Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		channel:    make(chan Message),
		games:      make(map[string]*Game),
		addGame:    make(chan *Game),
		removeGame: make(chan *Game),
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
			if h.clients[message.Destination] != nil {
				h.clients[message.Destination].send <- message
			}
		case game := <-h.addGame:
			h.games[game.ID] = game
			h.clients[game.Players[0].ID].CurrentGameId = game.ID
			if h.clients[game.Players[1].ID] != nil {
				h.clients[game.Players[1].ID].CurrentGameId = game.ID
			}
		case game := <-h.removeGame:
			delete(h.games, game.ID)
			h.clients[game.Players[0].ID].CurrentGameId = ""
			if h.clients[game.Players[1].ID] != nil {
				h.clients[game.Players[1].ID].CurrentGameId = ""
			}
		}
	}
}

func (h *Hub) Register(uid string, conn *websocket.Conn) {
	user, err := h.updateUserStatus(uid, "online")
	if err != nil {
		fmt.Println(logPrefix + "Errore DB: " + err.Error())
		return
	}

	client := &Client{
		uid:               uid,
		username:          user.Username,
		hub:               h,
		conn:              conn,
		send:              make(chan Message, 4),
		PendingChallenges: make(map[string]bool),
		CurrentGameId:     "",
	}

	// Segnala il nuovo client alla goroutine Listen() dell'hub
	h.register <- client

	// Goroutine che gestiscono scrittura e lettura sul canale del nuovo client
	go client.Writer()
	go client.Reader()

	fmt.Printf(logPrefix+"Nuovo client registrato con uid %s\n", uid)

	h.SendWelcome(client.uid)
}

func (h *Hub) updateUserStatus(uid string, status string) (*models.User, error) {
	id, _ := primitive.ObjectIDFromHex(uid)
	filter := bson.M{"_id": id}

	var update bson.D
	if status == "offline" {
		update = bson.D{
			{
				Key: "$set", Value: bson.D{
					{Key: "status", Value: status},
					{Key: "last-seen", Value: time.Now()},
				},
			}}
	} else {
		update = bson.D{
			{
				Key: "$set", Value: bson.D{
					{Key: "status", Value: status},
				}},
		}
	}

	user := &models.User{}
	err := h.db.Collection(db.UsersCollectionName).FindOneAndUpdate(context.TODO(), filter, update).Decode(user)

	return user, err
}

func (h *Hub) Unregister(c *Client) {
	h.unregister <- c
	c.conn.Close()

	h.updateUserStatus(c.uid, "offline")
}
