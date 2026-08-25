package troubleshoot

import (
	"fmt"
	"sort"

	"callcentertroubleshooter/internal/domain"
)

type BatchSummary struct {
	Requested int
	Returned  int
	Missing   []string
	Locked    []string
}

func SummarizeBatch(results []QueryResult) BatchSummary {
	summary := BatchSummary{Requested: len(results), Returned: len(results)}
	missing := make(map[string]bool)
	locked := make(map[string]bool)
	for _, result := range results {
		for _, name := range result.Account.MissingDomains() {
			missing[string(name)] = true
		}
		for _, name := range result.Account.LockedDomains() {
			locked[string(name)] = true
		}
	}
	for name := range missing {
		summary.Missing = append(summary.Missing, name)
	}
	for name := range locked {
		summary.Locked = append(summary.Locked, name)
	}
	sort.Strings(summary.Missing)
	sort.Strings(summary.Locked)
	return summary
}

func ValidateResult(result QueryResult) error {
	if result.Account.EmployeeID == "" {
		return fmt.Errorf("result has no employee")
	}
	if !domain.ScopeEqual(result.SourceScopes, domain.AuthorizedDomains) {
		return fmt.Errorf("result has unexpected scope")
	}
	return nil
}

func ResultIDs(results []QueryResult) []string {
	ids := make([]string, 0, len(results))
	for _, result := range results {
		ids = append(ids, result.Account.EmployeeID)
	}
	sort.Strings(ids)
	return ids
}
