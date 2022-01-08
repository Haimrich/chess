#ifndef ENGINE_POSITION_H_
#define ENGINE_POSITION_H_

#include <string>
#include <vector>

#include "global.hpp"
#include "bitboard.hpp"
#include "move.hpp"
#include "tables.hpp"

namespace engine {

    class Position {
        private:

            Bitboard bitboards[SIDES][PIECES];

            Bitboard enPassantSquare;

            bool castlingRights[SIDES][CASTLINGS];
        
            bool kingInCheck;

        public:
            int score;

            Position(std::string fen);

            Position() : Position("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1") {}

            std::vector<Move> GetMoves();

            Position Flip();

            int Evaluate(Move m);

            Position MakeMove(Move m);

            std::string MoveToString(Move m);

        private:

            bool squareInCheck(Bitboard square);


        public:

       
        bool operator== (const Position &other) const
        { 
            for (size_t s = 0; s < SIDES; s++)
                for (size_t p = 0; p < PIECES; p++)
                    if (bitboards[SIDES][PIECES] != other.bitboards[SIDES][PIECES])
                        return false;

            if (enPassantSquare != other.enPassantSquare)
                return false;

            for(size_t s = 0; s < SIDES; s++)
                for (size_t c = 0; c < CASTLINGS; c++)
                    if (castlingRights[s][c] != other.castlingRights[s][c])
                        return false;

            return true;
        }

        friend struct PositionHash;
    };


    // per transposition table
    struct PositionHash
    {
        std::size_t operator()(Position k) const {
            Tables& tables = Tables::Instance();

            size_t hash = 0;
            for (size_t s = 0; s < SIDES; s++)
                for (size_t p = 0; p < PIECES; p++)
                    for (auto b : k.bitboards[s][p].Split())
                        hash ^= tables.GetPieceRandom((Piece)p,b);

            hash ^= tables.GetEnPassantRandom(k.enPassantSquare);

            Bitboard castling(0);
            for(size_t s = 0; s < SIDES; s++)
                for (size_t c = 0; c < CASTLINGS; c++)
                    if (k.castlingRights[s][c])
                        castling.Set(s*SIDES+c);

            hash ^= tables.GetCastlingRandom(castling);

            return hash;
        }
    };
}
#endif