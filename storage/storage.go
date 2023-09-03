package storage

import (
	"lifeofsems-go/models"
)

type Storage interface {
	GetPost(int) (*models.BlogPost, error)
	GetPosts() []*models.BlogPost
	CreatePost(post *models.BlogPost) *models.BlogPost
	DeletePost(int)

	GetUser(username string) (*models.User, error)
	GetUsers() []*models.User
	AddUser(user *models.User) *models.User
	DeleteUser(user *models.User)

	GetSession(session string) (string, error)
	AddSession(session string, username string)
	DeleteSession(session string)
}
