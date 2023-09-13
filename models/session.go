package models

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"
)

type Session struct {
	ID           string    `json:"id"`
	Username     string    `json:"username"`
	LastActivity time.Time `json:"lastActivity"`
}

type SessionModel struct {
	DB *sql.DB
}

func (m *SessionModel) Get(session_id string) (*Session, error) {
	row := m.DB.QueryRow("SELECT * FROM cms.session WHERE id=$1", session_id)
	var session Session
	err := row.Scan(&session.ID, &session.Username, &session.LastActivity)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("could not find session %s", session_id))
	}
	return &session, nil
}

func (m *SessionModel) All() ([]*Session, error) {
	rows, err := m.DB.Query("SELECT * FROM cms.session")
	if err != nil {
		return nil, err
	}
	sessions := make([]*Session, 0)
	for rows.Next() {
		var session Session
		err := rows.Scan(&session.ID, &session.Username, &session.LastActivity)
		if err != nil {
			return nil, err
		}
	}
	return sessions, nil
}

func (m *SessionModel) Create(session_id string, username string) {
	lastActivity := time.Now()
	_, err := m.DB.Exec("INSERT INTO cms.session (id,username,lastactivity) VALUES ($1,$2,$3)", session_id, username, lastActivity)
	if err != nil {
		log.Default().Fatalln(err.Error())
	}
}

func (m *SessionModel) Delete(session_id string) {
	_, err := m.DB.Exec("DELETE FROM cms.session WHERE id=$1", session_id)
	if err != nil {
		log.Default().Fatalln(err.Error())
	}
}
