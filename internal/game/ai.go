package game

import (
	"sort"

	"example.com/othello/internal/board"
)

type MoveChoice struct {
	Position board.Position
	Score    int
	Reason   string
}

func ChooseComputerMove(state Game) (MoveChoice, bool) {
	moves := state.LegalMoves()
	if len(moves) == 0 {
		return MoveChoice{}, false
	}
	choices := make([]MoveChoice, 0, len(moves))
	for _, position := range moves {
		captures := board.Captures(state.Board, position, state.Current)
		score := len(captures) * 10
		if position.Row == 0 || position.Row == board.Size-1 {
			score += 3
		}
		if position.Col == 0 || position.Col == board.Size-1 {
			score += 3
		}
		if (position.Row == 0 || position.Row == board.Size-1) && (position.Col == 0 || position.Col == board.Size-1) {
			score += 20
		}
		choices = append(choices, MoveChoice{Position: position, Score: score, Reason: "captures most pieces with edge preference"})
	}
	sort.Slice(choices, func(i, j int) bool {
		if choices[i].Score == choices[j].Score {
			if choices[i].Position.Row == choices[j].Position.Row {
				return choices[i].Position.Col < choices[j].Position.Col
			}
			return choices[i].Position.Row < choices[j].Position.Row
		}
		return choices[i].Score > choices[j].Score
	})
	return choices[0], true
}

func (g *Game) PlayComputerTurn() (Move, error) {
	if g.Mode != ModeComputer {
		return Move{}, ErrComputerModeOnly
	}
	choice, ok := ChooseComputerMove(*g)
	if !ok {
		if err := g.Pass(); err != nil {
			return Move{}, err
		}
		return Move{}, nil
	}
	return g.PlayPosition(choice.Position)
}
