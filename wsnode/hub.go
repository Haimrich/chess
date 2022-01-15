package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/Shopify/sarama"
	"github.com/gorilla/websocket"
)

type Message struct {
	// UID è destinazione per i messaggi outbound
	// UID è sorgente per i messaggi provenienti dai socket
	UID         string                 `json:"uid,omitempty"`
	MessageType string                 `json:"type"`
	Content     map[string]interface{} `json:"content"`
}

type Hub struct {
	// Client connessi con WS, la chiave è l'_id dell'utente.
	clients    map[string]*Client
	register   chan *Client
	unregister chan *Client

	// Canale d'uscita dei messaggi verso i client
	outbound chan Message

	// Produttore messaggi di heartbeat e di quelli ricevuti dai client
	kafkaProducer sarama.AsyncProducer
	// Consumatore messaggi indirizzati a websocket
	kafkaConsumer sarama.PartitionConsumer
}

func NewHub() *Hub {

	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategySticky
	config.ClientID = KAFKA_WSNODE_INSTANCE_ID
	config.Metadata.AllowAutoTopicCreation = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	// Create new consumer

	consumer, err := sarama.NewConsumer(strings.Split(KAFKA_ADDRESS, ","), config)
	if err != nil {
		log.Fatalf("ERRORE: Kafka cluster non raggiungibile\n%s\n", err.Error())
	}
	kafkaConsumer, _ := consumer.ConsumePartition(KAFKA_OUTBOUND_TOPIC, 0, sarama.OffsetNewest)

	kafkaProducer, _ := sarama.NewAsyncProducer(strings.Split(KAFKA_ADDRESS, ","), config)

	/*
	   // Trap SIGINT to trigger a shutdown. TODO
	   signals := make(chan os.Signal, 1)
	   signal.Notify(signals, os.Interrupt)
	*/

	return &Hub{
		clients:       make(map[string]*Client),
		register:      make(chan *Client),
		unregister:    make(chan *Client),
		outbound:      make(chan Message, 2),
		kafkaProducer: kafkaProducer,
		kafkaConsumer: kafkaConsumer,
	}
}

func (h *Hub) WebsocketListener() {
	defer h.kafkaProducer.Close()
	for {
		select {
		case client := <-h.register:
			// Disconnetti altra sessione se esiste
			if old_client, ok := h.clients[client.uid]; ok {
				delete(h.clients, old_client.uid)
				close(old_client.send)
			}
			h.clients[client.uid] = client

		case client := <-h.unregister:
			if c, ok := h.clients[client.uid]; ok && c == client {
				delete(h.clients, client.uid)
				close(client.send)
			}
		case message := <-h.outbound:
			if h.clients[message.UID] != nil {
				h.clients[message.UID].send <- message
			}
		}
	}
}

func (h *Hub) Register(uid string, conn *websocket.Conn) {

	client := &Client{
		uid:  uid,
		hub:  h,
		conn: conn,
		send: make(chan Message, 2),
	}

	// Segnala il nuovo client alla goroutine Listen() dell'hub
	h.register <- client

	// Goroutine che gestiscono scrittura e lettura sul canale del nuovo client
	go client.Writer()
	go client.Reader()

	fmt.Printf("Nuovo client registrato con uid %s\n", uid)
}

func (h *Hub) Unregister(c *Client) {
	h.unregister <- c
	c.conn.Close()
}

/*
type Heartbeat struct {
	NodeId    string `json:"node-id"`
	Partition int    `json:"partition"`
}
*/

func (h *Hub) KafkaMessageConsumer() {
	for {
		ev := <-h.kafkaConsumer.Messages()

		fmt.Printf("[WSNODE] Kafka -> WS: %s\n", string(ev.Value))

		var message Message
		if err := json.Unmarshal(ev.Value, &message); err == nil {
			h.outbound <- message
		}
	}
}
