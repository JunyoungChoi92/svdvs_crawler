package pkg

import (
	"testing"
)

func TestCrawler_Command(t *testing.T) {

	t.Run("test case 1", func(t *testing.T) {
		Crawling("https://av19.org/korea/6463")
	})

	t.Run("test case 2", func(t *testing.T) {
		Crawling("https://av19.org/korea")
	})

	t.Run("test case 3", func(t *testing.T) {
		Crawling("https://kr21.mysstv.com/Video-List-1/-/-/1/35/new.html")
	})
}
