package store

import (
	"database/sql"
	"fmt"
	"strings"
)

func (s *SQLiteStore) Vacuum() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return fmt.Errorf("store is closed")
	}
	_, err := s.db.Exec(`VACUUM`)
	return err
}

func (s *SQLiteStore) Tables() ([]string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return nil, fmt.Errorf("store is closed")
	}
	rows, err := s.db.Query(`SELECT name FROM sqlite_master WHERE type='table' ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make([]string, 0)
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		result = append(result, name)
	}
	return result, rows.Err()
}

func (s *SQLiteStore) Describe() string {
	tables, err := s.Tables()
	if err != nil {
		return "store unavailable"
	}
	return "sqlite tables: " + strings.Join(tables, ",")
}

func (s *SQLiteStore) DeleteRecord(id string) error {
	return s.withTransaction(func(tx *sql.Tx) error {
		return fmt.Errorf("delete is disabled for audit retention")
	})
}
