package fixtures

import (
	"sort"
	"strings"

	"callcentertroubleshooter/internal/domain"
)

type Catalog struct {
	repository *Repository
}

func NewCatalog(repository *Repository) *Catalog {
	if repository == nil {
		repository = NewRepository()
	}
	return &Catalog{repository: repository}
}

func (c *Catalog) Query(employeeID string) (DirectoryEntry, error) {
	return c.repository.Lookup(employeeID)
}

func (c *Catalog) HasEmployee(employeeID string) bool {
	_, err := c.Query(employeeID)
	return err == nil
}

func (c *Catalog) AuthorizedSnapshot(entry DirectoryEntry) domain.AgentAccount {
	return domain.AgentAccount{EmployeeID: entry.EmployeeID, Office: entry.Office, Telephony: entry.Telephony, Quality: entry.Quality, QueriedAt: "2026-08-24T13:00:00Z"}
}

func (c *Catalog) SearchByDisplayName(fragment string) []DirectoryEntry {
	needle := strings.ToLower(strings.TrimSpace(fragment))
	results := make([]DirectoryEntry, 0)
	for _, id := range c.repository.IDs() {
		entry, err := c.Query(id)
		if err == nil && strings.Contains(strings.ToLower(entry.DisplayName), needle) {
			results = append(results, entry)
		}
	}
	sort.Slice(results, func(i, j int) bool { return results[i].EmployeeID < results[j].EmployeeID })
	return results
}

func (c *Catalog) DomainCoverage(entry DirectoryEntry) map[domain.DomainName]bool {
	return map[domain.DomainName]bool{domain.OfficeDomain: entry.Office.Exists, domain.TelephonyDomain: entry.Telephony.Exists, domain.QualityDomain: entry.Quality.Exists}
}

func (c *Catalog) Health() error {
	if c.repository == nil || c.repository.Count() == 0 {
		return errEmptyCatalog
	}
	return nil
}

var errEmptyCatalog = &catalogError{"fixture catalog is empty"}

type catalogError struct{ message string }

func (e *catalogError) Error() string { return e.message }
