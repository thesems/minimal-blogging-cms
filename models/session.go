package models

import "time"

type Session struct {
	ID           string    `json:"id"`
	Username     string    `json:"username"`
	LastActivity time.Time `json:"lastActivity"`
}
