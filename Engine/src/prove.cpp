/*
#include <iostream>
#include <bitset>
#include <vector>

using namespace std;

int main()
{
    typedef uint_fast64_t U64;
    U64 x = 395;
    U64 original = x;
    std::vector<U64> moves;
    int count = 0;
    while (x) {
        count++;
        U64 ls1b = x & -x; // isolate LS1B
        if (ls1b) {
            moves.push_back(ls1b);
        }
        x &= x - 1;
    }
    cout << count << std::endl;
    cout << std::bitset<64>(original) << std::endl;
    for (auto y : moves)
        cout << std::bitset<64>(y) << std::endl;

    return 0;
}*/