package domain

import (
	"fmt"
	"strings"
)

type AccountSummary struct {
	EmployeeID string
	Health     string
	Existing   int
	Missing    int
	Locked     int
	Coverage   float64
}

func SummarizeAccount(account AgentAccount) AccountSummary {
	counts := CountStatuses(account)
	return AccountSummary{EmployeeID: account.EmployeeID, Health: AccountHealth(account), Existing: counts.Existing, Missing: counts.Missing, Locked: counts.Locked, Coverage: LoginCoverage(account)}
}

func (s AccountSummary) String() string {
	return fmt.Sprintf("employee=%s health=%s existing=%d missing=%d locked=%d coverage=%.2f", s.EmployeeID, s.Health, s.Existing, s.Missing, s.Locked, s.Coverage)
}

func AccountLabels(account AgentAccount) []string {
	labels := make([]string, 0, len(account.Statuses()))
	for _, status := range account.Statuses() {
		labels = append(labels, DisplayLabel(status.Domain)+":"+statusState(status))
	}
	return labels
}

func statusState(status DomainStatus) string {
	if !status.Exists {
		return "missing"
	}
	if status.Locked {
		return "locked"
	}
	return "active"
}

func JoinLabels(account AgentAccount) string {
	return strings.Join(AccountLabels(account), ", ")
}

func AnyLocked(accounts []AgentAccount) bool {
	for _, account := range accounts {
		if len(account.LockedDomains()) > 0 {
			return true
		}
	}
	return false
}

func CountByHealth(accounts []AgentAccount) map[string]int {
	counts := make(map[string]int)
	for _, account := range accounts {
		counts[AccountHealth(account)]++
	}
	return counts
}
