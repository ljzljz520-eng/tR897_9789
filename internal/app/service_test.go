package app

import (
	"path/filepath"
	"testing"

	"callcentertroubleshooter/internal/fixtures"
	"callcentertroubleshooter/internal/store"
	"callcentertroubleshooter/internal/troubleshoot"
)

func testService(t *testing.T) *Service {
	t.Helper()
	db, err := store.Open(filepath.Join(t.TempDir(), "service.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })
	catalog := fixtures.NewCatalog(fixtures.NewRepository())
	queryer, err := troubleshoot.NewQueryer(catalog, troubleshoot.DefaultScopePolicy())
	if err != nil {
		t.Fatal(err)
	}
	service, err := NewService(queryer, db, catalog)
	if err != nil {
		t.Fatal(err)
	}
	return service
}

func TestPrimaryWorkflow(t *testing.T) {
	service := testService(t)
	result, err := service.RunDiagnosisWorkflow("100001", "Admin.One")
	if err != nil {
		t.Fatal(err)
	}
	if result.Record.ID != "TR-000001" || result.Record.EmployeeID != "100001" {
		t.Fatalf("unexpected record: %#v", result.Record)
	}
	if len(result.Record.Scope) != 3 {
		t.Fatalf("unexpected scope: %#v", result.Record.Scope)
	}
}

func TestServiceRejectsMalformedRequest(t *testing.T) {
	service := testService(t)
	if _, err := service.Diagnose("bad", "admin"); err == nil {
		t.Fatal("expected malformed request error")
	}
}
