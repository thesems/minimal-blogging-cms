package storage

import (
	"lifeofsems-go/models"
)

type Storage interface {
	GetPost(int) (*models.BlogPost, error)
	GetPostBy(map[string]string) (*models.BlogPost, error)
	GetPosts() ([]*models.BlogPost, error)
	CreatePost(post *models.BlogPost) int
	DeletePost(int) error
	UpdatePost(id int, setAttrs map[string]string) error

	GetUser(id int) (*models.User, error)
	GetUserBy(map[string]string) (*models.User, error)
	GetUsers() ([]*models.User, error)
	CreateUser(user *models.User) int
	DeleteUser(user *models.User) error
	UpdateUser(id int, setAttrs map[string]string) error

	GetSession(session_id string) (*models.Session, error)
	GetSessions() ([]*models.Session, error)
	CreateSession(session_id string, username string)
	DeleteSession(session_id string)
}
