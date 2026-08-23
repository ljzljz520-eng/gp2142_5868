package game

import (
	"fmt"
	"strings"

	"example.com/othello/internal/board"
)

type Mode string

const (
	ModeTwoPlayer Mode = "two-player"
	ModeComputer  Mode = "computer"
)

type Status string

const (
	StatusActive Status = "active"
	StatusEnded  Status = "ended"
)

type Move struct {
	Number  int         `json:"number"`
	Player  board.Piece `json:"player"`
	At      string      `json:"at"`
	Flipped int         `json:"flipped"`
}

type Game struct {
	ID          string      `json:"id"`
	Mode        Mode        `json:"mode"`
	Status      Status      `json:"status"`
	Current     board.Piece `json:"current"`
	Board       board.Board `json:"board"`
	Moves       []Move      `json:"moves"`
	BlackScore  int         `json:"black_score"`
	WhiteScore  int         `json:"white_score"`
	Winner      board.Piece `json:"winner"`
	Passes      int         `json:"passes"`
	LastMessage string      `json:"last_message"`
}

func New(id string, mode Mode) (Game, error) {
	clean := strings.TrimSpace(id)
	if clean == "" {
		return Game{}, fmt.Errorf("game id is required")
	}
	if mode != ModeTwoPlayer && mode != ModeComputer {
		return Game{}, fmt.Errorf("unsupported mode %q", mode)
	}
	game := Game{ID: clean, Mode: mode, Status: StatusActive, Current: board.Black, Board: board.NewBoard(), Moves: []Move{}}
	game.refreshScore()
	return game, nil
}

func FromBoard(id string, mode Mode, current board.Piece, position board.Board) (Game, error) {
	game, err := New(id, mode)
	if err != nil {
		return Game{}, err
	}
	if current != board.Black && current != board.White {
		return Game{}, fmt.Errorf("current piece must be black or white")
	}
	game.Board = position
	game.Current = current
	game.refreshScore()
	return game, nil
}

func (g *Game) refreshScore() {
	g.BlackScore, g.WhiteScore = g.Board.Counts()
}

func (g Game) LegalMoves() []board.Position {
	if g.Status != StatusActive {
		return nil
	}
	return board.LegalMoves(g.Board, g.Current)
}

func (g *Game) Play(value string) (Move, error) {
	position, err := board.ParseCoordinate(value)
	if err != nil {
		return Move{}, err
	}
	return g.PlayPosition(position)
}

func (g *Game) PlayPosition(position board.Position) (Move, error) {
	if g.Status != StatusActive {
		return Move{}, fmt.Errorf("game %s is already ended", g.ID)
	}
	next, flipped, err := board.ApplyMove(g.Board, position, g.Current)
	if err != nil {
		return Move{}, fmt.Errorf("illegal move at %s: %w", position, err)
	}
	g.Board = next
	move := Move{Number: len(g.Moves) + 1, Player: g.Current, At: position.String(), Flipped: len(flipped)}
	g.Moves = append(g.Moves, move)
	g.Passes = 0
	g.refreshScore()
	g.advance()
	return move, nil
}

func (g *Game) advance() {
	next := board.Opponent(g.Current)
	if board.GameOver(g.Board, next) {
		g.Status = StatusEnded
		g.Current = next
		g.Winner = g.winner()
		g.LastMessage = "game ended"
		return
	}
	if !board.HasMove(g.Board, next) {
		g.Current = board.Opponent(next)
		g.Passes++
		g.LastMessage = fmt.Sprintf("%c has no legal move and passes", next)
		return
	}
	g.Current = next
	g.LastMessage = fmt.Sprintf("%c to move", g.Current)
}

func (g *Game) Pass() error {
	if g.Status != StatusActive {
		return fmt.Errorf("game %s is already ended", g.ID)
	}
	if board.HasMove(g.Board, g.Current) {
		return fmt.Errorf("%c has a legal move and cannot pass", g.Current)
	}
	g.Passes++
	if g.Passes >= 2 {
		g.Status = StatusEnded
		g.Winner = g.winner()
		g.LastMessage = "both players passed"
		return nil
	}
	g.Current = board.Opponent(g.Current)
	return nil
}

func (g Game) winner() board.Piece {
	if g.BlackScore > g.WhiteScore {
		return board.Black
	}
	if g.WhiteScore > g.BlackScore {
		return board.White
	}
	return board.Empty
}

func (g Game) Summary() string {
	result := "draw"
	if g.Winner == board.Black {
		result = "black wins"
	} else if g.Winner == board.White {
		result = "white wins"
	}
	return fmt.Sprintf("%s: %s, black %d white %d", g.ID, result, g.BlackScore, g.WhiteScore)
}

func (g Game) Transcript() string {
	parts := make([]string, 0, len(g.Moves))
	for _, move := range g.Moves {
		parts = append(parts, fmt.Sprintf("%d.%c%s", move.Number, move.Player, move.At))
	}
	return strings.Join(parts, " ")
}
