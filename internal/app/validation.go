package app

import (
	"fmt"
	"strings"

	"callcentertroubleshooter/internal/domain"
)

type Request struct {
	EmployeeID string
	Actor      string
}

func ParseRequest(employeeID, actor string) (Request, error) {
	request := Request{EmployeeID: domain.NormalizeEmployeeID(employeeID), Actor: strings.TrimSpace(actor)}
	if err := domain.ValidateEmployeeID(request.EmployeeID); err != nil {
		return Request{}, err
	}
	if request.Actor == "" {
		return Request{}, fmt.Errorf("actor is required")
	}
	return request, nil
}

func ValidateServiceState(service *Service) error {
	if service == nil || service.queryer == nil || service.store == nil || service.catalog == nil {
		return fmt.Errorf("service is not ready")
	}
	return nil
}

func DescribeError(err error) string {
	if err == nil {
		return ""
	}
	return strings.TrimSpace(err.Error())
}

func IsRecoverable(err error) bool {
	if err == nil {
		return false
	}
	text := strings.ToLower(err.Error())
	return strings.Contains(text, "required") || strings.Contains(text, "invalid")
}
