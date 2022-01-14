package main

import (
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

const (
	accessTokenCookieName = "access-token"
)

func (h *Hub) Handler(c *gin.Context) {
	accessTokenCookie, err := c.Request.Cookie(accessTokenCookieName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	accessTokenString := accessTokenCookie.Value
	accessToken, _, err := new(jwt.Parser).ParseUnverified(accessTokenString, jwt.MapClaims{})
	uid := ""

	if err == nil {
		if claims, ok := accessToken.Claims.(jwt.MapClaims); ok {
			uid = claims["uid"].(string)
		}
	}

	if uid == "" {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Token non conforme.",
		})
		return
	}

	conn, err := wsupgrader.Upgrade(c.Writer, c.Request, c.Writer.Header())

	if err == nil {
		h.Register(uid, conn)
	} else {
		fmt.Printf("Failed to set websocket upgrade: %+v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
	}
}
