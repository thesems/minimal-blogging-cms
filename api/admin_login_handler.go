package api

import (
	"fmt"
	"log"
	"net/http"
)

func (s *Server) HandleAdminLogin(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "text/html")

	if req.Method == http.MethodGet {

		auth := IsAuthorized(req)
		if !auth {
			s.tpl.ExecuteTemplate(w, "admin-login.gohtml", nil)
			return
		} else {
			http.Redirect(w, req, "/admin", http.StatusSeeOther)
		}

	} else if req.Method == http.MethodPost {
		err := req.ParseForm()
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println("Got:", req.Form)

		username := req.Form.Get("username")
		password := req.Form.Get("password")

		if username == "admin" && password == "123" {
			http.SetCookie(w, &http.Cookie{
				Name:  "session",
				Value: "topsecret",
			})
			http.Redirect(w, req, "/admin", http.StatusSeeOther)
		} else {
			data := struct {
				Text string
			}{
				Text: "Invalid username or password.",
			}
			s.tpl.ExecuteTemplate(w, "admin-login.gohtml", data)
		}
	}
}
