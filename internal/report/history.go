package report

import (
	"fmt"
	"sort"
	"strings"

	"callcentertroubleshooter/internal/domain"
)

type HistoryItem struct {
	RecordID    string
	EmployeeID  string
	RequestedBy string
	Health      string
	Locked      []domain.DomainName
	Missing     []domain.DomainName
	AuditCount  int
	CreatedAt   string
}

func BuildHistory(records []domain.TroubleshootingRecord, audits []domain.AuditEvent) []HistoryItem {
	auditCounts := make(map[string]int)
	for _, event := range audits {
		auditCounts[event.RecordID]++
	}
	items := make([]HistoryItem, 0, len(records))
	for _, record := range records {
		items = append(items, HistoryItem{RecordID: record.ID, EmployeeID: record.EmployeeID, RequestedBy: record.RequestedBy, Health: domain.AccountHealth(record.Snapshot), Locked: record.Snapshot.LockedDomains(), Missing: record.Snapshot.MissingDomains(), AuditCount: auditCounts[record.ID], CreatedAt: record.CreatedAt})
	}
	sort.Slice(items, func(i, j int) bool { return items[i].CreatedAt < items[j].CreatedAt })
	return items
}

func FormatHistory(items []HistoryItem) string {
	if len(items) == 0 {
		return "no troubleshooting records"
	}
	lines := make([]string, 0, len(items)+1)
	lines = append(lines, "troubleshooting history")
	for _, item := range items {
		lines = append(lines, fmt.Sprintf("%s employee=%s actor=%s health=%s audits=%d", item.RecordID, item.EmployeeID, item.RequestedBy, item.Health, item.AuditCount))
	}
	return strings.Join(lines, "\n")
}

func FilterHistory(items []HistoryItem, employeeID string) []HistoryItem {
	filtered := make([]HistoryItem, 0, len(items))
	for _, item := range items {
		if employeeID == "" || item.EmployeeID == employeeID {
			filtered = append(filtered, item)
		}
	}
	return filtered
}

func HistoryHealth(items []HistoryItem) map[string]int {
	counts := map[string]int{"ready": 0, "attention": 0, "incomplete": 0}
	for _, item := range items {
		counts[item.Health]++
	}
	return counts
}
