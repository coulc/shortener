package main

import (
	"log"
	"net/http"
	"shortener/utils"
)

func main() {
	file := utils.InitLogger()
	defer file.Close()

	shortenerHandler,cleanUp := dbInit()
	defer cleanUp()

	routesInit(shortenerHandler)

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080",nil))
}
