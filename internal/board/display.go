package board

import (
	"fmt"
	"strings"
)

func Render(position Board) string {
	lines := []string{"  a b c d e f g h"}
	for row, cells := range position.Rows() {
		lines = append(lines, fmt.Sprintf("%d %s", row+1, strings.Join(strings.Split(cells, ""), " ")))
	}
	black, white := position.Counts()
	return strings.Join(lines, "\n") + fmt.Sprintf("\nblack=%d white=%d", black, white)
}

func Coordinates(values []Position) []string {
	result := make([]string, 0, len(values))
	for _, value := range values {
		result = append(result, value.String())
	}
	return result
}

func ExplainIllegal(position Board, coordinate string, piece Piece) string {
	parsed, err := ParseCoordinate(coordinate)
	if err != nil {
		return err.Error()
	}
	if err := ValidateMove(position, parsed, piece); err != nil {
		return err.Error()
	}
	return "coordinate is legal"
}

func PieceLegend() string {
	return fmt.Sprintf("%c=black %c=white %c=empty", Black, White, Empty)
}
