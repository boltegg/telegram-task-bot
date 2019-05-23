package main

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

const (
	CONFIG_FILE = "config.yaml"
)

type Config struct {
	TelegramBotApiKey string `yaml:"telegram_bot_api_key"`
}

var config Config

func InitConfig() {
	b, err := ioutil.ReadFile(CONFIG_FILE)
	if err != nil {
		logrus.Fatal(err)
	}

	err = yaml.Unmarshal(b, &config)
	if err != nil {
		logrus.Fatal(err)
	}
}