package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"lifeofsems-go/models"
	"log"
	"strconv"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type PostgresStorage struct {
	db *sql.DB
}

func NewPostgresStorage(connStr string, driver string) *PostgresStorage {
	db, err := sql.Open(driver, connStr)
	if err != nil {
		log.Fatal("Cannot connect to the database.")
	}
	fmt.Printf("SQL %s storage initialized.\n", driver)
	return &PostgresStorage{db}
}

func (s *PostgresStorage) GetPost(id int) (*models.BlogPost, error) {
	rows := s.db.QueryRow("SELECT * FROM cms.post WHERE id=$1", id)
	if rows.Err() != nil {
		return nil, errors.New("Failed to create a SQL query for fetching the post.")
	}

	var post models.BlogPost
	switch err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.ShortDescription, &post.CreatedAt, &post.UrlTitle, &post.Draft); err {
	case sql.ErrNoRows:
		log.Default().Println("No such post row was found.")
	case nil:
	default:
		log.Fatalln(err)
		return nil, errors.New("Failed to fetch post.")
	}

	return &post, nil
}

func (s *PostgresStorage) GetPostBy(query map[string]string) (*models.BlogPost, error) {
	i := 1
	queryStr := ""
	values := make([]any, 0)
	for key, val := range query {
		queryStr += fmt.Sprintf("%s=$%d", key, i)
		values = append(values, val)
		i++
	}

	rows := s.db.QueryRow("SELECT * FROM cms.post WHERE "+queryStr, values...)
	if rows.Err() != nil {
		return nil, errors.New("Failed to create a SQL query for fetching the post.")
	}

	var post models.BlogPost
	switch err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.ShortDescription, &post.CreatedAt, &post.UrlTitle, &post.Draft); err {
	case sql.ErrNoRows:
		log.Default().Println("No such post row was found.")
	case nil:
	default:
		log.Fatalln(err)
		return nil, errors.New("Failed to fetch post.")
	}

	return &post, nil
}

func (s *PostgresStorage) GetPosts() ([]*models.BlogPost, error) {
	rows, err := s.db.Query("SELECT * FROM cms.post")
	if err != nil {
		return nil, errors.New("Failed to query all posts.")
	}
	defer rows.Close()
	posts := make([]*models.BlogPost, 0)
	for rows.Next() {
		var post models.BlogPost
		err = rows.Scan(&post.ID, &post.Title, &post.Content, &post.ShortDescription, &post.CreatedAt, &post.UrlTitle, &post.Draft)
		if err != nil {
			return nil, err
		}
		posts = append(posts, &post)
	}
	err = rows.Err()
	if err != nil {
		return nil, errors.New("Failed to scan the whole row during fetching posts.")
	}
	return posts, nil
}
func (s *PostgresStorage) CreatePost(post *models.BlogPost) int {
	id := int(uuid.New().ID())
	_, err := s.db.Exec("INSERT INTO cms.post(id, title, content, shortdescription, createdat, urltitle, draft) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		id, post.Title, post.Content, post.ShortDescription, post.CreatedAt, post.UrlTitle, post.Draft,
	)
	if err != nil {
		log.Default().Println("Failed to insert new post. Error:", err.Error())
		return -1
	}
	return id
}
func (s *PostgresStorage) DeletePost(id int) error {
	_, err := s.db.Exec("DELETE FROM cms.post WHERE id=$1", id)
	if err != nil {
		log.Default().Printf("Failed to delete the post %d.\n Error: %s", id, err.Error())
		return err
	}
	return nil
}

func (s *PostgresStorage) updateTable(table string, id int, setAttrs map[string]string) error {
	i := 1
	queryStr := ""
	values := make([]any, 0)
	for key, val := range setAttrs {
		queryStr += fmt.Sprintf("%s=$%d,", key, i)
		values = append(values, val)
		i++
	}
	queryStr = queryStr[:len(queryStr)-1]
	values = append(values, id)
	_, err := s.db.Exec("UPDATE "+table+" SET "+queryStr+" WHERE id=$"+strconv.Itoa(i), values...)
	if err != nil {
		return err
	}
	return nil
}

func (s *PostgresStorage) UpdatePost(id int, setAttrs map[string]string) error {
	return s.updateTable("cms.post", id, setAttrs)
}

func (s *PostgresStorage) GetUser(id int) (*models.User, error) {
	rows := s.db.QueryRow("SELECT * FROM cms.user WHERE id=$1", id)
	if rows.Err() != nil {
		return nil, errors.New("Failed to create a SQL query for fetching the post.")
	}

	var user models.User
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

func (s *PostgresStorage) GetUserBy(query map[string]string) (*models.User, error) {
	i := 1
	queryStr := ""
	values := make([]any, 0)
	for key, val := range query {
		queryStr += fmt.Sprintf("%s=$%d", key, i)
		values = append(values, val)
	}

	rows := s.db.QueryRow("SELECT * FROM cms.user WHERE "+queryStr, values...)
	if rows.Err() != nil {
		return nil, errors.New("Failed to create a SQL query for fetching the user.")
	}

	var user models.User
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

func (s *PostgresStorage) GetUsers() ([]*models.User, error) {
	rows, err := s.db.Query("SELECT * FROM cms.user")
	if err != nil {
		return nil, errors.New("Failed to query all users.")
	}
	defer rows.Close()
	users := make([]*models.User, 0)
	for rows.Next() {
		var user models.User
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

func (s *PostgresStorage) CreateUser(user *models.User) int {
	id := int(uuid.New().ID())
	_, err := s.db.Exec("INSERT INTO cms.user(id, username, firstname, lastname, password, email, createdat, role) VALUES ($1, $2, $3, $4, $5, $6)",
		id, user.Username, user.Firstname, user.Lastname, user.Password, user.Email, user.CreatedAt, user.Role,
	)
	if err != nil {
		log.Default().Println("Failed to insert new user. Error:", err.Error())
		return -1
	}
	return id
}

func (s *PostgresStorage) UpdateUser(id int, setAttrs map[string]string) error {
	return s.updateTable("cms.user", id, setAttrs)
}

func (s *PostgresStorage) DeleteUser(user *models.User) error {
	_, err := s.db.Exec("DELETE FROM cms.user WHERE id=$1", user.ID)
	if err != nil {
		log.Default().Printf("Failed to delete the user %d.\n Error: %s", user.ID, err.Error())
		return err
	}
	return nil
}

func (s *PostgresStorage) GetSession(session_id string) (*models.Session, error) {
	row := s.db.QueryRow("SELECT * FROM cms.session WHERE id=$1", session_id)
	var session models.Session
	err := row.Scan(&session.ID, &session.Username, &session.LastActivity)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("could not find session %s", session_id))
	}
	return &session, nil
}

func (s *PostgresStorage) CreateSession(session_id string, username string) {
	lastActivity := time.Now()
	_, err := s.db.Exec("INSERT INTO cms.session (id,username,lastactivity) VALUES ($1,$2,$3)", session_id, username, lastActivity)
	if err != nil {
		log.Default().Fatalln(err.Error())
	}
}

func (s *PostgresStorage) DeleteSession(session_id string) {
	_, err := s.db.Exec("DELETE FROM cms.session WHERE id=$1", session_id)
	if err != nil {
		log.Default().Fatalln(err.Error())
	}
}
