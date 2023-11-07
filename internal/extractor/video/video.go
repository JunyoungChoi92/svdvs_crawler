package video

import (
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
	Source       string
	Header       string
	UnEscapeHtml string
	ParentUrl    string
}

func (videoSource *VideoSource) GetTags(tag string, attribute string, method string) ([]string, error) {
	unEscapeHtml := videoSource.UnEscapeHtml
	var tags []string

	// Prepare the HTML based on the specified method
	if method != cdn {
		spaceReg := regexp.MustCompile(`\s`)
		unEscapeHtml = spaceReg.ReplaceAllString(unEscapeHtml, ``)
	}
	// Parse the HTML document
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(unEscapeHtml))
	if err != nil {
		return nil, err
	}
	// Use GoQuery to select all elements with the given tag name
	doc.Find(tag).Each(func(i int, s *goquery.Selection) {
		// Extract the attribute value from the selected tag
		attr, exists := s.Attr(attribute)
		if exists {
			tags = append(tags, attr)
		}
	})

	return tags, nil
}

func (videoSource *VideoSource) FindVideos() error {
	playerReg := regexp.MustCompile("player_html5|fluidPlayer|sendvid|aliplayer|youtubecomsplayer|DPlayer|video-player")
	res := playerReg.FindString(videoSource.UnEscapeHtml)
	switch res {
	case htmlPlayer:
		return videoSource.findHtmlPlayer()
	case fluidPlayer:
		return videoSource.findFluidPlayer()
	case sendvid:
		return videoSource.findSendVid()
	case aliPlayer:
		return videoSource.findAliPlayer()
	case youtube:
		return videoSource.findYoutube()
	case dplayer:
		return videoSource.findDPlayer()
	case cdn:
		return videoSource.findCDNPlayer()
	default:
		return videoSource.findCDNPlayer()
	}
}

func (videoSource *VideoSource) IsVideoSourceLink() bool {
	if strings.Compare(videoSource.Source, htmlPlayer) == 0 || strings.Compare(videoSource.Source, youtube) == 0 {
		return true
	}

	return false
}
