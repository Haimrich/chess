package main

import (
	"log"
	"user/auth"
	"user/db"
	"user/handlers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	dbc := db.Connect()
	defer db.Disconnect(dbc)

	h := handlers.NewHandler(dbc)

	router := gin.Default()
	router.Use(cors.Default())

	router.POST("/signup", h.Signup)
	router.POST("/login", auth.TokenAuthMiddleware(false), h.Login)
	router.GET("/logout", auth.TokenAuthMiddleware(true), h.Logout)

	router.GET("/user/username/:username", auth.TokenAuthMiddleware(true), h.UserByUsername)
	router.GET("/user/id/:uid", auth.TokenAuthMiddleware(true), h.UserById)

	router.GET("/users/online", auth.TokenAuthMiddleware(true), h.OnlineUsers)

	router.GET("/", auth.TokenAuthMiddleware(false), h.Home)

	router.Static("/avatar", "avatar")
	log.Fatal(router.Run(":8080"))
}
