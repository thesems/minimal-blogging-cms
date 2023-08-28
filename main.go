package main

import (
	"flag"
	"html/template"
	"lifeofsems-go/api"
	"lifeofsems-go/storage"
	"lifeofsems-go/types"
	"log"

	"golang.org/x/crypto/bcrypt"
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
		bs, err := bcrypt.GenerateFromPassword([]byte("123"), bcrypt.DefaultCost)
		if err != nil {
			log.Fatalln("Internal server error during password encryption.")
		}

		store = storage.NewMemoryStorage(
			[]*types.User{
				types.NewUser("admin", bs, "admin@fe.com", types.Admin),
				types.NewUser("user", bs, "user@fe.com", types.Normal),
			},
			[]*types.BlogPost{
				{ID: 1, Title: "Post 1", Content: "Content 1"},
				{ID: 2, Title: "Post 2", Content: "Content 2"},
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
