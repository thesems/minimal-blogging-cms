package api

import (
	"lifeofsems-go/models"
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

func (s *Server) GetUser(req *http.Request) *models.User {
	c, err := req.Cookie("session")
	if err != nil {
		return nil
	}

	var user *models.User
	session, err := s.store.GetSession(c.Value)
	if err == nil {
		user, err = s.store.GetUserBy(map[string]string{"username": session.Username})
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
			Name:   "session",
			Value:  uuid.NewString(),
			MaxAge: 60 * 60 * 2, // 2 hours
		}
	}

	return c
}

func (s *Server) BuildNavigationItems(req *http.Request) []*types.Page {

	navigation := make([]*types.Page, 0)

	user := s.GetUser(req)
	if user != nil && user.Role == models.Admin {
		navigation = append(navigation, types.NewPage("Admin", "/admin", types.NORMAL))
	}

	if s.isLoggedIn(req) {
		navigation = append(navigation, types.NewPage("Logout", "/logout", types.NORMAL))
	} else {
		navigation = append(navigation, types.NewPage("Login", "/login", types.NORMAL))
	}

	return navigation
}
