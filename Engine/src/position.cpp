#include "position.hpp"

#include <iostream>
#include <sstream>
#include <string>
#include <unordered_map>
#include <vector>
#include <numeric>
#include <functional>
#include <utility>

#include "global.hpp"
#include "bitboard.hpp"
#include "tables.hpp"
#include "move.hpp"

namespace engine {

Position::Position(std::string fen) {
    score = 0;

    std::unordered_map<char, Piece> charToPiece = {
        {'K', Piece::KING},
        {'Q', Piece::QUEEN},
        {'R', Piece::ROOK},
        {'B', Piece::BISHOP},
        {'N', Piece::KNIGHT},
        {'P', Piece::PAWN},
    };

    for (size_t s = 0; s < SIDES; s++)
        for (size_t p = 0; p < PIECES; p++)
            bitboards[s][p] = Bitboard();

    fen = trim(fen,'[');
    fen = trim(fen,']');

    std::istringstream iss(fen);
    
	// Campo 1 - Posizione

    std::string position;
    getline(iss, position, ' ');

	size_t rank = 7;
    size_t file = 0;
    for (char c : position) {
        if (isDigit(c)) {
            file += c - '0';
        } else if (c == '/') {
            rank--;
            file = 0;
        } else {
            Piece p = charToPiece[toUpper(c)];
            Side s = isLower(c) ? BLACK : WHITE;
            bitboards[s][p].Set(rank*file);
            file++;
        }
	}

    // Campo 2 - Turno giocatore

	std::string playerToMove;
    getline(iss, playerToMove, ' ');

	// Campo 3 - Possibilità di arrocco
    for (auto& s : castlingRights) for (bool& r : s) r = false;

    std::string castlingOpp;
    getline(iss, castlingOpp, ' ');

	for (char c : castlingOpp) {
		switch (c) {
            case 'K': castlingRights[WHITE][KING_CASTLING] = true; break;
            case 'Q': castlingRights[WHITE][QUEEN_CASTLING] = true; break;
            case 'k': castlingRights[BLACK][KING_CASTLING] = true; break;
            case 'q': castlingRights[BLACK][QUEEN_CASTLING] = true; break;
        }
	}

	// Campo 4 - En Passant

    std::string enPassant;
    getline(iss, enPassant, ' ');

    enPassantSquare = 0;
	if (enPassant.compare("-") != 0) {
		enPassantSquare = Bitboard(enPassant);
	}

	// Campo 5 e 6 - Halfmove e fullmove clock
	// TODO


    // Flip se non è il turno del bianco
    if (playerToMove.compare("w") == 0) {
		// TODO
	} else {
		std::swap(bitboards[WHITE], bitboards[BLACK]);
        std::swap(castlingRights[WHITE], castlingRights[BLACK]);
        enPassantSquare.Flip();

        for (auto& bs : bitboards)
            for (auto& b : bs)
                b.Flip();

        score = -score;
	}
}


std::vector<Move> Position::GetMoves() {
    std::vector<Move> moves;

    Bitboard ourPieces = std::accumulate(std::begin(bitboards[WHITE]), std::end(bitboards[WHITE]), Bitboard(0));

    // PEDONI
    Bitboard pawns = bitboards[WHITE][PAWN];
    Bitboard pawnCaptureSquares = enPassantSquare + std::accumulate(std::begin(bitboards[BLACK]), std::end(bitboards[BLACK]), Bitboard(0));
    Bitboard occupiedSquares = pawnCaptureSquares + ourPieces;
    for (Bitboard p = pawns.LS1B(); !pawns.IsZero(); pawns.Clear(p)) {
        // Pedoni - Movimento
        Bitboard up = p.North();
        if (!occupiedSquares.Has(up)) {
            moves.emplace_back(up+p, PAWN);
            up = up.North();
            if (p.IsRank(1) && !occupiedSquares.Has(up)) {
                moves.emplace_back(up+p, PAWN);
            }
        }
        // Pedoni - Cattura
        Bitboard possibleDestinations = Tables::Instance().GetPawnPattern(p);
        possibleDestinations.Has(pawnCaptureSquares);
        auto possibleMoves = possibleDestinations.Split(p);
        for (auto m : possibleMoves)
            moves.emplace_back(m, PAWN);
    }

    // RE

    {
        Bitboard king = bitboards[WHITE][KING];
        Bitboard possibleDestinations = Tables::Instance().GetKingPattern(king);
        possibleDestinations.Clear(ourPieces);
        auto possibleMoves = possibleDestinations.Split(king);
        for (auto m : possibleMoves)
            moves.emplace_back(m, KING);

        // TODO Castling
        if (!kingInCheck) 
        {
            if (castlingRights[WHITE][KING_CASTLING]) 
            {
                Bitboard kCastlingPathMask(0b01100000);
                if (!kCastlingPathMask.Has(occupiedSquares)) {
                    Bitboard pathSquare = bitboards[WHITE][KING];
                    bool pathInCheck = false;
                    for (size_t i = 0; i < 2; i++) {
                        pathSquare = pathSquare.East();
                        if (squareInCheck(pathSquare)) {
                            pathInCheck = true;
                            break;
                        }
                    }
                    if (!pathInCheck) {
                        moves.emplace_back(pathSquare + bitboards[WHITE][KING], KING);
                    }
                }
            }
            if (castlingRights[WHITE][QUEEN_CASTLING]) 
            {
                Bitboard qCastlingPathMask(0b00001110);
                if (!qCastlingPathMask.Has(occupiedSquares)) {
                    Bitboard pathSquare = bitboards[WHITE][KING];
                    bool pathInCheck = false;
                    for (size_t i = 0; i < 3; i++) {
                        pathSquare = pathSquare.West();
                        if (squareInCheck(pathSquare)) {
                            pathInCheck = true;
                            break;
                        }
                    }
                    if (!pathInCheck) {
                        moves.emplace_back(pathSquare.East() + bitboards[WHITE][KING], KING);
                    }
                }
            }
        }
    }

    // CAVALLI

    Bitboard knights = bitboards[WHITE][KNIGHT];
    for (Bitboard k = knights.LS1B(); !knights.IsZero(); knights.Clear(k)) {
        Bitboard possibleDestinations = Tables::Instance().GetKnightPattern(k);
        possibleDestinations.Clear(ourPieces);
        auto possibleMoves = possibleDestinations.Split(k);
        for (auto m : possibleMoves)
            moves.emplace_back(m, KNIGHT);
    }

    // Raggi
    // Con questa cosa faccio i loop e non riscrivo 10 volte la stessa cosa
    std::function<Bitboard(Bitboard&)> dirsFunc[] = {
        &Bitboard::North, &Bitboard::East, &Bitboard::South, &Bitboard::West,
        &Bitboard::NorthWest, &Bitboard::NorthEast, &Bitboard::SouthWest, &Bitboard::SouthEast
    };

    // Regina, Torri e alfiere

    Bitboard queens = bitboards[WHITE][QUEEN];
    for (Bitboard q = queens.LS1B(); !queens.IsZero(); queens.Clear(q)) 
        for (size_t d = 0; d < 8; d++) 
            for (Bitboard s = dirsFunc[d](q); !s.IsZero() && !ourPieces.Has(s); s = dirsFunc[d](s)) 
                moves.emplace_back(s + q, QUEEN);

    Bitboard rooks = bitboards[WHITE][ROOK];
    for (Bitboard r = rooks.LS1B(); !rooks.IsZero(); rooks.Clear(r)) 
        for (size_t d = 0; d < 4; d++) 
            for (Bitboard s = dirsFunc[d](r); !s.IsZero() && !ourPieces.Has(s); s = dirsFunc[d](s)) 
                moves.emplace_back(s + r, ROOK);

    Bitboard bishops = bitboards[WHITE][BISHOP];
    for (Bitboard b = bishops.LS1B(); !bishops.IsZero(); bishops.Clear(b)) 
        for (size_t d = 0; d < 8; d++) 
            for (Bitboard s = dirsFunc[d](b); !s.IsZero() && !ourPieces.Has(s); s = dirsFunc[d](s)) 
                moves.emplace_back(s + b, BISHOP);


    return moves;
}


bool Position::squareInCheck(Bitboard square) {
    Bitboard pawn = Tables::Instance().GetPawnPattern(square);
    Bitboard enemyPawns = bitboards[BLACK][PAWN];
    if (enemyPawns.Has(pawn)) return true;

    Bitboard knight = Tables::Instance().GetKnightPattern(square);
    Bitboard enemyKnights = bitboards[BLACK][KNIGHT];
    if (enemyKnights.Has(knight)) return true;

    Bitboard king = Tables::Instance().GetKingPattern(square);
    Bitboard enemyKing = bitboards[BLACK][KING];
    if (enemyKing.Has(king)) return true;

    // Proiettare raggi
    Bitboard ourPieces = std::accumulate(std::begin(bitboards[WHITE]), std::end(bitboards[WHITE]), Bitboard(0));
    Bitboard enemyRookAndQueen = bitboards[BLACK][ROOK] + bitboards[BLACK][QUEEN];
    Bitboard enemyBishopAndQueen = bitboards[BLACK][BISHOP] + bitboards[BLACK][QUEEN];

    for (Bitboard s = square.North(); !s.IsZero() && !ourPieces.Has(s); s = s.North())
        if (enemyRookAndQueen.Has(s)) return true;

    for (Bitboard s = square.South(); !s.IsZero() && !ourPieces.Has(s); s = s.South())
        if (enemyRookAndQueen.Has(s)) return true;
    
    for (Bitboard s = square.East(); !s.IsZero() && !ourPieces.Has(s); s = s.East())
        if (enemyRookAndQueen.Has(s)) return true;

    for (Bitboard s = square.West(); !s.IsZero() && !ourPieces.Has(s); s = s.West())
        if (enemyRookAndQueen.Has(s)) return true;

    for (Bitboard s = square.NorthEast(); !s.IsZero() && !ourPieces.Has(s); s = s.NorthEast())
        if (enemyBishopAndQueen.Has(s)) return true;

    for (Bitboard s = square.NorthWest(); !s.IsZero() && !ourPieces.Has(s); s = s.NorthWest())
        if (enemyBishopAndQueen.Has(s)) return true;

    for (Bitboard s = square.SouthEast(); !s.IsZero() && !ourPieces.Has(s); s = s.SouthEast())
        if (enemyBishopAndQueen.Has(s)) return true;

    for (Bitboard s = square.SouthWest(); !s.IsZero() && !ourPieces.Has(s); s = s.SouthWest())
        if (enemyBishopAndQueen.Has(s)) return true;


    return false;
}

Position Position::Flip() {
    Position np = *this;
    
    std::swap(np.bitboards[WHITE], np.bitboards[BLACK]);
    std::swap(np.castlingRights[WHITE], np.castlingRights[BLACK]);

    np.enPassantSquare.Flip();

    for (auto& bs : np.bitboards)
        for (auto& b : bs)
            b.Flip();

    np.score = -np.score;

    return np;
}

int Position::Evaluate(Move m) {
    Bitboard startPiece = bitboards[WHITE][m.piece].Intersect(m.bitboard);
    Bitboard endPiece = startPiece.Invert(m.bitboard);

    // Togli valore pezzo in pos vecchio e aggiungi quello in pos nuova
    Tables& tables = Tables::Instance();
    score = tables.GetPieceValue(m.piece, endPiece) - tables.GetPieceValue(m.piece, startPiece); 

    // Se cattura aggiungi pure valore del pezzo catturato
    for (size_t p = 0; p < PIECES; p++)
        if (bitboards[BLACK][p].Has(endPiece)) 
            score += tables.GetPieceValue((Piece)p, endPiece);

    // Castling
	if (startPiece == bitboards[WHITE][KING] && (m.bitboard == Bitboard(0b00001010) || m.bitboard == Bitboard(0b00101000) ) ) {
        if (m.bitboard == Bitboard(0b00001010)) {
            score += tables.GetPieceValue(ROOK, Bitboard(0b000000100));
            score -= tables.GetPieceValue(ROOK, Bitboard(0b000000001));
        } else if (m.bitboard == Bitboard(0b00101000)) {
            score += tables.GetPieceValue(ROOK, Bitboard(0b000100000));
            score -= tables.GetPieceValue(ROOK, Bitboard(0b100000000));
        }
	}

    // Pedoni
    if (bitboards[WHITE][PAWN].Has(startPiece)) {
        // Ultimo rank -> Promozione
        if (Bitboard(0xFF00000000000000).Has(endPiece)) { 
            score += tables.GetPieceValue(QUEEN, endPiece) - tables.GetPieceValue(PAWN, startPiece);
        }
        // EnPassant
        if(endPiece == enPassantSquare) {
            score += tables.GetPieceValue(PAWN, endPiece.South());
        }
    }

    return score;
}

Position Position::MakeMove(Move m) {
    Bitboard startPiece = bitboards[WHITE][m.piece].Intersect(m.bitboard);
    Bitboard endPiece = startPiece.Invert(m.bitboard);

    Position np = *this;
    np.enPassantSquare = Bitboard(0);
	np.score = score + Evaluate(m);

    for (size_t p = 0; p < PIECES; p++)
        if (p == m.piece)
            np.bitboards[WHITE][p].Invert(m.bitboard);
        else 
            np.bitboards[WHITE][p].Clear(m.bitboard);
    
    if (startPiece.Has(0b00000001)) // Torre lato re mossa
        np.castlingRights[WHITE][KING_CASTLING] = false;
	else if (startPiece.Has(0b10000000)) // Torre lato regina mossa
        np.castlingRights[WHITE][QUEEN_CASTLING] = false;
    
    if (endPiece.Has(1UL << 63))
        np.castlingRights[BLACK][QUEEN_CASTLING] = false;
	else if (endPiece.Has(1UL << (63-8)))
        np.castlingRights[BLACK][KING_CASTLING] = false;

    if (startPiece == bitboards[WHITE][KING]) {
		np.castlingRights[WHITE][QUEEN_CASTLING] = false;
		np.castlingRights[WHITE][KING_CASTLING] = false;

        if (m.bitboard == Bitboard(0b00001010)) 
            np.bitboards[WHITE][ROOK].Invert(Bitboard(0b000000101));
        
        if (m.bitboard == Bitboard(0b00101000)) {
            np.bitboards[WHITE][ROOK].Invert(Bitboard(0b100100000));
        }
	}
    
    if (m.piece == PAWN) {
		// Promozione perchè arrivato in ultima traversa
        if (Bitboard(0xFF00000000000000).Has(endPiece)) { 
            np.bitboards[WHITE][QUEEN].Invert(endPiece);
            np.bitboards[WHITE][PAWN].Clear(endPiece);
        // Aggiornare en passant perchè doppio passo
        } else if (Bitboard(0x00000000FF000000).Has(endPiece)) {
            np.enPassantSquare = endPiece.South();
        }
        // Cattura enpassant
        if (endPiece == enPassantSquare) {
            Bitboard enPassantTarget = endPiece.South();
            for (size_t p = 0; p < PIECES; p++)
                if (np.bitboards[BLACK][p].Has(enPassantTarget)) 
                    np.bitboards[BLACK][p].Invert(enPassantTarget);
        }
	}

	return np.Flip();
}


std::string Position::MoveToString(Move m) {
    Bitboard startPiece = bitboards[WHITE][m.piece].Intersect(m.bitboard);
    Bitboard endPiece = startPiece.Invert(m.bitboard);

    std::cout << startPiece;
    std::cout << endPiece;
    
    return startPiece.ToString() + endPiece.ToString();
}

}