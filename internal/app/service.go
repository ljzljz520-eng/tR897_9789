package app

import (
	"fmt"

	"callcentertroubleshooter/internal/domain"
	"callcentertroubleshooter/internal/fixtures"
	"callcentertroubleshooter/internal/report"
	"callcentertroubleshooter/internal/store"
	"callcentertroubleshooter/internal/troubleshoot"
)

type Service struct {
	queryer *troubleshoot.Queryer
	store   *store.SQLiteStore
	catalog *fixtures.Catalog
	seq     int
}

func NewService(queryer *troubleshoot.Queryer, storage *store.SQLiteStore, catalog *fixtures.Catalog) (*Service, error) {
	if queryer == nil || storage == nil || catalog == nil {
		return nil, fmt.Errorf("service dependencies are required")
	}
	return &Service{queryer: queryer, store: storage, catalog: catalog}, nil
}

func (s *Service) Diagnose(employeeID, actor string) (domain.TroubleshootingRecord, error) {
	if err := troubleshoot.ValidateRequest(employeeID, actor); err != nil {
		return domain.TroubleshootingRecord{}, err
	}
	result, err := s.queryer.QueryAccount(domain.NormalizeEmployeeID(employeeID))
	if err != nil {
		return domain.TroubleshootingRecord{}, err
	}
	s.seq++
	record := domain.TroubleshootingRecord{ID: fmt.Sprintf("TR-%06d", s.seq), EmployeeID: result.Account.EmployeeID, RequestedBy: troubleshoot.NormalizeActor(actor), Scope: domain.CanonicalScope(), Snapshot: result.Account, CreatedAt: fmt.Sprintf("2026-08-24T13:%02d:00Z", s.seq)}
	event := domain.AuditEvent{ID: fmt.Sprintf("AU-%06d", s.seq), RecordID: record.ID, Actor: record.RequestedBy, Action: "diagnose", Occurred: record.CreatedAt}
	if err := s.store.SaveRecordAndAudit(record, event); err != nil {
		return domain.TroubleshootingRecord{}, err
	}
	entry, err := s.catalog.Query(record.EmployeeID)
	if err != nil {
		return domain.TroubleshootingRecord{}, err
	}
	for _, name := range domain.AuthorizedDomains {
		credential, ok := entry.EncryptedPasswords[name]
		if !ok {
			return domain.TroubleshootingRecord{}, fmt.Errorf("encrypted credential missing for %s", name)
		}
		if err := s.store.SaveCredential(record.EmployeeID, credential); err != nil {
			return domain.TroubleshootingRecord{}, err
		}
	}
	return record, nil
}

func (s *Service) History(employeeID string) ([]report.HistoryItem, error) {
	records, err := s.store.ListRecords(employeeID)
	if err != nil {
		return nil, err
	}
	audits, err := s.store.ListAudits("")
	if err != nil {
		return nil, err
	}
	return report.BuildHistory(records, audits), nil
}

func (s *Service) HealthSummary() (report.HealthSnapshot, error) {
	storageOK := s.store.Ping() == nil
	fixturesOK := s.catalog.Health() == nil
	records, err := s.store.ListRecords("")
	if err != nil {
		return report.HealthSnapshot{}, err
	}
	recordCount, err := s.store.CountRecords()
	if err != nil {
		return report.HealthSnapshot{}, err
	}
	auditCount, err := s.store.CountAudits()
	if err != nil {
		return report.HealthSnapshot{}, err
	}
	return report.BuildHealth(storageOK, fixturesOK, len(s.catalog.SearchByDisplayName("")), recordCount, auditCount, records), nil
}

func (s *Service) ExportSummary() (string, error) {
	snapshot, err := s.HealthSummary()
	if err != nil {
		return "", err
	}
	return report.FormatHealth(snapshot), nil
}

func (s *Service) QueryPolicy() string {
	return s.queryer.PolicyDescription()
}
