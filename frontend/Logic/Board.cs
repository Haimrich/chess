using System;
using System.Collections.Generic;

namespace Client.Logic
{
    public class Board
    {
        // Matrice 8x8 di pezzi (o null)
        private Piece[,] _Board;
        
        // Indexers
        public Piece this[int rank, int file]
        {
            get => _Board[rank, file];
        }
        public Piece this[Square s]
        {
            get => s is not null ? _Board[s.Rank, s.File] : null;
        }

        // Cose per la partita

        Square blackKingPosition;
        Square whiteKingPosition;
        Side SideToMove;

        public Castling CastlingOpportunities;
        public Square EnPassantSquare;

        int HalfmoveClock;
        int FullmoveClock;

        // Lato del giocatore nella gui
        Side view_side;
        public string ViewSideClass { get => view_side == Side.White ? "white" : "black"; }

        // Lista dei pezzi sulla _Board
        public List<Piece> Piecies
        {
            get
            {
                List<Piece> pieces = new List<Piece>();
                for (int i = 0; i < 8; i++)
                    for (int j = 0; j < 8; j++)
                        if (_Board[i, j] is Piece)
                            pieces.Add(_Board[i, j]);
                return pieces;
            }
        }

        // Costruttori

        public Board()
        {
            _Board = new Piece[8, 8];

            SideToMove = Side.White;
            CastlingOpportunities = Castling.None;

            HalfmoveClock = 0;
            FullmoveClock = 1;

            view_side = Side.White;

            EnPassantSquare = null;
        }


        public Board(string fen) : this()
        {
            SetPositionFromFEN(fen);
        }

        // Carica posizione a partire da stringa in Notazione Forsyth-Edwards
        // Es. "[rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1]"
        public void SetPositionFromFEN(string fen)
        {
            fen = fen.Replace("[", "").Replace("]", "");
            string[] fen_fields = fen.Split(' ');

            // Campo 1 - Posizione
            int rank = 7;
            int file = 0;
            int skip;
            foreach (char c in fen_fields[0])
            {
                if (int.TryParse(c.ToString(), out skip))
                {
                    file += skip;
                }
                else
                {
                    switch (c)
                    {
                        case '/':
                            rank--;
                            if (file != 8 || rank < 0) throw new FormatException("FEN format error.");
                            file = 0;
                            break;
                        case 'k':
                        case 'K':
                            Side color = Char.IsUpper(c) ? Side.White : Side.Black;
                            _Board[rank, file] = new King(color, file, rank);
                            SetKingPosition(color, new Square(_Board[rank, file].Square));
                            file++;
                            break;
                        case 'q':
                        case 'Q':
                            _Board[rank, file] = new Queen(Char.IsUpper(c) ? Side.White : Side.Black, file, rank);
                            file++;
                            break;
                        case 'r':
                        case 'R':
                            _Board[rank, file] = new Rook(Char.IsUpper(c) ? Side.White : Side.Black, file, rank);
                            file++;
                            break;
                        case 'b':
                        case 'B':
                            _Board[rank, file] = new Bishop(Char.IsUpper(c) ? Side.White : Side.Black, file, rank);
                            file++;
                            break;
                        case 'n':
                        case 'N':
                            _Board[rank, file] = new Knight(Char.IsUpper(c) ? Side.White : Side.Black, file, rank);
                            file++;
                            break;
                        case 'p':
                        case 'P':
                            _Board[rank, file] = new Pawn(Char.IsUpper(c) ? Side.White : Side.Black, file, rank);
                            file++;
                            break;
                        default:
                            throw new FormatException($"FEN format error. Unrecognized character: {c}");
                    }
                }
            }

            // Campo 2 - Turno giocatore
            if (fen_fields[1][0] == 'w')
                SideToMove = Side.White;
            else if (fen_fields[1][0] == 'b')
                SideToMove = Side.Black;
            else
                throw new FormatException($"FEN format error. Unrecognized player turn: {fen_fields[1][0]}");

            // Campo 3 - Possibilità di arrocco
            foreach (char c in fen_fields[2])
            {
                switch (c)
                {
                    case 'K':
                        CastlingOpportunities |= (Castling.White & Castling.King);
                        break;
                    case 'k':
                        CastlingOpportunities |= (Castling.Black & Castling.King);
                        break;
                    case 'Q':
                        CastlingOpportunities |= (Castling.White & Castling.Queen);
                        break;
                    case 'q':
                        CastlingOpportunities |= (Castling.Black & Castling.Queen);
                        break;
                    default:
                        throw new FormatException($"FEN format error. Unrecognized character: {c}");
                }
            }

            // Campo 4 - En Passant
            if (fen_fields[3] != "-")
            {
                EnPassantSquare = new Square(fen_fields[3]);
            }

            // Campo 5 e 6 - Halfmove e fullmove clock
            HalfmoveClock = int.Parse(fen_fields[4]);
            FullmoveClock = int.Parse(fen_fields[5]);
        }


