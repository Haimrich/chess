package main

import (
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var KAFKA_ADDRESS = os.Getenv("KAFKA_ADDRESS")

//var KAFKA_WSNODE_GROUP_ID = os.Getenv("KAFKA_WSNODE_GROUP_ID")
var KAFKA_WSNODE_INSTANCE_ID = os.Getenv("KAFKA_WSNODE_INSTANCE_ID")

var KAFKA_INBOUND_TOPIC = os.Getenv("KAFKA_INBOUND_TOPIC")
var KAFKA_OUTBOUND_TOPIC = os.Getenv("KAFKA_OUTBOUND_TOPIC")

func main() {
	hub := NewHub()
	go hub.WebsocketListener()
	go hub.KafkaMessageConsumer()

	router := gin.Default()
	router.Use(cors.Default())

	router.GET("/", hub.Handler)

	log.Fatal(router.Run(":8081"))
}
