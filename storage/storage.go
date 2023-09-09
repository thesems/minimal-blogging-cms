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

	GetUser(id int) (*models.User, error)
	GetUserByUsername(username string) (*models.User, error)
	GetUsers() ([]*models.User, error)
	CreateUser(user *models.User) int
	DeleteUser(user *models.User) error

	GetSession(session string) (string, error)
	CreateSession(session string, username string)
	DeleteSession(session string)
}
