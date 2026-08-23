package board

import (
	"fmt"
	"strconv"
	"strings"
)

const Size = 8

type Position struct {
	Row int
	Col int
}

func NewPosition(row, col int) (Position, error) {
	position := Position{Row: row, Col: col}
	if !position.InBounds() {
		return Position{}, fmt.Errorf("position %d,%d is outside the board", row, col)
	}
	return position, nil
}

func ParseCoordinate(value string) (Position, error) {
	clean := strings.TrimSpace(strings.ToLower(value))
	if len(clean) < 2 {
		return Position{}, fmt.Errorf("coordinate %q must look like d3", value)
	}
	column := clean[0]
	if column < 'a' || column > 'h' {
		return Position{}, fmt.Errorf("column %q must be between a and h", string(column))
	}
	row, err := strconv.Atoi(clean[1:])
	if err != nil || row < 1 || row > Size {
		return Position{}, fmt.Errorf("row in %q must be between 1 and 8", value)
	}
	return Position{Row: row - 1, Col: int(column - 'a')}, nil
}

func (p Position) InBounds() bool {
	return p.Row >= 0 && p.Row < Size && p.Col >= 0 && p.Col < Size
}

func (p Position) String() string {
	if !p.InBounds() {
		return "off-board"
	}
	return fmt.Sprintf("%c%d", 'a'+byte(p.Col), p.Row+1)
}

func (p Position) Offset(row, col int) Position {
	return Position{Row: p.Row + row, Col: p.Col + col}
}

func AllPositions() []Position {
	positions := make([]Position, 0, Size*Size)
	for row := 0; row < Size; row++ {
		for col := 0; col < Size; col++ {
			positions = append(positions, Position{Row: row, Col: col})
		}
	}
	return positions
}
