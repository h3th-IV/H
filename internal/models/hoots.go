package models

import (
	"database/sql"
	"time"
)

type Hoot struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type HootsModel struct {
	DB *sql.DB
}

// insert new chat
func (hm *HootsModel) Insert(title, content string, expires int) (int, error) {
	return 0, nil
}

// get chat by ID
func (hm *HootsModel) Get(id int) (*Hoot, error) {
	return nil, nil
}

func (hm *HootsModel) Latest() (*[]Hoot, error) {
	return nil, nil
}
