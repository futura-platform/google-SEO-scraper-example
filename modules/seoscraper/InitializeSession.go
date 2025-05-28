package seoscraper

import (
	"context"
	"fmt"
	"math/rand/v2"
	"net/http"
	"net/url"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/futura-platform/protocol/flowprotocol"
	"github.com/stack1ng/chromedp"
)

func (t *Task) InitializeSession() flowprotocol.TaskStepResult {
	b, cancel, err := t.SpawnSingleTabBrowser(t, t.GetProxy())
	if err != nil {
		return t.ReturnFatalErrorf("failed to spawn browser: %w", err)
	}
	defer cancel()

	u, _ := url.Parse("https://www.google.com/search?" + url.Values{
		"q": []string{fmt.Sprint(rand.Int())},
	}.Encode())

	// retrieve the session cookies
	var cookies []*network.Cookie
	err = chromedp.Run(b.CTX,
		chromedp.Navigate(u.String()),
		chromedp.Sleep(2*time.Second), // wait for challenge to complete
		chromedp.ActionFunc(func(ctx context.Context) error {
			var err error
			cookies, err = network.GetCookies().Do(ctx)
			return err
		}),
	)
	if err != nil {
		return t.ReturnSmallErrorf("failed to navigate to Google: %w", err)
	}

	httpCookies := make([]*http.Cookie, len(cookies))
	for i, cookie := range cookies {
		httpCookies[i] = &http.Cookie{
			Name:    cookie.Name,
			Value:   cookie.Value,
			Domain:  cookie.Domain,
			Path:    cookie.Path,
			Expires: time.Unix(int64(cookie.Expires), 0),
		}
	}
	t.GetCookieJar().SetCookies(u, httpCookies)

	return t.ReturnBasicStepSuccess()
}
