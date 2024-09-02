package models

import (
	"net/http"
	"time"
)

type AuthHandler func(w http.ResponseWriter, r *http.Request, user *User)

type User struct {
	ID        string    `json:"id"`
	Aud       string    `json:"aud"`
	Role      string    `json:"role"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type ApiKey struct {
	ID        string    `json:"id"`
	ApiKey    string    `json:"api_key"`
	Owner     string    `json:"owner"`
	Scopes    []string  `json:"scopes"`
	CreatedAt time.Time `json:"created_at"`
}
