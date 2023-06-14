package models

import (
	"time"
)

type Token struct {
	Token     string    `json:"token"`
	Name      string    `json:"name"`
	Limit     int       `json:"limit"`
	Timestamp time.Time `json:"timestamp"`
}

type TokenHistory struct {
	ID          int       `json:"id"`
	Token       string    `json:"token"`
	RequestTime time.Time `json:"requestTime"`
	Result      string    `json:"result"`
	Request     string    `json:"request"`
}
