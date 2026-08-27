package app

import (
	"fmt"
	"sort"

	"callcentertroubleshooter/internal/domain"
)

type QueueItem struct {
	EmployeeID string
	Actor      string
	Priority   int
}

func BuildQueue(records []domain.TroubleshootingRecord) []QueueItem {
	queue := make([]QueueItem, 0, len(records))
	for _, record := range records {
		priority := 1
		if len(record.Snapshot.LockedDomains()) > 0 {
			priority = 3
		} else if len(record.Snapshot.MissingDomains()) > 0 {
			priority = 2
		}
		queue = append(queue, QueueItem{EmployeeID: record.EmployeeID, Actor: record.RequestedBy, Priority: priority})
	}
	sort.Slice(queue, func(i, j int) bool {
		if queue[i].Priority == queue[j].Priority {
			return queue[i].EmployeeID < queue[j].EmployeeID
		}
		return queue[i].Priority > queue[j].Priority
	})
	return queue
}

func FormatQueue(queue []QueueItem) string {
	result := ""
	for index, item := range queue {
		if index > 0 {
			result += "\n"
		}
		result += fmt.Sprintf("%s actor=%s priority=%d", item.EmployeeID, item.Actor, item.Priority)
	}
	return result
}

func (s *Service) EscalationQueue() (string, error) {
	records, err := s.store.ListRecords("")
	if err != nil {
		return "", err
	}
	return FormatQueue(BuildQueue(records)), nil
}

func QueueNeedsAttention(queue []QueueItem) bool {
	for _, item := range queue {
		if item.Priority > 1 {
			return true
		}
	}
	return false
}
