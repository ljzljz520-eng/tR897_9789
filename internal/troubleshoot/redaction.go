package troubleshoot

import (
	"fmt"
	"sort"
	"strings"

	"callcentertroubleshooter/internal/domain"
)

type RedactedView struct {
	EmployeeID string
	Fields     map[string]string
}

func Redact(result QueryResult) RedactedView {
	fields := map[string]string{
		"office":    state(result.Account.Office),
		"telephony": state(result.Account.Telephony),
		"quality":   state(result.Account.Quality),
	}
	return RedactedView{EmployeeID: result.Account.EmployeeID, Fields: fields}
}

func state(status domain.DomainStatus) string {
	if !status.Exists {
		return "missing"
	}
	if status.Locked {
		return "locked"
	}
	return "active"
}

func FieldList(view RedactedView) []string {
	keys := make([]string, 0, len(view.Fields))
	for key := range view.Fields {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func FormatRedacted(view RedactedView) string {
	parts := make([]string, 0, len(view.Fields)+1)
	parts = append(parts, "employee="+view.EmployeeID)
	for _, key := range FieldList(view) {
		parts = append(parts, fmt.Sprintf("%s=%s", key, view.Fields[key]))
	}
	return strings.Join(parts, " ")
}

func ValidateRedaction(view RedactedView) error {
	if view.EmployeeID == "" {
		return fmt.Errorf("redacted view has no employee")
	}
	for key := range view.Fields {
		if key != "office" && key != "telephony" && key != "quality" {
			return fmt.Errorf("redacted view includes %s", key)
		}
	}
	return nil
}

func (q *Queryer) RedactedAccount(employeeID string) (RedactedView, error) {
	result, err := q.QueryAccount(employeeID)
	if err != nil {
		return RedactedView{}, err
	}
	view := Redact(result)
	if err := ValidateRedaction(view); err != nil {
		return RedactedView{}, err
	}
	return view, nil
}
