package link

import (
	"net/url"
	"regexp"
)

type UrlParser struct {
	urls        []string
	links       []string
	domains     []string
	formatLinks []string
}

func (urlParser *UrlParser) Urls() []string {
	return urlParser.urls
}

func (urlParser *UrlParser) Links() []string {
	return urlParser.links
}

func (urlParser *UrlParser) Domains() []string {
	return urlParser.domains
}

func (urlParser *UrlParser) FormatLinks() []string {
	return urlParser.formatLinks
}

func (urlParser *UrlParser) SetUrls(urls []string) {
	urlParser.urls = urls
}

func (urlParser *UrlParser) SetLinks(links []string) {
	urlParser.links = links
}

func (urlParser *UrlParser) SetDomains(domains []string) {
	urlParser.domains = domains
}

func (urlParser *UrlParser) SetFormatLinks(formatLinks []string) {
	urlParser.formatLinks = formatLinks
}

func (urlParser *UrlParser) CheckDuplicate() *UrlParser {
	urlParser.urls = RemoveDuplicates(urlParser.urls)

	return urlParser
}

func (urlParser *UrlParser) SeparateUrls(domain string) error {
	var links, domains []string
	domainParsed, err := url.Parse(domain)
	if err != nil {
		return err
	}
	wwwReg := regexp.MustCompile(`www.`)
	etcReg := regexp.MustCompile(`.json|.css|.js|favicon.ico|.png|.jepg|.webp|login|singin|singup|DMCA|privacy|policy|notice|register|password`)
	domainHost := wwwReg.ReplaceAllString(domainParsed.Host, ``)

	for _, uri := range urlParser.Urls() {
		if etcReg.MatchString(uri) {
			continue
		}
		urlParsed, err := url.Parse(uri)
		if err != nil {
			continue
		}
		urlHost := wwwReg.ReplaceAllString(urlParsed.Host, ``)

		if urlHost == domainHost {
			links = append(links, uri)
			continue
		}
		domains = append(domains, uri)
	}

	urlParser.SetLinks(links)
	urlParser.SetDomains(domains)

	return nil
}

func (urlParser *UrlParser) Formatting(format, parentUrl string) *UrlParser {
	var formatLinks []string
	// If format is empty, simply copy links to formatLinks
	if format == "" {
		if parentUrl == "" {
			formatLinks = urlParser.links
		}
		// Add logic here if there's any behavior for when format == "" and parentURL != ""
	} else {
		formatReg := formatToRegexp(format)
		for _, link := range urlParser.links {
			if formatReg.MatchString(link) {
				formatLinks = append(formatLinks, link)
			}
		}
	}

	urlParser.formatLinks = formatLinks
	return urlParser
}

func formatToRegexp(format string) *regexp.Regexp {
	reg := regexp.MustCompile(`{d}`)
	regexpQuestionMark := regexp.MustCompile(`\?`)

	formatToRegexpString := regexpQuestionMark.ReplaceAllString(format, `\?`)
	formatToRegexpString = reg.ReplaceAllString(formatToRegexpString, `(.*?)`)
	return regexp.MustCompile(formatToRegexpString)
}

func RemoveDuplicates(urls []string) []string {
	dupMap := make(map[string]bool, len(urls))
	var uniqueSlice []string

	for _, item := range urls {
		if !dupMap[item] {
			dupMap[item] = true
			uniqueSlice = append(uniqueSlice, item)
		}
	}
	return uniqueSlice
}
