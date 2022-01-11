package main

import (
	"backend/auth"
	"backend/db"
	"backend/handlers"
	"backend/websocket"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	dbc := db.Connect()
	defer db.Disconnect(dbc)

	wsHub := websocket.NewHub(dbc)
	go wsHub.Listen()

	h := handlers.NewHandler(dbc, wsHub)

	router := gin.Default()
	router.Use(cors.Default())
	/*
		router.Use(cors.New(cors.Config{
			//AllowOrigins:     []string{"http://localhost"},
			//AllowOriginFunc: func(origin string) bool {
			//	return strings.Contains(origin, "http://localhost")
			//},
			AllowAllOrigins: true,

			AllowMethods:     []string{"GET", "POST"},
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
			AllowHeaders:     []string{"Origin", "Content-Type", "Cookie"},
		}))
	*/

	router.POST("/signup", h.Signup)
	router.POST("/login", auth.TokenAuthMiddleware(false), h.Login)
	router.GET("/logout", auth.TokenAuthMiddleware(true), h.Logout)

	router.GET("/user/username/:username", auth.TokenAuthMiddleware(true), h.UserByUsername)
	router.GET("/user/id/:uid", auth.TokenAuthMiddleware(true), h.UserById)

	router.GET("/users/online", auth.TokenAuthMiddleware(true), h.OnlineUsers)

	router.GET("/", auth.TokenAuthMiddleware(false), h.Home)
	router.GET("/ws", auth.TokenAuthMiddleware(true), h.Websocket)

	router.Static("/avatar", "avatar")
	log.Fatal(router.Run(":8080"))
}
