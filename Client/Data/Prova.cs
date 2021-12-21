using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace Client.Data
{
    public class Prova
    {
        public int counter { get; set; }

        public Prova()
        {
            this.counter = 50;
        }

        public void increment() {
            counter--;
        }
    }
}
