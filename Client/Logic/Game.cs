﻿using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using System.Timers;

namespace Client.Logic
{
    public class Game
    {
        public List<Piece> Piecies;
        public Piece SelectedPiece;
        public List<Move> PossibleMoves;
        public List<Square> LastMoveSquares;

        public Side PlayingSide;
        private PromotionMove PromotionMove;

        public bool DisplayPromotionModal { get => PromotionMove is not null; }
        public User opponent;

        public Timer ticker;
        public int[] _PlayerSeconds;
        public string TimerPlayerA { get => TimeSpan.FromSeconds(_PlayerSeconds[0]).ToString(@"hh\:mm\:ss").TrimStart(' ', '0', ':'); }
        public string TimerPlayerB { get => TimeSpan.FromSeconds(_PlayerSeconds[1]).ToString(@"hh\:mm\:ss").TrimStart(' ', '0', ':'); }


        private Board _board;
        private Side SideToMove;
        public int PlayerToMove;

        int HalfmoveClock;
        int FullmoveClock;

        public Game(Side playingSide, User opponent, int time) {
            _board = new Board("[rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1]");

            Piecies = new List<Piece>();
            SelectedPiece = null;

            PossibleMoves = new List<Move>();
            LastMoveSquares = new List<Square>();

            PromotionMove = null;

            Piecies = _board.Piecies;
            SideToMove = Side.White;
            UpdatePossibleMoves();


            // ----

            PlayingSide = playingSide;
            PlayerToMove = playingSide == SideToMove ? 0 : 1;
            this.opponent = opponent;
            _PlayerSeconds = new int[2] { time, time };
            ticker = new System.Timers.Timer();
            ticker.Interval = 1000;
            ticker.Elapsed += UpdateTimers;

            ticker.Start();
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
                _board[m.Start.Rank, m.Start.File + post_rook_offset].MoveTo(new Square(m.Start.Rank, m.Start.File + post_rook_offset));
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


        private void UpdateTimers(Object source, System.Timers.ElapsedEventArgs e) {
            _PlayerSeconds[PlayerToMove] -= 1;

            if (_PlayerSeconds[PlayerToMove] == 0)
                ((Timer)source).Enabled = false;

            Application.Instance.updateUI.Invoke();
        }

    }
  
}
