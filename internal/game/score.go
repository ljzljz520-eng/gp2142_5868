package game

import (
	"fmt"
	"strings"

	"example.com/othello/internal/board"
)

type Result struct {
	GameID      string
	Status      Status
	Winner      board.Piece
	BlackScore  int
	WhiteScore  int
	Difference  int
	MoveCount   int
	Transcript  string
	Explanation string
}

func (g Game) Result() Result {
	difference := g.BlackScore - g.WhiteScore
	if difference < 0 {
		difference = -difference
	}
	explanation := "game is still active"
	if g.Status == StatusEnded {
		explanation = "final score recorded"
	}
	return Result{GameID: g.ID, Status: g.Status, Winner: g.Winner, BlackScore: g.BlackScore, WhiteScore: g.WhiteScore, Difference: difference, MoveCount: len(g.Moves), Transcript: g.Transcript(), Explanation: explanation}
}

func (g Game) ScoreLine() string {
	result := g.Result()
	return fmt.Sprintf("%s %s %d:%d (%d moves)", result.GameID, result.Explanation, result.BlackScore, result.WhiteScore, result.MoveCount)
}

func (g Game) NeedsPass() bool {
	return g.Status == StatusActive && len(g.LegalMoves()) == 0
}

func (g Game) HasMoves() bool {
	return g.Status == StatusActive && len(g.LegalMoves()) > 0
}

func (g Game) Clone() Game {
	clone := g
	clone.Moves = append([]Move(nil), g.Moves...)
	return clone
}

func Replay(id string, mode Mode, transcript string) (Game, error) {
	values := strings.Fields(strings.TrimSpace(transcript))
	return ParseMoves(id, mode, values)
}

func CompareResults(left, right Result) int {
	if left.Difference != right.Difference {
		if left.Difference > right.Difference {
			return 1
		}
		return -1
	}
	if left.MoveCount > right.MoveCount {
		return 1
	}
	if left.MoveCount < right.MoveCount {
		return -1
	}
	return 0
}
