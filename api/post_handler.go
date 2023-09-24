package api

import (
	"fmt"
	"html/template"
	"lifeofsems-go/models"
	"lifeofsems-go/types"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
)

func (s *Server) PostGet(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
	err := req.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse post form.", http.StatusBadRequest)
		log.Default().Println(err.Error())
		return
	}

	postParam := params.ByName("post")
	postId, err := strconv.Atoi(postParam)
	isTitle := err != nil

	var post *models.Post

	if isTitle {
		post, err = s.appEnv.Posts.GetBy(map[string]string{"urltitle": postParam})
		if err != nil {
			http.Error(w, "Could not find post with such title", http.StatusBadRequest)
			log.Default().Println(err.Error())
			return
		}
	} else {
		post, err = s.appEnv.Posts.Get(postId)
		if err != nil {
			http.Error(w, "Could not find post.", http.StatusBadRequest)
			log.Default().Println(err.Error())
			return
		}
	}

	isEditPost := req.URL.Query().Has("edit")
	if isEditPost {
		s.GetPostEditPage(w, req, post)
		return
	}

	hxReq := req.Header.Get("Hx-Request") == "true"
	isEditRow := req.URL.Query().Has("row")

	if hxReq && isEditRow {
		s.CreatePostRow(w, req, post)
		return
	}

	s.GetPostPage(w, req, post.ID)
}

func (s *Server) PostPost(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
	err := req.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse post form.", http.StatusBadRequest)
		log.Default().Println(err.Error())
		return
	}

	hxReq := req.Header.Get("Hx-Request") == "true"
	if !hxReq {
		http.Error(w, "Could not decode request intention.", http.StatusBadRequest)
		return
	}

	title := req.Form.Get("title")
	description := req.Form.Get("description")
	urlTitle := req.Form.Get("url")

	post := &models.Post{
		Title: title, Content: "to-do", CreatedAt: time.Now(), ShortDescription: description,
		Url: urlTitle, Draft: true,
	}

	valid := models.ValidatePost(post)
	if !valid {
		http.Error(w, "Not a valid post.", http.StatusBadRequest)
		log.Default().Println("[error] post is not valid.")
		return
	}

	postId, err := s.appEnv.Posts.Create(post)
	if postId == -1 {
		http.Error(w, "Failed to create the post.", http.StatusInternalServerError)
		log.Default().Println(err.Error())
		return
	}
	post.ID = postId

	log.Default().Printf("Blog post %s created.\n", title)
	s.CreatePostRow(w, req, post)
}

func (s *Server) PostPut(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
	postId, err := strconv.Atoi(params.ByName("postId"))
	if err != nil {
		http.Error(w, "Internal server error.", http.StatusInternalServerError)
		log.Default().Println(err.Error())
		return
	}

	post, err := s.appEnv.Posts.Get(postId)
	if err != nil {
		log.Default().Println(err.Error())
		http.Error(w, "Post not found.", http.StatusBadRequest)
		return
	}

	err = req.ParseForm()
	if err != nil {
		log.Default().Println(err.Error())
		http.Error(w, "Failed to parse form.", http.StatusBadRequest)
		return
	}

	title := req.Form.Get("title")
	description := req.Form.Get("description")
	url := req.Form.Get("url")
	content := req.Form.Get("content")
	draft := req.Form.Get("draft")

	if title == "" && description == "" && url == "" && content == "" && draft == "" {
		http.Error(w, "Internal server error.", http.StatusInternalServerError)
		log.Default().Println("[error] failed to parse form data. Form data empty.")
		return
	}

	if title != "" {
		post.Title = title
	}
	if description != "" {
		post.ShortDescription = description
	}
	if url != "" {
		post.Url = url
	}
	if content != "" {
		post.Content = content
	}
	if draft != "" {
		post.Draft = draft == "true"
	}

	attrs := map[string]string{
		"title":            post.Title,
		"shortdescription": post.ShortDescription,
		"urltitle":         post.Url,
		"content":          post.Content,
		"draft":            strconv.FormatBool(post.Draft),
	}
	err = s.appEnv.Posts.Update(post.ID, attrs)
	if err != nil {
		http.Error(w, "Internal server error.", http.StatusInternalServerError)
		log.Default().Println("[error] failed to update the post. Error: ", err.Error())
		return
	}

	log.Default().Printf("Post %s saved.\n", post.Url)
	if post == nil {
		http.Error(w, "Failed to parse put request.", http.StatusBadRequest)
		log.Default().Println(err.Error())
		return
	}
	hxReq := req.Header.Get("Hx-Request") == "true"
	isEditPage := req.URL.Query().Has("edit")

	if !hxReq {
		http.Error(w, "Could not decode request intention.", http.StatusBadRequest)
		return
	}

	if isEditPage {
		w.Header().Set("HX-Redirect", fmt.Sprintf("/post/%s", post.Url))
		return
	}

	s.CreatePostRow(w, req, post)
}

