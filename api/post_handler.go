package api

import (
	"encoding/json"
	"fmt"
	"html/template"
	"lifeofsems-go/models"
	"lifeofsems-go/types"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func (s *Server) HandleBlogPage(w http.ResponseWriter, req *http.Request) {
	tokens := strings.Split(req.URL.Path, "/")
	if len(tokens) < 3 {
		s.HandleErrorPage(w, req, http.StatusNotFound)
		return
	}

	hxReq := req.Header.Get("Hx-Request") == "true"
	// hxCurrUrl := req.Header.Get("Hx-Current-Url")
	err := req.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse URL.", http.StatusBadRequest)
		log.Default().Println(err.Error())
		return
	}

	// POST on user/create
	if tokens[2] == "create" {
		if req.Method == http.MethodPost {
			if hxReq == true {
				post := s.ParseCreatePost(w, req)
				s.CreatePostRow(w, req, post)
			} else {
				s.CreatePost(w, req)
			}
		}
		return
	}

	// GET, PUT, DELETE on blog/{postId}
	postId, err := strconv.Atoi(tokens[2])
	isNum := err == nil

	edit := req.Form.Get("edit")

	if req.Method == http.MethodGet {
		var post *models.Post
		if !isNum {
			post, err = s.appEnv.Posts.GetBy(map[string]string{"urltitle": tokens[2]})
			if err != nil {
				http.Error(w, "Could not find post with such title", http.StatusBadRequest)
				log.Default().Println(err.Error())
				return
			}
		} else {
			post, err = s.appEnv.Posts.Get(postId)
			if err != nil {
				http.Error(w, "Could not find post with such ID", http.StatusBadRequest)
				log.Default().Println(err.Error())
				return
			}
		}

		if edit != "" {
			s.GetPostEditPage(w, req, post)
			return
		}

		row := req.Form.Get("row") == "1"

		if hxReq == true && row {
			s.CreatePostRow(w, req, post)
		} else {
			s.GetPostPage(w, req, post.ID)
		}
		return
	}

	if err != nil {
		s.HandleErrorPage(w, req, http.StatusNotFound)
		return
	}
	if req.Method == http.MethodPut {
		post := s.ParsePutPost(w, req)
		if post == nil {
			http.Error(w, "Failed to parse put request.", http.StatusBadRequest)
			log.Default().Println(err.Error())
			return
		}
		if hxReq == true && edit == "" {
			s.CreatePostRow(w, req, post)
		} else if hxReq == true && edit != "" {
			w.Header().Set("HX-Redirect", fmt.Sprintf("/blog/%s", post.UrlTitle))
			w.WriteHeader(http.StatusOK)
			// http.Redirect(w, req, fmt.Sprintf("/blog/%s", post.UrlTitle), http.StatusSeeOther)
		}
	} else if req.Method == http.MethodDelete {
		s.appEnv.Posts.Delete(postId)
	} else {
		s.HandleErrorPage(w, req, http.StatusMethodNotAllowed)
	}
}

func (s *Server) GetPostPage(w http.ResponseWriter, req *http.Request, postId int) {
	post, err := s.appEnv.Posts.Get(postId)
	if err != nil {
		s.HandleErrorPage(w, req, http.StatusNotFound)
		return
	}

	user := s.GetUser(w, req)
	admin := false
	if user != nil {
		admin = user.Role == models.Admin
	}

	data := struct {
		Header      types.Header
		Post        *models.Post
		ContentHtml template.HTML
		Admin       bool
	}{
		Header: types.Header{
			Navigation: s.BuildNavigationItems(w, req),
			User:       "",
		},
		Post:        post,
		ContentHtml: template.HTML(post.Content),
		Admin:       admin,
	}

	w.Header().Add("Content-Type", "text/html")
	s.renderTemplate(w, req, "blog-post", data)
}

