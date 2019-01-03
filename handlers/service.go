package handlers

import (
	"tpark_db/errors"
	"tpark_db/helpers"

	"github.com/valyala/fasthttp"
)

// ClearHandler handles POST request /api/service/clear.
func ClearHandler(ctx *fasthttp.RequestCtx) {
	helpers.ClearHelper()
	responseCustomError(ctx, fasthttp.StatusOK, errors.New("null"))
}

// StatusHandler handles GET request /api/service/status.
func StatusHandler(ctx *fasthttp.RequestCtx) {
	response(ctx, fasthttp.StatusOK, helpers.StatusHelper())
}
