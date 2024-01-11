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

// Handler r viewing hoot with particular ID
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

// SignUpForm Handler
func (hb *hootBox) signUp(w http.ResponseWriter, r *http.Request) {
	data := hb.newTemplateData(r)
	data.Form = userSignupForm{}
	hb.render(w, http.StatusOK, "signup.tmpl", data)
}

// Post signUp data
func (hb *hootBox) postSignUp(w http.ResponseWriter, r *http.Request) {
	var form userSignupForm

	err := hb.decodePostForm(r, &form)
	if err != nil {
		hb.clientErr(w, http.StatusBadRequest)
		return
	}

	//validate form contents
	form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.ValidateEmail(form.Email), "email", "This field must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "This field must be 8 characters long")

	//check for erros whn validating form fields
	//if errors such as incorrect user input, render sign up form again
	if !form.Valid() {
		data := hb.newTemplateData(r)
		data.Form = form
		hb.render(w, http.StatusUnprocessableEntity, "signup.tmpl", data)
		return
	}

	//create new user in dB
	err = hb.users.Insert(form.Name, form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrExsistingCrednetials) {
			form.AddFieldError("email", "Email has already been taken.")

			data := hb.newTemplateData(r)
			data.Form = form
			hb.render(w, http.StatusUnprocessableEntity, "signup.tmpl", data)
		} else {
			hb.serverErr(w, err)
		}
		return
	}
	//flash account creation succesfull
	hb.sessionManager.Put(r.Context(), "flash", "Account Succesfully Created. Proceed to Login")
	//redirect to login page
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

// Login form Handler
func (hb *hootBox) logIn(w http.ResponseWriter, r *http.Request) {
	data := hb.newTemplateData(r)
	data.Form = Login{}
	hb.render(w, http.StatusOK, "login.tmpl", data)
}

// Post Login data
func (hb *hootBox) postLogIn(w http.ResponseWriter, r *http.Request) {
	//use form to decode login data
	var form Login

	err := hb.decodePostForm(r, &form)
	if err != nil {
		hb.clientErr(w, http.StatusBadRequest)
	}

	//validate all form inut
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.ValidateEmail(form.Email), "email", "This field must contain a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")

	if !form.Valid() {
		data := hb.newTemplateData(r)
		data.Form = form
		hb.render(w, http.StatusUnprocessableEntity, "login.tmpl", data)
		return
	}

	//check if user credentails match
	id, err := hb.users.Authenticate(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddNonFieldError("Email or psssword is incorrect")

			data := hb.newTemplateData(r)
			data.Form = form
			hb.render(w, http.StatusUnprocessableEntity, "login.tmpl", data)
		} else {
			hb.serverErr(w, err)
		}
		return
	}

	//generate new sessdion ID for the User after login (good pratice)
	hb.sessionManager.Put(r.Context(), "authenticatedUSerID", id)

	//rdirect to new page
	http.Redirect(w, r, "/hoot/create", http.StatusSeeOther)

}

// Logout
func (hb *hootBox) logOut(w http.ResponseWriter, r *http.Request) {
	//renew session token
	err := hb.sessionManager.RenewToken(r.Context())
	if err != nil {
		hb.serverErr(w, err)
		return
	}

	//remove sesion ID for loged out user
	hb.sessionManager.Remove(r.Context(), "authenticatedUserID")

	//flash notifcation to let user to log out successfully
	hb.sessionManager.Put(r.Context(), "flash", "logOut Successfully")
}
