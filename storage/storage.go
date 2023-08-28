package storage

import "lifeofsems-go/types"

type Storage interface {
	Get(int) (*types.BlogPost, error)
	GetAll() []*types.BlogPost
	Add(post *types.BlogPost)
	Delete(post *types.BlogPost)
}
