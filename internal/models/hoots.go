package models

import (
	"database/sql"
	"errors"
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
	//query for creating new cxat
	query := `INSERT INTO hoots (title, content, created, expires)
VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	//made use of transactions instead to ensure proper release of resources
	trx, err := hm.DB.Begin()
	if err != nil {
		return 0, err
	}
	defer trx.Rollback()

	stmt, err := trx.Prepare(query) //hm.DB.Exec(query, title, content, expires)
	if err != nil {
		return 0, err
	}
	result, err := stmt.Exec(title, content, expires)
	if err != nil {
		return 0, err
	}

	//get inserted chat
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	err = trx.Commit()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

// get chat by ID
func (hm *HootsModel) Get(id int) (*Hoot, error) {
	query := `SELECT id, title, content, created, expires FROM hoots WHERE expires > UTC_TIMESTAMP() AND id = ?`

	trx, err := hm.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer trx.Rollback()

	stmt, err := trx.Prepare(query)
	if err != nil {
		return nil, err
	}
	row := stmt.QueryRow(id) //hm.DB.QueryRow(query, id)
	//use a new hoot to pass value
	hoot := &Hoot{}
	//scan into corresonding fields into hoot struct
	err = row.Scan(&hoot.ID, &hoot.Title, &hoot.Content, &hoot.Created, &hoot.Expires)
	if err != nil {
		//check if error (err)returned is no row errors
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			//if it's another type of error then return that error err
			return nil, err
		}
	}
	err = trx.Commit()
	if err != nil {
		return nil, err
	}
	return hoot, nil
}

func (hm *HootsModel) Latest() ([]*Hoot, error) {
	query := `SELECT id, title, content, created, expires FROM hoots WHERE expires > UTC_TIMESTAMP() ORDER BY id DESC LIMIT 10`

	trx, err := hm.DB.Begin()
	if err != nil {
		return nil, err
	}

	defer trx.Rollback()

	stmt, err := trx.Prepare(query)
	if err != nil {
		return nil, err
	}
	rows, err := stmt.Query() //.DB.Query(query)
	if err != nil {
		return nil, err

	}
	defer rows.Close()

	Hoots := []*Hoot{}
	//loop through retrieved rows
	for rows.Next() {
		hoot := &Hoot{}
		//get fields for each row
		err := rows.Scan(&hoot.ID, &hoot.Title, &hoot.Content, &hoot.Created, &hoot.Expires)
		if err != nil {
			return nil, err
		}
		Hoots = append(Hoots, hoot)
	}

	//check for eny erros while looping throug the rows
	if err = rows.Err(); err != nil {
		return nil, err
	}

	err = trx.Commit()
	if err != nil {
		return nil, err
	}
	return Hoots, nil
}
