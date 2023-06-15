package controllers

import (
	"arz/utils"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type UserController struct{}

const (
	session_name    = "secret"
	tokenExpiration = time.Hour * 24
)

type Claims struct {
	UserID string `json:"user_id"`
	jwt.StandardClaims
}

func (uc *UserController) Login(c *gin.Context) {
	session := sessions.Default(c)
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
		token, err := generateToken(username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  http.StatusInternalServerError,
				"message": "Failed to generate token",
			})
			return
		}
		session.Set("session", token)
		session.Save()
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": "Login successful",
			"token":   token,
		})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  http.StatusUnauthorized,
			"message": "Invalid username or password",
		})
	}
}

func (uc *UserController) GetMessage(c *gin.Context) {
	db := &utils.Database{}
	err := db.Initialize()
	if err != nil {
		log.Fatal("Failed to initialize the database:", err)
	}
	defer db.Close()
	token := c.PostForm("token")

	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "token is empty",
		})
		return
	}

	isValid, err := db.ValidateToken(token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Failed to validate token",
		})
		return
	}

	if isValid {
		lastID, err := db.GetLastID(token)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  http.StatusInternalServerError,
				"message": "Failed to retrieve last ID",
			})
			return
		}

		if lastID == nil {
			message, err := db.GetLastMessage()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"status":  http.StatusInternalServerError,
					"message": "Failed to retrieve message",
				})
				return
			}

			err = db.UpdateLastID(token, message.ID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"status":  http.StatusInternalServerError,
					"message": "Failed to update last ID",
				})
				return
			}

			err = db.AddHistory(token, []string{message.Coin}, "GetLastMessage")
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"status":  http.StatusInternalServerError,
					"message": "Failed to add history",
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"status":  http.StatusOK,
				"message": "last ID updated successfully",
				"data":    message,
			})
		} else {
			messages, err := db.GetMessagesAfterID(*lastID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"status":  http.StatusInternalServerError,
					"message": "Failed to retrieve messages",
				})
				return
			}

			if len(messages) > 0 {
				lastMessage := messages[len(messages)-1]

				err = db.UpdateLastID(token, lastMessage.MessageID)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"status":  http.StatusInternalServerError,
						"message": "Failed to update last ID",
					})
					return
				}
				coins := make([]string, len(messages))
				for i, msg := range messages {
					coins[i] = msg.Coin
				}

				err = db.AddHistory(token, coins, "GetLastMessage")
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"status":  http.StatusInternalServerError,
						"message": "Failed to add history",
					})
					return
				}

				c.JSON(http.StatusOK, gin.H{
					"status":  http.StatusOK,
					"message": "Token validated, messages retrieved, and last ID updated successfully",
					"data":    messages,
				})
			} else {
				c.JSON(http.StatusOK, gin.H{
					"status":  http.StatusOK,
					"message": "No new messages found",
					"data":    nil,
				})
			}
		}
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  http.StatusUnauthorized,
			"message": "Invalid token",
		})
	}
}

func (uc *UserController) GetTokenByMessageID(c *gin.Context) {
	db := &utils.Database{}
	err := db.Initialize()
	if err != nil {
		log.Fatal("Failed to initialize the database:", err)
	}
	defer db.Close()
	token := c.PostForm("token")

	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "token is empty",
		})
		return
	}

	isValid, err := db.ValidateToken(token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Failed to validate token",
		})
		return
	}

	if isValid {
		if err != nil {
			log.Fatal("Failed to initialize the database:", err)
		}
		defer db.Close()

		messageID := c.Param("messageID")

		token, err := db.GetTokenByMessageID(messageID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  http.StatusInternalServerError,
				"message": "Failed to retrieve token",
			})
			return
		}

		if token != "" {
			c.JSON(http.StatusOK, gin.H{
				"status":  http.StatusOK,
				"message": "Token retrieved successfully",
				"data":    token,
			})
		} else {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  http.StatusNotFound,
				"message": "Token not found",
			})
		}
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  http.StatusUnauthorized,
			"message": "Invalid token",
		})
	}
}

func generateToken(userID string) (string, error) {
	claims := &Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenExpiration).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(session_name))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func (uc *UserController) Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Logout successful",
	})
}
