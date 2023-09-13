package main

import (
	"flag"
	"html/template"
	"lifeofsems-go/api"
	"lifeofsems-go/env"
	"log"
	"path/filepath"
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
	flag.Parse()

	connUrl := "postgresql://postgres:postgres@localhost:5432/postgres?sslmode=disable"
	appEnv := env.New(connUrl, "postgres")

	server := api.NewServer(*listenAddr, *appEnv, templates)
	err := server.Start()
	if err != nil {
		log.Fatalln("HTTP server failed with", err)
	}
}
