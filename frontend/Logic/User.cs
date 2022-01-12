using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace Client.Logic
{
    public class User
    {
 
        public string UID;
        public string Username;
        
        public int Elo;

        private string _avatar;

        public string Avatar { 
            get {
                if (_avatar == null || _avatar == "") 
                    return "images/default_user_image.png";

                return Application.SERVER_URL + "/avatar/" + _avatar;
            }
            set => _avatar = value;
        }
        
    }
}
