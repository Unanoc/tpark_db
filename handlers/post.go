package handlers

import (
	"strings"
	"tpark_db/errors"
	"tpark_db/helpers"
	"tpark_db/models"

	"github.com/valyala/fasthttp"
)

// PostUpdateHandler handles POST request /api/post/:id/details
func PostUpdateHandler(ctx *fasthttp.RequestCtx) {
	postID := ctx.UserValue("id").(string)
	postUpdate := models.PostUpdate{}
	if err := postUpdate.UnmarshalJSON(ctx.PostBody()); err != nil {
		responseDefaultError(ctx, fasthttp.StatusBadRequest, err) // 400
		return
	}

	result, err := helpers.PostUpdateHelper(&postUpdate, postID)
	switch err {
	case nil:
		response(ctx, fasthttp.StatusOK, result) // 200
	case errors.PostNotFound:
		responseCustomError(ctx, fasthttp.StatusNotFound, errors.PostNotFound) // 404
	default:
		responseDefaultError(ctx, fasthttp.StatusInternalServerError, err) // 500
	}
}

// PostGetOneHandler handles GET request /api/post/:id/details
func PostGetOneHandler(ctx *fasthttp.RequestCtx) {
	postID := ctx.UserValue("id").(string)
	relatedString := ctx.FormValue("related")
	relatedParams := []string{"post"}
	relatedParams = append(relatedParams, strings.Split(string(relatedString), ",")...)

	result, err := helpers.PostFullHelper(postID, relatedParams)
	switch err {
	case nil:
		response(ctx, fasthttp.StatusOK, result) // 200
	case errors.PostNotFound:
		responseCustomError(ctx, fasthttp.StatusNotFound, errors.PostNotFound) // 404
	default:
		responseDefaultError(ctx, fasthttp.StatusInternalServerError, err) // 500
	}
}
