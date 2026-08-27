package app

import (
	"fmt"
	"strings"

	"callcentertroubleshooter/internal/domain"
)

func AccessRequestFor(employeeID, actor string) domain.AccessRequest {
	return domain.AccessRequest{EmployeeID: strings.TrimSpace(employeeID), Actor: strings.ToLower(strings.TrimSpace(actor)), Scope: domain.CanonicalScope(), Level: domain.AccessRead}
}

func (s *Service) Authorize(employeeID, actor string) error {
	request := AccessRequestFor(employeeID, actor)
	if err := request.Validate(); err != nil {
		return err
	}
	if !s.catalog.HasEmployee(request.EmployeeID) {
		return fmt.Errorf("employee %s is not in the fixture directory", request.EmployeeID)
	}
	return nil
}

func (s *Service) ScopedAccount(employeeID, actor string) (domain.AgentAccount, error) {
	if err := s.Authorize(employeeID, actor); err != nil {
		return domain.AgentAccount{}, err
	}
	result, err := s.queryer.QueryAccount(employeeID)
	if err != nil {
		return domain.AgentAccount{}, err
	}
	return domain.RestrictAccount(result.Account, AccessRequestFor(employeeID, actor)), nil
}

func RequireReadAccess(request domain.AccessRequest, account domain.AgentAccount) error {
	if err := request.Validate(); err != nil {
		return err
	}
	for _, status := range account.Statuses() {
		if !domain.CanRead(request, status.Domain) {
			return fmt.Errorf("missing read access for %s", status.Domain)
		}
	}
	return nil
}

func (s *Service) AuthorizationSummary(employeeID, actor string) string {
	request := AccessRequestFor(employeeID, actor)
	return domain.ExplainAccess(request)
}
