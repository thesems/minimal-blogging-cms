package api

import (
	"lifeofsems-go/models"
	"lifeofsems-go/types"
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
	data := struct {
		BlogPosts []*models.BlogPost
		Header    types.Header
	}{
		BlogPosts: s.store.GetPosts(),
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
