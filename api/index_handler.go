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
	posts, err := s.store.GetPosts()
	if err != nil {
		log.Default().Println(err.Error())
		s.HandleErrorPage(w, req, 404)
		return
	}

	data := struct {
		BlogPosts []*models.BlogPost
		Header    types.Header
	}{
		BlogPosts: posts,
		Header: types.Header{
			Navigation: s.BuildNavigationItems(req),
			User:       "",
		},
	}

	if user != nil {
		data.Header.User = user.Username
	}

	w.Header().Add("Content-Type", "text/html")
	s.renderTemplate(w, req, "index", data)
}
