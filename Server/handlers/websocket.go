package handlers

import (
	"fmt"
	"net/http"

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

func (h *Handler) Websocket(c *gin.Context) {
	const logPrefix = "      WebSocket - "

	if uid := c.GetString("uid"); uid != "" {
		conn, err := wsupgrader.Upgrade(c.Writer, c.Request, c.Writer.Header())

		if err == nil {
			h.WS.Register(uid, conn)
		} else {
			fmt.Printf(logPrefix+"Failed to set websocket upgrade: %+v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": err.Error(),
			})
		}
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Accesso non autorizzato.",
		})
	}

}
