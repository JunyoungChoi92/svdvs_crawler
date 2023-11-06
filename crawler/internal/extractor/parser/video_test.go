package parser

import (
	"net/url"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegexpVideo(t *testing.T) {

	t.Run("test find video tag", func(t *testing.T) {
		expectVideoSource := "https://124fdsf6dsf.worldcup2022.icu/cupcup8/n3/0422/BJ 카이 친구의 남친이랑 해봤다는 게스트 221202.mp4/index.js"
		htmlVideo, err := os.ReadFile("../data/av19_video1.txt")
		if err != nil {
			t.Error(err)
		}
		htmlParser := NewParseHtml(string(htmlVideo))
		htmlParser.SetParentUrl("https://av19.org")

		actualVideoSource, err := htmlParser.FindVideos()
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, expectVideoSource, actualVideoSource.Source, "error")
	})

	t.Run("test find fluidplayer", func(t *testing.T) {
		expectVideoSource := "https://jp.thisiscdn.life/cupcup8//miss/90474789475269.mp4/index.js"
		htmlVideo, err := os.ReadFile("../data/av19_video2.txt")
		if err != nil {
			t.Error(err)
		}
		htmlParser := NewParseHtml(string(htmlVideo))
		htmlParser.SetParentUrl("https://av19.org")

		actualVideoSource, err := htmlParser.FindVideos()
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, expectVideoSource, actualVideoSource.Source)
	})

	t.Run("test find video tag", func(t *testing.T) {
		expectVideoSource := "https://124fdsf6dsf.worldcup2022.icu/cupcup8/1111/98호. 짭수현 고물상 20191110.mp4/index.js"
		htmlVideo, err := os.ReadFile("../data/av19_video3.txt")
		if err != nil {
			t.Error(err)
		}
		htmlParser := NewParseHtml(string(htmlVideo))
		htmlParser.SetParentUrl("https://av19.org")

		actualVideoSource, err := htmlParser.FindVideos()
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, expectVideoSource, actualVideoSource.Source)
	})

	t.Run("test find video tag", func(t *testing.T) {
		expectVideoSource := "https://124fdsf6dsf.worldcup2022.icu/cupcup8/n3/0422/BJ 박민정 아프리카 제대로 꼭노한다 230422.mp4/index.js"
		htmlVideo, err := os.ReadFile("../data/yadong_video.txt")
		if err != nil {
			t.Error(err)
		}
		htmlParser := NewParseHtml(string(htmlVideo))
		htmlParser.SetParentUrl("https://yadong.best/")

		actualVideoSource, err := htmlParser.FindVideos()
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, expectVideoSource, actualVideoSource.Source)
	})

	t.Run("test find DPlayer video tag", func(t *testing.T) {
		expectVideoSource := "https://124fdsf6dsf.worldcup2022.icu/cupcup8/n3/0427/세희 화보 영상.mp4/index.js"
		htmlVideo, err := os.ReadFile("../data/yadong.txt")
		if err != nil {
			t.Error(err)
		}
		htmlParser := NewParseHtml(string(htmlVideo))
		htmlParser.SetParentUrl("https://av19.org")

		actualVideoSource, err := htmlParser.FindVideos()
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, expectVideoSource, actualVideoSource.Source)
	})

	t.Run("test find Plyr video tag", func(t *testing.T) {
		expectVideoSource := "https://video-u1-cdn.cccb03.com/20230422/CEDneqvO/index.m3u8?sign=54c5a5d6133838b6bdf9f5503fb4745bebf7dfcb5c2cd7a722ee0d8d27aaef3f82b2bc627af0c5f01e262c64d4334010"
		htmlVideo, err := os.ReadFile("../data/kr13.txt")
		if err != nil {
			t.Error(err)
		}
		htmlParser := NewParseHtml(string(htmlVideo))
		htmlParser.SetParentUrl("https://av19.org")

		actualVideoSource, err := htmlParser.FindVideos()
		if err != nil || actualVideoSource == nil {
			t.Error(err)
		}
		assert.Equal(t, expectVideoSource, actualVideoSource.Source)
	})

	// TODO: "FAIL" => 개발자 도구 실행시 접속 차단 당함; + iframe 주소로 http 요청해도 File 값을 전달받지 못하는 문제
	t.Run("test link find", func(t *testing.T) {
		expectVideoSource := "https://video-u1-cdn.cccb03.com/20230422/CEDneqvO/index.m3u8?sign=54c5a5d6133838b6bdf9f5503fb4745bebf7dfcb5c2cd7a722ee0d8d27aaef3f82b2bc627af0c5f01e262c64d4334010"
		htmlVideo, err := os.ReadFile("../data/yadongtube_video.txt")
		if err != nil {
			t.Error(err)
		}
		htmlParser := NewParseHtml(string(htmlVideo))
		htmlParser.SetParentUrl("https://yadongtube.net")

		actualVideoSource, err := htmlParser.FindVideos()
		if err != nil || actualVideoSource == nil {
			t.Error(err)
		}
		assert.Equal(t, expectVideoSource, actualVideoSource.Source)
	})

	// TODO: 사이트 서버 에러로 테스트 실패 + <source> Tag 에 동영상 소스 이름이 존재함
	t.Run("test find video tag", func(t *testing.T) {
		expectVideoSource := "https://play2.sewobofang.com/20221113/LXfHZ2iw/index.m3u8"
		htmlVideo, err := os.ReadFile("../data/avlove4_video.txt")
		if err != nil {
			t.Error(err)
		}
		htmlParser := NewParseHtml(string(htmlVideo))
		htmlParser.SetParentUrl("https://avlove11.com")

		actualVideoSource, err := htmlParser.FindVideos()
		if err != nil || actualVideoSource == nil {
			t.Error(err)
		}
		assert.Equal(t, expectVideoSource, actualVideoSource.Source)
	})

	// TODO: FAIL
	t.Run("test find yadongpang video tag", func(t *testing.T) {
		expectVideoSource := "?"
		htmlVideo, err := os.ReadFile("../data/koreansexvid05_video.html")
		if err != nil {
			t.Error(err)
		}
		htmlParser := NewParseHtml(string(htmlVideo))
		htmlParser.SetParentUrl("https://koreansexvid05.com/")

		actualVideoSource, err := htmlParser.FindVideos()
		if err != nil || actualVideoSource == nil {
			t.Error(err)
		}
		assert.Equal(t, expectVideoSource, actualVideoSource.Source)
	})

	t.Run("test parse url", func(t *testing.T) {
		u, err := url.Parse("https://david.cdnbuzz.buzz/i.php?poster=https://yadong.best/data/file/korea/5481f222f6c5498a413608f0c9e69adb_gQ1NjZlM_38346a33fe05fadac317861f75d4291b4ae54a56.jpg&amp;vvv=n3/0427/세희 화보 영상.mp4")
		if err != nil {
			t.Error(err)
		}

		domain := u.Scheme + "://" + u.Hostname()
		t.Error(domain)
	})

	t.Run("test parse url", func(t *testing.T) {
		u, err := url.Parse("https://koreansexvid05.com/wp-content/plugins/clean-tube-player/public/player-x.php?q=cG9zdF9pZD0xMzMyNiZ0eXBlPWlmcmFtZSZ0YWc9JTNDaWZyYW1lJTIwc3JjJTNEJTIyaHR0cHMlM0ElMkYlMkZ3d3cueHZpZGVvcy5jb20lMkZlbWJlZGZyYW1lJTJGNzQwMzM1MTklMjIlMjBmcmFtZWJvcmRlciUzRCUyMjAlMjIlMjB3aWR0aCUzRCUyMjUxMCUyMiUyMGhlaWdodCUzRCUyMjQwMCUyMiUyMHNjcm9sbGluZyUzRCUyMm5vJTIyJTIwYWxsb3dmdWxsc2NyZWVuJTNEJTIyYWxsb3dmdWxsc2NyZWVuJTIyJTNFJTNDJTJGaWZyYW1lJTNF")
		if err != nil {
			t.Error(err)
		}

		domain := u.Scheme + "://" + u.Hostname()
		t.Error(domain)
	})
}
