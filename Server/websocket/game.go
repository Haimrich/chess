package websocket

import (
	"math/rand"
	"server/chess"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Player struct {
	ID            string
	Color         string
	Timer         *time.Timer
	RemainingTime time.Duration
}

type Game struct {
	ID           string
	Players      [2]Player
	playerToMove int // 0 o 1
	board        chess.Board
	lastMoveTime time.Time

	PlayMove chan *chess.Move
	Resign   chan string

	h  *Hub
	db *mongo.Database
}

const timePerPlayer = 10 * time.Minute

func NewGame(h *Hub, uidPlayerA string, uidPlayerB string) {
	gameId := primitive.NewObjectID().Hex()
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
		PlayMove:     make(chan *chess.Move, 2),
		Resign:       make(chan string),
		lastMoveTime: time.Time{},
		playerToMove: playerToMove,
		h:            h,
		db:           h.db,
	}
	go game.Game()
	h.addGame <- game

	h.SendGameStart(uidPlayerA, uidPlayerB, gameId, colorA, timePerPlayer)
	h.SendGameStart(uidPlayerB, uidPlayerA, gameId, colorB, timePerPlayer)
}

// Goroutine che gestisce timer ecc.
func (g *Game) Game() {

	for {
		select {
		case move := <-g.PlayMove:
			if move.Player == g.Players[g.playerToMove].ID && g.board.PlayMove(move) {
				g.updateTimers()
				g.sendMessages(move.Move)
				g.updatePlayerToMove()
			}
		case <-g.Players[0].Timer.C:
			// TODO player a perde
		case <-g.Players[1].Timer.C:
			// TODO player b perde
		case <-g.Resign:
			g.h.removeGame <- g
			return
		}
	}
}

func (g *Game) updateTimers() {
	otherPlayer := (g.playerToMove + 1) % 2
	g.Players[g.playerToMove].Timer.Stop()
	g.Players[otherPlayer].Timer.Reset(g.Players[otherPlayer].RemainingTime)
	if !g.lastMoveTime.IsZero() {
		g.Players[g.playerToMove].RemainingTime -= time.Since(g.lastMoveTime)
	}
	g.lastMoveTime = time.Now()
}

func (g *Game) sendMessages(move string) {
	otherPlayer := (g.playerToMove + 1) % 2
	g.h.SendMovePlayed(g.Players[otherPlayer].ID, g.Players[g.playerToMove].Color, move, g.Players[g.playerToMove].RemainingTime)
	g.h.SendMovePlayed(g.Players[g.playerToMove].ID, g.Players[g.playerToMove].Color, move, g.Players[g.playerToMove].RemainingTime)
}

func (g *Game) updatePlayerToMove() {
	g.playerToMove = (g.playerToMove + 1) % 2
}
