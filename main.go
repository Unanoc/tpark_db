package main

import (
	"log"
	"tpark_db/database"
	"tpark_db/logger"
	"tpark_db/router"

	"github.com/valyala/fasthttp"
)

func main() {
	db := database.Database{}
	err := db.Connect()
	if err != nil {
		log.Fatal(err)
	}
	if err = db.CreateDB(); err != nil {
		log.Fatal(err)
	}
	// db.Disconnect()

	router := router.NewRouter()
	log.Fatal(fasthttp.ListenAndServe(":5000", logger.Logger(router.Handler)))
}
