package link

import (
	"html"
	"log"
	"net/url"
	"regexp"
)

type HtmlParser struct {
	ContentParser
	UrlParser

	html         string
	escapeHtml   string
	unescapeHtml string
	format       string
	parentUrl    string
}

func NewHtmlParser() *HtmlParser {
	return &HtmlParser{}
}

func NewParseHtml(html string) *HtmlParser {
	htmlParser := &HtmlParser{
		html: html,
	}
	htmlParser.htmlEscape()
	htmlParser.htmlUnescape()

	return htmlParser
}

func (htmlParser *HtmlParser) Html() string {
	return htmlParser.html
}

func (htmlParser *HtmlParser) EscapeHtml() string {
	return htmlParser.escapeHtml
}

func (htmlParser *HtmlParser) UnEscapeHtml() string {
	return htmlParser.unescapeHtml
}

func (htmlParser *HtmlParser) Format() string {
	return htmlParser.format
}

func (htmlParser *HtmlParser) ParentUrl() string {
	return htmlParser.parentUrl
}

func (htmlParser *HtmlParser) SetParentUrl(parentUrl string) {
	htmlParser.parentUrl = parentUrl
}

func (htmlParser *HtmlParser) SetHtml(html string) {
	htmlParser.html = html
}

func (htmlParser *HtmlParser) SetEscapeHtml(escapeHtml string) {
	htmlParser.escapeHtml = escapeHtml
}

func (htmlParser *HtmlParser) SetUnEscapeHtml(unEscapeHtml string) {
	htmlParser.unescapeHtml = unEscapeHtml
}

func (htmlParser *HtmlParser) SetFormat(format string) {
	htmlParser.format = format
}

func (htmlParser *HtmlParser) ExtractURLsFromContent() *HtmlParser {
	htmlParser.extractURLsFromContent(htmlParser.UnEscapeHtml())

	return htmlParser
}

func (htmlParser *HtmlParser) ExtractAndCombineURLs(domain string) (*HtmlParser, error) {
	err := htmlParser.extractAndNormalizeURLs(domain)
	if err != nil {
		return htmlParser, err
	}
	links := htmlParser.extractURLsFromContent(htmlParser.UnEscapeHtml())

	htmlParser.SetUrls(CombineSlice(htmlParser.Urls(), links))

	return htmlParser, nil
}

func (htmlParser *HtmlParser) extractAndNormalizeURLs(domain string) error {
	regexpPath := regexp.MustCompile(`href="(.*?)"`)
	regexpAttribute := regexp.MustCompile(`href=|"|\s`)
	urlList := regexpPath.FindAllString(htmlParser.UnEscapeHtml(), -1)

	for idx, uri := range urlList {
		uri = regexpAttribute.ReplaceAllString(uri, "")
		uri, err := checkFullLink(domain, uri)
		if err != nil {
			return err
		}

		urlList[idx] = uri
	}
	htmlParser.SetUrls(urlList)

	return nil
}

func (htmlParser *HtmlParser) htmlUnescape() {
	// Unescape URL path component
	originHtml, err := url.PathUnescape(htmlParser.html)
	if err != nil {
		htmlParser.SetUnEscapeHtml(htmlParser.Html())
		return
	}

	newlineReg := regexp.MustCompile(`<br>`)
	ampReg := regexp.MustCompile(`&amp;`)
	// Unescape HTML entities and other patterns
	originHtml = html.UnescapeString(originHtml)
	originHtml = newlineReg.ReplaceAllString(originHtml, `\n`)
	originHtml = ampReg.ReplaceAllString(originHtml, `&`)

	htmlParser.SetUnEscapeHtml(originHtml)
}

func (htmlParser *HtmlParser) htmlEscape() {
	var newlineReg = regexp.MustCompile(`\n`)

	escapedHtml := newlineReg.ReplaceAllString(
		url.QueryEscape(
			url.PathEscape(
				html.EscapeString(htmlParser.Html()),
			),
		), "<br",
	)
	htmlParser.SetEscapeHtml(escapedHtml)
}

func checkFullLink(domain string, uri string) (string, error) {
	urlParsed, err := url.Parse(uri)
	if err != nil {
		log.Println(err)
		return "", err
	}
	domainParsed, err := url.Parse(domain)
	if err != nil {
		log.Println(err)
		return "", err
	}

	if len(urlParsed.Host) != 0 {
		return urlParsed.String(), nil
	}

	domainParsed.Path = urlParsed.Path
	domainParsed.RawQuery = urlParsed.RawQuery
	return domainParsed.String(), nil
}
