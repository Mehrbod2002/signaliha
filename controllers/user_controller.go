package controllers

import (
	"arz/utils"
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserController struct{}

func (uc *UserController) Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	if username == "" || password == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Username or password is empty",
		})
		return
	}

	db := &utils.Database{}
	err := db.Initialize()
	if err != nil {
		log.Fatal("Failed to initialize the database:", err)
	}
	defer db.Close()

	isValid, err := db.Login(username, password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Failed to authenticate user",
		})
		return
	}

	if isValid {
		session := sessions.Default(c)
		session.Set("user_id", uuid.New().String())
		session.Save()

		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": "Login successful",
		})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  http.StatusUnauthorized,
			"message": "Invalid username or password",
		})
	}
}
