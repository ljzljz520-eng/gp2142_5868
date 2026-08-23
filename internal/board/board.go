package board

import "strings"

type Piece byte

const (
	Empty Piece = '.'
	Black Piece = 'B'
	White Piece = 'W'
)

type Board struct {
	Cells [Size][Size]Piece
}

func NewBoard() Board {
	board := Board{}
	for row := 0; row < Size; row++ {
		for col := 0; col < Size; col++ {
			board.Cells[row][col] = Empty
		}
	}
	board.Cells[3][3] = White
	board.Cells[3][4] = Black
	board.Cells[4][3] = Black
	board.Cells[4][4] = White
	return board
}

func BoardFromRows(rows []string) (Board, error) {
	board := Board{}
	if len(rows) != Size {
		return board, ErrInvalidShape
	}
	for row, line := range rows {
		if len(line) != Size {
			return Board{}, ErrInvalidShape
		}
		for col := 0; col < Size; col++ {
			piece := Piece(line[col])
			if piece != Empty && piece != Black && piece != White {
				return Board{}, ErrInvalidPiece
			}
			board.Cells[row][col] = piece
		}
	}
	return board, nil
}

func (b Board) Clone() Board {
	return b
}

func (b Board) PieceAt(position Position) Piece {
	if !position.InBounds() {
		return Empty
	}
	return b.Cells[position.Row][position.Col]
}

func (b *Board) Set(position Position, piece Piece) {
	if position.InBounds() {
		b.Cells[position.Row][position.Col] = piece
	}
}

func (b Board) Count(piece Piece) int {
	count := 0
	for row := 0; row < Size; row++ {
		for col := 0; col < Size; col++ {
			if b.Cells[row][col] == piece {
				count++
			}
		}
	}
	return count
}

func (b Board) Counts() (int, int) {
	return b.Count(Black), b.Count(White)
}

func (b Board) Full() bool {
	return b.Count(Empty) == 0
}

func (b Board) Rows() []string {
	rows := make([]string, 0, Size)
	for row := 0; row < Size; row++ {
		var builder strings.Builder
		for col := 0; col < Size; col++ {
			builder.WriteByte(byte(b.Cells[row][col]))
		}
		rows = append(rows, builder.String())
	}
	return rows
}

func (b Board) String() string {
	return strings.Join(b.Rows(), "\n")
}
