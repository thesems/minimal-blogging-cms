package api

import (
	"log"
	"net/http"
)

func (s *Server) HandleLogout(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "text/html")

	if !isLoggedIn(s.appEnv, req) {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}

	if req.Method == http.MethodGet {
		c, err := req.Cookie("session")
		if err != nil {
			log.Fatalln("User should be logged in, but cookie not found.")
		}
		s.appEnv.Sessions.Delete(c.Value)
		c.Value = ""
		c.MaxAge = -1
		http.SetCookie(w, c)
		http.Redirect(w, req, "/", http.StatusSeeOther)
	}
}
