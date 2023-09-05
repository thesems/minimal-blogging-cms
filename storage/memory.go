package storage

import (
	"errors"
	"fmt"
	"lifeofsems-go/models"
	"log"

	"github.com/google/uuid"
)

type MemoryStorage struct {
	posts   map[int]*models.BlogPost
	users   []*models.User
	session map[string]string
}

func NewMemoryStorage(users []*models.User, posts map[int]*models.BlogPost) *MemoryStorage {
	fmt.Println("Initialized in-memory storage.")
	return &MemoryStorage{posts, users, make(map[string]string)}
}

func (ms *MemoryStorage) GetUser(username string) (*models.User, error) {
	for _, user := range ms.users {
		if user.Username == username {
			return user, nil
		}
	}
	return nil, errors.New("user not found")
}

func (ms *MemoryStorage) GetUsers() []*models.User {
	return ms.users
}

func (ms *MemoryStorage) AddUser(user *models.User) *models.User {
	ms.users = append(ms.users, user)
	return user
}

func (ms *MemoryStorage) DeleteUser(user *models.User) {}

func (ms *MemoryStorage) GetPost(id int) (*models.BlogPost, error) {
	post, ok := ms.posts[id]
	if !ok {
		return nil, errors.New("post not found")

	}
	return post, nil
}

func (ms *MemoryStorage) GetPosts() []*models.BlogPost {
	posts := make([]*models.BlogPost, len(ms.posts))

	i := 0
	for _, post := range ms.posts {
		posts[i] = post
		i++
	}
	return posts
}

func (ms *MemoryStorage) CreatePost(post *models.BlogPost) *models.BlogPost {
	post.ID = int(uuid.New().ID())

	_, ok := ms.posts[post.ID]
	if ok {
		log.Default().Printf("Post of ID %d already exists.\n", post.ID)
		return nil
	}

	ms.posts[post.ID] = post
	return post
}

func (ms *MemoryStorage) DeletePost(id int) {
	_, ok := ms.posts[id]
	if !ok {
		log.Default().Printf("Post of ID %d does not exists.\n", id)
		return
	}

	delete(ms.posts, id)
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
