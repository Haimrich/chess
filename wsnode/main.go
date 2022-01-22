package main

import (
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/penglongli/gin-metrics/ginmetrics"
)

var KAFKA_ADDRESS = os.Getenv("KAFKA_ADDRESS")

//var KAFKA_WSNODE_GROUP_ID = os.Getenv("KAFKA_WSNODE_GROUP_ID")
var KAFKA_WSNODE_INSTANCE_ID = os.Getenv("KAFKA_WSNODE_INSTANCE_ID")

var KAFKA_INBOUND_TOPIC = os.Getenv("KAFKA_INBOUND_TOPIC")
var KAFKA_OUTBOUND_TOPIC = os.Getenv("KAFKA_OUTBOUND_TOPIC")

func main() {
	// Prometheus cose
	wsMetric := &ginmetrics.Metric{
		Type:        ginmetrics.Gauge,
		Name:        "websocket_connections",
		Description: "Current connected clients.",
		Labels:      []string{"wsnode"},
	}
	inboundMetric := &ginmetrics.Metric{
		Type:        ginmetrics.Counter,
		Name:        "websocket_inbound_messages",
		Description: "Number of received websocket messages.",
		Labels:      []string{"wsnode"},
	}
	outboundMetric := &ginmetrics.Metric{
		Type:        ginmetrics.Counter,
		Name:        "websocket_outbound_messages",
		Description: "Number of sent websocket messages.",
		Labels:      []string{"wsnode"},
	}
	m := ginmetrics.GetMonitor()
	m.AddMetric(wsMetric)
	m.AddMetric(inboundMetric)
	m.AddMetric(outboundMetric)
	m.SetMetricPath("/metrics")
	m.SetMetricPrefix("wsnode_")
	metricsRouter := gin.New()

	// WsNode
	hub := NewHub(m)
	go hub.WebsocketListener()
	go hub.KafkaMessageConsumer()

	router := gin.Default()
	router.Use(cors.Default())
	m.UseWithoutExposingEndpoint(router)

	router.GET("/", hub.Handler)

	m.Expose(metricsRouter)
	go func() { metricsRouter.Run(":2112") }()

	log.Fatal(router.Run(":8081"))
}
