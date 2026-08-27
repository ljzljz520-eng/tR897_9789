package main

import (
	"path/filepath"
	"testing"
)

func TestRunHealthCommand(t *testing.T) {
	if err := run([]string{"health", "--store", filepath.Join(t.TempDir(), "health.db")}); err != nil {
		t.Fatal(err)
	}
}

func TestRunRejectsMissingCommand(t *testing.T) {
	if err := run(nil); err != nil {
		t.Fatal(err)
	}
}
