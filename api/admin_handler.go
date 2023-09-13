package api

import (
	"lifeofsems-go/models"
	"lifeofsems-go/types"
	"net/http"
	"strconv"
)

func (s *Server) HandleAdmin(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "text/html")

	if !s.isLoggedIn(req) {
		http.Redirect(w, req, "/login", http.StatusSeeOther)
		return
	}

	user := s.GetUser(w, req)
	if user == nil {
		http.Redirect(w, req, "/login", http.StatusSeeOther)
		return
	}

	if user.Role != models.Admin {
		s.HandleErrorPage(w, req, 401)
		return
	}

	if req.Method == http.MethodGet {
		posts, err := s.appEnv.Posts.All()
		if err != nil {
			return
		}
		users, err := s.appEnv.Users.All()
		if err != nil {
			return
		}

		data := struct {
			Header    types.Header
			ActiveTab string
			Posts     []*models.Post
			Users     []*models.User
		}{
			Header: types.Header{
				Navigation: s.BuildNavigationItems(w, req),
				User:       "",
			},
			ActiveTab: "posts",
			Posts:     posts,
			Users:     users,
		}

		if user != nil {
			data.Header.User = user.Username
		}

		req.ParseForm()
		tab := req.Form.Get("view")
		if tab == "users" {
			data.ActiveTab = "users"
		}

		edit := req.Form.Get("edit")
		if edit != "" {
			resId, err := strconv.Atoi(edit)
			if err != nil {
				http.Error(w, "Failed to parse the edit resource id.", http.StatusBadRequest)
				return
			}

			if data.ActiveTab == "posts" {
				post, err := s.appEnv.Posts.Get(resId)
				if err != nil {
					http.Error(w, "Failed to find the post with id.", http.StatusBadRequest)
					return
				}
				s.CreatePostRowEdit(w, req, post)
			} else if data.ActiveTab == "users" {
				user, err := s.appEnv.Users.Get(resId)
				if err != nil {
					http.Error(w, "Failed to find the user with id.", http.StatusBadRequest)
					return
				}
				s.RenderUserEditRow(w, req, user)
			}
			return
		}

		s.renderTemplate(w, req, "admin", data)

	} else if req.Method == http.MethodPost {
	}
}
