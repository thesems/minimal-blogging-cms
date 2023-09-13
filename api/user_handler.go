package api

import (
	"fmt"
	"html/template"
	"lifeofsems-go/models"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func (s *Server) HandleUser(w http.ResponseWriter, req *http.Request) {
	tokens := strings.Split(req.URL.Path, "/")
	fmt.Println(req.Method, req.URL.String())

	if len(tokens) < 3 {
		s.HandleErrorPage(w, req, http.StatusNotFound)
		return
	}

	hxReq := req.Header.Get("Hx-Request") == "true"
	// hxCurrUrl := req.Header.Get("Hx-Current-Url")

	if tokens[2] == "create" && req.Method == http.MethodPost {
		if hxReq {
			user := s.ParseUser(w, req)
			if user == nil {
				return
			}
			userId := s.appEnv.Users.Create(user)
			if userId == -1 {
				return
			}
			user.ID = userId
			s.RenderUserRow(w, req, user)
		} else {
			// Create user
		}
		return
	}

	userId, err := strconv.Atoi(tokens[2])
	if err != nil {
		s.HandleErrorPage(w, req, http.StatusNotFound)
		return
	}

	switch req.Method {
	case http.MethodGet:
		user, err := s.appEnv.Users.Get(userId)
		if err != nil {
			http.Error(w, fmt.Sprintf("Could not find user of ID %d.", userId), http.StatusBadRequest)
			return
		}
		if hxReq {
			s.RenderUserRow(w, req, user)
		} else {
			// TODO: Get user as json
			w.Write([]byte(user.Username))
		}
	case http.MethodPut:
		if hxReq {
			user := s.ParseUser(w, req)
			if user == nil {
				return
			}
			password, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
			if err != nil {
				http.Error(w, "Could not parse user put request.", http.StatusInternalServerError)
				log.Default().Println(err.Error())
				return
			}
			err = s.appEnv.Users.Update(user.ID, map[string]string{
				"username": user.Username,
				"password": string(password),
				"email":    user.Email,
				"role":     string(user.Role),
			})
			if err != nil {
				http.Error(w, "Could not parse user put request.", http.StatusInternalServerError)
				log.Default().Println(err.Error())
				return
			}
			s.RenderUserRow(w, req, user)
		}
	case http.MethodDelete:
		user, err := s.appEnv.Users.Get(userId)
		if err != nil {
			http.Error(w, fmt.Sprintf("Could not find user with ID %d.", userId), http.StatusBadRequest)
			return
		}
		err = s.appEnv.Users.Delete(user)
		if err != nil {
			http.Error(w, fmt.Sprintf("Could not delete user with ID %d.", userId), http.StatusBadRequest)
			return
		}
	default:
		s.HandleErrorPage(w, req, http.StatusMethodNotAllowed)
	}
}

func (s *Server) ParseUser(w http.ResponseWriter, req *http.Request) *models.User {
	err := req.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse user form.", http.StatusInternalServerError)
		log.Default().Println("Failed to parse user form.")
		return nil
	}

	id := req.Form.Get("ID")
	username := req.Form.Get("username")
	password := req.Form.Get("password")
	email := req.Form.Get("email")
	role := req.Form.Get("role")

	var user *models.User
	if id != "" {
		idInt, err := strconv.Atoi(id)
		if err != nil {
			log.Default().Println("Failed to convert user ID to integer.")
			return nil
		}
		user, err = s.appEnv.Users.Get(idInt)
		if err != nil {
			log.Default().Println("Could find the user with ID.")
			return nil
		}
		if username != "" {
			user.Username = username
		}
		if password != "" {
			pw, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
			if err != nil {
				http.Error(w, "Failed to generate password.", http.StatusInternalServerError)
				log.Default().Println(err.Error())
			}
			user.Password = pw
		}
		if email != "" {
			user.Email = email
		}
		if role != "" {
			user.Role = models.ToRole(role)
		}
	} else {
		pw, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Failed to generate password.", http.StatusInternalServerError)
			log.Default().Println(err.Error())
		}
		user = &models.User{
			ID:        0,
			Username:  username,
			Password:  pw,
			Email:     email,
			CreatedAt: time.Now(),
			Role:      models.ToRole(role),
		}
	}

	if !models.ValidateUser(user) {
		http.Error(w, "User values are not correctly defined.", http.StatusBadRequest)
		log.Default().Println("User did not pass validation.")
		return nil
	}

	return user
}

func (s *Server) RenderUserEditRow(w http.ResponseWriter, req *http.Request, user *models.User) {
	t, err := template.New("users-table-row-edit").Parse(`
		<tr hx-target="closest tr" hx-swap="outerHTML">
			<td>
				<input type="hidden" id="ID" name="ID" value="{{.ID}}" form="admin-users-edit-{{.ID}}"/>
				<input type="text" placeholder="Username" name="username" id="username"
					class="input input-bordered w-full max-w-xs" value="{{.Username}}" autofocus form="admin-users-edit-{{.ID}}"/>
			</td>
			<td>
				<input type="text" placeholder="Email" name="email" id="email"
					class="input input-bordered w-full max-w-xs" value="{{.Email}}" autofocus form="admin-users-edit-{{.ID}}"/>
			</td>
			<td>
				<select id="role" name="role" class="select select-bordered" form="admin-users-edit-{{.ID}}">
					<option value="normal" {{if eq .Role "normal"}}selected{{end}}>Normal</option>
					<option value="admin" {{if eq .Role "admin"}}selected{{end}}>Admin</option>
				</select>
			</td>
			<td>
				<span>{{.CreatedAt.Format "2006-01-02 15:04:05"}}</span>
			</td>
			<td>
				<button class="btn btn-outline btn-xs btn-success" form="admin-users-edit-{{.ID}}">Save</button>
				<button class="btn btn-outline btn-xs btn-error" hx-get="user/{{.ID}}?row">Discard</button>
			</td>
			<form hx-put="user/{{.ID}}" id="admin-users-edit-{{.ID}}"></form>
		</tr>
	`)

	err = t.Execute(w, user)
	if err != nil {
		http.Error(w, "[error] failed to generate the edit user row", http.StatusInternalServerError)
	}
}

func (s *Server) RenderUserRow(w http.ResponseWriter, req *http.Request, user *models.User) {
	t, err := template.New("users-table-row").Parse(`
		<tr hx-target="closest tr" hx-swap="outerHTML">
			<td><span>{{.Username}}</span></td>
			<td><span>{{.Email}}</span></td>
			<td><span>{{.Role}}</span></td>
			<td><span>{{.CreatedAt.Format "2006-01-02 15:04:05"}}</span></td>
			<td>
           		<div class="flex gap-4">
					<button class="btn btn-outline btn-ghost btn-xs" hx-get="admin?view=users&edit={{.ID}}"
						hx-target="closest tr">Edit</button>
					<button class="btn btn-outline btn-error btn-xs" hx-delete="user/{{.ID}}" hx-target="closest tr">Delete</button>
				</div>
			</td>
		</tr>
	`)
	err = t.Execute(w, user)
	if err != nil {
		http.Error(w, "[error] failed to generate the new user row", http.StatusInternalServerError)
	}
}
