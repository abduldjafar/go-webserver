package api

import (
	"go-webserver/config"
	"go-webserver/controller"
	"go-webserver/middleware"

	"github.com/gogearbox/gearbox"
)

type GearboxEndpoint struct {
	Router        gearbox.Gearbox
	StoragePath   string
	Configuration *config.Configuration
}

var (
	//gearboxRouter = gearbox.New()

	fileController controller.FileController = controller.GearboxImplFileController()
	newMiddleware  middleware.Middleware     = middleware.ImplMiddleware()
)

func (g *GearboxEndpoint) ALL() {

	g.Router.Use(newMiddleware.Logger)

	g.Router.Group("/v1", []*gearbox.Route{

		g.Router.Post("/idx/upload", fileController.Create(g.StoragePath).(func(ctx gearbox.Context))),
		g.Router.Get("/idx/token/:topic", fileController.GenerateFileToken().(func(ctx gearbox.Context))),
	})
	g.Router.Get("/static/:filename/:token", fileController.Get(g.StoragePath).(func(ctx gearbox.Context)))

	//g.Router = gearboxRouter

}
func (g *GearboxEndpoint) SetupStorage(path string) {
	g.StoragePath = path
}

func (g *GearboxEndpoint) SetupConfig(config *config.Configuration, tls bool, crt_file string, crt_key string) {
	g.Configuration = config
	g.Router = gearbox.New(&gearbox.Settings{
		ServerName:  "gearbox",
		TLSEnabled:  tls,
		TLSCertPath: crt_file,
		TLSKeyPath:  crt_key,
	})
	fileController.SetupConfig(config)
}

func (g *GearboxEndpoint) Api() gearbox.Gearbox {
	return g.Router
}
