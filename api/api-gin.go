package api

import (
	"go-webserver/config"
	"go-webserver/controller"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	ginlogrus "github.com/toorop/gin-logrus"
)

type GinEndpoint struct {
	Router        *gin.Engine
	StoragePath   string
	Configuration *config.Configuration
}

var (
	ginRouter                                 = gin.Default()
	v1                                        = ginRouter.Group("/v1")
	log                                       = logrus.New()
	fileControllers controller.FileController = controller.GinImplFileController()
)

func (g *GinEndpoint) ALL() {

	gin.SetMode(gin.ReleaseMode)

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}

	ginRouter.Use(ginlogrus.Logger(log), gin.Recovery(), cors.New(config))

	v1.POST("/idx/upload",
		cors.New(config),
		fileControllers.Create(g.StoragePath).((func(*gin.Context))),
	)

	ginRouter.GET("/static/:filename/:token",

		cors.New(config),
		fileControllers.Get(g.StoragePath).((func(*gin.Context))),
	)

	ginRouter.GET("/idx/token/:topic",

		cors.New(config),
		fileControllers.GenerateFileToken().((func(*gin.Context))),
	)
	g.Router = ginRouter
}

func (g *GinEndpoint) SetupStorage(path string) {
	g.StoragePath = path
}

func (g *GinEndpoint) SetupConfig(config *config.Configuration) {
	g.Configuration = config
	fileControllers.SetupConfig(config)
}
func (g *GinEndpoint) Api() *gin.Engine {
	return g.Router
}
