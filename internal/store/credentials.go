package store

import (
	"database/sql"
	"fmt"
	"sort"

	"callcentertroubleshooter/internal/domain"
	"callcentertroubleshooter/internal/security"
)

func (s *SQLiteStore) SaveCredential(employeeID string, credential security.EncryptedCredential) error {
	if employeeID == "" || credential.Domain == "" || credential.Ciphertext == "" {
		return fmt.Errorf("credential identity and ciphertext are required")
	}
	if !domain.IsAuthorizedDomain(domain.DomainName(credential.Domain)) {
		return fmt.Errorf("credential domain is not authorized")
	}
	return s.withTransaction(func(tx *sql.Tx) error {
		_, err := tx.Exec(`INSERT INTO domain_credentials (employee_id, domain, ciphertext) VALUES (?, ?, ?) ON CONFLICT(employee_id, domain) DO UPDATE SET ciphertext=excluded.ciphertext`, employeeID, credential.Domain, credential.Ciphertext)
		return err
	})
}

func (s *SQLiteStore) CredentialDomains(employeeID string) ([]string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return nil, fmt.Errorf("store is closed")
	}
	rows, err := s.db.Query(`SELECT domain FROM domain_credentials WHERE employee_id = ? ORDER BY domain`, employeeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	domains := make([]string, 0)
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		domains = append(domains, name)
	}
	sort.Strings(domains)
	return domains, rows.Err()
}

func (s *SQLiteStore) CredentialStored(employeeID string, domainName domain.DomainName) (bool, error) {
	domains, err := s.CredentialDomains(employeeID)
	if err != nil {
		return false, err
	}
	for _, name := range domains {
		if name == string(domainName) {
			return true, nil
		}
	}
	return false, nil
}
