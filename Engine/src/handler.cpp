#include  "handler.hpp"

#include "pistache/http.h"

#include "rapidjson/document.h"
#include "rapidjson/writer.h"
#include "rapidjson/stringbuffer.h"

#include "position.hpp"
#include "explorer.hpp"
#include "move.hpp"

using namespace Pistache;
using namespace rapidjson;

void Handler::onRequest(const Http::Request& request, Http::ResponseWriter response)
{
    auto body = request.body().c_str();

    Document data;
    data.Parse(body);

    if (data.HasMember("fen")) {
        std::string fen = data["fen"].GetString();

        engine::Position position(fen);
        engine::Move move = engine::Explorer().Search(position, 1e5, 1000);

        std::string uciMove = position.MoveToString(move);
        std::cout << "MOSSA TROVATA: " << uciMove << std::endl;

        response.send(Http::Code::Ok, uciMove);
    } else {
        response.send(Http::Code::Ok, "Ciao.\n");
    }
    
}
