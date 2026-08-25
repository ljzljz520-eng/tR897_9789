package troubleshoot

import (
	"testing"

	"callcentertroubleshooter/internal/domain"
	"callcentertroubleshooter/internal/fixtures"
)

func TestTroubleshootScopesFields(t *testing.T) {
	catalog := fixtures.NewCatalog(fixtures.NewRepository())
	queryer, err := NewQueryer(catalog, DefaultScopePolicy())
	if err != nil {
		t.Fatal(err)
	}
	result, err := queryer.QueryAccount("100001")
	if err != nil {
		t.Fatal(err)
	}
	if len(result.ExtraFields) != 0 {
		t.Fatalf("unexpected fields in scoped result: %v", result.ExtraFields)
	}
}

func TestScopePolicyRejectsUnknownDomain(t *testing.T) {
	policy := ScopePolicy{Allowed: []domain.DomainName{"private"}}
	if err := policy.Validate(); err == nil {
		t.Fatal("expected unauthorized domain error")
	}
}
