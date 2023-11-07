package extractor

import (
	"log"
	"net/url"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/dev-zipida-com/spo-vdvs-crawler/internal/extractor/link"
)

func TestRegexpVideo(t *testing.T) {

	t.Run("test find video tag", func(t *testing.T) {
		expectVideoSource := "https://124fdsf6dsf.worldcup2022.icu/cupcup8/n3/0422/BJ 카이 친구의 남친이랑 해봤다는 게스트 221202.mp4/index.js"
		htmlVideo, err := os.ReadFile("./data/av19_video1.txt")
		if err != nil {
			t.Error(err)
		}
		htmlParser := link.NewParseHtml(string(htmlVideo))
		htmlParser.SetParentUrl("https://av19.org")

		actualVideoSource := HtmlParserToVideoSource(htmlParser)

		err = actualVideoSource.FindVideos()
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, expectVideoSource, actualVideoSource.Source, "error")
	})

	t.Run("test find fluidplayer", func(t *testing.T) {
		expectVideoSource := "https://jp.thisiscdn.life/cupcup8//miss/90474789475269.mp4/index.js"
		htmlVideo, err := os.ReadFile("./data/av19_video2.txt")
		if err != nil {
			t.Error(err)
		}
		htmlParser := link.NewParseHtml(string(htmlVideo))
		htmlParser.SetParentUrl("https://av19.org")

		actualVideoSource := HtmlParserToVideoSource(htmlParser)

		err = actualVideoSource.FindVideos()
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, expectVideoSource, actualVideoSource.Source)
	})

	t.Run("test find video tag", func(t *testing.T) {
		expectVideoSource := "https://124fdsf6dsf.worldcup2022.icu/cupcup8/1111/98호. 짭수현 고물상 20191110.mp4/index.js"
		htmlVideo, err := os.ReadFile("./data/av19_video3.txt")
		if err != nil {
			t.Error(err)
		}
		htmlParser := link.NewParseHtml(string(htmlVideo))
		htmlParser.SetParentUrl("https://av19.org")

		actualVideoSource := HtmlParserToVideoSource(htmlParser)

		err = actualVideoSource.FindVideos()
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, expectVideoSource, actualVideoSource.Source)
	})

	t.Run("test find video tag", func(t *testing.T) {
		expectVideoSource := "https://124fdsf6dsf.worldcup2022.icu/cupcup8/n3/0422/BJ 박민정 아프리카 제대로 꼭노한다 230422.mp4/index.js"
		htmlVideo, err := os.ReadFile("./data/yadong_video.txt")
		if err != nil {
			t.Error(err)
		}
		htmlParser := link.NewParseHtml(string(htmlVideo))
		htmlParser.SetParentUrl("https://yadong.best/")

		actualVideoSource := HtmlParserToVideoSource(htmlParser)

		err = actualVideoSource.FindVideos()
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, expectVideoSource, actualVideoSource.Source)
	})

	t.Run("test find DPlayer video tag", func(t *testing.T) {
		expectVideoSource := "https://124fdsf6dsf.worldcup2022.icu/cupcup8/n3/0427/세희 화보 영상.mp4/index.js"
		htmlVideo, err := os.ReadFile("./data/yadong.txt")
		if err != nil {
			t.Error(err)
		}
		htmlParser := link.NewParseHtml(string(htmlVideo))
		htmlParser.SetParentUrl("https://av19.org")

		actualVideoSource := HtmlParserToVideoSource(htmlParser)

		err = actualVideoSource.FindVideos()
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, expectVideoSource, actualVideoSource.Source)
	})

	t.Run("test find Plyr video tag", func(t *testing.T) {
		expectVideoSource := "https://video-cdn-u2.cccb03.com/20231102/2Y74XBJ6/index.m3u8?sign=c30141940577c80dda87e84c3c476afb89fdae5ab8db3543a2b337a1432b4e32ecc9b40c576fd38cee7dffe6cfcd9465"
		htmlVideo, err := os.ReadFile("./data/kr13_video.txt")
		if err != nil {
			t.Error(err)
		}
		htmlParser := link.NewParseHtml(string(htmlVideo))
		htmlParser.SetParentUrl("https://av19.org")

		actualVideoSource := HtmlParserToVideoSource(htmlParser)

		err = actualVideoSource.FindVideos()
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, expectVideoSource, actualVideoSource.Source)
	})

	// <source> Tag 에 동영상 소스 이름이 존재
	t.Run("https://avlove11.com", func(t *testing.T) {
		expectVideoSource := "https://cdn.sdfj923rjsdg23.com/2211/13/MDBK-269.mp4"
		htmlVideo, err := os.ReadFile("./data/avlove4_video.txt")
		if err != nil {
			t.Error(err)
		}
		htmlParser := link.NewParseHtml(string(htmlVideo))
		htmlParser.SetParentUrl("https://avlove11.com")

		actualVideoSource := HtmlParserToVideoSource(htmlParser)

		err = actualVideoSource.FindVideos()
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, expectVideoSource, actualVideoSource.Source)
	})

	// TODO: FAIL
	t.Run("https://koreansexvid05.com", func(t *testing.T) {
		expectVideoSource := "https://s3t3d2y8.afcdn.net/library/611708/c3b51547f926dd0b183905e2c8586be6ea60775b.mp4"
		htmlVideo, err := os.ReadFile("./data/koreansexvid05_video.html")
		if err != nil {
			t.Error(err)
		}
		htmlParser := link.NewParseHtml(string(htmlVideo))
		htmlParser.SetParentUrl("https://koreansexvid05.com/")

		actualVideoSource := HtmlParserToVideoSource(htmlParser)

		err = actualVideoSource.FindVideos()
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, expectVideoSource, actualVideoSource.Source)
	})
	// TODO: "FAIL" => 개발자 도구 실행시 접속 차단 당함; + iframe 주소로 http 요청해도 File 값을 전달받지 못하는 문제
	t.Run("https://yadongtube.net", func(t *testing.T) {
		expectVideoSource := "https://hellocdn1.net/stream/?pc=true&title=%5BBHG-046%5D&v=6148523063484d364c79397159585a77624746355953356a623230765a5339344d484d795a4773344e5755314e6d4d756148527462413d3d&img=https%3A%2F%2Fimg.hellocdn1.net%2Fjimg%2F2f7a3ac93c186711faa77a2f69047b32.jpg&s=5a33567964513d3d&h=6557466b6232356e644856695a5335755a58513d&m=h&t=0&g=s"
		htmlVideo, err := os.ReadFile("./data/yadongtube_video.txt")
		if err != nil {
			t.Error(err)
		}
		htmlParser := link.NewParseHtml(string(htmlVideo))
		htmlParser.SetParentUrl("https://yadongtube.net")

		actualVideoSource := HtmlParserToVideoSource(htmlParser)

		err = actualVideoSource.FindVideos()
		if err != nil {
			log.Println(err)
		}
		assert.Equal(t, expectVideoSource, actualVideoSource.Source)
	})

	t.Run("test parse url", func(t *testing.T) {
		u, err := url.Parse("https://david.cdnbuzz.buzz/i.php?poster=https://yadong.best/data/file/korea/5481f222f6c5498a413608f0c9e69adb_gQ1NjZlM_38346a33fe05fadac317861f75d4291b4ae54a56.jpg&amp;vvv=n3/0427/세희 화보 영상.mp4")
		if err != nil {
			t.Error(err)
		}

		domain := u.Scheme + "://" + u.Hostname()
		log.Println(domain)
	})

	t.Run("test parse url", func(t *testing.T) {
		u, err := url.Parse("https://koreansexvid05.com/wp-content/plugins/clean-tube-player/public/player-x.php?q=cG9zdF9pZD0xMzMyNiZ0eXBlPWlmcmFtZSZ0YWc9JTNDaWZyYW1lJTIwc3JjJTNEJTIyaHR0cHMlM0ElMkYlMkZ3d3cueHZpZGVvcy5jb20lMkZlbWJlZGZyYW1lJTJGNzQwMzM1MTklMjIlMjBmcmFtZWJvcmRlciUzRCUyMjAlMjIlMjB3aWR0aCUzRCUyMjUxMCUyMiUyMGhlaWdodCUzRCUyMjQwMCUyMiUyMHNjcm9sbGluZyUzRCUyMm5vJTIyJTIwYWxsb3dmdWxsc2NyZWVuJTNEJTIyYWxsb3dmdWxsc2NyZWVuJTIyJTNFJTNDJTJGaWZyYW1lJTNF")
		if err != nil {
			t.Error(err)
		}

		domain := u.Scheme + "://" + u.Hostname()
		log.Println(domain)
	})
}
