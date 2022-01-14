package chess

type MoveMessage struct {
	Player string
	Move   string
}

type Move struct {
	playedByWhite bool
	start         Square
	end           Square
	capturedPiece Piece
	isEnPassant   bool
	castling      Castling
}

// Mappa che associa ad ogni pezzo un array in cui il primo valore è un booleano
// che dice se si deve proseguire lungo le direzioni o no
var DIRECTIONS = map[Piece][]interface{}{
	King:   {false, Square{-1, -1}, Square{-1, 0}, Square{-1, 1}, Square{0, -1}, Square{0, 1}, Square{1, -1}, Square{1, 0}, Square{1, 1}},
	Knight: {false, Square{2, -1}, Square{2, 1}, Square{-2, -1}, Square{-2, 1}, Square{-1, -2}, Square{-1, 2}, Square{1, -2}, Square{1, 2}},
	Queen:  {true, Square{-1, -1}, Square{-1, 0}, Square{-1, 1}, Square{0, -1}, Square{0, 1}, Square{1, -1}, Square{1, 0}, Square{1, 1}},
	Rook:   {true, Square{-1, 0}, Square{0, -1}, Square{0, 1}, Square{1, 0}},
	Bishop: {true, Square{-1, -1}, Square{-1, 1}, Square{1, -1}, Square{1, 1}},
	// Mosse di cattura pedone
	Pawn | White: {false, Square{-1, -1}, Square{-1, 1}},
	Pawn:         {false, Square{1, -1}, Square{1, 1}},
}
