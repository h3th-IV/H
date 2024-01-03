package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (hb *hootBox) routes() *mux.Router {
	//create new router
	router := mux.NewRouter()
	//serve static files
	fileserver := http.FileServer(neuteredFileSystem{fs: http.Dir("./ui/static/")})
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fileserver))

	router.HandleFunc("/", hb.home)
	router.HandleFunc("/H/create", hb.createHoot)
	router.HandleFunc("/H/view", hb.viewHoot)

	return router
}
