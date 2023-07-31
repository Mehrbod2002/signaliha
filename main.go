package main

import (
	"log"

	"arz/routes"
)

func main() {
	router := routes.SetupRouter()

	err := router.Run(":80")
	if err != nil {
		log.Fatal("Error starting server:", err.Error())
	}
}
