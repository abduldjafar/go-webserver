package cmd

import (
	"flag"
	"go-webserver/api"
	"go-webserver/config"
	"go-webserver/cron"
	"log"
)

type Cmd interface {
	Run()
	InitCommand()
}
type cmd struct {
	Port       string
	Ssl        bool
	Storage    string
	CrtFile    string
	CrtKeyFile string
	ConfigFile string
}

var (
	endpoints                     = api.GinEndpoint{}
	cronToken cron.CronController = cron.CronScheduler()
)

func (c *cmd) InitCommand() {
	Port := flag.String("port", "8000", "port for listening")
	Ssl := flag.Bool("ssl", false, "use ssl or not (default 'false')")
	Storage := flag.String("storage", "public", "path for save file that being uploaded")
	CrtFile := flag.String("crt_file", "example.crt", "ssl crt file")
	CrtKeyFile := flag.String("crt_key_file", "example.key", "ssl key file")
	ConfigFile := flag.String("config", "config.toml", "file config for runing service")

	flag.Parse()

	c.CrtFile = *CrtFile
	c.CrtKeyFile = *CrtKeyFile
	c.Storage = *Storage
	c.Port = *Port
	c.Ssl = *Ssl
	c.ConfigFile = *ConfigFile

}

func (c *cmd) Run() {

	initConfig := config.Configuration{}
	config.GetConfigUseFileParam(&initConfig, c.ConfigFile)

	endpoints.SetupConfig(&initConfig)
	cronToken.SetupConfig(&initConfig)
	go cronToken.Scheduler()
	endpoints.SetupStorage(c.Storage)
	endpoints.ALL()

	if c.Ssl {
		log.Println("runing with ssl in port " + c.Port)
		endpoints.Api().RunTLS(":"+c.Port, c.CrtFile, c.CrtKeyFile)
	} else {
		log.Println("runing without ssl in port " + c.Port)
		endpoints.Api().Run(":" + c.Port)
	}

}

func CmdServices() Cmd {
	return &cmd{}
}
