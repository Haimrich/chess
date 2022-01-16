package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"game/chess"
	"game/db"
	"game/models"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
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
	board        *chess.Board
	lastMoveTime time.Time

	PlayMove chan *chess.MoveMessage
	Resign   chan string

	h  *Hub
	db *mongo.Database
}

const timePerPlayer = 10 * time.Minute

// Goroutine che gestisce timer ecc.
func (g *Game) Game() {
	go g.addGameToDatabase()
	for {
		select {
		case move := <-g.PlayMove:
			if move.Player == g.Players[g.playerToMove].ID && g.board.ParseMove(g.Players[g.playerToMove].Color, move.Move) {
				g.board.Print()
				g.updateTimers()
				g.sendMessages(move.Move)
				if g.checkEndGame() {
					g.h.RemoveGame(g)
					return
				}
				g.updatePlayerToMove()
				go g.updateDatabase(move.Move)
				go g.queryChessEngine()
			}
		case <-g.Players[0].Timer.C:
			g.h.SendEndGame(g.Players[1].ID, "victory", +10)
			g.h.SendEndGame(g.Players[0].ID, "lose", -10)
			g.h.RemoveGame(g)
			return
		case <-g.Players[1].Timer.C:
			g.h.SendEndGame(g.Players[0].ID, "victory", +10)
			g.h.SendEndGame(g.Players[1].ID, "lose", -10)
			g.h.RemoveGame(g)
			return
		case <-g.Resign:
			g.h.RemoveGame(g)
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

func (g *Game) checkEndGame() bool {
	g.board.UpdatePossibleMoves()
	whitePlayed := g.Players[g.playerToMove].Color == "white"

	opponentKingInCheck := g.board.KingInCheck(!whitePlayed)
	isStalemate := !g.board.HasPossibleMoves(!whitePlayed)

	if isStalemate {
		if opponentKingInCheck {
			// Vittoria  del tizio che ha giocato
			g.h.SendEndGame(g.Players[g.playerToMove].ID, "victory", +10)
			g.h.SendEndGame(g.Players[(g.playerToMove+1)%2].ID, "lose", -10)
			fmt.Println("Vittoria.")
			return true
		}
		// Pareggio
		g.h.SendEndGame(g.Players[g.playerToMove].ID, "draw", +2)
		g.h.SendEndGame(g.Players[(g.playerToMove+1)%2].ID, "draw", +2)
		fmt.Println("Pareggio.")
		return true
	}
	return false

}

// Database

func (g *Game) addGameToDatabase() {
	players := []models.Player{
		{
			ID:            g.Players[0].ID,
			Color:         g.Players[0].Color,
			RemainingTime: g.Players[0].RemainingTime,
		},
		{
			ID:            g.Players[1].ID,
			Color:         g.Players[1].Color,
			RemainingTime: g.Players[1].RemainingTime,
		},
	}

	gameId, _ := primitive.ObjectIDFromHex(g.ID)

	game := models.Game{
		ID:              gameId,
		Players:         players,
		PlayerToMove:    g.playerToMove,
		CurrentPosition: g.board.GenerateFEN(g.Players[g.playerToMove].Color),
	}

	g.h.db.Collection(db.GamesCollectionName).InsertOne(context.TODO(), game)
}

func (g *Game) updateDatabase(moveNotation string) {
	gameId, _ := primitive.ObjectIDFromHex(g.ID)
	filter := bson.M{"_id": gameId}

	fen := g.board.GenerateFEN(g.Players[g.playerToMove].Color)

	move := models.Move{
		Move:   moveNotation,
		Time:   time.Now(),
		Result: fen,
	}

	update := bson.D{
		{
			Key: "$set", Value: bson.D{
				{Key: "players.0.remaining-time", Value: g.Players[0].RemainingTime},
				{Key: "players.1.remaining-time", Value: g.Players[1].RemainingTime},
				{Key: "player-to-move", Value: g.playerToMove},
				{Key: "current-position", Value: fen},
			}},
		{
			Key: "$push", Value: bson.M{"moves": move},
		},
	}

	g.h.db.Collection(db.GamesCollectionName).UpdateOne(context.TODO(), filter, update)
}

// Interrogazione dell'engine

var ENGINE_ENDPOINT string = os.Getenv("ENGINE_ENDPOINT")

func (g *Game) queryChessEngine() {
	if g.Players[g.playerToMove].ID != "computer" {
		return
	}

	fen := g.board.GenerateFEN(g.Players[g.playerToMove].Color)

	queryValues := map[string]interface{}{"fen": fen, "budget": g.Players[g.playerToMove].RemainingTime}
	queryJson, _ := json.Marshal(queryValues)

	response, _ := http.Post(ENGINE_ENDPOINT, "application/json", bytes.NewBuffer(queryJson))
	move, _ := ioutil.ReadAll(response.Body)

	fmt.Println("Mossa Engine: " + string(move))
	if g.Players[g.playerToMove].Color == "black" {
		move[0] = 'a' + 'h' - move[0]
		move[2] = 'a' + 'h' - move[2]
		move[1] = '1' + '8' - move[1]
		move[3] = '1' + '8' - move[3]
	}
	fmt.Println("Mossa Engine Dopo: " + string(move))

	g.PlayMove <- &chess.MoveMessage{Player: "computer", Move: string(move)}
}
