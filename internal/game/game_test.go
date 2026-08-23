package game

import "testing"

func TestGamePlayAdvancesAndRecordsTranscript(t *testing.T) {
	state, err := New("game-1", ModeTwoPlayer)
	if err != nil {
		t.Fatal(err)
	}
	move, err := state.Play("d3")
	if err != nil {
		t.Fatal(err)
	}
	if move.Flipped != 1 || state.Current != 'W' || state.Transcript() != "1.Bd3" {
		t.Fatalf("unexpected game state: %+v %s", state, state.Transcript())
	}
}

func TestComputerChoiceIsDeterministic(t *testing.T) {
	state, err := New("computer-1", ModeComputer)
	if err != nil {
		t.Fatal(err)
	}
	choice, ok := ChooseComputerMove(state)
	if !ok || choice.Position.String() != "d3" {
		t.Fatalf("unexpected choice: %+v %t", choice, ok)
	}
	if _, err := state.PlayComputerTurn(); err != nil {
		t.Fatal(err)
	}
}
