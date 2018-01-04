package main

import (
	"github.com/jasonlvhit/gocron"
	_ "github.com/skyitachi/anyway2pocket/common"
	"github.com/skyitachi/anyway2pocket/zhihu"
)

func main() {
	zhihu.PullCollection("https://www.zhihu.com/collection/119397553")

	// gocron.Every(1).Minute().Do(PullCollection, "https://www.zhihu.com/collection/119397553")
	<-gocron.Start()
}
