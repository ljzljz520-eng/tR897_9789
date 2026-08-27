package domain

import (
	"fmt"
	"strings"
)

type AccessLevel string

const (
	AccessRead  AccessLevel = "read"
	AccessWrite AccessLevel = "write"
)

type AccessRequest struct {
	Actor      string
	EmployeeID string
	Scope      []DomainName
	Level      AccessLevel
}

func (r AccessRequest) Validate() error {
	if strings.TrimSpace(r.Actor) == "" {
		return fmt.Errorf("actor is required")
	}
	if err := ValidateEmployeeID(r.EmployeeID); err != nil {
		return err
	}
	if len(r.Scope) == 0 {
		return fmt.Errorf("access scope is empty")
	}
	for _, name := range r.Scope {
		if !IsAuthorizedDomain(name) {
			return fmt.Errorf("scope %s is not authorized", name)
		}
	}
	if r.Level != AccessRead && r.Level != AccessWrite {
		return fmt.Errorf("unsupported access level %s", r.Level)
	}
	return nil
}

func CanRead(request AccessRequest, name DomainName) bool {
	if request.Level != AccessRead && request.Level != AccessWrite {
		return false
	}
	for _, allowed := range request.Scope {
		if allowed == name {
			return true
		}
	}
	return false
}

func AuthorizedFieldNames() map[string]bool {
	return map[string]bool{"employee_id": true, "office": true, "telephony": true, "quality": true, "queried_at": true}
}

func IsAllowedField(name string) bool {
	return AuthorizedFieldNames()[name]
}

func ExplainAccess(request AccessRequest) string {
	parts := make([]string, 0, len(request.Scope))
	for _, name := range request.Scope {
		parts = append(parts, string(name))
	}
	return strings.Join([]string{strings.ToLower(request.Actor), request.EmployeeID, string(request.Level), strings.Join(parts, ",")}, "|")
}

func RestrictAccount(account AgentAccount, request AccessRequest) AgentAccount {
	result := AgentAccount{EmployeeID: account.EmployeeID, QueriedAt: account.QueriedAt}
	for _, status := range account.Statuses() {
		if CanRead(request, status.Domain) {
			switch status.Domain {
			case OfficeDomain:
				result.Office = status
			case TelephonyDomain:
				result.Telephony = status
			case QualityDomain:
				result.Quality = status
			}
		}
	}
	return result
}
