package main

import (
	"log"
	"tpark_db/database"
	"tpark_db/logger"
	"tpark_db/router"

	"github.com/valyala/fasthttp"
)

func main() {
	db := database.DBConn

	if err := db.Connect(); err != nil {
		log.Fatal(err)
	}

	defer db.Disconnect()

	if err := db.CreateDB("sql/create_tables.sql"); err != nil {
		log.Fatal(err)
	}

	router := router.NewRouter()
	log.Fatal(fasthttp.ListenAndServe(":5000", logger.Logger(router.Handler)))
}
