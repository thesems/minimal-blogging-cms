package api

import (
	"fmt"
	"html/template"
	"lifeofsems-go/storage"
	"net/http"
)

type Server struct {
	listenAddr string
	store      storage.Storage
	tpl        *template.Template
}

func NewServer(listenAddr string, storage storage.Storage, tpl *template.Template) *Server {
	fmt.Println("Start HTTP server on port", listenAddr)

	return &Server{
		listenAddr: listenAddr,
		store:      storage,
		tpl:        tpl,
	}
}

func (s *Server) Start() error {
	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/public/", http.StripPrefix("/public/", fs))

	http.HandleFunc("/", s.HandleIndex)
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, req *http.Request) {
		s.HandleErrorPage(w, req, 404)
	})
	http.HandleFunc("/blog/", s.HandleBlogPage)
	http.HandleFunc("/login", s.HandleAdminLogin)
	http.HandleFunc("/logout", s.HandleAdminLogout)
	http.HandleFunc("/admin", s.HandleAdmin)
	return http.ListenAndServe(":"+s.listenAddr, nil)
}
