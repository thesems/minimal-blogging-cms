package models

import (
	"log"
	"time"
)

type Role string

const (
	Normal Role = "normal"
	Admin       = "admin"
)

type User struct {
	ID        int
	Username  string
	Firstname string
	Lastname  string
	Password  []byte
	Email     string
	CreatedAt time.Time
	Role      Role
}

func ValidateUser(user *User) bool {
	if len(user.Username) == 0 {
		log.Default().Println("Username not set.")
		return false
	}

	if len(user.Password) <= 3 {
		log.Default().Println("Password too short.")
		return false
	}

	if len(user.Email) == 0 {
		log.Default().Println("Email not set.")
		return false
	}
	return true
}

func ToRole(role string) Role {
	switch role {
	case string(Normal):
		return Normal
	case string(Admin):
		return Admin
	default:
		log.Fatalf("Role %s not found.\n", role)
		return Normal
	}
}
