package main

import (
	"context"
	"encoding/json"
	"fmt"
	"game/db"
	"log"
	"os"
	"strings"

	"github.com/Shopify/sarama"
)

var KAFKA_ADDRESS = os.Getenv("KAFKA_ADDRESS")

var KAFKA_GAME_GROUP_ID = os.Getenv("KAFKA_GAME_GROUP_ID")
var KAFKA_GAME_INSTANCE_ID = os.Getenv("KAFKA_GAME_INSTANCE_ID")

var KAFKA_GAME_TOPIC = os.Getenv("KAFKA_GAME_TOPIC")

var KAFKA_OUTBOUND_TOPIC = os.Getenv("KAFKA_OUTBOUND_TOPIC")

func main() {
	dbc := db.Connect()
	defer db.Disconnect(dbc)

	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.ClientID = KAFKA_GAME_INSTANCE_ID
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Consumer.Offsets.AutoCommit.Enable = false

	group, err := sarama.NewConsumerGroup(strings.Split(KAFKA_ADDRESS, ","), KAFKA_GAME_GROUP_ID, config)
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

	hub := NewHub(dbc, kafkaProducer)

	ctx := context.Background()
	for {
		topics := []string{KAFKA_GAME_TOPIC}
		handler := ConsumerGroupHandler{hub: hub}

		err := group.Consume(ctx, topics, handler)
		if err != nil {
			panic(err)
		}
	}
}

type MoveMessage struct {
	GameId string `json:"game-id"`
	UID    string `json:"uid"`
	Move   string `json:"move"`
}

type GameStartMessage struct {
	GameId  string `json:"game-id"`
	PlayerA string `json:"player-a"`
	PlayerB string `json:"player-b"`
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
		fmt.Printf("[GAME] Messaggio ricevuto dal topic: %s\n", string(msg.Value))

		var mm MoveMessage
		var gm GameStartMessage

		if err := json.Unmarshal(msg.Value, &gm); err == nil && gm.PlayerA != "" {
			h.hub.GameStart(gm.GameId, gm.PlayerA, gm.PlayerB)
		} else if err := json.Unmarshal(msg.Value, &mm); err == nil && mm.Move != "" {
			h.hub.PlayMove(mm.GameId, mm.UID, mm.Move)
		} else {
			fmt.Printf("[GAME] Invalid Message.\n")
			continue
		}
	}
	return nil
}
