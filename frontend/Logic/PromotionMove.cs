using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace Client.Logic
{
    public class PromotionMove : Move
    {
        private Piece _promotionPiece;

        public PromotionMove(Side player, Square start, Square end, Castling oldCastling, Square oldEnPassantSquare = null, Piece capturedPiece = null, Square enPassantSquare = null) : base(player, start, end, oldCastling, oldEnPassantSquare, capturedPiece, enPassantSquare)
        {
            _promotionPiece = null;
        }

        public PromotionMove(Side player, Square start, Square end, Castling oldCastling, Castling newCastling, Square oldEnPassantSquare = null, Piece capturedPiece = null, Square enPassantSquare = null) : base(player, start, end, oldCastling, newCastling, oldEnPassantSquare, capturedPiece, enPassantSquare)
        {
            _promotionPiece = null;
        }

        public Piece PromotionPiece { 
            get { return _promotionPiece; }
            set { 
                _promotionPiece = value;
                _promotionPiece.Square.X = End.X;
                _promotionPiece.Square.Y = End.Y;
            }
        }

        public override string GenerateNotation(Board board)
        {
            var pt = _promotionPiece.Type[1].ToString().ToUpper();

            if (CapturedPiece == null)
                return End.Position + "=" + pt;

            return Start.Position[0] + "x" + End.Position + "=" + pt;
        }

    }
}