func (s *Server) PostDelete(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
	postId, err := strconv.Atoi(params.ByName("postId"))
	if err != nil {
		return
	}
	s.appEnv.Posts.Delete(postId)
}

func (s *Server) GetPostPage(w http.ResponseWriter, req *http.Request, postId int) {
	post, err := s.appEnv.Posts.Get(postId)
	if err != nil {
		s.HandleErrorPage(w, req, http.StatusNotFound)
		return
	}

	user := GetUser(s.appEnv, w, req)
	admin := false
	if user != nil {
		admin = user.Role == models.Admin
	}

	data := struct {
		Header      types.Header
		Post        *models.Post
		ContentHtml template.HTML
		Admin       bool
		Meta        []types.Meta
	}{
		Header: types.Header{
			Navigation: BuildNavigationItems(s.appEnv, w, req),
			User:       "",
		},
		Post:        post,
		ContentHtml: template.HTML(post.Content),
		Admin:       admin,
		Meta:        []types.Meta{},
	}

	w.Header().Add("Content-Type", "text/html")
	renderTemplate(s.appEnv, w, "post", data)
}

func (s *Server) GetPostEditPage(w http.ResponseWriter, req *http.Request, post *models.Post) {
	data := struct {
		Header      types.Header
		Post        *models.Post
		ContentHtml template.HTML
		Meta        []types.Meta
	}{
		Header: types.Header{
			Navigation: BuildNavigationItems(s.appEnv, w, req),
			User:       "",
		},
		Post:        post,
		ContentHtml: template.HTML(post.Content),
		Meta:        []types.Meta{},
	}

	w.Header().Add("Content-Type", "text/html")
	renderTemplate(s.appEnv, w, "post-edit", data)
}

func (s *Server) CreatePostRow(w http.ResponseWriter, req *http.Request, post *models.Post) {
	t, err := template.New("posts-table-row").Parse(`
		<tr hx-target="closest tr" hx-swap="outerHTML">
			<td><span>{{if eq .Draft true}}Yes{{else}}No{{end}}</span></td>
			<td><span>{{.Title}}</span></td>
			<td><span>{{.ShortDescription}}</span></td>
			<td><span>{{.Url}}</span></td>
			<td><span>{{.CreatedAt.Format "2006-01-02 15:04:05"}}</span></td>
			<td>
           		<div class="flex gap-4">
					<button class="btn btn-outline btn-ghost btn-xs">
						<a href="post/{{.Url}}">View</a>
					</button>
					<button class="btn btn-outline btn-ghost btn-xs" hx-get="admin?edit={{.ID}}"
						hx-target="closest tr">Edit</button>
					<button class="btn btn-outline btn-error btn-xs" hx-delete="post/{{.ID}}"
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
					class="input input-bordered w-full max-w-xs" value="{{.Url}}" autofocus form="admin-posts-edit-{{.ID}}"/>
			</td>
			<td><span>{{.CreatedAt.Format "2006-01-02 15:04:05"}}</span></td>
			<td>
				<button class="btn btn-outline btn-xs btn-success" form="admin-posts-edit-{{.ID}}">Save</button>
				<button class="btn btn-outline btn-xs btn-error" hx-get="post/{{.ID}}?row">Discard</button>
			</td>
			<form hx-put="post/{{.ID}}" id="admin-posts-edit-{{.ID}}"></form>
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
