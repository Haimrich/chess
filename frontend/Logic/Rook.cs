using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace Client.Logic
{
    public class Rook : Piece
    {
        public override string Type { get => (Color == Side.White ? "w" : "b") + "r"; }
        
        public Rook(Side color, int X, int Y) : base(color, X, Y)
        {
        }

        public override List<Move> GetMoves(Board board)
        {
            var moves = new List<Move>();
            var dirs = new[] { (0, 1), (0, -1), (1, 0), (-1, 0) };

            Castling newCastlingOpportunities = board.CastlingOpportunities & ( Castling.All ^
                    ( (Color == Side.White  ? Castling.White    : Castling.Black    ) & 
                      (Square.File == 0     ? Castling.Queen     : Castling.King    ) ) );

            for (int d = 0; d < 4; d++)
            {
                Square dest = Square.CopyTranslate(dirs[d].Item1, dirs[d].Item2);
                while (dest is not null)
                {
                    Piece targetPiece = board[dest];

                    if (targetPiece is null)
                    {
                        moves.Add(new Move(Color, Square, dest, board.CastlingOpportunities, newCastlingOpportunities, board.EnPassantSquare));
                    }
                    else
                    {
                        if (targetPiece.Color != Color)
                            moves.Add(new Move(Color, Square, dest, board.CastlingOpportunities, newCastlingOpportunities, board.EnPassantSquare, targetPiece));

                        break;
                    }

                    dest = dest.CopyTranslate(dirs[d].Item1, dirs[d].Item2);
                }
            }

            return moves;
        }

    }
}
