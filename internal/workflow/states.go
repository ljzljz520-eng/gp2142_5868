package workflow

import "fmt"

const (
	StateDraft     = "draft"
	StateReview    = "review"
	StateConfirmed = "confirmed"
	StatePublished = "published"
	StateArchived  = "archived"
)

func IsKnownState(value string) bool {
	switch value {
	case StateDraft, StateReview, StateConfirmed, StatePublished, StateArchived:
		return true
	default:
		return false
	}
}

func CanMove(from, to string) bool {
	if !IsKnownState(from) || !IsKnownState(to) {
		return false
	}
	switch from {
	case StateDraft:
		return to == StateReview
	case StateReview:
		return to == StateConfirmed
	case StateConfirmed:
		return to == StatePublished
	case StatePublished:
		return to == StateArchived
	default:
		return false
	}
}

func ValidateState(state string) error {
	if !IsKnownState(state) {
		return fmt.Errorf("unknown workflow state %q", state)
	}
	return nil
}
