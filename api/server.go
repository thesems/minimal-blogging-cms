package api

import (
	"errors"
	"fmt"
	"lifeofsems-go/env"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Server struct {
	listenAddr string
	appEnv     env.Env
}

func NewServer(listenAddr string, appEnv env.Env) *Server {
	fmt.Println("Start HTTP server on port", listenAddr)

	return &Server{
		listenAddr: listenAddr,
		appEnv:     appEnv,
	}
}

func (s *Server) Start() error {
	router := httprouter.New()
	router.NotFound = http.HandlerFunc(s.NotFound)
	router.ServeFiles("/public/*filepath", http.Dir("./public"))
	router.GET("/", s.IndexGet)
	router.GET("/post/:post", s.PostGet)
	router.POST("/post", s.PostPost)
	router.PUT("/post/:postId", s.PostPut)
	router.DELETE("/post/:postId", s.PostDelete)
	router.GET("/login", s.LoginGet)
	router.POST("/login", s.LoginPost)
	router.GET("/admin", s.AdminGet)
	router.GET("/user/:userId", s.UserGet)
	router.POST("/user", s.UserPost)
	router.PUT("/user/:userId", s.UserPut)
	router.DELETE("/user/:userId", s.UserDelete)

	go s.CleanSessions()
	return http.ListenAndServe(":"+s.listenAddr, router)
}

func renderTemplate(appEnv env.Env, w http.ResponseWriter, name string, data interface{}) error {
	tmpl, ok := appEnv.Tpl[name+".gohtml"]
	if !ok {
		err := errors.New("template not found")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	err := tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	return nil
}
