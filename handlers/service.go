package handlers

import (
	"tpark_db/database"
	"tpark_db/helpers"

	"github.com/valyala/fasthttp"
)

func ClearHandler(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")
	if err := helpers.ExecSQLScript(database.DB.SchemaPath); err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError) // 500
		ctx.SetBodyString(err.Error())
	} else {
		ctx.SetStatusCode(fasthttp.StatusOK) // 200
		ctx.SetBodyString("null")
	}
}

func StatusHandler(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")
	if status, err := helpers.StatusHelper(); err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError) // 500
		ctx.SetBodyString(err.Error())
	} else {
		ctx.SetStatusCode(fasthttp.StatusOK) // 200
		buf, _ := status.MarshalJSON()
		ctx.SetBody(buf)
	}
}
