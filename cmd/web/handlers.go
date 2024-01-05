package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/h3th-IV/H/internal/models"
)

// just for theory create struct that satifies ServeHTTP and pass it as the Handler

func (hb *hootBox) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		hb.notFoundErr(w)
		return
	}

	hoots, err := hb.dataBox.Latest()
	if err != nil {
		hb.serverErr(w, err)
		return
	}

	for _, hoot := range hoots {
		fmt.Fprintf(w, "%+v\n\n", hoot)
	}
	// Use the template.ParseFiles() function to read the template(html)
	//file into a template set

	// files := []string{
	// 	"./ui/html/base.tmpl",
	// 	"./ui/html/partials/nav.tmpl",
	// 	"./ui/html/pages/home.tmpl",
	// }
	// ts, err := template.ParseFiles(files...)
	// if err != nil {
	// 	hb.serverErr(w, err)
	// 	return
	// }
	// // Use the ExecuteTemplate() method to write the content of the "base"
	// // template as the response body.
	// err = ts.ExecuteTemplate(w, "base", nil)
	// if err != nil {
	// 	hb.serverErr(w, err)
	// }

}

func (hb *hootBox) viewHoot(w http.ResponseWriter, r *http.Request) {
	//try convert id query to int type
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		hb.notFoundErr(w)
		return
	}

	//get chat with id
	hoot, err := hb.dataBox.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			hb.notFoundErr(w)
		} else {
			hb.serverErr(w, err)
		}
		return
	}
	fmt.Fprintf(w, "%+v", hoot) //use + in %+v to include field names
}

func (hb *hootBox) createHoot(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		hb.clientErr(w, http.StatusMethodNotAllowed)
		return
	}
	title := "Time Traveler"
	content := "O Man\nTraverse these path of Infiniteness,\nBut slowly, slowly!\n\n– Drunk Man"
	expires := 7

	id, err := hb.dataBox.Insert(title, content, expires)
	if err != nil {
		hb.serverErr(w, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/H/view?id=%d", id), http.StatusSeeOther)
}

// func directoryHandler(w http.ResponseWriter, r *http.Request) {
// 	path := filepath.Join("ui", "static", r.URL.Path)

// 	fi, err := os.Stat(path)
// 	if fi.IsDir() && err == nil {
// 		path = filepath.Join(path, "index.html")
// 	}
// 	http.ServeFile(w, r, path)
// }
