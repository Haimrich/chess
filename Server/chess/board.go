package chess

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

const BOARD_SIZE = 8

type Board struct {
	board           [BOARD_SIZE][BOARD_SIZE]Piece
	enPassantSquare *Square
	sideToPlay      string
	castling        Castling
	halfmoveClock   int
	fullmoveClock   int

	// White == true
	kingPositions map[bool]Square
	possibleMoves map[bool][]Move
}

const regex string = `(?:(?P<castle>^[O0o]-[O0o](?P<long>-[O0o])?)|` +
	`(?P<pawn_move>^(?P<pawn_file_start>[a-h])` +
	`(?:x(?P<pawn_file_end>[a-h]))?(?P<pawn_rank_end>[1-8])` +
	`(?:=(?P<promotion_piece>[QNRB]))?)|` +
	`(?P<move>(?P<piece>[RBQKPN])` +
	`(?P<file_start>[a-h])?(?P<rank_start>[1-8])?` +
	`(?P<capture>[x])?(?P<file_end>[a-h])(?P<rank_end>[1-8]))|` +
	`(?P<uci_move>^(?P<uci_file_start>[a-h])(?P<uci_rank_start>[1-8])(?P<uci_file_end>[a-h])(?P<uci_rank_end>[1-8])))` +
	`(?:[+#])?$`

var moveRegex *regexp.Regexp = regexp.MustCompile(regex)

func (b *Board) ParseMove(color string, move string) bool {
	matches := moveRegex.FindStringSubmatch(move)
	if matches == nil {
		return false
	}

	group_names := moveRegex.SubexpNames()
	result := make(map[string]string)

	for i, match := range matches {
		if match != "" {
			result[group_names[i]] = match
		}
	}

	if _, isCastle := result["castle"]; isCastle {
		_, isLongCastling := result["long"]
		return b.PlayCastling(color, isLongCastling)
	}

	if _, isPawnMove := result["pawn_move"]; isPawnMove {
		startFile := result["pawn_file_start"]
		endRank := result["pawn_rank_end"]
		endFile := result["pawn_file_end"]
		promotionPiece := result["promotion_piece"]
		return b.PlayPawnMove(color, startFile, endRank, endFile, promotionPiece)
	}

	if _, isMove := result["move"]; isMove {
		piece := result["piece"]
		startFile := result["file_start"]
		startRank := result["rank_start"]
		endRank := result["rank_end"]
		endFile := result["file_end"]
		return b.PlayMove(color, piece, endRank, endFile, startRank, startFile)
	}

	if _, isUciMove := result["uci_move"]; isUciMove {
		startFile := result["uci_file_start"]
		startRank := result["uci_rank_start"]
		endRank := result["uci_rank_end"]
		endFile := result["uci_file_end"]
		return b.PlayUciMove(color, endRank, endFile, startRank, startFile)
	}

	return false
}

func (b *Board) PlayCastling(color string, long bool) bool {
	isWhite := color == "white"
	castlingColor, castlingSide := CastlingWhite, CastlingKing
	if !isWhite {
		castlingColor = CastlingBlack
	}
	if long {
		castlingSide = CastlingQueen
	}
	castling := castlingColor & castlingSide
	for _, move := range b.possibleMoves[isWhite] {
		if move.castling == castling {
			b.MovePieces(&move)
			b.enPassantSquare = nil
			b.castling = b.castling &^ castlingColor

			return true
		}
	}
	return false
}

func (b *Board) PlayPawnMove(color string, startFile string, endRank string, endFile string, promotionPiece string) bool {
	isWhite := color == "white"

	for _, m := range b.possibleMoves[isWhite] {
		if m.start.file == FileToIdx(startFile) && m.end.rank == RankToIdx(endRank) &&
			b.GetPieceInSquare(&m.start).Has(Pawn) &&
			b.GetPieceInSquare(&m.start).Has(White) == isWhite &&
			(endFile == "" || m.end.file == FileToIdx(endFile)) {

			b.enPassantSquare = nil
			upOrDown := +1
			if !isWhite {
				upOrDown = -1
			}
			if m.end.rank-m.start.rank == upOrDown*2 {
				b.enPassantSquare = m.start.Translate(&Square{upOrDown, 0})
			} else if promotionPiece != "" {
				newPiece := LETTER_TO_PIECE[[]rune(strings.ToLower(promotionPiece))[0]]
				if isWhite {
					newPiece.Set(White)
				}
				b.board[m.end.rank][m.end.file] = newPiece
			}
			b.board[m.end.rank][m.end.file].Set(Moved)

			b.MovePieces(&m)
			return true
		}
	}
	return false
}

