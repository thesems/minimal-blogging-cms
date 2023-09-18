package api

import (
	"lifeofsems-go/types"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"golang.org/x/crypto/bcrypt"
)

func (s *Server) LoginGet(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
	w.Header().Add("Content-Type", "text/html")

	c := GetSessionCookie(req)
	http.SetCookie(w, c)
	user := GetUser(s.appEnv, w, req)

	data := struct {
		Header types.Header
		Text   string
		Meta   []types.Meta
	}{
		Header: types.Header{
			Navigation: BuildNavigationItems(s.appEnv, w, req),
			User:       "",
		},
		Text: "",
		Meta: []types.Meta{},
	}

	if user == nil {
		renderTemplate(s.appEnv, w, "login", data)
		return
	}

	http.Redirect(w, req, "/admin", http.StatusSeeOther)
}

func (s *Server) LoginPost(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
	w.Header().Add("Content-Type", "text/html")

	c := GetSessionCookie(req)
	http.SetCookie(w, c)
	user := GetUser(s.appEnv, w, req)

	data := struct {
		Header types.Header
		Text   string
		Meta   []types.Meta
	}{
		Header: types.Header{
			Navigation: BuildNavigationItems(s.appEnv, w, req),
			User:       "",
		},
		Text: "",
		Meta: []types.Meta{},
	}

	err := req.ParseForm()
	if err != nil {
		log.Fatalln(err)
	}

	data.Text = "Invalid username or password."

	username := req.Form.Get("username")
	password := req.Form.Get("password")

	user, err = s.appEnv.Users.GetBy(map[string]string{"username": username})
	if err != nil {
		data.Text = "User does not exist."
		renderTemplate(s.appEnv, w, "login", data)
		return
	}

	if bcrypt.CompareHashAndPassword(user.Password, []byte(password)) != nil {
		renderTemplate(s.appEnv, w, "login", data)
		return
	}

	s.appEnv.Sessions.Create(c.Value, username)
	http.Redirect(w, req, "/admin", http.StatusSeeOther)
}
