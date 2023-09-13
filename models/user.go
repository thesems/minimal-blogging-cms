package models

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
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

type UserModel struct {
	DB *sql.DB
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

func (m *UserModel) Get(id int) (*User, error) {
	rows := m.DB.QueryRow("SELECT * FROM cms.user WHERE id=$1", id)
	if rows.Err() != nil {
		return nil, errors.New("Failed to create a SQL query for fetching the post.")
	}

	var user User
	switch err := rows.Scan(&user.ID, &user.Username, &user.Firstname, &user.Lastname, &user.Password, &user.Email, &user.CreatedAt, &user.Role); err {
	case sql.ErrNoRows:
		log.Default().Println("No such user row was found.")
	case nil:
	default:
		log.Fatalln(err)
		return nil, errors.New("Failed to fetch user.")
	}
	return &user, nil
}

func (m *UserModel) GetBy(query map[string]string) (*User, error) {
	i := 1
	queryStr := ""
	values := make([]any, 0)
	for key, val := range query {
		queryStr += fmt.Sprintf("%s=$%d", key, i)
		values = append(values, val)
	}

	rows := m.DB.QueryRow("SELECT * FROM cms.user WHERE "+queryStr, values...)
	if rows.Err() != nil {
		return nil, errors.New("Failed to create a SQL query for fetching the user.")
	}

	var user User
	switch err := rows.Scan(&user.ID, &user.Username, &user.Firstname, &user.Lastname, &user.Password, &user.Email, &user.CreatedAt, &user.Role); err {
	case sql.ErrNoRows:
		log.Default().Println("No such user row was found.")
	case nil:
	default:
		log.Fatalln(err)
		return nil, errors.New("Failed to fetch post.")
	}

	return &user, nil
}

func (m *UserModel) All() ([]*User, error) {
	rows, err := m.DB.Query("SELECT * FROM cms.user")
	if err != nil {
		return nil, errors.New("Failed to query all users.")
	}
	defer rows.Close()
	users := make([]*User, 0)
	for rows.Next() {
		var user User
		err = rows.Scan(&user.ID, &user.Username, &user.Firstname, &user.Lastname, &user.Password, &user.Email, &user.CreatedAt, &user.Role)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	err = rows.Err()
	if err != nil {
		return nil, errors.New("Failed to scan the whole row during fetching users.")
	}
	return users, nil
}

func (m *UserModel) Create(user *User) int {
	id := int(uuid.New().ID())
	_, err := m.DB.Exec("INSERT INTO cms.user(id, username, firstname, lastname, password, email, createdat, role) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
		id, user.Username, user.Firstname, user.Lastname, user.Password, user.Email, user.CreatedAt, user.Role,
	)
	if err != nil {
		log.Default().Println("Failed to insert new user. Error:", err.Error())
		return -1
	}
	return id
}

func (m *UserModel) Update(id int, setAttrs map[string]string) error {
	return UpdateTable(m.DB, "cms.user", id, setAttrs)
}

func (m *UserModel) Delete(user *User) error {
	_, err := m.DB.Exec("DELETE FROM cms.user WHERE id=$1", user.ID)
	if err != nil {
		log.Default().Printf("Failed to delete the user %d.\n Error: %s", user.ID, err.Error())
		return err
	}
	return nil
}
