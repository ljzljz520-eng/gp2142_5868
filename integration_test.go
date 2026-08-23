package othello_test

import (
	"path/filepath"
	"testing"

	"example.com/othello/internal/game"
	"example.com/othello/internal/importer"
	"example.com/othello/internal/records"
	"example.com/othello/internal/service"
	"example.com/othello/internal/store"
)

func TestWorkflowCreateReviewConfirmArchive(t *testing.T) {
	database, err := store.Open(filepath.Join(t.TempDir(), "workflow.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	application, err := service.New(database)
	if err != nil {
		t.Fatal(err)
	}
	state, err := service.FinishedFixture("integration-record")
	if err != nil {
		t.Fatal(err)
	}
	record, err := application.RegisterGame(state, "east", "player")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := application.ReviewRecord(record.ID, "judge", "all moves checked"); err != nil {
		t.Fatal(err)
	}
	if _, err := application.ConfirmRecord(record.ID, "judge", "accepted"); err != nil {
		t.Fatal(err)
	}
	if _, err := application.PublishRecord(record.ID, "judge", "published"); err != nil {
		t.Fatal(err)
	}
	if _, err := application.ArchiveRecord(record.ID, "judge", "stored"); err != nil {
		t.Fatal(err)
	}
	status, err := application.Status(record.ID)
	if err != nil || status != "archived" {
		t.Fatalf("status=%s err=%v", status, err)
	}
}

func TestWorkflowSearchUpdatePublish(t *testing.T) {
	database, err := store.Open(filepath.Join(t.TempDir(), "history.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	application, err := service.New(database)
	if err != nil {
		t.Fatal(err)
	}
	state, err := service.FinishedFixture("history-record")
	if err != nil {
		t.Fatal(err)
	}
	record, err := application.RegisterGame(state, "west", "player")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := application.UpdateRecord(record.ID, "first", records.UpdateRequest{Notes: "first"}); err != nil {
		t.Fatal(err)
	}
	if _, err := application.UpdateRecord(record.ID, "second", records.UpdateRequest{Notes: "second"}); err != nil {
		t.Fatal(err)
	}
	if _, err := application.ReviewRecord(record.ID, "judge", "reviewed"); err != nil {
		t.Fatal(err)
	}
	if _, err := application.ConfirmRecord(record.ID, "judge", "accepted"); err != nil {
		t.Fatal(err)
	}
	if _, err := application.PublishRecord(record.ID, "judge", "published"); err != nil {
		t.Fatal(err)
	}
	history, err := application.Search(records.Filter{Region: "west", Status: "published"})
	if err != nil || history.Total != 1 || history.Records[0].Notes != "first: first | second: second" {
		t.Fatalf("history=%+v err=%v", history, err)
	}
}

func TestWorkflowImportValidatePersistReport(t *testing.T) {
	database, err := store.Open(filepath.Join(t.TempDir(), "import.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	application, err := service.New(database)
	if err != nil {
		t.Fatal(err)
	}
	result, err := application.Import(importer.Request{ID: "import-record", Mode: game.ModeTwoPlayer, Region: "east", Source: "opening", Content: "1.d3 2.c3"}, "importer")
	if err != nil {
		t.Fatal(err)
	}
	if result.Moves != 2 || result.Attachment.ID != "attachment-import-record" || result.Report == "" {
		t.Fatalf("import result=%+v", result)
	}
	attachment, err := database.GetAttachment(result.Attachment.ID)
	if err != nil || attachment.Content != "1.d3 2.c3" {
		t.Fatalf("attachment=%+v err=%v", attachment, err)
	}
}
