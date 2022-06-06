package controller

import (
	idxauth "go-webserver/auth"
	"go-webserver/config"
	"go-webserver/services"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gogearbox/gearbox"
)

type gearboxFile struct {
	config *config.Configuration
}

var (
	auths      idxauth.Auth      = idxauth.ImplAuthService()
	idxservice services.Services = services.ImplementServices()
)

func (g *gearboxFile) SetupConfig(config *config.Configuration) {
	g.config = config
	log.Println(g.config)
}

func (g *gearboxFile) GenerateFileToken() interface{} {
	return func(ctx gearbox.Context) {
		token := auths.CreateToken("idx", 7)
		topic := ctx.Param("topic")

		initConfig := g.config

		url := initConfig.Kafka.UrlProducer + "/cron"

		idxservice.PostTokenToKafkaCLient(token, topic, url)

		ctx.SendJSON(map[string]interface{}{
			"token": token,
		})

	}
}
func (g *gearboxFile) Get(download_path string) interface{} {
	return func(ctx gearbox.Context) {
		fileName := ctx.Param("filename")
		queryFileName := ctx.Query("file_name")
		log.Println(queryFileName)

		token := ctx.Param("token")

		_, err := auths.Validate(token, "idx")

		if err != nil {
			ctx.SendJSON(map[string]interface{}{
				"error": err.Error(),
				"code":  500,
			})
		}

		fileName = strings.Replace(fileName, "lic/", "", -1)

		targetPath := filepath.Join(download_path, fileName)
		log.Println(download_path)
		log.Println(fileName)

		if _, err := os.Stat(targetPath); err != nil {
			ctx.SendJSON(map[string]interface{}{
				"error": err,
				"code":  500,
			})
		}

		log.Println(targetPath)
		ctx.Context().Response.Header.Set("Content-Description", "File Transfer")
		ctx.Context().Response.Header.Set("Content-Transfer-Encoding", "binary")
		ctx.Context().Response.Header.Set("Content-Disposition", "attachment; filename="+fileName)
		ctx.Context().Response.Header.Set("Content-Type", "application/octet-stream")
		ctx.Context().SendFile(targetPath)

	}
}

func (g *gearboxFile) GetWithQuery(download_path string) interface{} {
	return func(ctx gearbox.Context) {
		fileName := ctx.Query("file_name")

		token := ctx.Query("token")

		_, err := auths.Validate(token, "idx")

		if err != nil {
			ctx.SendJSON(map[string]interface{}{
				"error": err.Error(),
				"code":  500,
			})
		}

		fileName = strings.Replace(fileName, "lic/", "", -1)

		targetPath := filepath.Join(download_path, fileName)
		log.Println(download_path)
		log.Println(fileName)

		if _, err := os.Stat(targetPath); err != nil {
			ctx.SendJSON(map[string]interface{}{
				"error": err,
				"code":  500,
			})
		}

		log.Println(targetPath)
		ctx.Context().Response.Header.Set("Content-Description", "File Transfer")
		ctx.Context().Response.Header.Set("Content-Transfer-Encoding", "binary")
		ctx.Context().Response.Header.Set("Content-Disposition", "attachment; filename="+fileName)
		ctx.Context().Response.Header.Set("Content-Type", "application/octet-stream")
		ctx.Context().SendFile(targetPath)

	}
}

func (g *gearboxFile) Create(path string) interface{} {
	return func(ctx gearbox.Context) {
		initConfig := g.config

		form, err := ctx.Context().Request.MultipartForm()
		if err != nil {
			ctx.SendJSON(map[string]interface{}{
				"error": err.Error(),
				"code":  500,
			})
		}

		fileHeader := form.File["file"][0]
		idxgroup := form.Value["idxgroup"][0]
		idxtotal := form.Value["idxtotal"][0]
		idxnumber := form.Value["idxnumber"][0]
		url := initConfig.Kafka.UrlProducer + "/produces"
		filename := strings.ToLower(fileHeader.Filename)
		fullpathname := initConfig.Kafka.HostUrl + filename

		file, err := fileHeader.Open()
		if err != nil {
			ctx.SendJSON(map[string]interface{}{
				"error": err.Error(),
				"code":  500,
			})
		}

		go idxservice.PostPathToKafkaClient(filename, fullpathname, idxgroup, initConfig.Kafka.Topic, idxtotal, idxnumber, file, path, url)
		go idxservice.PostPathToKafkaClient(filename, initConfig.Kafka.HostUrl+"/?file_name="+filename, idxgroup, initConfig.Kafka.Topic+"_query_path", idxtotal, idxnumber, file, path, url)

		//go g.sendTokafkaCLient(fileHeader.Filename, idxgroup, initConfig.Kafka.Topic, idxtotal, idxnumber, file, path)

	}
}

func GearboxImplFileController() FileController {
	return &gearboxFile{}
}
