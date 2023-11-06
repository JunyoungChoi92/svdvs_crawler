package validator

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strings"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

// how to wait until loading JS process is done?
func CheckRedirection(tctx context.Context, targetUrl string) (string, error) {
	// Waiting until the browser has finished following all consecutive redirects and has loaded the final destination page.
	var currentURL string

	err := chromedp.Run(tctx,
		chromedp.Navigate(targetUrl),
		chromedp.WaitReady("body", chromedp.ByQuery), // Wait for the body to be loaded
		chromedp.Location(&currentURL),
	)
	if err == nil {
		return currentURL, nil
	}

	return currentURL, err
}

func CheckHttpsClientError(tctx context.Context, currentURL string) error {
	var statusCode int64
	// By setting up the listener before navigating, we ensure that we're ready to capture the network.EventResponseReceived event as soon as it occurs.
	// ev is a interface type argument of the callback function on ListenTarget(). Check type of ev using type switch.
	chromedp.ListenTarget(tctx, func(ev interface{}) {
		switch ev := ev.(type) {
		case *network.EventResponseReceived:
			statusCode = ev.Response.Status
		}
	})

	// Once the listener is set up, we then navigate to the target URL.
	err := chromedp.Run(tctx, chromedp.Navigate(currentURL))
	if err != nil {
		return err
	}

	// if statusCode >= 400, return error message.
	if statusCode >= 400 {
		message := fmt.Sprintf("HTTP client error with status code: %d", statusCode)
		return errors.New(message)
	}

	return nil
}

// My assumption is here. If some commercially available captcha API calls appear in the HTML code, the page is probably a captcha page.
func CheckBlockedByCaptcha(tctx context.Context, currentURL string, indicators []string) error {
	// consider recaptcha, hcaptcha(cloudflare has been using this after 2020),
	// ref. https://blog.cloudflare.com/moving-from-recaptcha-to-hcaptcha/

	var htmlContent string
	err := chromedp.Run(tctx,
		chromedp.Navigate(currentURL),
		chromedp.OuterHTML("html", &htmlContent),
	)

	if err != nil {
		return err
	}

	for _, indicator := range indicators {
		if strings.Contains(strings.ToLower(htmlContent), indicator) {
			return errors.New("blocked by captcha")
		}
	}

	return nil
}

func CheckURLConnectivity(url string) error {
	_, err := net.LookupHost(url)
	if err != nil {
		return err
	}

	return nil
}
