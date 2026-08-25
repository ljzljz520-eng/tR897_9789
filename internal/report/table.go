package report

import (
	"fmt"
	"strings"

	"callcentertroubleshooter/internal/domain"
)

type Table struct {
	Headers []string
	Rows    [][]string
}

func HistoryTable(items []HistoryItem) Table {
	table := Table{Headers: []string{"record", "employee", "health", "audits"}, Rows: make([][]string, 0, len(items))}
	for _, item := range items {
		table.Rows = append(table.Rows, []string{item.RecordID, item.EmployeeID, item.Health, fmt.Sprintf("%d", item.AuditCount)})
	}
	return table
}

func (t Table) Render() string {
	lines := make([]string, 0, len(t.Rows)+1)
	lines = append(lines, strings.Join(t.Headers, "\t"))
	for _, row := range t.Rows {
		lines = append(lines, strings.Join(row, "\t"))
	}
	return strings.Join(lines, "\n")
}

func StatusTable(account domain.AgentAccount) Table {
	table := Table{Headers: []string{"domain", "exists", "locked", "last_login"}, Rows: make([][]string, 0, 3)}
	for _, status := range account.Statuses() {
		table.Rows = append(table.Rows, []string{string(status.Domain), fmt.Sprintf("%t", status.Exists), fmt.Sprintf("%t", status.Locked), status.LastLogin})
	}
	return table
}

func SummaryLines(snapshot HealthSnapshot) []string {
	return []string{
		fmt.Sprintf("storage: %t", snapshot.StorageOK),
		fmt.Sprintf("fixtures: %t", snapshot.FixturesOK),
		fmt.Sprintf("fixture accounts: %d", snapshot.FixtureAccounts),
		fmt.Sprintf("records: %d", snapshot.RecordCount),
		fmt.Sprintf("audits: %d", snapshot.AuditCount),
	}
}

func JoinSummary(snapshot HealthSnapshot) string {
	return strings.Join(SummaryLines(snapshot), "; ")
}
