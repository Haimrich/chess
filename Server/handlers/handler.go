package handlers

import "go.mongodb.org/mongo-driver/mongo"

type Handler struct {
	DB *mongo.Client
}

func NewHandler(dbm *mongo.Client) *Handler {
	return &Handler{
		DB: dbm,
	}
}
