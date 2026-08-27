package store

import (
	"database/sql"
	"fmt"

	"callcentertroubleshooter/internal/domain"
)

func (s *SQLiteStore) SaveRecord(record domain.TroubleshootingRecord) error {
	if err := record.Validate(); err != nil {
		return err
	}
	scope, err := domain.EncodeScope(record.Scope)
	if err != nil {
		return err
	}
	snapshot, err := domain.EncodeAccount(record.Snapshot)
	if err != nil {
		return err
	}
	return s.withTransaction(func(tx *sql.Tx) error {
		_, err := tx.Exec(`INSERT INTO troubleshooting_records (id, employee_id, requested_by, scope_json, snapshot_json, created_at) VALUES (?, ?, ?, ?, ?, ?)`, record.ID, record.EmployeeID, record.RequestedBy, scope, snapshot, record.CreatedAt)
		if err != nil {
			return fmt.Errorf("save record: %w", err)
		}
		return nil
	})
}

func (s *SQLiteStore) GetRecord(id string) (domain.TroubleshootingRecord, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return domain.TroubleshootingRecord{}, fmt.Errorf("store is closed")
	}
	row := s.db.QueryRow(`SELECT id, employee_id, requested_by, scope_json, snapshot_json, created_at FROM troubleshooting_records WHERE id = ?`, id)
	return scanRecord(row)
}

func (s *SQLiteStore) ListRecords(employeeID string) ([]domain.TroubleshootingRecord, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return nil, fmt.Errorf("store is closed")
	}
	query := `SELECT id, employee_id, requested_by, scope_json, snapshot_json, created_at FROM troubleshooting_records ORDER BY created_at, id`
	args := []any{}
	if employeeID != "" {
		query = `SELECT id, employee_id, requested_by, scope_json, snapshot_json, created_at FROM troubleshooting_records WHERE employee_id = ? ORDER BY created_at, id`
		args = append(args, employeeID)
	}
	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("list records: %w", err)
	}
	defer rows.Close()
	results := make([]domain.TroubleshootingRecord, 0)
	for rows.Next() {
		record, err := scanRecord(rows)
		if err != nil {
			return nil, err
		}
		results = append(results, record)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return results, nil
}

type scanner interface{ Scan(...any) error }

func scanRecord(row scanner) (domain.TroubleshootingRecord, error) {
	var record domain.TroubleshootingRecord
	var scopeJSON, snapshotJSON string
	if err := row.Scan(&record.ID, &record.EmployeeID, &record.RequestedBy, &scopeJSON, &snapshotJSON, &record.CreatedAt); err != nil {
		return domain.TroubleshootingRecord{}, err
	}
	scope, err := domain.DecodeScope(scopeJSON)
	if err != nil {
		return domain.TroubleshootingRecord{}, err
	}
	snapshot, err := domain.DecodeAccount(snapshotJSON)
	if err != nil {
		return domain.TroubleshootingRecord{}, err
	}
	record.Scope, record.Snapshot = scope, snapshot
	return record, nil
}

func (s *SQLiteStore) CountRecords() (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return 0, fmt.Errorf("store is closed")
	}
	var count int
	if err := s.db.QueryRow(`SELECT COUNT(*) FROM troubleshooting_records`).Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}
