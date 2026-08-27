package report

import (
	"fmt"
	"sort"

	"callcentertroubleshooter/internal/domain"
)

type Trend struct {
	EmployeeID string
	Records    int
	Locked     int
	Missing    int
}

func BuildTrends(records []domain.TroubleshootingRecord) []Trend {
	byEmployee := make(map[string]*Trend)
	for _, record := range records {
		trend := byEmployee[record.EmployeeID]
		if trend == nil {
			trend = &Trend{EmployeeID: record.EmployeeID}
			byEmployee[record.EmployeeID] = trend
		}
		trend.Records++
		trend.Locked += len(record.Snapshot.LockedDomains())
		trend.Missing += len(record.Snapshot.MissingDomains())
	}
	result := make([]Trend, 0, len(byEmployee))
	for _, trend := range byEmployee {
		result = append(result, *trend)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].EmployeeID < result[j].EmployeeID })
	return result
}

func FormatTrends(trends []Trend) string {
	lines := make([]string, 0, len(trends))
	for _, trend := range trends {
		lines = append(lines, fmt.Sprintf("%s records=%d locked=%d missing=%d", trend.EmployeeID, trend.Records, trend.Locked, trend.Missing))
	}
	return joinTrendLines(lines)
}

func joinTrendLines(lines []string) string {
	result := ""
	for index, line := range lines {
		if index != 0 {
			result += "\n"
		}
		result += line
	}
	return result
}

func EscalationNeeded(trend Trend) bool {
	return trend.Locked > 0 || trend.Missing > 0
}

func EscalationList(trends []Trend) []string {
	result := make([]string, 0)
	for _, trend := range trends {
		if EscalationNeeded(trend) {
			result = append(result, trend.EmployeeID)
		}
	}
	return result
}
