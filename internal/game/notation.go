package game

import (
	"errors"
	"fmt"
	"strings"

	"example.com/othello/internal/board"
)

var ErrComputerModeOnly = errors.New("computer turn requires computer mode")

func ParseMode(value string) (Mode, error) {
	clean := strings.ToLower(strings.TrimSpace(value))
	if clean == "two" || clean == "two-player" || clean == "human" {
		return ModeTwoPlayer, nil
	}
	if clean == "computer" || clean == "ai" {
		return ModeComputer, nil
	}
	return "", fmt.Errorf("mode %q must be two-player or computer", value)
}

func ParseMoves(id string, mode Mode, moves []string) (Game, error) {
	game, err := New(id, mode)
	if err != nil {
		return Game{}, err
	}
	for _, value := range moves {
		if strings.EqualFold(strings.TrimSpace(value), "pass") {
			if err := game.Pass(); err != nil {
				return Game{}, err
			}
			continue
		}
		if _, err := game.Play(value); err != nil {
			return Game{}, err
		}
	}
	return game, nil
}

func ValidOpeningMoves() []string {
	game, _ := New("opening", ModeTwoPlayer)
	moves := game.LegalMoves()
	values := make([]string, 0, len(moves))
	for _, move := range moves {
		values = append(values, move.String())
	}
	return values
}

func PieceName(piece board.Piece) string {
	switch piece {
	case board.Black:
		return "black"
	case board.White:
		return "white"
	default:
		return "empty"
	}
}
