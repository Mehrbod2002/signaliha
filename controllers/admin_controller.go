package controllers

import (
	"arz/models"
	"arz/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AdminController struct {
	// Define any additional fields or dependencies needed by the admin controllers
}

// CreateToken handles the creation of a token
func (ac *AdminController) CreateToken(c *gin.Context) {
	var token models.Token
	if err := c.ShouldBindJSON(&token); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call the utility function to create the token
	if err := utils.CreateToken(&token); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Token created successfully"})
}

// GetTokens retrieves all tokens
func (ac *AdminController) GetTokens(c *gin.Context) {
	tokens, err := utils.GetTokens()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tokens)
}

// GetTokenHistory retrieves the history of a specific token
func (ac *AdminController) GetTokenHistory(c *gin.Context) {
	token := c.Param("token")

	history, err := utils.GetTokenHistory(token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, history)
}
