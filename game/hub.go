package main

import (
	"context"
	"game/chess"
	"game/db"
	"game/models"
	"math/rand"
	"sync"
	"time"

	"github.com/Shopify/sarama"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
		board:        chess.NewBoard(""),
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

func (h *Hub) CheckAndReloadGame(gameIdStr string) bool {
	gameId, _ := primitive.ObjectIDFromHex(gameIdStr)
	filter := bson.M{"_id": gameId}
	result := h.db.Collection(db.GamesCollectionName).FindOne(context.TODO(), filter)
	if result.Err() != nil {
		return false
	}

	var gameData models.Game
	if err := result.Decode(&gameData); err != nil {
		return false
	}

	playerA := Player{
		ID:            gameData.Players[0].ID,
		Color:         gameData.Players[0].Color,
		Timer:         time.NewTimer(gameData.Players[0].RemainingTime),
		RemainingTime: gameData.Players[0].RemainingTime,
	}

	playerB := Player{
		ID:            gameData.Players[1].ID,
		Color:         gameData.Players[1].Color,
		Timer:         time.NewTimer(gameData.Players[1].RemainingTime),
		RemainingTime: gameData.Players[1].RemainingTime,
	}

	if len(gameData.Moves) == 0 {
		playerA.Timer.Stop()
		playerB.Timer.Stop()
	} else if gameData.PlayerToMove == 1 {
		playerA.Timer.Stop()
	} else {
		playerB.Timer.Stop()
	}

	lastMoveTime := time.Time{}
	if len(gameData.Moves) == 0 {
		lastMoveTime = gameData.Moves[len(gameData.Moves)-1].Time
	}

	game := &Game{
		ID:           gameIdStr,
		Players:      [2]Player{playerA, playerB},
		PlayMove:     make(chan *chess.MoveMessage, 2),
		Resign:       make(chan string),
		lastMoveTime: lastMoveTime,
		playerToMove: gameData.PlayerToMove,
		board:        chess.NewBoard(gameData.CurrentPosition),
		h:            h,
		db:           h.db,
	}
	go game.Game()
	h.AddGame(game)

	h.SendGameStart(gameData.Players[0].ID, gameData.Players[1].ID, gameIdStr, gameData.Players[0].Color, timePerPlayer)
	h.SendGameStart(gameData.Players[1].ID, gameData.Players[0].ID, gameIdStr, gameData.Players[1].Color, timePerPlayer)

	if gameData.Players[1].ID == "computer" && gameData.Players[1].Color == "white" {
		go game.queryChessEngine()
	}
	return true
}
