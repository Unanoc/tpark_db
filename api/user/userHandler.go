package user

import (
	"encoding/json"
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
	err := json.Unmarshal(ctx.PostBody(), &user)
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
		buf, err := json.Marshal(user)
		if err != nil {
			log.Println(err)
		}
		ctx.SetBody(buf)
	case errors.UserIsExist:
		ctx.SetStatusCode(fasthttp.StatusConflict) // 409
		buf, err := json.Marshal(result)
		if err != nil {
			log.Println(err)
		}
		ctx.SetBody(buf)
	default:
		ctx.SetStatusCode(fasthttp.StatusInternalServerError) // 500
		errMsg := err.Error()
		log.Println(err)
		ctx.SetBody([]byte(errMsg))
	}
}

func UserGetOneHandler(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")

	nickname := ctx.UserValue("nickname").(string)

	result, err := helpers.UserGetOneHelper(nickname)
	switch err {
	case nil:
		ctx.SetStatusCode(fasthttp.StatusOK) // 200
		buf, err := json.Marshal(result)
		if err != nil {
			log.Println(err)
		}
		ctx.SetBody(buf)
	case errors.UserNotFound:
		ctx.SetStatusCode(fasthttp.StatusNotFound) // 404

		errorResp := errors.Error{
			Message: fmt.Sprintf("Can't find user with nickname: %s", nickname)}

		buf, err := json.Marshal(errorResp)
		if err != nil {
			log.Println(err)
		}
		ctx.SetBody(buf)
	default:
		ctx.SetStatusCode(fasthttp.StatusInternalServerError) // 500
		errMsg := err.Error()
		log.Println(err)
		ctx.SetBody([]byte(errMsg))
	}
}

func UserUpdateHandler(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")

	user := models.User{}
	err := json.Unmarshal(ctx.PostBody(), &user)
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
		buf, err := json.Marshal(user)
		if err != nil {
			log.Println(err)
		}
		ctx.SetBody(buf)
	case errors.UserNotFound:
		ctx.SetStatusCode(fasthttp.StatusNotFound) // 404
		errorResp := errors.Error{Message: fmt.Sprintf("Can't find user with nickname: %s", user.Nickname)}

		buf, err := json.Marshal(errorResp)
		if err != nil {
			log.Println(err)
		}
		ctx.SetBody(buf)
	case errors.UserUpdateConflict:
		ctx.SetStatusCode(fasthttp.StatusConflict) // 409

		errorResp := errors.Error{Message: "New user profile data conflicts with existing users."}
		buf, err := json.Marshal(errorResp)
		if err != nil {
			log.Println(err)
		}
		ctx.SetBody(buf)
	default:
		ctx.SetStatusCode(fasthttp.StatusInternalServerError) // 500
		errMsg := err.Error()
		log.Println(err)
		ctx.SetBody([]byte(errMsg))
	}
}
