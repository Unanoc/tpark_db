package main

import (
	"log"
	
	"github.com/valyala/fasthttp"

	"tpark_db/router"
	"tpark_db/logger"
)

func main() {	
	router := router.NewRouter()
	log.Fatal(fasthttp.ListenAndServe(":5000", logger.Logger(router.Handler)))
}
