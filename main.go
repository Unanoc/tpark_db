package main

import (
	"log"

	"tpark_db/router"
	"tpark_db/logger"
	
	"github.com/valyala/fasthttp"
)

func main() {	
	router := router.NewRouter()
	log.Fatal(fasthttp.ListenAndServe(":5000", logger.Logger(router.Handler)))
}
