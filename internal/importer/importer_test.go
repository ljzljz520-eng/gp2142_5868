package importer

import (
	"path/filepath"
	"testing"

	"example.com/othello/internal/game"
	"example.com/othello/internal/store"
)

func TestParseValidateAndSaveTranscript(t *testing.T) {
	request := Request{ID: "import-game", Mode: game.ModeTwoPlayer, Source: "lesson-1", Content: "1.d3 2.c3"}
	parsed, message, err := Parse(request)
	if err != nil {
		t.Fatal(err)
	}
	database, err := store.Open(filepath.Join(t.TempDir(), "import.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	attachment, err := SaveAttachment(database, request, parsed, message)
	if err != nil {
		t.Fatal(err)
	}
	if err := ValidateAttachment(attachment); err != nil {
		t.Fatal(err)
	}
	report := BuildReport(attachment, parsed, message)
	if !report.Valid || report.MoveCount != 2 || report.AttachmentID != attachment.ID {
		t.Fatalf("report = %+v", report)
	}
}
