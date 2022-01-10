using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace Client.Logic
{
    public class CastlingMove : Move
    {
        public Castling CastlingType;

        public CastlingMove(Side player, Square start, Square end, Castling oldCastling, Castling newCastling, Castling castlingType, Square oldEnPassantSquare = null) : base(player, start, end, oldCastling, newCastling, oldEnPassantSquare, null, null)
        {
            this.CastlingType = castlingType;
        }

        public override string GenerateNotation(Board board) => (CastlingType & Castling.King) != Castling.None ? "O-O" : "O-O-O";

    }
}
