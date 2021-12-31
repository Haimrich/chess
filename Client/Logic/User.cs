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
                if (_avatar == null) 
                    return "TODO";

                return Application.SERVER_URL + "/avatar/" + _avatar;
            }
            set => _avatar = value;
        }
        
    }
}
