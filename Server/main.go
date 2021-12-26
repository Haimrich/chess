package main

import (
	"log"
	"server/auth"
	"server/db"
	"server/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	h := handlers.NewHandler(db.Connect())
	defer db.Disconnect(h.DB)

	router := gin.Default()
	router.POST("/login", auth.TokenAuthMiddleware(false), h.Login)
	router.GET("/", auth.TokenAuthMiddleware(false), h.Home)
	router.GET("/logout", auth.TokenAuthMiddleware(true), h.Logout)

	log.Fatal(router.Run(":8080"))
}
