package records

import (
	"errors"
	"fmt"
	"strings"

	"example.com/othello/internal/game"
	"example.com/othello/internal/store"
	"example.com/othello/internal/workflow"
)

type Repository struct {
	Store *store.Store
}

type UpdateRequest struct {
	Region string
	Notes  string
}

type Filter struct {
	Region   string
	Status   string
	Text     string
	Archived *bool
}

func NewRepository(database *store.Store) (*Repository, error) {
	if database == nil {
		return nil, errors.New("store is required")
	}
	return &Repository{Store: database}, nil
}

func (r *Repository) Register(result game.Game, region, createdBy string) (store.Record, error) {
	if r == nil || r.Store == nil {
		return store.Record{}, errors.New("record repository is not ready")
	}
	if result.Status != game.StatusEnded {
		return store.Record{}, errors.New("only ended games can be registered")
	}
	if strings.TrimSpace(region) == "" {
		return store.Record{}, errors.New("region is required")
	}
	if strings.TrimSpace(createdBy) == "" {
		return store.Record{}, errors.New("creator is required")
	}
	record := store.Record{
		ID:         "record-" + result.ID,
		GameID:     result.ID,
		Winner:     string(result.Winner),
		BlackScore: result.BlackScore,
		WhiteScore: result.WhiteScore,
		Status:     workflow.StateDraft,
		Region:     strings.TrimSpace(region),
		Summary:    result.Summary(),
		CreatedBy:  strings.TrimSpace(createdBy),
		Version:    1,
	}
	if err := r.Store.SaveRecord(record); err != nil {
		return store.Record{}, err
	}
	initial := store.Workflow{ID: "workflow-" + record.ID, RecordID: record.ID, State: workflow.StateDraft, Version: 1}
	if err := r.Store.SaveWorkflow(initial); err != nil {
		return store.Record{}, err
	}
	if err := r.appendAudit(record.ID, createdBy, "register", "", workflow.StateDraft, "game result registered"); err != nil {
		return store.Record{}, err
	}
	return record, nil
}

func (r *Repository) Get(id string) (store.Record, error) {
	if r == nil || r.Store == nil {
		return store.Record{}, errors.New("record repository is not ready")
	}
	return r.Store.GetRecord(id)
}

func (r *Repository) Search(filter Filter) ([]store.Record, error) {
	if r == nil || r.Store == nil {
		return nil, errors.New("record repository is not ready")
	}
	all, err := r.Store.ListRecords()
	if err != nil {
		return nil, err
	}
	result := make([]store.Record, 0, len(all))
	for _, record := range all {
		if filter.Region != "" && !strings.EqualFold(record.Region, strings.TrimSpace(filter.Region)) {
			continue
		}
		if filter.Status != "" && record.Status != filter.Status {
			continue
		}
		if filter.Archived != nil && record.Archived != *filter.Archived {
			continue
		}
		if filter.Text != "" && !containsRecordText(record, filter.Text) {
			continue
		}
		result = append(result, record)
	}
	return result, nil
}

func containsRecordText(record store.Record, text string) bool {
	needle := strings.ToLower(strings.TrimSpace(text))
	return strings.Contains(strings.ToLower(record.ID), needle) || strings.Contains(strings.ToLower(record.Summary), needle) || strings.Contains(strings.ToLower(record.Notes), needle)
}

func (r *Repository) Update(id, actor string, request UpdateRequest) (store.Record, error) {
	if err := ValidateUpdate(actor, request); err != nil {
		return store.Record{}, err
	}
	record, err := r.Get(id)
	if err != nil {
		return store.Record{}, err
	}
	if record.Archived {
		return store.Record{}, errors.New("archived record cannot be changed")
	}
	if strings.TrimSpace(request.Region) != "" {
		record.Region = strings.TrimSpace(request.Region)
	}
	if strings.TrimSpace(request.Notes) != "" {
		record.Notes = appendNote(record.Notes, actor, request.Notes)
	}
	record.Version++
	if err := r.Store.SaveRecord(record); err != nil {
		return store.Record{}, err
	}
	if err := r.appendAudit(record.ID, actor, "update", record.Status, record.Status, request.Notes); err != nil {
		return store.Record{}, err
	}
	return record, nil
}

func appendNote(existing, actor, note string) string {
	entry := fmt.Sprintf("%s: %s", strings.TrimSpace(actor), strings.TrimSpace(note))
	if strings.TrimSpace(existing) == "" {
		return entry
	}
	return existing + " | " + entry
}

func (r *Repository) Publish(id, actor string) (store.Record, error) {
	record, err := r.Get(id)
	if err != nil {
		return store.Record{}, err
	}
	if record.Status != workflow.StateConfirmed {
		return store.Record{}, fmt.Errorf("record %s must be confirmed before publishing", id)
	}
	record.Status = workflow.StatePublished
	record.Published = true
	record.Version++
	if err := r.Store.SaveRecord(record); err != nil {
		return store.Record{}, err
	}
	if err := r.appendAudit(record.ID, actor, "publish", workflow.StateConfirmed, workflow.StatePublished, "published to history"); err != nil {
		return store.Record{}, err
	}
	return record, nil
}

func (r *Repository) Archive(id, actor string) (store.Record, error) {
	record, err := r.Get(id)
	if err != nil {
		return store.Record{}, err
	}
	if record.Status != workflow.StatePublished {
		return store.Record{}, fmt.Errorf("record %s must be published before archiving", id)
	}
	record.Status = workflow.StateArchived
	record.Archived = true
	record.Version++
	if err := r.Store.SaveRecord(record); err != nil {
		return store.Record{}, err
	}
	if err := r.appendAudit(record.ID, actor, "archive", workflow.StatePublished, workflow.StateArchived, "archived result"); err != nil {
		return store.Record{}, err
	}
	return record, nil
}

func (r *Repository) appendAudit(recordID, actor, action, from, to, note string) error {
	events, err := r.Store.ListAuditEvents(recordID)
	if err != nil {
		return err
	}
	event := store.AuditEvent{ID: fmt.Sprintf("audit-%s-%d", recordID, len(events)+1), RecordID: recordID, Actor: actor, Action: action, FromStatus: from, ToStatus: to, Note: note, Sequence: len(events) + 1}
	return r.Store.SaveAuditEvent(event)
}
