package main

import (
	"log"
	"time"

	"github.com/jasonlvhit/gocron"
	"github.com/motemen/go-pocket/api"
	"github.com/skyitachi/anyway2pocket/common"
	"github.com/skyitachi/anyway2pocket/zhihu"
)

func main() {
	pocketClient := common.NewPocketClient()
	dbClient := common.PocketDBClient{
		DBHost: "localhost",
		DBUser: "skyitachi",
		DBName: "pocket",
	}
	dbClient.Init()

	zh := zhihu.Zhihu{
		Name: "zhihu",
		CanNext: func(url string) bool {
			if !dbClient.URLExists(url) {
				return true
			}
			lastUpdated, err := dbClient.GetDateByURL(url)
			if err != nil {
				log.Println("get url date error: " + err.Error())
				return true
			}
			if common.GetSecondsDiff(time.Now(), lastUpdated) > 10 {
				return true
			}
			log.Println("success prevent the url " + url)
			return false
		},
		OnGetURL: func(url string) {
			if dbClient.URLExists(url) {
				log.Println("url has stored in pocket " + url)
				return
			}
			log.Println("Got URL: " + url)
			option := &api.AddOption{
				URL:   url,
				Title: "zhihu",
				Tags:  "zhihu",
			}
			err := pocketClient.Add(option)
			if err != nil {
				log.Fatal("add to pocket error: " + err.Error())
			}
			dbClient.AddURL(url, common.URLStatusFinished)
			log.Println("add to pocket success: " + url)
		},
		PageDone: func(url string) {
			dbClient.UpdateURL(url)
		},
		URLExists: func(url string) bool {
			return dbClient.URLExists(url)
		},
	}

	zh.Start()
	<-gocron.Start()
}
