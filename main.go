package main

import (
	"flag"
	"html/template"
	"lifeofsems-go/api"
	"lifeofsems-go/models"
	"lifeofsems-go/storage"
	"log"
	"path/filepath"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var templates map[string]*template.Template

func init() {
	if templates == nil {
		templates = make(map[string]*template.Template)
	}

	layoutTemplates, err := filepath.Glob("templates/layouts/*.gohtml")
	if err != nil {
		log.Fatal(err)
	}
	includeTemplates, err := filepath.Glob("templates/pages/*.gohtml")
	if err != nil {
		log.Fatal(err)
	}

	mainTmpl := `{{define "main" }} {{ template "base" . }} {{ end }}`
	mainTemplate := template.New("main")
	mainTemplate, err = mainTemplate.Parse(mainTmpl)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range includeTemplates {
		fileName := filepath.Base(file)
		files := append(layoutTemplates, file)
		templates[fileName], err = mainTemplate.Clone()
		if err != nil {
			log.Fatal(err)
		}
		templates[fileName] = template.Must(templates[fileName].ParseFiles(files...))
	}
}

func main() {
	listenAddr := flag.String("listenaddr", "49999", "HTTP listen port.")
	storeType := flag.String("storage", "mongo", "Storage type: mongo, memory.")
	flag.Parse()

	var store interface{}

	switch *storeType {
	case "mongo":
		store = storage.NewMongoStorage()
	case "memory":
		bs, err := bcrypt.GenerateFromPassword([]byte("123"), bcrypt.DefaultCost)
		if err != nil {
			log.Fatalln("Internal server error during password encryption.")
		}

		store = storage.NewMemoryStorage(
			[]*models.User{
				models.NewUser("admin", bs, "admin@fe.com", models.Admin),
				models.NewUser("user", bs, "user@fe.com", models.Normal),
			},
			map[int]*models.BlogPost{
				1: {ID: 1, Title: "Post 1", Content: "Content 1", CreatedAt: time.Now()},
				2: {ID: 2, Title: "Post 2", Content: "Content 2", CreatedAt: time.Now()},
			},
		)
	default:
		log.Fatal("Store type not found.")
	}

	server := api.NewServer(*listenAddr, store.(storage.Storage), templates)
	err := server.Start()
	if err != nil {
		log.Fatalln("HTTP server failed with", err)
	}
}
