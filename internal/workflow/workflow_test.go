package workflow_test

import (
	"path/filepath"
	"testing"

	"example.com/othello/internal/service"
	"example.com/othello/internal/store"
	"example.com/othello/internal/workflow"
)

func TestManagerRunsValidLifecycle(t *testing.T) {
	database, err := store.Open(filepath.Join(t.TempDir(), "workflow.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	application, err := service.New(database)
	if err != nil {
		t.Fatal(err)
	}
	state, err := service.FinishedFixture("workflow-game")
	if err != nil {
		t.Fatal(err)
	}
	record, err := application.RegisterGame(state, "south", "owner")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := application.ReviewRecord(record.ID, "reviewer", "score verified"); err != nil {
		t.Fatal(err)
	}
	if _, err := application.ConfirmRecord(record.ID, "reviewer", "accepted"); err != nil {
		t.Fatal(err)
	}
	if _, err := application.PublishRecord(record.ID, "reviewer", "publish"); err != nil {
		t.Fatal(err)
	}
	if _, err := application.ArchiveRecord(record.ID, "reviewer", "archive"); err != nil {
		t.Fatal(err)
	}
	status, err := application.Status(record.ID)
	if err != nil || status != workflow.StateArchived {
		t.Fatalf("status=%s err=%v", status, err)
	}
}
