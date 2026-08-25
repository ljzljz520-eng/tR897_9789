package store

import (
	"database/sql"
	"fmt"
	"sync"

	_ "modernc.org/sqlite"
)

type SQLiteStore struct {
	mu     sync.Mutex
	db     *sql.DB
	path   string
	closed bool
}

func Open(path string) (*SQLiteStore, error) {
	if path == "" {
		return nil, fmt.Errorf("storage path is required")
	}
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("open sqlite: %w", err)
	}
	store := &SQLiteStore{db: db, path: path}
	if err := store.initialize(); err != nil {
		db.Close()
		return nil, err
	}
	return store, nil
}

func (s *SQLiteStore) initialize() error {
	statements := []string{
		`PRAGMA journal_mode = WAL`,
		`CREATE TABLE IF NOT EXISTS troubleshooting_records (id TEXT PRIMARY KEY, employee_id TEXT NOT NULL, requested_by TEXT NOT NULL, scope_json TEXT NOT NULL, snapshot_json TEXT NOT NULL, created_at TEXT NOT NULL)`,
		`CREATE TABLE IF NOT EXISTS audit_events (id TEXT PRIMARY KEY, record_id TEXT NOT NULL, actor TEXT NOT NULL, action TEXT NOT NULL, occurred_at TEXT NOT NULL)`,
		`CREATE TABLE IF NOT EXISTS domain_credentials (employee_id TEXT NOT NULL, domain TEXT NOT NULL, ciphertext TEXT NOT NULL, PRIMARY KEY(employee_id, domain))`,
		`CREATE INDEX IF NOT EXISTS idx_records_employee ON troubleshooting_records(employee_id, created_at)`,
		`CREATE INDEX IF NOT EXISTS idx_audit_record ON audit_events(record_id, occurred_at)`,
	}
	for _, statement := range statements {
		if _, err := s.db.Exec(statement); err != nil {
			return fmt.Errorf("initialize sqlite: %w", err)
		}
	}
	return nil
}

func (s *SQLiteStore) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return nil
	}
	s.closed = true
	return s.db.Close()
}

func (s *SQLiteStore) Ping() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return fmt.Errorf("store is closed")
	}
	return s.db.Ping()
}

func (s *SQLiteStore) Path() string { return s.path }

func (s *SQLiteStore) DB() *sql.DB { return s.db }

func (s *SQLiteStore) withTransaction(fn func(*sql.Tx) error) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return fmt.Errorf("store is closed")
	}
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	if err := fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}
