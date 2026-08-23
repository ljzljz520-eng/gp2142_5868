package board

import (
	"fmt"
	"strings"
)

type Analysis struct {
	Piece        Piece
	Legal        []Position
	Capturable   int
	Occupied     int
	Empty        int
	HasPass      bool
	GameFinished bool
}

func Analyze(position Board, piece Piece) Analysis {
	legal := LegalMoves(position, piece)
	occupied := Size*Size - position.Count(Empty)
	return Analysis{Piece: piece, Legal: legal, Capturable: len(legal), Occupied: occupied, Empty: Size*Size - occupied, HasPass: len(legal) == 0 && HasMove(position, Opponent(piece)), GameFinished: GameOver(position, piece)}
}

func ValidatePosition(position Board) error {
	black, white := position.Counts()
	if black == 0 && white == 0 {
		return fmt.Errorf("board cannot be empty")
	}
	if black > Size*Size || white > Size*Size || black+white > Size*Size {
		return fmt.Errorf("board piece count exceeds capacity")
	}
	return nil
}

func EmptyPositions(position Board) []Position {
	values := make([]Position, 0, position.Count(Empty))
	for _, candidate := range AllPositions() {
		if position.PieceAt(candidate) == Empty {
			values = append(values, candidate)
		}
	}
	return values
}

func DescribeAnalysis(analysis Analysis) string {
	legal := make([]string, 0, len(analysis.Legal))
	for _, position := range analysis.Legal {
		legal = append(legal, position.String())
	}
	status := "active"
	if analysis.GameFinished {
		status = "finished"
	} else if analysis.HasPass {
		status = "pass available"
	}
	return fmt.Sprintf("%c: %s, legal=%s, empty=%d", analysis.Piece, status, strings.Join(legal, ","), analysis.Empty)
}
