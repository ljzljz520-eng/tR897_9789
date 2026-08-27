package store

import (
	"fmt"
	"strings"
)

type RecordFilter struct {
	EmployeeID string
	Actor      string
	Limit      int
}

func (f RecordFilter) normalized() RecordFilter {
	copy := f
	copy.EmployeeID = strings.TrimSpace(copy.EmployeeID)
	copy.Actor = strings.TrimSpace(copy.Actor)
	if copy.Limit <= 0 || copy.Limit > 100 {
		copy.Limit = 100
	}
	return copy
}

func (s *SQLiteStore) Search(filter RecordFilter) ([]string, error) {
	filter = filter.normalized()
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return nil, fmt.Errorf("store is closed")
	}
	query := `SELECT id FROM troubleshooting_records WHERE 1=1`
	args := make([]any, 0, 3)
	if filter.EmployeeID != "" {
		query += ` AND employee_id = ?`
		args = append(args, filter.EmployeeID)
	}
	if filter.Actor != "" {
		query += ` AND requested_by = ?`
		args = append(args, filter.Actor)
	}
	query += ` ORDER BY created_at, id LIMIT ?`
	args = append(args, filter.Limit)
	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	ids := make([]string, 0)
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

func (s *SQLiteStore) LatestRecord(employeeID string) (string, error) {
	ids, err := s.Search(RecordFilter{EmployeeID: employeeID, Limit: 1})
	if err != nil {
		return "", err
	}
	if len(ids) == 0 {
		return "", fmt.Errorf("no records for %s", employeeID)
	}
	return ids[len(ids)-1], nil
}

func (s *SQLiteStore) RecordExists(id string) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return false, fmt.Errorf("store is closed")
	}
	var count int
	if err := s.db.QueryRow(`SELECT COUNT(*) FROM troubleshooting_records WHERE id = ?`, id).Scan(&count); err != nil {
		return false, err
	}
	return count > 0, nil
}

func (s *SQLiteStore) AuditForActor(actor string) ([]string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return nil, fmt.Errorf("store is closed")
	}
	rows, err := s.db.Query(`SELECT id FROM audit_events WHERE actor = ? ORDER BY occurred_at, id`, actor)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	ids := make([]string, 0)
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}
