package records_test

import (
	"path/filepath"
	"testing"

	"example.com/othello/internal/records"
	"example.com/othello/internal/service"
	"example.com/othello/internal/store"
)

func TestRepositoryRegistersSearchesAndAppendsUpdates(t *testing.T) {
	database, err := store.Open(filepath.Join(t.TempDir(), "records.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	repository, err := records.NewRepository(database)
	if err != nil {
		t.Fatal(err)
	}
	state, err := service.FinishedFixture("records-game")
	if err != nil {
		t.Fatal(err)
	}
	record, err := repository.Register(state, "north", "owner")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := repository.Update(record.ID, "alice", records.UpdateRequest{Notes: "first"}); err != nil {
		t.Fatal(err)
	}
	updated, err := repository.Update(record.ID, "bob", records.UpdateRequest{Notes: "second"})
	if err != nil {
		t.Fatal(err)
	}
	if updated.Version != 3 || updated.Notes != "alice: first | bob: second" {
		t.Fatalf("updated record = %+v", updated)
	}
	history, err := repository.History(records.Filter{Region: "north"})
	if err != nil || history.Total != 1 {
		t.Fatalf("history = %+v err=%v", history, err)
	}
}
