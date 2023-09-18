package api

import (
	"lifeofsems-go/models"
	"lifeofsems-go/types"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (s *Server) IndexGet(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	user := GetUser(s.appEnv, w, req)
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
			Navigation: BuildNavigationItems(s.appEnv, w, req),
			User:       "",
		},
		Meta: []types.Meta{},
	}

	if user != nil {
		data.Header.User = user.Username
	}

	w.Header().Add("Content-Type", "text/html")
	renderTemplate(s.appEnv, w, "index", data)
}
