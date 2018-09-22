package router

import (
	"tpark_db/handlers"

	"github.com/buaazp/fasthttprouter"
)

var Router = fasthttprouter.New()

func Init() {
	Router.GET("/", handlers.Index)
	Router.GET("/hello/:name", handlers.Hello)
	return
}