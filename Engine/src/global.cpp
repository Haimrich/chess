#include "global.hpp"

#include <string>

namespace engine {

std::string trim(std::string& str, char c)
{
    size_t last = str.find_last_not_of(c)+1;
    if (last < str.length())
        str.erase(last);

    size_t first = str.find_first_not_of(c);
    if (first > 0)
        str.erase(0, first);

    return str;
}

bool isLower(char c) {
    return c >='a' && c <= 'z';
}

bool toUpper(char c) {
    return c - 'a' + 'A';
}

bool isDigit(char c) {
    return c >='0' && c <= '9';
}

}
