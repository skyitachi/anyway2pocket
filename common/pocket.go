package common

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"path"
	"runtime"

	"github.com/motemen/go-pocket/api"
)

// PocketConfig pocket config
type PocketConfig struct {
	ConsumerKey string `json:"consumerKey"`
	AccessToken string `json:"accessToken"`
}

// NewPocketClient return new pocket client
func NewPocketClient() *api.Client {
	var config PocketConfig
	_, filename, _, ok := runtime.Caller(1)
	if ok {
		content, err := ioutil.ReadFile(path.Join(path.Dir(filename), "pocket.json"))
		if err != nil {
			log.Fatal("[NewPocketClient]: " + err.Error())
		}
		err = json.Unmarshal(content, &config)
		if err != nil {
			log.Fatal("[NewPocketClient]: " + err.Error())
		}
		return api.NewClient(config.ConsumerKey, config.AccessToken)
	}
	log.Fatal("[NewPocketClient]: cannot find pocket.json path")
	return nil
}
