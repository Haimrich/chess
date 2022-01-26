using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace Client.Logic
{
    public class Knight : Piece
    {
        public override string Type { get => (Color == Side.White ? "w" : "b") + "n"; }
        
        public Knight(Side color, int X, int Y) : base(color, X, Y)
        {
        }

        public override List<Move> GetMoves(Board board)
        {
            var moves = new List<Move>();

            for (int i = 0; i < 8; i++) 
            {
                int x = ((i / 2) % 2 != 0 ? +1 : -1) * (i / 4 != 0 ? 2 : 1);
                int y = (i % 2 != 0 ? -1 : +1) * (i / 4 != 0 ? 1 : 2);

                var dest = Square.CopyTranslate(x,y);
                if (dest is not null) 
                {
                    Piece targetPiece = board[dest];

                    if (targetPiece is null)
                        moves.Add(new Move(Color, Square, dest, board.CastlingOpportunities, board.EnPassantSquare));
                    else if (targetPiece.Color != Color)
                        moves.Add(new Move(Color, Square, dest, board.CastlingOpportunities, board.EnPassantSquare, targetPiece));
                }
             }

            return moves;
        }
    }
}
