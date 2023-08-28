package types

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
