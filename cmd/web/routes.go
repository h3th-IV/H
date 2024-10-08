package main

import (
	"net/http"
	"path/filepath"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

// check if request is a directory
func (nfs neuteredFileSystem) Open(path string) (http.File, error) {
	httpFile, err := nfs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	//get file info
	neuter, _ := httpFile.Stat()
	//check if file is a directory
	if neuter.IsDir() {
		//if path is directory join index.html to path
		index := filepath.Join(path, "index.html")
		if _, err := nfs.fs.Open(index); err != nil {
			closeErr := httpFile.Close()
			if closeErr != nil {
				return nil, closeErr
			}
			return nil, err
		}
	}
	return httpFile, nil
}

func (hb *hootBox) routes() *mux.Router {
	//create new router
	router := mux.NewRouter()

	//serve static files
	fileserver := http.FileServer(neuteredFileSystem{fs: http.Dir("./ui/static/")})
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fileserver))
	//router.Handle("/static/*filepath", fileserver).Methods(http.MethodGet)

	//create a dynamic middlware chain to monitor sessions for specific routes
	dynamicMWchain := alice.New(hb.sessionManager.LoadAndSave, noCSRF, hb.authUser)

	//use the dynamic chain for the routes
	router.Handle("/", dynamicMWchain.ThenFunc(hb.home)).Methods(http.MethodGet)
	router.Handle("/hoot/view/{id:[0-9]+}", dynamicMWchain.ThenFunc(hb.viewHoot)).Methods(http.MethodGet)
	router.Handle("/user/signup", dynamicMWchain.ThenFunc(hb.signUp)).Methods(http.MethodGet)
	router.Handle("/user/signup", dynamicMWchain.ThenFunc(hb.postSignUp)).Methods(http.MethodPost)
	router.Handle("/user/login", dynamicMWchain.ThenFunc(hb.logIn)).Methods(http.MethodGet)
	router.Handle("/user/login", dynamicMWchain.ThenFunc(hb.postLogIn)).Methods(http.MethodPost)

	//protected middleware for routes that require auth
	protectedMWchain := dynamicMWchain.Append(hb.requireAuth)
	router.Handle("/hoot/create", protectedMWchain.ThenFunc(hb.createHoot)).Methods(http.MethodGet)
	router.Handle("/hoot/create", protectedMWchain.ThenFunc(hb.postHoot)).Methods(http.MethodPost)
	router.Handle("/user/logout", protectedMWchain.ThenFunc(hb.logOut)).Methods(http.MethodPost)

	//set the middlewares chained togeter with alice
	middlewareChain := alice.New(hb.recoverPanic, hb.requestLogger, secureHaedersMW)
	router.Use(middlewareChain.Then)
	return router
}
