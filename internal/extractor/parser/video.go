package parser

import (
	"errors"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var (
	htmlPlayer  = "player_html5"
	fluidPlayer = "fluidPlayer"
	sendvid     = "sendvid"
	aliPlayer   = "aliplayer"
	youtube     = "youtubecomsplayer"
	dplayer     = "DPlayer"
	cdn         = "video-player"
)

type VideoSource struct {
	Source string
	Header string
}

func (htmlParser *HtmlParser) FindVideos() (*VideoSource, error) {
	playerReg := regexp.MustCompile("player_html5|fluidPlayer|sendvid|aliplayer|youtubecomsplayer|DPlayer|video-player")
	res := playerReg.FindString(htmlParser.UnEscapeHtml())
	switch res {
	case htmlPlayer:
		return htmlParser.findHtmlPlayer()
	case fluidPlayer:
		return htmlParser.findFluidPlayer()
	case sendvid:
		return htmlParser.findSendVid()
	case aliPlayer:
		return htmlParser.findAliPlayer()
	case youtube:
		return htmlParser.findYoutube()
	case dplayer:
		return htmlParser.findDPlayer()
	case cdn:
		return htmlParser.findCDNPlayer()
	default:
		return htmlParser.findCDNPlayer()
	}
}

func (htmlParser *HtmlParser) IsVideoSourceLink(videoSource string) bool {
	if strings.Compare(videoSource, htmlPlayer) == 0 || strings.Compare(videoSource, youtube) == 0 {
		return true
	}

	return false
}

func (htmlParser *HtmlParser) findHtmlPlayer() (*VideoSource, error) {
	return &VideoSource{Source: htmlPlayer}, nil
}

func (htmlParser *HtmlParser) findYoutube() (*VideoSource, error) {
	return &VideoSource{Source: youtube}, nil
}

func (htmlParser *HtmlParser) findFluidPlayer() (*VideoSource, error) {
	// Assume you've already parsed your HTML content into a *goquery.Document named doc
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlParser.UnEscapeHtml())) // Replace htmlParser.HTML with your method to get the HTML content
	if err != nil {
		return nil, err
	}

	var videoSource string
	doc.Find("iframe").Each(func(index int, item *goquery.Selection) {
		// Get src attribute of iframe
		src, exists := item.Attr("src")
		if !exists || !strings.HasPrefix(src, "http") {
			return
		}
		videoSource = src
	})

	if videoSource == "" {
		return nil, errors.New("not found video source")
	}

	return &VideoSource{Source: videoSource}, nil
}

func (htmlParser *HtmlParser) findSendVid() (*VideoSource, error) {
	sources := htmlParser.GetTags("source", sendvid)
	for _, source := range sources {
		srcReg := regexp.MustCompile(fmt.Sprintf(`%s="(.*?)"`, "src"))
		srcReplaceReg := regexp.MustCompile(`src(=?)|"`)

		srcAttr := srcReg.FindString(source)
		if len(srcAttr) == 0 || !strings.HasPrefix(srcAttr, "http") {
			continue
		}

		source := srcReplaceReg.ReplaceAllString(srcAttr, "")
		return &VideoSource{Source: source}, nil
	}
	return nil, errors.New("not found video source")
}

func (htmlParser *HtmlParser) findAliPlayer() (*VideoSource, error) {
	iframes := htmlParser.GetTags("iframe", aliPlayer)
	for _, iframe := range iframes {
		srcReg := regexp.MustCompile(fmt.Sprintf(`%s="(.*?)sewobofang(.*?)"`, "src"))
		srcReplaceReg := regexp.MustCompile(`src(=?)|"`)

		srcAttr := srcReg.FindString(iframe)
		if len(srcAttr) == 0 {
			continue
		}

		source := srcReplaceReg.ReplaceAllString(srcAttr, "")
		videoSourceParsed, _ := url.Parse(source)
		return &VideoSource{Source: videoSourceParsed.Query().Get("url")}, nil
	}

	return nil, errors.New("not found video source")
}

func (htmlParser *HtmlParser) findDPlayer() (*VideoSource, error) {
	videoSource := &VideoSource{}
	reg := regexp.MustCompile(`((([A-Za-z]{3,9}:(?://)?)(?:[-;:&=+$,\w]+@)?[A-Za-z0-9.-]+|(?:www.|[-;:&=+$,\w]+@)[A-Za-z0-9.-]+)((?:/[+~%/.\w-_]*)?\??[-+=&;%@.\w_]*#?[.!/\\w]*)?)`)
	res := reg.FindAll([]byte(htmlParser.UnEscapeHtml()), -1)

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

	return videoSource, nil
}

// TODO: Test Case FAIL 해결 필요
func (htmlParser *HtmlParser) findCDNPlayer() (*VideoSource, error) {
	videoSource := &VideoSource{}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlParser.UnEscapeHtml()))
	if err != nil {
		log.Println("Error creating document from reader:", err)
		return nil, err
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
		req.Header.Add("referer", htmlParser.ParentUrl())

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
			domain, err := ExtractDomain(videoSource.Source)
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
		return nil, errors.New("not found video source")
	}

	return videoSource, nil
}