func (s *Server) GetPostEditPage(w http.ResponseWriter, req *http.Request, post *models.Post) {
	data := struct {
		Header      types.Header
		Post        *models.Post
		ContentHtml template.HTML
	}{
		Header: types.Header{
			Navigation: s.BuildNavigationItems(w, req),
			User:       "",
		},
		Post:        post,
		ContentHtml: template.HTML(post.Content),
	}

	w.Header().Add("Content-Type", "text/html")
	s.renderTemplate(w, req, "blog-post-edit", data)
}

func (s *Server) ParseCreatePost(w http.ResponseWriter, req *http.Request) *models.Post {
	title := req.Form.Get("title")
	description := req.Form.Get("description")
	urlTitle := req.Form.Get("url")

	post := &models.Post{
		Title: title, Content: "to-do", CreatedAt: time.Now(), ShortDescription: description,
		UrlTitle: urlTitle, Draft: true,
	}

	valid := models.ValidatePost(post)
	if !valid {
		http.Error(w, "Not a valid post.", http.StatusBadRequest)
		log.Default().Println("[error] post is not valid.")
		return nil
	}

	postId := s.appEnv.Posts.Create(post)
	if postId == -1 {
		return nil
	}
	post.ID = postId

	log.Default().Printf("Blog post %s created.\n", title)
	return post
}

func (s *Server) ParsePutPost(w http.ResponseWriter, req *http.Request) *models.Post {
	postIdStr := req.Form.Get("ID")
	if postIdStr == "" {
		http.Error(w, "Internal server error.", http.StatusInternalServerError)
		log.Default().Println("[error] failed to find ID from form data.")
		return nil
	}
	postId, err := strconv.Atoi(postIdStr)
	if err != nil {
		http.Error(w, "Internal server error.", http.StatusInternalServerError)
		log.Default().Println("[error] failed to parse ID from form data.")
		return nil
	}

	post, err := s.appEnv.Posts.Get(postId)
	title := req.Form.Get("title")
	description := req.Form.Get("description")
	url := req.Form.Get("url")
	content := req.Form.Get("content")
	if title == "" && description == "" && url == "" && content == "" {
		http.Error(w, "Internal server error.", http.StatusInternalServerError)
		log.Default().Println("[error] failed to parse form data. Form data empty.")
		return nil
	}

	if title != "" {
		post.Title = title
	}
	if description != "" {
		post.ShortDescription = description
	}
	if url != "" {
		post.UrlTitle = url
	}
	if content != "" {
		post.Content = content
	}

	attrs := map[string]string{
		"title":            post.Title,
		"shortdescription": post.ShortDescription,
		"urltitle":         post.UrlTitle,
		"content":          post.Content,
	}
	err = s.appEnv.Posts.Update(post.ID, attrs)
	if err != nil {
		http.Error(w, "Internal server error.", http.StatusInternalServerError)
		log.Default().Println("[error] failed to update the post. Error: ", err.Error())
		return nil
	}

	log.Default().Printf("Post %s saved.\n", post.UrlTitle)
	return post
}

func (s *Server) CreatePost(w http.ResponseWriter, req *http.Request) {
	post := s.ParseCreatePost(w, req)
	message, err := json.Marshal(post)
	if err != nil {
		http.Error(w, "Internal server error.", http.StatusInternalServerError)
		log.Default().Println("[error] failed to marshal post as json.")
		return
	}

	log.Default().Println("render post in json")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(message)
	if err != nil {
		http.Error(w, "Internal server error.", http.StatusInternalServerError)
		log.Default().Println("[error] failed to write json message to HTTP response.")
		return
	}
}

