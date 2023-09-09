package api

import (
	"fmt"
	"lifeofsems-go/models"
	"lifeofsems-go/types"
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
	fmt.Println(req.Method, req.URL.String())
	// if user == nil {
	// 	http.Redirect(w, req, "/login", http.StatusSeeOther)
	// 	return
	// }

	// if user.Role != types.Admin {
	// 	s.HandleErrorPage(w, req, 401)
	// 	return
	// }

	if req.Method == http.MethodGet {
		posts, err := s.store.GetPosts()
		if err != nil {
			return
		}
		users, err := s.store.GetUsers()
		if err != nil {
			return
		}

		data := struct {
			Header    types.Header
			ActiveTab string
			Posts     []*models.BlogPost
			Users     []*models.User
		}{
			Header: types.Header{
				Navigation: s.BuildNavigationItems(req),
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
				post, err := s.store.GetPost(resId)
				if err != nil {
					http.Error(w, "Failed to find the post with id.", http.StatusBadRequest)
					return
				}
				s.CreatePostRowEdit(w, req, post)
			} else if data.ActiveTab == "users" {
				user, err := s.store.GetUser(resId)
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
