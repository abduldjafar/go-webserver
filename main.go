package main

import (
	"go-webserver/cmd"
)

var (
	command cmd.Cmd = cmd.CmdServicesGearbox()
)

func main() {
	command.InitCommand()
	command.Run()
}
