package api

import (
	"net/http"
)

func (s *Server) HandleAdminLogout(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "text/html")

	if req.Method == http.MethodGet {
		return
	} else if req.Method == http.MethodPost {
		c, err := req.Cookie("session")
		if err == nil {
			c.MaxAge = -1
			http.SetCookie(w, c)
		}
		http.Redirect(w, req, "/", http.StatusSeeOther)
	}
}
