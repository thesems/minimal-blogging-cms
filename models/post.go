package models

import "time"

type BlogPost struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdat"`
}

func ValidateBlogPost(bp *BlogPost) bool {
	titleLen := len(bp.Title)
	contentLen := len(bp.Content)

	if titleLen == 0 || titleLen > 120 {
		return false
	}

	if contentLen == 0 {
		return false
	}

	currYear := time.Now().Year()
	if bp.CreatedAt.Before(time.Date(currYear, 1, 1, 1, 1, 1, 1, nil)) {
		return false
	}

	return true
}
