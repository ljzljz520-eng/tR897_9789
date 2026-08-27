package troubleshoot

import (
	"fmt"
	"sort"

	"callcentertroubleshooter/internal/domain"
)

func ScopeResult(account domain.AgentAccount, policy ScopePolicy) (domain.AgentAccount, error) {
	if err := policy.Validate(); err != nil {
		return domain.AgentAccount{}, err
	}
	statuses := FilterStatuses(account.Statuses(), policy)
	if len(statuses) != len(domain.AuthorizedDomains) {
		return domain.AgentAccount{}, fmt.Errorf("incomplete authorized result")
	}
	return domain.AgentAccount{EmployeeID: account.EmployeeID, Office: statuses[0], Telephony: statuses[1], Quality: statuses[2], QueriedAt: account.QueriedAt}, nil
}

func AllowedFieldNames() []string {
	fields := []string{"employee_id", "office", "telephony", "quality", "queried_at"}
	sort.Strings(fields)
	return fields
}

func IsScoped(result QueryResult) bool {
	return len(result.ExtraFields) == 0
}

func DescribeScope(policy ScopePolicy) map[string]string {
	description := make(map[string]string)
	for _, name := range policy.Allowed {
		description[string(name)] = domain.DisplayLabel(name)
	}
	return description
}

func MergeScopedResults(results []QueryResult) map[string]domain.AgentAccount {
	merged := make(map[string]domain.AgentAccount, len(results))
	for _, result := range results {
		merged[result.Account.EmployeeID] = result.Account
	}
	return merged
}
