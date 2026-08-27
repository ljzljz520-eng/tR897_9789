package store

import (
	"database/sql"
	"fmt"

	"callcentertroubleshooter/internal/domain"
)

func (s *SQLiteStore) SaveAudit(event domain.AuditEvent) error {
	if err := event.Validate(); err != nil {
		return err
	}
	return s.withTransaction(func(tx *sql.Tx) error {
		_, err := tx.Exec(`INSERT INTO audit_events (id, record_id, actor, action, occurred_at) VALUES (?, ?, ?, ?, ?)`, event.ID, event.RecordID, event.Actor, event.Action, event.Occurred)
		if err != nil {
			return fmt.Errorf("save audit: %w", err)
		}
		return nil
	})
}

func (s *SQLiteStore) ListAudits(recordID string) ([]domain.AuditEvent, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return nil, fmt.Errorf("store is closed")
	}
	query := `SELECT id, record_id, actor, action, occurred_at FROM audit_events ORDER BY occurred_at, id`
	args := []any{}
	if recordID != "" {
		query = `SELECT id, record_id, actor, action, occurred_at FROM audit_events WHERE record_id = ? ORDER BY occurred_at, id`
		args = append(args, recordID)
	}
	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make([]domain.AuditEvent, 0)
	for rows.Next() {
		var event domain.AuditEvent
		if err := rows.Scan(&event.ID, &event.RecordID, &event.Actor, &event.Action, &event.Occurred); err != nil {
			return nil, err
		}
		result = append(result, event)
	}
	return result, rows.Err()
}

func (s *SQLiteStore) CountAudits() (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return 0, fmt.Errorf("store is closed")
	}
	var count int
	if err := s.db.QueryRow(`SELECT COUNT(*) FROM audit_events`).Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

func (s *SQLiteStore) SaveRecordAndAudit(record domain.TroubleshootingRecord, event domain.AuditEvent) error {
	if err := record.Validate(); err != nil {
		return err
	}
	if err := event.Validate(); err != nil {
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
		if _, err := tx.Exec(`INSERT INTO troubleshooting_records (id, employee_id, requested_by, scope_json, snapshot_json, created_at) VALUES (?, ?, ?, ?, ?, ?)`, record.ID, record.EmployeeID, record.RequestedBy, scope, snapshot, record.CreatedAt); err != nil {
			return err
		}
		if _, err := tx.Exec(`INSERT INTO audit_events (id, record_id, actor, action, occurred_at) VALUES (?, ?, ?, ?, ?)`, event.ID, event.RecordID, event.Actor, event.Action, event.Occurred); err != nil {
			return err
		}
		return nil
	})
}
