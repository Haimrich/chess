package chess

type Castling uint8

const (
	CastlingWhite Castling = 3  // 0011
	CastlingBlack Castling = 12 // 1100
	CastlingQueen Castling = 10 // 1010
	CastlingKing  Castling = 5  // 0101
)

var CastlingLookup = map[rune]Castling{
	'K': CastlingWhite & CastlingKing,
	'Q': CastlingWhite & CastlingQueen,
	'k': CastlingBlack & CastlingKing,
	'q': CastlingBlack & CastlingQueen,
}
