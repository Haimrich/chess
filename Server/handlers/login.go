package handlers

import (
	"context"
	"fmt"
	"net/http"
	"server/auth"
	"server/db"
	"server/helpers"
	"server/models"

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
		c.Redirect(http.StatusMovedPermanently, "/")
		c.Abort()
		return
	}

	var loginData LoginForm
	if err := c.ShouldBind(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, "Riempi tutti i campi.")
		return
	}

	result := h.DB.Collection(db.UsersCollectionName).FindOne(context.TODO(), bson.D{{Key: "username", Value: loginData.Username}})
	var user models.User

	if err := result.Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			c.String(http.StatusNotFound, "L'utente %s non esiste.", loginData.Username)
		} else {
			c.String(http.StatusInternalServerError, err.Error())
		}
		return
	}

	if err := helpers.PasswordCompare(loginData.Password, user.Password); err != nil {
		c.JSON(http.StatusUnauthorized, "Password errata.")
		return
	}

	uid := user.ID.Hex()

	accessToken, refreshToken, err := auth.GenerateTokensAndSetCookies(uid, c)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"auth": gin.H{
			"uid":           uid,
			"access-token":  accessToken,
			"refresh-token": refreshToken,
		},
	})
}

func (h *Handler) Logout(c *gin.Context) {
	auth.DeleteTokens(c)
}
