package api

import (
	"encoding/json"
	"fmt"
	"lifeofsems-go/models"
	"lifeofsems-go/types"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func (s *Server) HandleBlogPage(w http.ResponseWriter, req *http.Request) {
	tokens := strings.Split(req.URL.Path, "/")

	if len(tokens) < 3 {
		s.HandleErrorPage(w, req, http.StatusNotFound)
		return
	}

	// POST on blog/create
	if tokens[2] == "create" {
		if req.Method == http.MethodPost {
			s.CreatePost(w, req)
			return
		} else {
			s.HandleErrorPage(w, req, http.StatusMethodNotAllowed)
			return
		}
	}

	// GET, PUT, DELETE on blog/{postId}
	postId, err := strconv.Atoi(tokens[2])
	if err != nil {
		s.HandleErrorPage(w, req, http.StatusNotFound)
		return
	}

	if req.Method == http.MethodGet {
		s.GetPostPage(w, req, postId)
	} else if req.Method == http.MethodPut {
		fmt.Println("Method put on blog/{:d}")
	} else if req.Method == http.MethodDelete {
		fmt.Println("Method delete on blog/{:d}")
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
		Header types.Header
		Post   *models.BlogPost
	}{
		Header: types.Header{
			Navigation: s.BuildNavigationItems(req),
			User:       "",
		},
		Post: blogPost,
	}

	w.Header().Add("Content-Type", "text/html")
	s.renderTemplate(w, req, "blog-post", data)
}

func (s *Server) CreatePost(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		http.Error(w, "Internal server error.", http.StatusInternalServerError)
		log.Default().Println("[error] failed to parse form data.")
		return
	}

	title := req.Form.Get("title")
	content := req.Form.Get("content")

	post := s.store.CreatePost(&models.BlogPost{
		Title: title, Content: content,
	})

	valid := models.ValidateBlogPost(post)
	if !valid {
		http.Error(w, "Not a valid post.", http.StatusBadRequest)
		log.Default().Println("[error] post is not valid.")
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
