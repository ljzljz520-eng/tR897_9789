package troubleshoot

import (
	"errors"
	"fmt"
	"strings"

	"callcentertroubleshooter/internal/domain"
)

type ScopePolicy struct {
	Allowed []domain.DomainName
}

func DefaultScopePolicy() ScopePolicy {
	return ScopePolicy{Allowed: domain.CanonicalScope()}
}

func (p ScopePolicy) Validate() error {
	if len(p.Allowed) == 0 {
		return errors.New("scope policy is empty")
	}
	seen := make(map[domain.DomainName]bool)
	for _, name := range p.Allowed {
		if !domain.IsAuthorizedDomain(name) {
			return fmt.Errorf("policy includes unauthorized domain %q", name)
		}
		if seen[name] {
			return fmt.Errorf("policy repeats domain %q", name)
		}
		seen[name] = true
	}
	return nil
}

func (p ScopePolicy) Allows(name domain.DomainName) bool {
	for _, candidate := range p.Allowed {
		if candidate == name {
			return true
		}
	}
	return false
}

func (p ScopePolicy) String() string {
	names := make([]string, 0, len(p.Allowed))
	for _, name := range p.Allowed {
		names = append(names, string(name))
	}
	return strings.Join(names, ",")
}

func FilterStatuses(statuses []domain.DomainStatus, policy ScopePolicy) []domain.DomainStatus {
	filtered := make([]domain.DomainStatus, 0, len(statuses))
	for _, status := range statuses {
		if policy.Allows(status.Domain) {
			filtered = append(filtered, status)
		}
	}
	return filtered
}
