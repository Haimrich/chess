package websocket

import (
	"context"
	"server/db"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (h *Hub) updateUserStatus(uid string, status string) {
	id, _ := primitive.ObjectIDFromHex(uid)

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

	h.db.Collection(db.UsersCollectionName).UpdateByID(context.TODO(), id, update)
}
