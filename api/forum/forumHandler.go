package forum

import (
	"encoding/json"
	"fmt"
	"log"
	"tpark_db/errors"
	"tpark_db/helpers"
	"tpark_db/models"

	"github.com/valyala/fasthttp"
)

func ForumCreateHandler(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")

	// Unmarshalling JSON from POST request
	forum := models.Forum{}
	err := json.Unmarshal(ctx.PostBody(), &forum)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest) // 400 Bad Request
		ctx.WriteString(err.Error())
		return
	}

	result, err := helpers.ForumCreateHelper(&forum)
	switch err {
	case nil:
		ctx.SetStatusCode(fasthttp.StatusCreated) // 201
		buf, err := json.Marshal(result)
		if err != nil {
			log.Println(err)
		}
		ctx.SetBody(buf)
	case errors.UserNotFound:
		ctx.SetStatusCode(fasthttp.StatusNotFound) // 404

		errorResp := errors.Error{
			Message: fmt.Sprintf("Can't find user with nickname: %s", forum.User)}

		buf, err := json.Marshal(errorResp)
		if err != nil {
			log.Println(err)
		}
		ctx.SetBody(buf)
	case errors.ForumIsExist:
		ctx.SetStatusCode(fasthttp.StatusConflict) // 409
		buf, err := json.Marshal(result)
		if err != nil {
			log.Println(err)
		}
		ctx.Write(buf)
	default:
		ctx.SetStatusCode(fasthttp.StatusInternalServerError) // 500
		log.Println(err)
		ctx.SetBody([]byte(err.Error()))
	}
}

func ForumGetOneHandler(ctx *fasthttp.RequestCtx) {

}

func ForumGetThreadsHandler(ctx *fasthttp.RequestCtx) {

}

func ForumGetUsersHandler(ctx *fasthttp.RequestCtx) {

}
