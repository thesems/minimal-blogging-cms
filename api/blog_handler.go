package api

import (
	"encoding/json"
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

	hxReq := req.Header.Get("Hx-Request")
	// hxCurrUrl := req.Header.Get("Hx-Current-Url")

	// POST on user/create
	if tokens[2] == "create" {
		if req.Method == http.MethodPost {
			if hxReq == "true" {
				post := s.ParseCreatePost(w, req)
				s.CreatePostRow(w, req, post)
			} else {
				s.CreatePost(w, req)
			}

			return
		} else {
			s.HandleErrorPage(w, req, http.StatusMethodNotAllowed)
			return
		}
	}

	// GET, PUT, DELETE on blog/{postId}
	postId, err := strconv.Atoi(tokens[2])
	isNum := err == nil

	if req.Method == http.MethodGet {
		if hxReq == "true" {
			if err != nil {
				s.HandleErrorPage(w, req, http.StatusNotFound)
				return
			}
			post, err := s.store.GetPost(postId)
			if err != nil {
				http.Error(w, "Could not find post with such ID", http.StatusBadRequest)
				return
			}
			s.CreatePostRow(w, req, post)
		} else {
			if !isNum {
				post, err := s.store.GetPostBy(map[string]string{"urltitle": tokens[2]})
				if err != nil {
					http.Error(w, "Could not find post with such title", http.StatusBadRequest)
					return
				}
				postId = post.ID
				s.GetPostPage(w, req, postId)
			} else {
				// JSON representation of user
			}
			return
		}
	}

	if err != nil {
		s.HandleErrorPage(w, req, http.StatusNotFound)
		return
	}
	if req.Method == http.MethodPut {
		post := s.ParsePutPost(w, req)
		if hxReq == "true" {
			s.CreatePostRow(w, req, post)
		}

	} else if req.Method == http.MethodDelete {
		s.store.DeletePost(postId)
	} else {
		s.HandleErrorPage(w, req, http.StatusMethodNotAllowed)
	}
}

func (s *Server) GetPostPage(w http.ResponseWriter, req *http.Request, postId int) {
	blogPost, err := s.store.GetPost(postId)
	if err != nil {
		s.HandleErrorPage(w, req, http.StatusNotFound)
		return
	}

	data := struct {
		Header      types.Header
		Post        *models.BlogPost
		ContentHtml template.HTML
	}{
		Header: types.Header{
			Navigation: s.BuildNavigationItems(req),
			User:       "",
		},
		Post:        blogPost,
		ContentHtml: template.HTML(blogPost.Content),
	}

	w.Header().Add("Content-Type", "text/html")
	s.renderTemplate(w, req, "blog-post", data)
}

func (s *Server) ParseCreatePost(w http.ResponseWriter, req *http.Request) *models.BlogPost {
	err := req.ParseForm()
	if err != nil {
		http.Error(w, "Internal server error.", http.StatusInternalServerError)
		log.Default().Println("[error] failed to parse form data.")
		return nil
	}

	title := req.Form.Get("title")
	content := req.Form.Get("content")

	post := &models.BlogPost{
		Title: title, Content: content, CreatedAt: time.Now(),
	}

	valid := models.ValidateBlogPost(post)
	if !valid {
		http.Error(w, "Not a valid post.", http.StatusBadRequest)
		log.Default().Println("[error] post is not valid.")
		return nil
	}

	postId := s.store.CreatePost(post)
	if postId == -1 {
		return nil
	}
	post.ID = postId

	log.Default().Printf("Blog post %s created.\n", title)
	return post
}

func (s *Server) ParsePutPost(w http.ResponseWriter, req *http.Request) *models.BlogPost {
	err := req.ParseForm()
	if err != nil {
		http.Error(w, "Internal server error.", http.StatusInternalServerError)
		log.Default().Println("[error] failed to parse form data.")
		return nil
	}

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

	post, err := s.store.GetPost(postId)
	title := req.Form.Get("title")
	if title == "" {
		http.Error(w, "Internal server error.", http.StatusInternalServerError)
		log.Default().Println("[error] failed to parse title from form data.")
		return nil
	}
	post.Title = title
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

func (s *Server) CreatePostRow(w http.ResponseWriter, req *http.Request, post *models.BlogPost) {
	t, err := template.New("posts-table-row").Parse(`
		<tr hx-target="closest tr" hx-swap="outerHTML">
			<td>
				<span>
					{{.Title}}
				</span>
			</td>
			<td>
				<span>{{.CreatedAt.Format "2006-01-02 15:04:05"}}</span>
			</td>
			<td>
				<button class="btn btn-outline btn-ghost btn-xs">
					<a href="blog/{{.ID}}">View</a>
				</button>
				<button class="btn btn-outline btn-ghost btn-xs" hx-get="admin?edit={{.ID}}"
					hx-target="closest tr">Edit</button>
				<button class="btn btn-outline btn-error btn-xs" hx-delete="blog/{{.ID}}"
					hx-target="closest tr">Delete</button>
			</td>
		</tr>
	`)
	err = t.Execute(w, post)
	if err != nil {
		http.Error(w, "[error] failed to generate the new post row", http.StatusInternalServerError)
	}
}

func (s *Server) CreatePostRowEdit(w http.ResponseWriter, req *http.Request, post *models.BlogPost) {
	t, err := template.New("posts-table-row-edit").Parse(`
		<tr hx-target="closest tr" hx-swap="outerHTML">
			<td>
				<input type="hidden" name="ID" value="{{.ID}}" form="admin-posts-edit-{{.ID}}"/>
				<input type="text" placeholder="Title" name="title" id="title"
					class="input input-bordered w-full max-w-xs" value="{{.Title}}" autofocus form="admin-posts-edit-{{.ID}}"/>
			</td>
			<td>
				<span>{{.CreatedAt.Format "2006-01-02 15:04:05"}}</span>
			</td>
			<td>
				<button class="btn btn-outline btn-xs btn-success" form="admin-posts-edit-{{.ID}}">Save</button>
				<button class="btn btn-outline btn-xs btn-error" hx-get="blog/{{.ID}}?row">Discard</button>
			</td>
			<form hx-put="blog/{{.ID}}" id="admin-posts-edit-{{.ID}}"></form>
		</tr>
	`)

	err = t.Execute(w, post)
	if err != nil {
		http.Error(w, "[error] failed to generate the edit post row", http.StatusInternalServerError)
	}
}

func (s *Server) GetPost(w http.ResponseWriter, req *http.Request, postId int) {
	post, err := s.store.GetPost(postId)
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
	s.store.DeletePost(postId)
	w.WriteHeader(http.StatusOK)
}
