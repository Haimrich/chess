package chess

type Square struct {
	rank int
	file int
}

func (s *Square) Translate(offset *Square) *Square {
	newRank := s.rank + offset.rank
	newFile := s.file + offset.file
	if newRank < 0 || newRank >= BOARD_SIZE || newFile < 0 || newFile >= BOARD_SIZE {
		return nil
	}
	return &Square{rank: newRank, file: newFile}
}

func FileToIdx(file string) int {
	return int(file[0] - 'a')
}

func RankToIdx(rank string) int {
	return int(rank[0] - '1')
}

func ParseSquare(square string) Square {
	return Square{rank: int(square[1] - '1'), file: int(square[0] - 'a')}
}
