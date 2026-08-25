package store

import (
	"fmt"
	"strings"

	"callcentertroubleshooter/internal/domain"
)

func (s *SQLiteStore) ExportIDs() (string, error) {
	records, err := s.ListRecords("")
	if err != nil {
		return "", err
	}
	ids := make([]string, 0, len(records))
	for _, record := range records {
		ids = append(ids, record.ID)
	}
	return strings.Join(ids, "\n"), nil
}

func (s *SQLiteStore) ImportMarker(recordID, actor string) error {
	if strings.TrimSpace(recordID) == "" || strings.TrimSpace(actor) == "" {
		return fmt.Errorf("record and actor are required")
	}
	return s.SaveAudit(domain.AuditEvent{ID: "IMPORT-" + recordID, RecordID: recordID, Actor: actor, Action: "import-marker", Occurred: "2026-08-24T13:59:00Z"})
}

func (s *SQLiteStore) VerifyRecord(recordID string) (bool, error) {
	record, err := s.GetRecord(recordID)
	if err != nil {
		return false, err
	}
	return record.EmployeeID != "" && len(record.Scope) == len(domain.AuthorizedDomains), nil
}

func (s *SQLiteStore) VerifyAll() (int, error) {
	records, err := s.ListRecords("")
	if err != nil {
		return 0, err
	}
	valid := 0
	for _, record := range records {
		if record.Validate() == nil {
			valid++
		}
	}
	return valid, nil
}

func (s *SQLiteStore) AuditTrail(recordID string) (string, error) {
	events, err := s.ListAudits(recordID)
	if err != nil {
		return "", err
	}
	lines := make([]string, 0, len(events))
	for _, event := range events {
		lines = append(lines, fmt.Sprintf("%s %s %s", event.Occurred, event.Actor, event.Action))
	}
	return strings.Join(lines, "\n"), nil
}
