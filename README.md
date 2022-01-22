# Chess
An online chess platform developed as Course Project UniCT 2022 

Advanced Programming Languages | Distributed Systems and Big Data

## Description
- [Backend](https://github.com/Haimrich/chess/tree/main/backend) written in Go using [Gin](https://github.com/gin-gonic/gin) web framework and [Gorilla](https://github.com/gorilla/websocket) websockets. This monolithic component has been splitted in the following microservices:
  - [User Service](https://github.com/Haimrich/chess/tree/main/user) handles user signup, login and authentication. It also generates JWT Tokens for authorization in other microservices.
  - [WebSocket Node](https://github.com/Haimrich/chess/tree/main/wsnode) handles WebSocket connections with clients.
  - [Dispatcher](https://github.com/Haimrich/chess/tree/main/dispatcher) routes messages from WS Nodes to Game and Challenge Services.
  - [Challenge Service](https://github.com/Haimrich/chess/tree/main/dispatcher) handles challenge sending and accepting.
  - [Game Service](https://github.com/Haimrich/chess/tree/main/dispatcher) handles games status, timer, legal moves, etc.
- [Frontend](https://github.com/Haimrich/chess/tree/main/frontend) written in C# using [.NET Blazor WebAssembly](https://dotnet.microsoft.com/en-us/apps/aspnet/web-apps/blazor)
- [Engine](https://github.com/Haimrich/chess/tree/main/engine) written in C++ using [Pistache](https://github.com/pistacheio/pistache) HTTP server. [Bitboard](https://www.chessprogramming.org/Bitboards) representation was adopted as board representation while search logic was inspired by [carnatus](https://github.com/zserge/carnatus) and [sunfish](https://github.com/thomasahle/sunfish).

## Usage