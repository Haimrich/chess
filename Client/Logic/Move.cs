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
    }
}
