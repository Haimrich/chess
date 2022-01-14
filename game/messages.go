package main

import (
	"encoding/json"
	"time"

	"github.com/Shopify/sarama"
)

type Message struct {
	// UID destinazione serve all'hub per inviare al giusto websocket
	Destination string                 `json:"uid"`
	MessageType string                 `json:"type"`
	Content     map[string]interface{} `json:"content"`
}

func (h *Hub) SendGameStart(uid string, opponentUid string, gameId string, color string, time time.Duration) {
	content := map[string]interface{}{
		"message":  "Partita avviata.",
		"opponent": opponentUid,
		"game-id":  gameId,
		"color":    color,
		"time":     time.Seconds(),
	}
	md, _ := json.Marshal(Message{
		Destination: uid,
		MessageType: "game-start",
		Content:     content,
	})
	kafkaMessage := sarama.ProducerMessage{
		Topic: KAFKA_OUTBOUND_TOPIC,
		Key:   nil,
		Value: sarama.ByteEncoder(md),
	}
	h.producer.Input() <- &kafkaMessage
}

func (h *Hub) SendMovePlayed(uid string, movedColor string, move string, remainingTime time.Duration) {
	content := map[string]interface{}{
		"move":  move,
		"color": movedColor,
		"time":  remainingTime.Seconds(),
	}
	md, _ := json.Marshal(Message{
		Destination: uid,
		MessageType: "move-played",
		Content:     content,
	})
	kafkaMessage := sarama.ProducerMessage{
		Topic: KAFKA_OUTBOUND_TOPIC,
		Key:   nil,
		Value: sarama.ByteEncoder(md),
	}
	h.producer.Input() <- &kafkaMessage
}

func (h *Hub) SendEndGame(uid string, result string, elo int) {
	content := map[string]interface{}{
		"result": result,
		"elo":    elo,
	}
	md, _ := json.Marshal(Message{
		Destination: uid,
		MessageType: "end-game",
		Content:     content,
	})
	kafkaMessage := sarama.ProducerMessage{
		Topic: KAFKA_OUTBOUND_TOPIC,
		Key:   nil,
		Value: sarama.ByteEncoder(md),
	}
	h.producer.Input() <- &kafkaMessage
}
