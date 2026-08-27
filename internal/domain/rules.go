package domain

import (
	"errors"
	"sort"
	"strings"
)

var ErrInvalidEmployeeID = errors.New("invalid employee id")
var ErrUnknownEmployee = errors.New("unknown employee id")

func ValidateEmployeeID(employeeID string) error {
	value := strings.TrimSpace(employeeID)
	if len(value) != 6 {
		return ErrInvalidEmployeeID
	}
	for _, ch := range value {
		if ch < '0' || ch > '9' {
			return ErrInvalidEmployeeID
		}
	}
	if value == "000000" {
		return ErrUnknownEmployee
	}
	return nil
}

func NormalizeEmployeeID(employeeID string) string {
	return strings.TrimSpace(employeeID)
}

func CanonicalScope() []DomainName {
	copyOfScope := append([]DomainName(nil), AuthorizedDomains...)
	return copyOfScope
}

func ScopeNames(scope []DomainName) []string {
	names := make([]string, 0, len(scope))
	for _, item := range scope {
		names = append(names, string(item))
	}
	sort.Strings(names)
	return names
}

func ScopeEqual(left, right []DomainName) bool {
	if len(left) != len(right) {
		return false
	}
	for _, name := range left {
		found := false
		for _, candidate := range right {
			if name == candidate {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func DisplayLabel(name DomainName) string {
	switch name {
	case OfficeDomain:
		return "办公域"
	case TelephonyDomain:
		return "话务域"
	case QualityDomain:
		return "质检域"
	default:
		return "未知域"
	}
}

func AccountHealth(account AgentAccount) string {
	if len(account.MissingDomains()) > 0 {
		return "incomplete"
	}
	if len(account.LockedDomains()) > 0 {
		return "attention"
	}
	return "ready"
}
