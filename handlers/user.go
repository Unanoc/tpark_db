package handlers

import (
	"fmt"
	"log"
	"tpark_db/errors"
	"tpark_db/helpers"
	"tpark_db/models"

	"github.com/valyala/fasthttp"
)

func UserCreateHandler(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")

	user := models.User{}
	err := user.UnmarshalJSON(ctx.PostBody())
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest) // 400 Bad Request
		ctx.SetBodyString(err.Error())
		return
	}
	user.Nickname = ctx.UserValue("nickname").(string) // getting value from URI (/api/user/{nickname}/profile")

	result, err := helpers.UserCreateHelper(&user)
	switch err {
	case nil:
		ctx.SetStatusCode(fasthttp.StatusCreated) // 201
		buf, _ := user.MarshalJSON()
		ctx.SetBody(buf)
	case errors.UserIsExist:
		ctx.SetStatusCode(fasthttp.StatusConflict) // 409
		buf, _ := result.MarshalJSON()
		ctx.SetBody(buf)
	default:
		ctx.SetStatusCode(fasthttp.StatusInternalServerError) // 500
		log.Println(err)
		ctx.SetBodyString(err.Error())
	}
}

func UserGetOneHandler(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")

	nickname := ctx.UserValue("nickname").(string)

	result, err := helpers.UserGetOneHelper(nickname)
	switch err {
	case nil:
		ctx.SetStatusCode(fasthttp.StatusOK) // 200
		buf, _ := result.MarshalJSON()
		ctx.SetBody(buf)
	case errors.UserNotFound:
		ctx.SetStatusCode(fasthttp.StatusNotFound) // 404
		errorResp := errors.Error{
			Message: fmt.Sprintf("Can't find user with nickname: %s", nickname),
		}
		buf, _ := errorResp.MarshalJSON()
		ctx.SetBody(buf)
	default:
		ctx.SetStatusCode(fasthttp.StatusInternalServerError) // 500
		ctx.SetBodyString(err.Error())
	}
}

func UserUpdateHandler(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")

	user := models.User{}
	err := user.UnmarshalJSON(ctx.PostBody())
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest) // 400 Bad Request
		ctx.SetBodyString(err.Error())
		return
	}
	user.Nickname = ctx.UserValue("nickname").(string)

	err = helpers.UserUpdateHelper(&user)
	switch err {
	case nil:
		ctx.SetStatusCode(fasthttp.StatusOK) // 200
		buf, _ := user.MarshalJSON()
		ctx.SetBody(buf)
	case errors.UserNotFound:
		ctx.SetStatusCode(fasthttp.StatusNotFound) // 404
		errorResp := errors.Error{
			Message: fmt.Sprintf("Can't find user with nickname: %s", user.Nickname),
		}
		buf, _ := errorResp.MarshalJSON()
		ctx.SetBody(buf)
	case errors.UserUpdateConflict:
		ctx.SetStatusCode(fasthttp.StatusConflict) // 409
		errorResp := errors.Error{
			Message: "New user profile data conflicts with existing users.",
		}
		buf, _ := errorResp.MarshalJSON()
		ctx.SetBody(buf)
	default:
		ctx.SetStatusCode(fasthttp.StatusInternalServerError) // 500
		ctx.SetBodyString(err.Error())
	}
}
