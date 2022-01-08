#include "tables.hpp"

#include <unordered_map>
#include <random>

#include "bitboard.hpp"

namespace engine {

Tables::Tables() {
    for (Bitboard s(1); !s.IsZero(); s = s.West()) {
        pawnLookup[s] = s.NorthWest() + s.NorthEast();
    }

    for (Bitboard s(1); !s.IsZero(); s = s.West()) {
        kingLookup[s] = s.North() + s.South() + s.East() + s.West() + s.NorthWest() + s.NorthEast() + s.SouthEast() + s.SouthWest();
    }

    for (Bitboard s(1); !s.IsZero(); s = s.West()) {
        knightLookup[s] = s.North().NorthEast() + s.North().NorthWest() + 
                          s.South().SouthEast() + s.South().SouthWest() + 
                          s.East().NorthEast() + s.East().SouthEast() + 
                          s.West().NorthWest() + s.West().SouthWest();
    }

    for(size_t p = 0; p < PIECES; p++) 
        for (int i = 0; i < 64; i++)
            pieceValueLookup[p][Bitboard(1UL << i)] = PIECE_SQUARE_TABLES[p][i];

    // RANDOM COSE PER FARE L'HASH DELLE POSIZIONI
    
    std::mt19937 gen(8101998);
    std::uniform_int_distribution<size_t> distrib(0,SIZE_MAX);

    for(size_t p = 0; p < PIECES; p++)
        for (Bitboard s(1); !s.IsZero(); s = s.West())
            randomPieceHash[p][s] = distrib(gen);

    for (Bitboard s(1); !s.IsZero(); s = s.West())
        randomEnPassantHash[s] = distrib(gen);

    for (size_t c = 0; c < 16; c++)
        randomCastlingHash[Bitboard(c)] = distrib(gen);
}

Bitboard Tables::GetPawnPattern(Bitboard square) {
    return pawnLookup[square];
}

Bitboard Tables::GetKingPattern(Bitboard square) {
    return kingLookup[square];
}

Bitboard Tables::GetKnightPattern(Bitboard square) {
    return knightLookup[square];
}

int Tables::GetPieceValue(Piece piece, Bitboard b) {
    return pieceValueLookup[piece][b];
}

size_t Tables::GetPieceRandom(Piece piece, Bitboard square) {
    return randomPieceHash[piece][square];
}

size_t Tables::GetEnPassantRandom(Bitboard square) {
    return randomEnPassantHash[square];
}

size_t Tables::GetCastlingRandom(Bitboard square) {
    return randomCastlingHash[square];
}

}