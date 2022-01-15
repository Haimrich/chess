package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/Shopify/sarama"
)

var KAFKA_ADDRESS = os.Getenv("KAFKA_ADDRESS")

var KAFKA_DISPATCHER_GROUP_ID = os.Getenv("KAFKA_DISPATCHER_GROUP_ID")
var KAFKA_DISPATCHER_INSTANCE_ID = os.Getenv("KAFKA_DISPATCHER_INSTANCE_ID")

var KAFKA_INBOUND_TOPIC = os.Getenv("KAFKA_INBOUND_TOPIC")

var KAFKA_CHALLENGE_TOPIC = os.Getenv("KAFKA_CHALLENGE_TOPIC")
var KAFKA_GAME_TOPIC = os.Getenv("KAFKA_GAME_TOPIC")

var USER_SERVICE_ADDRESS = os.Getenv("USER_SERVICE_ADDRESS")

func main() {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.ClientID = KAFKA_DISPATCHER_INSTANCE_ID
	config.Metadata.AllowAutoTopicCreation = true
	config.Producer.RequiredAcks = sarama.WaitForAll

	group, err := sarama.NewConsumerGroup(strings.Split(KAFKA_ADDRESS, ","), KAFKA_DISPATCHER_GROUP_ID, config)
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

	ctx := context.Background()
	for {
		topics := []string{KAFKA_INBOUND_TOPIC}
		handler := ConsumerGroupHandler{kafkaProducer}

		err := group.Consume(ctx, topics, handler)
		if err != nil {
			panic(err)
		}
	}
}

type InboundMessage struct {
	// UID mittente messaggi da websocket
	UID         string                 `json:"uid"`
	MessageType string                 `json:"type"`
	Content     map[string]interface{} `json:"content"`
}

type MoveMessage struct {
	GameId string `json:"game-id"`
	UID    string `json:"uid"`
	Move   string `json:"move"`
}

type ChallengeMessage struct {
	Type string `json:"type"`
	From string `json:"from"`
	To   string `json:"to"`
}

type ConsumerGroupHandler struct {
	kafkaProducer sarama.AsyncProducer
}

func (ConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (ConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (h ConsumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		fmt.Printf("[DISPATCHER] WS -> Dispatcher: %s\n", string(msg.Value))

		var message InboundMessage
		if err := json.Unmarshal(msg.Value, &message); err != nil {
			fmt.Printf("[DISPATCHER] Invalid Message.\n")
			continue
		}

		switch message.MessageType {

		case "play-move":
			gameId, _ := message.Content["game-id"].(string)
			move, _ := message.Content["move"].(string)

			m := MoveMessage{
				UID:    message.UID,
				GameId: gameId,
				Move:   move,
			}
			h.SendKafkaMessage(KAFKA_GAME_TOPIC, gameId, m)

		case "challenge-request":
			if uid, ok := message.Content["uid"].(string); ok {
				m := ChallengeMessage{
					Type: "request",
					From: message.UID,
					To:   uid,
				}
				h.SendKafkaMessage(KAFKA_CHALLENGE_TOPIC, message.UID, m)
			}

		case "challenge-accept":
			if uid, ok := message.Content["uid"].(string); ok {
				m := ChallengeMessage{
					Type: "accept",
					From: uid,
					To:   message.UID,
				}
				h.SendKafkaMessage(KAFKA_CHALLENGE_TOPIC, uid, m)
			}

		case "challenge-computer":
			m := ChallengeMessage{
				Type: "computer",
				From: message.UID,
				To:   "computer",
			}
			h.SendKafkaMessage(KAFKA_CHALLENGE_TOPIC, message.UID, m)

		case "connect":
			h.UpdateUserStatus(message.UID, "online")
		case "disconnect":
			h.UpdateUserStatus(message.UID, "offline")
		}

		sess.MarkMessage(msg, "")
	}
	return nil
}

func (h ConsumerGroupHandler) SendKafkaMessage(topic string, key string, data interface{}) {
	md, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("[DISPATCHER] Marshalling error: %v\n", err)
		return
	}
	kafkaMessage := sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.ByteEncoder(md),
	}
	h.kafkaProducer.Input() <- &kafkaMessage
}

type StatusUpdateForm struct {
	UID    string `json:"uid"`
	Status string `json:"status"`
}

func (h ConsumerGroupHandler) UpdateUserStatus(uid string, status string) {
	formData := StatusUpdateForm{UID: uid, Status: status}
	jsonValue, _ := json.Marshal(formData)

	resp, err := http.Post(USER_SERVICE_ADDRESS+"status", "application/json", bytes.NewBuffer(jsonValue))
	if err != nil || resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Printf("[DISPATCHER] User service update status error: %v\n%s", err, string(body))
	}
}
