package handlers

import (
	"tpark_db/errors"
	"tpark_db/helpers"
	"tpark_db/models"

	"github.com/valyala/fasthttp"
)

// UserCreateHandler handles POST request /api/user/:nickname/create.
func UserCreateHandler(ctx *fasthttp.RequestCtx) {
	user := models.User{}
	user.Nickname = ctx.UserValue("nickname").(string)
	if err := user.UnmarshalJSON(ctx.PostBody()); err != nil {
		responseDefaultError(ctx, fasthttp.StatusBadRequest, err) // 400
		return
	}

	result, err := helpers.UserCreateHelper(&user)
	switch err {
	case nil:
		response(ctx, fasthttp.StatusCreated, user) // 201
	case errors.UserIsExist:
		response(ctx, fasthttp.StatusConflict, result) // 409
	default:
		responseDefaultError(ctx, fasthttp.StatusInternalServerError, err) // 500
	}
}

// UserUpdateHandler handles POST request /api/user/:nickname/profile.
func UserUpdateHandler(ctx *fasthttp.RequestCtx) {
	user := models.User{}
	user.Nickname = ctx.UserValue("nickname").(string)
	if err := user.UnmarshalJSON(ctx.PostBody()); err != nil {
		responseDefaultError(ctx, fasthttp.StatusBadRequest, err) // 400
		return
	}

	err := helpers.UserUpdateHelper(&user)
	switch err {
	case nil:
		response(ctx, fasthttp.StatusOK, user) // 200
	case errors.UserNotFound:
		responseCustomError(ctx, fasthttp.StatusNotFound, errors.UserNotFound) // 404
	case errors.UserUpdateConflict:
		responseCustomError(ctx, fasthttp.StatusConflict, errors.UserUpdateConflict) // 404
	default:
		responseDefaultError(ctx, fasthttp.StatusInternalServerError, err) // 500
	}
}

// UserGetOneHandler handles GET request /api/user/:nickname/profile.
func UserGetOneHandler(ctx *fasthttp.RequestCtx) {
	nickname := ctx.UserValue("nickname").(string)

	result, err := helpers.UserGetOneHelper(nickname)
	switch err {
	case nil:
		response(ctx, fasthttp.StatusOK, result) // 200
	case errors.UserNotFound:
		responseCustomError(ctx, fasthttp.StatusNotFound, errors.UserNotFound) // 404
	default:
		responseDefaultError(ctx, fasthttp.StatusInternalServerError, err) // 500
	}
}
