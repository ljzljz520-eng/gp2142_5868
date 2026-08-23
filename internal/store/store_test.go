package store

import (
	"path/filepath"
	"testing"
)

func TestStorePersistsAllEntities(t *testing.T) {
	path := filepath.Join(t.TempDir(), "records.db")
	database, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	record := Record{ID: "record-1", Status: "draft", Region: "north", Version: 1}
	if err := database.SaveRecord(record); err != nil {
		t.Fatal(err)
	}
	if err := database.SaveWorkflow(Workflow{ID: "workflow-record-1", RecordID: record.ID, State: "draft", Version: 1}); err != nil {
		t.Fatal(err)
	}
	if err := database.SaveAuditEvent(AuditEvent{ID: "audit-1", RecordID: record.ID, Sequence: 1}); err != nil {
		t.Fatal(err)
	}
	if err := database.SaveAttachment(Attachment{ID: "attachment-1", RecordID: record.ID, Content: "d3"}); err != nil {
		t.Fatal(err)
	}
	got, err := database.GetRecord(record.ID)
	if err != nil || got.ID != record.ID {
		t.Fatalf("record = %+v err=%v", got, err)
	}
	attachments, err := database.ListAttachments(record.ID)
	if err != nil || len(attachments) != 1 {
		t.Fatalf("attachments = %+v err=%v", attachments, err)
	}
}
