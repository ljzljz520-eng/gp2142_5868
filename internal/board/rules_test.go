package board

import (
	"errors"
	"testing"
)

func TestApplyMoveFlipsLine(t *testing.T) {
	board := NewBoard()
	position, err := ParseCoordinate("d3")
	if err != nil {
		t.Fatal(err)
	}
	updated, captures, err := ApplyMove(board, position, Black)
	if err != nil {
		t.Fatal(err)
	}
	if len(captures) != 1 || updated.PieceAt(position) != Black || updated.Count(Black) != 4 {
		t.Fatalf("unexpected move result: captures=%d black=%d", len(captures), updated.Count(Black))
	}
}

func TestIllegalMoveExplainsReason(t *testing.T) {
	board := NewBoard()
	position, _ := ParseCoordinate("a1")
	if !errors.Is(ValidateMove(board, position, Black), ErrNoCapture) {
		t.Fatalf("expected no-capture error")
	}
	occupied, _ := ParseCoordinate("d4")
	if !errors.Is(ValidateMove(board, occupied, Black), ErrOccupied) {
		t.Fatalf("expected occupied error")
	}
}
