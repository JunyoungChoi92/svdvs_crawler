package scrapper

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/chromedp/cdproto/browser"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

func TestScrapWebSite(t *testing.T) {
	t.Run("test to scrap URL 네이버", func(t *testing.T) {
		targetUrl := "https://www.naver.com"
		// scrap logic 추가
		ScrapWebsiteHelper(t, targetUrl)
	})

	t.Run("test to scrap URL 속사정", func(t *testing.T) {
		targetUrl := "https://scsj-30.com/bbs/board.php?bo_table=kr"
		// scrap logic 추가
		ScrapWebsiteHelper(t, targetUrl)
	})

	t.Run("test to scrap URL 섹포츠", func(t *testing.T) {
		targetUrl := "https://sexports28.com/bbs/board.php?bo_table=kr"
		// scrap logic 추가
		ScrapWebsiteHelper(t, targetUrl)
	})

	t.Run("test to scrap URL av19", func(t *testing.T) {
		targetUrl := "https://av19.org/korea"
		// scrap logic 추가
		ScrapWebsiteHelper(t, targetUrl)
	})

	t.Run("test to scrap URL 밤허브", func(t *testing.T) {
		targetUrl := "https://www.bamhub14.me/index.php?so=korea&v=s"
		// scrap logic 추가
		ScrapWebsiteHelper(t, targetUrl)
	})

	t.Run("test to scrap URL 칵테일", func(t *testing.T) {
		targetUrl := "https://cocktail41.com/bbs/board.php?bo_table=kr"
		// scrap logic 추가
		ScrapWebsiteHelper(t, targetUrl)
	})
}

func ScrapWebsiteHelper(t *testing.T, targetUrl string) {
	log.Println("target URL: ", targetUrl)
	ctx, cancel, err := InitializeBrowserContext()
	if err != nil {
		t.Errorf("Error occurred while initializing browser context. \n error log is here\n%s", err)
	}
	defer cancel()

	pageContent, err := ScrapHTML(ctx, targetUrl, 20*time.Second, false)
	if err != nil {
		t.Errorf("Error occurred while scraping %s. \n error log is here\n%s", targetUrl, err)
	}

	screenshot, err := CaptureScreenshot(ctx, targetUrl)
	if err != nil {
		t.Errorf("Error occurred while capturing screenshot %s. \n error log is here\n%s", targetUrl, err)
	}

	if err == nil {
		log.Println("success to scrap: ", targetUrl)
	}

	if pageContent != "" {
		log.Println("success to get the pageContent : ", targetUrl)
	}

	if screenshot != "" {
		log.Println("success to get the screenshot : ", targetUrl)
	}

	// assert.Empty(t, pageContent, fmt.Sprintf("HTML content is empty for %s", targetUrl))
	// assert.NotEmpty(t, screenshot, fmt.Sprintf("Getting pageContent is success for %s", targetUrl))
	// assert.Empty(t, screenshot, fmt.Sprintf("Screenshot is empty for %s", targetUrl))
	// assert.NotEmpty(t, pageContent, fmt.Sprintf("capturing Screenshot is success for %s", targetUrl))
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
		network.Enable(),
		browser.SetDownloadBehavior(browser.SetDownloadBehaviorBehaviorDeny),
		// To ensure that we're always fetching the latest version of resources,
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
