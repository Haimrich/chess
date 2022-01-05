package chess

type Piece uint8

const (
	White            Piece                                = 1 << iota // 00000001
	Moved                                                             // 00000010
	MovesDiagonally                                                   // 00000100
	MovesOrthogonaly                                                  // 00001000
	Knight                                                            // 00010000
	Pawn                                                              // 00100000
	King                                                              // 01000000
	Queen            = MovesDiagonally | MovesOrthogonaly             // 00001100
	Bishop           = MovesDiagonally
	Rook             = MovesOrthogonaly
)

var LETTER_TO_PIECE = map[rune]Piece{
	'R': White | Rook,
	'B': White | Bishop,
	'Q': White | Queen,
	'K': White | King,
	'P': White | Pawn,
	'N': White | Knight,
	'r': Rook,
	'b': Bishop,
	'q': Queen,
	'k': King,
	'p': Pawn,
	'n': Knight,
}

func (p *Piece) Set(flag Piece)     { *p = *p | flag }
func (p *Piece) Clear(flag Piece)   { *p = *p &^ flag }
func (p *Piece) Toggle(flag Piece)  { *p ^= flag }
func (p Piece) Has(flag Piece) bool { return p&flag == flag }
