package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

// just for theory create struct that satifies ServeHTTP and pass it as the Handler

func (hb *hootBox) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// Use the template.ParseFiles() function to read the template(html)
	//file into a template set

	files := []string{
		"./ui/html/base.tmpl",
		"./ui/html/partials/nav.tmpl",
		"./ui/html/pages/home.tmpl",
	}
	ts, err := template.ParseFiles(files...)
	if err != nil {
		hb.serverErr(w, err)
		return
	}
	// Use the ExecuteTemplate() method to write the content of the "base"
	// template as the response body.
	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		hb.serverErr(w, err)
	}

}

func (hb *hootBox) viewHoot(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		hb.notFoundErr(w)
		return
	}
	fmt.Fprintf(w, "snippet with id: %v", id)
}

func (hb *hootBox) createHoot(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		hb.clientErr(w, http.StatusMethodNotAllowed)
		return
	}
	fmt.Fprintf(w, "New snip created")
}

// func directoryHandler(w http.ResponseWriter, r *http.Request) {
// 	path := filepath.Join("ui", "static", r.URL.Path)

// 	fi, err := os.Stat(path)
// 	if fi.IsDir() && err == nil {
// 		path = filepath.Join(path, "index.html")
// 	}
// 	http.ServeFile(w, r, path)
// }
