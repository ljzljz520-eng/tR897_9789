package app

import (
	"fmt"
	"sort"
	"strings"

	"callcentertroubleshooter/internal/domain"
	"callcentertroubleshooter/internal/report"
)

type OperatorView struct {
	RecordID string
	Employee string
	Actor    string
	Status   map[domain.DomainName]string
}

func BuildOperatorView(record domain.TroubleshootingRecord) OperatorView {
	return OperatorView{RecordID: record.ID, Employee: record.EmployeeID, Actor: record.RequestedBy, Status: RecordStatus(record)}
}

func (s *Service) OperatorViews(employeeID string) ([]OperatorView, error) {
	records, err := s.store.ListRecords(employeeID)
	if err != nil {
		return nil, err
	}
	views := make([]OperatorView, 0, len(records))
	for _, record := range records {
		views = append(views, BuildOperatorView(record))
	}
	return views, nil
}

func SummarizeOperators(views []OperatorView) string {
	lines := make([]string, 0, len(views))
	for _, view := range views {
		keys := make([]string, 0, len(view.Status))
		for key := range view.Status {
			keys = append(keys, string(key))
		}
		sort.Strings(keys)
		states := make([]string, 0, len(keys))
		for _, key := range keys {
			states = append(states, key+"="+view.Status[domain.DomainName(key)])
		}
		lines = append(lines, fmt.Sprintf("%s %s %s %s", view.RecordID, view.Employee, view.Actor, strings.Join(states, ",")))
	}
	return strings.Join(lines, "\n")
}

func (s *Service) ExportHistory(employeeID string) (string, error) {
	history, err := s.History(employeeID)
	if err != nil {
		return "", err
	}
	return report.EncodeExport(history)
}

func (s *Service) Healthy() bool {
	snapshot, err := s.HealthSummary()
	return err == nil && snapshot.StorageOK && snapshot.FixturesOK
}

func (s *Service) EnsureReady() error {
	if err := ValidateServiceState(s); err != nil {
		return err
	}
	if !s.Healthy() {
		return fmt.Errorf("service dependencies are degraded")
	}
	return nil
}
