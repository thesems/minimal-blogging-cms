package models

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

type Post struct {
	ID               int       `json:"id"`
	Title            string    `json:"title"`
	Content          string    `json:"content"`
	ShortDescription string    `json:"shortdescription"`
	CreatedAt        time.Time `json:"createdat"`
	UrlTitle         string    `json:"urltitle"`
	Draft            bool      `json:"draft"`
}

type PostModel struct {
	DB *sql.DB
}

func ValidatePost(bp *Post) bool {
	titleLen := len(bp.Title)
	if titleLen == 0 || titleLen > 120 {
		log.Default().Printf("[error] title is %d characters long.\n", titleLen)
		return false
	}

	if len(bp.UrlTitle) == 0 {
		log.Default().Printf("[error] urlTitle is %d characters long.\n", len(bp.UrlTitle))
		return false
	}

	currYear := time.Now().Year()
	if bp.CreatedAt.Before(time.Date(currYear, 1, 1, 1, 1, 1, 1, time.UTC)) {
		return false
	}

	return true
}

func (m PostModel) Get(id int) (*Post, error) {
	rows := m.DB.QueryRow("SELECT * FROM cms.post WHERE id=$1", id)
	if rows.Err() != nil {
		return nil, errors.New("Failed to create a SQL query for fetching the post.")
	}

	var post Post
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

func (m PostModel) GetBy(query map[string]string) (*Post, error) {
	i := 1
	queryStr := ""
	values := make([]any, 0)
	for key, val := range query {
		queryStr += fmt.Sprintf("%s=$%d", key, i)
		values = append(values, val)
		i++
	}

	rows := m.DB.QueryRow("SELECT * FROM cms.post WHERE "+queryStr, values...)
	if rows.Err() != nil {
		return nil, errors.New("Failed to create a SQL query for fetching the post.")
	}

	var post Post
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

func (m PostModel) All() ([]*Post, error) {
	rows, err := m.DB.Query("SELECT * FROM cms.post")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	posts := make([]*Post, 0)
	for rows.Next() {
		var post Post
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
func (m PostModel) Create(post *Post) int {
	id := int(uuid.New().ID())
	_, err := m.DB.Exec("INSERT INTO cms.post(id, title, content, shortdescription, createdat, urltitle, draft) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		id, post.Title, post.Content, post.ShortDescription, post.CreatedAt, post.UrlTitle, post.Draft,
	)
	if err != nil {
		log.Default().Println("Failed to insert new post. Error:", err.Error())
		return -1
	}
	return id
}

func (m PostModel) Delete(id int) error {
	_, err := m.DB.Exec("DELETE FROM cms.post WHERE id=$1", id)
	if err != nil {
		log.Default().Printf("Failed to delete the post %d.\n Error: %s", id, err.Error())
		return err
	}
	return nil
}

func (m PostModel) Update(id int, setAttrs map[string]string) error {
	return UpdateTable(m.DB, "cms.post", id, setAttrs)
}