        public Side SquareOccupiedBy(Square s) => _Board[s.Y, s.X] is Piece p ? p.Color : Side.None;

        public void PlayMove(Move m)
        {
            if (m.CapturedPiece is not null && EnPassantSquare == m.End)
                _Board[m.End.Rank + (m.Player == Side.White ? -1 : +1), m.End.File] = null;

            if (m is PromotionMove pm && pm.PromotionPiece is not null)
                _Board[m.End.Rank, m.End.File] = pm.PromotionPiece;
            else
                _Board[m.End.Rank, m.End.File] = _Board[m.Start.Rank, m.Start.File];

            if (m is CastlingMove cm)
            {
                int pre_rook_offset = cm.CastlingType == Castling.King ? +3 : -4;
                int post_rook_offset = cm.CastlingType == Castling.King ? +1 : -1;
                _Board[m.Start.Rank, m.Start.File + post_rook_offset] = _Board[m.Start.Rank, m.Start.File + pre_rook_offset];
                _Board[m.Start.Rank, m.Start.File + pre_rook_offset] = null;
            }

            _Board[m.Start.Rank, m.Start.File] = null;

            EnPassantSquare = m.EnPassantSquare;
            CastlingOpportunities = m.CastlingOpportunities;

            if (_Board[m.End.Rank, m.End.File] is King k)
                SetKingPosition(k.Color, new Square(m.End));
        }

        public void UndoMove(Move m)
        {
            if (m is CastlingMove cm)
            {
                int pre_rook_offset = cm.CastlingType == Castling.King ? +3 : -4;
                int post_rook_offset = cm.CastlingType == Castling.King ? +1 : -1;
                _Board[m.Start.Rank, m.Start.File + pre_rook_offset] = _Board[m.Start.Rank, m.Start.File + post_rook_offset];
                _Board[m.Start.Rank, m.Start.File + post_rook_offset] = null;
            }

            _Board[m.Start.Rank, m.Start.File] = _Board[m.End.Rank, m.End.File];

            if (m.CapturedPiece is not null)
            {
                if (m.CapturedPiece.Square == m.End)
                {
                    _Board[m.End.Rank, m.End.File] = m.CapturedPiece;
                }
                else
                {
                    _Board[m.End.Rank + (m.Player == Side.White ? -1 : +1), m.End.File] = m.CapturedPiece;
                    _Board[m.End.Rank, m.End.File] = null;
                }
            }
            else
            {
                _Board[m.End.Rank, m.End.File] = null;
            }


            EnPassantSquare = m.OldEnPassantSquare;
            CastlingOpportunities = m.OldCastlingOpportunities;

            if (_Board[m.Start.Rank, m.Start.File] is King k)
                SetKingPosition(k.Color, new Square(m.Start));
        }

        public bool IsKingInCheck(Side color)
        {
            Square kingPosition = GetKingPosition(color);
            return IsSquareInCheck(color, kingPosition);
        }


