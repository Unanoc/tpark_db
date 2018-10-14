package main

import (
	"log"
	"tpark_db/database"
	"tpark_db/logger"
	"tpark_db/router"
	"go.uber.org/zap"
	"github.com/valyala/fasthttp"
)

const (
	port = ":5002"
	addr = "localhost"
)

func main() {
	// creating a connection with DB and loading of schema
	db := database.DBConn
	if err := db.Connect(); err != nil {
		logger.Logger.Error(err.Error())
	}
	if err := db.ExecSqlScript("sql/create_tables.sql"); err != nil {
		logger.Logger.Error(err.Error())
	}
	defer db.Disconnect()
	defer logger.Logger.Sync() // calling Sync before letting your process exit

	logger.Logger.Info("Starting server",
		zap.String("host", addr),
		zap.String("port", port),
	)
	// creating a router and starting of server
	router := router.NewRouter()
	log.Fatal(fasthttp.ListenAndServe(port, logger.LoggerHandler(router.Handler)))
}
