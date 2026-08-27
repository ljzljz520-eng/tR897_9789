package fixtures

import (
	"fmt"
	"sync"

	"callcentertroubleshooter/internal/domain"
	"callcentertroubleshooter/internal/security"
)

type DirectoryEntry struct {
	EmployeeID         string
	DisplayName        string
	Office             domain.DomainStatus
	Telephony          domain.DomainStatus
	Quality            domain.DomainStatus
	PayrollGroup       string
	SecurityGroups     []string
	HomeDirectory      string
	EncryptedPasswords map[domain.DomainName]security.EncryptedCredential
}

type Repository struct {
	mu      sync.RWMutex
	entries map[string]DirectoryEntry
}

func NewRepository() *Repository {
	repo := &Repository{entries: make(map[string]DirectoryEntry)}
	repo.seed()
	return repo
}

func (r *Repository) seed() {
	vault, _ := security.NewVault("call-center-fixture-vault-v1")
	r.entries["100001"] = DirectoryEntry{EmployeeID: "100001", DisplayName: "Li Wei", Office: status(domain.OfficeDomain, true, false, "2026-08-20T08:10:00Z"), Telephony: status(domain.TelephonyDomain, true, true, "2026-08-20T08:11:00Z"), Quality: status(domain.QualityDomain, true, false, "2026-08-19T17:40:00Z"), PayrollGroup: "CC-NORTH", SecurityGroups: []string{"AD-STAFF", "AD-CC"}, HomeDirectory: "\\fileserver\\li.wei", EncryptedPasswords: encryptedPasswords(vault, "100001")}
	r.entries["100002"] = DirectoryEntry{EmployeeID: "100002", DisplayName: "Chen Yu", Office: status(domain.OfficeDomain, true, false, "2026-08-21T09:00:00Z"), Telephony: status(domain.TelephonyDomain, true, false, "2026-08-21T09:01:00Z"), Quality: status(domain.QualityDomain, false, false, ""), PayrollGroup: "CC-SOUTH", SecurityGroups: []string{"AD-STAFF"}, HomeDirectory: "\\fileserver\\chen.yu", EncryptedPasswords: encryptedPasswords(vault, "100002")}
	r.entries["100003"] = DirectoryEntry{EmployeeID: "100003", DisplayName: "Zhao Lin", Office: status(domain.OfficeDomain, true, true, "2026-08-18T07:20:00Z"), Telephony: status(domain.TelephonyDomain, true, false, "2026-08-18T07:21:00Z"), Quality: status(domain.QualityDomain, true, false, "2026-08-17T15:35:00Z"), PayrollGroup: "CC-EAST", SecurityGroups: []string{"AD-STAFF", "AD-QA"}, HomeDirectory: "\\fileserver\\zhao.lin", EncryptedPasswords: encryptedPasswords(vault, "100003")}
}

func encryptedPasswords(vault *security.Vault, employeeID string) map[domain.DomainName]security.EncryptedCredential {
	result := make(map[domain.DomainName]security.EncryptedCredential)
	for _, name := range domain.AuthorizedDomains {
		credential, _ := vault.Encrypt(string(name), employeeID, employeeID+"-"+string(name)+"-password")
		result[name] = credential
	}
	return result
}

func status(name domain.DomainName, exists, locked bool, login string) domain.DomainStatus {
	item, _ := domain.NewDomainStatus(name, exists, locked, login)
	return item
}

func (r *Repository) Lookup(employeeID string) (DirectoryEntry, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	entry, ok := r.entries[employeeID]
	if !ok {
		return DirectoryEntry{}, fmt.Errorf("fixture account %s not found", employeeID)
	}
	return cloneEntry(entry), nil
}

func (r *Repository) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.entries)
}

func (r *Repository) IDs() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	ids := make([]string, 0, len(r.entries))
	for id := range r.entries {
		ids = append(ids, id)
	}
	return ids
}

func cloneEntry(entry DirectoryEntry) DirectoryEntry {
	entry.SecurityGroups = append([]string(nil), entry.SecurityGroups...)
	entry.EncryptedPasswords = cloneCredentials(entry.EncryptedPasswords)
	return entry
}

func cloneCredentials(input map[domain.DomainName]security.EncryptedCredential) map[domain.DomainName]security.EncryptedCredential {
	output := make(map[domain.DomainName]security.EncryptedCredential, len(input))
	for name, credential := range input {
		output[name] = credential
	}
	return output
}
