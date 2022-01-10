using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace Client.Logic
{
    public class Queen : Piece
    {
        public override string Type { get => (Color == Side.White ? "w" : "b") + "q"; }
        public Queen(Side color, int X, int Y) : base(color, X, Y)
        {
        }

        public override List<Move> GetMoves(Board board)
        {
            var moves = new List<Move>();
            var dirs = new[] { (1, 1), (-1, -1), (1, -1), (-1, 1), (0, 1), (0, -1), (1, 0), (-1, 0) };

            for (int d = 0; d < 8; d++)
            {
                Square dest = Square.CopyTranslate(dirs[d].Item1, dirs[d].Item2);
                while (dest is not null)
                {
                    Piece targetPiece = board[dest];

                    if (targetPiece is null)
                    {
                        moves.Add(new Move(Color, Square, dest, board.CastlingOpportunities, board.EnPassantSquare));
                    }
                    else
                    {
                        if (targetPiece.Color != Color)
                            moves.Add(new Move(Color, Square, dest, board.CastlingOpportunities, board.EnPassantSquare, targetPiece));

                        break;
                    }

                    dest = dest.CopyTranslate(dirs[d].Item1, dirs[d].Item2);
                }
            }

            return moves;
        }

    }
}
