using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using Timer = System.Timers.Timer;
using System.Text.RegularExpressions;

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

        public async Task PlayMove(Move m)
        {
            if (SideToMove != m.Player) return;

            if (m is PromotionMove pm && pm.PromotionPiece is null) {
                PromotionMove = pm;
                return;
            }

            string notation = m.GenerateNotation(_board);
            await Application.Instance.SendMove(notation);
        }

        public async Task ChoosePromotion(string type)
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

            await PlayMove(PromotionMove);
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


        public void PlayReceivedMove(string color, string move, int time) {
            if (!ticker.Enabled) ticker.Start();
            _PlayerSeconds[PlayerToMove] = time;

            Move m = ParseReceivedMove(color, move);

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
            PlayerToMove = (PlayerToMove + 1) % 2;
        }

        private Move ParseReceivedMove(string color, string move) {
            try
            {
                const string pattern = "(?:(?<castle>^[O0o]-[O0o](?<long>-[O0o])?)|" +
                                     "(?<pawn_move>^(?<pawn_file_start>[a-h])" +
                                     "(?:x(?<pawn_file_end>[a-h]))?(?<pawn_rank_end>[1-8])" +
                                     "(?:=(?<promotion_piece>[QNRB]))?)|" +
                                     "(?<move>(?<piece>[RBQKPN])" +
                                     "(?<file_start>[a-h])?(?<rank_start>[1-8])?" +
                                     "(?<capture>[x])?(?<file_end>[a-h])(?<rank_end>[1-8])))" +
                                     "(?:[+#])?$";

                Regex regex = new Regex(pattern);
                Match match = regex.Match(move);
                if (!match.Success) return null;

                Side side = color == "white" ? Side.White : Side.Black;

                if (match.Groups["castle"].Success)
                {
                    Castling ctype = match.Groups["long"].Success ? Castling.Queen : Castling.King;
                    return Piecies.OfType<King>().Where(k => k.Color == side).First().PossibleMoves.OfType<CastlingMove>().Where(cm => cm.CastlingType == ctype).First();
                }

                if (match.Groups["pawn_move"].Success)
                {
                    int fileStart = Square.FileToInt(match.Groups["pawn_file_start"].Value);
                    int rankEnd = Square.RankToInt(match.Groups["pawn_rank_end"].Value);
                    bool isCapture = match.Groups["pawn_file_end"].Success;
                    string fileEndS = match.Groups["pawn_file_end"].Value;

                    Move rMove = Piecies.OfType<Pawn>()
                            .Where(p => p.Color == side && p.Square.File == fileStart)
                            .SelectMany(p => p.PossibleMoves)
                            .Where(m => m.End.Rank == rankEnd && (!isCapture || Square.FileToInt(fileEndS) == m.End.File))
                            .First();

                    System.Diagnostics.Debug.Assert(rMove != null);

                    if (rMove is PromotionMove promMove)
                    {
                        System.Diagnostics.Debug.Assert(match.Groups["promotion"].Success);
                        string promotionPieceType = match.Groups["promotion"].Value;
                        switch (promotionPieceType)
                        {
                            case "Q":
                                promMove.PromotionPiece = new Queen(side, promMove.End.X, promMove.End.Y);
                                break;
                            case "R":
                                promMove.PromotionPiece = new Rook(side, promMove.End.X, promMove.End.Y);
                                break;
                            case "B":
                                promMove.PromotionPiece = new Bishop(side, promMove.End.X, promMove.End.Y);
                                break;
                            case "N":
                                promMove.PromotionPiece = new Knight(side, promMove.End.X, promMove.End.Y);
                                break;
                        }
                    }

                    return rMove;
                }

                if (match.Groups["move"].Success)
                {
                    Type pieceType = typeof(Piece);
                    switch (match.Groups["piece"].Value)
                    {
                        case "Q":
                            pieceType = typeof(Queen);
                            break;
                        case "R":
                            pieceType = typeof(Rook);
                            break;
                        case "B":
                            pieceType = typeof(Bishop);
                            break;
                        case "N":
                            pieceType = typeof(Knight);
                            break;
                        case "K":
                            pieceType = typeof(King);
                            break;
                    }

                    string endPosition = match.Groups["file_end"].Value + match.Groups["rank_end"].Value;

                    Func<Move, bool> filterMoves = (m) =>
                    {
                        bool fileStartSpecified = match.Groups["file_start"].Success;
                        bool rankStartSpecified = match.Groups["rank_start"].Success;

                        return m.End.Position == endPosition &&
                                (!fileStartSpecified || (fileStartSpecified && m.Start.File == Square.FileToInt(match.Groups["file_start"].Value))) &&
                                (!rankStartSpecified || (rankStartSpecified && m.Start.Rank == Square.RankToInt(match.Groups["rank_start"].Value)));
                    };

                    return Piecies.Where(p => p.Color == side && p.GetType() == pieceType)
                        .SelectMany(p => p.PossibleMoves).Where(filterMoves).First();
                }
            }
            catch (Exception ex) {
                System.Diagnostics.Debug.WriteLine(ex.ToString());
            }
            return null;
        }
    }
  
}
