package zhihu

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/url"
	"path"
	"regexp"
	"runtime"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/jasonlvhit/gocron"
	"github.com/skyitachi/anyway2pocket/common"
)

func checkError(err error) {
	if err != nil {
		log.Println(err)
	}
}

const ZHIHU_HOST = "https://www.zhihu.com"

// Zhihu implements Crawler interface
type Zhihu struct {
	Name      string
	CanNext   func(string) bool
	OnGetURL  func(string)
	PageDone  func(string)
	URLExists func(string) bool
}

// PullCollection pull collection url
func (z Zhihu) PullCollection(collectionURL string) {
	logger := common.GetLogger()
	if !z.CanNext(collectionURL) {
		logger.Info("[PullCollection]: had searched the url " + collectionURL)
		return
	}
	logger.Info("[PullCollection]: current url " + collectionURL)
	go func() {
		z.PageDone(collectionURL)
	}()
	doc, err := goquery.NewDocument(collectionURL)
	if err != nil {
		logger.Error("[PullCollection]: " + err.Error())
		return
	}
	doc.Find(".zm-item").Each(func(i int, s *goquery.Selection) {
		dataType, ok := s.Attr("data-type")
		if !ok {
			return
		}
		switch dataType {
		case "Answer":
			child := s.Find(".zm-item-rich-text")
			if child != nil {
				url, ok := child.Attr("data-entry-url")
				if ok {
					z.OnGetURL(z.buildURL(url))
				}
			}
		case "Post":
			child := s.Find(".post-link")
			if child != nil {
				url, ok := child.Attr("href")
				if ok {
					z.OnGetURL(z.buildURL(url))
				}
			}
		}
	})
	// 判断是否有分页器
	doc.Find(".border-pager .zm-invite-pager a").Each(func(i int, s *goquery.Selection) {
		text := s.Text()
		if text == "下一页" {
			link, ok := s.Attr("href")
			if ok {
				nextPageURL := z.buildNextPageURL(collectionURL, link)
				go func() {
					time.Sleep(1 * time.Second)
					z.PullCollection(nextPageURL)
				}()
			}
		}
	})
}

func (z Zhihu) buildURL(url string) string {
	matches, err := regexp.Match("zhihu.com", []byte(url))
	if err != nil {
		return url
	}
	if !matches {
		return ZHIHU_HOST + url
	}
	return url
}

func (z Zhihu) buildNextPageURL(currentURL string, nextURL string) string {
	cURL, err := url.Parse(currentURL)
	checkError(err)
	nURL, err := url.Parse(nextURL)
	checkError(err)
	if nURL.Query().Get("page") == "" {
		log.Println("[NextURL]: next url is wrong " + nextURL)
		return currentURL
	}
	q := cURL.Query()
	q.Set("page", nURL.Query().Get("page"))
	cURL.RawQuery = q.Encode()
	return cURL.String()
}

// GetLatestCollection get all uncrawled url
func (z Zhihu) GetLatestCollection(collectionURL string) {
	logger := common.GetLogger()
	found := false
	doc, err := goquery.NewDocument(collectionURL)
	if err != nil {
		logger.Error(err)
		return
	}
	doc.Find(".zm-item").Each(func(i int, s *goquery.Selection) {
		if found {
			return
		}
		dataType, ok := s.Attr("data-type")
		if !ok {
			return
		}
		switch dataType {
		case "Answer":
			child := s.Find(".zm-item-rich-text")
			if child != nil {
				url, ok := child.Attr("data-entry-url")
				if ok {
					santizedURL := z.buildURL(url)
					if z.URLExists(santizedURL) {
						found = true
					} else {
						z.OnGetURL(santizedURL)
					}
				}
			}
		case "Post":
			child := s.Find(".post-link")
			if child != nil {
				url, ok := child.Attr("href")
				if ok {
					santizedURL := z.buildURL(url)
					if z.URLExists(santizedURL) {
						found = true
					} else {
						z.OnGetURL(santizedURL)
					}
				}
			}
		}
	})
	if found {
		return
	}
	// 判断是否有分页器
	doc.Find(".border-pager .zm-invite-pager a").Each(func(i int, s *goquery.Selection) {
		text := s.Text()
		if text == "下一页" {
			link, ok := s.Attr("href")
			if ok {
				nextPageURL := z.buildNextPageURL(collectionURL, link)
				go func() {
					time.Sleep(1 * time.Second)
					z.GetLatestCollection(nextPageURL)
				}()
			}
		}
	})

}

// Start start cron tasks
func (z Zhihu) Start() {
	logger := common.GetLogger()
	var collectionList []string
	_, filename, _, ok := runtime.Caller(1)
	if ok {
		raw, err := ioutil.ReadFile(path.Join(path.Dir(filename), "zhihu/collection.json"))
		if err != nil {
			logger.Fatal("[Zhihu.Start]: " + err.Error())
		}
		err = json.Unmarshal(raw, &collectionList)
		for _, url := range collectionList {
			task := func(startUrl string) func() {
				return func() {
					logger.Info("[Zhihu.Start]: start url " + startUrl)
					z.GetLatestCollection(startUrl)
				}
			}(url)
			gocron.Every(5).Minutes().Do(task)
		}
	} else {
		logger.Fatal("[Zhihu.Start]: cannot get collection.json path")
	}
}
