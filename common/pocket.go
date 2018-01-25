package common

import (
	"encoding/json"
	"io/ioutil"
	"log"

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
	content, err := ioutil.ReadFile("../pocket.json")
	if err != nil {
		log.Fatal("[NewPocketClient]: " + err.Error())
	}
	err = json.Unmarshal(content, &config)
	if err != nil {
		log.Fatal("[NewPocketClient]: " + err.Error())
	}
	return api.NewClient(config.ConsumerKey, config.AccessToken)
}
