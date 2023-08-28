package api

import (
	"lifeofsems-go/types"
	"log"
	"net/http"

	"github.com/google/uuid"
)

func (s *Server) isLoggedIn(req *http.Request) bool {
	c, err := req.Cookie("session")
	if err != nil {
		return false
	}

	_, err = s.store.GetSession(c.Value)
	return err == nil
}

func (s *Server) GetUser(w http.ResponseWriter, req *http.Request) *types.User {
	c, err := req.Cookie("session")
	if err != nil {
		return nil
	}

	var user *types.User
	username, err := s.store.GetSession(c.Value)
	if err == nil {
		user, err = s.store.GetUser(username)
		if err != nil {
			log.Fatalln("Session exists for username but user does not.")
		}
	}

	return user
}

func GetSessionCookie(req *http.Request) *http.Cookie {
	c, err := req.Cookie("session")

	if err != nil {
		return &http.Cookie{
			Name:  "session",
			Value: uuid.NewString(),
		}
	}

	return c
}
