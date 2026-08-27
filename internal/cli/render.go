package cli

import (
	"fmt"
	"strings"

	"callcentertroubleshooter/internal/app"
	"callcentertroubleshooter/internal/report"
)

func RenderWorkflow(result app.WorkflowResult) string {
	if result.Output != "" {
		return result.Output
	}
	if len(result.History) > 0 {
		return report.FormatHistory(result.History)
	}
	return "no result"
}

func RenderError(err error) string {
	if err == nil {
		return ""
	}
	return "error: " + strings.TrimSpace(err.Error())
}

func RenderHelp() string {
	return fmt.Sprintf("%s\nCommands inspect only office, telephony, and quality status.", Usage())
}

func RenderHealth(snapshot report.HealthSnapshot) string {
	return report.FormatHealth(snapshot)
}
