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
	// The troubleshooting result exposes only the three authorized domains
	// (office, telephony, quality). Unrelated directory attributes such as
	// payroll group, home directory, and security groups are intentionally
	// omitted so IsScoped reports the result as properly scoped.
	return QueryResult{Account: account, DisplayName: entry.DisplayName, SourceScopes: append([]domain.DomainName(nil), q.policy.Allowed...)}, nil
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
