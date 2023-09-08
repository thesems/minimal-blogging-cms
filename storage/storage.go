package storage

import (
	"lifeofsems-go/models"
)

type Storage interface {
	GetPost(int) (*models.BlogPost, error)
	GetPosts() []*models.BlogPost
	CreatePost(post *models.BlogPost) *models.BlogPost
	DeletePost(int) error

	GetUser(id int) (*models.User, error)
	GetUserByUsername(username string) (*models.User, error)
	GetUsers() []*models.User
	AddUser(user *models.User) *models.User
	DeleteUser(user *models.User) error

	GetSession(session string) (string, error)
	AddSession(session string, username string)
	DeleteSession(session string)
}