func (b *Board) PlayMove(color string, piece string, endRank string, endFile string, startRankS string, startFileS string) bool {
	end := ParseSquare(endFile + endRank)
	var filteredMoves []Move
	isWhite := color == "white"

	startRank, startFile := -1, -1
	if startRankS != "" {
		startRank = RankToIdx(startRankS)
	}
	if startFileS != "" {
		startFile = FileToIdx(startFileS)
	}

	for _, m := range b.possibleMoves[isWhite] {
		pieceType := LETTER_TO_PIECE[[]rune(strings.ToLower(piece))[0]]
		pMovePiece := b.GetPieceInSquare(&m.start)
		pMovePieceType := pMovePiece
		pMovePieceType.Clear(White | Moved)

		if m.end.file == end.file && m.end.rank == end.rank &&
			pMovePieceType == pieceType &&
			pMovePiece.Has(White) == isWhite &&
			(startRank < 0 || m.start.rank == startRank) &&
			(startFile < 0 || m.start.file == startFile) {
			filteredMoves = append(filteredMoves, m)
		}
	}

	if len(filteredMoves) != 1 {
		return false
	}

	move := filteredMoves[0]
	if b.GetPieceInSquare(&move.start).Has(King) {
		myCastlingColor := CastlingWhite
		if !isWhite {
			myCastlingColor = CastlingBlack
		}
		b.castling = b.castling &^ myCastlingColor

	} else if b.GetPieceInSquare(&move.start).Has(Rook) {
		myCastlingColor := CastlingWhite
		if !isWhite {
			myCastlingColor = CastlingBlack
		}
		rookSide := CastlingKing
		if move.start.file == 0 {
			rookSide = CastlingQueen
		}
		b.castling = b.castling &^ (myCastlingColor & rookSide)
	}

	b.MovePieces(&move)
	b.enPassantSquare = nil

	return true
}

func (b *Board) PlayUciMove(color string, endRank string, endFile string, startRank string, startFile string) bool {
	end := ParseSquare(endFile + endRank)
	start := ParseSquare(startFile + startRank)

	/*
		var filteredMoves []Move
		isWhite := color == "white"

		for _, m := range b.possibleMoves[isWhite] {
			if m.end.file == end.file && m.end.rank == end.rank &&
				m.start.file == start.file && m.start.rank == start.rank {
				filteredMoves = append(filteredMoves, m)
			}
		}

		if len(filteredMoves) != 1 {
			return false
		}
	*/

	// Mossa pedone
	if b.GetPieceInSquare(&start).Has(Pawn) {
		if end.rank == 7 {
			return b.PlayPawnMove(color, startFile, endRank, endFile, "Q")
		} else {
			return b.PlayPawnMove(color, startFile, endRank, endFile, "")
		}
	}

	// Arrocco
	if b.GetPieceInSquare(&start).Has(King) && start.file-end.file == 2 || end.file-start.file == 2 {
		b.PlayCastling(color, start.file > end.file)
	}

	// Mossa normale
	for c, p := range LETTER_TO_PIECE {
		piece := b.GetPieceInSquare(&start)
		piece.Clear(Moved)
		if piece == p {
			return b.PlayMove(color, strings.ToUpper(string(c)), endRank, endFile, startRank, startFile)
		}
	}

	return false
}

