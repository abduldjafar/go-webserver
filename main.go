package main

import (
	"go-webserver/cmd"
)

var (
	command cmd.Cmd = cmd.CmdServices()
)

func main() {
	command.InitCommand()
	command.Run()
}
