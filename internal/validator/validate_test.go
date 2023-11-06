package validator

import (
	"context"
	"testing"

	"github.com/chromedp/cdproto/browser"
	"github.com/chromedp/chromedp"
)

func TestValidateCheck(t *testing.T) {
	//# Redirection
	ctx, cancel, err := InitializeBrowserContext()
	if err != nil {
		t.Errorf("Error occurred while initializing browser context. \n\n%s", err)
	}
	defer cancel()

	t.Run("test to redirection check with pronhub", func(t *testing.T) {
		targetUrl := `https://pronhub.com/`
		expectedRedirectedURL := "https://www.pornhub.com/?utm_source=pronhub.com&utm_medium=redirects&utm_campaign=tldredirects"

		redirectedURL, err := redirectionURLHelper(t, ctx, targetUrl)
		if err != nil {
			t.Errorf("Error occurred while redirection checking %s. \n\n%s", targetUrl, err)
		}

		if redirectedURL != expectedRedirectedURL {
			t.Errorf("redirected URL and expected Redirected URL is not same. \nredirectedURL: %s\nexpectedRedirectedURL:%s", redirectedURL, expectedRedirectedURL)
		}

		err = validationHelper(t, ctx, redirectedURL)
		if err != nil {
			t.Errorf("Error occurred while validation checking %s. \n\n%s", redirectedURL, err)
		}
	})
	t.Run("test to redirection check with xandar", func(t *testing.T) {
		targetUrl := `http://xandar.com/`
		expectedRedirectedURL := "https://xkcorp.com/"

		redirectedURL, err := redirectionURLHelper(t, ctx, targetUrl)
		if err != nil {
			t.Errorf("Error occurred while redirection checking %s. \n\n%s", targetUrl, err)
		}

		if redirectedURL != expectedRedirectedURL {
			t.Errorf("redirected URL and expected Redirected URL is not same. \nredirectedURL: %s\nexpectedRedirectedURL:%s", redirectedURL, expectedRedirectedURL)
		}

		err = validationHelper(t, ctx, redirectedURL)
		if err != nil {
			t.Errorf("Error occurred while validation checking %s. \n\n%s", redirectedURL, err)
		}

	})
}

func redirectionURLHelper(t *testing.T, ctx context.Context, targetUrl string) (string, error) {
	currentURL, err := CheckRedirection(ctx, targetUrl)
	if err != nil {
		return currentURL, err
	}

	return currentURL, err
}

func validationHelper(t *testing.T, ctx context.Context, targetUrl string) error {
	// err := CheckURLConnectivity(targetUrl)
	// if err != nil {
	// 	return err
	// }

	err := CheckHttpsClientError(ctx, targetUrl)
	if err != nil {
		return err
	}

	err = CheckBlockedByCaptcha(ctx, targetUrl, []string{`script src="https://js.hcaptcha.com/1/api.js"`, `script src="https://www.google.com/recaptcha/api.js"`})
	if err != nil {
		return err
	}

	return err
}

func InitializeBrowserContext() (context.Context, context.CancelFunc, error) {
	// Initiate a browser context.
	config := chromedp.DefaultExecAllocatorOptions[:]
	config = append(config,
		chromedp.Flag("headless", false),
		chromedp.NoDefaultBrowserCheck,
		chromedp.Flag("disable-blink-features", "AutomationControlled"),
		chromedp.Flag("enable-automation", false),
		chromedp.NoSandbox,
	)

	execCtx, _ := chromedp.NewExecAllocator(context.Background(), config...)

	ctx, cancel := chromedp.NewContext(execCtx)

	err := chromedp.Run(ctx,
		// add network.Enable() for Monitoring Network Traffic and Performance Analysis
		// network.Enable(),
		browser.SetDownloadBehavior(browser.SetDownloadBehaviorBehaviorDeny),
		// To ensure that we're always fetching the latest version of resources,
		// chromedp.ActionFunc(func(ctx context.Context) error {
		// 	err := network.SetCacheDisabled(true).Do(ctx)
		// 	return err
		// }),
	)
	if err != nil {
		return nil, nil, err
	}

	return ctx, cancel, err
}
