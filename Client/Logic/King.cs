using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace Client.Logic
{
    public class King : Piece
    {
        public override string Type { get => (Color == Side.White ? "w" : "b") + "k"; }
        public King(Side color, int X, int Y) : base(color, X, Y)
        {
        }

        public override List<Move> GetMoves(Board board)
        {
            var moves = new List<Move>();
            var dirs = new[] { (1, 1), (-1, -1), (1, -1), (-1, 1), (0, 1), (0, -1), (1, 0), (-1, 0) };

            Castling newCastlingOpportunities = board.CastlingOpportunities & (Color == Side.White ? Castling.Black : Castling.White);

            for (int d = 0; d < 8; d++)
            {
                Square dest = Square.CopyTranslate(dirs[d].Item1, dirs[d].Item2);
                if (dest is not null)
                {
                    Piece targetPiece = board[dest];

                    if (targetPiece is null)
                        moves.Add(new Move(Color, Square, dest, board.CastlingOpportunities, newCastlingOpportunities, board.EnPassantSquare));
                    else if (targetPiece.Color != Color)
                        moves.Add(new Move(Color, Square, dest, board.CastlingOpportunities, newCastlingOpportunities, board.EnPassantSquare, targetPiece));
                }
            }

            if (board.IsSquareInCheck(Color, Square)) {
                return moves;
            }

            // Castling
            Castling castlingColor = Color == Side.White ? Castling.White : Castling.Black;
            Castling castlingOpportunities = board.CastlingOpportunities & castlingColor;

            if ((castlingOpportunities & Castling.King) > 0)
            {
                bool kingCastling = true;
                Square castlingPath = Square.CopyTranslate(dirs[6].Item1, dirs[6].Item2);
                for (int i = 0; i < 2; i++) {
                    if (board.IsSquareInCheck(Color, castlingPath) || board.SquareOccupiedBy(castlingPath) != Side.None)
                    {
                        kingCastling = false;
                        break;
                    }
                    castlingPath = castlingPath.CopyTranslate(dirs[6].Item1, dirs[6].Item2);
                }
                System.Diagnostics.Debug.WriteLine("KingCast: " + kingCastling);
                if (kingCastling) {
                    moves.Add(new CastlingMove(Color, Square, Square.CopyTranslate(+2,0), board.CastlingOpportunities, newCastlingOpportunities, Castling.King, board.EnPassantSquare));
                }
            }
            
            if ((castlingOpportunities & Castling.Queen) > 0)
            {
                bool queenCastling = true;
                Square castlingPath = Square.CopyTranslate(dirs[7].Item1, dirs[7].Item2);
                for (int i = 0; i < 3; i++)
                {
                    if (board.IsSquareInCheck(Color, castlingPath) || board.SquareOccupiedBy(castlingPath) != Side.None)
                    {
                        queenCastling = false;
                        break;
                    }
                    castlingPath = castlingPath.CopyTranslate(dirs[7].Item1, dirs[7].Item2);
                }
                System.Diagnostics.Debug.WriteLine("QuuenCast: " + queenCastling);
                if (queenCastling)
                {
                    moves.Add(new CastlingMove(Color, Square, Square.CopyTranslate(-2, 0), board.CastlingOpportunities, newCastlingOpportunities, Castling.Queen, board.EnPassantSquare));
                }
            }

            return moves;
        }
    }
}
