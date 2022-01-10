#ifndef ENGINE_EXPLORER_H_
#define ENGINE_EXPLORER_H_

#include <unordered_map>

#include "bitboard.hpp"
#include "position.hpp"

namespace engine {

class Explorer {
    private:

        struct Result {
            int depth;
            int score;
            int gamma;
            Move move;
        };

        int nodes;

        std::unordered_map<Position, Result, PositionHash> transpositionTable;

    public:

        Move Search(Position pos, int maxNodes, double seconds);

    private:

        int Bound(Position pos, int gamma, int depth);

};

}

#endif