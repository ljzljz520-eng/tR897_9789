package report

import (
	"testing"

	"callcentertroubleshooter/internal/domain"
)

func TestHistoryFormatting(t *testing.T) {
	status, _ := domain.NewDomainStatus(domain.OfficeDomain, true, false, "2026-08-20T08:00:00Z")
	record := domain.TroubleshootingRecord{ID: "TR-1", EmployeeID: "100001", RequestedBy: "admin", Snapshot: domain.AgentAccount{EmployeeID: "100001", Office: status, Telephony: status, Quality: status}, CreatedAt: "2026-08-24T13:00:00Z"}
	items := BuildHistory([]domain.TroubleshootingRecord{record}, []domain.AuditEvent{{ID: "AU-1", RecordID: "TR-1", Actor: "admin", Action: "diagnose"}})
	if len(items) != 1 || items[0].AuditCount != 1 {
		t.Fatalf("unexpected history: %#v", items)
	}
	if FormatHistory(items) == "" {
		t.Fatal("expected formatted history")
	}
}

func TestHealthFormatDeterministic(t *testing.T) {
	snapshot := BuildHealth(true, true, 3, 2, 2, nil)
	if OverallHealth(snapshot) != "healthy" || FormatHealth(snapshot) == "" {
		t.Fatalf("unexpected health: %#v", snapshot)
	}
}
