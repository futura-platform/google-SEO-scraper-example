package seoscraper

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/futura-platform/protocol/flowprotocol"
)

func (t *Task) FetchSearchResults() flowprotocol.TaskStepResult {
	resp, err := t.Get("https://www.google.com/search?"+url.Values{
		"q": []string{string(*t.Params.SearchTerm)},
	}.Encode(), getHeaders())
	if err != nil {
		return t.ReturnSmallErrorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return t.ReturnFatalErrorf("bad status code: %d", resp.StatusCode)
	}

	// parse the response body with goquery
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return t.ReturnSmallErrorf("failed to parse response body: %w", err)
	}

	// detect if we received the JS challenge page, if so then retry initialization
	if doc.Find("script").FilterFunction(func(i int, s *goquery.Selection) bool {
		return strings.Contains(s.Text(), `SG_SS=`)
	}).Length() > 0 {
		r := t.ReturnSmallErrorf("encountered JS challenge page")
		r.NextStepLabel = "InitializeSession"
		return r
	}

	// find the top-level search result anchors
	resultAnchors := doc.Find("a[href]").FilterFunction(func(i int, s *goquery.Selection) bool {
		first := s.Children().First()
		return goquery.NodeName(first) == "h3"
	})
	t.topLevelSearchResults = make([]*url.URL, resultAnchors.Length())
	for i, a := range resultAnchors.EachIter() {
		href, exists := a.Attr("href")
		if !exists {
			return t.ReturnSmallErrorf("search result anchor does not have href attribute")
		}
		u, err := url.Parse(href)
		if err != nil {
			return t.ReturnSmallErrorf("failed to parse search result URL: %w", err)
		}
		t.topLevelSearchResults[i] = u
	}

	t.SmallSuccessf("fetched %d top-level search results for term '%s'", len(t.topLevelSearchResults), t.Params.SearchTerm)
	for _, u := range t.topLevelSearchResults {
		fmt.Println(" -", u.String())
	}

	return t.ReturnBasicStepSuccess()
}
