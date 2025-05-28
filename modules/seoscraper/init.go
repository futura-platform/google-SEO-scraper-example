package seoscraper

import (
	"net/url"

	"github.com/futura-platform/protocol"
	basicgroupsprotocol "github.com/futura-platform/protocol/basicgroups/protocol"
	"github.com/futura-platform/protocol/flowprotocol"
	"github.com/futura-platform/protocol/netprotocol"
)

type Params struct {
	SearchTerm     basicgroupsprotocol.EntryProvided[SearchTerm]
	ResultsWebhook string
}

type Task struct {
	*protocol.Task[Params]

	topLevelSearchResults     []*url.URL
	lastTopLevelSearchResults []*url.URL
}

func getHeaders() netprotocol.OrderedHeaders {
	return netprotocol.OrderedHeaders{
		// sec-ch-ua headers will be automatically added by the netprotocol package based on the selected browser profile.
		// (the default is the latest stable chrome for mac)
		{`sec-ch-ua`},
		{`sec-ch-ua-mobile`},
		{`sec-ch-ua-platform`},
		{`upgrade-insecure-requests`, `1`},
		{`user-agent`},
		{`accept`, `text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7`},
		{`x-browser-channel`, `stable`},
		{`x-browser-year`, `2025`},
		{`x-browser-validation`, `xFGgWt/wguugeH2wFxmiwYRqxZo=`}, // values chrome sends to google for unknown reasons, hardcoding doesnt appear to cause issues so lets do it
		{`x-browser-copyright`, `Copyright 2025 Google LLC. All rights reserved.`},
		{`x-client-data`, `CIe2yQEIpLbJAQipncoBCITvygEIlKHLAQiSo8sBCIagzQEIhe3OAQjk7c4BCN7uzgEIkfHOAQit8c4BCLTxzgEYpPDOAQ==`},
		{`sec-fetch-site`, `none`},
		{`sec-fetch-mode`, `navigate`},
		{`sec-fetch-user`, `?1`},
		{`sec-fetch-dest`, `document`},
		{`accept-encoding`, `gzip, deflate, br, zstd`},
		{`accept-language`, `en-US,en;q=0.9`},
		{`cookie`},
		{`priority`, `u=0, i`},
	}
}

func ConstructTask(base *protocol.Task[Params]) (protocol.BaseTask, []flowprotocol.TaskStep, error) {
	t := &Task{
		Task: base,
	}

	return t,
		[]flowprotocol.TaskStep{
			{StepFunc: t.InitializeSession},
			{StepFunc: t.FetchSearchResults},
			{StepFunc: t.ReportResults},
		},
		nil
}
