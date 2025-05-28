package seoscraper

import (
	"fmt"
	"strings"

	"github.com/futura-platform/protocol/flowprotocol"
	"github.com/go-resty/resty/v2"
)

const (
	upSymbol      = "⬆️"
	downSymbol    = "⬇️"
	neutralSymbol = "➖"
	newSymbol     = "🆕"
)

var webhookClient = resty.New()

func (t *Task) ReportResults() flowprotocol.TaskStepResult {
	var report strings.Builder
	didChange := t.lastTopLevelSearchResults == nil
	for i1, latestResult := range t.topLevelSearchResults {
		var rankChange *int
		for i2, lastResult := range t.lastTopLevelSearchResults {
			if latestResult.String() == lastResult.String() {
				delta := i2 - i1
				rankChange = &delta
				break
			}
		}
		statusText := newSymbol
		if rankChange != nil {
			rc := *rankChange
			if rc != 0 {
				didChange = true
			}

			if rc < 0 {
				statusText = downSymbol
			} else if rc > 0 {
				statusText = upSymbol
			} else {
				statusText = neutralSymbol
			}
		}
		report.WriteString(fmt.Sprintf("%s - <%s>\n", statusText, latestResult.String()))
	}

	if didChange {
		// report the change
		resp, err := webhookClient.R().
			SetHeader("Content-Type", "application/json").
			SetBody(map[string]string{
				"content": report.String(),
			}).
			Post(t.Params.ResultsWebhook)
		if err != nil {
			return t.ReturnSmallErrorf("failed to request webhook: %w", err)
		} else if resp.StatusCode() != 204 {
			return t.ReturnSmallErrorf("bad webhook status code %d", resp.StatusCode())
		}
	}

	t.lastTopLevelSearchResults = t.topLevelSearchResults

	// wait and scrape again
	t.Sleep(t.GetErrorDelay())
	return flowprotocol.TaskStepResult{
		NextStepLabel: "FetchSearchResults",
	}
}
