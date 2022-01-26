#ifndef ENGINE_TABLES_H_
#define ENGINE_TABLES_H_

#include <unordered_map>

#include "global.hpp"
#include "bitboard.hpp"

namespace engine {

// Singleton
// Crea lookup tables alla prima instanziazione
class Tables {
    private:
        Tables();

        std::unordered_map<Bitboard, Bitboard, BitboardHash> pawnLookup;
        std::unordered_map<Bitboard, Bitboard, BitboardHash> kingLookup;
        std::unordered_map<Bitboard, Bitboard, BitboardHash> knightLookup;

        std::unordered_map<Bitboard, int, BitboardHash> pieceValueLookup[PIECES];

        std::unordered_map<Bitboard, size_t, BitboardHash> randomPieceHash[PIECES];
        std::unordered_map<Bitboard, size_t, BitboardHash> randomEnPassantHash;
        std::unordered_map<Bitboard, size_t, BitboardHash> randomCastlingHash;


    public:

        static Tables& Instance()
        {
            static Tables instance;
            return instance;
        }

        // Questi restituiscono caselle di destinazione per un pezzo nella casa square
        Bitboard GetPawnPattern(Bitboard square);
        Bitboard GetKingPattern(Bitboard square);
        Bitboard GetKnightPattern(Bitboard square);

        // restituisce valore di un pezzo in una determinata posizione
        int GetPieceValue(Piece piece, Bitboard b);

        // Restituisce valore random per un pezzo in una casa square, serve per zobrisk hashing
        size_t GetPieceRandom(Piece piece, Bitboard square);
        size_t GetEnPassantRandom(Bitboard square);
        size_t GetCastlingRandom(Bitboard square);
};

}

#endif