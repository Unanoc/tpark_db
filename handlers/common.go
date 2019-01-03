package handlers

import (
	"encoding/json"

	"github.com/valyala/fasthttp"
)

func toJSON(responseMessage json.Marshaler) []byte {
	responseMessageJSON, _ := responseMessage.MarshalJSON()
	return responseMessageJSON
}

func response(ctx *fasthttp.RequestCtx, status int, responseStruct json.Marshaler) {
	ctx.SetContentType("application/json")
	ctx.SetStatusCode(status)
	ctx.Write(toJSON(responseStruct))
}

func responseCustomError(ctx *fasthttp.RequestCtx, status int, err error) {
	ctx.SetContentType("application/json")
	ctx.SetStatusCode(status)
	ctx.SetBodyString(err.Error())
}

func responseDefaultError(ctx *fasthttp.RequestCtx, status int, err error) {
	ctx.SetContentType("application/json")
	ctx.SetStatusCode(status)
	ctx.SetBodyString(err.Error())
}
