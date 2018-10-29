package post

import (
	"fmt"
	"strings"
	"tpark_db/errors"
	"tpark_db/helpers"
	"tpark_db/models"

	"github.com/valyala/fasthttp"
)

func PostGetOneHandler(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")

	postID := ctx.UserValue("id").(string)
	relatedString := ctx.FormValue("related")
	relatedParams := strings.Split(string(relatedString), ",")
	result, err := helpers.PostFullHelper(postID, relatedParams)

	switch err {
	case nil:
		ctx.SetStatusCode(fasthttp.StatusOK) // 200
		buf, _ := result.MarshalJSON()
		ctx.SetBody(buf)
	case errors.PostNotFound:
		ctx.SetStatusCode(fasthttp.StatusNotFound) // 404
		errorResp := errors.Error{
			Message: fmt.Sprintf("Can't find post by id: %s", postID),
		}
		buf, _ := errorResp.MarshalJSON()
		ctx.SetBody(buf)
	default:
		ctx.SetStatusCode(fasthttp.StatusInternalServerError) // 500
		ctx.SetBodyString(err.Error())
	}
}

func PostUpdateHandler(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")

	postID := ctx.UserValue("id").(string)
	postUpdate := models.PostUpdate{}
	err := postUpdate.UnmarshalJSON(ctx.PostBody())
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest) // 400 Bad Request
		ctx.SetBodyString(err.Error())
		return
	}
	result, err := helpers.PostUpdateHelper(&postUpdate, postID)

	switch err {
	case nil:
		ctx.SetStatusCode(fasthttp.StatusOK) // 200
		buf, _ := result.MarshalJSON()
		ctx.SetBody(buf)
	case errors.PostNotFound:
		ctx.SetStatusCode(fasthttp.StatusNotFound) // 404
		errorResp := errors.Error{
			Message: fmt.Sprintf("Can't find post by id: %s", postID),
		}
		buf, _ := errorResp.MarshalJSON()
		ctx.SetBody(buf)
	default:
		ctx.SetStatusCode(fasthttp.StatusInternalServerError) // 500
		ctx.SetBodyString(err.Error())
	}
}
