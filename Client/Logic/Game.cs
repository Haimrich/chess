using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace Client.Logic
{
    public class Game
    {
        // Cose per GUI
        public List<Piece> Piecies;
        public Piece SelectedPiece;
        public List<Move> PossibleMoves;
        public List<Square> LastMoveSquares;

        private Side _viewSide;
        private PromotionMove PromotionMove;
        
        public string ViewSideClass { get => _viewSide == Side.White ? "white" : "black"; }
        public string PromotionModalClass { get => PromotionMove is null ? "hidden" : ""; }


        // Altre cose
        private Board _board;
        private Side SideToMove;

        int HalfmoveClock;
        int FullmoveClock;

        public Game() {
            _board = new Board("[rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1]");
            _viewSide = Side.White;

            Piecies = new List<Piece>();
            SelectedPiece = null;
            
            PossibleMoves = new List<Move>();
            LastMoveSquares = new List<Square>();

            PromotionMove = null;

            Piecies = _board.Piecies;
            SideToMove = Side.White;
            UpdatePossibleMoves();
        }

        private void ClearSelectedSquare()
        {
            SelectedPiece = null;
            PossibleMoves = new List<Move>();
        }

        public void SelectPiece(Piece p)
        {
            if (PromotionMove is not null) return;

            if (SelectedPiece is null || SelectedPiece.Square != p.Square)
            {
                SelectedPiece = p;
                PossibleMoves = p.PossibleMoves;
            }
            else
            {
                ClearSelectedSquare();
            }
        }

        public void PlayMove(Move m)
        {
            if (SideToMove != m.Player) return;

            if (m is PromotionMove pm && pm.PromotionPiece is null) {
                PromotionMove = pm;
                return;
            }

            _board.PlayMove(m);
            _board[m.End].MoveTo(m.End);

            if (m is CastlingMove cm)
            {
                int post_rook_offset = cm.CastlingType == Castling.King ? +1 : -1;
                _board[m.Start.Rank, m.Start.File+post_rook_offset].MoveTo(new Square(m.Start.Rank, m.Start.File + post_rook_offset));
            }

            ClearSelectedSquare();
            Piecies = _board.Piecies;
            UpdatePossibleMoves();
            SetLastMoveSquares(m);

            SideToMove = SideToMove == Side.White ? Side.Black : Side.White;
        }

        public void ChoosePromotion(string type)
        {
            switch (type) {
                case "queen": 
                    PromotionMove.PromotionPiece = new Queen(SideToMove, PromotionMove.End.X, PromotionMove.End.Y);
                    break;
                case "rook":
                    PromotionMove.PromotionPiece = new Rook(SideToMove, PromotionMove.End.X, PromotionMove.End.Y);
                    break;
                case "bishop":
                    PromotionMove.PromotionPiece = new Bishop(SideToMove, PromotionMove.End.X, PromotionMove.End.Y);
                    break;
                case "knight":
                    PromotionMove.PromotionPiece = new Knight(SideToMove, PromotionMove.End.X, PromotionMove.End.Y);
                    break;
                default:
                    break;
            }
            
            PlayMove(PromotionMove);
            PromotionMove = null;
        }

        private void UpdatePossibleMoves()
        {
            foreach (var p in Piecies)
            {
                Func<Move, bool> moveIsLegal = move =>
                {
                    _board.PlayMove(move);
                    bool isCheck = _board.IsKingInCheck(p.Color);
                    _board.UndoMove(move);

                    return !isCheck;
                };

                p.PossibleMoves = p.GetMoves(_board).Where(moveIsLegal).ToList<Move>();
            }
        }

        private void SetLastMoveSquares(Move m)
        {
            LastMoveSquares.Clear();
            LastMoveSquares.Add(m.Start);
            LastMoveSquares.Add(m.End);
        }

    }
  
}
