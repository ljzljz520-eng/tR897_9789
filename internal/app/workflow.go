package app

import (
	"fmt"
	"strings"

	"callcentertroubleshooter/internal/domain"
	"callcentertroubleshooter/internal/report"
)

type WorkflowResult struct {
	Record  domain.TroubleshootingRecord
	History []report.HistoryItem
	Output  string
}

func (s *Service) RunDiagnosisWorkflow(employeeID, actor string) (WorkflowResult, error) {
	record, err := s.Diagnose(employeeID, actor)
	if err != nil {
		return WorkflowResult{}, err
	}
	return WorkflowResult{Record: record, Output: FormatRecord(record)}, nil
}

func (s *Service) RunHistoryWorkflow(employeeID string) (WorkflowResult, error) {
	history, err := s.History(employeeID)
	if err != nil {
		return WorkflowResult{}, err
	}
	return WorkflowResult{History: history, Output: report.FormatHistory(history)}, nil
}

func (s *Service) RunHealthWorkflow() (WorkflowResult, error) {
	output, err := s.ExportSummary()
	if err != nil {
		return WorkflowResult{}, err
	}
	return WorkflowResult{Output: output}, nil
}

func FormatRecord(record domain.TroubleshootingRecord) string {
	statuses := record.Snapshot.Statuses()
	parts := make([]string, 0, len(statuses))
	for _, status := range statuses {
		parts = append(parts, fmt.Sprintf("%s exists=%t locked=%t last_login=%s", status.Domain, status.Exists, status.Locked, status.LastLogin))
	}
	return strings.Join([]string{record.ID, "employee=" + record.EmployeeID, "actor=" + record.RequestedBy, strings.Join(parts, " ")}, " ")
}

func RecordStatus(record domain.TroubleshootingRecord) map[domain.DomainName]string {
	result := make(map[domain.DomainName]string)
	for _, status := range record.Snapshot.Statuses() {
		state := "missing"
		if status.Exists {
			state = "active"
		}
		if status.Locked {
			state = "locked"
		}
		result[status.Domain] = state
	}
	return result
}