func (b *Board) LoadFEN(fen string) {
	fen = strings.Trim(fen, "[]")
	fields := strings.Split(fen, " ")

	// Campo 1 - Posizione
	rank, file := 7, 0
	for _, c := range fields[0] {
		switch {
		case unicode.IsDigit(c): // Skip
			file += int(c - '0')
		case c == '/': // Nuova riga
			rank--
			// if file != 8 || rank < 0 { // TODO errore }
			file = 0
		case LETTER_TO_PIECE[c] != 0: // Pezzo
			b.board[rank][file] = LETTER_TO_PIECE[c]
			if b.board[rank][file].Has(King) {
				b.kingPositions[b.board[rank][file].Has(White)] = Square{rank: rank, file: file}
			} else if b.board[rank][file].Has(Pawn) {
				if (b.board[rank][file].Has(White) && rank != 1) || (!b.board[rank][file].Has(White) && rank != 6) {
					b.board[rank][file].Set(Moved)
				}
			}
			file++
		default:
			// TODO Error
		}
	}

	// Campo 2 - Turno giocatore
	if fields[1] == "w" {
		b.sideToPlay = "white"
	} else {
		b.sideToPlay = "black"
	}

	// Campo 3 - Possibilità di arrocco
	for _, c := range fields[2] {
		if _, ok := CastlingLookup[c]; ok {
			b.castling |= CastlingLookup[c]
		}
	}

	// Campo 4 - En Passant
	if fields[3] != "-" {
		eps := ParseSquare(fields[3])
		b.enPassantSquare = &eps
	}

	// Campo 5 e 6 - Halfmove e fullmove clock
	b.halfmoveClock, _ = strconv.Atoi(fields[4])
	b.fullmoveClock, _ = strconv.Atoi(fields[5])
}

func NewBoard() *Board {
	board := &Board{
		kingPositions: make(map[bool]Square),
		possibleMoves: make(map[bool][]Move),
	}
	board.LoadFEN("[rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1]")
	board.UpdatePossibleMoves()
	board.Print()
	return board
}

func (h *Board) GetPieceInSquare(s *Square) Piece {
	return h.board[s.rank][s.file]
}

func GetPieceInSquareFromBoard(s *Square, board *[BOARD_SIZE][BOARD_SIZE]Piece) Piece {
	return board[s.rank][s.file]
}

func (h *Board) KingInCheck(kingIsWhite bool) bool {
	inCheck := h.SquareInCheck(&h.board, kingIsWhite, h.kingPositions[kingIsWhite])
	return inCheck
}

func (h *Board) SquareInCheck(board *[BOARD_SIZE][BOARD_SIZE]Piece, imWhite bool, square Square) bool {
	pawnDir := Pawn
	if !imWhite {
		pawnDir.Set(White)
	}
	var dirs = []Piece{pawnDir, Knight, MovesDiagonally, MovesOrthogonaly}

	for _, p := range dirs {
		for _, d := range DIRECTIONS[p][1:] {
			offset := d.(Square)
			translatedSquare := square.Translate(&offset)
			// Do while
			for ok := true; ok && translatedSquare != nil; ok = DIRECTIONS[p][0].(bool) {
				piece := h.GetPieceInSquare(translatedSquare)
				if piece != 0 {
					// Se il pezzo trovato può muoversi nella direzione considerata ed è dell'avversario, siamo in pericolo
					if piece.Has(p) && piece.Has(White) != imWhite {
						return true
					}
					break
				}
				translatedSquare = translatedSquare.Translate(&offset)
			}
		}
	}
	return false
}

