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
	query := "SELECT COUNT(*) FROM admins WHERE username = ? AND password = ?"

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
	query := "SELECT token, name, `limit`, created_at FROM users"
	rows, err := d.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tokens := []models.Token{}
	for rows.Next() {
		var token models.Token
		var createdAt []uint8
		err := rows.Scan(&token.Token, &token.Name, &token.Limit, &createdAt)
		if err != nil {
			return nil, err
		}

		token.Timestamp, err = time.Parse("2006-01-02 15:04:05", string(createdAt))
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

func (d *Database) DeleteToken(tokenValue string) error {
	query := "DELETE FROM users WHERE token = ?"
	_, err := d.db.Exec(query, tokenValue)
	if err != nil {
		return err
	}

	return nil
}

func (d *Database) ValidateToken(token string) (bool, error) {
	query := "SELECT COUNT(*) FROM users WHERE token = ?"

	var count int
	err := d.db.QueryRow(query, token).Scan(&count)
	if err != nil {
		return false, err
	}

	if count > 0 {
		return true, nil
	} else {
		return false, nil
	}
}

func (d *Database) GetLastID(token string) (*int, error) {
	query := "SELECT last_id FROM users WHERE token = ?"

	var lastID sql.NullInt64
	err := d.db.QueryRow(query, token).Scan(&lastID)
	if err != nil {
		return nil, err
	}

	if lastID.Valid {
		id := int(lastID.Int64)
		return &id, nil
	} else {
		return nil, nil
	}
}

func (d *Database) GetLastMessage() (*models.Message, error) {
	query := "SELECT message_id, coin, base_currency, platform, leverage, side, entries, margin, sl, timestamp, `exit`, risk FROM messages ORDER BY message_id DESC LIMIT 1"

	var message models.Message
	err := d.db.QueryRow(query).Scan(&message.MessageID, &message.Coin, &message.BaseCurrency, &message.Platform, &message.Leverage, &message.Side, &message.Entries, &message.Margin, &message.SL, &message.Timestamp, &message.Exit, &message.Risk)
	if err != nil {
		return nil, err
	}

	_, err = d.db.Exec("UPDATE users SET `limit` = `limit` - 1")
	if err != nil {
		return nil, err
	}

	return &message, nil
}

func (d *Database) UpdateLastID(token string, lastID int) error {
	query := "UPDATE users SET last_id = ? WHERE token = ?"

	_, err := d.db.Exec(query, lastID, token)
	if err != nil {
		return err
	}

	return nil
}

func (d *Database) GetMessagesAfterID(lastID int) ([]models.Message, error) {
	query := "SELECT message_id, coin, base_currency, platform, leverage, side, entries, margin, sl, timestamp, `exit`, risk FROM messages WHERE message_id > ?"

	rows, err := d.db.Query(query, lastID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	messages := []models.Message{}
	for rows.Next() {
		var message models.Message
		err := rows.Scan(&message.MessageID, &message.Coin, &message.BaseCurrency, &message.Platform, &message.Leverage, &message.Side, &message.Entries, &message.Margin, &message.SL, &message.Timestamp, &message.Exit, &message.Risk)
		if err != nil {
			return nil, err
		}

		messages = append(messages, message)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	decreaseQuery := "UPDATE users SET `limit` = `limit` - 1"
	_, err = d.db.Exec(decreaseQuery)
	if err != nil {
		return nil, err
	}

	return messages, nil
}

func (d *Database) GetTokenByMessageID(messageID string) (string, error) {
	query := "SELECT token FROM messages WHERE message_id = ?"

	var token string
	err := d.db.QueryRow(query, messageID).Scan(&token)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", err
	}

	decreaseQuery := "UPDATE users SET `limit` = `limit` - 1"
	_, err = d.db.Exec(decreaseQuery)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (d *Database) AddHistory(token string, coins []string, request string) error {
	query := "INSERT INTO history (token, time, result, request) VALUES (?, NOW(), ?, ?)"

	for _, coin := range coins {
		_, err := d.db.Exec(query, token, coin, request)
		if err != nil {
			return err
		}
	}

	return nil
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
