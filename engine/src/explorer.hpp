#ifndef ENGINE_EXPLORER_H_
#define ENGINE_EXPLORER_H_

#include <unordered_map>

#include "bitboard.hpp"
#include "position.hpp"

namespace engine {

class Explorer {
    private:

        struct Entry {
            int depth;
            int score;
            int gamma;
            Move move;

            Entry(int d, int s, int g, Move m) : depth(d), score(s), gamma(g), move(m) {}
            Entry() : Entry(0,0,0,Move()) {}
        };

        int nodes;

        std::unordered_map<Position, Entry, PositionHash> transpositionTable;

    public:

        Move Search(Position pos, int maxNodes, double seconds);

    private:

        int Bound(Position pos, int gamma, int depth);

};

}

#endif