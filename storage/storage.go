package storage

import (
	"lifeofsems-go/models"
	"lifeofsems-go/types"
)

type Storage interface {
	GetPost(int) (*models.BlogPost, error)
	GetPosts() []*models.BlogPost
	CreatePost(post *models.BlogPost) *models.BlogPost
	DeletePost(int)

	GetUser(username string) (*types.User, error)
	GetUsers() []*types.User
	AddUser(user *types.User) *types.User
	DeleteUser(user *types.User)

	GetSession(session string) (string, error)
	AddSession(session string, username string)
	DeleteSession(session string)
}
