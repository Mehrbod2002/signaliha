package arz

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/your-app/controllers"
	"github.com/your-app/database"
)

func main() {
	err := godotenv.Load("local.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db := database.Database{}
	err = db.Initialize()
	if err != nil {
		log.Fatal("Error initializing database:", err.Error())
	}
	defer db.Close()

	router := gin.Default()

	adminController := controllers.NewAdminController(&db)

	router.POST("/admin/token", adminController.CreateToken)
	router.GET("/admin/history/:token", adminController.GetTokenHistory)

	err = router.Run(":8080")
	if err != nil {
		log.Fatal("Error starting server:", err.Error())
	}
}
