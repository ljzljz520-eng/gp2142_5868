package workflow

import (
	"fmt"
	"strings"

	"example.com/othello/internal/store"
)

type AuditSummary struct {
	RecordID string
	Count    int
	Actors   []string
	Actions  []string
}

func (m *Manager) Audit(recordID string) (AuditSummary, error) {
	events, err := m.Store.ListAuditEvents(recordID)
	if err != nil {
		return AuditSummary{}, err
	}
	actors := make([]string, 0, len(events))
	actions := make([]string, 0, len(events))
	for _, event := range events {
		actors = append(actors, event.Actor)
		actions = append(actions, event.Action)
	}
	return AuditSummary{RecordID: recordID, Count: len(events), Actors: actors, Actions: actions}, nil
}

func FormatAudit(events []store.AuditEvent) string {
	lines := make([]string, 0, len(events))
	for _, event := range events {
		lines = append(lines, fmt.Sprintf("%02d %s %s %s->%s %s", event.Sequence, event.Actor, event.Action, event.FromStatus, event.ToStatus, strings.TrimSpace(event.Note)))
	}
	return strings.Join(lines, "\n")
}
