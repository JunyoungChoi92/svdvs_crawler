package extractor

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// html 데이터를 db에서 호출해야하지만 테스트 케이스에서는 file read 기능으로 대체함

func TestExtractSubLinks(t *testing.T) {
	t.Run("test link find", func(t *testing.T) {

		urls := []string{
			"https://av19.org/korea", "https://av19.org/koreayadong", "https://av19.org/korea/12407", "https://av19.org/korea/12406", "https://av19.org/korea/12405",
			"https://av19.org/korea/12404", "https://av19.org/korea/12403", "https://av19.org/korea/12402", "https://av19.org/korea/12401", "https://av19.org/korea/12400",
			"https://av19.org/korea/12399", "https://av19.org/korea/12398", "https://av19.org/korea/12397", "https://av19.org/korea/12396", "https://av19.org/korea/12395",
			"https://av19.org/korea/12394", "https://av19.org/korea/12393", "https://av19.org/korea/12392", "https://av19.org/korea/12391", "https://av19.org/korea/12390",
			"https://av19.org/korea/12389", "https://av19.org/korea/12388", "https://av19.org/korea/12387", "https://av19.org/korea/12386", "https://av19.org/korea/12385",
			"https://av19.org/korea/12321", "https://av19.org/korea?page=2", "https://av19.org/korea?page=3", "https://av19.org/korea?page=4", "https://av19.org/korea?page=5",
			"https://av19.org/korea?page=6", "https://av19.org/korea?page=7", "https://av19.org/korea?page=8", "https://av19.org/korea?page=9", "https://av19.org/korea?page=10",
			"https://av19.org/korea?page=11", "https://av19.org/korea?page=473",
		}
		htmlData, err := os.ReadFile(`./data/av19korea.txt`)
		if err != nil {
			t.Error("os.ReadFile error: ", err)
		}

		extractedURLs, err := ExtractSubLinks(string(htmlData), "https://av19.org/korea", "https://av19.org")
		if err != nil {
			t.Error("ExtractSubLinks Error: ", err)
		}
		assert.Equal(t, urls, extractedURLs, "extractURLs is not equal 'urls'")
	})

	t.Run("test link find", func(t *testing.T) {

		urls := []string{
			"https://av19.org/korea", "https://av19.org/koreayadong", "https://av19.org/korea/12407", "https://av19.org/korea/12406", "https://av19.org/korea/12405",
			"https://av19.org/korea/12404", "https://av19.org/korea/12403", "https://av19.org/korea/12402", "https://av19.org/korea/12401", "https://av19.org/korea/12400",
			"https://av19.org/korea/12399", "https://av19.org/korea/12398", "https://av19.org/korea/12397", "https://av19.org/korea/12396", "https://av19.org/korea/12395",
			"https://av19.org/korea/12394", "https://av19.org/korea/12393", "https://av19.org/korea/12392", "https://av19.org/korea/12391", "https://av19.org/korea/12390",
			"https://av19.org/korea/12389", "https://av19.org/korea/12388", "https://av19.org/korea/12387", "https://av19.org/korea/12386", "https://av19.org/korea/12385",
			"https://av19.org/korea/12321", "https://av19.org/korea?page=2", "https://av19.org/korea?page=3", "https://av19.org/korea?page=4", "https://av19.org/korea?page=5",
			"https://av19.org/korea?page=6", "https://av19.org/korea?page=7", "https://av19.org/korea?page=8", "https://av19.org/korea?page=9", "https://av19.org/korea?page=10",
			"https://av19.org/korea?page=11", "https://av19.org/korea?page=473",
		}
		htmlData, err := os.ReadFile(`./data/av19korea.txt`)
		if err != nil {
			t.Error("os.ReadFile error: ", err)
		}

		extractedURLs, err := ExtractSubLinks(string(htmlData), "https://av19.org/korea", "https://av19.org")
		if err != nil {
			t.Error("ExtractSubLinks Error: ", err)
		}
		assert.Equal(t, urls, extractedURLs, "extractURLs is not equal 'urls'")
	})
}

func TestExtractVideoSource(t *testing.T) {

	t.Run("av19 video tag", func(t *testing.T) {
		referer := "https://david.cdnbuzz.buzz"
		videosource := "https://124fdsf6dsf.worldcup2022.icu/cupcup8/1111/98호. 짭수현 고물상 20191110.mp4/index.js"
		domain := "https://av19.org"

		htmlData, err := os.ReadFile(`./data/av19_video3.txt`)
		if err != nil {
			t.Error(err)
		}

		extractedVideoSource, err := CreateVideoSource(string(htmlData), domain)
		if err != nil {
			t.Error(err)
		}

		assert.Equal(t, referer, *extractedVideoSource.VideoDownloadHeader, "Header should be 'https://david.cdnbuzz.buzz/' ")
		assert.Equal(t, videosource, *extractedVideoSource.VideoDownloadSource, "VideoSource is not equal")
	})

	t.Run("yadong.best video tag", func(t *testing.T) {
		referer := "https://david.cdnbuzz.buzz"
		videosource := "https://124fdsf6dsf.worldcup2022.icu/cupcup8/n3/0422/BJ 박민정 아프리카 제대로 꼭노한다 230422.mp4/index.js"
		domain := "https://yadong.best"

		htmlData, err := os.ReadFile("./data/yadong_video.txt")
		if err != nil {
			t.Error(err)
		}

		extractedVideoSource, err := CreateVideoSource(string(htmlData), domain)
		if err != nil {
			t.Error(err)
		}

		assert.Equal(t, referer, *extractedVideoSource.VideoDownloadHeader, "Header should be 'https://david.cdnbuzz.buzz/' ")
		assert.Equal(t, videosource, *extractedVideoSource.VideoDownloadSource, "VideoSource is not equal")
	})

	t.Run("yadong.best video tag", func(t *testing.T) {
		referer := ""
		videosource := "https://video-u1-cdn.cccb03.com/20230422/CEDneqvO/index.m3u8?sign=54c5a5d6133838b6bdf9f5503fb4745bebf7dfcb5c2cd7a722ee0d8d27aaef3f82b2bc627af0c5f01e262c64d4334010"
		domain := "https://kr13.mysstv.com"

		htmlData, err := os.ReadFile("./data/kr13.txt")
		if err != nil {
			t.Error(err)
		}

		extractedVideoSource, err := CreateVideoSource(string(htmlData), domain)
		if err != nil {
			t.Error(err)
		}

		assert.Equal(t, referer, *extractedVideoSource.VideoDownloadHeader, "Header should be 'https://david.cdnbuzz.buzz/' ")
		assert.Equal(t, videosource, *extractedVideoSource.VideoDownloadSource, "VideoSource is not equal")
	})
}
