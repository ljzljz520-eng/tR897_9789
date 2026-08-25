package store

import (
	"path/filepath"
	"testing"

	"callcentertroubleshooter/internal/domain"
)

func sampleRecord() domain.TroubleshootingRecord {
	office, _ := domain.NewDomainStatus(domain.OfficeDomain, true, false, "2026-08-20T08:10:00Z")
	telephony, _ := domain.NewDomainStatus(domain.TelephonyDomain, true, true, "2026-08-20T08:11:00Z")
	quality, _ := domain.NewDomainStatus(domain.QualityDomain, true, false, "2026-08-19T17:40:00Z")
	return domain.TroubleshootingRecord{ID: "TR-000001", EmployeeID: "100001", RequestedBy: "admin", Scope: domain.CanonicalScope(), Snapshot: domain.AgentAccount{EmployeeID: "100001", Office: office, Telephony: telephony, Quality: quality, QueriedAt: "2026-08-24T13:00:00Z"}, CreatedAt: "2026-08-24T13:01:00Z"}
}

func TestPersistenceSurvivesReopen(t *testing.T) {
	path := filepath.Join(t.TempDir(), "persist.db")
	first, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	record := sampleRecord()
	if err := first.SaveRecordAndAudit(record, domain.AuditEvent{ID: "AU-000001", RecordID: record.ID, Actor: "admin", Action: "diagnose", Occurred: record.CreatedAt}); err != nil {
		t.Fatal(err)
	}
	if err := first.Close(); err != nil {
		t.Fatal(err)
	}
	second, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer second.Close()
	loaded, err := second.GetRecord(record.ID)
	if err != nil {
		t.Fatal(err)
	}
	if loaded.EmployeeID != record.EmployeeID || loaded.Snapshot.Telephony.Locked != true {
		t.Fatalf("unexpected reopened record: %#v", loaded)
	}
}

func TestStoreListsAudits(t *testing.T) {
	path := filepath.Join(t.TempDir(), "audits.db")
	db, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	record := sampleRecord()
	if err := db.SaveRecordAndAudit(record, domain.AuditEvent{ID: "AU-000001", RecordID: record.ID, Actor: "admin", Action: "diagnose", Occurred: record.CreatedAt}); err != nil {
		t.Fatal(err)
	}
	audits, err := db.ListAudits(record.ID)
	if err != nil || len(audits) != 1 {
		t.Fatalf("audits: %v %#v", err, audits)
	}
}
