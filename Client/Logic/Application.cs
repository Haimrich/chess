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

namespace Client.Logic
{
    public class Application
    {
        public const string SERVER_URL = "http://localhost:8080";

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

        public delegate void UpdateUI();
        public UpdateUI updateUI;

        public Application()
        {
            auth = new Auth();
            game = new Game();

            var cookies = new CookieContainer();
            var httpHandler = new HttpClientHandler();
            httpHandler.CookieContainer = cookies;
            http = new HttpClient(httpHandler);

            websocket = new ClientWebSocket();
            websocket.Options.Cookies = cookies;

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
        }

        public async Task AuthSuccess() {
            var wsc = WebSocketConnect("ws://localhost:8080/ws");
            var udr = http.GetAsync(SERVER_URL + "/user/" + user.Username);

            var response = await udr;
            bool wsStatus = await wsc;
            if (wsStatus && response.IsSuccessStatusCode)
            {
                _ = WebSocketListener();
                _ = OnlineUserUpdater(TimeSpan.FromSeconds(5));


                string jsonString = response.Content.ReadAsStringAsync().Result;
                JsonNode resultNode = JsonNode.Parse(jsonString);

                // TODO avatar alternativo se non c'è
                user.Avatar = resultNode["data"]["avatar"].GetValue<string>();
                user.Elo = resultNode["data"]["elo"].GetValue<int>();

                ChangePanel("home");

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
                    using (var reader = new StreamReader(ms, Encoding.UTF8))
                        System.Diagnostics.Debug.WriteLine(await reader.ReadToEndAsync()); // TODO
                }
            } 
        }

        private async Task OnlineUserUpdater(TimeSpan interval) {
            var timer = new PeriodicTimer(interval);
            do
            {
                var response = await http.GetAsync(SERVER_URL + "/users/online");
                if (response.IsSuccessStatusCode)
                {
                    onlineUsers.Clear();
                    string jsonString = response.Content.ReadAsStringAsync().Result;
                    System.Diagnostics.Debug.WriteLine(jsonString);
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
                        System.Diagnostics.Debug.WriteLine(onlineUsers.Count);
                    }

                    updateUI.Invoke();
                }

            } while (await timer.WaitForNextTickAsync());
        }
    }

}
