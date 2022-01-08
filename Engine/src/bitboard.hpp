#ifndef ENGINE_BITBOARD_H_
#define ENGINE_BITBOARD_H_

#include <string>
#include <vector>
#include <iostream>
#include <bitset>

namespace engine {

class Bitboard {
    private:
        typedef uint_fast64_t U64;
        U64 bitboard;

        static const U64 notAFile = 0xfefefefefefefefe;
        static const U64 notHFile = 0x7f7f7f7f7f7f7f7f;
    public:
        Bitboard() : bitboard(0) {}
        Bitboard(U64 v) : bitboard(v) {}
        Bitboard(std::string square);

        void Set(size_t pos);
        void Clear(Bitboard b);
        bool Has(Bitboard b);
        bool IsZero();
        bool IsRank(size_t rank);

        Bitboard North();
        Bitboard South();
        Bitboard East();
        Bitboard West();

        Bitboard NorthEast();
        Bitboard NorthWest();
        Bitboard SouthEast();
        Bitboard SouthWest();

        Bitboard LS1B();

        Bitboard Intersect(Bitboard b);
        Bitboard Invert(Bitboard b);

        std::vector<Bitboard> Split(Bitboard add);
        std::vector<Bitboard> Split();

        /*
        00000000
        00000000
        00000000
        00000000
        00000000
        00000000
        00000000 a3
        00000000 a2
     h1 00000000 a1
        */

        void Flip();

        std::string ToString();

        inline Bitboard operator+ (Bitboard const &b) {
            return Bitboard(bitboard | b.bitboard);
        }
        
        // per usare bitboard come chiave nelle mappe
        bool operator== (const Bitboard &other) const
        { 
            return bitboard == other.bitboard;
        }

        bool operator!= (const Bitboard &other) const
        { 
            return !(bitboard == other.bitboard);
        }

        friend std::ostream& operator<< (std::ostream& os, const Bitboard& b)
        { 
            os << std::bitset<64>(b.bitboard) << std::endl;
            return os;
        }

        friend struct BitboardHash;
};

// per usare bitboard come chiave nelle mappe
struct BitboardHash
{
  std::size_t operator()(Bitboard k) const {
    return k.bitboard;
  }
};

}

#endif