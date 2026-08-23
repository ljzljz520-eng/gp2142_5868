package game

import (
	"errors"
	"fmt"

	"example.com/othello/internal/board"
)

func ValidateState(state Game) error {
	if state.ID == "" {
		return errors.New("game id is required")
	}
	if state.Mode != ModeTwoPlayer && state.Mode != ModeComputer {
		return fmt.Errorf("unsupported game mode %q", state.Mode)
	}
	if state.Current != board.Black && state.Current != board.White {
		return errors.New("current player must be black or white")
	}
	if err := board.ValidatePosition(state.Board); err != nil {
		return err
	}
	black, white := state.Board.Counts()
	if state.BlackScore != black || state.WhiteScore != white {
		return fmt.Errorf("stored score %d:%d differs from board %d:%d", state.BlackScore, state.WhiteScore, black, white)
	}
	if state.Status != StatusActive && state.Status != StatusEnded {
		return fmt.Errorf("unknown game status %q", state.Status)
	}
	return nil
}

func ValidateMoveHistory(state Game) error {
	if len(state.Moves) == 0 {
		return nil
	}
	for index, move := range state.Moves {
		if move.Number != index+1 {
			return fmt.Errorf("move number %d is out of sequence", move.Number)
		}
		if move.Player != board.Black && move.Player != board.White {
			return fmt.Errorf("move %d has an unknown player", move.Number)
		}
		if move.Flipped < 1 {
			return fmt.Errorf("move %d did not capture a piece", move.Number)
		}
	}
	return nil
}

func LegalMoveCount(state Game) int {
	if state.Status != StatusActive {
		return 0
	}
	return len(state.LegalMoves())
}
