package handlers

import (
	"server/auth"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Logout(c *gin.Context) {
	auth.DeleteTokens(c)
}
