package main

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"runtime/debug"
	"time"

	"github.com/h3th-IV/H/internal/models"
)

// write error err and stack trace  to the errlog attribute,
// sends 500 internal server error to user
// It returns a formatted stack trace of the current goroutine.
// A stack trace is a list of the function calls that have been
// executed up to a certain point in the program.
func (hb *hootBox) serverErr(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	hb.errlog.Output(2, trace) // track the debug stack 2 steps backward
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// sends stsatus code to user and corresponding status message
func (hb *hootBox) clientErr(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// 404-like responses  to user
func (hb *hootBox) notFoundErr(w http.ResponseWriter) {
	hb.clientErr(w, http.StatusNotFound)
}

// template data to be rendered
type templateData struct {
	Hoot        *models.Hoot
	Hoots       []*models.Hoot
	CurrentYear int
}

func (hb *hootBox) newTemplateData(r *http.Request) *templateData {
	return &templateData{
		CurrentYear: time.Now().Year(),
	}
}

// function to return humanreadable date
func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

// Initialize a template.FuncMap object and store it in a global variable. This is
// essentially a string-keyed map which acts as a lookup between the names of our
// custom template functions and the functions themselves.
var functions = template.FuncMap{
	"humanDate": humanDate,
}

// cache(store) parsed templates in an in-memory cache. instead of loading 'em everytime
func newTemplateCache() (map[string]*template.Template, error) {
	//init new map to use as cache
	cache := map[string]*template.Template{}

	//use filepath.Glob to get all template files
	pages, err := filepath.Glob("./ui/html/pages/*.tmpl")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		//return the last element of the string (which is file name e.g "home.tmpl")
		file := filepath.Base(page)

		//parse template files
		ts, err := template.New(file).Funcs(functions).ParseFiles("./ui/html/base.tmpl")
		if err != nil {
			return nil, err
		}

		// Call ParseGlob() *on this template set * to add any partials.
		ts, err = ts.ParseGlob("./ui/html/partials/*.tmpl")
		if err != nil {
			return nil, err
		}

		// Call ParseFiles() *on this template set* to add the page template.
		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		// Add the template set to the map as normal...
		cache[file] = ts
	}
	return cache, nil
}

func (hb *hootBox) render(w http.ResponseWriter, status int, page string, data *templateData) {
	//retrive the appropriate template(page) matching from the cache
	ts, ok := hb.templateCache[page]
	if !ok {
		err := fmt.Errorf("template does not exist: %v", page)
		hb.serverErr(w, err)
		return
	}

	//init a new buffer to check for errors before writing to response writer
	buffer := new(bytes.Buffer)

	// Write the template to the buffer, instead of straight to the
	// http.ResponseWriter. If there's an error, call our serverError() helper
	// and then return.
	err := ts.ExecuteTemplate(buffer, "base", data) //(w, "base", data)
	if err != nil {
		hb.serverErr(w, err)
		return
	}

	//if template is written to buffer succesfully we can write header
	//write out the Header 200 if page is found and 404 if not found
	w.WriteHeader(status)

	//then write from buffer to
	buffer.WriteTo(w)

}
