package logger

import (
	"time"

	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

var Logger, _ = zap.NewProduction()

func LoggerHandler(handler fasthttp.RequestHandler) fasthttp.RequestHandler {
	return fasthttp.RequestHandler(func(ctx *fasthttp.RequestCtx) {
		start := time.Now()

		handler(ctx)

		Logger.Info("Request",
			zap.String("date", time.Now().String()),
			zap.String("method", string(ctx.Method())),
			zap.String("URI", string(ctx.RequestURI())),
			zap.Duration("duration", time.Since(start)),
		)
	})
}

func LoggerInfo(info string) {
	Logger.Info(info,
		zap.String("date", time.Now().String()),
	)
}
