package store

import (
	"path/filepath"
	"testing"
)

func TestPersistenceSurvivesReopen(t *testing.T) {
	path := filepath.Join(t.TempDir(), "reopen.db")
	database, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := database.SaveRecord(Record{ID: "record-reopen", Status: "published"}); err != nil {
		t.Fatal(err)
	}
	if err := database.SaveWorkflow(Workflow{ID: "workflow-record-reopen", RecordID: "record-reopen", State: "published"}); err != nil {
		t.Fatal(err)
	}
	if err := database.SaveAuditEvent(AuditEvent{ID: "audit-reopen", RecordID: "record-reopen", Sequence: 1}); err != nil {
		t.Fatal(err)
	}
	if err := database.SaveAttachment(Attachment{ID: "attachment-reopen", RecordID: "record-reopen", Content: "d3 c3"}); err != nil {
		t.Fatal(err)
	}
	if err := database.Close(); err != nil {
		t.Fatal(err)
	}
	database, err = Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	if _, err := database.GetRecord("record-reopen"); err != nil {
		t.Fatal(err)
	}
	workflow, err := database.GetWorkflow("workflow-record-reopen")
	if err != nil || workflow.State != "published" {
		t.Fatalf("workflow = %+v err=%v", workflow, err)
	}
	events, err := database.ListAuditEvents("record-reopen")
	if err != nil || len(events) != 1 {
		t.Fatalf("events = %+v err=%v", events, err)
	}
	attachments, err := database.ListAttachments("record-reopen")
	if err != nil || len(attachments) != 1 {
		t.Fatalf("attachments = %+v err=%v", attachments, err)
	}
}
