using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace Client.Logic
{
    public abstract class Piece
    {
        public Side Color;
        public Square Square;

        protected bool _moved;
        public bool Moved { get => _moved; }

        public abstract string Type { get; }

        public List<Move> PossibleMoves;

        public Piece(Side color, int x, int y)
        {
            this.Color = color;
            this.Square = new Square(file: x, rank: y);
            _moved = false;
        }

        public Piece(Piece p)
        {
            this.Color = p.Color;
            this.Square = new Square(p.Square);
            this._moved = p.Moved;
        }

        public void MoveTo(Square destination) { 
            _moved = true;
            Square.Rank = destination.Rank;
            Square.File = destination.File;
        }

        // Da implementare per ciascun tipo di pezzo

        public abstract List<Move> GetMoves(Board board);
        /*
        public virtual List<Move> GetMoves(Board board) { 
            return new List<Move>();
        }
        */

    }
}