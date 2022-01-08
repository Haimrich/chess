#include <thread>
#include <iostream>
/*
#include "pistache/net.h"
#include "pistache/http.h"
#include "pistache/endpoint.h"

#include "handler.hpp"
*/

#include "bitboard.hpp"
#include "position.hpp"
#include "explorer.hpp"
#include "move.hpp"

#include <exception>

//using namespace Pistache;


int main()
{
    std::cout << "Bitboard: " << sizeof(engine::Bitboard) << std::endl;
    std::cout << "Position: " << sizeof(engine::Position) << std::endl;

    engine::Position pos;
    engine::Explorer exp;

    engine::Move m = exp.Search(pos, 20, 99999);

    std::cout << pos.MoveToString(m) << std::endl;

    //Address addr(Ipv4::any(), Port(9080));

    //int threads = std::thread::hardware_concurrency();
    //auto opts = Http::Endpoint::options().threads(threads);

    //Http::Endpoint server(addr);
    //server.init(opts);
    //server.setHandler(Http::make_handler<Handler>());
    //server.serve();
}