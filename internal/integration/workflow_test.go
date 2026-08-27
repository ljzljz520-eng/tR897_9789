package integration

import (
	"path/filepath"
	"testing"

	"callcentertroubleshooter/internal/app"
	"callcentertroubleshooter/internal/fixtures"
	"callcentertroubleshooter/internal/store"
	"callcentertroubleshooter/internal/troubleshoot"
)

func setup(t *testing.T) *app.Service {
	t.Helper()
	db, err := store.Open(filepath.Join(t.TempDir(), "integration.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })
	catalog := fixtures.NewCatalog(fixtures.NewRepository())
	queryer, err := troubleshoot.NewQueryer(catalog, troubleshoot.DefaultScopePolicy())
	if err != nil {
		t.Fatal(err)
	}
	service, err := app.NewService(queryer, db, catalog)
	if err != nil {
		t.Fatal(err)
	}
	return service
}

func TestSecondaryWorkflow(t *testing.T) {
	service := setup(t)
	if _, err := service.Diagnose("100002", "admin"); err != nil {
		t.Fatal(err)
	}
	history, err := service.RunHistoryWorkflow("100002")
	if err != nil || len(history.History) != 1 {
		t.Fatalf("history: %v %#v", err, history)
	}
}

func TestTertiaryWorkflow(t *testing.T) {
	service := setup(t)
	if _, err := service.Diagnose("100003", "admin"); err != nil {
		t.Fatal(err)
	}
	result, err := service.RunHealthWorkflow()
	if err != nil || result.Output == "" {
		t.Fatalf("health: %v %#v", err, result)
	}
}
