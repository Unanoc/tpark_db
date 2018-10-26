package thread

import (
	"fmt"
	"tpark_db/errors"
	"tpark_db/helpers"
	"tpark_db/models"

	"github.com/valyala/fasthttp"
)

func ThreadCreateHandler(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")

	slugOrID := ctx.UserValue("slug_or_id").(string)
	posts := models.Posts{}
	err := posts.UnmarshalJSON(ctx.PostBody())
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest) // 400 Bad Request
		ctx.SetBodyString(err.Error())
		return
	}
	result, err := helpers.ThreadCreateHelper(&posts, slugOrID)

	switch err {
	case nil:
		ctx.SetStatusCode(fasthttp.StatusCreated) // 201
		buf, _ := result.MarshalJSON()
		ctx.SetBody(buf)
	case errors.NoPostsForCreate:
		ctx.SetStatusCode(fasthttp.StatusCreated) // 201
		ctx.SetBody(ctx.PostBody())
	case errors.ThreadNotFound:
		ctx.SetStatusCode(fasthttp.StatusNotFound) // 404
		errorResp := errors.Error{
			Message: fmt.Sprintf("Can't find thread by slug or id: %s", slugOrID),
		}
		buf, _ := errorResp.MarshalJSON()
		ctx.SetBody(buf)
	case errors.PostParentNotFound:
		ctx.SetStatusCode(fasthttp.StatusConflict)
		errorResp := errors.Error{
			Message: fmt.Sprintf("Can't find parent of post"),
		}
		buf, _ := errorResp.MarshalJSON()
		ctx.SetBody(buf)
	default:
		ctx.SetStatusCode(fasthttp.StatusInternalServerError) // 500
		ctx.SetBodyString(err.Error())
	}
}

func ThreadGetOneHandler(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")

	slugOrID := ctx.UserValue("slug_or_id").(string)

	result, err := helpers.GetThreadBySlugOrId(slugOrID)
	switch err {
	case nil:
		ctx.SetStatusCode(fasthttp.StatusOK) // 200
		buf, _ := result.MarshalJSON()
		ctx.SetBody(buf)
	case errors.ThreadNotFound:
		ctx.SetStatusCode(fasthttp.StatusNotFound) // 404
		errorResp := errors.Error{
			Message: fmt.Sprintf("Can't find thread by slug or id: %s", slugOrID),
		}
		buf, _ := errorResp.MarshalJSON()
		ctx.SetBody(buf)
	default:
		ctx.SetStatusCode(fasthttp.StatusInternalServerError) // 500
		ctx.SetBodyString(err.Error())
	}
}

func ThreadGetPosts(ctx *fasthttp.RequestCtx) {

}

func ThreadUpdateHandler(ctx *fasthttp.RequestCtx) {

}

func ThreadVoteHandler(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")

	slug_or_id := ctx.UserValue("slug_or_id").(string)

	vote := models.Vote{}
	err := vote.UnmarshalJSON(ctx.PostBody())
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest) // 400 Bad Request
		ctx.SetBodyString(err.Error())
		return
	}

	result, err := helpers.ThreadVoteHelper(&vote, slug_or_id)
	switch err {
	case nil:
		ctx.SetStatusCode(fasthttp.StatusOK) // 200
		buf, _ := result.MarshalJSON()
		ctx.SetBody(buf)
	case errors.ThreadNotFound:
		ctx.SetStatusCode(fasthttp.StatusNotFound) // 404
		errorResp := errors.Error{
			Message: fmt.Sprintf("Can't find thread by slug or id: %s", slug_or_id),
		}
		buf, _ := errorResp.MarshalJSON()
		ctx.SetBody(buf)
	default:
		ctx.SetStatusCode(fasthttp.StatusInternalServerError) // 500
		ctx.SetBodyString(err.Error())
	}
}
