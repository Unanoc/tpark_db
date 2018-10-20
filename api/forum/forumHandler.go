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
		ctx.WriteString(err.Error())
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
		ctx.Write(buf)
	default:
		ctx.SetStatusCode(fasthttp.StatusInternalServerError) // 500
		ctx.SetBody([]byte(err.Error()))
	}
}

func ForumGetOneHandler(ctx *fasthttp.RequestCtx) {

}

func ForumGetThreadsHandler(ctx *fasthttp.RequestCtx) {

}

func ForumGetUsersHandler(ctx *fasthttp.RequestCtx) {

}
