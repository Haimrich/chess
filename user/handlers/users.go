package handlers

import (
	"context"
	"net/http"
	"time"
	"user/db"
	"user/models"

	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (h *Handler) UserByUsername(c *gin.Context) {

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

func (h *Handler) UserById(c *gin.Context) {
	uid, err := primitive.ObjectIDFromHex(c.Param("uid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
	}

	var user *models.User
	err = h.DB.Collection(db.UsersCollectionName).FindOne(context.TODO(), bson.M{"_id": uid}).Decode(&user)

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

type UpdateStatusForm struct {
	Uid    string `json:"uid" form:"uid" binding:"required"`
	Status string `json:"status" form:"status" binding:"required"`
}

func (h *Handler) UpdateUserStatus(c *gin.Context) {
	var statusData UpdateStatusForm
	if err := c.ShouldBind(&statusData); err == nil {
		id, _ := primitive.ObjectIDFromHex(statusData.Uid)
		filter := bson.M{"_id": id}

		var update bson.D
		if statusData.Status == "offline" {
			update = bson.D{
				{
					Key: "$set", Value: bson.D{
						{Key: "status", Value: statusData.Status},
						{Key: "last-seen", Value: time.Now()},
					},
				}}
		} else if statusData.Status == "online" {
			update = bson.D{
				{
					Key: "$set", Value: bson.D{
						{Key: "status", Value: statusData.Status},
					}},
			}
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Stato utente non valido.",
			})
			return
		}

		user := &models.User{}
		err := h.DB.Collection(db.UsersCollectionName).FindOneAndUpdate(context.TODO(), filter, update).Decode(user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": err.Error(),
			})
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data":    user,
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

}

type UpdateCurrentGameForm struct {
	Uid           string `json:"uid" form:"uid" binding:"required"`
	CurrentGameId string `json:"current-game-id" form:"current-game-id" binding:"required"`
}

func (h *Handler) UpdateCurrentGame(c *gin.Context) {
	var gameData UpdateCurrentGameForm
	if err := c.ShouldBind(&gameData); err == nil {
		id, _ := primitive.ObjectIDFromHex(gameData.Uid)
		filter := bson.M{"_id": id}

		update := bson.D{{
			Key: "$set", Value: bson.D{
				{Key: "current-game-id", Value: gameData.CurrentGameId},
			}},
		}

		user := &models.User{}
		err := h.DB.Collection(db.UsersCollectionName).FindOneAndUpdate(context.TODO(), filter, update).Decode(user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": err.Error(),
			})
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data":    user,
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
}
