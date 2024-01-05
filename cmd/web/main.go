package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/h3th-IV/H/internal/models"
	"github.com/joho/godotenv"
)

type hootBox struct {
	infolog *log.Logger
	errlog  *log.Logger
	dataBox *models.HootsModel
}

func main() {
	//logger for wrting informational message
	InfoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	//logger for writing error related messages
	ErrorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	err := godotenv.Load()
	if err != nil {
		ErrorLog.Printf("Error loading Environment variables: %v", err)
	}
	dB, err := Init()
	if err != nil {
		ErrorLog.Printf("%v", err)
	}
	owl := &hootBox{
		infolog: InfoLog,
		errlog:  ErrorLog,
		dataBox: &models.HootsModel{DB: dB},
	}

	defer owl.dataBox.DB.Close()

	//define a cli flag with default addr and help message
	addr := flag.String("addr", ":8000", "http network addr(port)")
	// This reads in the command-line flag value and assigns it to the addr
	// variable. You need to call this *before* you use the addr variable
	// otherwise it will always contain the default value of ":4000". If any errors are
	// encountered during parsing the application will be terminated.
	flag.Parse()

	server := &http.Server{
		Addr:     *addr,
		Handler:  owl.routes(),
		ErrorLog: ErrorLog,
	}
	InfoLog.Printf("Listening and serving %s", *addr)
	ErrorLog.Fatal(server.ListenAndServe())
}
