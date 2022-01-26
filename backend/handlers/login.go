package handlers

import (
	"backend/auth"
	"backend/db"
	"backend/helpers"
	"backend/models"
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type LoginForm struct {
	Username string `json:"username" form:"username" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
}

func (h *Handler) Login(c *gin.Context) {
	if _, logged := c.Get("uid"); logged {
		fmt.Println("Sei già loggato zio")
		//c.Redirect(http.StatusMovedPermanently, "/")
		//c.Abort()
		//return
	}

	var loginData LoginForm
	if err := c.ShouldBind(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Riempi tutti i campi." + err.Error(),
		})
		return
	}

	result := h.DB.Collection(db.UsersCollectionName).FindOne(context.TODO(), bson.D{{Key: "username", Value: loginData.Username}})
	var user models.User

	if err := result.Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": fmt.Sprintf("L'utente %s non esiste.", loginData.Username),
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": err.Error(),
			})
		}
		return
	}

	if err := helpers.PasswordCompare(loginData.Password, user.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Password errata.",
		})
		return
	}

	uid := user.ID.Hex()

	accessToken, refreshToken, err := auth.GenerateTokensAndSetCookies(uid, c)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"uid":           uid,
			"username":      user.Username,
			"access-token":  accessToken,
			"refresh-token": refreshToken,
		},
	})
}

func (h *Handler) Logout(c *gin.Context) {
	auth.DeleteTokens(c)
}
