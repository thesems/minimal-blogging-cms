package types

import (
	"time"
)

type PageType int

const (
	NORMAL PageType = iota
	BLOG
)

type Page struct {
	Name      string
	Path      string
	Visible   bool
	Typ       PageType
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewPage(name string, path string, typ PageType) *Page {
	moment := time.Now()
	return &Page{
		Name:      name,
		Path:      path,
		Visible:   false,
		Typ:       typ,
		CreatedAt: moment,
		UpdatedAt: moment,
	}
}
