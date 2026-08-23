package records

import (
	"errors"
	"fmt"
	"strings"

	"example.com/othello/internal/store"
	"example.com/othello/internal/workflow"
)

type ChangeSummary struct {
	RecordID string
	Version  int
	Actor    string
	Notes    string
}

func ValidateUpdate(actor string, request UpdateRequest) error {
	if strings.TrimSpace(actor) == "" {
		return errors.New("actor is required")
	}
	if strings.TrimSpace(request.Region) == "" && strings.TrimSpace(request.Notes) == "" {
		return errors.New("an update must change region or notes")
	}
	return nil
}

func (r *Repository) ApplyChange(id, actor string, request UpdateRequest) (ChangeSummary, error) {
	if err := ValidateUpdate(actor, request); err != nil {
		return ChangeSummary{}, err
	}
	record, err := r.Update(id, actor, request)
	if err != nil {
		return ChangeSummary{}, err
	}
	return ChangeSummary{RecordID: record.ID, Version: record.Version, Actor: actor, Notes: record.Notes}, nil
}

func (r *Repository) PendingReviews() ([]store.Record, error) {
	return r.Search(Filter{Status: workflow.StateReview})
}

func (r *Repository) RequireVersion(record store.Record, expected int) error {
	if expected <= 0 {
		return errors.New("expected version must be positive")
	}
	if record.Version != expected {
		return fmt.Errorf("record %s is version %d, expected %d", record.ID, record.Version, expected)
	}
	return nil
}
