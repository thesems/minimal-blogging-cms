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
	users   map[int]*models.User
	session map[string]string
}

func NewMemoryStorage(users map[int]*models.User, posts map[int]*models.BlogPost) *MemoryStorage {
	fmt.Println("Initialized in-memory storage.")
	return &MemoryStorage{posts, users, make(map[string]string)}
}

func (ms *MemoryStorage) GetUser(id int) (*models.User, error) {
	user, ok := ms.users[id]
	if !ok {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (ms *MemoryStorage) GetUserByUsername(username string) (*models.User, error) {
	for _, user := range ms.users {
		if user.Username == username {
			return user, nil
		}
	}
	return nil, errors.New("user not found")
}

func (ms *MemoryStorage) GetUsers() []*models.User {
	users := make([]*models.User, len(ms.users))

	i := 0
	for _, user := range ms.users {
		users[i] = user
		i++
	}

	return users
}

func (ms *MemoryStorage) AddUser(user *models.User) *models.User {
	user.ID = int(uuid.New().ID())

	_, ok := ms.users[user.ID]
	if ok {
		log.Default().Printf("User of ID %d already exists.\n", user.ID)
		return nil
	}

	ms.users[user.ID] = user
	return user
}

func (ms *MemoryStorage) DeleteUser(user *models.User) error {
	user, ok := ms.users[user.ID]
	if !ok {
		return errors.New("user not found")
	}
	delete(ms.users, user.ID)
	return nil
}

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

func (ms *MemoryStorage) DeletePost(id int) error {
	_, ok := ms.posts[id]
	if !ok {
		log.Default().Printf("Post of ID %d does not exists.\n", id)
		return errors.New(fmt.Sprintf("Post of ID %d does not exists.", id))
	}

	delete(ms.posts, id)
	return nil
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
