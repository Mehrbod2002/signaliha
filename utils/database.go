package utils

import (
	"arz/models"
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

type Database struct {
	db *sql.DB
}

func (d *Database) Initialize() error {
	err := godotenv.Load("local.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	dbURI := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPassword, dbHost, dbPort, dbName)

	db, err := sql.Open("mysql", dbURI)
	if err != nil {
		return err
	}

	err = db.Ping()
	if err != nil {
		return err
	}

	d.db = db
	return nil
}

func (d *Database) Login(username string, password string) (bool, error) {
	query := "SELECT COUNT(*) FROM users WHERE username = ? AND password = ?"

	var count int
	err := d.db.QueryRow(query, username, password).Scan(&count)
	if err != nil {
		return false, err
	}

	if count > 0 {
		return true, nil
	} else {
		return false, nil
	}
}

func (d *Database) CreateToken(name string, limit int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const tokenLength = 10

	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)

	token := make([]byte, tokenLength)
	for i := 0; i < tokenLength; i++ {
		token[i] = charset[random.Intn(len(charset))]
	}

	var limitValue interface{}
	if limit == -1 {
		limitValue = nil
	} else {
		limitValue = limit
	}

	query := "INSERT INTO users (name, token, `limit`) VALUES (?, ?, ?)"
	_, err := d.db.Exec(query, name, string(token), limitValue)
	if err != nil {
		return "", err
	}

	return string(token), nil
}

func (d *Database) GetTokens() ([]models.Token, error) {
	query := "SELECT token, name, `limit`, timestamp FROM tokens"
	rows, err := d.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tokens := []models.Token{}
	for rows.Next() {
		var token models.Token
		err := rows.Scan(&token.Token, &token.Name, &token.Limit, &token.Timestamp)
		if err != nil {
			return nil, err
		}

		tokens = append(tokens, token)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tokens, nil
}

func (d *Database) GetTokenHistory(token string) ([]models.TokenHistory, error) {
	query := "SELECT id, token, time, result, request FROM history WHERE token = ?"
	rows, err := d.db.Query(query, token)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	history := []models.TokenHistory{}
	for rows.Next() {
		var record models.TokenHistory
		err := rows.Scan(&record.ID, &record.Token, &record.RequestTime, &record.Result, &record.Request)
		if err != nil {
			return nil, err
		}

		history = append(history, record)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return history, nil
}

func (d *Database) Close() {
	d.db.Close()
}
