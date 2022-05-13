package middleware

import "github.com/gogearbox/gearbox"

type Middleware interface {
	Logger(gearboxCtx gearbox.Context)
}
