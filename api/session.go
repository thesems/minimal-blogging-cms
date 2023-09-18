package api

import (
	"lifeofsems-go/env"
	"lifeofsems-go/models"
	"lifeofsems-go/types"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

const expiryTime = 2 * 60 * 60
const cleaningTime = 1 * 60 * 60

func isLoggedIn(appEnv env.Env, req *http.Request) bool {
	c, err := req.Cookie("session")
	if err != nil {
		return false
	}

	_, err = appEnv.Sessions.Get(c.Value)
	return err == nil
}

func GetUser(appEnv env.Env, w http.ResponseWriter, req *http.Request) *models.User {
	c, err := req.Cookie("session")
	if err != nil {
		return nil
	}

	c.MaxAge = expiryTime
	http.SetCookie(w, c)

	var user *models.User
	session, err := appEnv.Sessions.Get(c.Value)
	if err == nil {
		user, err = appEnv.Users.GetBy(map[string]string{"username": session.Username})
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
			MaxAge: expiryTime,
		}
	}
	return c
}

func (s *Server) CleanSessions() {
	if time.Now().Sub(s.appEnv.LastSessionCleaning) > cleaningTime {
		sessions, err := s.appEnv.Sessions.All()
		if err != nil {
			log.Default().Println("No cleaning was performed. Error:", err.Error())
			return
		}
		for _, session := range sessions {
			if time.Now().Sub(session.LastActivity) > (time.Second * expiryTime) {
				s.appEnv.Sessions.Delete(session.ID)
			}
		}

		s.appEnv.LastSessionCleaning = time.Now()
	}

	time.Sleep(time.Second * 60)
	go s.CleanSessions()
}

func BuildNavigationItems(appEnv env.Env, w http.ResponseWriter, req *http.Request) []*types.Page {

	navigation := make([]*types.Page, 0)

	user := GetUser(appEnv, w, req)
	if user != nil && user.Role == models.Admin {
		navigation = append(navigation, types.NewPage("Admin", "/admin", types.NORMAL))
	}

	if isLoggedIn(appEnv, req) {
		navigation = append(navigation, types.NewPage("Logout", "/logout", types.NORMAL))
	} else {
		navigation = append(navigation, types.NewPage("Login", "/login", types.NORMAL))
	}

	return navigation
}
