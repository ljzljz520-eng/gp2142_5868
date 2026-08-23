package workflow

import (
	"errors"
	"strings"
)

func ValidateTransitionInput(target, actor, note string) error {
	if strings.TrimSpace(actor) == "" {
		return errors.New("workflow actor is required")
	}
	switch target {
	case StateReview:
		if strings.TrimSpace(note) == "" {
			return errors.New("review note is required")
		}
	case StateArchived:
		if strings.TrimSpace(note) == "" {
			return errors.New("archive note is required")
		}
	case StatePublished:
		if strings.TrimSpace(note) == "" {
			return errors.New("publish note is required")
		}
	case StateConfirmed:
		return nil
	}
	return nil
}

func RequiredFields(target string) []string {
	switch target {
	case StateReview:
		return []string{"actor", "note"}
	case StateConfirmed:
		return []string{"actor", "note"}
	case StatePublished:
		return []string{"actor", "note"}
	case StateArchived:
		return []string{"actor", "note"}
	default:
		return []string{}
	}
}

func IsTerminal(state string) bool {
	return state == StateArchived
}
