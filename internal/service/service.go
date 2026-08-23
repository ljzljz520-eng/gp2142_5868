package service

import (
	"errors"
	"fmt"
	"strings"

	"example.com/othello/internal/game"
	"example.com/othello/internal/importer"
	"example.com/othello/internal/records"
	"example.com/othello/internal/store"
	"example.com/othello/internal/workflow"
)

type Service struct {
	Store    *store.Store
	Records  *records.Repository
	Workflow *workflow.Manager
}

func New(database *store.Store) (*Service, error) {
	if database == nil {
		return nil, errors.New("service store is required")
	}
	repository, err := records.NewRepository(database)
	if err != nil {
		return nil, err
	}
	manager, err := workflow.NewManager(database)
	if err != nil {
		return nil, err
	}
	return &Service{Store: database, Records: repository, Workflow: manager}, nil
}

func (s *Service) StartGame(id string, mode game.Mode) (game.Game, error) {
	return game.New(id, mode)
}

func (s *Service) PlayTurn(state *game.Game, coordinate string) (game.Move, error) {
	if state == nil {
		return game.Move{}, errors.New("game is required")
	}
	return state.Play(coordinate)
}

func (s *Service) ComputerTurn(state *game.Game) (game.Move, error) {
	if state == nil {
		return game.Move{}, errors.New("game is required")
	}
	return state.PlayComputerTurn()
}

func (s *Service) RegisterGame(state game.Game, region, actor string) (store.Record, error) {
	if err := game.ValidateState(state); err != nil {
		return store.Record{}, err
	}
	if err := game.ValidateMoveHistory(state); err != nil {
		return store.Record{}, err
	}
	return s.Records.Register(state, region, actor)
}

func (s *Service) ReviewRecord(recordID, reviewer, note string) (store.Workflow, error) {
	if strings.TrimSpace(note) == "" {
		return store.Workflow{}, errors.New("review note is required")
	}
	return s.Workflow.Review(recordID, reviewer, note)
}

func (s *Service) ConfirmRecord(recordID, reviewer, note string) (store.Workflow, error) {
	return s.Workflow.Confirm(recordID, reviewer, note)
}

func (s *Service) PublishRecord(recordID, reviewer, note string) (store.Workflow, error) {
	return s.Workflow.Publish(recordID, reviewer, note)
}

func (s *Service) ArchiveRecord(recordID, actor, note string) (store.Workflow, error) {
	return s.Workflow.Archive(recordID, actor, note)
}

func (s *Service) UpdateRecord(recordID, actor string, change records.UpdateRequest) (store.Record, error) {
	return s.Records.Update(recordID, actor, change)
}

func (s *Service) Search(filter records.Filter) (records.History, error) {
	return s.Records.History(filter)
}

func (s *Service) Audit(recordID string) (workflow.AuditSummary, error) {
	return s.Workflow.Audit(recordID)
}

func (s *Service) Import(request importer.Request, actor string) (importer.Result, error) {
	if strings.TrimSpace(actor) == "" {
		return importer.Result{}, errors.New("import actor is required")
	}
	parsed, message, err := importer.Parse(request)
	if err != nil {
		return importer.Result{}, err
	}
	attachment, err := importer.SaveAttachment(s.Store, request, parsed, message)
	if err != nil {
		return importer.Result{}, err
	}
	if err := importer.ValidateAttachment(attachment); err != nil {
		return importer.Result{}, err
	}
	result := importer.Result{Game: parsed, Attachment: attachment, Moves: len(parsed.Moves), Report: importer.BuildReport(attachment, parsed, message).String()}
	if parsed.Status == game.StatusEnded {
		if _, err := s.RegisterGame(parsed, request.Region, actor); err != nil {
			return importer.Result{}, fmt.Errorf("save imported result: %w", err)
		}
	}
	return result, nil
}

func (s *Service) Status(recordID string) (string, error) {
	record, err := s.Records.Get(recordID)
	if err != nil {
		return "", err
	}
	return record.Status, nil
}
