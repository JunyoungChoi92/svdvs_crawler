package video

import (
	"errors"
	"html"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"github.com/dev-zipida-com/spo-vdvs-crawler/internal/extractor/link"
)

func (videoSource *VideoSource) findHtmlPlayer() error {
	videoSource.Source = htmlPlayer
	return nil
}

func (videoSource *VideoSource) findYoutube() error {
	videoSource.Source = youtube
	return nil
}

func (videoSource *VideoSource) findFluidPlayer() error {
	// Assume you've already parsed your HTML content into a *goquery.Document named doc
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(videoSource.UnEscapeHtml)) // Replace videoSource.HTML with your method to get the HTML content
	if err != nil {
		return err
	}

	doc.Find("iframe").Each(func(index int, item *goquery.Selection) {
		// Get src attribute of iframe
		src, exists := item.Attr("src")
		if !exists || !strings.HasPrefix(src, "http") {
			return
		}
		videoSource.Source = src
	})

	if videoSource.Source == "" {
		return errors.New("not found video source")
	}

	return nil
}

func (videoSource *VideoSource) findSendVid() error {
	sources, err := videoSource.GetTags("source", "src", sendvid)
	if err != nil {
		return err
	}
	for _, source := range sources {
		if len(source) == 0 || !strings.HasPrefix(source, "http") {
			continue
		}
		videoSource.Source = source
		return nil
	}
	return errors.New("not found video source")
}

func (videoSource *VideoSource) findAliPlayer() error {
	iframes, err := videoSource.GetTags("iframe", "src", aliPlayer)
	if err != nil {
		return err
	}
	for _, iframe := range iframes {
		if len(iframe) == 0 || !strings.HasPrefix(iframe, "http") {
			continue
		}

		videoSourceParsed, _ := url.Parse(iframe)
		videoSource.Source = videoSourceParsed.Query().Get("url")
		return nil
	}

	return errors.New("not found video source")
}

func (videoSource *VideoSource) findDPlayer() error {
	reg := regexp.MustCompile(`((([A-Za-z]{3,9}:(?://)?)(?:[-;:&=+$,\w]+@)?[A-Za-z0-9.-]+|(?:www.|[-;:&=+$,\w]+@)[A-Za-z0-9.-]+)((?:/[+~%/.\w-_]*)?\??[-+=&;%@.\w_]*#?[.!/\\w]*)?)`)
	res := reg.FindAll([]byte(videoSource.UnEscapeHtml), -1)

	for _, v := range res {
		videoSource.Source = string(v)
		if strings.Contains(videoSource.Source, `.m3u8`) {
			_, err := url.Parse(videoSource.Source)
			if err != nil {
				continue
			}
			break
		}
	}

	return nil
}

// TODO: Test Case FAIL 해결 필요
func (videoSource *VideoSource) findCDNPlayer() error {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(videoSource.UnEscapeHtml))
	if err != nil {
		log.Println("Error creating document from reader:", err)
		return err
	}

	doc.Find("iframe").EachWithBreak(func(index int, item *goquery.Selection) bool {
		src, exists := item.Attr("src")
		if !exists || !strings.HasPrefix(src, "http") {
			return true // Continue to the next iframe
		}

		videoSource.Source = html.UnescapeString(src)
		req, err := http.NewRequest("GET", videoSource.Source, nil)
		if err != nil {
			log.Println("http NewRequest error:", err)
			return false // Break the loop
		}
		req.Header.Add("referer", videoSource.ParentUrl)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Println("http client error:", err)
			return false // Break the loop
		}

		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				log.Println(err)
			}
		}(resp.Body)

		innerDoc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			log.Println("Error creating inner document:", err)
			return false // Break the loop
		}

		fileURL := innerDoc.Find("script").Text()
		re := regexp.MustCompile(`file:\s*"(.*?)"`)
		m3u8Target := re.FindStringSubmatch(fileURL)
		if len(m3u8Target) > 1 {
			domain, err := link.ExtractDomain(videoSource.Source)
			if err != nil {
				return false
			}
			videoSource.Header = domain
			videoSource.Source = m3u8Target[1]
			return false // Found target, break the loop
		} else {
			videoSource.Source = ""
		}

		return true
	})

	if videoSource.Source == "" {
		// Find the source tag and extract the src attribute
		doc.Find("source").Each(func(index int, item *goquery.Selection) {
			src, _ := item.Attr("src")
			if !strings.HasPrefix(src, "http") {
				return
			}
			videoSource.Source = src
		})
	}

	if videoSource.Source == "" {
		return errors.New("not found video source")
	}

	return nil
}
