package service

import (
	"example.com/othello/internal/board"
	"example.com/othello/internal/game"
)

func FinishedFixture(id string) (game.Game, error) {
	position, err := board.BoardFromRows([]string{
		"BWW.BBBB",
		"BBBBBBBB",
		"BBBBBBBB",
		"BBBBBBBB",
		"BBBBBBBB",
		"BBBBBBBB",
		"BBBBBBBB",
		"BBBBBBBB",
	})
	if err != nil {
		return game.Game{}, err
	}
	state, err := game.FromBoard(id, game.ModeTwoPlayer, board.Black, position)
	if err != nil {
		return game.Game{}, err
	}
	if _, err := state.Play("d1"); err != nil {
		return game.Game{}, err
	}
	return state, nil
}
