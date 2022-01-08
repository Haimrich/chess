#include "directions.hpp"

#include <unordered_map>

#include "bitboard.hpp"

namespace engine {

Directions::Directions() {
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
}

Bitboard Directions::GetPawnPattern(Bitboard square) {
    return pawnLookup[square];
}

Bitboard Directions::GetKingPattern(Bitboard square) {
    return kingLookup[square];
}

Bitboard Directions::GetKnightPattern(Bitboard square) {
    return knightLookup[square];
}

}