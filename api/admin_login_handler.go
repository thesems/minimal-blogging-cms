package api

import (
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func (s *Server) HandleAdminLogin(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "text/html")

	c := GetSessionCookie(req)
	http.SetCookie(w, c)
	user := s.GetUser(w, req)

	if req.Method == http.MethodGet {
		if user == nil {
			s.tpl.ExecuteTemplate(w, "admin-login.gohtml", nil)
			return
		}

		http.Redirect(w, req, "/admin", http.StatusSeeOther)
	} else if req.Method == http.MethodPost {
		err := req.ParseForm()
		if err != nil {
			log.Fatalln(err)
		}

		data := struct {
			Text string
		}{
			Text: "Invalid username or password.",
		}

		username := req.Form.Get("username")
		password := req.Form.Get("password")

		user, err = s.store.GetUser(username)
		if err != nil {
			data.Text = "User does not exist."
			s.tpl.ExecuteTemplate(w, "admin-login.gohtml", data)
			return
		}

		if bcrypt.CompareHashAndPassword(user.Password, []byte(password)) != nil {
			s.tpl.ExecuteTemplate(w, "admin-login.gohtml", data)
			return
		}

		s.store.AddSession(c.Value, username)
		http.Redirect(w, req, "/admin", http.StatusSeeOther)
	}
}
