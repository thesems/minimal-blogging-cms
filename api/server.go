package api

import (
	"errors"
	"fmt"
	"html/template"
	"lifeofsems-go/storage"
	"net/http"
)

type Server struct {
	listenAddr string
	store      storage.Storage
	tpl        map[string]*template.Template
}

func NewServer(listenAddr string, storage storage.Storage, tpl map[string]*template.Template) *Server {
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
	http.HandleFunc("/login", s.HandleLogin)
	http.HandleFunc("/logout", s.HandleLogout)
	http.HandleFunc("/admin", s.HandleAdmin)
	return http.ListenAndServe(":"+s.listenAddr, nil)
}

func (s *Server) renderTemplate(w http.ResponseWriter, req *http.Request, name string, data interface{}) error {
	tmpl, ok := s.tpl[name+".gohtml"]
	if !ok {
		err := errors.New("template not found")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	err := tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		err := errors.New("template execution failed")
		return err
	}

	return nil
}
