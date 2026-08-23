package service

import (
	"path/filepath"
	"testing"

	"example.com/othello/internal/records"
	"example.com/othello/internal/store"
)

func TestServiceUpdatesKeepBothCollaborators(t *testing.T) {
	database, err := store.Open(filepath.Join(t.TempDir(), "service.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	application, err := New(database)
	if err != nil {
		t.Fatal(err)
	}
	state, err := FinishedFixture("service-game")
	if err != nil {
		t.Fatal(err)
	}
	record, err := application.RegisterGame(state, "north", "owner")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := application.UpdateRecord(record.ID, "alice", records.UpdateRequest{Notes: "first change"}); err != nil {
		t.Fatal(err)
	}
	updated, err := application.UpdateRecord(record.ID, "bob", records.UpdateRequest{Notes: "second change"})
	if err != nil {
		t.Fatal(err)
	}
	if updated.Notes != "alice: first change | bob: second change" {
		t.Fatalf("notes = %q", updated.Notes)
	}
}

func TestBusiness037Regression(t *testing.T) {
	database, err := store.Open(filepath.Join(t.TempDir(), "regression.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	application, err := New(database)
	if err != nil {
		t.Fatal(err)
	}
	state, err := FinishedFixture("K-137")
	if err != nil {
		t.Fatal(err)
	}
	record, err := application.RegisterGame(state, "north", "owner")
	if err != nil {
		t.Fatal(err)
	}
	history, err := application.Search(records.Filter{Region: "north"})
	if err != nil || history.Total != 1 {
		t.Fatalf("history = %+v err=%v", history, err)
	}
	if _, err := application.UpdateRecord(history.Records[0].ID, "alice", records.UpdateRequest{Notes: "first change"}); err != nil {
		t.Fatal(err)
	}
	if _, err := application.UpdateRecord(history.Records[0].ID, "bob", records.UpdateRequest{Notes: "second change"}); err != nil {
		t.Fatal(err)
	}
	if _, err := application.ReviewRecord(record.ID, "reviewer", "changes checked"); err != nil {
		t.Fatal(err)
	}
	if _, err := application.ConfirmRecord(record.ID, "reviewer", ""); err == nil {
		t.Fatalf("confirmation without a review note was accepted")
	}
}
