#include  "handler.hpp"

#include "pistache/http.h"

using namespace Pistache;

void Handler::onRequest(const Http::Request& request, Http::ResponseWriter response)
{
    std::string body = request.body();

    

    std::this_thread::sleep_for(std::chrono::milliseconds(10000));
    response.send(Pistache::Http::Code::Ok, "Hello World\n");
}
