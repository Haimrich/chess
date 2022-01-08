#ifndef ENGINE_GLOBAL_H_
#define ENGINE_GLOBAL_H_

#include <string>

namespace engine {

enum Side : size_t {
    WHITE = 0,
    BLACK,

    SIDES
};

enum Castling : size_t {
    KING_CASTLING = 0,
    QUEEN_CASTLING,

    CASTLINGS
};

enum Piece : size_t {
    KING    = 0,
    QUEEN,
    ROOK,
    BISHOP,
    KNIGHT,
    PAWN,

    PIECES,
};


std::string trim(std::string& str, char c);

bool isLower(char c);

bool toUpper(char c);

bool isDigit(char c);

}


#endif