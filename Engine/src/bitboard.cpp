#include "bitboard.hpp"

#include <string>
#include <vector>

namespace engine {

    Bitboard::Bitboard(std::string square) {
        size_t file = std::string("abcdefgh").find(square.at(0));
        size_t rank = square.at(1) - '0';

        bitboard = 1 << rank*file;
    }

    void Bitboard::Set(size_t pos) {
        bitboard |= 1 << pos;
    }

    void Bitboard::Clear(Bitboard b) {
        bitboard &= ~b.bitboard;
    }

    bool Bitboard::Has(Bitboard b) {
        return bitboard & b.bitboard;
    }

    bool Bitboard::IsZero() {
        return bitboard == 0;
    }

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

    Bitboard Bitboard::LS1B() {
        return Bitboard(bitboard & -bitboard);
    }

    std::vector<Bitboard> Bitboard::Split(Bitboard add = 0) {
        std::vector<Bitboard> v;
        while (!IsZero()) {
            auto b = LS1B();
            v.push_back(b + add);
            Clear(b);
        }
        return v;
    };

}
