#ifndef ENGINE_MOVE_H_
#define ENGINE_MOVE_H_

#include "global.hpp"
#include "bitboard.hpp"

namespace engine {

class Move
{
    public:
    Bitboard bitboard;
    Piece piece;

    Move(Bitboard bitboard, Piece piece) : bitboard(bitboard), piece(piece) {}
    Move() : bitboard(0), piece(KING) {}

    Move(const Move& other) : bitboard(other.bitboard), piece(other.piece) {}

    Move& operator=(const Move& other) {
        bitboard = other.bitboard;
        piece = other.piece;
        return *this;
    }


};

}


#endif