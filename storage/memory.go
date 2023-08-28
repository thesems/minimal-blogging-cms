package storage

import (
	"errors"
	"fmt"
	"lifeofsems-go/types"
)

type MemoryStorage struct {
	posts []*types.BlogPost
}

func NewMemoryStorage(posts []*types.BlogPost) *MemoryStorage {
	fmt.Println("Initialized in-memory storage.")
	return &MemoryStorage{posts}
}

func (ms *MemoryStorage) Get(id int) (*types.BlogPost, error) {
	for _, post := range ms.posts {
		if post.ID == id {
			return post, nil
		}
	}
	return nil, errors.New("post not found")
}

func (ms *MemoryStorage) GetAll() []*types.BlogPost {
	return ms.posts
}

func (ms *MemoryStorage) Add(post *types.BlogPost) {
	ms.posts = append(ms.posts, post)
}

func (ms *MemoryStorage) Delete(post *types.BlogPost) {
}
