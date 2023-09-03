package models

type Role int

const (
	Normal Role = iota
	Admin
)

type User struct {
	Username string
	Password []byte
	Email    string
	Role     Role
}

func NewUser(username string, password []byte, email string, role Role) *User {
	return &User{
		Username: username, Password: password, Email: email, Role: role,
	}
}

func ValidateUser(user *User) bool {
	if len(user.Username) == 0 {
		return false
	}

	if len(user.Password) <= 6 {
		return false
	}

	if len(user.Email) == 0 {
		return false
	}
	return true
}
