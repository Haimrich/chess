using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace Client.Logic
{
    public class Move
    {
        public Square Start;
        public Square End;

        public Piece CapturedPiece;
        public Side Player;

        public Square OldEnPassantSquare;
        public Square EnPassantSquare;

        public Castling OldCastlingOpportunities;
        public Castling CastlingOpportunities;

        public Move(Side player, Square start, Square end, Castling oldCastling, Castling newCastling, Square oldEnPassantSquare = null, Piece capturedPiece = null, Square enPassantSquare = null) { 
            this.Player = player;
            this.Start = new Square(start);    
            this.End = new Square(end);
            this.CapturedPiece = capturedPiece;
            this.OldEnPassantSquare = oldEnPassantSquare;
            this.EnPassantSquare = enPassantSquare;
            this.OldCastlingOpportunities = oldCastling;
            this.CastlingOpportunities = newCastling;
        }

        public Move(Side player, Square start, Square end, Castling oldCastling, Square oldEnPassantSquare = null, Piece capturedPiece = null, Square enPassantSquare = null)
        {
            this.Player = player;
            this.Start = new Square(start);
            this.End = new Square(end);
            this.CapturedPiece = capturedPiece;
            this.OldEnPassantSquare = oldEnPassantSquare;
            this.EnPassantSquare = enPassantSquare;
            this.OldCastlingOpportunities = oldCastling;
            this.CastlingOpportunities = oldCastling;
        }

        public virtual string GenerateNotation(Board board)
        {
            var piece = board[Start];

            string DisambiguateFileOrRank<T>() where T : Piece
            {
                var otherPieces = board.GetSimiliarPiecesTargetingSameSquare<T>(Player, End);
                System.Diagnostics.Debug.Assert(otherPieces.Any());

                if (otherPieces.Count() == 1)
                {
                    return "";
                }

                if (otherPieces.Where((T p) => p.Square.File == Start.File).Count() > 1)
                {
                    if (otherPieces.Where((T p) => p.Square.Rank == Start.Rank).Count() > 1)
                    {
                        return Start.Position;
                    }
                    return Start.Position[1].ToString();
                }
                return Start.Position[0].ToString();
            }

            var notation = "";
            switch (piece) {
                case Pawn:
                    if (CapturedPiece == null) return End.Position;
                    return Start.Position[0] + "x" + End.Position;
                case Knight:
                    notation = "N" + DisambiguateFileOrRank<Knight>();
                    break;
                case Bishop:
                    notation = "B" + DisambiguateFileOrRank<Bishop>();
                    break;
                case Rook:
                    notation = "R" + DisambiguateFileOrRank<Rook>();
                    break;
                case Queen:
                    notation = "Q" + DisambiguateFileOrRank<Queen>();
                    break;
                case King:
                    notation = "K";
                    break;
                default:
                    System.Diagnostics.Debug.Assert(false);
                    break;
            }

            if (CapturedPiece == null)  return notation + End.Position;

            return notation + "x" + End.Position;
        }
    }
}
