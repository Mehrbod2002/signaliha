package main

import (
	"log"

	"arz/routes"
)

func main() {
	router := routes.SetupRouter()

	err := router.Run(":8080")
	if err != nil {
		log.Fatal("Error starting server:", err.Error())
	}
}
