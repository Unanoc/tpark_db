package main

import (
	// "tpark_db/database"
	"log"
	"tpark_db/logger"
	"tpark_db/router"

	"github.com/valyala/fasthttp"
)

func main() {
	// db := database.DBConn
	// err := db.Connect()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// if err = db.CreateDB("sql/create_tables.sql"); err != nil {
	// 	log.Fatal(err)
	// }
	// defer db.Disconnect()

	router := router.NewRouter()
	log.Fatal(fasthttp.ListenAndServe(":5001", logger.Logger(router.Handler)))
}
