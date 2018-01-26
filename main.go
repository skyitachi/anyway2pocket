package main

import (
	"time"

	logging "github.com/hhkbp2/go-logging"
	"github.com/jasonlvhit/gocron"
	"github.com/motemen/go-pocket/api"
	"github.com/skyitachi/anyway2pocket/common"
	"github.com/skyitachi/anyway2pocket/zhihu"
)

func main() {
	common.InitLogger("./pocket.log")
	logger := common.GetLogger()

	pocketClient := common.NewPocketClient()
	defer logging.Shutdown()
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
				logger.Error("[CanNext]: get url date error " + err.Error())
				return true
			}
			if common.GetSecondsDiff(time.Now(), lastUpdated) > 10 {
				return true
			}
			logger.Info("[CanNext]: success prevent the url " + url)
			return false
		},
		OnGetURL: func(url string) {
			if dbClient.URLExists(url) {
				logger.Info("[OnGetUrl]: url has stored in pocket " + url)
				return
			}
			logger.Info("[OnGetUrl]: Got URL: " + url)
			option := &api.AddOption{
				URL:   url,
				Title: "zhihu",
				Tags:  "zhihu",
			}
			err := pocketClient.Add(option)
			if err != nil {
				logger.Error("[]add to pocket error: " + err.Error())
				return
			}
			err = dbClient.AddURL(url, common.URLStatusFinished)
			if err == nil {
				logger.Info("add to pocket success: " + url)
			}
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
