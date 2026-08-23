package board

import "errors"

var (
	ErrInvalidShape = errors.New("board must contain eight rows of eight cells")
	ErrInvalidPiece = errors.New("board contains an unknown piece")
	ErrOccupied     = errors.New("coordinate is already occupied")
	ErrNoCapture    = errors.New("move must capture at least one opposing piece")
	ErrWrongTurn    = errors.New("piece does not match the current turn")
)

var directions = [...]Position{
	{Row: -1, Col: -1}, {Row: -1, Col: 0}, {Row: -1, Col: 1},
	{Row: 0, Col: -1}, {Row: 0, Col: 1},
	{Row: 1, Col: -1}, {Row: 1, Col: 0}, {Row: 1, Col: 1},
}

func Opponent(piece Piece) Piece {
	if piece == Black {
		return White
	}
	if piece == White {
		return Black
	}
	return Empty
}

func Captures(board Board, position Position, piece Piece) []Position {
	if !position.InBounds() || piece == Empty || board.PieceAt(position) != Empty {
		return nil
	}
	opponent := Opponent(piece)
	if opponent == Empty {
		return nil
	}
	captures := make([]Position, 0, 8)
	for _, direction := range directions {
		line := make([]Position, 0, Size)
		cursor := position.Offset(direction.Row, direction.Col)
		for cursor.InBounds() && board.PieceAt(cursor) == opponent {
			line = append(line, cursor)
			cursor = cursor.Offset(direction.Row, direction.Col)
		}
		if len(line) > 0 && cursor.InBounds() && board.PieceAt(cursor) == piece {
			captures = append(captures, line...)
		}
	}
	return captures
}

func LegalMoves(board Board, piece Piece) []Position {
	moves := make([]Position, 0)
	for _, position := range AllPositions() {
		if len(Captures(board, position, piece)) > 0 {
			moves = append(moves, position)
		}
	}
	return moves
}

func ValidateMove(board Board, position Position, piece Piece) error {
	if piece != Black && piece != White {
		return ErrWrongTurn
	}
	if !position.InBounds() {
		return errors.New("coordinate is outside the board")
	}
	if board.PieceAt(position) != Empty {
		return ErrOccupied
	}
	if len(Captures(board, position, piece)) == 0 {
		return ErrNoCapture
	}
	return nil
}

func ApplyMove(board Board, position Position, piece Piece) (Board, []Position, error) {
	if err := ValidateMove(board, position, piece); err != nil {
		return board, nil, err
	}
	captures := Captures(board, position, piece)
	board.Set(position, piece)
	for _, captured := range captures {
		board.Set(captured, piece)
	}
	return board, captures, nil
}

func HasMove(board Board, piece Piece) bool {
	return len(LegalMoves(board, piece)) > 0
}

func GameOver(board Board, piece Piece) bool {
	return board.Full() || (!HasMove(board, piece) && !HasMove(board, Opponent(piece)))
}
