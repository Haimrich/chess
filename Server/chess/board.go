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
	BOARD_SIZE       = 8
)

func (p Piece) Set(flag Piece) Piece    { return p | flag }
func (p Piece) Clear(flag Piece) Piece  { return p &^ flag }
func (p Piece) Toggle(flag Piece) Piece { return p ^ flag }
func (p Piece) Has(flag Piece) bool     { return p&flag != 0 }

type Board struct {
	board [BOARD_SIZE][BOARD_SIZE]uint8
}

func (b *Board) PlayMove(move *Move) bool {
	return true
}
