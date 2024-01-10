package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

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

}

// handler to parse (html)form for creatingHoot
func (hb *hootBox) createHoot(w http.ResponseWriter, r *http.Request) {
	data := hb.newTemplateData(r)

	//init a new create HootForm and pass it to template
	data.Form = hootCreateForm{
		Expires: 365,
	}
	hb.render(w, http.StatusOK, "create.tmpl", data)

}

type hootCreateForm struct {
	Title               string     `form:"title"`
	Content             string     `form:"content"`
	Expires             int        `form:"expires"`
	validator.Validator `form:"-"` //"-" will be ignored when decoding
}

// handler for creating new hoot
func (hb *hootBox) postHoot(w http.ResponseWriter, r *http.Request) {
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

	id, err := hb.dataBox.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		hb.serverErr(w, err)
		return
	}

	//use .Put() to store success message in the session data with the key 'flash'
	hb.sessionManager.Put(r.Context(), "flash", "Hoot succesfully created")
	//then redirect to the newly created page *wink *wink
	http.Redirect(w, r, fmt.Sprintf("/hoot/view/%d", id), http.StatusSeeOther)

}

// func directoryHandler(w http.ResponseWriter, r *http.Request) {
// 	path := filepath.Join("ui", "static", r.URL.Path)

// 	fi, err := os.Stat(path)
// 	if fi.IsDir() && err == nil {
// 		path = filepath.Join(path, "index.html")
// 	}
// 	http.ServeFile(w, r, path)
// }

type userSignupForm struct {
	Name                string `form:"name"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

// SignUpForm Handler
func (hb *hootBox) signUp(w http.ResponseWriter, r *http.Request) {
	data := hb.newTemplateData(r)
	data.Form = userSignupForm{}
	hb.render(w, http.StatusOK, "signup.tmpl", data)
}

// Post signUp data
func (hb *hootBox) postSignUp(w http.ResponseWriter, r *http.Request) {
	var form userSignupForm

	err := hb.decodePostForm(r, form)
	if err != nil {
		hb.clientErr(w, http.StatusBadRequest)
		return
	}

	//validate form contents
	form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.ValidateEmail(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "This field must be 8 characters long")

	//check for erros whn validating form fields
	//if errors such as incorrect user input, render sign up form again
	if !form.Valid() {
		data := hb.newTemplateData(r)
		data.Form = form
		hb.render(w, http.StatusUnprocessableEntity, "sign.tmpl", data)
		return
	}

	//Create user in db
	fmt.Fprintln(w, "SignUp Succesfull")
}

// Login form Handler
func (hb *hootBox) logIn(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "PLease Provide your Details to Login")
}

// Post Login data
func (hb *hootBox) postLogIn(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "LogIn Successfull")
}

// Logout
func (hb *hootBox) logOut(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "LogOut Succesfull")
}
