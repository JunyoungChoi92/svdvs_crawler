package parser

import (
	"fmt"
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

func (htmlParser *HtmlParser) ExtractAndCombineURLs(domain string) *HtmlParser {
	htmlParser.extractAndNormalizeURLs(domain)
	links := htmlParser.extractURLsFromContent(htmlParser.UnEscapeHtml())

	htmlParser.SetUrls(CombineSlice(htmlParser.Urls(), links))

	return htmlParser
}

func (htmlParser *HtmlParser) GetTags(tag string, method string) []string {
	unEscapeHtml := htmlParser.UnEscapeHtml()

	var tags []string
	if method == cdn {
		tagReg := regexp.MustCompile(fmt.Sprintf(`<%s(.*?)>`, tag))
		tags = tagReg.FindAllString(unEscapeHtml, -1)
	} else {
		spaceReg := regexp.MustCompile(`\s`)
		tagReg := regexp.MustCompile(fmt.Sprintf(`<%s(.*?)>`, tag))

		unEscapeHtml = spaceReg.ReplaceAllString(unEscapeHtml, ``)
		tags = tagReg.FindAllString(unEscapeHtml, -1)
	}

	return tags
}

func (htmlParser *HtmlParser) extractAndNormalizeURLs(domain string) *HtmlParser {
	regexpPath := regexp.MustCompile(`href="(.*?)"`)
	regexpAttribute := regexp.MustCompile(`href=|"|\s`)
	urlList := regexpPath.FindAllString(htmlParser.UnEscapeHtml(), -1)

	for idx, uri := range urlList {
		uri = regexpAttribute.ReplaceAllString(uri, "")
		uri = checkFullLink(domain, uri)

		urlList[idx] = uri
	}
	htmlParser.SetUrls(urlList)

	return htmlParser
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

func checkFullLink(domain string, uri string) string {
	urlParsed, err := url.Parse(uri)
	if err != nil {
		log.Println(err)
		return ""
	}
	domainParsed, err := url.Parse(domain)
	if err != nil {
		log.Println(err)
		return ""
	}

	if len(urlParsed.Host) != 0 {
		return urlParsed.String()
	}

	domainParsed.Path = urlParsed.Path
	domainParsed.RawQuery = urlParsed.RawQuery
	return domainParsed.String()
}
