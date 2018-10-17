package forum

import (
	"encoding/json"
	"log"
	"tpark_db/errors"
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

	result, err := forum.CreateForum()
	switch err {
	case nil:
		ctx.SetStatusCode(201)
		buf, err := json.Marshal(result)
		if err != nil {
			log.Println(err)
		}
		ctx.Write(buf)
	case errors.UserNotFound:
		ctx.SetStatusCode(404)
		ctx.Write([]byte(err.Error()))
	case errors.ForumIsExist:
		ctx.SetStatusCode(409)
		buf, err := json.Marshal(result)
		if err != nil {
			log.Println(err)
		}
		ctx.Write(buf)
	}
}

func ForumGetOneHandler(ctx *fasthttp.RequestCtx) {

}

func ForumGetThreadsHandler(ctx *fasthttp.RequestCtx) {

}

func ForumGetUsersHandler(ctx *fasthttp.RequestCtx) {

}
