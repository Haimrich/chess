using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace Client.Logic
{
    public class Pawn : Piece
    {
        public override string Type { get => (Color == Side.White ? "w" : "b") + "p"; }
        public Pawn(Side color, int X, int Y) : base(color, X, Y)
        {
        }

        public override List<Move> GetMoves(Board board) { 
            var moves = new List<Move>();

            // Movimento

            Square dest = Square.CopyTranslate(0, Color == Side.White ? 1 : -1);
            Piece targetPiece = board[dest];

            if (targetPiece is null)
            {
                if (dest.Y == 7 || dest.Y == 0)
                    moves.Add(new PromotionMove(Color, Square, dest, board.CastlingOpportunities, board.EnPassantSquare));
                else
                    moves.Add(new Move(Color, Square, dest, board.CastlingOpportunities, board.EnPassantSquare));

                if (!_moved)
                {
                    dest = dest.CopyTranslate(0, Color == Side.White ? 1 : -1);
                    targetPiece = board[dest];
                    if (dest is not null && targetPiece is null)
                    {
                        Square eps = Square.CopyTranslate(0, Color == Side.White ? 1 : -1);
                        moves.Add(new Move(Color, Square, dest, board.CastlingOpportunities, board.EnPassantSquare, enPassantSquare: eps));
                    }
                }
            }

            // Catture 
            Action<Square> addCaptureMove = delegate (Square dest) {
                if (dest is not null)
                {
                    targetPiece = board[dest];
                    if (targetPiece is not null && targetPiece.Color != Color)
                    {
                        if (dest.Y == 7 || dest.Y == 0)
                            moves.Add(new PromotionMove(Color, Square, dest, board.CastlingOpportunities, board.EnPassantSquare, targetPiece));
                        else
                            moves.Add(new Move(Color, Square, dest, board.CastlingOpportunities, board.EnPassantSquare, targetPiece));
                    }
                    else if (dest == board.EnPassantSquare)
                    {
                        targetPiece = board[dest.Rank + (Color == Side.White ? -1 : +1), dest.File];
                        if (targetPiece is not null && targetPiece.Color != Color)
                            moves.Add(new Move(Color, Square, dest, board.CastlingOpportunities, board.EnPassantSquare, targetPiece));
                    }
                }
            };

            addCaptureMove(Square.CopyTranslate(-1, Color == Side.White ? 1 : -1));
            addCaptureMove(Square.CopyTranslate(1, Color == Side.White ? 1 : -1));

            return moves;
        }

    }
}