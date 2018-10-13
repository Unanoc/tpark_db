package logger

import (
	"log"
	"time"

	"github.com/valyala/fasthttp"
)

func Logger(handler fasthttp.RequestHandler) fasthttp.RequestHandler {
	return fasthttp.RequestHandler(func(ctx *fasthttp.RequestCtx) {
		start := time.Now()

		handler(ctx)

		log.Printf(
			"%s %s %s",
			ctx.Method(),
			ctx.RequestURI(),
			time.Since(start),
		)
	})
}
