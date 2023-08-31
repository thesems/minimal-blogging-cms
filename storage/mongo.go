package storage

import (
	"fmt"
	"lifeofsems-go/models"
)

type MongoStorage struct{}

func NewMongoStorage() *MongoStorage {
	fmt.Println("Initialized mongodb storage.")
	return &MongoStorage{}
}

func (s *MongoStorage) Get(id int) (*models.BlogPost, error) {
	return &models.BlogPost{
		ID: 0, Title: "Title", Content: "Content",
	}, nil
}
