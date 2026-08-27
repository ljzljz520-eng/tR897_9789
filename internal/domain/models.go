package domain

import (
	"errors"
	"fmt"
)

type DomainName string

const (
	OfficeDomain    DomainName = "office"
	TelephonyDomain DomainName = "telephony"
	QualityDomain   DomainName = "quality"
)

var AuthorizedDomains = []DomainName{OfficeDomain, TelephonyDomain, QualityDomain}

type DomainStatus struct {
	Domain    DomainName `json:"domain"`
	Exists    bool       `json:"exists"`
	Locked    bool       `json:"locked"`
	LastLogin string     `json:"last_login"`
}

type AgentAccount struct {
	EmployeeID string
	Office     DomainStatus
	Telephony  DomainStatus
	Quality    DomainStatus
	QueriedAt  string
}

type TroubleshootingRecord struct {
	ID          string
	EmployeeID  string
	RequestedBy string
	Scope       []DomainName
	Snapshot    AgentAccount
	CreatedAt   string
}

type AuditEvent struct {
	ID       string
	RecordID string
	Actor    string
	Action   string
	Occurred string
}

func NewDomainStatus(name DomainName, exists, locked bool, lastLogin string) (DomainStatus, error) {
	if !IsAuthorizedDomain(name) {
		return DomainStatus{}, fmt.Errorf("unsupported domain %q", name)
	}
	if exists && lastLogin == "" {
		return DomainStatus{}, errors.New("existing account requires last login")
	}
	if !exists && locked {
		return DomainStatus{}, errors.New("missing account cannot be locked")
	}
	return DomainStatus{Domain: name, Exists: exists, Locked: locked, LastLogin: lastLogin}, nil
}

func IsAuthorizedDomain(name DomainName) bool {
	for _, allowed := range AuthorizedDomains {
		if name == allowed {
			return true
		}
	}
	return false
}

func (a AgentAccount) Statuses() []DomainStatus {
	return []DomainStatus{a.Office, a.Telephony, a.Quality}
}

func (a AgentAccount) MissingDomains() []DomainName {
	missing := make([]DomainName, 0, len(AuthorizedDomains))
	for _, status := range a.Statuses() {
		if !status.Exists {
			missing = append(missing, status.Domain)
		}
	}
	return missing
}

func (a AgentAccount) LockedDomains() []DomainName {
	locked := make([]DomainName, 0, len(AuthorizedDomains))
	for _, status := range a.Statuses() {
		if status.Locked {
			locked = append(locked, status.Domain)
		}
	}
	return locked
}

func (r TroubleshootingRecord) Validate() error {
	if r.ID == "" || r.EmployeeID == "" || r.RequestedBy == "" {
		return errors.New("record identity is required")
	}
	if len(r.Scope) != len(AuthorizedDomains) {
		return errors.New("record scope must include all authorized domains")
	}
	for _, name := range r.Scope {
		if !IsAuthorizedDomain(name) {
			return fmt.Errorf("record scope contains %q", name)
		}
	}
	return nil
}

func (e AuditEvent) Validate() error {
	if e.ID == "" || e.RecordID == "" || e.Actor == "" || e.Action == "" {
		return errors.New("audit event fields are required")
	}
	return nil
}
