﻿using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using System.Threading;
using System.IO;

using System.Net;
using System.Net.Http;
using System.Net.WebSockets;
using System.Text.Json.Nodes;
using System.Text.Json;
using Microsoft.AspNetCore.Components.WebAssembly.Http;

namespace Client.Logic
{
    public class MyHttpClientHandler : HttpClientHandler
    {
        // Serve per includere i cookie nelle richieste

        protected override async Task<HttpResponseMessage>
        SendAsync(HttpRequestMessage request, CancellationToken cancellationToken)
        {
            request.SetBrowserRequestCredentials(BrowserRequestCredentials.Include);

            return await base.SendAsync(request, cancellationToken);
        }
    }
    public class Application
    {
        //public const string SERVER_URL = "http://localhost:8080";
        public const string SERVER_URL = "http://localhost/api";
        public const string WEBSOCKET_URL = "ws://localhost/ws";


        private static Application instance = null;

        public static Application Instance
        {
            get
            {
                if (instance == null)
                {
                    instance = new Application();
                }
                return instance;
            }
        }

        public Auth auth;
        public Game game;

        public HttpClient http;
        private ClientWebSocket websocket;

        public User user;
        public List<User> onlineUsers { get; set; }

        public string currentPanelClass;

        public class Request {
            public string uid;
            public string message;
        }
        public Request incomingRequest;

        public class GameResult
        {
            public string message;
            public string elo;
            public string type;
        }
        public GameResult gameResult;

        public string modalMessage;

        public delegate void UpdateUI();
        public UpdateUI updateUI;

        public AudioService audioService;

        public Application()
        {
            auth = new Auth();

            //var cookies = new CookieContainer();
            var httpHandler = new MyHttpClientHandler();

            //httpHandler.CookieContainer = cookies;
            http = new HttpClient(httpHandler);

            websocket = new ClientWebSocket();
            //websocket.Options.Cookies = cookies;

            currentPanelClass = "login";
            user = new User {
                UID = "",
                Username = "",
                Avatar = "",
                Elo = 0,
            };

            onlineUsers = new List<User>();
        }


        public void ChangePanel(string panel)
        {
            currentPanelClass = panel;
            System.Diagnostics.Debug.WriteLine("Panel: "+panel);
        }

        public async Task AuthSuccess() {
            var wsc = WebSocketConnect(WEBSOCKET_URL);
            var udr = http.GetAsync(SERVER_URL + "/user/username/" + user.Username);

            var response = await udr;
            bool wsStatus = await wsc;
            if (wsStatus && response.IsSuccessStatusCode)
            {
                _ = WebSocketListener();

                ChangePanel("home");
                _ = OnlineUserUpdater(TimeSpan.FromSeconds(5));


                string jsonString = response.Content.ReadAsStringAsync().Result;
                JsonNode resultNode = JsonNode.Parse(jsonString);

                user.UID = resultNode["data"]["uid"].GetValue<string>();
                // TODO avatar alternativo se non c'è
                user.Avatar = resultNode["data"]["avatar"].GetValue<string>();
                user.Elo = resultNode["data"]["elo"].GetValue<int>();

                updateUI.Invoke();

            }
        }

        private async Task<bool> WebSocketConnect(string uri)
        {
            for (int i = 0; i < 5; i++)
            {
                try {
                    await websocket.ConnectAsync(new Uri(uri), CancellationToken.None);
                    return true;
                }
                catch (Exception ex) {
                    System.Diagnostics.Debug.WriteLine($"Errore WS: {ex.Message}");
                }
            }
            return false;
        }

        private async Task WebSocketSend(string data) => 
            await websocket.SendAsync(Encoding.UTF8.GetBytes(data), WebSocketMessageType.Text, true, CancellationToken.None);

        private async Task WebSocketListener()
        {
            var buffer = WebSocket.CreateServerBuffer(2048);
            while (true)
            {
                WebSocketReceiveResult result;
                using (var ms = new MemoryStream())
                {
                    do
                    {
                        result = await websocket.ReceiveAsync(buffer, CancellationToken.None);
                        ms.Write(buffer.Array, buffer.Offset, result.Count);
                    } while (!result.EndOfMessage);

                    if (result.MessageType == WebSocketMessageType.Close)
                        break;

                    ms.Seek(0, SeekOrigin.Begin);

                    var jsonData = JsonNode.Parse(ms);

                    System.Diagnostics.Debug.WriteLine(jsonData);

                    switch (jsonData["type"].GetValue<string>())
                    {
                        case "challenge-request":
                            HandleChallengeRequest(jsonData["content"]["message"].GetValue<string>(), jsonData["content"]["uid"].GetValue<string>());
                            break;
                        case "game-start":
                            if (game is not null) continue;
                            await HandleGameStart(
                                jsonData["content"]["opponent"].GetValue<string>(),
                                jsonData["content"]["game-id"].GetValue<string>(),
                                jsonData["content"]["color"].GetValue<string>(),
                                jsonData["content"]["time"].GetValue<int>()
                            );
                            break;
                        case "move-played":
                            var color = jsonData["content"]["color"].GetValue<string>();
                            var move = jsonData["content"]["move"].GetValue<string>();
                            var time = jsonData["content"]["time"].GetValue<double>();

                            game.PlayReceivedMove(color,move,(int)Math.Round(time));
                            await audioService.PlaySound("move");
                            updateUI.Invoke();
                            break;
                        case "end-game":
                            var gameResult = jsonData["content"]["result"].GetValue<string>();
                            var elo = jsonData["content"]["elo"].GetValue<int>();

                            HandleEndGame(gameResult, elo);
                            break;
                    }
                }
            } 
        }

