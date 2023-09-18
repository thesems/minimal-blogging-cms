package api

import (
	"lifeofsems-go/types"
	"net/http"
)

func (s *Server) NotFound(w http.ResponseWriter, req *http.Request) {
	data := struct {
		Status int
		Text   string
		Header *types.Header
	}{
		Status: 404,
		Text:   "Status Not Found",
		Header: types.NewHeader(nil, ""),
	}

	w.Header().Add("Content-Type", "text/html")
	w.WriteHeader(http.StatusNotFound)
	renderTemplate(s.appEnv, w, "error-page", data)
}

func (s *Server) HandleErrorPage(w http.ResponseWriter, req *http.Request, status int) {
	data := struct {
		Status int
		Text   string
		Header *types.Header
	}{
		Status: status,
		Text:   "",
		Header: types.NewHeader(nil, ""),
	}

	switch status {
	case http.StatusNotFound:
		data.Text = "Page not found"
	case http.StatusMethodNotAllowed:
		data.Text = "Method not allowed"
	case http.StatusUnauthorized:
		data.Text = "Unauthorized"
	case http.StatusBadRequest:
		data.Text = "Bad request"
	case http.StatusForbidden:
		data.Text = "Forbidden"
	default:
		data.Text = "Internal error"
	}

	w.Header().Add("Content-Type", "text/html")
	w.WriteHeader(status)
	renderTemplate(s.appEnv, w, "error-page", data)
}
