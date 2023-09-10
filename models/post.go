package models

import "time"

type BlogPost struct {
	ID               int       `json:"id"`
	Title            string    `json:"title"`
	Content          string    `json:"content"`
	ShortDescription string    `json:"shortdescription"`
	CreatedAt        time.Time `json:"createdat"`
	UrlTitle         string    `json:"urltitle"`
	Draft            bool      `json:"draft"`
}

func ValidateBlogPost(bp *BlogPost) bool {
	titleLen := len(bp.Title)
	if titleLen == 0 || titleLen > 120 {
		return false
	}

	if len(bp.UrlTitle) == 0 {
		return false
	}

	currYear := time.Now().Year()
	if bp.CreatedAt.Before(time.Date(currYear, 1, 1, 1, 1, 1, 1, time.UTC)) {
		return false
	}

	return true
}
