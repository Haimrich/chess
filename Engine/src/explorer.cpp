#include "explorer.hpp"

#include <chrono>
#include <iostream>

#include "global.hpp"

namespace engine {

// References
// http://people.csail.mit.edu/plaat/mtdf.html
// https://zserge.com/posts/carnatus/
// https://github.com/thomasahle/sunfish/blob/master/sunfish.py


Move Explorer::Search(Position pos, int maxNodes, double timeBudget) {
    std::cout << "Nuova ricerca avviata..." << std::endl;
    using namespace std::chrono;

    auto startTime = high_resolution_clock::now();
    nodes = 0;

    for (int depth = 1; depth < MAX_DEPTH; depth++) {
        int lower = - MATE_VALUE * 3;
        int upper = + MATE_VALUE * 3;

        int score = 0;

        while (lower < upper - EVAL_ROUGHNESS) {
            int gamma = (lower + upper + 1) / 2;
            score = Bound(pos, gamma, depth);
            std::cout << "Gamma: " << gamma << std::endl;

            if (score >= gamma) 
                lower = score;
            if (score < gamma)
                upper = score;
        }
        
        auto currentTime = high_resolution_clock::now();
        double elapsedTime = duration_cast<seconds>(currentTime - startTime).count();

        if (std::abs(score) >= MATE_VALUE || nodes >= maxNodes || elapsedTime >= timeBudget) break;
    }

    std::cout << "Dim: " << transpositionTable.size() << " - Trovato: " << (transpositionTable.find(pos) != transpositionTable.end()) << std::endl;
    /*
    for (auto m : transpositionTable) {
        Position pos = m.first;
        std::cout << "Hash: " << PositionHash()(pos) << " - Mossa: " << pos.MoveToString(m.second.move) << std::endl;
    }
    */
    return transpositionTable[pos].move;
}



int Explorer::Bound(Position pos, int gamma, int depth) {
    nodes++;
    //std::cout << nodes << std::endl;
    //std::cout << pos << std::endl;

    auto tpEntry = transpositionTable.find(pos);
    bool foundPositionInTable = (tpEntry != transpositionTable.end());
    bool foundShallower = true;

    if (foundPositionInTable) {
        Result r = tpEntry->second;
        foundShallower = depth >= r.depth;
        if (r.depth >= depth && ( (r.score < r.gamma && r.score < gamma) || (r.score >= r.gamma && r.score >= gamma) ))
            return r.score;
    }

    if (std::abs(pos.score) >= MATE_VALUE) return pos.score;

    int nullScore = pos.score;

    if (depth > 0) nullScore = - Bound(pos.Flip(), 1-gamma, depth-3);

    if (nullScore >= gamma) return nullScore;

    int bestScore = -3 * MATE_VALUE;
    Move bestMove;

    for (auto& m : pos.GetMoves()) {
        if (depth <= 0 && pos.Evaluate(m) < 150) break;

        int score = - Bound(pos.MakeMove(m), 1-gamma, depth-1);

        if (score > bestScore) {
            bestScore = score;
            bestMove = m;
        }

        if (score >= gamma) break;
    }

    if (depth <= 0 && bestScore < nullScore) return nullScore;

    // Stalemate check: best move loses king + null move is better
    if (depth > 0 && bestScore <= -MATE_VALUE && nullScore > -MATE_VALUE) bestScore = 0;

    if (!foundPositionInTable || (foundShallower && bestScore >= gamma)) {
        transpositionTable[pos] = {.depth= depth, .score= bestScore, .gamma= gamma, .move= bestMove};

        //if (transpositionTable.size() > MAX_TRANSPOSITION_TABLE_SIZE)
        //    transpositionTable.clear();
    }

    return bestScore; 
}

}