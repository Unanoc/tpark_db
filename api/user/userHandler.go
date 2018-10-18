package user

import (
	"encoding/json"
	"log"
	"tpark_db/errors"
	"tpark_db/helpers"
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

	result, err := helpers.UserCreateHelper(&user)
	switch err {
	case nil:
		ctx.SetStatusCode(fasthttp.StatusCreated) // 201
		buf, err := json.Marshal(user)
		if err != nil {
			log.Println(err)
		}
		ctx.Write(buf)
	case errors.UserIsExist:
		ctx.SetStatusCode(fasthttp.StatusConflict) // 409
		buf, err := json.Marshal(result)
		if err != nil {
			log.Println(err)
		}
		ctx.Write(buf)
	default:
		ctx.SetStatusCode(fasthttp.StatusInternalServerError) // 500
		ctx.Write([]byte(err.Error()))
	}
}

func UserGetOneHandler(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")

	nickname := ctx.UserValue("nickname").(string)

	result, err := helpers.UserGetOneHelper(nickname)
	switch err {
	case nil:
		ctx.SetStatusCode(fasthttp.StatusOK) // 200
	case errors.UserNotFound:
		ctx.SetStatusCode(fasthttp.StatusNotFound) // 404
	default:
		ctx.SetStatusCode(fasthttp.StatusInternalServerError) // 500
		ctx.Write([]byte(err.Error()))
	}

	buf, err := json.Marshal(result)
	if err != nil {
		log.Println(err)
	}
	ctx.Write(buf)
}

func UserUpdateHandler(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")

	user := models.User{}
	err := json.Unmarshal(ctx.PostBody(), &user)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest) // 400 Bad Request
		ctx.WriteString(err.Error())
		return
	}
	user.Nickname = ctx.UserValue("nickname").(string)

	result, err := helpers.UserUpdateHelper(&user)
}
