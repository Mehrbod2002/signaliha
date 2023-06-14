package controllers

import (
	"arz/utils"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AdminController struct {
	db *utils.Database
}

func NewAdminController(db *utils.Database) *AdminController {
	return &AdminController{
		db: db,
	}
}

func (ac *AdminController) CreateToken(c *gin.Context) {
	db := &utils.Database{}

	err := db.Initialize()
	if err != nil {
		log.Fatal("Failed to initialize the database:", err)
	}

	defer db.Close()
	name := c.PostForm("name")
	limit := c.PostForm("limit")

	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid name"})
		return
	}

	var limitValue int
	if limit == "" || limit == "-1" {
		limitValue = -1
	} else {
		limitValue, err = strconv.Atoi(limit)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit"})
			return
		}
	}

	token, err := db.CreateToken(name, limitValue)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK, "token": token, "message": "Token created successfully"})
}

func (ac *AdminController) GetTokens(c *gin.Context) {
	db := &utils.Database{}

	err := db.Initialize()
	if err != nil {
		log.Fatal("Failed to initialize the database:", err)
	}

	defer db.Close()
	tokens, err := db.GetTokens()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tokens)
}

func (ac *AdminController) GetTokenHistory(c *gin.Context) {
	db := &utils.Database{}

	err := db.Initialize()
	if err != nil {
		log.Fatal("Failed to initialize the database:", err)
	}

	defer db.Close()

	token := c.Param("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "error": "Invalid token"})
		return
	}

	history, err := db.GetTokenHistory(token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, history)
}
