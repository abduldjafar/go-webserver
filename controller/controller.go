package controller

import "go-webserver/config"

type FileController interface {
	Create(path string) interface{}
	Get(download_path string) interface{}
	SetupConfig(config *config.Configuration)
	GenerateFileToken() interface{}
}
