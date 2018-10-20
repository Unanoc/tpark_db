package forum

import (
	"fmt"
	"tpark_db/errors"
	"tpark_db/helpers"
	"tpark_db/models"

	"github.com/valyala/fasthttp"
)

func ForumCreateHandler(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")

	forum := models.Forum{}
	err := forum.UnmarshalJSON(ctx.PostBody())
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest) // 400 Bad Request
		ctx.SetBodyString(err.Error())
		return
	}
	result, err := helpers.ForumCreateHelper(&forum)

	switch err {
	case nil:
		ctx.SetStatusCode(fasthttp.StatusCreated) // 201
		buf, _ := result.MarshalJSON()
		ctx.SetBody(buf)
	case errors.UserNotFound:
		ctx.SetStatusCode(fasthttp.StatusNotFound) // 404
		errorResp := errors.Error{
			Message: fmt.Sprintf("Can't find user with nickname: %s", forum.User),
		}
		buf, _ := errorResp.MarshalJSON()
		ctx.SetBody(buf)
	case errors.ForumIsExist:
		ctx.SetStatusCode(fasthttp.StatusConflict) // 409
		buf, _ := result.MarshalJSON()
		ctx.SetBody(buf)
	default:
		ctx.SetStatusCode(fasthttp.StatusInternalServerError) // 500
		ctx.SetBody([]byte(err.Error()))
	}
}

func ForumGetOneHandler(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")

	slug := ctx.UserValue("slug")
	forum, err := helpers.ForumGetBySlug(slug.(string))

	switch err {
	case nil:
		ctx.SetStatusCode(fasthttp.StatusOK) // 200
		buf, _ := forum.MarshalJSON()
		ctx.SetBody(buf)
	case errors.ForumNotFound:
		ctx.SetStatusCode(fasthttp.StatusNotFound) // 404
		errorResp := errors.Error{
			Message: fmt.Sprintf("Can't find forum with slug: %s", slug),
		}
		buf, _ := errorResp.MarshalJSON()
		ctx.SetBody(buf)
	default:
		ctx.SetStatusCode(fasthttp.StatusInternalServerError) // 500
		ctx.SetBody([]byte(err.Error()))
	}
}

func ForumCreateThreadHandler(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")

	thread := models.Thread{}
	err := thread.UnmarshalJSON(ctx.PostBody())
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest) // 400 Bad Request
		ctx.SetBodyString(err.Error())
		return
	}
	slug := ctx.UserValue("slug").(string)
	err = helpers.ForumCreateThreadHelper(&thread)

	switch err {
	case nil:
		ctx.SetStatusCode(fasthttp.StatusCreated) // 201
		buf, _ := thread.MarshalJSON()
		ctx.SetBody(buf)
	case errors.ForumOrAuthorNotFound:
		ctx.SetStatusCode(fasthttp.StatusNotFound) // 404
		errorResp := errors.Error{
			Message: fmt.Sprintf("Can't find forum by slug: %s", slug),
		}
		buf, _ := errorResp.MarshalJSON()
		ctx.SetBody(buf)
	case errors.ThreadIsExist:
		ctx.SetStatusCode(fasthttp.StatusConflict) // 409
		buf, _ := thread.MarshalJSON()
		ctx.SetBody(buf)
	default:
		ctx.SetStatusCode(fasthttp.StatusInternalServerError) // 500
		ctx.SetBody([]byte(err.Error()))
	}
}

func ForumGetThreadsHandler(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")
	slug := ctx.UserValue("slug").(string)
	limit := ctx.FormValue("limit")
	since := ctx.FormValue("since")
	desc := ctx.FormValue("desc")
	threads, err := helpers.ForumGetThreadsHelper(slug, limit, since, desc)

	switch err {
	case nil:
		ctx.SetStatusCode(fasthttp.StatusOK) // 200
		buf, _ := threads.MarshalJSON()
		ctx.SetBody(buf)
	case errors.ForumNotFound:
		ctx.SetStatusCode(fasthttp.StatusNotFound) // 404
		errorResp := errors.Error{
			Message: fmt.Sprintf("Can't find forum by slug: %s", slug),
		}
		buf, _ := errorResp.MarshalJSON()
		ctx.SetBody(buf)
	default:
		ctx.SetStatusCode(fasthttp.StatusInternalServerError) // 500
		ctx.SetBody([]byte(err.Error()))
	}
}

func ForumGetUsersHandler(ctx *fasthttp.RequestCtx) {

}
