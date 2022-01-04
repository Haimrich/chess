using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace Client.Logic
{

    public enum Side { None, White = 0, Black = 1 };

    public enum Castling { White = 0b_1100, Black = 0b_0011, King = 0b_0101, Queen = 0b_1010, None = 0b_0000 , All = 0b_1111 }

}
