package cmd

import (
	"flag"
	"go-webserver/api"
	"go-webserver/config"
	"log"
)

type CmdGearbox interface {
	Run()
	InitCommand()
}
type cmdGearbox struct {
	Port       string
	Ssl        bool
	Storage    string
	CrtFile    string
	CrtKeyFile string
	ConfigFile string
}

var (
	endpointsGearbox = api.GearboxEndpoint{}
)

func (c *cmdGearbox) InitCommand() {
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

func (c *cmdGearbox) Run() {

	initConfig := config.Configuration{}
	config.GetConfigUseFileParam(&initConfig, c.ConfigFile)

	endpointsGearbox.SetupConfig(&initConfig, c.Ssl, c.CrtFile, c.CrtKeyFile)
	endpointsGearbox.SetupStorage(c.Storage)
	endpointsGearbox.ALL()
	log.Println(c.Port)
	endpointsGearbox.Api().Start(":" + c.Port)
}

func CmdServicesGearbox() CmdGearbox {
	return &cmdGearbox{}
}
