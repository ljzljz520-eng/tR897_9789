package domain

import (
	"fmt"
	"strings"
)

type AuditAction string

const (
	ActionDiagnose AuditAction = "diagnose"
	ActionHistory  AuditAction = "history"
	ActionExport   AuditAction = "export"
)

func ParseAuditAction(value string) (AuditAction, error) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case string(ActionDiagnose):
		return ActionDiagnose, nil
	case string(ActionHistory):
		return ActionHistory, nil
	case string(ActionExport):
		return ActionExport, nil
	default:
		return "", fmt.Errorf("unknown audit action %s", value)
	}
}

func (e AuditEvent) ActionType() (AuditAction, error) {
	return ParseAuditAction(e.Action)
}

func NewAuditEvent(id, recordID, actor, action, occurred string) (AuditEvent, error) {
	parsed, err := ParseAuditAction(action)
	if err != nil {
		return AuditEvent{}, err
	}
	event := AuditEvent{ID: id, RecordID: recordID, Actor: actor, Action: string(parsed), Occurred: occurred}
	if err := event.Validate(); err != nil {
		return AuditEvent{}, err
	}
	return event, nil
}

func AuditLabel(action AuditAction) string {
	switch action {
	case ActionDiagnose:
		return "Account diagnosis"
	case ActionHistory:
		return "History review"
	case ActionExport:
		return "Summary export"
	default:
		return "Unknown action"
	}
}

func GroupAudits(events []AuditEvent) map[string][]AuditEvent {
	groups := make(map[string][]AuditEvent)
	for _, event := range events {
		groups[event.RecordID] = append(groups[event.RecordID], event)
	}
	return groups
}
