package main

import (
	"net/http"
	"path/filepath"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

type neuteredFileSystem struct {
	fs http.FileSystem
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

func (hb *hootBox) routes() *mux.Router {
	//create new router
	router := mux.NewRouter()

	//serve static files
	fileserver := http.FileServer(neuteredFileSystem{fs: http.Dir("./ui/static/")})
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fileserver))

	router.HandleFunc("/", hb.home)
	router.HandleFunc("/{path:[Hh]}/create", hb.createHoot).Methods(http.MethodGet)
	router.HandleFunc("/{path:[Hh]}/create", hb.postHoot).Methods(http.MethodPost)
	router.HandleFunc("/{path:[Hh]}/view/{id:[0-9]+}", hb.viewHoot).Methods(http.MethodGet)

	//set the middlewares chained togeter with alice
	middlewareChain := alice.New(hb.recoverPanic, hb.requestLogger, secureHaedersMW)
	router.Use(middlewareChain.Then)
	return router
}
