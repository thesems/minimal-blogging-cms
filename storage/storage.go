package storage

import "lifeofsems-go/types"

type Storage interface {
	GetPost(int) (*types.BlogPost, error)
	GetPosts() []*types.BlogPost
	AddPost(post *types.BlogPost)
	DeletePost(post *types.BlogPost)

	GetUser(username string) (*types.User, error)
	GetUsers() []*types.User
	AddUser(user *types.User)
	DeleteUser(user *types.User)

	GetSession(session string) (string, error)
	AddSession(session string, username string)
	DeleteSession(session string)
}
