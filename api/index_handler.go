package api

import (
	"lifeofsems-go/models"
	"lifeofsems-go/types"
	"log"
	"net/http"
	"strings"
)

func (s *Server) HandleIndex(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		s.HandleErrorPage(w, req, 405)
		return
	}

	// Disallow any URL except /
	tokens := strings.Split(req.URL.Path, "/")
	if tokens[1] != "" {
		s.HandleErrorPage(w, req, 404)
		return
	}

	user := s.GetUser(w, req)
	posts, err := s.appEnv.Posts.All()
	if err != nil {
		log.Default().Println(err.Error())
		s.HandleErrorPage(w, req, 404)
		return
	}

	releasedPosts := make([]*models.Post, 0)
	for _, post := range posts {
		if !post.Draft {
			releasedPosts = append(releasedPosts, post)
		}
	}

	data := struct {
		Posts  []*models.Post
		Header types.Header
		Meta   []types.Meta
	}{
		Posts: releasedPosts,
		Header: types.Header{
			Navigation: s.BuildNavigationItems(w, req),
			User:       "",
		},
		Meta: []types.Meta{
			{Src: "hello"},
		},
	}

	if user != nil {
		data.Header.User = user.Username
	}

	w.Header().Add("Content-Type", "text/html")
	s.renderTemplate(w, req, "index", data)
}
