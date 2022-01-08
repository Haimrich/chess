#include "position.hpp"

#include <iostream>
#include <sstream>
#include <string>
#include <unordered_map>
#include <vector>
#include <numeric>

#include "global.hpp"
#include "bitboard.hpp"
#include "directions.hpp"
#include "move.hpp"

namespace engine {

Position::Position(std::string fen) {
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

	if (playerToMove.compare("w") == 0) {
		// TODO
	} else {
		// TODO FLIP
	}

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
}


std::vector<Move> Position::GetMoves() {
    std::vector<Move> moves;

    Bitboard ourPieces = std::accumulate(std::begin(bitboards[WHITE]), std::end(bitboards[WHITE]), Bitboard(0));

    Bitboard knights = bitboards[WHITE][KNIGHT];
    for (Bitboard k = knights.LS1B(); !knights.IsZero(); knights.Clear(k)) {
        Bitboard possibleDestinations = Directions::Instance().GetKnightPattern(k);
        possibleDestinations.Clear(ourPieces);
        auto possibleMoves = possibleDestinations.Split(k);
        for (auto m : possibleMoves)
            moves.emplace_back(m, KNIGHT);
    }

    {
        Bitboard king = bitboards[WHITE][KING];
        Bitboard possibleDestinations = Directions::Instance().GetKingPattern(king);
        possibleDestinations.Clear(ourPieces);
        auto possibleMoves = possibleDestinations.Split(king);
        for (auto m : possibleMoves)
            moves.emplace_back(m, KING);

        // TODO Castling
    }



    return moves;
}


bool Position::squareInCheck(Bitboard square) {
    Bitboard pawn = Directions::Instance().GetPawnPattern(square);
    Bitboard enemyPawns = bitboards[BLACK][PAWN];
    if (enemyPawns.Has(pawn)) return true;

    Bitboard knight = Directions::Instance().GetKnightPattern(square);
    Bitboard enemyKnights = bitboards[BLACK][KNIGHT];
    if (enemyKnights.Has(knight)) return true;

    Bitboard king = Directions::Instance().GetKingPattern(square);
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

}