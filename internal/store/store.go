package store

import (
	"errors"
	"fmt"
	"path/filepath"
	"sort"

	bolt "go.etcd.io/bbolt"
)

type Record struct {
	ID         string `json:"id"`
	GameID     string `json:"game_id"`
	Winner     string `json:"winner"`
	BlackScore int    `json:"black_score"`
	WhiteScore int    `json:"white_score"`
	Status     string `json:"status"`
	Region     string `json:"region"`
	Summary    string `json:"summary"`
	Notes      string `json:"notes"`
	CreatedBy  string `json:"created_by"`
	Archived   bool   `json:"archived"`
	Published  bool   `json:"published"`
	Version    int    `json:"version"`
}

type AuditEvent struct {
	ID         string `json:"id"`
	RecordID   string `json:"record_id"`
	Actor      string `json:"actor"`
	Action     string `json:"action"`
	FromStatus string `json:"from_status"`
	ToStatus   string `json:"to_status"`
	Note       string `json:"note"`
	Sequence   int    `json:"sequence"`
}

type Workflow struct {
	ID         string `json:"id"`
	RecordID   string `json:"record_id"`
	State      string `json:"state"`
	Reviewer   string `json:"reviewer"`
	ReviewNote string `json:"review_note"`
	Version    int    `json:"version"`
}

type Attachment struct {
	ID       string `json:"id"`
	RecordID string `json:"record_id"`
	Name     string `json:"name"`
	Kind     string `json:"kind"`
	Content  string `json:"content"`
}

type Store struct {
	db *bolt.DB
}

func Open(path string) (*Store, error) {
	if path == "" {
		return nil, errors.New("store path is required")
	}
	database, err := bolt.Open(filepath.Clean(path), 0o600, nil)
	if err != nil {
		return nil, fmt.Errorf("open bbolt store: %w", err)
	}
	store := &Store{db: database}
	if err := store.initialize(); err != nil {
		database.Close()
		return nil, err
	}
	return store, nil
}

func (s *Store) initialize() error {
	return s.db.Update(func(tx *bolt.Tx) error {
		for _, name := range bucketNames {
			if _, err := tx.CreateBucketIfNotExists(name); err != nil {
				return fmt.Errorf("create bucket: %w", err)
			}
		}
		return nil
	})
}

func (s *Store) Close() error {
	if s == nil || s.db == nil {
		return nil
	}
	return s.db.Close()
}

func (s *Store) save(bucketName, key string, value any) error {
	if key == "" {
		return errors.New("entity key is required")
	}
	data, err := encode(value)
	if err != nil {
		return err
	}
	return s.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(bucket(bucketName))
		if bucket == nil {
			return fmt.Errorf("bucket %s is missing", bucketName)
		}
		return bucket.Put([]byte(key), data)
	})
}

func (s *Store) load(bucketName, key string, target any) error {
	if key == "" {
		return errors.New("entity key is required")
	}
	return s.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(bucket(bucketName))
		if bucket == nil {
			return fmt.Errorf("bucket %s is missing", bucketName)
		}
		data := bucket.Get([]byte(key))
		if data == nil {
			return fmt.Errorf("%s %q not found", bucketName, key)
		}
		return decode(cloneBytes(data), target)
	})
}

func (s *Store) list(bucketName string, factory func() any, read func(any) error) error {
	return s.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(bucket(bucketName))
		if bucket == nil {
			return fmt.Errorf("bucket %s is missing", bucketName)
		}
		return bucket.ForEach(func(key, data []byte) error {
			value := factory()
			if err := decode(cloneBytes(data), value); err != nil {
				return err
			}
			return read(value)
		})
	})
}

func (s *Store) SaveRecord(value Record) error {
	return s.save("records", value.ID, value)
}

func (s *Store) GetRecord(id string) (Record, error) {
	var value Record
	err := s.load("records", id, &value)
	return value, err
}

func (s *Store) ListRecords() ([]Record, error) {
	values := make([]Record, 0)
	err := s.list("records", func() any { return &Record{} }, func(value any) error {
		values = append(values, *(value.(*Record)))
		return nil
	})
	sort.Slice(values, func(i, j int) bool { return values[i].ID < values[j].ID })
	return values, err
}

func (s *Store) SaveWorkflow(value Workflow) error {
	return s.save("workflows", value.ID, value)
}

func (s *Store) GetWorkflow(id string) (Workflow, error) {
	var value Workflow
	err := s.load("workflows", id, &value)
	return value, err
}

func (s *Store) SaveAuditEvent(value AuditEvent) error {
	return s.save("audits", value.ID, value)
}

func (s *Store) ListAuditEvents(recordID string) ([]AuditEvent, error) {
	values := make([]AuditEvent, 0)
	err := s.list("audits", func() any { return &AuditEvent{} }, func(value any) error {
		event := *(value.(*AuditEvent))
		if recordID == "" || event.RecordID == recordID {
			values = append(values, event)
		}
		return nil
	})
	sort.Slice(values, func(i, j int) bool { return values[i].Sequence < values[j].Sequence })
	return values, err
}

func (s *Store) SaveAttachment(value Attachment) error {
	return s.save("attachments", value.ID, value)
}

func (s *Store) GetAttachment(id string) (Attachment, error) {
	var value Attachment
	err := s.load("attachments", id, &value)
	return value, err
}

func (s *Store) ListAttachments(recordID string) ([]Attachment, error) {
	values := make([]Attachment, 0)
	err := s.list("attachments", func() any { return &Attachment{} }, func(value any) error {
		attachment := *(value.(*Attachment))
		if recordID == "" || attachment.RecordID == recordID {
			values = append(values, attachment)
		}
		return nil
	})
	sort.Slice(values, func(i, j int) bool { return values[i].ID < values[j].ID })
	return values, err
}

func (s *Store) RemoveRecord(id string) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(bucket("records"))
		if bucket == nil {
			return errors.New("records bucket is missing")
		}
		return bucket.Delete([]byte(id))
	})
}
