#include "bitboard.hpp"

#include <string>
#include <vector>

#include <iostream>

namespace engine {

    Bitboard::Bitboard(std::string square) {
        size_t file = std::string("abcdefgh").find(square.at(0));
        size_t rank = square.at(1) - '1';

        bitboard = 1UL << (rank*8+file);
    }

    void Bitboard::Set(size_t pos) {
        bitboard |= 1UL << pos;
    }

    void Bitboard::Clear(Bitboard b) {
        bitboard &= ~b.bitboard;
    }

    bool Bitboard::Has(Bitboard b) const {
        return bitboard & b.bitboard;
    }

    bool Bitboard::IsZero() {
        return bitboard == 0;
    }
    
    // Restituisce true se c'è un bit settato nella traversa "rank"
    bool Bitboard::IsRank(size_t rank) {
        return Has(Bitboard(0xFFULL << rank*8));
    }

    // Restituiscono Bitboard corrispondenti ad una casa adiacente
    Bitboard Bitboard::North() {
        return Bitboard(bitboard << 8);
    }

    Bitboard Bitboard::South() {
        return Bitboard(bitboard >> 8);
    }

    Bitboard Bitboard::East() {
        return Bitboard((bitboard << 1) & notAFile);
    }

    Bitboard Bitboard::West() {
        return Bitboard((bitboard >> 1) & notHFile);
    }

    Bitboard Bitboard::NorthEast() {
        return Bitboard((bitboard << 9) & notAFile);
    }

    Bitboard Bitboard::NorthWest() {
        return Bitboard((bitboard << 7) & notHFile);
    }

    Bitboard Bitboard::SouthEast() {
        return Bitboard((bitboard >> 7) & notAFile);
    }

    Bitboard Bitboard::SouthWest() {
        return Bitboard((bitboard >> 9) & notHFile);
    }

    // Restituisce il bit meno significativo: 0110 -> 0010
    Bitboard Bitboard::LS1B() {
        return Bitboard(bitboard & -bitboard);
    }

    // Splitta una bitboard tipo così: 0110 -> {0100, 0010}
    std::vector<Bitboard> Bitboard::Split() {
        std::vector<Bitboard> v;
        while (!IsZero()) {
            auto b = LS1B();
            v.push_back(b);
            Clear(b);
        }
        return v;
    };

    // Splitta una bitboard e fa l'OR con un'altra contemporaneamente: 0110.Add(1000) -> {1100,1010}
    std::vector<Bitboard> Bitboard::Split(Bitboard add) {
        std::vector<Bitboard> v;
        while (!IsZero()) {
            auto b = LS1B();
            v.push_back(b + add);
            Clear(b);
        }
        return v;
    };


    // Ruota una bitboard di 180° e la riflette orizzontalmente
    // Praticamente è come guardare la scacchiera dal lato opposto
    // https://www.chessprogramming.org/Flipping_Mirroring_and_Rotating
    void Bitboard::Flip() {
        U64& x = bitboard;

        const U64 k1 = 0x5555555555555555ULL;
        const U64 k2 = 0x3333333333333333ULL;
        const U64 k4 = 0x0f0f0f0f0f0f0f0fULL;
        const U64 k5 = 0x00FF00FF00FF00FFULL;
        const U64 k6 = 0x0000FFFF0000FFFFULL;

        x = ((x >> 1) & k1) | ((x & k1) << 1);
        x = ((x >> 2) & k2) | ((x & k2) << 2);
        x = ((x >> 4) & k4) | ((x & k4) << 4);
        x = ((x >>  8) & k5)| ((x & k5) <<  8);
        x = ((x >> 16) & k6)| ((x & k6) << 16);
        x = ( x >> 32)      | ( x       << 32);
    }

    Bitboard Bitboard::Intersect(Bitboard b) {
        return Bitboard(bitboard & b.bitboard);
    }
    
    Bitboard Bitboard::Invert(Bitboard b) {
        return Bitboard(bitboard ^ b.bitboard);
    }

    std::string Bitboard::ToString() {
       // std::cout << "Si rompe qui?" << std::endl;
        size_t v;
        for(v = 0; v < 64; v++)
            if (Has(1UL << v)) break;
        
        //std::cout << v << std::endl;

        char file = std::string("abcdefgh").at(v % 8);
        char rank = std::string("12345678").at(v / 8);

        return std::string(1, file) + std::string(1, rank);
    }
}
