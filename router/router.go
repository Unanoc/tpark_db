package router

import (
	"github.com/buaazp/fasthttprouter"

	"github.com/Unanoc/tpark_db/handlers"
)

var Router = fasthttprouter.New()

func Init() {
	Router.GET("/", handlers.Index)
	Router.GET("/hello/:name", handlers.Hello)
}