func (h *Board) UpdatePossibleMoves() {
	h.possibleMoves[true], h.possibleMoves[false] = nil, nil

	for rank, row := range h.board {
		for file, piece := range row {
			isWhite := piece.Has(White)
			moved := piece.Has(Moved)
			piece.Clear(White | Moved)
			square := Square{rank, file}

			// Pedoni...
			if piece.Has(Pawn) {
				upOrDown := +1
				if !isWhite {
					upOrDown = -1
				}
				// Movimento Pedoni
				moveOffset := Square{upOrDown, 0}
				jumpOffset := Square{upOrDown * 2, 0}

				if moveSquare := square.Translate(&moveOffset); moveSquare != nil && h.GetPieceInSquare(moveSquare) == 0 {
					move := Move{playedByWhite: isWhite, start: square, end: *moveSquare}
					if h.MoveIsLegal(&move) {
						h.possibleMoves[isWhite] = append(h.possibleMoves[isWhite], move)
					}
					if !moved {
						if jumpSquare := square.Translate(&jumpOffset); h.GetPieceInSquare(jumpSquare) == 0 {
							move = Move{playedByWhite: isWhite, start: square, end: *jumpSquare}
							if h.MoveIsLegal(&move) {
								h.possibleMoves[isWhite] = append(h.possibleMoves[isWhite], move)
							}
						}
					}
				}
				// Cattura Pedoni
				captureOffsets := []Square{{upOrDown, -1}, {upOrDown, +1}}
				for _, cOffset := range captureOffsets {
					if targetSquare := square.Translate(&cOffset); targetSquare != nil {
						targetPiece := h.GetPieceInSquare(targetSquare)
						if targetPiece != 0 && targetPiece.Has(White) != isWhite {
							move := Move{
								playedByWhite: isWhite,
								start:         square,
								end:           *targetSquare,
								capturedPiece: targetPiece,
							}
							if h.MoveIsLegal(&move) {
								h.possibleMoves[isWhite] = append(h.possibleMoves[isWhite], move)
							}
						} else if h.enPassantSquare != nil && targetSquare.file == h.enPassantSquare.file && targetSquare.rank == h.enPassantSquare.rank {
							// En Passant
							capturedPiece := h.GetPieceInSquare(&Square{square.rank, targetSquare.file})
							if capturedPiece.Has(White) != isWhite && capturedPiece.Has(Pawn) {
								move := Move{
									playedByWhite: isWhite,
									start:         square,
									end:           *targetSquare,
									capturedPiece: capturedPiece,
									isEnPassant:   true,
								}
								if h.MoveIsLegal(&move) {
									h.possibleMoves[isWhite] = append(h.possibleMoves[isWhite], move)
								}
							}
						}
					}
				}

				// Mosse normali
			} else if piece != 0 {
				for _, d := range DIRECTIONS[piece][1:] {
					offset := d.(Square)
					translatedSquare := square.Translate(&offset)

					for ok := true; ok && translatedSquare != nil; ok = DIRECTIONS[piece][0].(bool) {
						targetPiece := h.GetPieceInSquare(translatedSquare)

						if targetPiece != 0 && targetPiece.Has(White) == isWhite {
							break
						}

						move := Move{
							playedByWhite: isWhite,
							start:         square,
							end:           *translatedSquare,
							capturedPiece: targetPiece,
						}

						if h.MoveIsLegal(&move) {
							h.possibleMoves[isWhite] = append(h.possibleMoves[isWhite], move)
						}

						if targetPiece != 0 {
							break
						}

						translatedSquare = translatedSquare.Translate(&offset)
					}
				}

				// Castling
				if piece.Has(King) && !moved && !h.KingInCheck(isWhite) {
					castlingColor := CastlingWhite
					if !isWhite {
						castlingColor = CastlingBlack
					}

					castlingPaths := map[Castling][]Square{
						CastlingQueen: {{rank, 0}, {rank, 1}, {rank, 2}, {rank, 3}},
						CastlingKing:  {{rank, 7}, {rank, 5}, {rank, 6}},
					}

					for castlingType, path := range castlingPaths {
						if h.castling&castlingType&castlingColor != 0 {
							shouldRook := h.GetPieceInSquare(&path[0])
							if !(shouldRook.Has(Rook) && shouldRook.Has(White) == isWhite) {
								h.castling = h.castling &^ (castlingType | castlingColor)
								break
							}

							canCastle := true

							for _, transitSquare := range path[1:] {
								if h.GetPieceInSquare(&transitSquare) != 0 || h.SquareInCheck(&h.board, isWhite, transitSquare) {
									canCastle = false
									break
								}
							}
							if canCastle {
								move := Move{
									playedByWhite: isWhite,
									start:         square,
									end:           path[2],
									castling:      castlingType & castlingColor,
								}
								h.possibleMoves[isWhite] = append(h.possibleMoves[isWhite], move)
							}
						}
					}
				}

			}
		}
	}

}

func (h *Board) MoveIsLegal(move *Move) bool {
	h.MovePieces(move)
	isLegal := !h.KingInCheck(move.playedByWhite)
	h.UndoMovePieces(move)
	return isLegal
}

func (h *Board) MovePieces(move *Move) {
	h.board[move.end.rank][move.end.file] = h.board[move.start.rank][move.start.file]
	h.board[move.start.rank][move.start.file] = 0

	if move.isEnPassant {
		h.board[move.start.rank][move.end.file] = 0
	} else if move.castling != 0 {
		if move.castling&CastlingKing != 0 {
			h.board[move.end.rank][move.end.file-1] = h.board[move.end.rank][7]
			h.board[move.end.rank][7] = 0
		} else if move.castling&CastlingQueen != 0 {
			h.board[move.end.rank][move.end.file+1] = h.board[move.end.rank][0]
			h.board[move.end.rank][0] = 0
		}
	}

	if h.board[move.end.rank][move.end.file].Has(King) {
		h.kingPositions[move.playedByWhite] = Square{move.end.rank, move.end.file}
	}
}

