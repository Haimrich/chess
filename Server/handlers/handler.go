package handlers

import "go.mongodb.org/mongo-driver/mongo"

type Handler struct {
	DB *mongo.Database
}

func NewHandler(dbm *mongo.Database) *Handler {
	return &Handler{
		DB: dbm,
	}
}
