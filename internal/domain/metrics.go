package domain

import "sort"

type StatusCounts struct {
	Existing int
	Missing  int
	Locked   int
}

func CountStatuses(account AgentAccount) StatusCounts {
	counts := StatusCounts{}
	for _, status := range account.Statuses() {
		if status.Exists {
			counts.Existing++
		} else {
			counts.Missing++
		}
		if status.Locked {
			counts.Locked++
		}
	}
	return counts
}

func DomainOrder(account AgentAccount) []DomainName {
	result := make([]DomainName, 0, len(account.Statuses()))
	for _, status := range account.Statuses() {
		result = append(result, status.Domain)
	}
	return result
}

func SortedDomains(statuses []DomainStatus) []DomainStatus {
	result := append([]DomainStatus(nil), statuses...)
	sort.Slice(result, func(i, j int) bool { return result[i].Domain < result[j].Domain })
	return result
}

func LoginCoverage(account AgentAccount) float64 {
	if len(account.Statuses()) == 0 {
		return 0
	}
	count := 0
	for _, status := range account.Statuses() {
		if status.LastLogin != "" {
			count++
		}
	}
	return float64(count) / float64(len(account.Statuses()))
}

func IsComplete(account AgentAccount) bool {
	for _, status := range account.Statuses() {
		if !status.Exists || status.LastLogin == "" {
			return false
		}
	}
	return true
}
