package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", homeHandler)
	router.HandleFunc("/hoot/create", createHoot)
	router.HandleFunc("/hoot/view", viewHoot)
	server := &http.Server{
		Addr:    ":8000",
		Handler: router,
	}
	fmt.Println("Listening and serving @ :8000")
	log.Fatal(server.ListenAndServe())
}
