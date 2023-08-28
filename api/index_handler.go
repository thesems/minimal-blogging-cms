package api

import (
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

	data := struct {
		BlogPosts []*types.BlogPost
		Pages     []*types.Page
	}{
		BlogPosts: s.store.GetAll(),
		Pages: []*types.Page{
			types.NewPage("Admin", "/login", types.NORMAL),
		},
	}

	w.Header().Add("Content-Type", "text/html")
	s.tpl.ExecuteTemplate(w, "layout", data)
}
