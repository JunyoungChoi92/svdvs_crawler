package pkg

import (
	"context"
	"log"
	"time"
)

func Crawling(targetUrl string) (*Crawler, error) {
	crawler := &Crawler{Id: "1"}

	// Set the targetURL
	crawler.TargetURL = targetUrl

	// Set the captcha indicators
	indicators := []string{`script src="https://js.hcaptcha.com/1/api.js"`, `script src="https://www.google.com/recaptcha/api.js"`}

	// Initialize the chromedp's BrowserContext
	ctx, cancel, err := InitializeBrowserContext()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer cancel()

	// the context for whole process will be quitted in 120 seconds.
	tctx, tcancel := context.WithTimeout(ctx, 120*time.Second)
	defer tcancel()

	// Set the BrowserContext and CancelFunc
	crawler.BrowserContext = tctx

	crawler.AddCommand(
		//&CheckUrlConnectivityCommand{
		//	TargetURL: crawler.TargetURL,
		//},
		&GetRedirectedUrlCommand{
			TimeOut: 15 * time.Second,
		},
		&CheckURLValidationCommand{
			Indicators: indicators,
		},
		&ScrapWebpageCommand{
			IsBlockedByCaptcha: false,
			IframeTimeOut:      15 * time.Second,
		},
		&ExtractSubLinkCommand{},
		&ExtractVideoSourceCommand{},
	)
	// insert extracted video download links into a new row in the video downloader table
	// Run the commands

	if err := crawler.Run(); err != nil {
		log.Println(err)
		return nil, err
	}

	return crawler, nil
}
