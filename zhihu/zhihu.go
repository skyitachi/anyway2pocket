package zhihu

import (
	"fmt"
	"log"
	"net/url"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/skyitachi/anyway2pocket/common"
)

const ZHIHU_HOST = "https://www.zhihu.com"

func checkError(err error) {
	if err != nil {
		log.Println(err)
	}
}

func buildURL(url string) string {
	matches, err := regexp.Match("zhihu.com", []byte(url))
	if err != nil {
		return url
	}
	if !matches {
		return ZHIHU_HOST + url
	}
	return url
}

func buildNextPageURL(currentURL string, nextURL string) string {
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

// PullCollection crawl zhihu page
func PullCollection(collectionURL string) {
	exists := common.URLExists(collectionURL)
	if exists {
		log.Println("[PullCollection]: had searched the url " + collectionURL)
		return
	}
	log.Println("[PullCollection]: current url " + collectionURL)
	doc, err := goquery.NewDocument(collectionURL)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(strings.TrimSpace(doc.Find("#zh-fav-head-title").Text()))
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
					fmt.Println(buildURL(url))
				}
			}
		case "Post":
			child := s.Find(".post-link")
			if child != nil {
				url, ok := child.Attr("href")
				if ok {
					fmt.Println(buildURL(url))
				}
			}
		}
	})
	// 判断是否有分页器
	doc.Find(".border-pager .zm-invite-pager a").Each(func(i int, s *goquery.Selection) {
		text := s.Text()
		if text != "下一页" {
			link, ok := s.Attr("href")
			if ok {
				nextPageURL := buildNextPageURL(collectionURL, link)
				PullCollection(nextPageURL)
			}
		}
	})
}
