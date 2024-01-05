package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func Init() (*sql.DB, error) {
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

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	infoLog.Printf("connected to dataBase succefully!")
	//pass the database object to hootBox
	return dataBase, nil
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
