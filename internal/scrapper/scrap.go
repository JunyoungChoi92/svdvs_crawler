package scrapper

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
	"github.com/google/uuid"
)

type IframeData struct {
	Content   string
	OuterHTML string
}

func CaptureScreenshot(ctx context.Context, targetUrl string) (string, error) {
	var screenshot []byte
	scrollScript := `window.scrollTo(0, document.body.scrollHeight);`

	err := chromedp.Run(ctx,
		chromedp.Navigate(targetUrl),
		chromedp.WaitReady(`body`, chromedp.ByQuery),
		chromedp.Sleep(2*time.Second),
		chromedp.Evaluate(scrollScript, nil),
		chromedp.Sleep(2*time.Second),
		chromedp.FullScreenshot(&screenshot, 90),
	)
	if err != nil {
		return "", err
	}

	screenshotDir := "../../screenshots"

	if _, err := os.Stat(screenshotDir); os.IsNotExist(err) {
		err := os.MkdirAll(screenshotDir, os.ModePerm)
		if err != nil {
			return "", err
		}
	}

	// Create a unique file name for the screenshot.
	uuid := uuid.New().String()
	screenshot_path := filepath.Join(screenshotDir, "screenshot-"+uuid+".png")

	file, err := os.Create(screenshot_path)
	if err != nil {
		return "", err
	}

	_, err = file.Write(screenshot)
	if err != nil {
		return "", err
	}

	defer file.Close()

	return screenshot_path, nil
}

// chromedp.Navigate() automatically waits for the load event to fire, which happens when the whole page has loaded.
// But nowadays, many web pages adopted ajax feature, which makes it hard to tell whether the page has been loaded from the user perspective.
// In this case, we have to study the target web page to find out what does "page loaded" means for that page.
func ScrapHTML(ctx context.Context, targetURL string, iframeTimeOut time.Duration, isBlockedByCaptcha bool) (string, error) {
	var pageContent string
	var iframeNodes []*cdp.Node

	var tasks chromedp.Tasks

	// a one "chromedp.Run" function is called for solving "Nested chromedp.Run" problem.
	tasks = append(tasks, chromedp.Navigate(targetURL))

	if isBlockedByCaptcha {
		tasks = append(tasks,
			chromedp.QueryAfter("#warning", func(ctx context.Context, eci runtime.ExecutionContextID, nodes ...*cdp.Node) error {
				for _, node := range nodes {
					if id, _ := node.Attribute("id"); id == "warning" {
						// tasks = append(tasks, chromedp.MouseClickNode(node))
						return chromedp.MouseClickNode(node).Do(ctx)
					}
				}
				return nil
			}, chromedp.AtLeast(0)),
		)
	}

	tasks = append(tasks,
		chromedp.WaitReady(`body`, chromedp.ByQuery),
		chromedp.OuterHTML(`html`, &pageContent),
		chromedp.Nodes(`iframe`, &iframeNodes, chromedp.AtLeast(0)),
	)

	if err := chromedp.Run(ctx, tasks); err != nil {
		return "", err
	}

	log.Println("The number of iframeNodes: ", len(iframeNodes))

	iframeContentsChan := make(chan IframeData, len(iframeNodes))

	// If there's no iframe, close the channel.
	if len(iframeNodes) == 0 {
		close(iframeContentsChan)
	} else {
		// If there's iframe, scrap the iframe contents concurrently.
		var waitingGroup sync.WaitGroup

		for _, iframeNode := range iframeNodes {
			waitingGroup.Add(1)

			go func(iframeNode *cdp.Node) {
				defer waitingGroup.Done()

				var iframeContent string
				var iframeOuterHTML string

				// if chromedp connects to the iframe but there's no response in long time, it must be time out.
				iframeCtx, cancel := context.WithTimeout(ctx, iframeTimeOut)
				defer cancel()

				if err := chromedp.Run(iframeCtx,
					chromedp.OuterHTML(iframeNode.FullXPath(), &iframeOuterHTML, chromedp.AtLeast(0)),
					chromedp.Navigate(iframeNode.AttributeValue("src")),
					chromedp.WaitReady(`body`, chromedp.ByQuery),
					chromedp.Sleep(2*time.Second),
					chromedp.OuterHTML(`html`, &iframeContent, chromedp.AtLeast(0))); err != nil {
					log.Println(err)
					return
				}

				iframeContentsChan <- IframeData{Content: iframeContent, OuterHTML: iframeOuterHTML}
			}(iframeNode)
		}

		// Close the channel after all goroutines are finished.
		go func() {
			waitingGroup.Wait()
			close(iframeContentsChan)
		}()
	}

	for iframeData := range iframeContentsChan {
		index := strings.Index(pageContent, iframeData.OuterHTML)
		if index != -1 {
			insertionPoint := index + len(iframeData.OuterHTML)
			pageContent = pageContent[:insertionPoint] + "\n" + iframeData.Content + pageContent[insertionPoint:]
		} else {
			log.Printf("Warning: iframeData.OuterHTML not found in pageContent")
		}
	}

	// Count the number of "iframe" in the pageContent.
	iframeCount := strings.Count(pageContent, "<iframe")
	log.Println("The number of iframe in the pageContent: ", iframeCount)
	return pageContent, nil
}
