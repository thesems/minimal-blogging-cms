package api

import (
	"fmt"
	"log"
	"net/http"
)

func (s *Server) HandleAdmin(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "text/html")

	if req.Method == http.MethodGet {

		auth := IsAuthorized(req)
		if !auth {
			s.HandleErrorPage(w, req, http.StatusForbidden)
			return
		}

		data := struct {
			accessKey string
		}{}

		s.tpl.ExecuteTemplate(w, "admin.gohtml", data)
	} else if req.Method == http.MethodPost {
		err := req.ParseForm()
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println("Got:", req.Form)
	}
}