func (s *Server) CreatePostRow(w http.ResponseWriter, req *http.Request, post *models.Post) {
	t, err := template.New("posts-table-row").Parse(`
		<tr hx-target="closest tr" hx-swap="outerHTML">
			<td><span>{{if eq .Draft true}}Yes{{else}}No{{end}}</span></td>
			<td><span>{{.Title}}</span></td>
			<td><span>{{.ShortDescription}}</span></td>
			<td><span>{{.UrlTitle}}</span></td>
			<td><span>{{.CreatedAt.Format "2006-01-02 15:04:05"}}</span></td>
			<td>
           		<div class="flex gap-4">
					<button class="btn btn-outline btn-ghost btn-xs">
						<a href="blog/{{.UrlTitle}}">View</a>
					</button>
					<button class="btn btn-outline btn-ghost btn-xs" hx-get="admin?edit={{.ID}}"
						hx-target="closest tr">Edit</button>
					<button class="btn btn-outline btn-error btn-xs" hx-delete="blog/{{.ID}}"
						hx-target="closest tr">Delete</button>
				</div>
			</td>
		</tr>
	`)
	err = t.Execute(w, post)
	if err != nil {
		http.Error(w, "[error] failed to generate the new post row", http.StatusInternalServerError)
	}
}

func (s *Server) CreatePostRowEdit(w http.ResponseWriter, req *http.Request, post *models.Post) {
	t, err := template.New("posts-table-row-edit").Parse(`
		<tr hx-target="closest tr" hx-swap="outerHTML">
			<td><span>{{if eq .Draft true}}Yes{{else}}No{{end}}</span></td>
			<td>
				<input type="hidden" name="ID" value="{{.ID}}" form="admin-posts-edit-{{.ID}}"/>
				<input type="text" placeholder="Title" name="title" id="title"
					class="input input-bordered w-full max-w-xs" value="{{.Title}}" autofocus form="admin-posts-edit-{{.ID}}"/>
			</td>
			<td>
				<input type="text" placeholder="Description" name="description" id="description"
					class="input input-bordered w-full max-w-xs" value="{{.ShortDescription}}" autofocus form="admin-posts-edit-{{.ID}}"/>
			</td>
			<td>
				<input type="text" placeholder="URL" name="url" id="url"
					class="input input-bordered w-full max-w-xs" value="{{.UrlTitle}}" autofocus form="admin-posts-edit-{{.ID}}"/>
			</td>
			<td><span>{{.CreatedAt.Format "2006-01-02 15:04:05"}}</span></td>
			<td>
				<button class="btn btn-outline btn-xs btn-success" form="admin-posts-edit-{{.ID}}">Save</button>
				<button class="btn btn-outline btn-xs btn-error" hx-get="blog/{{.ID}}?row=1">Discard</button>
			</td>
			<form hx-put="blog/{{.ID}}" id="admin-posts-edit-{{.ID}}"></form>
		</tr>
	`)

	err = t.Execute(w, post)
	if err != nil {
		http.Error(w, "[error] failed to generate the edit post row", http.StatusInternalServerError)
		log.Default().Println(err.Error())
	}
}

func (s *Server) CreatePostContent(w http.ResponseWriter, req *http.Request, post *models.Post) {
	t, err := template.New("post-content").Parse(`<div>{{.Content}}</div>`)
	err = t.Execute(w, post)
	if err != nil {
		http.Error(w, "[error] failed to generate the post content", http.StatusInternalServerError)
	}
}

func (s *Server) GetPost(w http.ResponseWriter, req *http.Request, postId int) {
	post, err := s.appEnv.Posts.Get(postId)
	if err != nil {
		http.Error(w, "Post could not be found.", http.StatusNotFound)
		log.Default().Printf("Post ID %d could not be found.\n", postId)
		return
	}

	message, err := json.Marshal(post)
	if err != nil {
		http.Error(w, "Internal server error.", http.StatusInternalServerError)
		log.Default().Println("[error] failed to marshal post as json.")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(message)
	if err != nil {
		http.Error(w, "Internal server error.", http.StatusInternalServerError)
		log.Default().Println("[error] failed to write json message to HTTP response.")
		return
	}
}

func (s *Server) UpdatePost(w http.ResponseWriter, req *http.Request) {}

func (s *Server) DeletePost(w http.ResponseWriter, req *http.Request, postId int) {
	s.appEnv.Posts.Delete(postId)
	w.WriteHeader(http.StatusOK)
}
