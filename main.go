package main

import (
	"log"
	"time"

	"github.com/jasonlvhit/gocron"
	"github.com/skyitachi/anyway2pocket/common"
	"github.com/skyitachi/anyway2pocket/zhihu"
)

func main() {
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
			log.Println(time.Now())
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
			dbClient.AddURL(url, common.URLStatusFinished)
			log.Println("add to pocket success: " + url)
		},
		PageDone: func(url string) {
			dbClient.UpdateURL(url)
		},
	}

	zh.Start()

	<-gocron.Start()
}
