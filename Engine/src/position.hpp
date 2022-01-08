#ifndef ENGINE_POSITION_H_
#define ENGINE_POSITION_H_

#include <string>
#include <vector>

#include "global.hpp"
#include "bitboard.hpp"
#include "move.hpp"

namespace engine {

    class Position {
        private:

            Bitboard bitboards[SIDES][PIECES];

            Bitboard enPassantSquare;

            bool castlingRights[SIDES][CASTLINGS];

            int score;
        
        public:

            Position(std::string fen);

            Position() : Position("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1") {}

            std::vector<Move> GetMoves();

        private:

            bool squareInCheck(Bitboard square);
    };
}
#endif