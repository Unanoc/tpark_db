package handlers

import (
	"tpark_db/errors"
	"tpark_db/helpers"
	"tpark_db/models"

	"github.com/valyala/fasthttp"
)

// ThreadCreateHandler handles POST request /api/thread/:slug_or_id/create.
func ThreadCreateHandler(ctx *fasthttp.RequestCtx) {
	slugOrID := ctx.UserValue("slug_or_id").(string)
	posts := models.Posts{}
	if err := posts.UnmarshalJSON(ctx.PostBody()); err != nil {
		responseDefaultError(ctx, fasthttp.StatusBadRequest, err) // 400
		return
	}

	result, err := helpers.ThreadCreateHelper(&posts, slugOrID)
	switch err {
	case nil:
		response(ctx, fasthttp.StatusCreated, result) // 201
	case errors.ThreadNotFound:
		responseCustomError(ctx, fasthttp.StatusNotFound, errors.ThreadNotFound) // 404
	case errors.UserNotFound:
		responseCustomError(ctx, fasthttp.StatusNotFound, errors.UserNotFound) // 404
	case errors.PostParentNotFound:
		responseCustomError(ctx, fasthttp.StatusConflict, errors.PostParentNotFound) // 404
	default:
		responseDefaultError(ctx, fasthttp.StatusInternalServerError, err) // 500
	}
}

// ThreadUpdateHandler handles POST request /api/thread/:slug_or_id/details.
func ThreadUpdateHandler(ctx *fasthttp.RequestCtx) {
	slugOrID := ctx.UserValue("slug_or_id").(string)
	threadUpdate := models.ThreadUpdate{}
	if err := threadUpdate.UnmarshalJSON(ctx.PostBody()); err != nil {
		responseDefaultError(ctx, fasthttp.StatusBadRequest, err) // 400
		return
	}

	result, err := helpers.ThreadUpdateHelper(&threadUpdate, slugOrID)
	switch err {
	case nil:
		response(ctx, fasthttp.StatusOK, result) // 200
	case errors.ThreadNotFound:
		responseCustomError(ctx, fasthttp.StatusNotFound, errors.ThreadNotFound) // 404
	default:
		responseDefaultError(ctx, fasthttp.StatusInternalServerError, err) // 500
	}
}

// ThreadVoteHandler handles POST request /api/thread/:slug_or_id/vote.
func ThreadVoteHandler(ctx *fasthttp.RequestCtx) {
	slugOrID := ctx.UserValue("slug_or_id").(string)
	vote := models.Vote{}
	if err := vote.UnmarshalJSON(ctx.PostBody()); err != nil {
		responseDefaultError(ctx, fasthttp.StatusBadRequest, err) // 400
		return
	}

	result, err := helpers.ThreadVoteHelper(&vote, slugOrID)
	switch err {
	case nil:
		response(ctx, fasthttp.StatusOK, result) // 200
	case errors.ThreadNotFound:
		responseCustomError(ctx, fasthttp.StatusNotFound, errors.ThreadNotFound) // 404
	default:
		responseDefaultError(ctx, fasthttp.StatusInternalServerError, err) // 500
	}
}

// ThreadGetOneHandler handles GET request /api/thread/:slug_or_id/details.
func ThreadGetOneHandler(ctx *fasthttp.RequestCtx) {
	slugOrID := ctx.UserValue("slug_or_id").(string)

	result, err := helpers.GetThreadBySlugOrID(slugOrID)
	switch err {
	case nil:
		response(ctx, fasthttp.StatusOK, result) // 200
	case errors.ThreadNotFound:
		responseCustomError(ctx, fasthttp.StatusNotFound, errors.ThreadNotFound) // 404
	default:
		responseDefaultError(ctx, fasthttp.StatusInternalServerError, err) // 500
	}
}

// ThreadGetPostsHandler handles GET request /api/thread/:slug_or_id/posts.
func ThreadGetPostsHandler(ctx *fasthttp.RequestCtx) {
	slugOrID := ctx.UserValue("slug_or_id").(string)
	limit := ctx.FormValue("limit")
	since := ctx.FormValue("since")
	sort := ctx.FormValue("sort")
	desc := ctx.FormValue("desc")

	result, err := helpers.ThreadGetPosts(slugOrID, limit, since, sort, desc)
	switch err {
	case nil:
		response(ctx, fasthttp.StatusOK, result) // 200
	case errors.ThreadNotFound:
		responseCustomError(ctx, fasthttp.StatusNotFound, errors.ThreadNotFound) // 404
	default:
		responseDefaultError(ctx, fasthttp.StatusInternalServerError, err) // 500
	}
}
