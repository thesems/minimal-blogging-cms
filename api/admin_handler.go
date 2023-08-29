package api

import (
	"fmt"
	"lifeofsems-go/types"
	"log"
	"net/http"
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

	if req.Method == http.MethodGet {

		data := struct {
			Header types.Header
		}{
			Header: types.Header{
				Navigation: s.BuildNavigationItems(req),
				User:       "",
			},
		}
		if user != nil {
			data.Header.User = user.Username
		}
		s.renderTemplate(w, req, "admin", data)

	} else if req.Method == http.MethodPost {
		err := req.ParseForm()
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println("Got:", req.Form)
	}
}
