package report

import (
	"fmt"
	"sort"
	"strings"

	"callcentertroubleshooter/internal/domain"
	"callcentertroubleshooter/internal/fixtures"
)

type HealthSnapshot struct {
	StorageOK       bool
	FixturesOK      bool
	FixtureAccounts int
	RecordCount     int
	AuditCount      int
	StatusCounts    map[string]int
}

func BuildHealth(storageOK, fixturesOK bool, fixtureAccounts, recordCount, auditCount int, records []domain.TroubleshootingRecord) HealthSnapshot {
	items := BuildHistory(records, nil)
	return HealthSnapshot{StorageOK: storageOK, FixturesOK: fixturesOK, FixtureAccounts: fixtureAccounts, RecordCount: recordCount, AuditCount: auditCount, StatusCounts: HistoryHealth(items)}
}

func OverallHealth(snapshot HealthSnapshot) string {
	if !snapshot.StorageOK || !snapshot.FixturesOK {
		return "degraded"
	}
	return "healthy"
}

func FormatHealth(snapshot HealthSnapshot) string {
	keys := []string{"attention", "incomplete", "ready"}
	parts := make([]string, 0, len(keys))
	for _, key := range keys {
		parts = append(parts, fmt.Sprintf("%s=%d", key, snapshot.StatusCounts[key]))
	}
	return fmt.Sprintf("health=%s storage=%t fixtures=%t fixture_accounts=%d records=%d audits=%d %s", OverallHealth(snapshot), snapshot.StorageOK, snapshot.FixturesOK, snapshot.FixtureAccounts, snapshot.RecordCount, snapshot.AuditCount, strings.Join(parts, " "))
}

func FixtureDiagnostics(catalog *fixtures.Catalog) map[string]string {
	result := map[string]string{"catalog": "unavailable", "accounts": "0"}
	if catalog == nil {
		return result
	}
	if err := catalog.Health(); err != nil {
		return result
	}
	result["catalog"] = "ready"
	result["accounts"] = fmt.Sprintf("%d", len(catalog.SearchByDisplayName("")))
	return result
}

func SortedHealthKeys(counts map[string]int) []string {
	keys := make([]string, 0, len(counts))
	for key := range counts {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}
