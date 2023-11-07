package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/dev-zipida-com/spo-vdvs-crawler/pkg"
)

func main() {
	// Check if a URL argument was provided
	if len(os.Args) != 3 {
		log.Println("Usage: go run crawler.go <crawler's ID> <url>")
		return
	}

	// Set the cralwer's ID (Status is Idle for now)
	crawler := &pkg.Crawler{Id: os.Args[1]}

	// Set the targetURL
	crawler.TargetURL = os.Args[2]

	// Set the captcha indicators
	indicators := []string{`script src="https://js.hcaptcha.com/1/api.js"`, `script src="https://www.google.com/recaptcha/api.js"`}

	configPath := "../config.yaml"
	// Connect to the database
	database := pkg.NewDatabase()

	err := database.Connect(configPath)
	if err != nil {
		log.Fatal(err)
	}

	defer database.Disconnect()

	// Initialize the chromedp's BrowserContext
	ctx, cancel, err := pkg.InitializeBrowserContext()
	if err != nil {
		log.Fatal(err)
	}
	defer cancel()

	// the context for whole process will be quitted in 120 seconds.
	tctx, tcancel := context.WithTimeout(ctx, 120*time.Second)
	defer tcancel()

	// Set the BrowserContext and CancelFunc
	crawler.BrowserContext = tctx

	crawler.AddCommand(
		&pkg.CheckUrlConnectivityCommand{
			TargetURL: crawler.TargetURL,
		},
		&pkg.GetRedirectedUrlCommand{
			TimeOut: 15 * time.Second,
		},
		&pkg.CheckURLValidationCommand{
			Indicators: indicators,
		},
		&pkg.ScrapWebpageCommand{
			IsBlockedByCaptcha: false,
			IframeTimeOut:      15 * time.Second,
		},
		&pkg.ExtractSubLinkCommand{},
		&pkg.ExtractVideoSourceCommand{},
	)

	if err := crawler.Run(); err != nil {
		log.Fatal(err)
	}
}
