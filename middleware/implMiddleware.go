package middleware

import (
	"log"
	"os"
	"time"

	"github.com/gogearbox/gearbox"
	"github.com/valyala/fasthttp"
)

type middleware struct{}

var (
	output = log.New(os.Stdout, "", 0)
)

func (*middleware) Logger(gearboxCtx gearbox.Context) {
	ctx := gearboxCtx.Context()
	begin := time.Now()
	end := time.Now()
	output.Printf("[%v] %v | %s | %s %s - %v - %v | %s",
		end.Format("2006/01/02 - 15:04:05"),
		ctx.RemoteAddr(),
		getHttp(ctx),
		ctx.Method(),
		ctx.RequestURI(),
		ctx.Response.Header.StatusCode(),
		end.Sub(begin),
		ctx.UserAgent(),
	)

	// Next is what allows the request to continue to the next
	// middleware/handler
	gearboxCtx.Next()
}
func ImplMiddleware() Middleware {
	return &middleware{}
}

func getHttp(ctx *fasthttp.RequestCtx) string {
	if ctx.Response.Header.IsHTTP11() {
		return "HTTP/1.1"
	}
	return "HTTP/1.0"
}
