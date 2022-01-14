package main

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/ReneKroon/ttlcache/v2"
	"github.com/Shopify/sarama"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Hub struct {
	Requests      map[string]ttlcache.SimpleCache // Chiave UID sorgente -> UID destinazione
	mu            sync.Mutex
	kafkaProducer sarama.AsyncProducer
}

type GameStartMessage struct {
	GameId  string `json:"game-id"`
	PlayerA string `json:"player-a"`
	PlayerB string `json:"player-b"`
}

func NewHub(producer sarama.AsyncProducer) *Hub {
	return &Hub{
		Requests:      make(map[string]ttlcache.SimpleCache),
		kafkaProducer: producer,
	}
}

func (h *Hub) NewRequest(from string, to string) {
	h.mu.Lock()
	_, ok := h.Requests[from]
	if !ok {
		h.Requests[from] = ttlcache.NewCache()
		h.Requests[from].SetTTL(REQUEST_TTL)
	}
	h.Requests[from].Set(to, true)
	h.mu.Unlock()

	// TODO trovare username tizio

	h.SendChallenge(to, from, "TODO")
}

func (h *Hub) AcceptRequest(from string, to string) {
	h.mu.Lock()
	_, exists := h.Requests[from]
	if !exists {
		return
	}

	_, err := h.Requests[from].Get(to)
	if err != nil {
		return
	}

	h.Requests[from].Purge()
	h.Requests[from].Close()
	delete(h.Requests, from)
	h.mu.Unlock()

	// TODO controllare in qualche modo partite del giocatore in corso

	// Inizia partita scrivendo sul topic dei game
	gameId := primitive.NewObjectID().Hex()
	md, _ := json.Marshal(GameStartMessage{
		GameId:  gameId,
		PlayerA: from,
		PlayerB: to,
	})
	kafkaMessage := sarama.ProducerMessage{
		Topic: KAFKA_GAME_TOPIC,
		Key:   sarama.StringEncoder(gameId),
		Value: sarama.ByteEncoder(md),
	}
	h.kafkaProducer.Input() <- &kafkaMessage
}

func (h *Hub) Computer(uid string) {
	h.mu.Lock()
	_, exists := h.Requests[uid]
	if exists {
		h.Requests[uid].Purge()
		h.Requests[uid].Close()
		delete(h.Requests, uid)
	}
	h.mu.Unlock()

	// TODO controllare in qualche modo partite dei giocatori in corso

	// Inizia partita contro il pc
	gameId := primitive.NewObjectID().Hex()
	md, _ := json.Marshal(GameStartMessage{
		GameId:  gameId,
		PlayerA: uid,
		PlayerB: "computer",
	})
	kafkaMessage := sarama.ProducerMessage{
		Topic: KAFKA_GAME_TOPIC,
		Key:   sarama.StringEncoder(gameId),
		Value: sarama.ByteEncoder(md),
	}
	h.kafkaProducer.Input() <- &kafkaMessage

}

func (h *Hub) Clear() {
	h.mu.Lock()
	for k := range h.Requests {
		h.Requests[k].Purge()
		h.Requests[k].Close()
		delete(h.Requests, k)
	}
	h.mu.Unlock()
}

type Message struct {
	// UID destinazione serve all'hub per inviare al giusto websocket
	Destination string                 `json:"uid"`
	MessageType string                 `json:"type"`
	Content     map[string]interface{} `json:"content"`
}

func (h *Hub) SendChallenge(uid string, sourceUid string, sourceUsername string) {
	content := map[string]interface{}{
		"message": fmt.Sprintf("%s ti ha sfidato.", sourceUsername),
		"uid":     sourceUid,
	}
	md, _ := json.Marshal(Message{
		Destination: uid,
		MessageType: "challenge-request",
		Content:     content,
	})
	kafkaMessage := sarama.ProducerMessage{
		Topic: KAFKA_OUTBOUND_TOPIC,
		Key:   nil,
		Value: sarama.ByteEncoder(md),
	}
	h.kafkaProducer.Input() <- &kafkaMessage
}
