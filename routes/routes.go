package routes

import (
	"arz/controllers"
	"net/http"

	"github.com/didip/tollbooth/v6"
	"github.com/didip/tollbooth_gin/v3"
	"github.com/gin-gonic/gin"
)

func LimitHandler(c *gin.Context) {
	c.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests"})
	c.Abort()
}

func SetupRouter() *gin.Engine {
	r := gin.Default()

	limiter := tollbooth.NewLimiter(10, nil)
	r.Use(tollbooth_gin.LimitHandler(limiter))
	r.Use(tollbooth_gin.TollboothMiddleware(limiter))

	admin := r.Group("/admin")
	admin.Use(AdminAuthMiddleware()) // Middleware for admin authentication
	{
		admin.POST("/tokens", controllers.CreateToken)
		admin.GET("/tokens", controllers.GetTokens)
		admin.GET("/tokens/:token/history", controllers.GetTokenHistory)
	}

	r.NoRoute(LimitHandler)
	return r
}

func AdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Implement your admin authentication logic here
		// You can check if the user is authenticated as an admin
		// If not, you can return an error or redirect to a login page
		// Example:
		// if !isAdminAuthenticated(c) {
		// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		// 	c.Abort()
		// 	return
		// }

		c.Next()
	}
}
