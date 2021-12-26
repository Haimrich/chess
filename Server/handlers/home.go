package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Home(c *gin.Context) {
	uid, authenticated := c.Get("uid")
	if authenticated {
		c.String(http.StatusOK, "Hello %s", uid)
	} else {
		c.String(http.StatusOK, "You are not logged in!")
	}
}
