package routes

import (
	"arz/controllers"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	store := cookie.NewStore([]byte("session"))
	store.Options(sessions.Options{
		Path:     "/",
		MaxAge:   86400,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
	})
	r.Use(sessions.Sessions("session", store))

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://127.0.0.1:3000"},
		AllowMethods:     []string{"PUT", "PATCH", "GET", "POST", "HEAD"},
		AllowHeaders:     []string{"Origin", "Cookie", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "http://127.0.0.1:3000"
		},
	}))

	userController := &controllers.UserController{}
	adminController := &controllers.AdminController{}
	r.POST("/login", userController.Login)
	r.POST("/logout", userController.Logout)
	r.POST("/signal", userController.GetTokenByMessageID)
	r.POST("/last_signals", userController.GetMessage)

	admin := r.Group("/admin")
	admin.Use(AdminAuthMiddleware())
	{
		admin.POST("/tokens", adminController.CreateToken)
		admin.GET("/tokens", adminController.GetTokens)
		admin.POST("/tokens/delete", adminController.DeleteTokens)
		admin.GET("/tokens/:token/history", adminController.GetTokenHistory)
	}

	return r
}

func AdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		token := session.Get("session")
		tokenString := c.GetHeader("Authorization")

		if token == nil && tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  http.StatusUnauthorized,
				"message": "Unauthorized",
			})
			c.Abort()
			return
		}
		if token == nil {
			token = tokenString
		}

		claims := jwt.MapClaims{}
		_, err := jwt.ParseWithClaims(token.(string), claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("secret"), nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  http.StatusUnauthorized,
				"message": "Invalid token",
			})
			c.Abort()
			return
		}

		if claims["user_id"] != "admin" {
			c.JSON(http.StatusForbidden, gin.H{
				"status":  http.StatusForbidden,
				"message": "Forbidden",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
