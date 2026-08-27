package fixtures

import "testing"

func TestCatalogProvidesDeterministicFixtures(t *testing.T) {
	catalog := NewCatalog(NewRepository())
	first, err := catalog.Query("100001")
	if err != nil {
		t.Fatal(err)
	}
	second, err := catalog.Query("100001")
	if err != nil {
		t.Fatal(err)
	}
	if first.Office.LastLogin != second.Office.LastLogin || first.PayrollGroup != "CC-NORTH" {
		t.Fatalf("fixture changed: %#v %#v", first, second)
	}
}

func TestCatalogSearchIsSorted(t *testing.T) {
	catalog := NewCatalog(NewRepository())
	entries := catalog.SearchByDisplayName("")
	if len(entries) != 3 || entries[0].EmployeeID != "100001" {
		t.Fatalf("unexpected entries: %#v", entries)
	}
}
