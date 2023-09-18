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
	fmt.Println("PostGet")
	postTitle := params.ByName("post")
	post, err := s.appEnv.Posts.GetBy(map[string]string{"urltitle": postTitle})
	if err != nil {
		http.Error(w, "Could not find post with such title", http.StatusBadRequest)
		log.Default().Println(err.Error())
		return
	}

	edit := req.Form.Get("edit")
	if edit != "" {
		s.GetPostEditPage(w, req, post)
		return
	}

	hxReq := req.Header.Get("Hx-Request") == "true"
	row := req.Form.Get("row") == "1"
	if hxReq == true && row {
		s.CreatePostRow(w, req, post)
	} else {
		s.GetPostPage(w, req, post.ID)
	}
}

func (s *Server) PostPost(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
	fmt.Println("PostPost")
	err := req.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse post form.", http.StatusBadRequest)
		log.Default().Println("[error] failed to parse post form.")
		return
	}

	hxReq := req.Header.Get("Hx-Request") == "true"
	if hxReq == true {
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
}

func (s *Server) PostPut(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
	fmt.Println("PostPut")

	postId, err := strconv.Atoi(params.ByName("ID"))
	if err != nil {
		http.Error(w, "Internal server error.", http.StatusInternalServerError)
		log.Default().Println("[error] failed to parse ID from form data.")
		return
	}

	post, err := s.appEnv.Posts.Get(postId)
	title := req.Form.Get("title")
	description := req.Form.Get("description")
	url := req.Form.Get("url")
	content := req.Form.Get("content")
	if title == "" && description == "" && url == "" && content == "" {
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
		return
	}

	log.Default().Printf("Post %s saved.\n", post.UrlTitle)
	if post == nil {
		http.Error(w, "Failed to parse put request.", http.StatusBadRequest)
		log.Default().Println(err.Error())
		return
	}
	hxReq := req.Header.Get("Hx-Request") == "true"
	edit := req.Form.Get("edit")
	if hxReq == true && edit == "" {
		s.CreatePostRow(w, req, post)
	} else if hxReq == true && edit != "" {

		draft := req.Form.Get("draft")
		if draft == "false" {
			setAttrs := map[string]string{
				"Draft": "false",
			}
			err := s.appEnv.Posts.Update(post.ID, setAttrs)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				log.Default().Println(err.Error())
			}
		} else {
			w.Header().Set("HX-Redirect", fmt.Sprintf("/post/%s", post.UrlTitle))
		}
	}
}

func (s *Server) PostDelete(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
	postId, err := strconv.Atoi(params.ByName("post"))
	if err != nil {
		return
	}
	fmt.Println("post delete hit,", postId)
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
			<td><span>{{.UrlTitle}}</span></td>
			<td><span>{{.CreatedAt.Format "2006-01-02 15:04:05"}}</span></td>
			<td>
           		<div class="flex gap-4">
					<button class="btn btn-outline btn-ghost btn-xs">
						<a href="post/{{.UrlTitle}}">View</a>
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
					class="input input-bordered w-full max-w-xs" value="{{.UrlTitle}}" autofocus form="admin-posts-edit-{{.ID}}"/>
			</td>
			<td><span>{{.CreatedAt.Format "2006-01-02 15:04:05"}}</span></td>
			<td>
				<button class="btn btn-outline btn-xs btn-success" form="admin-posts-edit-{{.ID}}">Save</button>
				<button class="btn btn-outline btn-xs btn-error" hx-get="post/{{.ID}}?row=1">Discard</button>
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
