package fixtures

import (
	"fmt"
	"sort"

	"callcentertroubleshooter/internal/domain"
)

type MatrixRow struct {
	EmployeeID string
	Domain     domain.DomainName
	Exists     bool
	Locked     bool
}

func (c *Catalog) Matrix() []MatrixRow {
	rows := make([]MatrixRow, 0)
	for _, employeeID := range c.repository.IDs() {
		entry, err := c.Query(employeeID)
		if err != nil {
			continue
		}
		for _, status := range []domain.DomainStatus{entry.Office, entry.Telephony, entry.Quality} {
			rows = append(rows, MatrixRow{EmployeeID: employeeID, Domain: status.Domain, Exists: status.Exists, Locked: status.Locked})
		}
	}
	sort.Slice(rows, func(i, j int) bool {
		if rows[i].EmployeeID == rows[j].EmployeeID {
			return rows[i].Domain < rows[j].Domain
		}
		return rows[i].EmployeeID < rows[j].EmployeeID
	})
	return rows
}

func FormatMatrix(rows []MatrixRow) string {
	lines := make([]string, 0, len(rows))
	for _, row := range rows {
		lines = append(lines, fmt.Sprintf("%s %s exists=%t locked=%t", row.EmployeeID, row.Domain, row.Exists, row.Locked))
	}
	return fmt.Sprintf("%d rows\n%s", len(lines), joinLines(lines))
}

func joinLines(lines []string) string {
	result := ""
	for index, line := range lines {
		if index > 0 {
			result += "\n"
		}
		result += line
	}
	return result
}

func MatrixCounts(rows []MatrixRow) map[string]int {
	counts := map[string]int{"exists": 0, "missing": 0, "locked": 0}
	for _, row := range rows {
		if row.Exists {
			counts["exists"]++
		} else {
			counts["missing"]++
		}
		if row.Locked {
			counts["locked"]++
		}
	}
	return counts
}
