package config

import (
	"fmt"
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

type Configuration struct {
	Kafka kafka
}

type ConfigFile struct {
	File string
}

type kafka struct {
	Topic       string
	FileConfig  string
	UrlProducer string
	HostUrl     string
	TokenTopic  string
	CronTime    string
}

func GetConfig(baseConfig *Configuration) {
	basePath, _ := os.Getwd()
	if _, err := toml.DecodeFile(basePath+"/config.toml", &baseConfig); err != nil {
		fmt.Println(err)
	}
}

func GetConfigUseFileParam(baseConfig *Configuration, Filename string) {
	if _, err := toml.DecodeFile(Filename, &baseConfig); err != nil {
		log.Fatalln(err)
	}
}
