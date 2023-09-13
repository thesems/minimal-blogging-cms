package env

import (
	"database/sql"
	"fmt"
	"lifeofsems-go/models"
	"log"

	_ "github.com/lib/pq"
)

type Env struct {
	Posts    *models.PostModel
	Users    *models.UserModel
	Sessions *models.SessionModel
}

func initDB(connUrl string, driver string) (*sql.DB, error) {
	db, err := sql.Open(driver, connUrl)
	if err != nil {
		return nil, err
	}

	fmt.Printf("SQL %s storage initialized.\n", driver)

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func New(connUrl string, driver string) *Env {
	db, err := initDB(connUrl, driver)
	if err != nil {
		log.Fatalln("Database setup failed: Error:", err.Error())
		return nil
	}
	return &Env{
		Posts:    &models.PostModel{DB: db},
		Users:    &models.UserModel{DB: db},
		Sessions: &models.SessionModel{DB: db},
	}
}
