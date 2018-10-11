package main

import (
	"log"
	"net/http"
	"tpark_db/router"
)

func main() {	
	router := router.NewRouter()
	log.Fatal(http.ListenAndServe(":8080", router))
}
