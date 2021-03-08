package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
)

// Config Holds app configuration
type Config struct {
	TwitterConsumerKey    string `json:"twitterConsumerKey"`
	TwitterConsumerSecret string `json:"twitterConsumerSecret"`
	TwitterAccessToken    string `json:"twitterAccessToken"`
	TwitterAccessSecret   string `json:"twitterAccessSecret"`
	Port                  string `json:"port"`
}

// Read the configuration file
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

	err = conf.check()

	if err != nil {
		return conf, err
	}

	conf.Port = ":" + conf.Port

	return conf, nil
}

// Check the config file
func (conf *Config) check() error {

	if conf.Port == "" {
		fmt.Println("No `port` specified in the configuration file (`config.json`), defaulting to 8080")
		conf.Port = "8080"
	}

	if conf.TwitterAccessSecret == "" || conf.TwitterAccessToken == "" || conf.TwitterConsumerKey == "" || conf.TwitterConsumerSecret == "" {
		return errors.New("twitter* fields not specified in the configuration file")
	}

	return nil
}
