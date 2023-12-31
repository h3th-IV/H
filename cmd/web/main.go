package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"github.com/gorilla/mux"
)

type neuteredFileSystem struct {
	fs http.FileSystem
}

func (nfs neuteredFileSystem) Open(path string) (http.File, error) {
	file, err := nfs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	//get file info
	neuter, err := file.Stat()
	//check if file is a directory
	if neuter.IsDir() {
		//if path is directory join index.html to path
		index := filepath.Join(path, "/../index.html")
		if _, err := nfs.fs.Open(index); err != nil {
			closeErr := file.Close()
			if closeErr != nil {
				return nil, closeErr
			}
			return nil, err
		}
	}
	return file, nil
}
func main() {
	router := mux.NewRouter()

	fileserver := http.FileServer(neuteredFileSystem{fs: http.Dir("./ui/static/")})
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fileserver))

	router.HandleFunc("/", homeHandler)
	router.HandleFunc("/H/create", createHoot)
	router.HandleFunc("/H/view", viewHoot)

	server := &http.Server{
		Addr:    ":8000",
		Handler: router,
	}
	fmt.Println("Listening and serving @ :8000")
	log.Fatal(server.ListenAndServe())
}

///static/img/logo.png
