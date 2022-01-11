package handlers

import (
	"backend/websocket"

	"go.mongodb.org/mongo-driver/mongo"
)

type Handler struct {
	DB *mongo.Database
	WS *websocket.Hub
}

func NewHandler(dbm *mongo.Database, ws *websocket.Hub) *Handler {
	return &Handler{
		DB: dbm,
		WS: ws,
	}
}
