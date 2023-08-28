package storage

import (
	"errors"
	"fmt"
	"lifeofsems-go/types"
)

type MemoryStorage struct {
	posts   []*types.BlogPost
	users   []*types.User
	session map[string]string
}

func NewMemoryStorage(users []*types.User, posts []*types.BlogPost) *MemoryStorage {
	fmt.Println("Initialized in-memory storage.")
	return &MemoryStorage{posts, users, make(map[string]string)}
}

func (ms *MemoryStorage) GetUser(username string) (*types.User, error) {
	for _, user := range ms.users {
		if user.Username == username {
			return user, nil
		}
	}
	return nil, errors.New("user not found")
}

func (ms *MemoryStorage) GetUsers() []*types.User {
	return ms.users
}

func (ms *MemoryStorage) AddUser(user *types.User) {
	ms.users = append(ms.users, user)
}

func (ms *MemoryStorage) DeleteUser(user *types.User) {
}

func (ms *MemoryStorage) GetPost(id int) (*types.BlogPost, error) {
	for _, post := range ms.posts {
		if post.ID == id {
			return post, nil
		}
	}
	return nil, errors.New("post not found")
}

func (ms *MemoryStorage) GetPosts() []*types.BlogPost {
	return ms.posts
}

func (ms *MemoryStorage) AddPost(post *types.BlogPost) {
	ms.posts = append(ms.posts, post)
}

func (ms *MemoryStorage) DeletePost(post *types.BlogPost) {
}

func (ms *MemoryStorage) GetSession(session string) (string, error) {
	username, ok := ms.session[session]
	if ok {
		return username, nil
	}
	return "", errors.New("no session uid")
}

func (ms *MemoryStorage) AddSession(session string, username string) {
	ms.session[session] = username
}

func (ms *MemoryStorage) DeleteSession(session string) {
	delete(ms.session, session)
}
