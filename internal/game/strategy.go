package game

import (
	"fmt"
	"sort"

	"example.com/othello/internal/board"
)

type Candidate struct {
	Coordinate string
	Flips      int
	Edge       bool
	Corner     bool
	Score      int
}

func RankMoves(state Game) []Candidate {
	values := make([]Candidate, 0)
	for _, position := range state.LegalMoves() {
		flips := len(board.Captures(state.Board, position, state.Current))
		edge := position.Row == 0 || position.Row == board.Size-1 || position.Col == 0 || position.Col == board.Size-1
		corner := (position.Row == 0 || position.Row == board.Size-1) && (position.Col == 0 || position.Col == board.Size-1)
		score := flips
		if edge {
			score += 2
		}
		if corner {
			score += 8
		}
		values = append(values, Candidate{Coordinate: position.String(), Flips: flips, Edge: edge, Corner: corner, Score: score})
	}
	sort.Slice(values, func(i, j int) bool {
		if values[i].Score == values[j].Score {
			return values[i].Coordinate < values[j].Coordinate
		}
		return values[i].Score > values[j].Score
	})
	return values
}

func StrategySummary(state Game) string {
	choices := RankMoves(state)
	if len(choices) == 0 {
		return fmt.Sprintf("%s has no legal choice", PieceName(state.Current))
	}
	best := choices[0]
	return fmt.Sprintf("%s prefers %s with score %d and %d flip(s)", PieceName(state.Current), best.Coordinate, best.Score, best.Flips)
}

func IsCorner(position board.Position) bool {
	return (position.Row == 0 || position.Row == board.Size-1) && (position.Col == 0 || position.Col == board.Size-1)
}

func IsEdge(position board.Position) bool {
	return position.Row == 0 || position.Row == board.Size-1 || position.Col == 0 || position.Col == board.Size-1
}
