package websocket

import "backend/chess"

// Metodi che gestiscono i messaggi ricevuti dai client

func (h *Hub) HandleChallengeRequest(source *Client, content map[string]interface{}) {
	if source.CurrentGameId != "" {
		return
	}

	uidm, ok := content["uid"]
	if !ok {
		return
	}

	uid := uidm.(string)
	if h.clients[uid].CurrentGameId != "" {
		h.SendChallengeBusyResponse(source.uid, h.clients[uid].username)
		return
	}

	source.PendingChallenges[uid] = true
	h.SendChallenge(uid, source.uid, source.username)
}

func (h *Hub) HandleChallengeDecline(source *Client, content map[string]interface{}) {
	uidm, ok := content["uid"]
	if !ok {
		return
	}
	uid := uidm.(string)

	_, exists := h.clients[uid].PendingChallenges[source.uid]
	if exists {
		delete(h.clients[uid].PendingChallenges, source.uid)
		h.SendChallengeDeclinedResponse(uid, source.username)
	}
}

func (h *Hub) HandleChallengeAccept(source *Client, content map[string]interface{}) {

	if source.CurrentGameId != "" {
		return
	}

	uidm, ok := content["uid"]
	if !ok {
		return
	}

	uid := uidm.(string)
	if _, online := h.clients[uid]; !online {
		h.SendChallengeOfflineResponse(source.uid)
		return
	}

	_, challengeExists := h.clients[uid].PendingChallenges[source.uid]
	if h.clients[uid].CurrentGameId != "" || !challengeExists {
		h.SendChallengeInvalidResponse(source.uid)
		return
	}

	NewGame(h, source.uid, uid)
}

func (h *Hub) HandlePlayMove(source *Client, content map[string]interface{}) {
	if source.CurrentGameId == "" {
		return
	}

	movec, ok := content["move"]
	if !ok {
		return
	}

	if move, ok := movec.(string); ok {
		h.games[source.CurrentGameId].PlayMove <- &chess.MoveMessage{Player: source.uid, Move: move}
	}

}

func (h *Hub) HandleResign(source *Client, content map[string]interface{}) {
	// TODO
}

func (h *Hub) HandleChallengeComputer(source *Client) {

	if source.CurrentGameId != "" {
		return
	}

	NewGame(h, source.uid, "computer")
}
