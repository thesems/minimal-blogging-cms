package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"lifeofsems-go/models"
	"log"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type PostgresStorage struct {
	db      *sql.DB
	session map[string]string
}

func NewPostgresStorage(connStr string, driver string) *PostgresStorage {
	db, err := sql.Open(driver, connStr)
	if err != nil {
		log.Fatal("Cannot connect to the database.")
	}
	fmt.Printf("SQL %s storage initialized.\n", driver)
	return &PostgresStorage{db, make(map[string]string)}
}

func (s *PostgresStorage) GetPost(id int) (*models.BlogPost, error) {
	rows := s.db.QueryRow("SELECT * FROM cms.post WHERE id=$1", id)
	if rows.Err() != nil {
		return nil, errors.New("Failed to create a SQL query for fetching the post.")
	}

	var post models.BlogPost
	switch err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.ShortDescription, &post.CreatedAt, &post.UrlTitle); err {
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
	}

	rows := s.db.QueryRow("SELECT * FROM cms.post WHERE "+queryStr, values...)
	if rows.Err() != nil {
		return nil, errors.New("Failed to create a SQL query for fetching the post.")
	}

	var post models.BlogPost
	switch err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.ShortDescription, &post.CreatedAt, &post.UrlTitle); err {
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
		err = rows.Scan(&post.ID, &post.Title, &post.Content, &post.ShortDescription, &post.CreatedAt, &post.UrlTitle)
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
	_, err := s.db.Exec("INSERT INTO cms.post(id, title, content, shortdescription, createdat, urltitle) VALUES ($1, $2, $3, $4, $5, $6)",
		id, post.Title, post.Content, post.ShortDescription, post.CreatedAt, post.UrlTitle,
	)
	if err != nil {
		log.Default().Println("Failed to insert new post.")
		return -1
	}
	return id
}
func (s *PostgresStorage) DeletePost(int) error {
	return nil
}

func (s *PostgresStorage) GetUser(id int) (*models.User, error) {
	return nil, nil
}

func (s *PostgresStorage) GetUserByUsername(username string) (*models.User, error) {
	return nil, nil
}

func (s *PostgresStorage) GetUsers() ([]*models.User, error) {
	return nil, nil
}

func (s *PostgresStorage) CreateUser(user *models.User) int {
	return -1
}

func (s *PostgresStorage) DeleteUser(user *models.User) error {
	return nil
}

func (s *PostgresStorage) GetSession(session string) (string, error) {
	username, ok := s.session[session]
	if ok {
		return username, nil
	}
	return "", errors.New("no session uid")
}

func (s *PostgresStorage) CreateSession(session string, username string) {
	s.session[session] = username
}

func (s *PostgresStorage) DeleteSession(session string) {
	delete(s.session, session)
}
