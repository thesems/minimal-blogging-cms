package types

type BlogPost struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

func ValidateBlogPost(bp *BlogPost) bool {
	return true
}
