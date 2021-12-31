package main

import (
	"log"
	"server/auth"
	"server/db"
	"server/handlers"
	"server/websocket"

	"github.com/gin-gonic/gin"
)

func main() {
	dbc := db.Connect()
	defer db.Disconnect(dbc)

	wsHub := websocket.NewHub(dbc)
	go wsHub.Listen()

	h := handlers.NewHandler(dbc, wsHub)

	router := gin.Default()
	router.POST("/signup", h.Signup)
	router.POST("/login", auth.TokenAuthMiddleware(false), h.Login)
	router.GET("/logout", auth.TokenAuthMiddleware(true), h.Logout)

	router.GET("/user/:username", auth.TokenAuthMiddleware(true), h.User)
	router.GET("/users/online", auth.TokenAuthMiddleware(true), h.OnlineUsers)

	router.GET("/", auth.TokenAuthMiddleware(false), h.Home)
	router.GET("/ws", auth.TokenAuthMiddleware(true), h.Websocket)

	router.Static("/avatar", "avatar")
	log.Fatal(router.Run(":8080"))
}
