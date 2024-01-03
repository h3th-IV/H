package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type neuteredFileSystem struct {
	fs http.FileSystem
}

type hootBox struct {
	infolog *log.Logger
	errlog  *log.Logger
}

func (nfs neuteredFileSystem) Open(path string) (http.File, error) {
	httpFile, err := nfs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	//get file info
	neuter, _ := httpFile.Stat()
	//check if file is a directory
	if neuter.IsDir() {
		//if path is directory join index.html to path
		index := filepath.Join(path, "index.html")
		if _, err := nfs.fs.Open(index); err != nil {
			closeErr := httpFile.Close()
			if closeErr != nil {
				return nil, closeErr
			}
			return nil, err
		}
	}
	return httpFile, nil
}
func main() {
	//logger for wrting informational message
	InfoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	//logger for writing error related messages
	ErrorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	owl := &hootBox{
		infolog: InfoLog,
		errlog:  ErrorLog,
	}
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
