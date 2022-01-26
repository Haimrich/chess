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

/*
più sign -> 0 0 0 0 0 0 0 0 a8
            0 0 0 0 0 0 0 0
            0 0 0 0 0 0 0 0
            0 0 0 0 0 0 0 0
            0 0 0 0 0 0 0 0
            0 0 0 0 0 0 0 0
            0 0 0 0 0 0 0 0 a3
            0 0 0 0 0 0 0 0 a2
         h1 0 0 0 0 0 0 0 0 a1 <- meno significativo
*/
        
        U64 bitboard;

        static const U64 notAFile = 0xfefefefefefefefeULL;
        static const U64 notHFile = 0x7f7f7f7f7f7f7f7fULL;
    public:
        // Costruttori
        Bitboard() : bitboard(0ULL) {}
        Bitboard(U64 v) : bitboard(v) {}
        Bitboard(std::string square);

        // Copia e assegnazione
        Bitboard(const Bitboard& b) { 
            bitboard = b.bitboard; 
        }

        Bitboard& operator=(const Bitboard& other) { 
            bitboard = other.bitboard; 
            return *this;
        }

        // Metodi vari
        void Set(size_t pos);
        void Clear(Bitboard b);
        bool Has(Bitboard b) const;
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

        void Flip();

        std::string ToString();

        // Ridefinisco operatore somma, in realtà è un or
        inline Bitboard operator+ (Bitboard const &b) {
            return Bitboard(bitboard | b.bitboard);
        }
        
        // Per stampare bitboard, utile per debug
        friend std::ostream& operator<< (std::ostream& os, const Bitboard& b)
        { 
            os << std::bitset<64>(b.bitboard);
            return os;
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