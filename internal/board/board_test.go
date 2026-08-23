package board

import "testing"

func TestNewBoardHasOpeningPieces(t *testing.T) {
	board := NewBoard()
	black, white := board.Counts()
	if black != 2 || white != 2 || board.PieceAt(Position{Row: 3, Col: 3}) != White {
		t.Fatalf("unexpected opening board: black=%d white=%d", black, white)
	}
}

func TestBoardRowsRoundTrip(t *testing.T) {
	want := []string{"BWW.BBBB", "BBBBBBBB", "BBBBBBBB", "BBBBBBBB", "BBBBBBBB", "BBBBBBBB", "BBBBBBBB", "BBBBBBBB"}
	board, err := BoardFromRows(want)
	if err != nil {
		t.Fatal(err)
	}
	got := board.Rows()
	for index := range want {
		if got[index] != want[index] {
			t.Fatalf("row %d = %q", index, got[index])
		}
	}
}
