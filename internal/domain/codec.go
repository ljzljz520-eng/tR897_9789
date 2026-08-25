package domain

import (
	"encoding/json"
	"fmt"
)

type accountWire struct {
	EmployeeID string       `json:"employee_id"`
	Office     DomainStatus `json:"office"`
	Telephony  DomainStatus `json:"telephony"`
	Quality    DomainStatus `json:"quality"`
	QueriedAt  string       `json:"queried_at"`
}

func EncodeAccount(account AgentAccount) (string, error) {
	wire := accountWire{EmployeeID: account.EmployeeID, Office: account.Office, Telephony: account.Telephony, Quality: account.Quality, QueriedAt: account.QueriedAt}
	data, err := json.Marshal(wire)
	if err != nil {
		return "", fmt.Errorf("encode account: %w", err)
	}
	return string(data), nil
}

func DecodeAccount(payload string) (AgentAccount, error) {
	var wire accountWire
	if err := json.Unmarshal([]byte(payload), &wire); err != nil {
		return AgentAccount{}, fmt.Errorf("decode account: %w", err)
	}
	if wire.EmployeeID == "" {
		return AgentAccount{}, fmt.Errorf("decode account: employee id missing")
	}
	return AgentAccount{EmployeeID: wire.EmployeeID, Office: wire.Office, Telephony: wire.Telephony, Quality: wire.Quality, QueriedAt: wire.QueriedAt}, nil
}

func EncodeScope(scope []DomainName) (string, error) {
	data, err := json.Marshal(scope)
	if err != nil {
		return "", fmt.Errorf("encode scope: %w", err)
	}
	return string(data), nil
}

func DecodeScope(payload string) ([]DomainName, error) {
	var scope []DomainName
	if err := json.Unmarshal([]byte(payload), &scope); err != nil {
		return nil, fmt.Errorf("decode scope: %w", err)
	}
	return scope, nil
}
