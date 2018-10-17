package user

import (
	"encoding/json"
	"log"
	"tpark_db/errors"
	"tpark_db/models"

	"github.com/valyala/fasthttp"
)

func UserCreateHandler(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")

	user := models.User{}
	err := json.Unmarshal(ctx.PostBody(), &user)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest) // 400 Bad Request
		ctx.WriteString(err.Error())
		return
	}
	user.Nickname = ctx.UserValue("nickname").(string) // getting value from URI (/api/user/{nickname}/profile")

	result, err := user.CreateUser()
	switch err {
	case nil:
		ctx.SetStatusCode(201)
		buf, err := json.Marshal(user)
		if err != nil {
			log.Println(err)
		}
		ctx.Write(buf)
	case errors.UserIsExist:
		ctx.SetStatusCode(409)
		buf, err := json.Marshal(result)
		if err != nil {
			log.Println(err)
		}
		ctx.Write(buf)
	default:
		ctx.SetStatusCode(500)
		ctx.Write([]byte(err.Error()))
	}
}

func UserGetOneHandler(ctx *fasthttp.RequestCtx) {

}

func UserUpdateHandler(ctx *fasthttp.RequestCtx) {

}
