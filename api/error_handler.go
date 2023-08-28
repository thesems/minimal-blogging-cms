package api

import (
	"fmt"
	"net/http"
)

func (s *Server) HandleErrorPage(w http.ResponseWriter, req *http.Request, status int) {
	data := struct {
		Status int
		Text   string
	}{
		Status: status,
		Text:   "",
	}

	switch status {
	case http.StatusNotFound:
		data.Text = "Page not found"
	case http.StatusMethodNotAllowed:
		data.Text = "Method not allowed"
	case http.StatusBadRequest:
		data.Text = "Bad request"
	case http.StatusForbidden:
		data.Text = "Forbidden"
	default:
		data.Text = "Internal error"
	}

	w.Header().Add("Content-Type", "text/html")
	w.WriteHeader(status)
	s.tpl.ExecuteTemplate(w, "error-page.gohtml", data)
	fmt.Println("hello", status, req.URL.Path)
}
