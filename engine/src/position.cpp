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

// Costruttore con notazione di Forsyth-Edwards per decodificare posizione ricevuta dal backend
Position::Position(std::string fen) : enPassantSquare(0), score(0) {

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
            bitboards[s][p].Set(rank*8+file);
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

	if (enPassant.compare("-") != 0) {
		enPassantSquare = Bitboard(enPassant);
	}

	// Campo 5 e 6 - Halfmove e fullmove clock
	// TODO ma non credo servano


    // Flip se non è il turno del bianco
    if (playerToMove.compare("b") == 0) {
		std::swap(bitboards[WHITE], bitboards[BLACK]);
        for (size_t s = 0; s < SIDES; s++)
            for (size_t p = 0; p < PIECES; p++)
                bitboards[s][p].Flip();

        std::swap(castlingRights[WHITE], castlingRights[BLACK]);
        enPassantSquare.Flip();

        score = -score;
	}
}

// Costruttore di copia
Position::Position(const Position &p) {
    bitboards = p.bitboards;
    enPassantSquare = p.enPassantSquare;
    castlingRights = p.castlingRights;
    score = p.score;
}

// Restituisce vettore di tutte le mosse possibili in questa posizione per il bianco
std::vector<Move> Position::GetMoves() {
    std::vector<Move> moves;

    Bitboard ourPieces = std::accumulate(std::begin(bitboards[WHITE]), std::end(bitboards[WHITE]), Bitboard(0));

    // PEDONI
    Bitboard pawns = bitboards[WHITE][PAWN];
    Bitboard pawnCaptureSquares = enPassantSquare + std::accumulate(std::begin(bitboards[BLACK]), std::end(bitboards[BLACK]), Bitboard(0));
    Bitboard occupiedSquares = pawnCaptureSquares + ourPieces;

    for (Bitboard p = pawns.LS1B(); !pawns.IsZero(); pawns.Clear(p), p = pawns.LS1B()) {
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
        auto possibleMoves = possibleDestinations.Intersect(pawnCaptureSquares).Split(p);
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

        // Castling - TODO qua ci sarà qualche prob visto che non arrocca mai
        if (castlingRights[WHITE][KING_CASTLING] && !squareInCheck(bitboards[WHITE][KING])) 
        {
            Bitboard kCastlingPathMask(0b01100000ULL);
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
        if (castlingRights[WHITE][QUEEN_CASTLING] && !squareInCheck(bitboards[WHITE][KING])) 
        {
            Bitboard qCastlingPathMask(0b00001110ULL);
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

    // CAVALLI

    Bitboard knights = bitboards[WHITE][KNIGHT];
    for (Bitboard k = knights.LS1B(); !knights.IsZero(); knights.Clear(k), k = knights.LS1B() ) {
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
    for (Bitboard q = queens.LS1B(); !queens.IsZero(); queens.Clear(q), q = queens.LS1B() ) 
        for (size_t d = 0; d < 8; d++) 
            for (Bitboard s = dirsFunc[d](q); !s.IsZero() && !ourPieces.Has(s); s = dirsFunc[d](s)) {
                moves.emplace_back(s + q, QUEEN);
                if (occupiedSquares.Has(s)) break;
            }

    Bitboard rooks = bitboards[WHITE][ROOK];
    for (Bitboard r = rooks.LS1B(); !rooks.IsZero(); rooks.Clear(r), r = rooks.LS1B() ) 
        for (size_t d = 0; d < 4; d++) 
            for (Bitboard s = dirsFunc[d](r); !s.IsZero() && !ourPieces.Has(s); s = dirsFunc[d](s)) {
                moves.emplace_back(s + r, ROOK);
                if (occupiedSquares.Has(s)) break;
            }


    Bitboard bishops = bitboards[WHITE][BISHOP];
    for (Bitboard b = bishops.LS1B(); !bishops.IsZero(); bishops.Clear(b), b = bishops.LS1B()) 
        for (size_t d = 4; d < 8; d++) 
            for (Bitboard s = dirsFunc[d](b); !s.IsZero() && !ourPieces.Has(s); s = dirsFunc[d](s)) {
                moves.emplace_back(s + b, BISHOP);
                if (occupiedSquares.Has(s)) break;
            }


    //std::cout << "Mosse possibili: " << moves.size()  << std::endl;
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
    Bitboard enemyPieces = std::accumulate(std::begin(bitboards[BLACK]), std::end(bitboards[BLACK]), Bitboard(0));
    Bitboard enemyRookAndQueen = bitboards[BLACK][ROOK] + bitboards[BLACK][QUEEN];
    Bitboard enemyBishopAndQueen = bitboards[BLACK][BISHOP] + bitboards[BLACK][QUEEN];

    for (Bitboard s = square.North(); !s.IsZero() && !ourPieces.Has(s); s = s.North()) {
        if (enemyRookAndQueen.Has(s)) return true;
        if (enemyPieces.Has(s)) break;
    }

    for (Bitboard s = square.South(); !s.IsZero() && !ourPieces.Has(s); s = s.South()) {
        if (enemyRookAndQueen.Has(s)) return true;
        if (enemyPieces.Has(s)) break;
    }
    
    for (Bitboard s = square.East(); !s.IsZero() && !ourPieces.Has(s); s = s.East()) {
        if (enemyRookAndQueen.Has(s)) return true;
        if (enemyPieces.Has(s)) break;
    }

    for (Bitboard s = square.West(); !s.IsZero() && !ourPieces.Has(s); s = s.West()) {
        if (enemyRookAndQueen.Has(s)) return true;
        if (enemyPieces.Has(s)) break;
    }

    for (Bitboard s = square.NorthEast(); !s.IsZero() && !ourPieces.Has(s); s = s.NorthEast()) {
        if (enemyBishopAndQueen.Has(s)) return true;
        if (enemyPieces.Has(s)) break;
    }

    for (Bitboard s = square.NorthWest(); !s.IsZero() && !ourPieces.Has(s); s = s.NorthWest()) {
        if (enemyBishopAndQueen.Has(s)) return true;
        if (enemyPieces.Has(s)) break;
    }

    for (Bitboard s = square.SouthEast(); !s.IsZero() && !ourPieces.Has(s); s = s.SouthEast()) {
        if (enemyBishopAndQueen.Has(s)) return true;
        if (enemyPieces.Has(s)) break;
    }

    for (Bitboard s = square.SouthWest(); !s.IsZero() && !ourPieces.Has(s); s = s.SouthWest()) {
        if (enemyBishopAndQueen.Has(s)) return true;
        if (enemyPieces.Has(s)) break;
    }

    return false;
}

// Scambia la posizione del bianco e del nero e guarda la scacchiera dal lato opposto
Position Position::Flip() {
    Position np = *this;
    
    std::swap(np.bitboards[WHITE], np.bitboards[BLACK]);
    for (size_t s = 0; s < SIDES; s++)
        for (size_t p = 0; p < PIECES; p++)
            np.bitboards[s][p].Flip();

    std::swap(np.castlingRights[WHITE], np.castlingRights[BLACK]);

    np.enPassantSquare.Flip();

    np.score = -score;

    return np;
}

// Restituisce punteggio della posizione dopo aver fatto la mossa m
// Non modifica la posizione attuale
int Position::Evaluate(Move m) {
    Bitboard startPiece = bitboards[WHITE][m.piece].Intersect(m.bitboard);
    Bitboard endPiece = startPiece.Invert(m.bitboard);

    // Togli valore pezzo in pos vecchio e aggiungi quello in pos nuova
    Tables& tables = Tables::Instance();
    int newScore = tables.GetPieceValue(m.piece, endPiece) - tables.GetPieceValue(m.piece, startPiece); 

    // Se cattura aggiungi pure valore del pezzo catturato dal punto di vista dell'avversario
    for (size_t p = 0; p < PIECES; p++)
        if (bitboards[BLACK][p].Has(endPiece)) {
            Bitboard flippedEndPiece = endPiece;
            flippedEndPiece.Flip();
            newScore += tables.GetPieceValue((Piece)p, flippedEndPiece);
        }

    // Castling
	if (startPiece == bitboards[WHITE][KING]) {
        if (m.bitboard == Bitboard(0b01010000ULL) ) { // Arrocco lato re
            newScore += tables.GetPieceValue(ROOK, Bitboard(0b00100000ULL));
            newScore -= tables.GetPieceValue(ROOK, Bitboard(0b10000000ULL));
        } else if (m.bitboard == Bitboard(0b00010100ULL) ) { // Arrocco lato regina
            newScore += tables.GetPieceValue(ROOK, Bitboard(0b00001000ULL));
            newScore -= tables.GetPieceValue(ROOK, Bitboard(0b00000001ULL));
        }
	}

    // Pedoni
    if (bitboards[WHITE][PAWN].Has(startPiece)) {
        // Ultimo rank -> Promozione
        if (endPiece.IsRank(7)) { 
            newScore += tables.GetPieceValue(QUEEN, endPiece) - tables.GetPieceValue(PAWN, startPiece);
        }
        // EnPassant
        if(endPiece == enPassantSquare) {
            newScore += tables.GetPieceValue(PAWN, endPiece.South());
        }
    }

    return newScore;
}

// Restituisce la posizione ottenuta giocando la mossa m e flippando tutto
Position Position::MakeMove(Move m) {
    Bitboard startPiece = bitboards[WHITE][m.piece].Intersect(m.bitboard);
    Bitboard endPiece = startPiece.Invert(m.bitboard);

    Position np = *this;
    np.enPassantSquare = Bitboard(0);
	np.score += np.Evaluate(m);

    for (size_t p = 0; p < PIECES; p++) {    
        // Eventuali catture
        for (size_t s = 0; s < SIDES; s++)
            np.bitboards[s][p].Clear(m.bitboard);

        // Sposta il pezzo
        if (p == m.piece)
            np.bitboards[WHITE][p] = np.bitboards[WHITE][p].Invert(endPiece);
    }
    
    if (startPiece.Has(0b00000001ULL)) // Torre lato regina mossa
        np.castlingRights[WHITE][QUEEN_CASTLING] = false;
	else if (startPiece.Has(0b10000000ULL)) // Torre lato re mossa
        np.castlingRights[WHITE][KING_CASTLING] = false;
    
    if (endPiece.Has(1UL << 63)) // Torre lato re mangiata
        np.castlingRights[BLACK][KING_CASTLING] = false;
	else if (endPiece.Has(1UL << (63-8))) // Torre lato regina mangiata
        np.castlingRights[BLACK][QUEEN_CASTLING] = false;

    if (startPiece == bitboards[WHITE][KING]) {
		np.castlingRights[WHITE][QUEEN_CASTLING] = false;
		np.castlingRights[WHITE][KING_CASTLING] = false;

        if (m.bitboard == Bitboard(0b00010100ULL))  // Arrocco lato regina bianco
            np.bitboards[WHITE][ROOK] = np.bitboards[WHITE][ROOK].Invert(Bitboard(0b00001001ULL));
        
        if (m.bitboard == Bitboard(0b01010000ULL)) // Arrocco lato re bianco
            np.bitboards[WHITE][ROOK] = np.bitboards[WHITE][ROOK].Invert(Bitboard(0b10100000ULL));
	}
    
    if (m.piece == PAWN) {
		// Promozione se arrivato in ultima traversa
        if (endPiece.IsRank(7)) { 
            np.bitboards[WHITE][QUEEN].Invert(endPiece);
            np.bitboards[WHITE][PAWN].Clear(endPiece);
        // Aggiornare en passant se doppio passo
        } else if (endPiece.IsRank(3) && startPiece.IsRank(1)) {
            np.enPassantSquare = endPiece.South();
        // Cattura enpassant  
        } else if (endPiece == enPassantSquare) {
            Bitboard enPassantTarget = endPiece.South();
            np.bitboards[BLACK][PAWN].Clear(enPassantTarget);
        }
	}

	return np.Flip();
}

// Restituisce notazione UCI di una mossa tipo e2e4 cioè (casa partenza)(casa destinazione)
std::string Position::MoveToString(Move m) {
    Bitboard startPiece = bitboards[WHITE][m.piece].Intersect(m.bitboard);
    Bitboard endPiece = startPiece.Invert(m.bitboard);

   // std::cout << startPiece << std::endl;
   // std::cout << endPiece << std::endl;

    return startPiece.ToString() + endPiece.ToString();
}

// Operatore assegnazione
Position& Position::operator=(const Position& p) { 
    bitboards = p.bitboards;
    enPassantSquare = p.enPassantSquare;
    castlingRights = p.castlingRights;
    score = p.score;

    return *this;
}

// Per debug
std::ostream& operator<< (std::ostream& os, const Position& p) {
    os << "POSITION - Hash: " << PositionHash()(p) << std::endl;
    std::string pieceToSymbol[] = {" ♚", " ♛", " ♜", " ♝", " ♞", " ♟︎", " ♔", " ♕", " ♖", " ♗", " ♘", " ♙"};

    for (int rank = 7; rank >= 0; rank--) {
        for (int file = 0; file < 8; file++) {
            bool pieceFound = false;
            for (size_t side = 0; side < SIDES; side++)
                for (size_t piece = 0; piece < PIECES; piece++) 
                    if (p.bitboards[side][piece].Has(Bitboard(1UL << (rank*8+file) ))) {
                        os << pieceToSymbol[side*PIECES+piece];
                        pieceFound = true;
                    }
            if (!pieceFound) os << " ·";
       }
       os << std::endl;
    }
    return os;
}

}