package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/jasonlvhit/gocron"
)

const ZHIHU_HOST = "https://www.zhihu.com"

type CollectionItem struct {
	Url   string
	Title string
}

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

func download(url string) []byte {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Panic(err)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Print(err)
	}
	defer resp.Body.Close()
	var buf bytes.Buffer
	buf.ReadFrom(resp.Body)
	return buf.Bytes()
}

func pullCollection(collectionURL string) {
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
				pullCollection(nextPageURL)
			}
		}
	})

}

func main() {
	//body := download("https://www.zhihu.com/collection/119397553")
	pullCollection("https://www.zhihu.com/collection/119397553")
	// gocron.Every(1).Minute().Do(PullCollection, "https://www.zhihu.com/collection/119397553")
	<-gocron.Start()
}
