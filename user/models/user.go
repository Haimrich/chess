package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"uid"`
	Username string             `json:"username"`
	Password string             `json:"-"`
	Avatar   string             `json:"avatar"`
	Elo      int                `json:"elo"`
	Status   string             `json:"status"`
	LastSeen time.Time          `bson:"last-seen" json:"last-seen"`
}
