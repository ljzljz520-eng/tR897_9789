package report

import (
	"fmt"
	"sort"
	"strings"

	"callcentertroubleshooter/internal/domain"
)

type DomainSummary struct {
	Domain  domain.DomainName
	Exists  int
	Locked  int
	Missing int
}

func SummarizeDomains(records []domain.TroubleshootingRecord) []DomainSummary {
	counts := make(map[domain.DomainName]*DomainSummary)
	for _, name := range domain.AuthorizedDomains {
		counts[name] = &DomainSummary{Domain: name}
	}
	for _, record := range records {
		for _, status := range record.Snapshot.Statuses() {
			summary := counts[status.Domain]
			if status.Exists {
				summary.Exists++
			} else {
				summary.Missing++
			}
			if status.Locked {
				summary.Locked++
			}
		}
	}
	result := make([]DomainSummary, 0, len(counts))
	for _, summary := range counts {
		result = append(result, *summary)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Domain < result[j].Domain })
	return result
}

func FormatDomainSummary(summaries []DomainSummary) string {
	lines := make([]string, 0, len(summaries))
	for _, summary := range summaries {
		lines = append(lines, fmt.Sprintf("%s exists=%d missing=%d locked=%d", summary.Domain, summary.Exists, summary.Missing, summary.Locked))
	}
	return strings.Join(lines, "\n")
}

func CriticalDomains(summaries []DomainSummary) []domain.DomainName {
	result := make([]domain.DomainName, 0)
	for _, summary := range summaries {
		if summary.Locked > 0 || summary.Missing > 0 {
			result = append(result, summary.Domain)
		}
	}
	return result
}

func SummaryByDomain(summaries []DomainSummary, name domain.DomainName) (DomainSummary, bool) {
	for _, summary := range summaries {
		if summary.Domain == name {
			return summary, true
		}
	}
	return DomainSummary{}, false
}
