package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/h3th-IV/H/internal/models"
)

func (hb *hootBox) Init() (*hootBox, error) {
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", os.Getenv("MY_USER"), os.Getenv("MY_PASSWORD"), os.Getenv("MY_HOST"), os.Getenv("MY_PORT"), os.Getenv("MY_DBNAME"))

	dataBase, err := sql.Open("mysql", connStr)
	if err != nil {
		return nil, fmt.Errorf("err creating db pool: %v", err)
	}
	err = dataBase.Ping()
	if err != nil {
		dataBase.Close()
		return nil, fmt.Errorf("err connecting to db: %v", err)
	}
	dB := models.HootsModel{
		DB: dataBase,
	}
	hb.infolog.Printf("Connected to database Successfully\n")
	//pass the database object to hootBox
	hootDB := hootBox{
		dataBox: &dB,
	}
	return &hootDB, nil
}

func (hb *hootBox) CloseDB() error {
	//check the database existence and close it succesfully
	if hb.dataBox != nil {
		err := hb.dataBox.DB.Close()
		if err != nil {
			return err
		}
	}
	fmt.Println("Database Closed Succesfully")
	return nil
}
