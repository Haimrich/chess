#ifndef ENGINE_HANDLER_H_
#define ENGINE_HANDLER_H_

#include "pistache/http.h"

using namespace Pistache;

class Handler : public Http::Handler
{
public:
    HTTP_PROTOTYPE(Handler)

    void onRequest(const Http::Request& request, Http::ResponseWriter response) override;
};

#endif