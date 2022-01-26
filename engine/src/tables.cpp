#include "tables.hpp"

#include <unordered_map>
#include <random>

#include "bitboard.hpp"

namespace engine {

// Il costruttore crea tutte le lookup tables

Tables::Tables() {
    // Pattern di cattura di pedoni, re e cavalli
    for (int i = 0; i < 64; i++) {
        Bitboard s(1UL << i);

        pawnLookup[s] = s.NorthWest() + s.NorthEast();

        kingLookup[s] = s.North() + s.South() + s.East() + s.West() + s.NorthWest() + s.NorthEast() + s.SouthEast() + s.SouthWest();

        knightLookup[s] = s.North().NorthEast() + s.North().NorthWest() + 
                          s.South().SouthEast() + s.South().SouthWest() + 
                          s.East().NorthEast() + s.East().SouthEast() + 
                          s.West().NorthWest() + s.West().SouthWest();
    }

    // Valore posizionale dei pezzi
    for(size_t p = 0; p < PIECES; p++) 
        for (int i = 0; i < 64; i++)
            pieceValueLookup[p][Bitboard(1UL << i)] = PIECE_SQUARE_TABLES[p][i];

    // RANDOM COSE PER FARE L'HASH ZOBRISK DELLE POSIZIONI
    
    std::mt19937 gen(8101998);
    std::uniform_int_distribution<size_t> distrib(0,SIZE_MAX);

    for(size_t p = 0; p < PIECES; p++)
        for (int i = 0; i < 64; i++)
            randomPieceHash[p][Bitboard(1UL << i)] = distrib(gen);

    for (int i = 0; i < 64; i++)
        randomEnPassantHash[Bitboard(1UL << i)] = distrib(gen);

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

int Tables::GetPieceValue(Piece piece, Bitboard square) {
    return pieceValueLookup[piece][square];
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