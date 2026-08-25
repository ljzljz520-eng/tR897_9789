package fixtures

import (
	"fmt"
	"sort"

	"callcentertroubleshooter/internal/domain"
)

type Diagnostic struct {
	EmployeeID string
	Domains    int
	Existing   int
	Locked     int
	Coverage   float64
}

func (c *Catalog) DiagnoseFixture(employeeID string) (Diagnostic, error) {
	entry, err := c.Query(employeeID)
	if err != nil {
		return Diagnostic{}, err
	}
	account := c.AuthorizedSnapshot(entry)
	counts := domain.CountStatuses(account)
	return Diagnostic{EmployeeID: employeeID, Domains: len(account.Statuses()), Existing: counts.Existing, Locked: counts.Locked, Coverage: domain.LoginCoverage(account)}, nil
}

func (c *Catalog) LockedEmployees() []string {
	ids := make([]string, 0)
	for _, id := range c.repository.IDs() {
		entry, err := c.Query(id)
		if err != nil {
			continue
		}
		if entry.Office.Locked || entry.Telephony.Locked || entry.Quality.Locked {
			ids = append(ids, id)
		}
	}
	sort.Strings(ids)
	return ids
}

func (c *Catalog) CoverageSummary() string {
	diagnostics := make([]string, 0)
	for _, id := range c.repository.IDs() {
		item, err := c.DiagnoseFixture(id)
		if err == nil {
			diagnostics = append(diagnostics, fmt.Sprintf("%s:%d/%d", item.EmployeeID, item.Existing, item.Domains))
		}
	}
	sort.Strings(diagnostics)
	return fmt.Sprintf("fixtures %v", diagnostics)
}

func (c *Catalog) DomainNames(entry DirectoryEntry) []string {
	names := make([]string, 0, len(domain.AuthorizedDomains))
	for _, name := range domain.AuthorizedDomains {
		names = append(names, string(name))
	}
	return names
}
