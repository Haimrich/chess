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

	wsHub := websocket.NewHub()
	go wsHub.Listen()

	h := handlers.NewHandler(db.Connect(), wsHub)
	defer db.Disconnect(h.DB)

	router := gin.Default()
	router.POST("/signup", h.Signup)
	router.POST("/login", auth.TokenAuthMiddleware(false), h.Login)
	router.GET("/logout", auth.TokenAuthMiddleware(true), h.Logout)

	router.GET("/", auth.TokenAuthMiddleware(false), h.Home)
	router.GET("/ws", auth.TokenAuthMiddleware(true), h.Websocket)

	router.Static("/avatar", "avatar")
	log.Fatal(router.Run(":8080"))
}
