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
	pageTemplates, err := filepath.Glob("templates/pages/*.gohtml")
	if err != nil {
		log.Fatal(err)
	}
	componentTemplates, err := filepath.Glob("templates/components/*.gohtml")
	if err != nil {
		log.Fatal(err)
	}

	mainTmpl := `{{define "main" }} {{ template "base" . }} {{ end }}`
	mainTemplate := template.New("main")
	mainTemplate, err = mainTemplate.Parse(mainTmpl)
	if err != nil {
		log.Fatal(err)
	}

	allTemplates := make([]string, 0)
	allTemplates = append(allTemplates, layoutTemplates...)
	allTemplates = append(allTemplates, componentTemplates...)

	for _, file := range pageTemplates {
		fileName := filepath.Base(file)
		files := append(allTemplates, file)
		templates[fileName], err = mainTemplate.Clone()
		if err != nil {
			log.Fatal(err)
		}
		templates[fileName] = template.Must(templates[fileName].ParseFiles(files...))
	}
}

func main() {
	listenAddr := flag.String("listenaddr", "49999", "HTTP listen port.")
	storeType := flag.String("storage", "postgres", "Storage type: postgres, memory.")
	flag.Parse()

	var store interface{}

	switch *storeType {
	case "postgres":
		connStr := "postgresql://postgres:postgres@localhost:5432/postgres?sslmode=disable"
		store = storage.NewPostgresStorage(connStr, "postgres")
	case "memory":
		genPass := func(pw string) []byte {
			bs, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
			if err != nil {
				log.Fatalln("Internal server error during password encryption.")
			}
			return bs
		}

		store = storage.NewMemoryStorage(
			map[int]*models.User{
				1: {ID: 1, Username: "admin", Password: genPass("1234"), Email: "admin@fe.com", Role: models.Admin, CreatedAt: time.Now()},
				2: {ID: 2, Username: "user", Password: genPass("4321"), Email: "user@fe.com", Role: models.Normal, CreatedAt: time.Now()},
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
