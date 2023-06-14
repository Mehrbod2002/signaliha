package routes

import (
	"arz/controllers"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func LimitHandler(c *gin.Context) {
	c.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests"})
	c.Abort()
}

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// limiter := tollbooth.NewLimiter(10, nil)
	// r.Use(tollbooth_gin.LimitHandler(limiter))
	// r.Use(tollbooth_gin.TollboothMiddleware(limiter))

	// Public routes
	userController := &controllers.UserController{}
	adminController := &controllers.AdminController{}
	r.POST("/login", userController.Login)

	admin := r.Group("/admin")
	admin.Use(AdminAuthMiddleware()) // Middleware for admin authentication
	{
		admin.POST("/tokens", adminController.CreateToken)
		admin.GET("/tokens", adminController.GetTokens)
		admin.GET("/tokens/:token/history", adminController.GetTokenHistory)
	}

	// r.NoRoute(LimitHandler) // Uncomment this line if LimitHandler is defined

	return r
}

func AdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		userID := session.Get("user_id")

		if userID == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		c.Next()
	}
}