func (h *Board) UndoMovePieces(move *Move) {
	h.board[move.start.rank][move.start.file] = h.board[move.end.rank][move.end.file]

	if move.isEnPassant {
		h.board[move.start.rank][move.end.file] = move.capturedPiece
	} else if move.castling != 0 {
		h.board[move.end.rank][move.end.file] = 0
		// Sposta torre
		if move.castling&CastlingKing != 0 {
			h.board[move.end.rank][7] = h.board[move.end.rank][move.end.file+1]
			h.board[move.end.rank][move.end.file+1] = 0
		} else if move.castling&CastlingQueen != 0 {
			h.board[move.end.rank][0] = h.board[move.end.rank][move.end.file-1]
			h.board[move.end.rank][move.end.file-1] = 0
		}
	} else {
		h.board[move.end.rank][move.end.file] = move.capturedPiece
	}

	if h.board[move.start.rank][move.start.file].Has(King) {
		h.kingPositions[move.playedByWhite] = Square{move.start.rank, move.start.file}
	}
}

func (h *Board) HasPossibleMoves(isWhite bool) bool {
	return len(h.possibleMoves[isWhite]) > 0
}

func (b *Board) Print() {
	fmt.Printf("DEBUG SCACCHIERA\nMosse possibili\nBianco: %d\t\tNero: %d\n", len(b.possibleMoves[true]), len(b.possibleMoves[false]))

	fmt.Print("Partenza mosse bianco: \n")
	for _, move := range b.possibleMoves[true] {
		fmt.Printf("(%d, %d) ", move.start.file, move.start.rank)
	}

	fmt.Printf("Posizione re - Bianco: (%d, %d) - Nero: (%d, %d)\n", b.kingPositions[true].file, b.kingPositions[true].rank, b.kingPositions[false].file, b.kingPositions[false].rank)
	fmt.Printf("Re bianco in scacco: %v\n", b.KingInCheck(true))
	fmt.Printf("Re nero in scacco: %v\n", b.KingInCheck(false))
	fmt.Printf("POSIZIONE\n")
	for rank := len(b.board) - 1; rank >= 0; rank-- {
		for file := 0; file < len(b.board[rank]); file++ {
			piece := b.board[rank][file]
			t := " ·"
			if piece.Has(White) {
				if piece.Has(King) {
					t = " ♔"
				} else if piece.Has(Queen) {
					t = " ♕"
				} else if piece.Has(Bishop) {
					t = " ♗"
				} else if piece.Has(Knight) {
					t = " ♘"
				} else if piece.Has(Rook) {
					t = " ♖"
				} else if piece.Has(Pawn) {
					t = " ♙"
				}
			} else {
				if piece.Has(King) {
					t = " ♚"
				} else if piece.Has(Queen) {
					t = " ♛"
				} else if piece.Has(Bishop) {
					t = " ♝"
				} else if piece.Has(Knight) {
					t = " ♞"
				} else if piece.Has(Rook) {
					t = " ♜"
				} else if piece.Has(Pawn) {
					t = " ♟︎"
				}
			}
			fmt.Print(t)
		}
		fmt.Print("\n")
	}

	fmt.Print("\n")

}

func (b *Board) GenerateFEN(sideToPlay string) (fen string) {
	fen = ""

	for rank := 7; rank >= 0; rank-- {
		skip := 0
		for file := 0; file < 8; file++ {
			piece := b.board[rank][file]
			piece.Clear(Moved)
			if b.board[rank][file] != 0 {
				if skip > 0 {
					fen += strconv.Itoa(skip)
				}
				for c, p := range LETTER_TO_PIECE {
					if p == piece {
						fen += string(c)
					}
				}
				skip = 0
			} else {
				skip++
			}
		}
		if skip > 0 {
			fen += strconv.Itoa(skip)
		}
		fen += "/"
	}

	fen = strings.TrimSuffix(fen, "/")
	fen += " " + sideToPlay[0:1] + " "

	for c, v := range CastlingLookup {
		if b.castling&v != 0 {
			fen += string(c)
		}
	}

	if b.enPassantSquare == nil {
		fen += " -"
	} else {
		fen += " " + b.enPassantSquare.String()
	}

	fen += " " + strconv.Itoa(b.halfmoveClock) + " " + strconv.Itoa(b.fullmoveClock)

	fmt.Println("FEN: " + fen)
	return
}