        public bool IsSquareInCheck(Side color, Square square)
        {
            var orto_dirs = new[] { (0, 1), (0, -1), (1, 0), (-1, 0) };
            var diag_dirs = new[] { (1, 1), (-1, -1), (1, -1), (-1, 1) };
            var knight_squares = new[] { (1, 2), (-1, 2), (-2, 1), (-2, -1), (2, 1), (2, -1), (-1, -2), (1, -2) };
            var king_squares = new[] { (0, 1), (0, -1), (1, 0), (-1, 0), (1, 1), (-1, -1), (1, -1), (-1, 1) };

            // Funzioni locali
            
            bool checkDirections<T, U>((int, int)[] dirs)
            {
                for (int i = 0; i < dirs.Length; i++)
                {
                    Square s = square.CopyTranslate(dirs[i].Item1, dirs[i].Item2);
                    while (s is not null)
                    {
                        if (this[s] is not null)
                        {
                            if (this[s].Color != color && (this[s] is T || this[s] is U))
                                return true;
                            else
                                break;
                        }

                        s = s.CopyTranslate(dirs[i].Item1, dirs[i].Item2);
                    }
                }
                return false;
            }

            bool checkSquares<T>((int, int)[] squares)
            {
                for (int i = 0; i < squares.Length; i++)
                {
                    Square s = square.CopyTranslate(squares[i].Item1, squares[i].Item2);

                    if (this[s] is not null && this[s].Color != color && this[s] is T)
                        return true;
                }
                return false;
            }

            if (checkDirections<Queen, Rook>(orto_dirs))
                return true;
            if (checkDirections<Queen, Bishop>(diag_dirs))
                return true;
            if (checkSquares<Knight>(knight_squares))
                return true;
            if (checkSquares<King>(king_squares))
                return true;

            int pawn_dir = color == Side.White ? +1 : -1;
            if (checkSquares<Pawn>(new[] { (-1, pawn_dir), (1, pawn_dir) }))
                return true;

            return false;
        }

        // Metodo generico per capire se la stessa casella è minacciata da più pezzi dello stesso tipo
        // Serve per disambiguare la notazione algebrica delle mosse
        public List<P> GetSimiliarPiecesTargetingSameSquare<P>(Side myColor, Square square) where P:Piece {
            var pieces = new List<P>();

            var orto_dirs = new[] { (0, 1), (0, -1), (1, 0), (-1, 0) };
            var diag_dirs = new[] { (1, 1), (-1, -1), (1, -1), (-1, 1) };
            var knight_squares = new[] { (1, 2), (-1, 2), (-2, 1), (-2, -1), (2, 1), (2, -1), (-1, -2), (1, -2) };

            // Funzioni locali

            void checkDirections((int, int)[] dirs)
            {
                for (int i = 0; i < dirs.Length; i++)
                {
                    Square s = square.CopyTranslate(dirs[i].Item1, dirs[i].Item2);
                    while (s is not null)
                    {
                        if (this[s] is not null)
                        {
                            if (this[s].Color == myColor && this[s] is P p)
                            {
                                pieces.Add(p);
                            }
                            else
                            {
                                break;
                            }
                        }

                        s = s.CopyTranslate(dirs[i].Item1, dirs[i].Item2);
                    }
                }
            }

            void checkSquares((int, int)[] squares)
            {
                for (int i = 0; i < squares.Length; i++)
                {
                    Square s = square.CopyTranslate(squares[i].Item1, squares[i].Item2);

                    if (this[s] is not null && this[s].Color == myColor && this[s] is P p)
                    {
                        pieces.Add(p);
                    }
                }
            }

            // Fai cose il base al tipo di pezzo

            switch (typeof(P))
            {
                case Type q when q == typeof(Queen):
                    checkDirections(orto_dirs);
                    checkDirections(diag_dirs);
                    break;
                case Type b when b == typeof(Bishop):
                    checkDirections(diag_dirs);
                    break;
                case Type r when r == typeof(Rook):
                    checkDirections(orto_dirs);
                    break;
                case Type n when n == typeof(Knight):
                    checkSquares(knight_squares);
                    break;
                case Type p when p == typeof(Pawn):
                    int pawn_dir = myColor == Side.White ? +1 : -1;
                    checkSquares(new[] { (-1, pawn_dir), (1, pawn_dir) });
                    break;
            }

            return pieces;
        }

        private void SetKingPosition(Side color, Square position)
        {
            if (color == Side.White)
                whiteKingPosition = position;
            else if (color == Side.Black)
                blackKingPosition = position;
        }

        private Square GetKingPosition(Side color)
        {
            if (color == Side.White)
                return whiteKingPosition;
            else if (color == Side.Black)
                return blackKingPosition;

            return null;
        }


    }
}
