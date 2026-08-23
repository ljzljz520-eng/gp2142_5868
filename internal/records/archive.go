package records

import (
	"errors"
	"fmt"

	"example.com/othello/internal/store"
	"example.com/othello/internal/workflow"
)

func (r *Repository) ReadyForArchive(id string) (bool, string, error) {
	record, err := r.Get(id)
	if err != nil {
		return false, "", err
	}
	if record.Archived {
		return false, "already archived", nil
	}
	if record.Status != workflow.StatePublished {
		return false, fmt.Sprintf("status is %s", record.Status), nil
	}
	return true, "published record is ready", nil
}

func (r *Repository) RemoveDraft(id string) error {
	record, err := r.Get(id)
	if err != nil {
		return err
	}
	if record.Status != workflow.StateDraft {
		return errors.New("only draft records can be removed")
	}
	return r.Store.RemoveRecord(id)
}

func CloneRecord(record store.Record) store.Record {
	return record
}
