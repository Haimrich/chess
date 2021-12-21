using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace Client.Logic
{

    public class Square : IEquatable<Square>
    {
        private const string letters = "abcdefgh";

        protected int _file;
        public int File
        {
            get { return _file; }
            set
            {
                if (value < 0 || value > 7)
                    throw new ArgumentOutOfRangeException("Piece X coordinate must be between 0 and 7.");
                _file = value;
            }
        }

        protected int _rank;
        public int Rank
        {
            get { return _rank; }
            set
            {
                if (value < 0 || value > 7)
                    throw new ArgumentOutOfRangeException("Piece Y coordinate must be between 0 and 7.");
                _rank = value;
            }
        }

        public int X { get => File; set => File = value; }
        public int Y { get => Rank; set => Rank = value; }

        public string Position
        {
            get
            {
                return $"{letters[_file]}{_rank + 1}";
            }
            set
            {
                File = letters.IndexOf(value[0]);
                Rank = int.Parse(value[1].ToString()) - 1;
            }
        }


        public Square(int rank, int file)
        {
            Rank = rank;
            File = file;
        }

        public Square() : this(0, 0) { }

        public Square(Square other)
        {
            _rank = other.Rank;
            _file = other.File;
        }

        public Square(string position) => Position = position;


        public override bool Equals(object obj) => this.Equals(obj as Square);

        public bool Equals(Square other)
        {
            if (other is null) 
                return false;

            if (Object.ReferenceEquals(this, other)) 
                return true;

            if (this.GetType() != other.GetType())
                return false;

            return (X == other.X) && (Y == other.Y);
        }

        public override int GetHashCode() => (X, Y).GetHashCode();

        public static bool operator ==(Square a, Square b)
        {
            if (a is null)
            {
                if (b is null) 
                    return true;

                return false;
            }

            return a.Equals(b);
        }

        public static bool operator !=(Square a, Square b) => !(a == b);

        public Square CopyTranslate(int x, int y) {
            int file = _file + x;
            int rank = _rank + y; 

            if (file >= 0 && file < 8 && rank >= 0 && rank < 8)
                return new Square(rank, file);
            else
                return null;
        } 
    }
}
