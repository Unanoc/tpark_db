package forum

import (
	"github.com/valyala/fasthttp"
)

func ForumCreateHandler(ctx *fasthttp.RequestCtx) {
	body := ctx.PostBody()
	ctx.SetBody(body)
	// ctx.SetContentType("application/json")

	// // Unmarshalling JSON from POST request
	// var forum models.Forum
	// err := json.Unmarshal(body, &forum)
	// if err != nil {
	// 	ctx.SetStatusCode(fasthttp.StatusBadRequest) // 400 Bad Request
	// 	ctx.WriteString(err.Error())
	// 	return
	// }

	// result, err := database.CreteOrGetExistingForum(&forum)

	// switch err {
	// case nil:

	// }
}

func ForumGetOneHandler(ctx *fasthttp.RequestCtx) {

}

func ForumGetThreadsHandler(ctx *fasthttp.RequestCtx) {

}

func ForumGetUsersHandler(ctx *fasthttp.RequestCtx) {

}
