#ifndef ENGINE_DIRECTIONS_H_
#define ENGINE_DIRECTIONS_H_

#include <unordered_map>

#include "bitboard.hpp"

namespace engine {

// Singleton
// Crea lookup tables alla prima instanziazione
class Directions {
    private:
        Directions();

        std::unordered_map<Bitboard, Bitboard, BitboardHash> pawnLookup;
        std::unordered_map<Bitboard, Bitboard, BitboardHash> kingLookup;
        std::unordered_map<Bitboard, Bitboard, BitboardHash> knightLookup;

    public:

        static Directions& Instance()
        {
            static Directions instance;
            return instance;
        }

        Bitboard GetPawnPattern(Bitboard square);
        Bitboard GetKingPattern(Bitboard square);
        Bitboard GetKnightPattern(Bitboard square);
};

}

#endif