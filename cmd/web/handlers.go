package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/h3th-IV/H/internal/models"
	"github.com/h3th-IV/H/internal/validator"
)

// just for theory create struct that satifies ServeHTTP and pass it as the Handler
// ghostmac#6861
func (hb *hootBox) home(w http.ResponseWriter, r *http.Request) {

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
		id, err := strconv.Atoi(vars["id"])
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

// handler to parse form for creatingHoot
func (hb *hootBox) createHoot(w http.ResponseWriter, r *http.Request) {
	//match all H/h routes with  path variable
	vars := mux.Vars(r)
	path := strings.ToLower(vars["path"])
	if path == "h" {
		data := hb.newTemplateData(r)

		//init a new create HootForm and pass it to template
		data.Form = hootCreateForm{
			Expires: 365,
		}
		hb.render(w, http.StatusOK, "create.tmpl", data)
	} else {
		hb.notFoundErr(w)
	}
}

type hootCreateForm struct {
	Title               string     `form:"title"`
	Content             string     `form:"content"`
	Expires             int        `form:"expires"`
	validator.Validator `form:"-"` //"-" will be ignored when decoding
}

// handler for creating new hoot
func (hb *hootBox) postHoot(w http.ResponseWriter, r *http.Request) {
	//match all H/h routes with  path variable
	vars := mux.Vars(r)
	path := strings.ToLower(vars["path"])
	if path == "h" {
		// Declare a new empty instance of the snippetCreateForm struct.
		var form hootCreateForm
		//limit request body size to 4069 bytes
		r.Body = http.MaxBytesReader(w, r.Body, 4069)
		err := hb.decodePostForm(r, &form)
		if err != nil {
			hb.clientErr(w, http.StatusBadRequest)
			return
		}

		form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
		form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")
		form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
		form.CheckField(validator.PermittedInt(form.Expires, 1, 7, 365), "expires", "This field must equal 1, 7 or 365")

		//use valid mathod to see if any of the baove heck failed
		if !form.Valid() {
			data := hb.newTemplateData(r)
			data.Form = form
			hb.render(w, http.StatusUnprocessableEntity, "create.tmpl", data)
			return
		}

		//if after all white space has been removed and we have empty string
		//use RuneCountInString to check number of characters returned

		id, err := hb.dataBox.Insert(form.Title, form.Content, form.Expires)
		if err != nil {
			hb.serverErr(w, err)
			return
		}
		http.Redirect(w, r, fmt.Sprintf("/H/view/%d", id), http.StatusSeeOther)
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
