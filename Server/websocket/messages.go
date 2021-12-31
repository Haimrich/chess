package websocket

import "fmt"

type MessageContent struct {
	MessageType string `json:"type"`
	Message     string `json:"message"`
}

type Message struct {
	Destination string // UID destinazione
	Content     MessageContent
}

func (h *Hub) SendWelcome(uid string) {
	h.channel <- Message{
		Destination: uid,
		Content: MessageContent{
			MessageType: "welcome",
			Message:     "Hello " + uid,
		},
	}
}

func (h *Hub) SendChallenge(uid string, my_username string, my_elo int) {
	h.channel <- Message{
		Destination: uid,
		Content: MessageContent{
			MessageType: "challenge-request",
			Message:     fmt.Sprintf("%s (%d pt.) ti ha sfidato.", my_username, my_elo),
		},
	}
}

func (h *Hub) SendChallengeOfflineResponse(uid string, username string) {
	h.channel <- Message{
		Destination: uid,
		Content: MessageContent{
			MessageType: "challenge-invalid",
			Message:     fmt.Sprintf("%s è offline.", username),
		},
	}
}

func (h *Hub) SendChallengeDeclinedResponse(uid string, username string) {
	h.channel <- Message{
		Destination: uid,
		Content: MessageContent{
			MessageType: "challenge-declined",
			Message:     fmt.Sprintf("%s ha rifiutato la sfida.", username),
		},
	}
}

func (h *Hub) SendChallengeInvalidResponse(uid string) {
	h.channel <- Message{
		Destination: uid,
		Content: MessageContent{
			MessageType: "challenge-invalid",
			Message:     "La richiesta di sfida non è più valida.",
		},
	}
}
