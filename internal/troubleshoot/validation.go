package troubleshoot

import (
	"fmt"
	"strings"

	"callcentertroubleshooter/internal/domain"
)

func ValidateRequest(employeeID, actor string) error {
	if err := domain.ValidateEmployeeID(employeeID); err != nil {
		return err
	}
	if strings.TrimSpace(actor) == "" {
		return fmt.Errorf("requesting administrator is required")
	}
	return nil
}

func NormalizeActor(actor string) string {
	return strings.ToLower(strings.TrimSpace(actor))
}

func RequestSummary(employeeID, actor string, policy ScopePolicy) string {
	return fmt.Sprintf("employee=%s actor=%s scope=%s", domain.NormalizeEmployeeID(employeeID), NormalizeActor(actor), policy.String())
}

func ValidateBatch(employeeIDs []string) error {
	if len(employeeIDs) == 0 {
		return fmt.Errorf("at least one employee is required")
	}
	for _, employeeID := range employeeIDs {
		if err := domain.ValidateEmployeeID(employeeID); err != nil {
			return fmt.Errorf("employee %s: %w", employeeID, err)
		}
	}
	return nil
}
