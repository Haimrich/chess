#include <thread>

#include "pistache/net.h"
#include "pistache/http.h"
#include "pistache/endpoint.h"

#include "handler.hpp"

using namespace Pistache;


int main()
{
    Address addr(Ipv4::any(), Port(9080));

    int threads = std::thread::hardware_concurrency();
    auto opts = Http::Endpoint::options().threads(threads);

    Http::Endpoint server(addr);
    server.init(opts);
    server.setHandler(Http::make_handler<Handler>());
    server.serve();
}