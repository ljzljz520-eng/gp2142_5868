package cli

import (
	"fmt"
	"strings"

	"example.com/othello/internal/board"
	"example.com/othello/internal/game"
	"example.com/othello/internal/store"
)

func renderBoard(state game.Game) string {
	black, white := state.Board.Counts()
	return board.Render(state.Board) + fmt.Sprintf(" turn=%s status=%s", game.PieceName(state.Current), state.Status) + fmt.Sprintf("\n%s", board.PieceLegend()) + fmt.Sprintf("\nscore=%d:%d", black, white)
}

func renderMove(move game.Move) string {
	if move.At == "" {
		return "computer passed"
	}
	return fmt.Sprintf("move %d: %s at %s, flipped %d", move.Number, game.PieceName(move.Player), move.At, move.Flipped)
}

func renderRecord(record store.Record) string {
	return fmt.Sprintf("%s [%s] region=%s black=%d white=%d %s notes=%s", record.ID, record.Status, record.Region, record.BlackScore, record.WhiteScore, record.Summary, record.Notes)
}

func renderHelp() string {
	return strings.TrimSpace(`othello commands:
  start <game-id> <two-player|computer>
  play <game-id> <two-player|computer> <coordinate>...
  history <database> [region]
  import <database> <game-id> <two-player|computer> <region> <coordinate>...
  demo <database>`)
}

func coordinateList(values []board.Position) string {
	return strings.Join(board.Coordinates(values), ",")
}
