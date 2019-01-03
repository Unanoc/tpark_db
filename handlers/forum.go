package handlers

import (
	"tpark_db/errors"
	"tpark_db/helpers"
	"tpark_db/models"

	"github.com/valyala/fasthttp"
)

// ForumCreateHandler handles POST request /forum/create.
func ForumCreateHandler(ctx *fasthttp.RequestCtx) {
	forum := models.Forum{}
	if err := forum.UnmarshalJSON(ctx.PostBody()); err != nil {
		responseDefaultError(ctx, fasthttp.StatusBadRequest, err) // 400
		return
	}

	result, err := helpers.ForumCreateHelper(&forum)
	switch err {
	case nil:
		response(ctx, fasthttp.StatusCreated, result) // 201
	case errors.UserNotFound:
		responseCustomError(ctx, fasthttp.StatusNotFound, errors.UserNotFound) // 404
	case errors.ForumIsExist:
		response(ctx, fasthttp.StatusConflict, result) // 409
	default:
		responseDefaultError(ctx, fasthttp.StatusInternalServerError, err) // 500
	}
}

// ForumCreateThreadHandler handles POST request /api/forum/:slug/create.
func ForumCreateThreadHandler(ctx *fasthttp.RequestCtx) {
	thread := models.Thread{}
	if err := thread.UnmarshalJSON(ctx.PostBody()); err != nil {
		responseDefaultError(ctx, fasthttp.StatusBadRequest, err) // 400
		return
	}

	forumSlug := ctx.UserValue("slug").(string)
	thread.Forum = forumSlug

	result, err := helpers.ForumCreateThreadHelper(&thread)
	switch err {
	case nil:
		response(ctx, fasthttp.StatusCreated, result) // 201
	case errors.ForumOrAuthorNotFound:
		responseCustomError(ctx, fasthttp.StatusNotFound, errors.ForumOrAuthorNotFound) // 404
	case errors.ThreadIsExist:
		response(ctx, fasthttp.StatusConflict, result) // 409
	default:
		responseDefaultError(ctx, fasthttp.StatusInternalServerError, err) // 500
	}
}

// ForumGetOneHandler handles GET request /api/forum/:slug/details.
func ForumGetOneHandler(ctx *fasthttp.RequestCtx) {
	slug := ctx.UserValue("slug")

	result, err := helpers.ForumGetBySlugHelper(slug.(string))
	switch err {
	case nil:
		response(ctx, fasthttp.StatusOK, result) // 200
	case errors.ForumNotFound:
		responseCustomError(ctx, fasthttp.StatusNotFound, errors.ForumNotFound) // 404
	default:
		responseDefaultError(ctx, fasthttp.StatusInternalServerError, err) // 500
	}
}

// ForumGetThreadsHandler handles GET request /api/forum/:slug/threads.
func ForumGetThreadsHandler(ctx *fasthttp.RequestCtx) {
	slug := ctx.UserValue("slug").(string)
	limit := ctx.FormValue("limit")
	since := ctx.FormValue("since")
	desc := ctx.FormValue("desc")

	result, err := helpers.ForumGetThreadsHelper(slug, limit, since, desc)
	switch err {
	case nil:
		response(ctx, fasthttp.StatusOK, result) // 200
	case errors.ForumNotFound:
		responseCustomError(ctx, fasthttp.StatusNotFound, errors.ForumNotFound) // 404
	default:
		responseDefaultError(ctx, fasthttp.StatusInternalServerError, err) // 500
	}
}

// ForumGetUsersHandler handles GET request /api/forum/:slug/users.
func ForumGetUsersHandler(ctx *fasthttp.RequestCtx) {
	slug := ctx.UserValue("slug").(string)
	limit := ctx.FormValue("limit")
	since := ctx.FormValue("since")
	desc := ctx.FormValue("desc")

	result, err := helpers.ForumGetUsersHelper(slug, limit, since, desc)
	switch err {
	case nil:
		response(ctx, fasthttp.StatusOK, result) // 200
	case errors.ForumNotFound:
		responseCustomError(ctx, fasthttp.StatusNotFound, errors.ForumNotFound) // 404
	default:
		responseDefaultError(ctx, fasthttp.StatusInternalServerError, err) // 500
	}
}
