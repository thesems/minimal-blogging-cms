package api

import (
	"lifeofsems-go/models"
	"lifeofsems-go/types"
	"log"
	"net/http"
	"strconv"
)

func (s *Server) HandleAdmin(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "text/html")

	// if !s.isLoggedIn(req) {
	// 	http.Redirect(w, req, "/login", http.StatusSeeOther)
	// 	return
	// }

	user := s.GetUser(w, req)
	// if user == nil {
	// 	http.Redirect(w, req, "/login", http.StatusSeeOther)
	// 	return
	// }

	// if user.Role != types.Admin {
	// 	s.HandleErrorPage(w, req, 401)
	// 	return
	// }

	log.Default().Println(req.Method)

	if req.Method == http.MethodGet {
		data := struct {
			Header     types.Header
			ActiveTab  string
			UpdatePost *models.BlogPost
			Posts      []*models.BlogPost
			Users      []*models.User
		}{
			Header: types.Header{
				Navigation: s.BuildNavigationItems(req),
				User:       "",
			},
			ActiveTab:  "posts",
			UpdatePost: nil,
			Posts:      s.store.GetPosts(),
			Users:      s.store.GetUsers(),
		}

		if user != nil {
			data.Header.User = user.Username
		}

		req.ParseForm()
		tab := req.Form.Get("tab")
		if tab == "users" {
			data.ActiveTab = "users"
		}
		edit := req.Form.Get("edit")
		if edit != "" {
			postId, err := strconv.Atoi(edit)
			if err != nil {
				http.Error(w, "Failed to parse the edit post id.", http.StatusBadRequest)
				return
			}
			post, err := s.store.GetPost(postId)
			if err != nil {
				http.Error(w, "Failed to find the post with id.", http.StatusBadRequest)
				return
			}
			data.UpdatePost = post
			s.CreatePostRowEdit(w, req, post)
			return
		}

		s.renderTemplate(w, req, "admin", data)

	} else if req.Method == http.MethodPost {
	}
}
