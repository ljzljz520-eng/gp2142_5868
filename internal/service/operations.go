package service

import (
	"errors"
	"fmt"
	"strings"

	"example.com/othello/internal/game"
	"example.com/othello/internal/records"
	"example.com/othello/internal/store"
	"example.com/othello/internal/workflow"
)

type LifecycleResult struct {
	Record  store.Record
	History records.History
	Audit   workflow.AuditSummary
}

func (s *Service) CompleteLifecycle(state game.Game, region, owner, reviewer string) (LifecycleResult, error) {
	if strings.TrimSpace(reviewer) == "" {
		return LifecycleResult{}, errors.New("reviewer is required")
	}
	record, err := s.RegisterGame(state, region, owner)
	if err != nil {
		return LifecycleResult{}, err
	}
	if _, err := s.ReviewRecord(record.ID, reviewer, "opening and score checked"); err != nil {
		return LifecycleResult{}, err
	}
	if _, err := s.ConfirmRecord(record.ID, reviewer, "result accepted"); err != nil {
		return LifecycleResult{}, err
	}
	if _, err := s.PublishRecord(record.ID, reviewer, "history published"); err != nil {
		return LifecycleResult{}, err
	}
	history, err := s.Search(records.Filter{Region: region, Status: workflow.StatePublished})
	if err != nil {
		return LifecycleResult{}, err
	}
	audit, err := s.Audit(record.ID)
	if err != nil {
		return LifecycleResult{}, err
	}
	record, err = s.Records.Get(record.ID)
	if err != nil {
		return LifecycleResult{}, err
	}
	return LifecycleResult{Record: record, History: history, Audit: audit}, nil
}

func (s *Service) SelectFirst(region, status string) (store.Record, error) {
	return s.Records.FindFirst(records.Filter{Region: strings.TrimSpace(region), Status: status})
}

func (s *Service) EnsureRecordReady(id string) error {
	record, err := s.Records.Get(id)
	if err != nil {
		return err
	}
	if record.ID == "" || record.GameID == "" {
		return fmt.Errorf("record %s is incomplete", id)
	}
	return s.Store.ValidateReferences()
}

func (s *Service) StorageSnapshot() (store.Snapshot, error) {
	return s.Store.Snapshot()
}

func (s *Service) OutcomeStats(filter records.Filter) (records.OutcomeStats, error) {
	return s.Records.Stats(filter)
}
