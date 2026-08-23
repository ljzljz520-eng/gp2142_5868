package store

import (
	"fmt"
	"sort"
)

type Snapshot struct {
	Records     int
	Workflows   int
	AuditEvents int
	Attachments int
}

func (s *Store) Snapshot() (Snapshot, error) {
	records, err := s.ListRecords()
	if err != nil {
		return Snapshot{}, err
	}
	workflows, err := s.listWorkflows()
	if err != nil {
		return Snapshot{}, err
	}
	audits, err := s.ListAuditEvents("")
	if err != nil {
		return Snapshot{}, err
	}
	attachments, err := s.ListAttachments("")
	if err != nil {
		return Snapshot{}, err
	}
	return Snapshot{Records: len(records), Workflows: len(workflows), AuditEvents: len(audits), Attachments: len(attachments)}, nil
}

func (s *Store) listWorkflows() ([]Workflow, error) {
	values := make([]Workflow, 0)
	err := s.list("workflows", func() any { return &Workflow{} }, func(value any) error {
		values = append(values, *(value.(*Workflow)))
		return nil
	})
	sort.Slice(values, func(i, j int) bool { return values[i].ID < values[j].ID })
	return values, err
}

func (s *Store) ValidateReferences() error {
	records, err := s.ListRecords()
	if err != nil {
		return err
	}
	for _, record := range records {
		workflow, err := s.GetWorkflow("workflow-" + record.ID)
		if err != nil {
			return fmt.Errorf("record %s has no workflow: %w", record.ID, err)
		}
		if workflow.RecordID != record.ID {
			return fmt.Errorf("workflow %s points to %s", workflow.ID, workflow.RecordID)
		}
	}
	return nil
}
