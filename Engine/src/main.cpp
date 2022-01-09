#include <thread>
#include <iostream>

#include "pistache/endpoint.h"

#include "handler.hpp"

using namespace Pistache;

int main()
{
    std::cout << "ENGINE START - Server listening..." << std::endl;

    Pistache::Address addr(Pistache::Ipv4::any(), Pistache::Port(9080));

    int threads = std::thread::hardware_concurrency();
    auto opts = Pistache::Http::Endpoint::options().threads(threads);

    Http::Endpoint server(addr);
    server.init(opts);
    server.setHandler(Http::make_handler<Handler>());
    server.serve();

}