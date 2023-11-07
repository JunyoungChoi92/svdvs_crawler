package pkg

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestScrapWebSite(t *testing.T) {
	ctx, cancel, err := InitializeBrowserContext()
	if err != nil {
		t.Errorf("Error occurred while initializing browser context. \n error log is here\n%s", err)
	}
	defer cancel()

	t.Run("test to scrap website https://www.naver.com", func(t *testing.T) {
		targetURL := "https://www.naver.com"
		scrapper(t, ctx, targetURL, 20*time.Second, false)
	})

	t.Run("test to scrap website https://scsj-30.com/bbs/board.php?bo_table=kr", func(t *testing.T) {
		targetURL := "https://scsj-30.com/bbs/board.php?bo_table=kr"
		scrapper(t, ctx, targetURL, 20*time.Second, false)
	})
	t.Run("test to scrap website https://sexports28.com/bbs/board.php?bo_table=kr", func(t *testing.T) {
		targetURL := "https://sexports28.com/bbs/board.php?bo_table=kr"
		scrapper(t, ctx, targetURL, 20*time.Second, false)
	})
	t.Run("test to scrap website https://av19.org/korea", func(t *testing.T) {
		targetURL := "https://av19.org/korea"
		scrapper(t, ctx, targetURL, 20*time.Second, false)
	})
	t.Run("test to scrap website https://www.bamhub14.me/index.php?so=korea&v=s", func(t *testing.T) {
		targetURL := "https://www.bamhub14.me/index.php?so=korea&v=s"
		scrapper(t, ctx, targetURL, 20*time.Second, false)
	})
	t.Run("https://cocktail41.com/bbs/board.php?bo_table=kr", func(t *testing.T) {
		targetURL := "https://cocktail41.com/bbs/board.php?bo_table=kr"
		scrapper(t, ctx, targetURL, 20*time.Second, false)
	})
}

func TestValidation(t *testing.T) {
	ctx, cancel, err := InitializeBrowserContext()
	if err != nil {
		t.Errorf("Error occurred while initializing browser context. \n\n%s", err)
	}
	defer cancel()

	indicators := []string{`script src="https://js.hcaptcha.com/1/api.js"`, `script src="https://www.google.com/recaptcha/api.js"`}

	t.Run("test to redirection check with https://pronhub.com/", func(t *testing.T) {
		targetURL := "https://pronhub.com/"
		expectedRedirectedURL := "https://www.pornhub.com/?utm_source=pronhub.com&utm_medium=redirects&utm_campaign=tldredirects"

		validator(t, ctx, targetURL, expectedRedirectedURL, indicators)
	})

	t.Run("test to redirection check with http://xandar.com/", func(t *testing.T) {
		targetURL := "http://xandar.com/"
		expectedRedirectedURL := "https://xkcorp.com/"

		validator(t, ctx, targetURL, expectedRedirectedURL, indicators)
	})
}

func scrapper(t *testing.T, ctx context.Context, targetUrl string, iframeTimeOut time.Duration, isBlockedByCaptcha bool) {
	pageContent, screenshotPath, err := ScrapWebpage(ctx, targetUrl, 20*time.Second, false)
	if err != nil {
		assert.Fail(t, err.Error())
	}

	if len(pageContent) == 0 {
		message := "pageContent is empty"
		assert.NotContains(t, pageContent, message)
	}

	if len(screenshotPath) == 0 {
		message := "screenshotPath is empty"
		assert.NotContains(t, pageContent, message)
	}
}

func validator(t *testing.T, ctx context.Context, targetUrl string, expectedRedirectedURL string, indicators []string) {
	redirectedURL, err := GetRedirectedUrl(ctx, targetUrl, 20*time.Second)
	if err != nil {
		assert.NotEqual(t, redirectedURL, expectedRedirectedURL)
	}

	if redirectedURL != expectedRedirectedURL {
		assert.Fail(t, "redirected URL and expected Redirected URL is not same. \nredirectedURL: %s\nexpectedRedirectedURL:%s", redirectedURL, expectedRedirectedURL)
	}

	err = CheckURLValidation(ctx, redirectedURL, indicators)
	if err != nil {
		assert.Fail(t, err.Error())
	}
}