        private async Task OnlineUserUpdater(TimeSpan interval) {
            var timer = new PeriodicTimer(interval);
            do
            {
                if (currentPanelClass != "home")
                {
                    await timer.WaitForNextTickAsync();
                    continue;
                }

                var response = await http.GetAsync(SERVER_URL + "/users/online");
                if (response.IsSuccessStatusCode)
                {
                    onlineUsers.Clear();
                    string jsonString = response.Content.ReadAsStringAsync().Result;
                    JsonNode resultNode = JsonNode.Parse(jsonString);

                    foreach (var user in resultNode["data"].AsArray())
                    {
                        onlineUsers.Add(new User
                        {
                            UID = user["uid"].GetValue<string>(),
                            Username = user["username"].GetValue<string>(),
                            Avatar = user["avatar"].GetValue<string>(),
                            Elo = user["elo"].GetValue<int>(),
                        });
                    }

                    updateUI.Invoke();
                }

            } while (await timer.WaitForNextTickAsync());
        }


        public async Task SendChallengeRequest(string uid) {
            var content = new JsonObject();
            content.Add("uid", uid);
            var jsonData = new JsonObject();
            jsonData.Add("type", "challenge-request");
            jsonData.Add("content", content);


            await WebSocketSend(jsonData.ToJsonString());
        }

       void HandleChallengeRequest(string message, string uid) { 
            incomingRequest = new Request{ message = message, uid = uid };
            updateUI.Invoke();
       }

        void HandleEndGame(string result, int elo) {
            var message = "";
            switch (result) {
                case "victory":
                    message = "Hai Vinto!";
                    break;
                case "lose":
                    message = "Hai Perso";
                    break;
                case "draw":
                    message = "Patta";
                    break;
            }
            gameResult = new GameResult { elo = (elo > 0 ? "+" : "") + elo.ToString(), message = message, type = result };
            updateUI.Invoke();
        }

        public void CloseGame()
        {
            ChangePanel("home");
            gameResult = null;
            updateUI.Invoke();
            game = null;
        }

        public async Task SendChallengeDecline(string uid)
        {
            var content = new JsonObject();
            content.Add("uid", uid);
            var jsonData = new JsonObject();
            jsonData.Add("type", "challenge-decline");
            jsonData.Add("content", content);


            await WebSocketSend(jsonData.ToJsonString());
            incomingRequest = null;
        }


        public async Task SendChallengeAccept(string uid)
        {
            var content = new JsonObject();
            content.Add("uid", uid);
            var jsonData = new JsonObject();
            jsonData.Add("type", "challenge-accept");
            jsonData.Add("content", content);


            await WebSocketSend(jsonData.ToJsonString());
            incomingRequest = null;
            updateUI.Invoke();
        }

        public async Task SendChallengeComputer()
        {
            var content = new JsonObject();
            var jsonData = new JsonObject();
            jsonData.Add("type", "challenge-computer");
            jsonData.Add("content", content);

            await WebSocketSend(jsonData.ToJsonString());
            incomingRequest = null;
            updateUI.Invoke();
        }

        async Task HandleGameStart(string opponentUid, string gameId, string color, int time)
        {
            gameResult = null;
            if (opponentUid == "computer")
            {
                User user = new User();
                user.Username = "Computer";
                user.Elo = 2000;
                game = new Game(color == "white" ? Side.White : Side.Black, user, time, gameId);
            }
            else
            {
                game = new Game(color == "white" ? Side.White : Side.Black, await GetUserData(opponentUid), time, gameId);
            }

            ChangePanel("game");
            updateUI.Invoke();
        }


        public async Task<User> GetUserData(string uid)
        {
            var response = await http.GetAsync(SERVER_URL + "/user/id/" + uid);
            if (!response.IsSuccessStatusCode) return null;

            string jsonString = response.Content.ReadAsStringAsync().Result;
            JsonNode resultNode = JsonNode.Parse(jsonString);

            User user = new User();
                
            user.Username = resultNode["data"]["username"].GetValue<string>();
            user.Avatar = resultNode["data"]["avatar"].GetValue<string>(); // TODO avatar alternativo se non c'è
            user.Elo = resultNode["data"]["elo"].GetValue<int>();

            return user;
        }

        public async Task SendMove(string move, string gameId)
        {
            var content = new JsonObject();
            content.Add("move", move);
            content.Add("game-id", gameId);
            var jsonData = new JsonObject();
            jsonData.Add("type", "play-move");
            jsonData.Add("content", content);

            await WebSocketSend(jsonData.ToJsonString());

            System.Console.WriteLine(jsonData.ToJsonString());
            //System.Diagnostics.Debug.WriteLine(jsonData.ToJsonString());
        }
    }

}
