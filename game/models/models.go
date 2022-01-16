package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Move struct {
	Time   time.Time
	Move   string
	Result string
}

type Player struct {
	ID            string
	Color         string
	RemainingTime time.Duration `bson:"remaining-time" json:"remaining-time"`
}

type Game struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"game-id"`
	Players         []Player
	PlayerToMove    int    `bson:"player-to-move" json:"player-to-move"`
	CurrentPosition string `bson:"current-position" json:"current-position"`
	Moves           []Move
}
