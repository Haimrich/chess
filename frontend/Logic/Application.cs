using System;
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
    public class Application
    {
        // Singleton
        private static Application instance = null;
        public static Application Instance
        {
            get
            {
                if (instance == null)
                {
                    instance = new Application("http://localhost/");
                }
                return instance;
            }
        }

        // Authorization
        public Auth auth;

        // Partita
        public Game game;

        // HTTP and WebSocket Clients

        public readonly Uri SERVER_URL;
        public readonly Uri WEBSOCKET_URL;

        public HttpClient http;
        private ClientWebSocket websocket;

        // Cose da visualizzare nella GUI

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

        // Delegato per aggiornare UI quando voglio

        public delegate void UpdateUI();
        public UpdateUI updateUI;

        // Audio 

        public AudioService audioService;

        // -----

        public Application(string baseUrl)
        {
            auth = new Auth();

            var httpHandler = new CookieHttpClientHandler();

            http = new HttpClient(httpHandler);
            SERVER_URL = new Uri(baseUrl + "api");
            WEBSOCKET_URL = new UriBuilder(baseUrl + "ws"){
                Scheme = Uri.UriSchemeWs,
                Port = SERVER_URL.IsDefaultPort ? -1 : SERVER_URL.Port
            }.Uri;

            websocket = new ClientWebSocket();

            currentPanelClass = "login";
            user = new User {
                UID = "",
                Username = "",
                Avatar = "",
                Elo = 0,
            };

            onlineUsers = new List<User>();
            instance = this;
        }

        // Cambia pannello della pagina web
        public void ChangePanel(string panel)
        {
            currentPanelClass = panel;
            System.Diagnostics.Debug.WriteLine("Panel: "+panel);
        }

        // Autorizzazione completata con successo
        public async Task AuthSuccess() {
            // Connetti websocket
            var wsc = WebSocketConnect(WEBSOCKET_URL);
            // Richiedi info sull'utente loggato
            var udr = http.GetAsync(SERVER_URL + "/user/username/" + user.Username);

            var response = await udr;
            bool wsStatus = await wsc;
            if (wsStatus && response.IsSuccessStatusCode)
            {
                // Task infinito che ascolta il canale WS
                _ = WebSocketListener();

                ChangePanel("home");

                // Task infinito che aggiorna la lista degli utenti online
                _ = OnlineUserUpdater(TimeSpan.FromSeconds(5));


                string jsonString = response.Content.ReadAsStringAsync().Result;
                JsonNode resultNode = JsonNode.Parse(jsonString);

                user.UID = resultNode["data"]["uid"].GetValue<string>();
                user.Avatar = resultNode["data"]["avatar"].GetValue<string>();
                user.Elo = resultNode["data"]["elo"].GetValue<int>();

                updateUI.Invoke();
            }
        }

        // ==== WEBSOCKETS ====

        private async Task<bool> WebSocketConnect(Uri uri)
        {
            for (int i = 0; i < 5; i++)
            {
                try {
                    await websocket.ConnectAsync(uri, CancellationToken.None);
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

        // Task che ascolta il canale websocket in attesa di messaggi ricevuti
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

                    // In base al tipo di messaggio ricevuto faccio cose
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

        // Task che aggiorna la lista degli utenti online
        private async Task OnlineUserUpdater(TimeSpan interval) {
            var timer = new PeriodicTimer(interval);
            do
            {
                // Se non è in home aspetta
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

        
        // ==== WEBSOCKET MESSAGE HANDLERS ==== //
        void HandleChallengeRequest(string message, string uid) { 
            incomingRequest = new Request{ message = message, uid = uid };
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

        // ==== WEBSOCKET MESSAGE SENDERS ==== //

        public async Task SendChallengeRequest(string uid) {
            var content = new JsonObject();
            content.Add("uid", uid);
            var jsonData = new JsonObject();
            jsonData.Add("type", "challenge-request");
            jsonData.Add("content", content);


            await WebSocketSend(jsonData.ToJsonString());
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

        // ==== CHIUDI PARTITA ==== //

        public void CloseGame()
        {
            ChangePanel("home");
            gameResult = null;
            updateUI.Invoke();
            game = null;
        }
    }


    // Serve per includere i cookie nelle richieste
    public class CookieHttpClientHandler : HttpClientHandler
    {
        protected override async Task<HttpResponseMessage>
        SendAsync(HttpRequestMessage request, CancellationToken cancellationToken)
        {
            request.SetBrowserRequestCredentials(BrowserRequestCredentials.Include);

            return await base.SendAsync(request, cancellationToken);
        }
    }
}
