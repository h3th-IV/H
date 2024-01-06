package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
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

	data := hb.newTemplateData(r)
	data.Hoots = hoots
	hb.render(w, http.StatusOK, "home.tmpl", data)

}

func (hb *hootBox) viewHoot(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	path := strings.ToLower(vars["path"])
	if path == "h" {
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

		data := hb.newTemplateData(r)
		data.Hoot = hoot

		hb.render(w, http.StatusOK, "view.tmpl", data)
	} else {
		hb.notFoundErr(w)
	}

}

func (hb *hootBox) createHoot(w http.ResponseWriter, r *http.Request) {
	//match all H/h routes with  path variable
	vars := mux.Vars(r)
	path := strings.ToLower(vars["path"])
	if path == "h" {
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
	} else {
		hb.notFoundErr(w)
	}

}

// func directoryHandler(w http.ResponseWriter, r *http.Request) {
// 	path := filepath.Join("ui", "static", r.URL.Path)

// 	fi, err := os.Stat(path)
// 	if fi.IsDir() && err == nil {
// 		path = filepath.Join(path, "index.html")
// 	}
// 	http.ServeFile(w, r, path)
// }
