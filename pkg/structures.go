package pkg

import (
	"context"
	"errors"
	"time"
)

type Command interface {
	Execute(crawler *Crawler) error
}

type Crawler struct {
	BrowserContext context.Context
	Commands       []Command

	Id             string
	HTML           string
	ScreenshotPath string
	TargetURL      string
	CurrentURL     string
	subLinks       []string
}

func (crawler *Crawler) AddCommand(command ...Command) {
	crawler.Commands = append(crawler.Commands, command...)
}

func (crawler *Crawler) Run() error {
	for _, command := range crawler.Commands {
		if err := command.Execute(crawler); err != nil {
			return err
		}
	}
	return nil
}

type GetRedirectedUrlCommand struct {
	TargetURL string
	TimeOut   time.Duration
}

func (getRedirectedUrlCommand *GetRedirectedUrlCommand) Execute(crawler *Crawler) error {
	RedirectedUrl, err := GetRedirectedUrl(crawler.BrowserContext, getRedirectedUrlCommand.TargetURL, getRedirectedUrlCommand.TimeOut)
	if err != nil {
		return err
	}

	crawler.CurrentURL = RedirectedUrl
	return nil
}

type CheckUrlConnectivityCommand struct {
	TargetURL string
}

func (checkUrlConnectivityCommand *CheckUrlConnectivityCommand) Execute(crawler *Crawler) error {
	url := checkUrlConnectivityCommand.TargetURL
	connectability, err := CheckURLConnectivity(url)
	if err != nil {
		return err
	}

	// if connectability is false, just end the program.
	if !connectability {
		return errors.New("connectivity to the URL failed")
	}

	return nil
}

type CheckURLValidationCommand struct {
	TargetURL  string
	Indicators []string // hcaptchaIndicator, recaptchaIndicator
}

func (checkURLValidationCommand *CheckURLValidationCommand) Execute(crawler *Crawler) error {
	err := CheckURLValidation(crawler.BrowserContext, checkURLValidationCommand.TargetURL, checkURLValidationCommand.Indicators)
	if err != nil {
		return err
	}

	return nil
}

type ScrapWebpageCommand struct {
	TargetURL          string
	IsBlockedByCaptcha bool
	IframeTimeOut      time.Duration
}

func (scrapWebpageCommand *ScrapWebpageCommand) Execute(crawler *Crawler) error {
	pageContent, base64TypeScreenShot, err := ScrapWebpage(crawler.BrowserContext, scrapWebpageCommand.TargetURL, scrapWebpageCommand.IframeTimeOut, scrapWebpageCommand.IsBlockedByCaptcha)
	if err != nil {
		return err
	}

	crawler.HTML = pageContent
	crawler.ScreenshotPath = base64TypeScreenShot

	return nil
}
