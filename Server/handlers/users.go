package handlers

import (
	"context"
	"net/http"
	"server/db"
	"server/models"

	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/bson"
)

func (h *Handler) User(c *gin.Context) {

	username := c.Param("username")
	var user *models.User
	err := h.DB.Collection(db.UsersCollectionName).FindOne(context.TODO(), bson.M{"username": username}).Decode(&user)

	if err == nil {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data":    user,
		})
	} else {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Utente non trovato.",
		})
	}
}

func (h *Handler) OnlineUsers(c *gin.Context) {

	filter := bson.D{{Key: "status", Value: "online"}}
	cursor, err := h.DB.Collection(db.UsersCollectionName).Find(context.TODO(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
	}

	var results []models.User
	if err = cursor.All(context.TODO(), &results); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    results,
	})
}
