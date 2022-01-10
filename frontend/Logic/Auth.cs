using Microsoft.AspNetCore.Components.Forms;
using System;
using System.Net;
using System.Net.Http;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Text.Json;
using System.Text.Json.Nodes;
using System.Threading.Tasks;

namespace Client.Logic
{
    public class Auth
    {
        public string username;
        public string password;
        public string confirm_password;
        private IBrowserFile avatar;
        const long maxAvatarFileSize = 2 << 20;


        public bool canSubmit;
        public string errorMessage;

        public Auth()
        {
            avatar = null;
            canSubmit = true;
            errorMessage = "";
        }

        public void LoadAvatar(InputFileChangeEventArgs e)
        {
            if (e.FileCount > 0)
            {
                avatar = e.File;
                if (avatar.Size > maxAvatarFileSize)
                    errorMessage = "Dimensione massima dell'avatar 2MB.<br \\>";
            }
            else
            {
                avatar = null;
            }
        }

        public async Task Signup()
        {
            if (!canSubmit) return;
            canSubmit = false;
            errorMessage = "";

            if (password != confirm_password)
                errorMessage += "Le password non coincidono.<br \\>";

            if (avatar != null && avatar.Size > maxAvatarFileSize)
                errorMessage += "Dimensione massima dell'avatar 2MB.<br \\>";

            if (errorMessage.Length > 0)
            {
                canSubmit = true;
                return;
            }


            MultipartFormDataContent form = new MultipartFormDataContent();

            form.Add(new StringContent(username), "username");
            form.Add(new StringContent(password), "password");
            form.Add(new StringContent(confirm_password), "confirm_password");

            if (avatar != null)
            {
                var fileContent = new StreamContent(avatar.OpenReadStream(maxAvatarFileSize));
                form.Add(fileContent, "avatar", avatar.Name);
            }

            HttpResponseMessage response = await Application.Instance.http.PostAsync(Application.SERVER_URL + "/signup", form);
            string sd = response.Content.ReadAsStringAsync().Result;
            System.Diagnostics.Debug.WriteLine(sd);

            if (response.IsSuccessStatusCode)
            {
                password = "";
                username = "";
                errorMessage = "";

                Application.Instance.ChangePanel("login");
            }

            canSubmit = true;
        }


        public async Task Login()
        {
            if (!canSubmit) return;
            canSubmit = false;

            var loginForm = new JsonObject();
            loginForm["username"] = username;
            loginForm["password"] = password;


            var response = await Application.Instance.http.PostAsync(
                Application.SERVER_URL + "/login",
                new StringContent(loginForm.ToJsonString(), Encoding.UTF8, "application/json")
            );


            if (response.IsSuccessStatusCode)
            {
                Application.Instance.user.Username = username;
                await Application.Instance.AuthSuccess();
            }
            else {
                string jsonString = response.Content.ReadAsStringAsync().Result;
                JsonNode resultNode = JsonNode.Parse(jsonString);
                errorMessage = resultNode["message"].GetValue<string>();
            }

            canSubmit = true;
            System.Diagnostics.Debug.WriteLine("Login Finito.");
        }
    }
}
