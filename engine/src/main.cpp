#include <thread>
#include <iostream>

#include <cstdlib>
#include <csignal>

#include "pistache/endpoint.h"

#include "handler.hpp"

using namespace Pistache;

void StopHandler(int s){
    std::cout << std::endl << "🛑 ENGINE STOPPED: " << s << std::endl;
    exit(1); 
}

int main()
{
    struct sigaction sigHandler;
    sigHandler.sa_handler = StopHandler;
    sigemptyset(&sigHandler.sa_mask);
    sigHandler.sa_flags = 0;
    sigaction(SIGINT, &sigHandler, NULL);
    sigaction(SIGTERM, &sigHandler, NULL);


    std::cout << "🟢 ENGINE START: Server listening..." << std::endl;

    Pistache::Address addr(Pistache::Ipv4::any(), Pistache::Port(9080));

    int threads = std::thread::hardware_concurrency();
    auto opts = Pistache::Http::Endpoint::options().threads(threads);

    Http::Endpoint server(addr);
    server.init(opts);
    server.setHandler(Http::make_handler<Handler>());
    server.serve();

}