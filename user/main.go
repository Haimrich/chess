package main

import (
	"log"
	"user/auth"
	"user/db"
	"user/handlers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/penglongli/gin-metrics/ginmetrics"
)

func main() {
	dbc := db.Connect()
	defer db.Disconnect(dbc)

	h := handlers.NewHandler(dbc)

	metrics := gin.New()
	m := ginmetrics.GetMonitor()
	m.SetMetricPath("/metrics")

	public := gin.Default()
	public.Use(cors.Default())
	m.UseWithoutExposingEndpoint(public)

	public.POST("/signup", h.Signup)
	public.POST("/login", auth.TokenAuthMiddleware(false), h.Login)
	public.GET("/logout", auth.TokenAuthMiddleware(true), h.Logout)

	public.GET("/user/username/:username", auth.TokenAuthMiddleware(true), h.UserByUsername)
	public.GET("/user/id/:uid", auth.TokenAuthMiddleware(true), h.UserById)

	public.GET("/users/online", auth.TokenAuthMiddleware(true), h.OnlineUsers)

	public.GET("/", auth.TokenAuthMiddleware(false), h.Home)

	public.Static("/avatar", "avatar")

	private := gin.Default()
	private.Use(cors.Default())
	m.UseWithoutExposingEndpoint(private)

	private.POST("/status", h.UpdateUserStatus)
	private.POST("/game", h.UpdateCurrentGame)

	m.Expose(metrics)
	go func() { log.Fatal(metrics.Run(":2112")) }()

	go func() { log.Fatal(private.Run(":8070")) }()
	log.Fatal(public.Run(":8080"))
}
