package workflow

import (
	"errors"
	"strings"

	"example.com/othello/internal/store"
)

type ReviewRequest struct {
	Reviewer string
	Note     string
}

func (m *Manager) ValidateReview(request ReviewRequest) error {
	if strings.TrimSpace(request.Reviewer) == "" {
		return errors.New("reviewer is required")
	}
	if strings.TrimSpace(request.Note) == "" {
		return errors.New("review note is required")
	}
	return nil
}

func (m *Manager) ReviewDetails(recordID string) (store.Workflow, string, error) {
	workflow, err := m.Get(recordID)
	if err != nil {
		return store.Workflow{}, "", err
	}
	if err := ValidateState(workflow.State); err != nil {
		return store.Workflow{}, "", err
	}
	fields := strings.Join(RequiredFields(workflow.State), ",")
	if fields != "" {
		return workflow, reviewDescription(workflow.State) + " requires " + fields, nil
	}
	return workflow, reviewDescription(workflow.State), nil
}

func reviewDescription(state string) string {
	switch state {
	case StateDraft:
		return "awaiting reviewer"
	case StateReview:
		return "review in progress"
	case StateConfirmed:
		return "review confirmed"
	case StatePublished:
		return "visible in history"
	case StateArchived:
		return "read-only archive"
	default:
		return "unknown state"
	}
}
