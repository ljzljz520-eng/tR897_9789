package report

import (
	"encoding/json"
	"fmt"
	"sort"

	"callcentertroubleshooter/internal/domain"
)

type ExportRow struct {
	RecordID   string `json:"record_id"`
	EmployeeID string `json:"employee_id"`
	Health     string `json:"health"`
	AuditCount int    `json:"audit_count"`
}

func BuildExportRows(items []HistoryItem) []ExportRow {
	rows := make([]ExportRow, 0, len(items))
	for _, item := range items {
		rows = append(rows, ExportRow{RecordID: item.RecordID, EmployeeID: item.EmployeeID, Health: item.Health, AuditCount: item.AuditCount})
	}
	sort.Slice(rows, func(i, j int) bool { return rows[i].RecordID < rows[j].RecordID })
	return rows
}

func EncodeExport(items []HistoryItem) (string, error) {
	data, err := json.Marshal(BuildExportRows(items))
	if err != nil {
		return "", fmt.Errorf("encode export: %w", err)
	}
	return string(data), nil
}

func ExplainAccount(account domain.AgentAccount) string {
	counts := domain.CountStatuses(account)
	return fmt.Sprintf("%s existing=%d missing=%d locked=%d", account.EmployeeID, counts.Existing, counts.Missing, counts.Locked)
}

func CompareHealth(left, right []HistoryItem) bool {
	return len(BuildExportRows(left)) == len(BuildExportRows(right))
}
