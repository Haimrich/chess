package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/Shopify/sarama"
)

var KAFKA_ADDRESS = os.Getenv("KAFKA_ADDRESS")

var KAFKA_CHALLENGE_GROUP_ID = os.Getenv("KAFKA_CHALLENGE_GROUP_ID")
var KAFKA_CHALLENGE_INSTANCE_ID = os.Getenv("KAFKA_CHALLENGE_INSTANCE_ID")

var KAFKA_CHALLENGE_TOPIC = os.Getenv("KAFKA_CHALLENGE_TOPIC")

var KAFKA_GAME_TOPIC = os.Getenv("KAFKA_GAME_TOPIC")

var KAFKA_OUTBOUND_TOPIC = os.Getenv("KAFKA_OUTBOUND_TOPIC")

var REQUEST_TTL = 1 * time.Minute

func main() {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.ClientID = KAFKA_CHALLENGE_INSTANCE_ID
	config.Metadata.AllowAutoTopicCreation = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Consumer.Offsets.AutoCommit.Enable = false

	group, err := sarama.NewConsumerGroup(strings.Split(KAFKA_ADDRESS, ","), KAFKA_CHALLENGE_GROUP_ID, config)
	if err != nil {
		log.Fatalf("ERRORE: Kafka cluster non raggiungibile\n%s\n", err.Error())
	}
	defer func() { _ = group.Close() }()

	go func() {
		for err := range group.Errors() {
			fmt.Println("ERROR", err)
		}
	}()

	kafkaProducer, _ := sarama.NewAsyncProducer(strings.Split(KAFKA_ADDRESS, ","), config)

	hub := NewHub(kafkaProducer)

	ctx := context.Background()
	for {
		topics := []string{KAFKA_CHALLENGE_TOPIC}
		handler := ConsumerGroupHandler{hub: hub}

		err := group.Consume(ctx, topics, handler)
		if err != nil {
			panic(err)
		}
	}
}

type ChallengeMessage struct {
	Type string `json:"type"`
	From string `json:"from"`
	To   string `json:"to"`
}

type ConsumerGroupHandler struct {
	hub *Hub
}

func (c ConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error {
	c.hub.Clear()
	return nil
}

func (c ConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	c.hub.Clear()
	return nil
}

func (h ConsumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		fmt.Printf("[CHALLENGE] Dispatcher -> Challenge: %s\n", string(msg.Value))

		var message ChallengeMessage
		if err := json.Unmarshal(msg.Value, &message); err != nil {
			fmt.Printf("[CHALLENGE] Invalid Message.\n")
			continue
		}

		switch message.Type {

		case "request":
			h.hub.NewRequest(message.From, message.To)

		case "accept":
			h.hub.AcceptRequest(message.From, message.To)

		case "computer":
			h.hub.Computer(message.From)

		}
	}
	return nil
}
