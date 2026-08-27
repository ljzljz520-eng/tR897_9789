package app

import (
	"fmt"
	"sort"

	"callcentertroubleshooter/internal/domain"
)

type ServiceMetrics struct {
	Records int
	Audits  int
	Health  string
	Policy  string
}

func (s *Service) Metrics() (ServiceMetrics, error) {
	records, err := s.store.CountRecords()
	if err != nil {
		return ServiceMetrics{}, err
	}
	audits, err := s.store.CountAudits()
	if err != nil {
		return ServiceMetrics{}, err
	}
	snapshot, err := s.HealthSummary()
	if err != nil {
		return ServiceMetrics{}, err
	}
	return ServiceMetrics{Records: records, Audits: audits, Health: reportHealth(snapshot.StorageOK, snapshot.FixturesOK), Policy: s.QueryPolicy()}, nil
}

func reportHealth(storageOK, fixturesOK bool) string {
	if storageOK && fixturesOK {
		return "healthy"
	}
	return "degraded"
}

func (s *Service) DomainCounts(record domain.TroubleshootingRecord) map[string]int {
	counts := make(map[string]int)
	for _, status := range record.Snapshot.Statuses() {
		if status.Exists {
			counts[string(status.Domain)]++
		}
	}
	return counts
}

func FormatMetrics(metrics ServiceMetrics) string {
	keys := []string{"audits", "health", "policy", "records"}
	sort.Strings(keys)
	return fmt.Sprintf("records=%d audits=%d health=%s policy=%s", metrics.Records, metrics.Audits, metrics.Health, metrics.Policy)
}
