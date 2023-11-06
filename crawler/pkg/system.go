package pkg

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/chromedp/cdproto/browser"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/dev-zipida-com/spo-vdvs-crawler/internal/db"
	"github.com/dev-zipida-com/spo-vdvs-crawler/internal/db/schema"
	scrap "github.com/dev-zipida-com/spo-vdvs-crawler/internal/scrapper"
	validate "github.com/dev-zipida-com/spo-vdvs-crawler/internal/validator"
	"github.com/jmoiron/sqlx"
)

func NewDatabase() *db.Database {
	return &db.Database{}
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
		network.Enable(),
		browser.SetDownloadBehavior(browser.SetDownloadBehaviorBehaviorDeny),
		// To ensure that we're always fetching the latest version of resources, init the cache.
		chromedp.ActionFunc(func(ctx context.Context) error {
			err := network.SetCacheDisabled(true).Do(ctx)
			return err
		}),
	)
	if err != nil {
		return nil, nil, err
	}

	return ctx, cancel, err
}

func CheckURLConnectivity(url string) (bool, error) {
	resp, err := http.Get(url)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return true, nil
	}

	return false, fmt.Errorf("server responded with status code: %d", resp.StatusCode)
}

func GetRedirectedUrl(ctx context.Context, targetURL string, timeOut time.Duration) (string, error) {
	currentURL, RedirectionError := validate.CheckRedirection(ctx, targetURL)
	if RedirectionError != nil {
		return currentURL, RedirectionError
	}

	return currentURL, nil
}

// Perform a "Connection Check" to ensure the URL is active and reachable.
func CheckURLValidation(ctx context.Context, currentURL string, indicators []string) error {
	HttpsClientError := validate.CheckHttpsClientError(ctx, currentURL)
	if HttpsClientError != nil {
		return HttpsClientError
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	if strings.Contains(currentURL, "warning.or.kr") {
		return errors.New("redirected to warning.or.kr.")
	}

	CheckBlockedByCaptchaError := validate.CheckBlockedByCaptcha(ctx, currentURL, indicators)
	if CheckBlockedByCaptchaError != nil {
		return CheckBlockedByCaptchaError
	}

	return nil
}

func ScrapWebpage(ctx context.Context, targetUrl string, iframeTimeOut time.Duration, isBlockedByCaptcha bool) (string, string, error) {
	pageContent, err := scrap.ScrapHTML(ctx, targetUrl, iframeTimeOut, isBlockedByCaptcha)
	if err != nil {
		return "", "", err
	}

	// Take a full-page screenshot.
	base64TypeScreenShot, err := scrap.CaptureScreenshot(ctx, targetUrl)
	if err != nil {
		return pageContent, "", err
	}

	return pageContent, base64TypeScreenShot, nil
}

func SaveCrawledData(database *sqlx.DB, data schema.CrawlerTable) error {

	return nil
}

func SaveDownloadedVideos() {
	// Store the downloaded videos on the PostgreSQL server.
}

func SaveExtractedInternalLinksIntoCrawlerTable() {

}
