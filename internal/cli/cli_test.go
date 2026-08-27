package cli

import "testing"

func TestParseDiagnoseCommand(t *testing.T) {
	command, err := Parse([]string{"diagnose", "--employee", "100001", "--actor", "admin", "--store", "x.db"})
	if err != nil {
		t.Fatal(err)
	}
	if command.Name != "diagnose" || command.Employee != "100001" || !IsMutating(command) {
		t.Fatalf("unexpected command: %#v", command)
	}
}

func TestParseRejectsUnknownArgument(t *testing.T) {
	if _, err := Parse([]string{"health", "--nope"}); err == nil {
		t.Fatal("expected argument error")
	}
}
