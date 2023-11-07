package pkg

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/dev-zipida-com/spo-vdvs-crawler/internal/db/schema"
	"github.com/dev-zipida-com/spo-vdvs-crawler/internal/extractor"
	"github.com/dev-zipida-com/spo-vdvs-crawler/internal/extractor/link"
)

type Command interface {
	Execute(crawler *Crawler) error
}

type Crawler struct {
	BrowserContext context.Context
	Commands       []Command

	Id             string
	Name           string
	HTML           string
	ScreenshotPath string
	TargetURL      string
	CurrentURL     string
	VideoSource    string
	VideoHeader    string
	SubLinks       []string
	SeedDomainID   string
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
	TimeOut time.Duration
}

func (getRedirectedUrlCommand *GetRedirectedUrlCommand) Execute(crawler *Crawler) error {
	RedirectedUrl, err := GetRedirectedUrl(crawler.BrowserContext, crawler.TargetURL, getRedirectedUrlCommand.TimeOut)
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
	Indicators []string // hcaptchaIndicator, recaptchaIndicator
}

func (checkURLValidationCommand *CheckURLValidationCommand) Execute(crawler *Crawler) error {
	err := CheckURLValidation(crawler.BrowserContext, crawler.TargetURL, checkURLValidationCommand.Indicators)
	if err != nil {
		return err
	}

	return nil
}

type ScrapWebpageCommand struct {
	IsBlockedByCaptcha bool
	IframeTimeOut      time.Duration
}

func (scrapWebpageCommand *ScrapWebpageCommand) Execute(crawler *Crawler) error {
	pageContent, base64TypeScreenShot, err := ScrapWebpage(crawler.BrowserContext, crawler.TargetURL, scrapWebpageCommand.IframeTimeOut, scrapWebpageCommand.IsBlockedByCaptcha)
	if err != nil {
		return err
	}

	crawler.HTML = pageContent
	crawler.ScreenshotPath = base64TypeScreenShot

	return nil
}

type ExtractSubLinkCommand struct {
	Html string
}

func (extractSubLinkCommand *ExtractSubLinkCommand) Execute(crawler *Crawler) error {
	links, err := extractor.ExtractSubLinks(crawler.HTML, crawler.TargetURL, crawler.CurrentURL)
	if err != nil {
		return err
	}
	crawler.SubLinks = links
	log.Println("SubLinks: ", links)

	return nil
}

type ExtractVideoSourceCommand struct {
}

func (extractVideoSourceCommand *ExtractVideoSourceCommand) Execute(crawler *Crawler) error {
	domain, err := link.ExtractDomain(crawler.CurrentURL)
	if err != nil {
		return err
	}
	videoDownloader, err := extractor.CreateVideoSource(crawler.HTML, domain)
	if err != nil {
		log.Println(err)
		return nil
	}

	crawler.VideoSource = *videoDownloader.VideoDownloadSource
	crawler.VideoHeader = *videoDownloader.VideoDownloadHeader
	log.Println("VideoSource: ", *videoDownloader.VideoDownloadSource)
	log.Println("VideoHeader: ", *videoDownloader.VideoDownloadHeader)

	return nil
}

type SaveVideoSourceCommand struct {
	conn *sqlx.DB
}

func (saveVideoSourceCommand *SaveVideoSourceCommand) Execute(crawler *Crawler) error {
	videoDownloader := &schema.VideoDownloader{
		VideoDownloadSource: &crawler.VideoSource,
		VideoDownloadHeader: &crawler.VideoHeader,
	}

	err := extractor.SaveVideoSource(saveVideoSourceCommand.conn, videoDownloader)
	if err != nil {
		return err
	}

	return nil
}

type SaveExtractLinksCommand struct {
	Conn *sqlx.DB
}

func (saveExtractLinksCommand *SaveExtractLinksCommand) Execute(crawler *Crawler) error {

	var crawlers []*schema.CrawlerTable
	for _, subLink := range crawler.SubLinks {
		crawlerStruct := &schema.CrawlerTable{
			Name:         crawler.Name,
			ParentUrl:    crawler.CurrentURL,
			CurrentUrl:   subLink,
			SeedDomainID: crawler.SeedDomainID,
		}
		crawlers = append(crawlers, crawlerStruct)
	}

	err := schema.BatchInsertCrawlers(saveExtractLinksCommand.Conn, crawlers)
	if err != nil {
		return err
	}

	return nil
}
