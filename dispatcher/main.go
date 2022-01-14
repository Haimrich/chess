package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
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

			md, _ := json.Marshal(MoveMessage{
				UID:    message.UID,
				GameId: gameId,
				Move:   move,
			})
			kafkaMessage := sarama.ProducerMessage{
				Topic: KAFKA_GAME_TOPIC,
				Key:   sarama.StringEncoder(gameId),
				Value: sarama.ByteEncoder(md),
			}
			h.kafkaProducer.Input() <- &kafkaMessage

		case "challenge-request":
			if uid, ok := message.Content["uid"].(string); ok {
				md, _ := json.Marshal(ChallengeMessage{
					Type: "request",
					From: message.UID,
					To:   uid,
				})
				kafkaMessage := sarama.ProducerMessage{
					Topic: KAFKA_CHALLENGE_TOPIC,
					Key:   sarama.StringEncoder(message.UID),
					Value: sarama.ByteEncoder(md),
				}
				h.kafkaProducer.Input() <- &kafkaMessage
			}

		case "challenge-accept":
			if uid, ok := message.Content["uid"].(string); ok {
				md, _ := json.Marshal(ChallengeMessage{
					Type: "accept",
					From: uid,
					To:   message.UID,
				})
				kafkaMessage := sarama.ProducerMessage{
					Topic: KAFKA_CHALLENGE_TOPIC,
					Key:   sarama.StringEncoder(uid),
					Value: sarama.ByteEncoder(md),
				}
				h.kafkaProducer.Input() <- &kafkaMessage
			}

		case "challenge-computer":
			md, _ := json.Marshal(ChallengeMessage{
				Type: "computer",
				From: message.UID,
				To:   "computer",
			})
			kafkaMessage := sarama.ProducerMessage{
				Topic: KAFKA_CHALLENGE_TOPIC,
				Key:   sarama.StringEncoder(message.UID),
				Value: sarama.ByteEncoder(md),
			}
			h.kafkaProducer.Input() <- &kafkaMessage

		}
		sess.MarkMessage(msg, "")
	}
	return nil
}
