package handlers

import (
	"fmt"
	"net/http"
	"server/auth"

	"github.com/gin-gonic/gin"
)

type LoginForm struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

//A sample use
var user = LoginForm{
	Username: "username",
	Password: "password",
}

func (h *Handler) Login(c *gin.Context) {
	if _, logged := c.Get("uid"); logged {
		fmt.Println("Sei già loggato zio")
		c.Redirect(http.StatusMovedPermanently, "/")
		c.Abort()
		return
	}

	var loginData LoginForm
	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}
	//compare the user from the request, with the one we defined:
	if user.Username != loginData.Username || user.Password != loginData.Password {
		c.JSON(http.StatusUnauthorized, "Please provide valid login details")
		return
	}
	uid := "UID-sdfjsodifjs"

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
