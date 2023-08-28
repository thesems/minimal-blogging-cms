package main

import (
	"flag"
	"html/template"
	"lifeofsems-go/api"
	"lifeofsems-go/storage"
	"lifeofsems-go/types"
	"log"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("templates/*"))
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
		store = storage.NewMemoryStorage(
			[]*types.BlogPost{
				{ID: 1, Title: "Post 1", Content: "Content 1"},
			},
		)
	default:
		log.Fatal("Store type not found.")
	}

	server := api.NewServer(*listenAddr, store.(storage.Storage), tpl)
	err := server.Start()
	if err != nil {
		log.Fatalln("HTTP server failed with", err)
	}
}
