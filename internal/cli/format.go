package cli

import (
	"fmt"
	"strings"

	"callcentertroubleshooter/internal/domain"
)

func RenderDomainStatus(status domain.DomainStatus) string {
	return fmt.Sprintf("%s exists=%t locked=%t last_login=%s", domain.DisplayLabel(status.Domain), status.Exists, status.Locked, status.LastLogin)
}

func RenderStatuses(statuses []domain.DomainStatus) string {
	parts := make([]string, 0, len(statuses))
	for _, status := range statuses {
		parts = append(parts, RenderDomainStatus(status))
	}
	return strings.Join(parts, " | ")
}

func RenderScope(scope []domain.DomainName) string {
	parts := make([]string, 0, len(scope))
	for _, name := range scope {
		parts = append(parts, domain.DisplayLabel(name))
	}
	return strings.Join(parts, ",")
}

func RenderEmployee(employeeID string, statuses []domain.DomainStatus) string {
	return fmt.Sprintf("employee=%s %s", employeeID, RenderStatuses(statuses))
}
