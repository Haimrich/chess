package main

import (
	"game/chess"
	"math/rand"
	"sync"
	"time"

	"github.com/Shopify/sarama"
	"go.mongodb.org/mongo-driver/mongo"
)

type Hub struct {
	// Partite in corso, la chiave è l'id della partita.
	games map[string]*Game

	db       *mongo.Database
	producer sarama.AsyncProducer
	mu       sync.Mutex
}

func NewHub(dbc *mongo.Database, producer sarama.AsyncProducer) *Hub {
	return &Hub{
		games:    make(map[string]*Game),
		db:       dbc,
		producer: producer,
	}
}

func (h *Hub) GameStart(gameId string, uidPlayerA string, uidPlayerB string) {
	// Ricarica partita dal database se esiste già
	if h.CheckAndReloadGame(gameId) {
		return
	}

	colorA, colorB, playerToMove := "white", "black", 0
	if rand.Intn(2) == 0 {
		colorA, colorB, playerToMove = colorB, colorA, 1
	}

	playerA := Player{
		ID:            uidPlayerA,
		Color:         colorA,
		Timer:         time.NewTimer(timePerPlayer),
		RemainingTime: timePerPlayer,
	}
	playerA.Timer.Stop()

	playerB := Player{
		ID:            uidPlayerB,
		Color:         colorB,
		Timer:         time.NewTimer(timePerPlayer),
		RemainingTime: timePerPlayer,
	}
	playerB.Timer.Stop()

	game := &Game{
		ID:           gameId,
		Players:      [2]Player{playerA, playerB},
		PlayMove:     make(chan *chess.MoveMessage, 2),
		Resign:       make(chan string),
		lastMoveTime: time.Time{},
		playerToMove: playerToMove,
		board:        chess.NewBoard(),
		h:            h,
		db:           h.db,
	}
	go game.Game()
	h.AddGame(game)

	h.SendGameStart(uidPlayerA, uidPlayerB, gameId, colorA, timePerPlayer)
	h.SendGameStart(uidPlayerB, uidPlayerA, gameId, colorB, timePerPlayer)

	if uidPlayerB == "computer" && colorB == "white" {
		go game.queryChessEngine()
	}
}

func (h *Hub) AddGame(game *Game) {
	h.mu.Lock()
	h.games[game.ID] = game
	h.mu.Unlock()
}

func (h *Hub) RemoveGame(game *Game) {
	h.mu.Lock()
	delete(h.games, game.ID)
	h.mu.Unlock()
}

func (h *Hub) PlayMove(gameId string, uidPlayer string, move string) {
	if _, ok := h.games[gameId]; ok {
		h.games[gameId].PlayMove <- &chess.MoveMessage{Player: uidPlayer, Move: move}
	}
}

func (h *Hub) Clear() {
	h.mu.Lock()
	for k := range h.games {
		delete(h.games, k)
	}
	h.mu.Unlock()
}

func (h *Hub) CheckAndReloadGame(gameId string) bool {
	return false
}
