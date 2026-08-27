package domain

import "testing"

func TestEmployeeValidation(t *testing.T) {
	if ValidateEmployeeID("bad") == nil || ValidateEmployeeID("000000") == nil {
		t.Fatal("expected validation errors")
	}
	if ValidateEmployeeID("100001") != nil {
		t.Fatal("expected valid employee")
	}
}

func TestAccountHealthAndCounts(t *testing.T) {
	status, _ := NewDomainStatus(OfficeDomain, true, false, "2026-08-20T08:00:00Z")
	account := AgentAccount{EmployeeID: "100001", Office: status, Telephony: status, Quality: status}
	counts := CountStatuses(account)
	if counts.Existing != 3 || AccountHealth(account) != "ready" || !IsComplete(account) {
		t.Fatalf("unexpected metrics: %#v", counts)
	}
}
