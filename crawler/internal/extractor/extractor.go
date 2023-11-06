package extractor

import (
	"log"

	"github.com/jmoiron/sqlx"

	"github.com/dev-zipida-com/spo-vdvs-crawler/internal/db/schema"
	"github.com/dev-zipida-com/spo-vdvs-crawler/internal/extractor/parser"
)

func ParseHtml(html string) *parser.HtmlParser {
	return parser.NewParseHtml(html)
}

func ExtractSubLinks(html, domain, parentUrl string) ([]string, error) {
	// Extract all necessary arguments and headers from the HTML code that are required by the yt-dlp module to download videos.
	// For example, If the video comes from a CDN, ensure that any additional header values are also extracted.
	htmlParser := ParseHtml(html)
	htmlParser.SetFormat(domain)
	htmlParser.ExtractAndCombineURLs(domain).SeparateUrls(domain)
	htmlParser.Formatting(htmlParser.Format(), parentUrl)

	links := htmlParser.FormatLinks()

	return links, nil
}

func CreateVideoSource(html, domain string) (*schema.VideoDownloader, error) {
	htmlParser := ParseHtml(html)
	htmlParser.SetParentUrl(domain)

	videoSource, err := htmlParser.FindVideos()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	videoDownloader := &schema.VideoDownloader{
		VideoDownloadSource: &videoSource.Source,
		VideoDownloadHeader: &videoSource.Header,
	}

	return videoDownloader, nil
}

func SaveVideoSource(db *sqlx.DB, videoDownloader *schema.VideoDownloader) error {
	err := videoDownloader.Create(db)
	if err != nil {
		return err
	}

	return nil
}

func CreateExtractedLinks(db *sqlx.DB, crawlerID string) error {
	crawler, err := schema.ReadCrawler(db, crawlerID)
	if err != nil {
		return err
	}
	domain, err := parser.ExtractDomain(crawler.CurrentUrl)
	if err != nil {
		return err
	}
	subLinks, err := ExtractSubLinks(crawler.HTML, domain, crawler.ParentUrl)
	if err != nil {
		return err
	}

	var crawlers []*schema.CrawlerTable
	for _, link := range subLinks {
		crawlerStruct := &schema.CrawlerTable{
			Name:         crawler.Name,
			ParentUrl:    crawler.CurrentUrl,
			CurrentUrl:   link,
			SeedDomainID: crawler.SeedDomainID,
		}
		crawlers = append(crawlers, crawlerStruct)
	}

	err = schema.BatchInsertCrawlers(db, crawlers)
	if err != nil {
		return err
	}

	return nil
}
