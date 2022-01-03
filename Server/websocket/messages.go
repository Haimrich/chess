package websocket

import (
	"encoding/json"
	"fmt"
	"time"
)

type Message struct {
	// UID destinazione serve all'hub per inviare al giusto websocket
	Destination string                 `json:"-"`
	MessageType string                 `json:"type"`
	Content     map[string]interface{} `json:"content"`
}

func (h *Hub) ParseIncomingMessage(source *Client, messageBytes []byte) {
	var message Message
	if err := json.Unmarshal(messageBytes, &message); err != nil {
		fmt.Println(logPrefix + "Formato messaggio non valido: " + string(messageBytes))
		return
	}
	switch message.MessageType {
	case "challenge-request":
		h.HandleChallengeRequest(source, message.Content)
	case "challenge-accept":
		h.HandleChallengeAccept(source, message.Content)
	case "challenge-decline":
		h.HandleChallengeDecline(source, message.Content)
	case "play-move":
		h.HandlePlayMove(source, message.Content)
	case "resign":
		h.HandleResign(source, message.Content)
	default:
		fmt.Println(logPrefix + "Tipo messaggio non valido: " + string(message.MessageType))
	}
}

func (h *Hub) SendWelcome(uid string) {
	content := map[string]interface{}{
		"message": fmt.Sprintf("Ciao %s.", uid),
	}
	h.channel <- Message{
		Destination: uid,
		MessageType: "welcome",
		Content:     content,
	}
}

func (h *Hub) SendChallenge(uid string, sourceUid string, sourceUsername string) {
	content := map[string]interface{}{
		"message": fmt.Sprintf("%s ti ha sfidato.", sourceUsername),
		"uid":     sourceUid,
	}
	h.channel <- Message{
		Destination: uid,
		MessageType: "challenge-request",
		Content:     content,
	}
}

func (h *Hub) SendChallengeOfflineResponse(uid string) {
	content := map[string]interface{}{
		"message": "L'utente è offline.",
	}
	h.channel <- Message{
		Destination: uid,
		MessageType: "challenge-invalid",
		Content:     content,
	}
}

func (h *Hub) SendChallengeDeclinedResponse(uid string, sourceUsername string) {
	content := map[string]interface{}{
		"message": fmt.Sprintf("%s ha rifiutato la sfida.", sourceUsername),
	}
	h.channel <- Message{
		Destination: uid,
		MessageType: "challenge-declined",
		Content:     content,
	}
}

func (h *Hub) SendChallengeInvalidResponse(uid string) {
	content := map[string]interface{}{
		"message": "La richiesta di sfida non è più valida.",
	}
	h.channel <- Message{
		Destination: uid,
		MessageType: "challenge-invalid",
		Content:     content,
	}
}

func (h *Hub) SendChallengeBusyResponse(uid string, busyUsername string) {
	content := map[string]interface{}{
		"message": busyUsername + " sta giocando un'altra partita.",
	}
	h.channel <- Message{
		Destination: uid,
		MessageType: "challenge-busy",
		Content:     content,
	}
}

func (h *Hub) SendGameStart(uid string, opponentUid string, gameId string, color string, time time.Duration) {
	content := map[string]interface{}{
		"message":  "Partita avviata.",
		"opponent": opponentUid,
		"game-id":  gameId,
		"color":    color,
		"time":     time.Seconds(),
	}
	h.channel <- Message{
		Destination: uid,
		MessageType: "game-start",
		Content:     content,
	}
}

func (h *Hub) SendMovePlayed(uid string, movedColor string, move string, remainingTime time.Duration) {
	content := map[string]interface{}{
		"move":  move,
		"color": movedColor,
		"time":  remainingTime.Seconds(),
	}
	h.channel <- Message{
		Destination: uid,
		MessageType: "move-played",
		Content:     content,
	}
}
