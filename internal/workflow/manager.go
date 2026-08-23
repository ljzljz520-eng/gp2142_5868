package workflow

import (
	"errors"
	"fmt"
	"strings"

	"example.com/othello/internal/store"
)

type Manager struct {
	Store *store.Store
}

func NewManager(database *store.Store) (*Manager, error) {
	if database == nil {
		return nil, errors.New("workflow store is required")
	}
	return &Manager{Store: database}, nil
}

func (m *Manager) Get(recordID string) (store.Workflow, error) {
	return m.Store.GetWorkflow("workflow-" + recordID)
}

func (m *Manager) Move(recordID, actor, target, note string) (store.Workflow, error) {
	if err := ValidateState(target); err != nil {
		return store.Workflow{}, err
	}
	if err := ValidateTransitionInput(target, actor, note); err != nil {
		return store.Workflow{}, err
	}
	workflow, err := m.Get(recordID)
	if err != nil {
		return store.Workflow{}, err
	}
	if !CanMove(workflow.State, target) {
		return store.Workflow{}, fmt.Errorf("cannot move workflow from %s to %s", workflow.State, target)
	}
	previous := workflow.State
	workflow.State = target
	workflow.Reviewer = strings.TrimSpace(actor)
	workflow.ReviewNote = strings.TrimSpace(note)
	workflow.Version++
	if err := m.Store.SaveWorkflow(workflow); err != nil {
		return store.Workflow{}, err
	}
	record, err := m.Store.GetRecord(recordID)
	if err != nil {
		return store.Workflow{}, err
	}
	record.Status = target
	record.Version++
	if target == StateArchived {
		record.Archived = true
	}
	if target == StatePublished {
		record.Published = true
	}
	if err := m.Store.SaveRecord(record); err != nil {
		return store.Workflow{}, err
	}
	if err := m.writeAudit(recordID, actor, "transition", previous, target, note); err != nil {
		return store.Workflow{}, err
	}
	return workflow, nil
}

func (m *Manager) Review(recordID, actor, note string) (store.Workflow, error) {
	return m.Move(recordID, actor, StateReview, note)
}

func (m *Manager) Confirm(recordID, actor, note string) (store.Workflow, error) {
	return m.Move(recordID, actor, StateConfirmed, note)
}

func (m *Manager) Publish(recordID, actor, note string) (store.Workflow, error) {
	return m.Move(recordID, actor, StatePublished, note)
}

func (m *Manager) Archive(recordID, actor, note string) (store.Workflow, error) {
	return m.Move(recordID, actor, StateArchived, note)
}

func (m *Manager) writeAudit(recordID, actor, action, from, to, note string) error {
	events, err := m.Store.ListAuditEvents(recordID)
	if err != nil {
		return err
	}
	event := store.AuditEvent{ID: fmt.Sprintf("audit-%s-%d", recordID, len(events)+1), RecordID: recordID, Actor: actor, Action: action, FromStatus: from, ToStatus: to, Note: note, Sequence: len(events) + 1}
	return m.Store.SaveAuditEvent(event)
}
