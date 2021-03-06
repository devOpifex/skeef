package config

import (
	"encoding/json"
	"io/ioutil"
)

// Config Holds app configuration
type Config struct {
	TwitterConsumerKey    string `json:"twitterConsumerKey"`
	TwitterConsumerSecret string `json:"twitterConsumerSecret"`
	TwitterCallbackPath   string `json:"twitterCallbackPath"`
	Port                  string `json:"port"`
}

// ReadConfig Reads the configguration file
func Read() (Config, error) {
	var path = "config.json"
	var conf Config

	fl, err := ioutil.ReadFile(path)

	if err != nil {
		return conf, nil
	}

	err = json.Unmarshal(fl, &conf)

	if err != nil {
		return conf, nil
	}

	return conf, nil
}
