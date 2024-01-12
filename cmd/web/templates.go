package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	"github.com/h3th-IV/H/internal/models"
	"github.com/h3th-IV/H/internal/validator"
)

// hootBox Application
type hootBox struct {
	infolog        *log.Logger
	errlog         *log.Logger
	dataBox        *models.HootsModel
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
	users          *models.UserModels
}

// directory checker
type neuteredFileSystem struct {
	fs http.FileSystem
}

// template data to be rendered
type templateData struct {
	Hoot            *models.Hoot
	Hoots           []*models.Hoot
	CurrentYear     int
	Form            any
	Flash           string
	IsAuthenticated bool
	XsRfToken       string
}

// new hoot form
type hootCreateForm struct {
	Title               string     `form:"title"`
	Content             string     `form:"content"`
	Expires             int        `form:"expires"`
	validator.Validator `form:"-"` //"-" will be ignored when decoding
}

// signUp fprm
type userSignupForm struct {
	Name                string `form:"name"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	UserDB              *models.UserModels
	validator.Validator `form:"-"`
}

// login form
type LoginForm struct {
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}
