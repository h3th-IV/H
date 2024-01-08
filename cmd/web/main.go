package main

import (
	"crypto/tls"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	"github.com/h3th-IV/H/internal/models"
	"github.com/joho/godotenv"
)

type hootBox struct {
	infolog        *log.Logger
	errlog         *log.Logger
	dataBox        *models.HootsModel
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
}

func main() {
	//logger for wrting informational message
	InfoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	//logger for writing error related messages
	ErrorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	err := godotenv.Load()
	if err != nil {
		ErrorLog.Printf("Error loading Environment variables: %v", err)
	}
	dB, err := Init()
	if err != nil {
		ErrorLog.Printf("%v", err)
	}
	cacheFiles, err := newTemplateCache()
	if err != nil {
		ErrorLog.Printf("Err Parsing template files")
	}

	//init new session manager, configure it to use mysql db as it store and set life span
	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(dB)
	sessionManager.Lifetime = 10 * time.Hour
	//set session cookie  secure field to true -- to serve all requests over HTTPS
	sessionManager.Cookie.Secure = true

	formDecoder := form.NewDecoder()
	owl := &hootBox{
		infolog:        InfoLog,
		errlog:         ErrorLog,
		dataBox:        &models.HootsModel{DB: dB},
		templateCache:  cacheFiles,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
	}

	defer owl.dataBox.DB.Close()

	//define a cli flag with default addr and help message
	addr := flag.String("addr", ":8000", "http network addr(port)")
	// This reads in the command-line flag value and assigns it to the addr
	// variable. You need to call this *before* you use the addr variable
	// otherwise it will always contain the default value of ":4000". If any errors are
	// encountered during parsing the application will be terminated.
	flag.Parse()

	//init a tls config struct that stores non-default tls configuration that we want to specify for our servers
	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	//server
	server := &http.Server{
		Addr:         *addr,
		Handler:      owl.routes(),
		ErrorLog:     ErrorLog,
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 4 * time.Second,
	}
	InfoLog.Printf("Listening and serving %s", *addr)
	//http.ListenAndServeTLS() -serves on TCP netaddr calls ServeTLS to handle all iincoming connectionns over TLS
	ErrorLog.Fatal(server.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem"))
}
