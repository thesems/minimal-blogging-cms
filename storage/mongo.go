package storage

import (
	"fmt"
	"lifeofsems-go/types"
)

type MongoStorage struct{}

func NewMongoStorage() *MongoStorage {
	fmt.Println("Initialized mongodb storage.")
	return &MongoStorage{}
}

func (s *MongoStorage) Get(id int) (*types.BlogPost, error) {
	return &types.BlogPost{
		ID: 0, Title: "Title", Content: "Content",
	}, nil
}
