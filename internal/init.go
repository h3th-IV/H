package internal

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type database struct {
	dB *sql.DB
}

func Init() (*database, error) {
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", os.Getenv("MY_USER"), os.Getenv("MY_PASSWORD"), os.Getenv("MY_HOST"), os.Getenv("MY_PORT"), os.Getenv("MY_DBNAME"))

	dB, err := sql.Open("mysql", connStr)
	if err != nil {
		return nil, fmt.Errorf("err creating db pool: %v", err)
	}
	err = dB.Ping()
	if err != nil {
		return nil, fmt.Errorf("err connecting to db: %v", err)
	}

	hootDB := database{
		dB: dB,
	}
	return &hootDB, nil
}

func (db *database) CloseDB(error) {
	if db.dB != nil {
		err := db.dB.Close()
		if err != nil {
			log.Printf("Error closing Database Connection")
		}
	}
	fmt.Println("Database Closed Succesfully")
}
