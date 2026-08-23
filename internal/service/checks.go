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

type CheckResult struct {
	RecordID string
	Checks   []string
	Passed   bool
}

func (s *Service) ValidateRecord(id string) (CheckResult, error) {
	record, err := s.Records.Get(id)
	if err != nil {
		return CheckResult{}, err
	}
	checks := []string{}
	if record.GameID != "" {
		checks = append(checks, "game-linked")
	}
	if record.BlackScore >= 0 && record.WhiteScore >= 0 {
		checks = append(checks, "score-valid")
	}
	if workflow.IsKnownState(record.Status) {
		checks = append(checks, "status-valid")
	}
	return CheckResult{RecordID: record.ID, Checks: checks, Passed: len(checks) == 3}, nil
}

func (s *Service) PrepareGame(id string, mode game.Mode, opening []string) (game.Game, error) {
	if strings.TrimSpace(id) == "" {
		return game.Game{}, errors.New("game id is required")
	}
	state, err := game.ParseMoves(id, mode, opening)
	if err != nil {
		return game.Game{}, err
	}
	if state.Status == game.StatusEnded && len(state.Moves) == 0 {
		return game.Game{}, errors.New("an ended game needs at least one move")
	}
	return state, nil
}

func (s *Service) UpdateAndCheck(id, actor string, request records.UpdateRequest) (store.Record, CheckResult, error) {
	change, err := s.Records.ApplyChange(id, actor, request)
	if err != nil {
		return store.Record{}, CheckResult{}, err
	}
	record, err := s.Records.Get(change.RecordID)
	if err != nil {
		return store.Record{}, CheckResult{}, err
	}
	check, err := s.ValidateRecord(record.ID)
	if err != nil {
		return store.Record{}, CheckResult{}, err
	}
	if !check.Passed {
		return store.Record{}, CheckResult{}, fmt.Errorf("record %s failed validation", id)
	}
	return record, check, nil
}
