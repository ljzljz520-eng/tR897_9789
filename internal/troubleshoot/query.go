package troubleshoot

import (
	"fmt"

	"callcentertroubleshooter/internal/domain"
	"callcentertroubleshooter/internal/fixtures"
)

type QueryResult struct {
	Account      domain.AgentAccount
	DisplayName  string
	ExtraFields  map[string]string
	SourceScopes []domain.DomainName
}

type Queryer struct {
	catalog *fixtures.Catalog
	policy  ScopePolicy
}

func NewQueryer(catalog *fixtures.Catalog, policy ScopePolicy) (*Queryer, error) {
	if catalog == nil {
		return nil, fmt.Errorf("catalog is required")
	}
	if err := policy.Validate(); err != nil {
		return nil, err
	}
	return &Queryer{catalog: catalog, policy: policy}, nil
}

func (q *Queryer) QueryAccount(employeeID string) (QueryResult, error) {
	entry, err := q.catalog.Query(employeeID)
	if err != nil {
		return QueryResult{}, err
	}
	account := q.catalog.AuthorizedSnapshot(entry)
	statuses := FilterStatuses(account.Statuses(), q.policy)
	if len(statuses) != len(domain.AuthorizedDomains) {
		return QueryResult{}, fmt.Errorf("query policy does not cover all domains")
	}
	account.Office, account.Telephony, account.Quality = statuses[0], statuses[1], statuses[2]
	result := QueryResult{Account: account, DisplayName: entry.DisplayName, SourceScopes: append([]domain.DomainName(nil), q.policy.Allowed...)}
	result.ExtraFields = map[string]string{
		"payroll_group":   entry.PayrollGroup,
		"home_directory":  entry.HomeDirectory,
		"security_groups": fmt.Sprintf("%v", entry.SecurityGroups),
	}
	return result, nil
}

func (q *Queryer) QueryMany(employeeIDs []string) ([]QueryResult, error) {
	results := make([]QueryResult, 0, len(employeeIDs))
	for _, employeeID := range employeeIDs {
		result, err := q.QueryAccount(employeeID)
		if err != nil {
			return nil, fmt.Errorf("query %s: %w", employeeID, err)
		}
		results = append(results, result)
	}
	return results, nil
}

func (q *Queryer) PolicyDescription() string {
	return q.policy.String()
}
